package web

import (
	"bytes"
	"compress/gzip"
	"io/fs"
	"path"
	"regexp"
	"strconv"
	"strings"
)

var hashedAssetName = regexp.MustCompile(`^.+-[A-Za-z0-9_-]{8,}\.[A-Za-z0-9]+$`)

// The map is populated only at startup from the finite embedded asset set.
// Retained gzip payload is strictly smaller than its source; no request adds keys.
func buildCompressedAssets(assets fs.FS) map[string][]byte {
	compressed := make(map[string][]byte)
	_ = fs.WalkDir(assets, ".", func(name string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() || !hashedAssetName.MatchString(path.Base(name)) {
			return nil
		}
		ext := strings.ToLower(path.Ext(name))
		if ext != ".js" && ext != ".css" {
			return nil
		}
		original, err := fs.ReadFile(assets, name)
		if err != nil {
			return nil
		}
		var buf bytes.Buffer
		writer := gzip.NewWriter(&buf)
		if _, err = writer.Write(original); err != nil {
			return nil
		}
		if err = writer.Close(); err != nil {
			return nil
		}
		if buf.Len() < len(original) {
			compressed[name] = bytes.Clone(buf.Bytes())
		}
		return nil
	})
	return compressed
}

func encodingQuality(values []string) (gzipQ, identityQ float64) {
	qualities := make(map[string]float64)
	for _, item := range strings.Split(strings.Join(values, ","), ",") {
		fields := strings.Split(item, ";")
		coding := strings.ToLower(strings.TrimSpace(fields[0]))
		if coding == "" {
			continue
		}
		quality := 1.0
		for _, param := range fields[1:] {
			key, value, ok := strings.Cut(strings.TrimSpace(param), "=")
			if !ok || !strings.EqualFold(strings.TrimSpace(key), "q") {
				continue
			}
			parsed, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
			if err != nil || !(parsed >= 0 && parsed <= 1) {
				quality = 0
			} else {
				quality = parsed
			}
		}
		if previous, ok := qualities[coding]; !ok || quality < previous {
			qualities[coding] = quality
		}
	}
	gzipQ = qualities["gzip"]
	if _, ok := qualities["gzip"]; !ok {
		gzipQ = qualities["*"]
	}
	identityQ = 1
	if value, ok := qualities["identity"]; ok {
		identityQ = value
	} else if value, ok := qualities["*"]; ok && value == 0 {
		identityQ = 0
	}
	return
}
