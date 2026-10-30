package schema

// PointFormatMode records operator intent independently of projected DataFormat.
type PointFormatMode string

const (
	// PointFormatLegacy retains the protocol's historical decoder policy.
	PointFormatLegacy PointFormatMode = "legacy"
	// PointFormatInherit follows applicable rule, connection, or native defaults.
	PointFormatInherit PointFormatMode = "inherit"
	// PointFormatOverride retains an explicitly configured scalar format.
	PointFormatOverride PointFormatMode = "override"
	// PointFormatVersion identifies the current persisted intent contract.
	PointFormatVersion = 1
)
