package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRuntimeSourcePageCursorRejectsChangedRevisionAndBindings(t *testing.T) {
	for _, change := range []string{"limit", "query", "kind", "position", "revision", "binding", "source"} {
		t.Run(change, func(t *testing.T) {
			f, q := runtimeSourcePageFixture(t, 120)
			first := runtimeSourcePageGet(t, f.h, q)
			require.NotNil(t, first.Next)
			q.Set("cursor", *first.Next)
			switch change {
			case "limit":
				q.Set("limit", "51")
			case "query":
				q.Set("q", "Meter")
			case "kind":
				q.Set("source_list", "mappings")
			case "source":
				f.h.sourceRuleReader.(runtimeSetupRuleReaderStub).rules[0].RevisionID = "changed"
			default:
				raw, err := base64.RawURLEncoding.DecodeString(*first.Next)
				require.NoError(t, err)
				var cursor map[string]any
				require.NoError(t, json.Unmarshal(raw, &cursor))
				if change == "position" {
					cursor["p"] = float64(100)
				} else if change == "revision" {
					cursor["r"] = strings.Repeat("a", 64)
				} else {
					cursor["b"] = strings.Repeat("a", 64)
				}
				raw, err = json.Marshal(cursor)
				require.NoError(t, err)
				q.Set("cursor", base64.RawURLEncoding.EncodeToString(raw))
			}
			rec := runtimeEvidenceScopeRequest(t, f.h, "?"+q.Encode())
			require.Equal(t, 409, rec.Code, rec.Body.String())
			require.NotContains(t, rec.Body.String(), `"source_page"`)
		})
	}
}
func TestRuntimeSourcePageCursorClosedSchemaAndSize(t *testing.T) {
	f, q := runtimeSourcePageFixture(t, 120)
	first := runtimeSourcePageGet(t, f.h, q)
	decoded, err := base64.RawURLEncoding.DecodeString(*first.Next)
	require.NoError(t, err)
	invalid := map[string]string{
		"padding": *first.Next + "=", "oversized-encoded": strings.Repeat("A", 1025), "oversized-decoded": base64.RawURLEncoding.EncodeToString([]byte(strings.Repeat(" ", 769))),
		"unknown":         base64.RawURLEncoding.EncodeToString([]byte(strings.TrimSuffix(string(decoded), "}") + `,"unsafe":"fixture-secret"}`)),
		"duplicate":       base64.RawURLEncoding.EncodeToString([]byte(strings.TrimSuffix(string(decoded), "}") + `,"v":1}`)),
		"fractional":      base64.RawURLEncoding.EncodeToString([]byte(strings.Replace(string(decoded), `"p":50`, `"p":50.5`, 1))),
		"string-position": base64.RawURLEncoding.EncodeToString([]byte(strings.Replace(string(decoded), `"p":50`, `"p":"50"`, 1))),
		"version":         base64.RawURLEncoding.EncodeToString([]byte(strings.Replace(string(decoded), `"v":1`, `"v":2`, 1))),
		"extra-value":     base64.RawURLEncoding.EncodeToString(append(decoded, []byte(` {}`)...)),
	}
	for name, token := range invalid {
		t.Run(name, func(t *testing.T) {
			q.Set("cursor", token)
			rec := runtimeEvidenceScopeRequest(t, f.h, "?"+q.Encode())
			require.Equal(t, 400, rec.Code, rec.Body.String())
			require.NotContains(t, rec.Body.String(), "fixture-secret")
		})
	}
	require.NotContains(t, string(decoded), "Meter")
	require.NotContains(t, string(decoded), f.record.ID)
}
func TestRuntimeSourcePageQueryUnicodeAndOrphanParameters(t *testing.T) {
	for _, text := range []string{strings.Repeat("界", 128), strings.Repeat("😀", 128), "  Meter  "} {
		f, q := runtimeSourcePageFixture(t, 1)
		q.Set("q", text)
		page := runtimeSourcePageGet(t, f.h, q)
		require.Equal(t, strings.TrimSpace(text), page.Query)
	}
	for _, text := range []string{strings.Repeat("界", 129), "Meter\nsecret", string([]byte{0xff}), "/private/fixture-secret"} {
		f, q := runtimeSourcePageFixture(t, 1)
		q.Set("q", text)
		rec := runtimeEvidenceScopeRequest(t, f.h, "?"+q.Encode())
		require.Equal(t, 400, rec.Code, rec.Body.String())
	}
	for _, key := range []string{"q", "limit", "cursor"} {
		f, q := runtimeSourcePageFixture(t, 1)
		q.Del("source_list")
		q.Set(key, "1")
		rec := runtimeEvidenceScopeRequest(t, f.h, "?"+q.Encode())
		require.Equal(t, 400, rec.Code, rec.Body.String())
	}
	for _, query := range []url.Values{{"source_list": {"rules", "rules"}}, {"source_list": {"rules"}, "q": {"a", "b"}}} {
		f, q := runtimeSourcePageFixture(t, 1)
		for k, v := range query {
			q[k] = v
		}
		rec := runtimeEvidenceScopeRequest(t, f.h, "?"+q.Encode())
		require.Equal(t, 400, rec.Code, rec.Body.String())
	}
}
