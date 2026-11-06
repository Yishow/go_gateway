package connector

import (
	"encoding/json"
	"slices"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func acquisitionWithRandomPoints(t *testing.T, profile string, randomPoints *int) string {
	t.Helper()
	declaration := map[string]any{"version": 1, "profile": profile}
	if randomPoints != nil {
		declaration["limits"] = map[string]any{"random_points": *randomPoints}
	}
	raw, err := json.Marshal(map[string]any{"acquisition": declaration})
	require.NoError(t, err)
	return string(raw)
}

func TestRandomProfilesCarryDocumentedLimits(t *testing.T) {
	limit := func(n int) *int { return &n }
	for _, tc := range []struct {
		name         string
		protocol     schema.ProtocolType
		profile      string
		randomPoints *int
		wantErr      bool
		wantLimit    int
	}{
		{"q-iqr default", schema.ProtocolMC3E, acquisitionProfileQIQR, nil, false, 192},
		{"q-iqr lowered", schema.ProtocolMC3E, acquisitionProfileQIQR, limit(120), false, 120},
		{"q-iqr raised", schema.ProtocolMC3E, acquisitionProfileQIQR, limit(300), true, 0},
		{"q-iqr zero", schema.ProtocolMC3E, acquisitionProfileQIQR, limit(0), true, 0},
		{"fx5 keeps 192", schema.ProtocolMC3E, acquisitionProfileFX5CPU, nil, false, 192},
		{"fx5 raised", schema.ProtocolMC3E, acquisitionProfileFX5CPU, limit(193), true, 0},
		{"fatek default", schema.ProtocolFatekFBs, acquisitionProfileFatek, nil, false, 64},
		{"fatek lowered", schema.ProtocolFatekFBs, acquisitionProfileFatek, limit(48), false, 48},
		{"fatek raised", schema.ProtocolFatekFBs, acquisitionProfileFatek, limit(65), true, 0},
		{"safe rejects random limit", schema.ProtocolMC3E, acquisitionProfileSafe, limit(10), true, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := ResolveAcquisitionConfig(tc.protocol, acquisitionWithRandomPoints(t, tc.profile, tc.randomPoints))
			if tc.wantErr {
				require.ErrorIs(t, err, ErrAcquisitionConfig)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.wantLimit, cfg.RandomPointLimit())
			capLimit := cfg.Capability(16).RandomProtocolLimit
			if tc.protocol == schema.ProtocolFatekFBs && tc.randomPoints == nil {
				// FATEK 預設由 planner 的 64 元件／64 word 加權預算約束；不寫入能力以保留既有計畫指紋。
				require.Zero(t, capLimit)
				return
			}
			require.Equal(t, tc.wantLimit, capLimit)
		})
	}
}

func TestQIQRProfileIsOfferedAndExcludesZR(t *testing.T) {
	require.True(t, slices.Contains(acquisitionProfiles(schema.ProtocolMC3E), acquisitionProfileQIQR))
	require.False(t, slices.Contains(acquisitionProfiles(schema.ProtocolFatekFBs), acquisitionProfileQIQR))

	// Q／iQ-R 是透過網路路由存取，不受 FX5 本機 CPU 路由限制。
	cfg, err := ResolveAcquisitionConfig(schema.ProtocolMC3E, `{"network_no":1,"acquisition":{"version":1,"profile":"mitsubishi-q-iqr-3e-binary-v1"}}`)
	require.NoError(t, err)
	capability := cfg.Capability(16)
	require.Equal(t, "supported", capability.RandomSupport)
	require.Equal(t, []string{"ZR"}, capability.RandomExcludedAreas)

	// bit 裝置以 16 點 word 單位走 random，共用同一個 192 上限。
	bit := cfg.Capability(1)
	require.Equal(t, "supported", bit.RandomSupport)
	require.Equal(t, 192, bit.RandomProtocolLimit)
}

func TestAcquisitionSchemaAdvertisesRandomPointsFromSameCatalog(t *testing.T) {
	var declaration map[string]any
	require.NoError(t, json.Unmarshal(AcquisitionConfigSchema(schema.ProtocolMC3E), &declaration))
	properties := declaration["properties"].(map[string]any)
	profiles := properties["profile"].(map[string]any)["enum"].([]any)
	require.Contains(t, profiles, acquisitionProfileQIQR)
	limits := properties["limits"].(map[string]any)["properties"].(map[string]any)
	require.Contains(t, limits, "random_points")
}
