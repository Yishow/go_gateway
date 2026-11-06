package web

import (
	"errors"
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/gin-gonic/gin"
)

type failingStaticFS struct {
	fs.FS
	failPath string
}

func (f failingStaticFS) Open(name string) (fs.File, error) {
	if name == f.failPath {
		return nil, &fs.PathError{Op: "open", Path: name, Err: errors.New("private backend failure detail")}
	}
	return f.FS.Open(name)
}

func TestEmbeddedCacheFailureDiscardsPartialState(t *testing.T) {
	for _, failure := range []struct{ name, path string }{
		{"file read after valid asset", "static/assets/z-failed-Ab12Cd34.js"},
		{"directory walk after valid asset", "static/assets/z-failed"},
		{"validator-only walk failure", "static/z-failed"},
	} {
		t.Run(failure.name, func(t *testing.T) {
			files := fstest.MapFS{
				"static/index.html":                         {Data: []byte("<!doctype html><title>available</title>")},
				"static/assets/a-valid-Ab12Cd34.js":         {Data: deliveryJS},
				"static/assets/z-failed-Ab12Cd34.js":        {Data: deliveryJS},
				"static/assets/z-failed/nested-Ab12Cd34.js": {Data: deliveryJS},
				"static/z-failed/file.txt":                  {Data: []byte("unreadable")},
			}
			faulty := failingStaticFS{FS: files, failPath: failure.path}
			static, err := fs.Sub(faulty, "static")
			if err != nil {
				t.Fatal(err)
			}
			assets, err := fs.Sub(static, "assets")
			if err != nil {
				t.Fatal(err)
			}
			compressed := buildCompressedAssets(assets)
			if strings.HasPrefix(failure.path, "static/assets/") && len(compressed) != 0 {
				t.Errorf("partial compressed cache survived failure: %d entries", len(compressed))
			}
			if validators := buildStaticValidators(static, compressed); len(validators) != 0 {
				t.Errorf("partial validators survived failure: %d entries", len(validators))
			}
			gin.SetMode(gin.TestMode)
			r := gin.New()
			setupStaticFS(r, faulty)
			for _, path := range []string{"/assets/a-valid-Ab12Cd34.js", "/studio/v2"} {
				w := deliveryRequest(r, "GET", path, map[string]string{"Accept-Encoding": "gzip"})
				if w.Code != 200 {
					t.Errorf("valid identity unavailable: %s status=%d", path, w.Code)
				}
				if w.Header().Get("Content-Encoding") != "" || w.Header().Get("ETag") != "" {
					t.Errorf("partial derived state served: %s %v", path, w.Header())
				}
				if strings.Contains(w.Body.String(), "private backend") {
					t.Error("raw failure leaked")
				}
				if path == "/assets/a-valid-Ab12Cd34.js" && w.Body.String() != string(deliveryJS) {
					t.Error("valid identity bytes changed")
				}
			}
			refused := deliveryRequest(r, "GET", "/assets/a-valid-Ab12Cd34.js", map[string]string{"Accept-Encoding": "gzip, identity;q=0"})
			if refused.Code != 406 {
				t.Errorf("fallback used rejected identity or partial gzip: %d", refused.Code)
			}
		})
	}
}
