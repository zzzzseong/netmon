package parser

import (
	"bufio"
	"regexp"
	"strconv"

	"netmon/formatter"
)

// TracerouteParser parses traceroute output from different operating systems.
type TracerouteParser struct {
	formatter *formatter.TracerouteFormatter
}

// NewTracerouteParser creates a new TracerouteParser instance.
func NewTracerouteParser() *TracerouteParser {
	return &TracerouteParser{
		formatter: formatter.NewTracerouteFormatter(),
	}
}

// ParseUnixTraceroute parses Unix traceroute output from the scanner.
// It reads lines from the scanner and formats each hop in real-time.
func (p *TracerouteParser) ParseUnixTraceroute(scanner *bufio.Scanner) {
	// 첫 줄 건너뛰기 (헤더)
	if scanner.Scan() {
		// traceroute to ... 줄
	}

	hopRegex := regexp.MustCompile(`^\s*(\d+)\s+(.*)$`)
	ipRegex := regexp.MustCompile(`(\d+\.\d+\.\d+\.\d+|(?:[0-9a-fA-F]{0,4}:){2,}[0-9a-fA-F]{0,4})`)
	rttRegex := regexp.MustCompile(`(\d+\.\d+)\s*ms`)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		matches := hopRegex.FindStringSubmatch(line)
		if len(matches) < 3 {
			continue
		}

		hopNum, _ := strconv.Atoi(matches[1])
		rest := matches[2]

		hop := formatter.TraceHop{
			Hop:    hopNum,
			RTT1:   "*",
			RTT2:   "*",
			RTT3:   "*",
			Status: "timeout",
		}

		// IP 주소 추출
		if ipMatch := ipRegex.FindString(rest); ipMatch != "" {
			hop.IP = ipMatch
			hop.Host = ipMatch
			hop.Status = "success"
		}

		// RTT 추출
		rttMatches := rttRegex.FindAllStringSubmatch(rest, -1)
		for i, match := range rttMatches {
			if i >= 3 {
				break
			}
			rttValue := match[1] + " ms"
			switch i {
			case 0:
				hop.RTT1 = rttValue
			case 1:
				hop.RTT2 = rttValue
			case 2:
				hop.RTT3 = rttValue
			}
		}

		// 실시간으로 hop 출력
		p.formatter.PrintHopLine(hop)
	}
}

// ParseWindowsTracert parses Windows tracert output from the scanner.
// It reads lines from the scanner and formats each hop in real-time.
func (p *TracerouteParser) ParseWindowsTracert(scanner *bufio.Scanner) {
	hopRegex := regexp.MustCompile(`^\s*(\d+)\s+(.*)$`)
	ipRegex := regexp.MustCompile(`(\d+\.\d+\.\d+\.\d+|(?:[0-9a-fA-F]{0,4}:){2,}[0-9a-fA-F]{0,4})`)
	rttRegex := regexp.MustCompile(`(<?\d+)\s*ms`)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		matches := hopRegex.FindStringSubmatch(line)
		if len(matches) < 3 {
			continue
		}

		hopNum, _ := strconv.Atoi(matches[1])
		rest := matches[2]

		hop := formatter.TraceHop{
			Hop:    hopNum,
			RTT1:   "*",
			RTT2:   "*",
			RTT3:   "*",
			Status: "timeout",
		}

		if ipMatch := ipRegex.FindString(rest); ipMatch != "" {
			hop.IP = ipMatch
			hop.Host = ipMatch
			hop.Status = "success"
		}

		rttMatches := rttRegex.FindAllStringSubmatch(rest, -1)
		for i, match := range rttMatches {
			if i >= 3 {
				break
			}
			rttValue := match[1] + " ms"
			switch i {
			case 0:
				hop.RTT1 = rttValue
			case 1:
				hop.RTT2 = rttValue
			case 2:
				hop.RTT3 = rttValue
			}
		}

		// 실시간으로 hop 출력
		p.formatter.PrintHopLine(hop)
	}
}
