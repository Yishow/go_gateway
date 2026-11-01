package fatek

import "errors"

// IsCommandRejected recognizes a parsed command-48 unsupported-function status.
// Appendix 1-3 assigns that meaning to 9; 4 is ambiguous format/command failure.
// The client validates framing, checksum, command, station and error-body length
// before returning this typed error. Text and malformed responses are not proof.
func IsCommandRejected(err error) bool {
	var protocolError *FatekProtocolError
	return errors.As(err, &protocolError) && protocolError.Command == "48" && protocolError.ErrorCode == "9"
}
