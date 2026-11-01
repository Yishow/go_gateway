package counterproof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCounterIdentityKeysPreserveCanonicalJSON(t *testing.T) {
	for _, tc := range []struct {
		name        string
		member      CounterBindingMemberIdentity
		key, sample string
	}{
		{name: "empty", key: `["","","","","",""]`, sample: `["","",""]`},
		{name: "escaped", member: CounterBindingMemberIdentity{WorkspaceID: "電", GroupID: `"`, AppliedRevision: `\`, DeviceID: "\n", PointID: "<>&", TagID: "\u2028"}, key: `["電","\"","\\","\n","\u003c\u003e\u0026","\u2028"]`, sample: `["\n","\u003c\u003e\u0026","\u2028"]`},
		{name: "control", member: CounterBindingMemberIdentity{DeviceID: "\x00\t\r", PointID: "point", TagID: "tag"}, key: `["","","","\u0000\t\r","point","tag"]`, sample: `["\u0000\t\r","point","tag"]`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.key, tc.member.Key())
			require.Equal(t, tc.sample, tc.member.SampleMemberKey())
		})
	}
	data, proof := mixedProofSnapshot(t, "delta", true)
	resolved, err := ResolveSnapshot(data, SnapshotScope{"w", "g", "r"})
	require.NoError(t, err)
	require.Equal(t, proof, resolved[proof.Member.Key()])
}

func TestSnapshotPreservesTypedUnsupportedProofError(t *testing.T) {
	_, err := ResolveSnapshot([]byte(`{"counter_bindings":{"member":{"version":99}}}`), SnapshotScope{"w", "g", "r"})
	require.ErrorIs(t, err, ErrCounterBindingProofUnsupported)
	var typed *CounterBindingProofError
	require.ErrorAs(t, err, &typed)
	require.Equal(t, 99, typed.Version)
}
