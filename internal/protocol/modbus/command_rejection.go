package modbus

import "errors"

// IsCommandRejected only accepts a typed, parsed illegal-function response.
// CRC, timeout, malformed or short frames never indicate unsupported commands.
func IsCommandRejected(err error) bool {
	var protocolError *ProtocolError
	return errors.As(err, &protocolError) && protocolError.Code() == ExceptionIllegalFunction
}
