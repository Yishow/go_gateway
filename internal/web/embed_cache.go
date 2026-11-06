package web

import (
	"crypto/sha256"
	"fmt"
	"io/fs"
	"path"
	"strings"
)

type representationValidators struct{ identity, gzip string }

func strongValidator(data []byte) string { return fmt.Sprintf(`"%x"`, sha256.Sum256(data)) }

// Validators are content-based, so rebuilding cannot reuse an old HTML tag.
func buildStaticValidators(root fs.FS, compressed map[string][]byte) map[string]representationValidators {
	validators := make(map[string]representationValidators)
	_ = fs.WalkDir(root, ".", func(name string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return nil
		}
		data, err := fs.ReadFile(root, name)
		if err != nil {
			return nil
		}
		tags := representationValidators{identity: strongValidator(data)}
		if strings.HasPrefix(name, "assets/") {
			if encoded, ok := compressed[strings.TrimPrefix(name, "assets/")]; ok {
				tags.gzip = strongValidator(encoded)
			}
		}
		validators[name] = tags
		return nil
	})
	return validators
}

func staticCachePolicy(name string) string {
	if hashedAssetName.MatchString(path.Base(name)) {
		return "public, max-age=31536000, immutable"
	}
	return "no-cache"
}
