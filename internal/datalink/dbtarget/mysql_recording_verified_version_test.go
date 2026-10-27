package dbtarget

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMySQLRecordingVersionVerifiedExactPairs(t *testing.T) {
	for _, tc := range []struct {
		name, family, version string
		want                  bool
	}{
		{name: "mysql verified", family: "mysql", version: "8.4.13", want: true},
		{name: "mariadb verified", family: "mariadb", version: "11.4.13-MariaDB-ubu2404", want: true},
		{name: "mysql older patch", family: "mysql", version: "8.4.4"},
		{name: "mysql newer patch", family: "mysql", version: "8.4.14"},
		{name: "mariadb older patch", family: "mariadb", version: "11.4.2-MariaDB-ubu2204"},
		{name: "mariadb newer patch", family: "mariadb", version: "11.4.14-MariaDB-ubu2404"},
		{name: "mariadb distro suffix", family: "mariadb", version: "11.4.13-MariaDB-otherdistro"},
		{name: "mysql debug suffix", family: "mysql", version: "8.4.13-debug"},
		{name: "mariadb debug suffix", family: "mariadb", version: "11.4.13-MariaDB-ubu2404-debug"},
		{name: "prefix alias is not accepted", family: "mariadb", version: "5.5.5-11.4.13-MariaDB-ubu2404"},
		{name: "leading whitespace", family: "mysql", version: " 8.4.13"},
		{name: "trailing whitespace", family: "mariadb", version: "11.4.13-MariaDB-ubu2404 "},
		{name: "uppercase family", family: "MySQL", version: "8.4.13"},
		{name: "wrong family", family: "postgres", version: "8.4.13"},
		{name: "empty version", family: "mysql", version: ""},
		{name: "empty family", family: "", version: "8.4.13"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, MySQLRecordingVersionVerified(tc.family, tc.version))
		})
	}
}
