package desktop

import (
	"errors"
	"strings"
)

// validateWindowsFilePath is platform-neutral so path-boundary negative tests run
// on every CI host. Call it before CreateFile: ADS identity can equal the base
// file's identity while its byte locks protect a different stream.
func validateWindowsFilePath(input string) error {
	path := strings.ReplaceAll(input, "/", `\`)
	var parts []string
	switch {
	case strings.HasPrefix(path, `\\`):
		parts = strings.Split(path[2:], `\`)
		if len(parts) < 3 {
			return errors.New("database requires a plain absolute UNC file path")
		}
	case len(path) >= 3 && path[1] == ':' && path[2] == '\\' && ((path[0] >= 'A' && path[0] <= 'Z') || (path[0] >= 'a' && path[0] <= 'z')):
		parts = strings.Split(path[3:], `\`)
	default:
		return errors.New("database requires a plain absolute Windows file path")
	}
	for _, part := range parts {
		if part == "" || part == "." || part == ".." || strings.HasSuffix(part, ".") || strings.HasSuffix(part, " ") {
			return errors.New("ambiguous Windows database path")
		}
		if strings.ContainsAny(part, `:<>"|?*`) {
			return errors.New("device or alternate-stream database path is unsupported")
		}
		for _, r := range part {
			if r < 32 {
				return errors.New("invalid Windows database path character")
			}
		}
		stem, _, _ := strings.Cut(part, ".")
		stem = strings.ToUpper(strings.TrimRight(stem, " "))
		switch stem {
		case "CON", "PRN", "AUX", "NUL", "CONIN$", "CONOUT$":
			return errors.New("reserved device name is not a database")
		}
		if strings.HasPrefix(stem, "COM") || strings.HasPrefix(stem, "LPT") {
			suffix := stem[3:]
			if len(suffix) == 1 && suffix[0] >= '1' && suffix[0] <= '9' || suffix == "¹" || suffix == "²" || suffix == "³" {
				return errors.New("reserved device name is not a database")
			}
		}
	}
	return nil
}

// Handle paths are validated again after conversion to ordinary Win32 syntax.
// Otherwise a reparse target created under extended-path semantics could change
// identity when the later SQLite open parses its trailing dots, spaces or stream.
func normalizeFinalWindowsPath(path string) (string, error) {
	if suffix, ok := strings.CutPrefix(path, `\\?\UNC\`); ok {
		path = `\\` + suffix
	} else {
		path = strings.TrimPrefix(path, `\\?\`)
	}
	if err := validateWindowsFilePath(path); err != nil {
		return "", err
	}
	return path, nil
}

// windowsAncestorPaths returns root-first directories, including the volume or
// share root. Inputs already describe a trusted ordinary absolute file path.
func windowsAncestorPaths(path string) ([]string, error) {
	if err := validateWindowsFilePath(path); err != nil {
		return nil, err
	}
	path = strings.ReplaceAll(path, "/", `\`)
	var root string
	var parts []string
	if strings.HasPrefix(path, `\\`) {
		parts = strings.Split(path[2:], `\`)
		root = `\\` + parts[0] + `\` + parts[1] + `\`
		parts = parts[2:]
	} else {
		root = path[:3]
		parts = strings.Split(path[3:], `\`)
	}
	directories := make([]string, 1, len(parts))
	directories[0] = root
	current := strings.TrimSuffix(root, `\`)
	for _, part := range parts[:len(parts)-1] {
		current += `\` + part
		directories = append(directories, current)
	}
	return directories, nil
}
