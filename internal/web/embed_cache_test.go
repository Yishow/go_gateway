package web

import (
	"strings"
	"testing"
	"testing/fstest"

	"github.com/gin-gonic/gin"
)

func TestEmbeddedValidators(t *testing.T) {
	r := deliveryRouter()
	path := "/assets/main-Ab12Cd34.js"
	identity := deliveryRequest(r, "GET", path, nil)
	compressed := deliveryRequest(r, "GET", path, map[string]string{"Accept-Encoding": "gzip"})
	identityTag, gzipTag := identity.Header().Get("ETag"), compressed.Header().Get("ETag")
	for _, tag := range []string{identityTag, gzipTag} {
		if !strings.HasPrefix(tag, "\"") || !strings.HasSuffix(tag, "\"") {
			t.Fatalf("missing strong quoted ETag %q", tag)
		}
	}
	if identityTag == gzipTag {
		t.Fatal("distinct representation bytes share validator")
	}
	for _, method := range []string{"GET", "HEAD"} {
		for _, tc := range []struct {
			encoding, tag string
			status        int
		}{
			{"", identityTag, 304}, {"gzip", gzipTag, 304}, {"gzip", identityTag, 200}, {"", gzipTag, 200},
			{"gzip", "W/" + gzipTag, 304}, {"gzip", `"unrelated", ` + gzipTag, 304}, {"gzip", `"stale-build"`, 200},
		} {
			w := deliveryRequest(r, method, path, map[string]string{"Accept-Encoding": tc.encoding, "If-None-Match": tc.tag})
			if w.Code != tc.status {
				t.Errorf("%s %s %s status=%d want=%d", method, tc.encoding, tc.tag, w.Code, tc.status)
			}
			if w.Header().Get("Vary") != "Accept-Encoding" || w.Header().Get("Cache-Control") != "public, max-age=31536000, immutable" {
				t.Errorf("cache headers=%v", w.Header())
			}
			expected := identityTag
			if tc.encoding == "gzip" {
				expected = gzipTag
			}
			if w.Header().Get("ETag") != expected {
				t.Error("wrong selected ETag")
			}
			if (tc.status == 304 || method == "HEAD") && w.Body.Len() != 0 {
				t.Error("conditional/HEAD body")
			}
		}
	}
	w := deliveryRequest(r, "GET", path, map[string]string{"Accept-Encoding": "gzip, identity;q=0", "Range": "bytes=0-6", "If-None-Match": gzipTag})
	if w.Code != 406 {
		t.Fatalf("conditional bypassed forbidden Range identity: %d", w.Code)
	}
	for _, tag := range []string{identityTag, gzipTag} {
		w = deliveryRequest(r, "GET", path, map[string]string{"Accept-Encoding": "gzip", "Range": "bytes=0-6", "If-Range": tag})
		expected := 200
		if tag == identityTag {
			expected = 206
		}
		if w.Code != expected || w.Header().Get("ETag") != identityTag {
			t.Errorf("If-Range %s status=%d", tag, w.Code)
		}
	}
}

func TestEmbeddedHTMLVersionRevalidation(t *testing.T) {
	r := deliveryRouter()
	for _, path := range []string{"/", "/studio/v2", "/studio/runtime", "/studio", "/assets/plain.js"} {
		w := deliveryRequest(r, "GET", path, nil)
		if w.Code != 200 || w.Header().Get("Cache-Control") != "no-cache" {
			t.Errorf("%s cache status=%d headers=%v", path, w.Code, w.Header())
		}
		tag := w.Header().Get("ETag")
		if tag == "" {
			t.Errorf("%s missing validator", path)
		}
		for _, method := range []string{"GET", "HEAD"} {
			conditional := deliveryRequest(r, method, path, map[string]string{"If-None-Match": tag})
			if conditional.Code != 304 || conditional.Body.Len() != 0 || conditional.Header().Get("Cache-Control") != "no-cache" {
				t.Errorf("%s conditional=%d %v", path, conditional.Code, conditional.Header())
			}
		}
	}
	old := deliveryRequest(r, "GET", "/", nil).Header().Get("ETag")
	upgraded := gin.New()
	setupStaticFS(upgraded, fstest.MapFS{"static/index.html": {Data: []byte("<!doctype html><title>version two</title>")}, "static/assets/placeholder": {Data: []byte("x")}})
	w := deliveryRequest(upgraded, "GET", "/studio/v2", map[string]string{"If-None-Match": old})
	if w.Code != 200 || !strings.Contains(w.Body.String(), "version two") || w.Header().Get("ETag") == old {
		t.Error("new HTML reused old validator")
	}
}
