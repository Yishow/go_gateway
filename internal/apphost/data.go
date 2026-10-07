package apphost

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

// DataInput supplies explicit environment values and launcher locations.
type DataInput struct {
	Mode                   Mode
	Executable, WorkingDir string
	Environment            map[string]string
}

// DataSelection is a read-only result; confirmation precedes storage reservation.
type DataSelection struct {
	DatabasePath, DataRoot, ConfigRoot string
	Explicit, ConfirmCreate            bool
}

// ResolveData performs no mkdir, migration, movement or database creation.
// CLI callers load their existing environment first; malformed cwd dotenv retains
// the old config loader's fallback rather than being reparsed by this resolver.
func ResolveData(in DataInput) (DataSelection, error) {
	root := in.WorkingDir
	if in.Mode == Desktop {
		root = filepath.Dir(in.Executable)
	}
	root, err := filepath.EvalSymlinks(root)
	if err != nil {
		return DataSelection{}, NewFault("startup.path_unusable", err)
	}
	root, err = filepath.Abs(root)
	if err != nil {
		return DataSelection{}, NewFault("startup.path_unusable", err)
	}
	selected := databaseAlias(in.Environment)
	if selected == "" && in.Mode == Desktop {
		values, readErr := godotenv.Read(filepath.Join(root, ".env"))
		if readErr != nil && !os.IsNotExist(readErr) {
			return DataSelection{}, NewFault("startup.config_invalid", readErr)
		}
		selected = databaseAlias(values)
	}
	explicit := selected != ""
	if in.Mode == Desktop && !filepath.IsAbs(selected) {
		cwd, e := filepath.EvalSymlinks(in.WorkingDir)
		if e != nil {
			return DataSelection{}, NewFault("startup.path_unusable", e)
		}
		cwd, e = filepath.Abs(cwd)
		if e != nil {
			return DataSelection{}, NewFault("startup.path_unusable", e)
		}
		if !sameDirectory(root, cwd) {
			for _, name := range []string{".env", "datalink.db"} {
				_, e = os.Stat(filepath.Join(cwd, name))
				if e == nil {
					return DataSelection{}, NewFault("startup.path_ambiguous", errors.New("configure the intended absolute database path"))
				}
				if !os.IsNotExist(e) {
					return DataSelection{}, NewFault("startup.path_unusable", e)
				}
			}
		}
	}
	if selected == "" {
		selected = "datalink.db"
	}
	if in.Mode == Desktop && (strings.HasPrefix(selected, "file:") || strings.ContainsAny(selected, "?\x00")) {
		return DataSelection{}, NewFault("startup.path_unusable", errors.New("configure a plain database file path"))
	}
	if !filepath.IsAbs(selected) {
		selected = filepath.Join(root, selected)
	}
	parent, err := filepath.EvalSymlinks(filepath.Dir(selected))
	if err != nil {
		return DataSelection{}, NewFault("startup.path_unusable", err)
	}
	selected = filepath.Join(parent, filepath.Base(selected))
	info, err := os.Stat(selected)
	if err != nil && !os.IsNotExist(err) {
		return DataSelection{}, NewFault("startup.path_unusable", err)
	}
	exists := err == nil
	if exists {
		if !info.Mode().IsRegular() {
			return DataSelection{}, NewFault("startup.path_unusable", errors.New("database is not a regular file"))
		}
		selected, err = filepath.EvalSymlinks(selected)
		if err != nil {
			return DataSelection{}, NewFault("startup.path_unusable", err)
		}
	}
	return DataSelection{DatabasePath: selected, DataRoot: filepath.Dir(selected), ConfigRoot: root, Explicit: explicit, ConfirmCreate: in.Mode == Desktop && !explicit && !exists}, nil
}
func databaseAlias(values map[string]string) string {
	for _, key := range []string{"GATEWAY_DB_PATH", "DB_PATH", "SQLITE_PATH"} {
		if value := strings.TrimSpace(values[key]); value != "" {
			return value
		}
	}
	return ""
}
func sameDirectory(a, b string) bool {
	left, err := os.Stat(a)
	if err != nil {
		return false
	}
	right, err := os.Stat(b)
	return err == nil && os.SameFile(left, right)
}
