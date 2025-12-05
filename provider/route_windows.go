//go:build windows

package provider

import (
	"fmt"
	"net"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modiphlpapi           = windows.NewLazySystemDLL("iphlpapi.dll")
	procGetIpForwardTable2 = modiphlpapi.NewProc("GetIpForwardTable2")
	procFreeMibTable      = modiphlpapi.NewProc("FreeMibTable")
)

// WindowsRouteProvider is a RouteProvider implementation using Windows GetIpForwardTable2 API.
type WindowsRouteProvider struct{}

// NewWindowsRouteProvider creates a new WindowsRouteProvider instance.
func NewWindowsRouteProvider() *WindowsRouteProvider {
	return &WindowsRouteProvider{}
}

// NewRouteProvider returns a RouteProvider implementation for the current OS.
// On Windows, it returns a WindowsRouteProvider.
func NewRouteProvider() RouteProvider {
	return NewWindowsRouteProvider()
}

// MIB_IPFORWARD_ROW2 구조체 (간소화 버전)
type mibIpforwardRow2 struct {
	InterfaceLuid      uint64
	InterfaceIndex     uint32
	DestinationPrefix  ipAddressPrefix
	NextHop            sockaddrInet
	SitePrefixLength   byte
	ValidLifetime      uint32
	PreferredLifetime  uint32
	Metric             uint32
	Protocol           uint32
	Loopback           byte
	AutoconfigureAddress byte
	Publish            byte
	Immortal           byte
	Age                uint32
	Origin             uint32
}

type ipAddressPrefix struct {
	Prefix       sockaddrInet
	PrefixLength uint8
	_            [2]byte // padding
}

type sockaddrInet struct {
	Family uint16
	_      [2]byte // padding
	Data   [26]byte
}

// MIB_IPFORWARD_TABLE2 구조체
type mibIpforwardTable2 struct {
	NumEntries uint32
	_          [4]byte // padding for alignment
	Table      [1]mibIpforwardRow2
}

// GetRoutes retrieves the routing table using Windows GetIpForwardTable2 API.
// Returns a slice of RouteEntry and an error if the operation fails.
func (p *WindowsRouteProvider) GetRoutes() ([]RouteEntry, error) {
	var table *mibIpforwardTable2

	// GetIpForwardTable2 호출 (AF_INET = IPv4)
	ret, _, _ := procGetIpForwardTable2.Call(
		uintptr(windows.AF_INET),
		uintptr(unsafe.Pointer(&table)),
	)

	if ret != 0 {
		return nil, fmt.Errorf("GetIpForwardTable2 failed with error code: %d", ret)
	}

	if table == nil {
		return nil, fmt.Errorf("GetIpForwardTable2 returned null table")
	}

	// 사용 후 메모리 해제
	defer procFreeMibTable.Call(uintptr(unsafe.Pointer(table)))

	// 테이블 항목 파싱
	numEntries := table.NumEntries
	entries := make([]RouteEntry, 0, numEntries)

	// table.Table은 [1]이지만 실제로는 NumEntries만큼의 배열
	tablePtr := uintptr(unsafe.Pointer(&table.Table[0]))
	rowSize := unsafe.Sizeof(mibIpforwardRow2{})

	for i := uint32(0); i < numEntries; i++ {
		row := (*mibIpforwardRow2)(unsafe.Pointer(tablePtr + uintptr(i)*rowSize))

		entry := RouteEntry{}

		// Destination 처리
		destIP := parseSockaddrInet(&row.DestinationPrefix.Prefix)
		prefixLen := row.DestinationPrefix.PrefixLength

		if destIP != nil {
			if prefixLen == 0 && destIP.Equal(net.IPv4zero) {
				entry.Destination = "default"
			} else {
				entry.Destination = fmt.Sprintf("%s/%d", destIP.String(), prefixLen)
			}
		}

		// Gateway 처리
		gwIP := parseSockaddrInet(&row.NextHop)
		if gwIP != nil && !gwIP.Equal(net.IPv4zero) {
			entry.Gateway = gwIP.String()
		}

		// Interface 처리
		if row.InterfaceIndex != 0 {
			iface, err := net.InterfaceByIndex(int(row.InterfaceIndex))
			if err == nil {
				entry.Interface = iface.Name
			} else {
				// 실패하면 인덱스 번호로 표시
				entry.Interface = fmt.Sprintf("if%d", row.InterfaceIndex)
			}
		}

		// Metric 처리
		entry.Metric = int(row.Metric)

		// Source IP는 Windows API에서 직접 제공하지 않으므로
		// 인터페이스에서 가져옴
		if entry.Interface != "" {
			if srcIP := getInterfaceIP(entry.Interface); srcIP != "" {
				entry.Source = srcIP
			}
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

// parseSockaddrInet은 sockaddrInet을 net.IP로 변환합니다
func parseSockaddrInet(sa *sockaddrInet) net.IP {
	if sa.Family == windows.AF_INET {
		// IPv4: Data[2:6]에 주소가 있음
		return net.IPv4(sa.Data[2], sa.Data[3], sa.Data[4], sa.Data[5])
	}
	return nil
}

