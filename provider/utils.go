package provider

import "net"

// getInterfaceIP는 인터페이스 이름으로부터 IP 주소를 가져옵니다
// BSD와 Windows에서 공통으로 사용됩니다
func getInterfaceIP(ifname string) string {
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		return ""
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			if ipv4 := ipnet.IP.To4(); ipv4 != nil {
				return ipv4.String()
			}
		}
	}

	return ""
}

