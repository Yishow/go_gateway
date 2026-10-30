package connector

import (
	"crypto/sha256"
	"encoding/hex"
	"go-gateway/internal/datalink/schema"
)

// ReadConnectionRevision binds a resolved default to its exact configuration
// input, without retaining or exposing connection secrets in the profile.
func ReadConnectionRevision(config string) string {
	digest := sha256.Sum256([]byte(config))
	return "read-connection-v1:" + hex.EncodeToString(digest[:])
}
func (p ReadProfile) MatchesConnection(protocol schema.ProtocolType, config string) bool {
	if !p.Present() {
		return true
	}
	if p.Protocol() != protocol {
		return false
	}
	revision := p.Resolution().ConnectionRevision
	return revision == "" || revision == ReadConnectionRevision(config)
}
