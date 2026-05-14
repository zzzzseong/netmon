package utils

import (
	"context"
	"fmt"
	"net"
	"time"
)

const (
	// DefaultDNSTimeout is the default timeout for DNS lookups
	DefaultDNSTimeout = 5 * time.Second
)

// DNSResult contains DNS lookup results.
type DNSResult struct {
	Query        string        // The query (domain or IP)
	ARecords     []string      // IPv4 addresses
	AAAARecords  []string      // IPv6 addresses
	PTRRecords   []string      // Reverse DNS records
	MXRecords    []string      // Mail exchange records
	NSRecords    []string      // Name server records
	TXTRecords   []string      // TXT records
	CNAMERecord  string        // CNAME record
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

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), DefaultDNSTimeout)
	defer cancel()

	// Create resolver with timeout
	resolver := &net.Resolver{
		PreferGo: true,
	}

	// Lookup IP addresses (both IPv4 and IPv6) with timeout
	ips, err := resolver.LookupIP(ctx, "ip", domain)
	result.ResponseTime = time.Since(startTime)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			result.Error = fmt.Errorf("DNS lookup timeout: %w", err)
		} else {
			result.Error = err
		}
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

	// Lookup CNAME record
	cname, err := resolver.LookupCNAME(ctx, domain)
	if err == nil && cname != domain+"." && cname != "" {
		result.CNAMERecord = cname
	}

	// Lookup MX records
	mxRecords, err := resolver.LookupMX(ctx, domain)
	if err == nil {
		for _, mx := range mxRecords {
			result.MXRecords = append(result.MXRecords, fmt.Sprintf("%s (priority: %d)", mx.Host, mx.Pref))
		}
	}

	// Lookup NS records
	nsRecords, err := resolver.LookupNS(ctx, domain)
	if err == nil {
		for _, ns := range nsRecords {
			result.NSRecords = append(result.NSRecords, ns.Host)
		}
	}

	// Lookup TXT records
	txtRecords, err := resolver.LookupTXT(ctx, domain)
	if err == nil {
		result.TXTRecords = txtRecords
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

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), DefaultDNSTimeout)
	defer cancel()

	// Create resolver with timeout
	resolver := &net.Resolver{
		PreferGo: true,
	}

	// Perform reverse lookup with timeout
	names, err := resolver.LookupAddr(ctx, ipAddr)
	result.ResponseTime = time.Since(startTime)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			result.Error = fmt.Errorf("DNS reverse lookup timeout: %w", err)
		} else {
			result.Error = err
		}
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

	// Check for invalid characters
	for _, char := range domain {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '.' || char == '-' || char == '_') {
			return fmt.Errorf("domain contains invalid character: %c", char)
		}
	}

	// Check that domain doesn't start or end with dot or hyphen
	if domain[0] == '.' || domain[0] == '-' || domain[len(domain)-1] == '.' || domain[len(domain)-1] == '-' {
		return fmt.Errorf("domain cannot start or end with dot or hyphen")
	}

	// Check for consecutive dots
	for i := 0; i < len(domain)-1; i++ {
		if domain[i] == '.' && domain[i+1] == '.' {
			return fmt.Errorf("domain cannot contain consecutive dots")
		}
	}

	return nil
}

// ValidateHostname validates a hostname or IP address for use in network commands.
// Returns an error if the hostname is invalid.
func ValidateHostname(hostname string) error {
	if hostname == "" {
		return fmt.Errorf("hostname cannot be empty")
	}

	// Check if it's a valid IP address
	if IsIPAddress(hostname) {
		return nil
	}

	// Validate as domain name
	return ValidateDomain(hostname)
}
