package measurement

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/schema"
)

// GatewayTimeOrigin identifies whether an observation time came from a
// trusted source timestamp or from the gateway completing the acquisition.
type GatewayTimeOrigin string

const (
	GatewayTimeOriginSource  = "source"
	GatewayTimeOriginGateway = "gateway"
)

// IsTrustedSourceTime reports whether an adapter explicitly identified its
// timestamp as a source-origin timestamp.
func IsTrustedSourceTime(origin string) bool {
	switch strings.ToLower(strings.TrimSpace(origin)) {
	case GatewayTimeOriginSource, "device", "trusted-source":
		return true
	default:
		return false
	}
}

// ResolveAcquisitionTime keeps a trusted source timestamp and otherwise uses
// the completion time captured immediately after the backend read.
func ResolveAcquisitionTime(sourceAt time.Time, sourceOrigin string, completedAt time.Time) (time.Time, GatewayTimeOrigin) {
	if !sourceAt.IsZero() && IsTrustedSourceTime(sourceOrigin) {
		return sourceAt.UTC(), GatewayTimeOrigin(GatewayTimeOriginSource)
	}
	return completedAt.UTC(), GatewayTimeOrigin(GatewayTimeOriginGateway)
}

// HashStringTuple hashes the JSON array of strings used by existing source
// and mapping revisions. Keeping this tuple encoding shared prevents typed
// acquisition identities from drifting from persisted workspace revisions.
func HashStringTuple(parts ...string) (string, error) {
	payload, err := json.Marshal(parts)
	if err != nil {
		return "", fmt.Errorf("encode string identity tuple: %w", err)
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}

// SourceRevision derives the source revision tuple used by write-group
// validation from the installed device, point and tag metadata.
func SourceRevision(device schema.Device, point schema.Point, tag schema.Tag) (string, error) {
	return HashStringTuple(
		device.ID, device.Name, string(device.Protocol), device.ConnectionConfig,
		point.ID, point.DeviceID, point.Name, point.Address, point.Function, string(point.DataType),
		point.DataFormat, string(point.Mode), fmt.Sprint(point.Enabled), tag.ID, tag.Key,
		tag.DisplayName, string(tag.DataType), tag.Unit, string(tag.Status), tag.Labels,
	)
}

// MappingRevision derives the mapping revision tuple used by write-group
// validation from the installed mapping metadata.
func MappingRevision(mapping schema.Mapping) (string, error) {
	return HashStringTuple(
		mapping.ID, mapping.PointID, mapping.TagID, mapping.TransformPipeline,
		string(mapping.Status), mapping.RuleCandidateID, mapping.ProposedSignature,
		mapping.LastAppliedSignature, mapping.BlockingReason, fmt.Sprint(mapping.Enabled),
	)
}

// AcquisitionConfigFingerprint returns an opaque identity for the exact
// device and point configuration used by a poll. The input can contain
// credentials, but the returned value never exposes them.
func AcquisitionConfigFingerprint(device schema.Device, point schema.Point) (string, error) {
	return HashStringTuple(
		"acquisition-config-v1", device.ID, string(device.Protocol), device.ConnectionConfig,
		point.ID, point.DeviceID, point.Address, point.Function, string(point.DataType),
		point.DataFormat, string(point.Mode), pollingGroupID(point.PollingGroupID),
	)
}

// StableSampleID derives a retry-stable identity for one mapped acquisition.
func StableSampleID(acquisitionID, workspaceID, deviceID, pointID, tagID, sourceRevision, mappingRevision string) (string, error) {
	return HashStringTuple(
		"sample-id-v1", acquisitionID, workspaceID, deviceID, pointID, tagID,
		sourceRevision, mappingRevision,
	)
}

func pollingGroupID(id *string) string {
	if id == nil {
		return ""
	}
	return *id
}
