package desktop

import (
	"slices"
	"testing"
)

func TestWindowsFilePathAllowsPlainAbsoluteFiles(t *testing.T) {
	for _, path := range []string{`C:\data\gateway.db`, `C:\data\COMLPT1.db`, `c:/data/資料#1%20.db`, `\\server\share\gateway.db`, `\\server\share\folder\gateway.db`} {
		if err := validateWindowsFilePath(path); err != nil {
			t.Errorf("plain file rejected %q: %v", path, err)
		}
	}
}
func TestWindowsFilePathRejectsAliasesOutsideOrdinaryFiles(t *testing.T) {
	for _, path := range []string{
		`C:gateway.db`, `gateway.db`, `\data\gateway.db`, `C:\data\gateway.db:secret`, `C:\data\gateway.db::$DATA`,
		`\\?\C:\data\gateway.db`, `\\.\C:\data\gateway.db`, `\??\C:\data\gateway.db`, `\\?\UNC\server\share\gateway.db`,
		`C:\data\gateway.db.`, `C:\data \gateway.db`, `C:\data\gateway.db `, `C:\data\.\gateway.db`, `C:\data\..\gateway.db`,
		`C:\data\NUL.db`, `C:\CON`, `C:\data\COM1.txt`, `C:\data\LPT¹.db`, `C:\data\CONOUT$`, `C:\data\AUX`,
		`C:\data\file?.db`, `C:\data\file|.db`, `C:\data\\gateway.db`, `\\server\share`, `\\server\share\`, `C:\`, "C:\\data\\bad\x00.db",
	} {
		if validateWindowsFilePath(path) == nil {
			t.Errorf("unsafe path accepted %q", path)
		}
	}
}

func TestFinalWindowsPathRejectsReparseTargetNormalization(t *testing.T) {
	for _, path := range []string{`\\?\C:\data\file.db.`, `\\?\C:\data\file.db `, `\\?\C:\data\file.db:stream`, `\\?\UNC\server\share\file.db.`} {
		if _, err := normalizeFinalWindowsPath(path); err == nil {
			t.Errorf("unsafe final path accepted %q", path)
		}
	}
	for input, want := range map[string]string{`\\?\C:\data\file.db`: `C:\data\file.db`, `\\?\UNC\server\share\file.db`: `\\server\share\file.db`} {
		got, err := normalizeFinalWindowsPath(input)
		if err != nil || got != want {
			t.Errorf("canonical path mismatch %q: %q, %v", input, got, err)
		}
	}
}

func TestWindowsAncestorPathsCoverWholeCanonicalNamespace(t *testing.T) {
	cases := []struct {
		path string
		want []string
	}{
		{`C:\gateway.db`, []string{`C:\`}},
		{`C:\data\nested\gateway.db`, []string{`C:\`, `C:\data`, `C:\data\nested`}},
		{`\\server\share\gateway.db`, []string{`\\server\share\`}},
		{`\\server\share\data\gateway.db`, []string{`\\server\share\`, `\\server\share\data`}},
	}
	for _, c := range cases {
		got, err := windowsAncestorPaths(c.path)
		if err != nil || !slices.Equal(got, c.want) {
			t.Errorf("%q: %v, %v", c.path, got, err)
		}
	}
	if _, err := windowsAncestorPaths(`C:relative.db`); err == nil {
		t.Fatal("untrusted path accepted")
	}
}
