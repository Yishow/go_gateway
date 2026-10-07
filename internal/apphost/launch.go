// Package apphost defines platform-independent startup and shutdown policy.
package apphost

import (
	"errors"
	"fmt"
	"net"
	"net/netip"
	"net/url"
)

// Mode distinguishes an explicit headless override from legacy console defaults.
type Mode string

const (
	Desktop  Mode = "desktop"
	Console  Mode = "console"
	Headless Mode = "headless"
)

// Launch is parsed before any database or runtime initialization.
type Launch struct {
	Mode          Mode
	Help, Version bool
}

// Parse implements the portable command contract without terminal guessing.
func Parse(args []string, desktopDefault bool) (Launch, error) {
	out := Launch{Mode: Console}
	if desktopDefault {
		out.Mode = Desktop
	}
	for _, arg := range args {
		switch arg {
		case "--headless":
			out.Mode = Headless
		case "--help", "-h":
			out.Help = true
		case "--version":
			out.Version = true
		default:
			return Launch{}, NewFault("startup.config_invalid", errors.New("unsupported command-line option"))
		}
	}
	return out, nil
}

// Informational reports commands which must not initialize services.
func (l Launch) Informational() bool { return l.Help || l.Version }

// AutoOpen preserves legacy opt-in while explicit headless always wins.
func (l Launch) AutoOpen(value string) bool {
	if l.Mode == Headless {
		return false
	}
	if l.Mode == Desktop {
		return value != "false" && value != "0"
	}
	return value == "true" || value == "1"
}

// FaultError separates safe classification from a console-only cause.
type FaultError struct {
	Code  string
	Cause error
}

func (e *FaultError) Error() string {
	if e.Cause == nil {
		return e.Code
	}
	return fmt.Sprintf("%s: %v", e.Code, e.Cause)
}
func (e *FaultError) Unwrap() error { return e.Cause }

// NewFault classifies errors before projection into managed diagnostics.
func NewFault(code string, cause error) error { return &FaultError{Code: code, Cause: cause} }

// Code never returns arbitrary error text.
func Code(err error) string {
	var fault *FaultError
	if errors.As(err, &fault) {
		return fault.Code
	}
	return "startup.failed"
}

// IsCode compares a safe classification.
func IsCode(err error, code string) bool { return Code(err) == code }

// ListenerURL constructs fixed-route origins from verified numeric listeners.
// Local logs are available only through the Host authorities allowed by their guard.
func ListenerURL(address string) (baseURL string, logsAvailable bool, err error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return "", false, err
	}
	ip, parseErr := netip.ParseAddr(host)
	if parseErr != nil {
		return "", false, errors.New("listener is not numeric")
	}
	if ip.IsUnspecified() {
		if ip.Is6() {
			host = "::1"
		} else {
			host = "127.0.0.1"
		}
	} else {
		host = ip.Unmap().String()
	}
	return (&url.URL{Scheme: "http", Host: net.JoinHostPort(host, port)}).String(), host == "127.0.0.1" || host == "::1", nil
}
