package fatek

import "errors"

// Control Characters
const (
	STX = 0x02
	ETX = 0x03
)

// Default Configuration
const (
	DefaultStation = 1
	DefaultTCPPort = 500
)

// Component Types (Symbol)
const (
	SymbolX   = "X"
	SymbolY   = "Y"
	SymbolM   = "M"
	SymbolS   = "S"
	SymbolT   = "T"
	SymbolC   = "C"
	SymbolR   = "R"
	SymbolD   = "D"
	SymbolRT  = "RT"
	SymbolRC  = "RC"
	SymbolDR  = "DR" // 32-bit D
	SymbolDW  = "DW" // 32-bit Discrete
	SymbolDWM = "DWM"
)

// Errors
var (
	ErrInvalidSTX      = errors.New("invalid STX")
	ErrInvalidETX      = errors.New("invalid ETX")
	ErrChecksumMismatch = errors.New("checksum mismatch")
	ErrResponseTooShort = errors.New("response too short")
	ErrInvalidCommand   = errors.New("unexpected command in response")
	ErrTimeout          = errors.New("timeout waiting for response")
	ErrConnectionClosed = errors.New("connection closed")
)

// Protocol Errors
type ProtocolError struct {
	Code    string
	Command string
	Msg     string
}

func (e *ProtocolError) Error() string {
	return e.Msg
}

func NewProtocolError(code, cmd string) *ProtocolError {
	msg := "PLC Error " + code + " on Command " + cmd
	switch code {
	case "2":
		msg += ": Illegal Value"
	case "4":
		msg += ": Illegal Format/Command or Checksum Error"
	case "A":
		msg += ": Illegal Address"
	}
	return &ProtocolError{Code: code, Command: cmd, Msg: msg}
}
