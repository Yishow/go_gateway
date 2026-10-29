package handlers

import (
	"go-gateway/internal/datalink/tag"
	"strings"
)

// A semantic Tag key is not an opaque entity ID: dotted and relative slash
// namespaces are supported by the Tag authority. DNS-looking keys cannot be
// distinguished from those namespaces by syntax alone; never reuse this for IDs.
func runtimeSafeTagKey(value string) string {
	if tag.ValidateTagKey(value) != nil || strings.HasPrefix(value, "/") {
		return ""
	}
	for _, part := range strings.Split(value, "/") {
		if part == "." || part == ".." {
			return ""
		}
	}
	return value
}
