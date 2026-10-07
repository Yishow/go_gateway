package web

import (
	"testing"
	"testing/fstest"
)

func TestIncompleteEmbedGraphFailsBeforeServiceStart(t *testing.T) {
	for _, files := range []fstest.MapFS{{"static/embed-placeholder.txt": {Data: []byte("placeholder")}}, {"static/index.html": {Data: []byte(`<script src="/assets/main.js"></script>`)}}, {"static/index.html": {Data: []byte(`<script src="/assets/main.js"></script>`)}, "static/assets/main.js": {Data: []byte(`import("./logs.js")`)}}} {
		if err := ValidateStaticFiles(files); err == nil {
			t.Fatal("incomplete graph accepted")
		}
	}
	files := fstest.MapFS{"static/index.html": {Data: []byte(`<script src="/assets/main.js"></script>`)}, "static/assets/main.js": {Data: []byte(`import("./logs.js")`)}, "static/assets/logs.js": {Data: []byte(`export default {}`)}}
	if err := ValidateStaticFiles(files); err != nil {
		t.Fatal(err)
	}
}
