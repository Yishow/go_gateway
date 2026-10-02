package measurement

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestHashStringTupleMatchesPersistedJSONTupleEncoding(t *testing.T) {
	parts := []string{"device-A", "point-A", "tag-A", "true"}
	payload, err := json.Marshal(parts)
	require.NoError(t, err)
	sum := sha256.Sum256(payload)

	got, err := HashStringTuple(parts...)
	require.NoError(t, err)
	require.Equal(t, hex.EncodeToString(sum[:]), got)
}

func TestResolveAcquisitionTimeUsesTrustedSourceOrGatewayCompletion(t *testing.T) {
	sourceAt := time.Date(2026, 1, 1, 0, 0, 3, 0, time.FixedZone("source", 8*60*60))
	completedAt := time.Date(2026, 1, 1, 0, 0, 8, 0, time.UTC)

	observedAt, origin := ResolveAcquisitionTime(sourceAt, "source", completedAt)
	require.Equal(t, sourceAt.UTC(), observedAt)
	require.Equal(t, GatewayTimeOrigin(GatewayTimeOriginSource), origin)

	observedAt, origin = ResolveAcquisitionTime(sourceAt, "", completedAt)
	require.Equal(t, completedAt, observedAt)
	require.Equal(t, GatewayTimeOrigin(GatewayTimeOriginGateway), origin)
}

func TestIdentityHashesSeparateSameAddressSourcesAndStableAcquisition(t *testing.T) {
	deviceA := schema.Device{ID: "device-A", Protocol: schema.ProtocolModbusTCP, ConnectionConfig: `{"host":"a"}`}
	deviceB := schema.Device{ID: "device-B", Protocol: schema.ProtocolModbusTCP, ConnectionConfig: `{"host":"b"}`}
	pointA := schema.Point{ID: "point-A", DeviceID: "device-A", Address: "40001", DataType: schema.DataTypeUint64}
	pointB := schema.Point{ID: "point-B", DeviceID: "device-B", Address: "40001", DataType: schema.DataTypeUint64}
	tagA := schema.Tag{ID: "tag-A", Key: "same-name", DisplayName: "Same"}
	tagB := schema.Tag{ID: "tag-B", Key: "same-name-b", DisplayName: "Same"}
	mappingA := schema.Mapping{ID: "mapping-A", PointID: pointA.ID, TagID: tagA.ID, Enabled: true}
	mappingB := schema.Mapping{ID: "mapping-B", PointID: pointB.ID, TagID: tagB.ID, Enabled: true}

	sourceA, err := SourceRevision(deviceA, pointA, tagA)
	require.NoError(t, err)
	sourceB, err := SourceRevision(deviceB, pointB, tagB)
	require.NoError(t, err)
	require.NotEqual(t, sourceA, sourceB)
	mappingRevisionA, err := MappingRevision(mappingA)
	require.NoError(t, err)
	mappingRevisionB, err := MappingRevision(mappingB)
	require.NoError(t, err)
	require.NotEqual(t, mappingRevisionA, mappingRevisionB)

	first, err := StableSampleID("acquisition-A", "workspace-A", deviceA.ID, pointA.ID, tagA.ID, sourceA, mappingRevisionA)
	require.NoError(t, err)
	second, err := StableSampleID("acquisition-A", "workspace-A", deviceA.ID, pointA.ID, tagA.ID, sourceA, mappingRevisionA)
	require.NoError(t, err)
	require.Equal(t, first, second)
	other, err := StableSampleID("acquisition-A", "workspace-A", deviceB.ID, pointB.ID, tagB.ID, sourceB, mappingRevisionB)
	require.NoError(t, err)
	require.NotEqual(t, first, other)
}
