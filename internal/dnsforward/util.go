package dnsforward

import (
	"net"
	"sort"
	"strings"

	"github.com/AdguardTeam/AdGuardHome/internal/aghnet"
)

// IPFromAddr gets IP address from addr.
func IPFromAddr(addr net.Addr) (ip net.IP) {
	switch addr := addr.(type) {
	case *net.UDPAddr:
		return addr.IP
	case *net.TCPAddr:
		return addr.IP
	}
	return nil
}

// IPStringFromAddr extracts IP address from net.Addr.
// Note: we can't use net.SplitHostPort(a.String()) because of IPv6 zone:
// https://github.com/AdguardTeam/AdGuardHome/internal/issues/1261
func IPStringFromAddr(addr net.Addr) (ipStr string) {
	if ip := IPFromAddr(addr); ip != nil {
		return ip.String()
	}

	return ""
}

// Find value in a sorted array
func findSorted(ar []string, val string) int {
	i := sort.SearchStrings(ar, val)
	if i == len(ar) || ar[i] != val {
		return -1
	}
	return i
}

func isWildcard(host string) bool {
	return len(host) >= 2 &&
		host[0] == '*' && host[1] == '.'
}

// Return TRUE if host name matches a wildcard pattern
func matchDomainWildcard(host, wildcard string) bool {
	return isWildcard(wildcard) &&
		strings.HasSuffix(host, wildcard[1:])
}

// Return TRUE if client's SNI value matches DNS names from certificate
func matchDNSName(dnsNames []string, sni string) bool {
	if aghnet.ValidateDomainName(sni) != nil {
		return false
	}

	if findSorted(dnsNames, sni) != -1 {
		return true
	}

	for _, dn := range dnsNames {
		if matchDomainWildcard(sni, dn) {
			return true
		}
	}
	return false
}
