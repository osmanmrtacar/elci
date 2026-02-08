package services

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// privateRanges contains CIDR ranges that should be blocked for SSRF prevention.
var privateRanges = []string{
	"127.0.0.0/8",    // Loopback
	"10.0.0.0/8",     // RFC 1918
	"172.16.0.0/12",  // RFC 1918
	"192.168.0.0/16", // RFC 1918
	"169.254.0.0/16", // Link-local / cloud metadata
	"::1/128",        // IPv6 loopback
	"fc00::/7",       // IPv6 unique local
	"fe80::/10",      // IPv6 link-local
}

var parsedPrivateRanges []*net.IPNet

func init() {
	for _, cidr := range privateRanges {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(fmt.Sprintf("invalid CIDR in privateRanges: %s", cidr))
		}
		parsedPrivateRanges = append(parsedPrivateRanges, network)
	}
}

// ValidateMediaURL checks that a URL is safe to fetch (not targeting internal services).
func ValidateMediaURL(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("empty URL")
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Only allow http and https
	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("unsupported URL scheme: %s (only http and https are allowed)", parsed.Scheme)
	}

	host := parsed.Hostname()
	if host == "" {
		return fmt.Errorf("URL has no host")
	}

	// Block localhost variants
	lowerHost := strings.ToLower(host)
	if lowerHost == "localhost" || lowerHost == "metadata.google.internal" {
		return fmt.Errorf("URL host is not allowed: %s", host)
	}

	// Resolve the hostname to IPs and check each one
	ips, err := net.LookupIP(host)
	if err != nil {
		return fmt.Errorf("cannot resolve host %s: %w", host, err)
	}

	for _, ip := range ips {
		for _, network := range parsedPrivateRanges {
			if network.Contains(ip) {
				return fmt.Errorf("URL resolves to a private/internal address")
			}
		}
	}

	return nil
}
