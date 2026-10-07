package api

import (
	"net"
	"net/http"
	"net/netip"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func runtimeLogPath(path, base string) bool {
	root := strings.TrimRight(base, "/") + "/system/logs"
	return path == root || strings.HasPrefix(path, root+"/")
}

func localLogGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		for key := range c.Writer.Header() {
			if strings.HasPrefix(strings.ToLower(key), "access-control-") {
				c.Writer.Header().Del(key)
			}
		}
		c.Header("Cache-Control", "no-store")
		c.Header("Cross-Origin-Resource-Policy", "same-origin")
		c.Header("X-Content-Type-Options", "nosniff")
		if !localLogAllowed(c.Request) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": "logs_local_only"})
			return
		}
		c.Next()
	}
}

func localLogAllowed(r *http.Request) bool {
	peer, peerPort, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil || !numericLoopback(peer) || !validLogPort(peerPort) {
		return false
	}
	local, ok := r.Context().Value(http.LocalAddrContextKey).(net.Addr)
	if !ok || local == nil {
		return false
	}
	localHost, port, err := net.SplitHostPort(local.String())
	if err != nil || !numericLoopback(localHost) || !validLogPort(port) {
		return false
	}
	scheme, defaultPort := "http", "80"
	if r.TLS != nil {
		scheme, defaultPort = "https", "443"
	}
	host, hostPort, err := net.SplitHostPort(r.Host)
	if err != nil {
		host, hostPort = r.Host, defaultPort
		if host != "127.0.0.1" && host != "localhost" && host != "[::1]" {
			return false
		}
		if host == "[::1]" {
			host = "::1"
		}
	} else if r.Host != net.JoinHostPort(host, hostPort) {
		return false
	}
	if host != "127.0.0.1" && host != "::1" && host != "localhost" {
		return false
	}
	if hostPort != port {
		return false
	}
	for key := range r.Header {
		lower := strings.ToLower(key)
		if lower == "forwarded" || lower == "x-real-ip" || strings.HasPrefix(lower, "x-forwarded") {
			return false
		}
	}
	origins, present := exactHeaderValues(r.Header, "Origin")
	if present && (len(origins) != 1 || origins[0] != scheme+"://"+r.Host) {
		return false
	}
	sites, present := exactHeaderValues(r.Header, "Sec-Fetch-Site")
	return !present || len(sites) == 1 && (sites[0] == "same-origin" || sites[0] == "none")
}

func validLogPort(port string) bool {
	n, err := strconv.Atoi(port)
	return err == nil && n > 0 && n <= 65535 && port == strconv.Itoa(n)
}
func numericLoopback(value string) bool {
	ip, err := netip.ParseAddr(value)
	return err == nil && ip.Zone() == "" && ip.IsLoopback()
}

// Header names are also case insensitive in requests constructed by middleware.
func exactHeaderValues(header http.Header, name string) ([]string, bool) {
	var values []string
	present := false
	for key, entries := range header {
		if strings.EqualFold(key, name) {
			present = true
			values = append(values, entries...)
		}
	}
	return values, present
}
