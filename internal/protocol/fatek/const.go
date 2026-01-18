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
	SymbolDD  = "DD" // 32-bit D (Double D)
	SymbolDF  = "DF" // 32-bit D (File)
	SymbolDW  = "DW" // 32-bit Discrete
	SymbolDWM = "DWM"
	SymbolWX  = "WX" // 16-bit access to discrete
	SymbolWY  = "WY"
	SymbolWM  = "WM"
	SymbolWS  = "WS"
	SymbolWT  = "WT"
	SymbolWC  = "WC"
)

// Errors
var (
	ErrInvalidSTX       = errors.New("invalid STX")
	ErrInvalidETX       = errors.New("invalid ETX")
	ErrChecksumMismatch = errors.New("checksum mismatch")
	ErrResponseTooShort = errors.New("response too short")
	ErrInvalidCommand   = errors.New("unexpected command in response")
	ErrTimeout          = errors.New("timeout waiting for response")
	ErrConnectionClosed = errors.New("connection closed")
)

// 注意：ProtocolError 已統一定義於 errors.go 中
// 請使用 errors.go 中的 NewFatekProtocolError 或 NewProtocolError
