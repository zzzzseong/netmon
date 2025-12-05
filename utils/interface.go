package utils

import (
	"fmt"
	"net"
)

// GetInterfaceIP retrieves the IPv4 address of a network interface by name.
// It returns the first IPv4 address found on the interface, or an error if
// the interface is not found or has no IPv4 addresses.
// This function is commonly used by BSD and Windows route providers.
func GetInterfaceIP(ifname string) (string, error) {
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		return "", fmt.Errorf("interface %s not found: %w", ifname, err)
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return "", fmt.Errorf("failed to get addresses for interface %s: %w", ifname, err)
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			if ipv4 := ipnet.IP.To4(); ipv4 != nil {
				return ipv4.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no IPv4 address found for interface %s", ifname)
}
