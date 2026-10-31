package measurement

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func TestMCFormatIntentChangesIdentityEvenWithSameLegacyDataFormat(t *testing.T) {
	d := schema.Device{ID: "D", Protocol: schema.ProtocolMC3E, ConnectionConfig: `{"data_format":"CDAB"}`}
	p := schema.Point{ID: "P", DeviceID: d.ID, Address: "D0", DataType: schema.DataTypeUint32, DataFormat: "ABCD"}
	tag := schema.Tag{ID: "T"}
	oldSource, err := SourceRevision(d, p, tag)
	require.NoError(t, err)
	oldFingerprint, err := AcquisitionConfigFingerprint(d, p)
	require.NoError(t, err)
	p.FormatMode = schema.PointFormatOverride
	p.FormatOverride = "ABCD"
	p.FormatVersion = 1
	source, sourceErr := SourceRevision(d, p, tag)
	fingerprint, fingerprintErr := AcquisitionConfigFingerprint(d, p)
	require.True(t, sourceErr != nil || source != oldSource, "modern intent cannot retain legacy source proof")
	require.True(t, fingerprintErr != nil || fingerprint != oldFingerprint, "modern intent cannot retain legacy acquisition proof")
}

func TestModernIdentityUsesFrozenPolicyAndDefaultIdentity(t *testing.T) {
	d := schema.Device{ID: "D", Protocol: schema.ProtocolMC3E, ConnectionConfig: `{"data_format":"CDAB","password":"private-secret"}`}
	p := schema.Point{ID: "P", DeviceID: d.ID, RevisionID: "p1", Address: "D0", DataType: schema.DataTypeUint64, DataFormat: "ABCD", FormatMode: schema.PointFormatInherit, FormatVersion: 1}
	defaults := decoding.FormatContext{Protocol: d.Protocol, RuleID: "R", RuleOwned: true, RuleRevision: "r1", RuleFormat: "ABCD", ConnectionFormat: "CDAB", ConnectionRevision: connector.ReadConnectionRevision(d.ConnectionConfig)}
	profile, err := connector.ResolveReadProfile(p, defaults)
	require.NoError(t, err)
	first, err := AcquisitionConfigFingerprint(d, p, profile)
	require.NoError(t, err)
	firstSource, err := SourceRevision(d, p, schema.Tag{ID: "T"}, profile)
	require.NoError(t, err)
	payload, err := json.Marshal(profile.IdentityView())
	require.NoError(t, err)
	require.NotEqual(t, "{}", string(payload))
	require.NotContains(t, string(payload), "private-secret")
	for name, edit := range map[string]func(*schema.Point, *decoding.FormatContext){
		"point revision":    func(p *schema.Point, c *decoding.FormatContext) { p.RevisionID = "p2" },
		"rule revision":     func(p *schema.Point, c *decoding.FormatContext) { c.RuleRevision = "r2" },
		"rule identity":     func(p *schema.Point, c *decoding.FormatContext) { c.RuleID = "R2" },
		"effective default": func(p *schema.Point, c *decoding.FormatContext) { c.RuleFormat = "DCBA" },
		"same order different source": func(p *schema.Point, c *decoding.FormatContext) {
			p.FormatMode = schema.PointFormatOverride
			p.FormatOverride = "ABCD"
		},
		"connection-selected": func(p *schema.Point, c *decoding.FormatContext) { c.RuleFormat = ""; c.ConnectionFormat = "ABCD" },
	} {
		t.Run(name, func(t *testing.T) {
			candidate, ctx := p, defaults
			edit(&candidate, &ctx)
			next, err := connector.ResolveReadProfile(candidate, ctx)
			require.NoError(t, err)
			fp, err := AcquisitionConfigFingerprint(d, candidate, next)
			require.NoError(t, err)
			require.NotEqual(t, first, fp)
			sr, err := SourceRevision(d, candidate, schema.Tag{ID: "T"}, next)
			require.NoError(t, err)
			require.NotEqual(t, firstSource, sr)
		})
	}
	unrelated := defaults
	unrelated.WorkspaceRevision = "unrelated-database-or-share-revision"
	unrelatedProfile, err := connector.ResolveReadProfile(p, unrelated)
	require.NoError(t, err)
	unchanged, err := AcquisitionConfigFingerprint(d, p, unrelatedProfile)
	require.NoError(t, err)
	require.Equal(t, first, unchanged)
	// Restart with unchanged durable input derives the same identity. Editing the
	// returned view or a later default never mutates the installed profile.
	view := profile.IdentityView()
	view.Resolution.EffectiveFormat = "DCBA"
	again, err := AcquisitionConfigFingerprint(d, p, profile)
	require.NoError(t, err)
	require.Equal(t, first, again)
	restarted, err := connector.ResolveReadProfile(p, defaults)
	require.NoError(t, err)
	again, err = AcquisitionConfigFingerprint(d, p, restarted)
	require.NoError(t, err)
	require.Equal(t, first, again)
	d.ConnectionConfig = `{"data_format":"ABCD"}`
	_, err = AcquisitionConfigFingerprint(d, p, profile)
	require.Error(t, err)
}
