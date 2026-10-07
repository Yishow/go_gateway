package web

import (
	"errors"
	"io/fs"
	"path"
	"regexp"
	"strings"
)

const assetJSExt = ".js"

var assetReference = regexp.MustCompile(`(?:["'(])((?:/assets/|\./)[^"'()\s]+\.(?:js|css|svg|png|woff2?))(?:["')])`)

// ValidateStaticFiles rejects placeholders and incomplete lazy asset graphs.
func ValidateStaticFiles(files fs.FS) error {
	root, err := fs.Sub(files, "static")
	if err != nil {
		return err
	}
	index, err := fs.ReadFile(root, "index.html")
	if err != nil || !strings.Contains(string(index), "/assets/") {
		return errors.New("embedded frontend entry is missing")
	}
	foundScript := false
	err = fs.WalkDir(root, ".", func(name string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		ext := path.Ext(name)
		if ext != ".html" && ext != assetJSExt && ext != ".css" {
			return nil
		}
		if ext == assetJSExt {
			foundScript = true
		}
		content, readErr := fs.ReadFile(root, name)
		if readErr != nil {
			return readErr
		}
		for _, match := range assetReference.FindAllSubmatch(content, -1) {
			ref := string(match[1])
			if strings.HasPrefix(ref, "/") {
				ref = strings.TrimPrefix(ref, "/")
			} else {
				ref = path.Join(path.Dir(name), ref)
			}
			info, e := fs.Stat(root, ref)
			if e != nil || !info.Mode().IsRegular() {
				return errors.New("embedded frontend asset graph is incomplete")
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if !foundScript {
		return errors.New("embedded frontend scripts are missing")
	}
	return nil
}
