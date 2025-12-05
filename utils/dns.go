package utils

import (
	"fmt"
	"net"
	"time"
)

// DNSResult contains DNS lookup results.
type DNSResult struct {
	Query        string        // The query (domain or IP)
	ARecords     []string      // IPv4 addresses
	AAAARecords  []string      // IPv6 addresses
	PTRRecords   []string      // Reverse DNS records
	ResponseTime time.Duration // Time taken for lookup
	Error        error         // Error if lookup failed
}

// LookupDomain performs a DNS lookup for the given domain.
// Returns A and AAAA records along with response time.
func LookupDomain(domain string) DNSResult {
	result := DNSResult{
		Query:       domain,
		ARecords:    []string{},
		AAAARecords: []string{},
	}

	startTime := time.Now()

	// Lookup IP addresses (both IPv4 and IPv6)
	ips, err := net.LookupIP(domain)
	result.ResponseTime = time.Since(startTime)

	if err != nil {
		result.Error = err
		return result
	}

	// Separate IPv4 and IPv6 addresses
	for _, ip := range ips {
		if ip.To4() != nil {
			// IPv4
			result.ARecords = append(result.ARecords, ip.String())
		} else {
			// IPv6
			result.AAAARecords = append(result.AAAARecords, ip.String())
		}
	}

	return result
}

// ReverseLookup performs a reverse DNS lookup for the given IP address.
// Returns PTR records along with response time.
func ReverseLookup(ipAddr string) DNSResult {
	result := DNSResult{
		Query:      ipAddr,
		PTRRecords: []string{},
	}

	startTime := time.Now()

	// Perform reverse lookup
	names, err := net.LookupAddr(ipAddr)
	result.ResponseTime = time.Since(startTime)

	if err != nil {
		result.Error = err
		return result
	}

	result.PTRRecords = names

	return result
}

// IsIPAddress checks if the given string is a valid IP address.
func IsIPAddress(s string) bool {
	return net.ParseIP(s) != nil
}

// ValidateDomain performs basic validation on a domain name.
// Returns an error if the domain is invalid.
func ValidateDomain(domain string) error {
	if domain == "" {
		return fmt.Errorf("domain cannot be empty")
	}
	if len(domain) > 253 {
		return fmt.Errorf("domain name too long (max 253 characters)")
	}
	return nil
}
