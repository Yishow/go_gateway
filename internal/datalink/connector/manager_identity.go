package connector

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"path/filepath"
	"strings"

	"go-gateway/internal/datalink/schema"
)

var errConnectionCleanup = errors.New("connection cleanup failed; endpoint remains reserved")

var errConnectionSuperseded = errors.New("connection generation superseded")
var errSerialAlias = errors.New("serial endpoint already owned; shared serial connections are not supported")

// connectionIdentity keeps the entire configuration significant, including credentials.
// JSON key order/whitespace do not change identity; the digest is never logged.
func connectionIdentity(typ schema.ProtocolType, config string) ([32]byte, string, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(config), &fields); err != nil || fields == nil {
		return [32]byte{}, "", errors.New("invalid connection configuration")
	}
	var value any
	decoder := json.NewDecoder(strings.NewReader(config))
	decoder.UseNumber()
	if err := decoder.Decode(&value); err != nil {
		return [32]byte{}, "", errors.New("invalid connection configuration")
	}
	canonical, err := json.Marshal(value)
	if err != nil {
		return [32]byte{}, "", errors.New("invalid connection configuration")
	}
	var mode, port string
	_ = json.Unmarshal(fields["mode"], &mode)
	serial := typ == schema.ProtocolModbusRTU || (typ == schema.ProtocolFatekFBs && strings.EqualFold(mode, "serial"))
	if serial {
		if json.Unmarshal(fields["serial_port"], &port) != nil || port == "" {
			return [32]byte{}, "", errors.New("serial endpoint required")
		}
		windows := strings.TrimPrefix(strings.ToUpper(port), `\\.\`)
		if strings.HasPrefix(windows, "COM") {
			port = windows
		} else {
			absolute, e := filepath.Abs(port)
			if e != nil {
				return [32]byte{}, "", errors.New("invalid serial endpoint")
			}
			port = absolute
			if resolved, e := filepath.EvalSymlinks(port); e == nil {
				port = resolved
			}
		}
	}
	return sha256.Sum256(append([]byte(string(typ)+"\x00"), canonical...)), port, nil
}
