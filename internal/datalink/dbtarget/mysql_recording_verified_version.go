package dbtarget

// MySQLRecordingVersionVerified reports whether a raw probed server identity is
// one of the exact versions whose managed recording operation set has passed
// acceptance. It is intentionally an exact family/version pair lookup: no
// trimming, case folding, prefix stripping, or distro suffix matching.
func MySQLRecordingVersionVerified(family, rawVersion string) bool {
	switch family {
	case mysqlManagedFamilyMySQL:
		return rawVersion == "8.4.13"
	case mysqlManagedFamilyMariaDB:
		return rawVersion == "11.4.13-MariaDB-ubu2404"
	default:
		return false
	}
}
