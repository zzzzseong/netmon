package provider

import (
	"bufio"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// FallbackRouteProvider is a RouteProvider implementation that uses command-line parsing.
// It serves as a fallback for unsupported operating systems.
type FallbackRouteProvider struct{}

// NewFallbackRouteProvider creates a new FallbackRouteProvider instance.
func NewFallbackRouteProvider() *FallbackRouteProvider {
	return &FallbackRouteProvider{}
}

// GetRoutes retrieves the routing table by executing system commands and parsing their output.
// Returns a slice of RouteEntry and an error if the operation fails.
func (p *FallbackRouteProvider) GetRoutes() ([]RouteEntry, error) {
	var cmd *exec.Cmd
	var parser func(string) []RouteEntry

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("ip", "route", "show")
		parser = parseLinuxRoute
	case "darwin", "freebsd", "openbsd", "netbsd":
		cmd = exec.Command("netstat", "-rn")
		parser = parseBSDRoute
	case "windows":
		cmd = exec.Command("route", "print")
		parser = parseWindowsRoute
	default:
		return nil, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute command: %w", err)
	}

	return parser(string(output)), nil
}

// parseLinuxRoute는 Linux "ip route" 출력을 파싱합니다
func parseLinuxRoute(output string) []RouteEntry {
	routes := make([]RouteEntry, 0, 32) // 일반적으로 라우트는 많지 않음
	scanner := bufio.NewScanner(strings.NewReader(output))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 1 {
			continue
		}

		route := RouteEntry{}

		// ip route show 형식 파싱
		if fields[0] == "default" {
			route.Destination = "default"
		} else if strings.Contains(fields[0], "/") || strings.Contains(fields[0], ".") {
			route.Destination = fields[0]
		}

		for i := 1; i < len(fields)-1; i++ {
			switch fields[i] {
			case "via":
				route.Gateway = fields[i+1]
				i++
			case "dev":
				route.Interface = fields[i+1]
				i++
			case "metric":
				if metric, err := strconv.Atoi(fields[i+1]); err == nil {
					route.Metric = metric
				}
				i++
			case "src":
				route.Source = fields[i+1]
				i++
			}
		}

		if route.Destination != "" {
			routes = append(routes, route)
		}
	}

	return routes
}

// parseBSDRoute는 BSD "netstat -rn" 출력을 파싱합니다
func parseBSDRoute(output string) []RouteEntry {
	routes := make([]RouteEntry, 0, 32) // 일반적으로 라우트는 많지 않음
	scanner := bufio.NewScanner(strings.NewReader(output))

	inIPv4Section := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// IPv4 섹션 시작 확인
		if strings.Contains(line, "Internet:") || strings.Contains(line, "Routing tables") {
			inIPv4Section = true
			continue
		}

		// IPv6 섹션 시작 시 종료
		if strings.Contains(line, "Internet6:") {
			break
		}

		// 헤더 라인 건너뛰기
		if strings.Contains(line, "Destination") || strings.Contains(line, "Gateway") {
			continue
		}

		if !inIPv4Section {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		dest := fields[0]
		gateway := fields[1]
		iface := fields[len(fields)-1] // 인터페이스는 보통 마지막

		// IPv6, 루프백, 시스템 라우트 제외
		if strings.Contains(dest, ":") || strings.Contains(gateway, ":") {
			continue
		}
		if strings.HasPrefix(iface, "lo") || strings.HasPrefix(iface, "utun") || strings.HasPrefix(iface, "llw") {
			continue
		}

		route := RouteEntry{
			Destination: dest,
			Interface:   iface,
		}

		// 게이트웨이 처리
		if !strings.HasPrefix(gateway, "link#") && gateway != "localhost" {
			route.Gateway = gateway
		}

		if dest == "default" {
			route.Destination = "default"
		}

		routes = append(routes, route)
	}

	return routes
}

// parseWindowsRoute는 Windows "route print" 출력을 파싱합니다
func parseWindowsRoute(output string) []RouteEntry {
	routes := make([]RouteEntry, 0, 32) // 일반적으로 라우트는 많지 않음
	scanner := bufio.NewScanner(strings.NewReader(output))

	inIPv4Section := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// IPv4 라우팅 테이블 섹션 확인
		if strings.Contains(line, "IPv4 Route Table") || strings.Contains(line, "Active Routes") {
			inIPv4Section = true
			continue
		}

		// 다른 섹션 시작 시 종료
		if strings.Contains(line, "IPv6 Route Table") || strings.Contains(line, "Persistent Routes") {
			break
		}

		// 헤더 라인 건너뛰기
		if strings.Contains(line, "Network Destination") || strings.Contains(line, "=====") {
			continue
		}

		if !inIPv4Section {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		dest := fields[0]
		netmask := fields[1]
		gateway := fields[2]
		iface := fields[3]
		metric := 0

		if len(fields) >= 5 {
			if m, err := strconv.Atoi(fields[4]); err == nil {
				metric = m
			}
		}

		// 기본 게이트웨이 처리
		if dest == "0.0.0.0" {
			dest = "default"
		}

		// CIDR 형식으로 변환
		if dest != "default" && netmask != "" {
			dest = convertToCIDR(dest, netmask)
		}

		route := RouteEntry{
			Destination: dest,
			Gateway:     gateway,
			Interface:   iface,
			Metric:      metric,
		}

		routes = append(routes, route)
	}

	return routes
}

// convertToCIDR은 IP와 넷마스크를 CIDR 형식으로 변환합니다
func convertToCIDR(ip, netmask string) string {
	// 넷마스크를 prefix length로 변환
	parts := strings.Split(netmask, ".")
	if len(parts) != 4 {
		return ip
	}

	prefixLen := 0
	for _, part := range parts {
		if n, err := strconv.Atoi(part); err == nil {
			for n > 0 {
				if n&128 != 0 {
					prefixLen++
				}
				n <<= 1
			}
		}
	}

	return fmt.Sprintf("%s/%d", ip, prefixLen)
}

