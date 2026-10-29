package handlers

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"go-gateway/internal/datalink/runtimeview"
)

type runtimeSourcePageRequest struct {
	Kind, Query, Cursor, RuleID, RuleRevision string
	Limit                                     int
}
type runtimeSourcePageCursor struct {
	Version  int    `json:"v"`
	Binding  string `json:"b"`
	Revision string `json:"r"`
	Position int    `json:"p"`
}

var runtimeSourcePageDigits = regexp.MustCompile(`^[0-9]+$`)
var runtimeSourcePageDigest = regexp.MustCompile(`^[0-9a-f]{64}$`)

func runtimeSourcePageRequestFromQuery(q url.Values, scope *runtimeview.Scope) (runtimeSourcePageRequest, bool) {
	request := runtimeSourcePageRequest{Kind: q.Get("source_list"), Query: strings.TrimSpace(q.Get("q")), Limit: 50, Cursor: q.Get("cursor"), RuleID: q.Get("rule_id"), RuleRevision: q.Get("rule_revision")}
	if scope == nil || (request.Kind != "rules" && request.Kind != "mappings") {
		return request, false
	}
	for _, key := range []string{"source_list", "q", "limit", "cursor", "rule_id", "rule_revision"} {
		if values, ok := q[key]; ok && len(values) != 1 {
			return request, false
		}
	}
	if _, ok := q["point_id"]; ok {
		return request, false
	}
	if _, ok := q["mapping_id"]; ok {
		return request, false
	}
	if !runtimeSourcePageSafeQuery(q.Get("q")) || utf8.RuneCountInString(request.Query) > 128 {
		return request, false
	}
	if raw, ok := q["limit"]; ok {
		if !runtimeSourcePageDigits.MatchString(raw[0]) {
			return request, false
		}
		digits := strings.TrimLeft(raw[0], "0")
		if digits == "" {
			return request, false
		}
		if len(digits) > 3 {
			request.Limit = 100
		} else {
			value, _ := strconv.Atoi(digits)
			request.Limit = min(value, 100)
		}
	}
	_, rule := q["rule_id"]
	_, revision := q["rule_revision"]
	if rule || revision {
		if request.Kind != "mappings" || !rule || !revision || runtimeSafeID(request.RuleID) == "" || runtimeSafeID(request.RuleRevision) == "" {
			return request, false
		}
	}
	if values, ok := q["cursor"]; ok && values[0] == "" {
		return request, false
	}
	return request, true
}
func runtimeSourcePageSafeQuery(query string) bool {
	if !utf8.ValidString(query) {
		return false
	}
	for _, r := range query {
		if unicode.IsControl(r) {
			return false
		}
	}
	return !runtimeUnsafeNetworkText(query) && !runtimeNetworkPattern.MatchString(query) && !runtimeIPv4Pattern.MatchString(query) && !runtimeSQLPattern.MatchString(query) && !runtimeCredentialPattern.MatchString(query) && !runtimeFileURIPattern.MatchString(query) && !runtimePathPattern.MatchString(query) && !strings.Contains(query, "://")
}
func runtimeSourcePageHash(value any) (string, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:]), nil
}
func runtimeSourcePageBinding(scope runtimeview.Scope, request runtimeSourcePageRequest, position int) string {
	binding, _ := runtimeSourcePageHash([]any{1, scope.WorkspaceID, scope.WorkspaceRevision, scope.DeviceID, request.Kind, request.Query, request.Limit, request.RuleID, request.RuleRevision, position})
	return binding
}
func runtimeSourcePageNext(scope runtimeview.Scope, request runtimeSourcePageRequest, revision string, position int) *string {
	cursor := runtimeSourcePageCursor{Version: 1, Binding: runtimeSourcePageBinding(scope, request, position), Revision: revision, Position: position}
	raw, _ := json.Marshal(cursor)
	encoded := base64.RawURLEncoding.EncodeToString(raw)
	return &encoded
}
func runtimeSourcePageDecode(raw string) (runtimeSourcePageCursor, bool) {
	var result runtimeSourcePageCursor
	if len(raw) > 1024 {
		return result, false
	}
	decoded, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil || len(decoded) > 768 || base64.RawURLEncoding.EncodeToString(decoded) != raw {
		return result, false
	}
	decoder := json.NewDecoder(bytes.NewReader(decoded))
	token, err := decoder.Token()
	if err != nil || token != json.Delim('{') {
		return result, false
	}
	fields := map[string]json.RawMessage{}
	for decoder.More() {
		key, err := decoder.Token()
		if err != nil {
			return result, false
		}
		name, ok := key.(string)
		if !ok {
			return result, false
		}
		if _, exists := fields[name]; exists {
			return result, false
		}
		if name != "v" && name != "b" && name != "r" && name != "p" {
			return result, false
		}
		var value json.RawMessage
		if decoder.Decode(&value) != nil {
			return result, false
		}
		fields[name] = value
	}
	if token, err = decoder.Token(); err != nil || token != json.Delim('}') || len(fields) != 4 || decoder.Decode(new(any)) != io.EOF {
		return result, false
	}
	if json.Unmarshal(fields["v"], &result.Version) != nil || json.Unmarshal(fields["p"], &result.Position) != nil || json.Unmarshal(fields["b"], &result.Binding) != nil || json.Unmarshal(fields["r"], &result.Revision) != nil {
		return result, false
	}
	return result, result.Version == 1 && result.Position > 0 && runtimeSourcePageDigest.MatchString(result.Binding) && runtimeSourcePageDigest.MatchString(result.Revision)
}
