package groupdelivery

import (
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/counterproof"

	"github.com/stretchr/testify/require"
)

const legacyProofSampleDigest = "8186935d152323ef94ee48e3f310ea0b4d3daf4d1bcf1264e0769170fbe9a7a4"

const proofSampleGolden = `{"version":2,"sample_id":"s-1","member_key":"m-1","observed_at":"2026-01-01T00:00:02Z","quality":"good","value":{"type":"int64","encoding":"decimal-string","value":"7"},"source_revision":"src-1","mapping_revision":"map-1","counter_proof":{"version":1,"member":{"workspace_id":"workspace-A","group_id":"group-G","applied_revision":"rev-1","device_id":"device-1","point_id":"point-1","tag_id":"tag-1"},"digest":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}}`
const proofSampleGoldenDigest = "958c8c664c373ee7631c616528cd66019551ba9fe65ed654cf3cce9d7d80f591"

func proofSampleReference() counterproof.CounterBindingReference {
	return counterproof.CounterBindingReference{Version: 1,
		Member: counterproof.CounterBindingMemberIdentity{WorkspaceID: "workspace-A", GroupID: "group-G", AppliedRevision: "rev-1", DeviceID: "device-1", PointID: "point-1", TagID: "tag-1"},
		Digest: strings.Repeat("a", 64)}
}

func TestCounterProofSampleCodecKeepsLegacyBytesAndDigest(t *testing.T) {
	sample, err := decodeSample([]byte(legacySamplePayload))
	require.NoError(t, err)
	payload, digest, err := encodeSample(sample)
	require.NoError(t, err)
	require.Equal(t, legacySamplePayload, string(payload))
	require.Equal(t, legacyProofSampleDigest, digest)
}

func TestCounterProofSampleCodecVersionsReferenceIdentity(t *testing.T) {
	sample, err := decodeSample([]byte(legacySamplePayload))
	require.NoError(t, err)
	sample.CounterProof = proofSampleReference()
	payload, digest, err := encodeSample(sample)
	require.NoError(t, err)
	require.Equal(t, proofSampleGolden, string(payload))
	require.Equal(t, proofSampleGoldenDigest, digest)
	require.Contains(t, string(payload), `"counter_proof":{"version":1`)
	require.NotEqual(t, legacyProofSampleDigest, digest)
	decoded, err := decodeSample(payload)
	require.NoError(t, err)
	require.Equal(t, sample.CounterProof, decoded.CounterProof)
	sample.ReceivedAt = at(5).Add(time.Nanosecond)
	_, resentDigest, err := encodeSample(sample)
	require.NoError(t, err)
	require.Equal(t, digest, resentDigest)
	sample.CounterProof.Digest = strings.Repeat("b", 64)
	_, changedDigest, err := encodeSample(sample)
	require.NoError(t, err)
	require.NotEqual(t, digest, changedDigest)
}

func TestCounterProofSampleCodecRefusesUnknownOrIncompleteFormat(t *testing.T) {
	for name, extra := range map[string]string{
		"unknown":                 `"version":99,`,
		"missing reference":       `"version":2,`,
		"legacy forged reference": `"counter_proof":{},`,
		"legacy null reference":   `"counter_proof":null,`,
		"case alias zero":         `"Version":0,`,
		"explicit zero version":   `"version":0,`,
	} {
		t.Run(name, func(t *testing.T) {
			_, err := decodeSample([]byte(`{` + extra + legacySamplePayload[1:]))
			require.Error(t, err)
		})
	}
}

func TestCounterProofSampleCodecMalformedReferenceIsTyped(t *testing.T) {
	for name, wire := range map[string]string{
		"unknown proof version": strings.Replace(proofSampleGolden, `"counter_proof":{"version":1`, `"counter_proof":{"version":99`, 1),
		"duplicate digest":      strings.Replace(proofSampleGolden, `"digest":`, `"digest":"bad","digest":`, 1),
		"field alias":           strings.Replace(proofSampleGolden, `"source_revision":`, `"ſource_revision":`, 1),
		"unknown field":         strings.Replace(proofSampleGolden, `"quality":"good"`, `"quality":"good","future":"hidden"`, 1),
		"empty digest":          strings.Replace(proofSampleGolden, strings.Repeat("a", 64), "", 1),
	} {
		t.Run(name, func(t *testing.T) {
			_, err := decodeSample([]byte(wire))
			var typed *counterproof.CounterBindingProofError
			require.ErrorAs(t, err, &typed)
		})
	}
}
