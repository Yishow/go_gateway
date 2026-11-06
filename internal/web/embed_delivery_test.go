package web

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/gin-gonic/gin"
)

var deliveryJS = []byte(strings.Repeat("console.log('embedded delivery');\n", 128))
var deliveryCSS = []byte(strings.Repeat("body { color: blue; }\n", 128))

func deliveryRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	setupStaticFS(r, fstest.MapFS{
		"static/index.html":               {Data: []byte("<!doctype html><title>version one</title>")},
		"static/assets/main-Ab12Cd34.js":  {Data: deliveryJS},
		"static/assets/main-Ef56Gh78.css": {Data: deliveryCSS},
		"static/assets/plain.js":          {Data: deliveryJS},
		"static/assets/tiny-Ab12Cd34.js":  {Data: []byte("x")},
	})
	r.GET("/api/check", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })
	r.GET("/api/events", func(c *gin.Context) { c.Header("Content-Type", "text/event-stream"); c.String(200, "data: ready\n\n") })
	return r
}

func deliveryRequest(r http.Handler, method, path string, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, http.NoBody)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestEmbeddedEncodingNegotiation(t *testing.T) {
	r := deliveryRouter()
	cases := []struct {
		name, accept, encoding string
		status                 int
	}{
		{"gzip", "gzip", "gzip", 200}, {"identity default", "", "", 200},
		{"gzip zero", "gzip;q=0", "", 200}, {"explicit beats wildcard", "gzip;q=0, *;q=1", "", 200},
		{"wildcard", "*;q=1", "gzip", 200}, {"identity forbidden", "identity;q=0, gzip", "gzip", 200},
		{"all forbidden", "identity;q=0, gzip;q=0", "", 406}, {"wildcard zero", "*;q=0", "", 406},
		{"identity exception", "*;q=0, identity;q=1", "", 200},
		{"prefer identity", "gzip;q=0.2, identity;q=0.8", "", 200},
	}
	for _, asset := range []struct {
		path, mime string
		data       []byte
	}{
		{"/assets/main-Ab12Cd34.js", "application/javascript; charset=utf-8", deliveryJS},
		{"/assets/main-Ef56Gh78.css", "text/css; charset=utf-8", deliveryCSS},
	} {
		for _, tc := range cases {
			t.Run(asset.path+"/"+tc.name, func(t *testing.T) {
				h := map[string]string{"Accept-Encoding": tc.accept}
				w := deliveryRequest(r, "GET", asset.path, h)
				if w.Code != tc.status {
					t.Fatalf("status = %d, want %d", w.Code, tc.status)
				}
				if tc.status != 200 {
					return
				}
				if w.Header().Get("Content-Encoding") != tc.encoding {
					t.Fatalf("encoding = %q, want %q", w.Header().Get("Content-Encoding"), tc.encoding)
				}
				if w.Header().Get("Vary") != "Accept-Encoding" {
					t.Errorf("Vary = %q", w.Header().Get("Vary"))
				}
				if w.Header().Get("Content-Type") != asset.mime {
					t.Errorf("MIME = %q", w.Header().Get("Content-Type"))
				}
				body := w.Body.Bytes()
				if tc.encoding == "gzip" {
					zr, err := gzip.NewReader(bytes.NewReader(body))
					if err != nil {
						t.Fatal(err)
					}
					body, err = io.ReadAll(zr)
					if err != nil {
						t.Fatal(err)
					}
					zr.Close()
					if w.Body.Len() >= len(asset.data) {
						t.Error("gzip not smaller")
					}
				}
				if !bytes.Equal(body, asset.data) {
					t.Error("body mismatch")
				}
				head := deliveryRequest(r, "HEAD", asset.path, h)
				if head.Code != w.Code || head.Body.Len() != 0 {
					t.Errorf("HEAD status/body = %d/%d", head.Code, head.Body.Len())
				}
				for _, key := range []string{"Content-Encoding", "Content-Length", "Content-Type", "Vary"} {
					if head.Header().Get(key) != w.Header().Get(key) {
						t.Errorf("HEAD %s differs", key)
					}
				}
			})
		}
	}
}

func TestEmbeddedRangeAndIsolation(t *testing.T) {
	r := deliveryRouter()
	for _, tc := range []struct {
		accept string
		status int
	}{{"gzip", 206}, {"identity;q=0, gzip", 406}, {"*;q=0", 406}} {
		w := deliveryRequest(r, "GET", "/assets/main-Ab12Cd34.js", map[string]string{"Accept-Encoding": tc.accept, "Range": "bytes=0-6"})
		if w.Code != tc.status {
			t.Errorf("range %q status=%d want=%d", tc.accept, w.Code, tc.status)
		}
		if w.Header().Get("Content-Encoding") != "" {
			t.Error("range compressed")
		}
		if tc.status == 206 && w.Body.String() != "console" {
			t.Errorf("range body=%q", w.Body.String())
		}
	}
	for _, path := range []string{"/assets/missing.js", "/assets/", "/api/missing"} {
		w := deliveryRequest(r, "GET", path, map[string]string{"Accept-Encoding": "gzip"})
		if w.Code != 404 || strings.Contains(w.Body.String(), "doctype") {
			t.Errorf("missing %s status/body=%d/%s", path, w.Code, w.Body.String())
		}
	}
	for _, path := range []string{"/api/check", "/api/events"} {
		w := deliveryRequest(r, "GET", path, map[string]string{"Accept-Encoding": "identity;q=0, gzip"})
		if w.Code != 200 || w.Header().Get("Content-Encoding") != "" || w.Header().Get("Vary") != "" {
			t.Errorf("API changed: %s %v", path, w.Result())
		}
	}
	for _, path := range []string{"/assets/plain.js", "/assets/tiny-Ab12Cd34.js"} {
		w := deliveryRequest(r, "GET", path, map[string]string{"Accept-Encoding": "gzip"})
		if w.Code != 200 || w.Header().Get("Content-Encoding") != "" {
			t.Errorf("ineligible asset changed %s", path)
		}
		w = deliveryRequest(r, "GET", path, map[string]string{"Accept-Encoding": "gzip, identity;q=0"})
		if w.Code != 406 {
			t.Errorf("unavailable representation status=%d for %s", w.Code, path)
		}
	}
}

func TestEmbeddedProtectedContract(t *testing.T) {
	r := deliveryRouter()
	for _, method := range []string{"GET", "HEAD"} {
		w := deliveryRequest(r, method, "/assets/main-Ab12Cd34.js", nil)
		if w.Code != 200 || w.Header().Get("Content-Type") != "application/javascript; charset=utf-8" {
			t.Fatalf("asset contract %d %v", w.Code, w.Header())
		}
		if method == "HEAD" && w.Body.Len() != 0 {
			t.Fatal("HEAD body")
		}
	}
	w := deliveryRequest(r, "GET", "/assets/missing.js", nil)
	if w.Code != 404 {
		t.Fatal("missing asset not404")
	}
	w = deliveryRequest(r, "GET", "/assets/main-Ab12Cd34.js", map[string]string{"Range": "bytes=0-6"})
	if w.Code != 206 || w.Body.String() != "console" {
		t.Fatal("identity Range changed")
	}
	w = deliveryRequest(r, "GET", "/api/events", nil)
	if w.Header().Get("Content-Type") != "text/event-stream" || w.Body.String() != "data: ready\n\n" {
		t.Fatal("SSE changed")
	}
}

func TestEmbeddedCompressionBudget(t *testing.T) {
	files := fstest.MapFS{
		"main-Ab12Cd34.js": {Data: deliveryJS}, "main-Ef56Gh78.css": {Data: deliveryCSS},
		"plain.js": {Data: deliveryJS}, "tiny-Ab12Cd34.js": {Data: []byte("x")},
	}
	cache := buildCompressedAssets(files)
	total := 0
	for _, payload := range cache {
		total += len(payload)
	}
	if len(cache) != 2 || total >= len(deliveryJS)+len(deliveryCSS) {
		t.Fatalf("cache entries=%d bytes=%d", len(cache), total)
	}
	t.Logf("fixed fixture gzip payload: %d bytes, eligible original: %d bytes", total, len(deliveryJS)+len(deliveryCSS)+1)
}
