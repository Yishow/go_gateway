package workspace

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// ErrWriteGroupCircuitProfileConflict reports a canonical circuit-profile key
// whose saved group does not match the requested intent, or a source that is
// already owned by another group of the same profile kind and destination.
var ErrWriteGroupCircuitProfileConflict = errors.New("circuit recording profile conflict")

// WriteGroupLayoutVersion names the storage layout contract of one group.
type WriteGroupLayoutVersion string

const (
	// WriteGroupLayoutManagedWideV1 is the layout of every group saved before
	// layout versions existed; an absent value always means this layout.
	WriteGroupLayoutManagedWideV1 WriteGroupLayoutVersion = "managed_wide_v1"
	// WriteGroupLayoutManagedLongV1 is the fixed typed long table shared by the
	// canonical profile groups of one storage scope. Any other value is unknown.
	WriteGroupLayoutManagedLongV1 WriteGroupLayoutVersion = "managed_long_v1"
	// WriteGroupLayoutManagedWideProfileV1 is the per-source wide layout that
	// carries recording profiles and member metadata. It is savable, but no
	// runtime proof serves it yet, so it never starts or restores a writer.
	WriteGroupLayoutManagedWideProfileV1 WriteGroupLayoutVersion = "managed_wide_profile_v1"
)

// WriteGroupStorageMode selects how a long layout is stored.
type WriteGroupStorageMode string

// WriteGroupStorageModePlain is ordinary SQL storage and the default.
const WriteGroupStorageModePlain WriteGroupStorageMode = "plain"

// WriteGroupStorageModeTimescale is the explicit TimescaleDB hypertable mode.
const WriteGroupStorageModeTimescale WriteGroupStorageMode = "timescale"

// WriteGroupStorageModeNativeMonthly selects MySQL-family native UTC month partitions.
const WriteGroupStorageModeNativeMonthly WriteGroupStorageMode = "native_monthly"

// WriteGroupProfileKind separates legacy profiles from explicit mixed intent.
type WriteGroupProfileKind string

const (
	WriteGroupProfileMixed         WriteGroupProfileKind = "mixed"
	WriteGroupProfileInstantaneous WriteGroupProfileKind = "instantaneous"
	WriteGroupProfileCumulative    WriteGroupProfileKind = "cumulative"
)

// WriteGroupIntervalUnit is the public unit of a profile interval.
type WriteGroupIntervalUnit string

const (
	WriteGroupIntervalUnitSecond WriteGroupIntervalUnit = "second"
	WriteGroupIntervalUnitMinute WriteGroupIntervalUnit = "minute"
)

// WriteGroupSampleReduction names how a profile reduces its samples for
// storage. It is an explicit operator selection: nothing infers it from names,
// units or tag types.
type WriteGroupSampleReduction string

const (
	// WriteGroupSampleReductionLatest keeps the newest eligible sample and is
	// also the semantics of every profile saved before reductions existed.
	WriteGroupSampleReductionLatest WriteGroupSampleReduction = "latest"
	// WriteGroupSampleReductionMean averages samples over the interval. Only a
	// numeric instantaneous gauge may use it, and no runtime proof exists yet,
	// so readiness keeps such a group unwritable.
	WriteGroupSampleReductionMean WriteGroupSampleReduction = "mean"
)

// WriteGroupRecordingProfile is the recording profile of one canonical group.
// Interval is expressed in Unit; see write_group_profile_interval.go.
type WriteGroupRecordingProfile struct {
	// ResultShape versions mixed result semantics; legacy profiles omit it.
	ResultShape WriteGroupResultShape  `json:"result_shape,omitempty"`
	Kind        WriteGroupProfileKind  `json:"kind"`
	Interval    int                    `json:"interval"`
	Unit        WriteGroupIntervalUnit `json:"unit"`
	// SampleReduction is empty for every profile saved before reductions
	// existed; read it through EffectiveSampleReduction so absence stays latest
	// and a legacy profile never grows the key on the wire.
	SampleReduction WriteGroupSampleReduction `json:"sample_reduction,omitempty"`
}

// EffectiveSampleReduction returns the stored reduction, or latest for the
// profiles saved before reductions existed. A missing profile is latest too:
// it never invents a mean.
func (p *WriteGroupRecordingProfile) EffectiveSampleReduction() WriteGroupSampleReduction {
	if p == nil || p.SampleReduction == "" {
		return WriteGroupSampleReductionLatest
	}
	return p.SampleReduction
}

const (
	circuitProfileRoleCircuit    = "circuit"
	circuitProfileRoleUnassigned = "unassigned"
	maxCircuitIdentityLength     = 128
)

// EffectiveLayoutVersion returns the stored layout, or the wide layout for
// groups saved before layout versions existed.
func (g *WriteGroup) EffectiveLayoutVersion() WriteGroupLayoutVersion {
	if g == nil || g.LayoutVersion == "" {
		return WriteGroupLayoutManagedWideV1
	}
	return g.LayoutVersion
}

// circuitProfileKey is the create-once identity of one canonical profile
// group. Exactly one of CircuitID and SelectionScope is set, matching Role.
type circuitProfileKey struct {
	WorkspaceID    string
	CircuitID      string
	SelectionScope string
	Kind           WriteGroupProfileKind
	Role           string
}

// circuitProfileKeyOf derives the key from a normalized group. An assigned
// group needs one shared circuit id; an unassigned raw selection belongs to
// one device and records only the instantaneous profile.
func circuitProfileKeyOf(group *WriteGroup) (circuitProfileKey, error) {
	if group.RecordingProfile == nil || len(group.Members) == 0 {
		return circuitProfileKey{}, fmt.Errorf("%w: a recording profile and members are required", ErrWriteGroupValidation)
	}
	if isMixedWriteGroup(group) {
		return mixedRecordProfileKey(group)
	}
	key := circuitProfileKey{WorkspaceID: group.WorkspaceID, Kind: group.RecordingProfile.Kind}
	circuitID := group.Members[0].CircuitID
	deviceID := group.Members[0].DeviceID
	for i, member := range group.Members {
		if member.CircuitID != circuitID {
			return circuitProfileKey{}, fmt.Errorf("%w: member %d belongs to a different circuit", ErrWriteGroupValidation, i)
		}
		if circuitID == "" && member.DeviceID != deviceID {
			return circuitProfileKey{}, fmt.Errorf("%w: an unassigned selection must belong to one device", ErrWriteGroupValidation)
		}
	}
	if circuitID != "" {
		key.CircuitID, key.Role = circuitID, circuitProfileRoleCircuit
		return key, nil
	}
	if key.Kind != WriteGroupProfileInstantaneous {
		return circuitProfileKey{}, fmt.Errorf("%w: a cumulative profile requires a circuit", ErrWriteGroupValidation)
	}
	key.SelectionScope, key.Role = deviceID, circuitProfileRoleUnassigned
	return key, nil
}

// normalizeWriteGroupProfile trims circuit identity and rejects shapes that
// have no canonical key. It never derives a circuit from names or addresses.
func normalizeWriteGroupProfile(group *WriteGroup) error {
	if !supportedLayoutVersion(group.LayoutVersion) {
		return fmt.Errorf("%w: unsupported layout version %q", ErrWriteGroupValidation, group.LayoutVersion)
	}
	if err := normalizeWriteGroupStorageMode(group); err != nil {
		return err
	}
	// The new wide profile layout serves canonical profile groups only; the
	// existing profile rules already require a managed destination for it.
	if group.EffectiveLayoutVersion() == WriteGroupLayoutManagedWideProfileV1 && group.RecordingProfile == nil {
		return fmt.Errorf("%w: the %s layout serves recording profile groups only", ErrWriteGroupValidation, WriteGroupLayoutManagedWideProfileV1)
	}
	if group.WritePolicy.RecreateMissing != nil && !managedLongGroup(group) {
		return fmt.Errorf("%w: recreate_missing only governs managed long storage", ErrWriteGroupValidation)
	}
	for i := range group.Members {
		member := &group.Members[i]
		member.DeviceID = strings.TrimSpace(member.DeviceID)
		member.CircuitID = strings.TrimSpace(member.CircuitID)
		member.CircuitName = strings.TrimSpace(member.CircuitName)
		if member.CircuitID == "" && member.CircuitName != "" {
			return fmt.Errorf("%w: member %d has a circuit name without a circuit id", ErrWriteGroupValidation, i)
		}
		if !validCircuitText(member.CircuitID) || !validCircuitText(member.CircuitName) {
			return fmt.Errorf("%w: member %d has an invalid circuit identity", ErrWriteGroupValidation, i)
		}
		if group.RecordingProfile == nil && member.CircuitID != "" {
			return fmt.Errorf("%w: member %d circuit identity requires a recording profile", ErrWriteGroupValidation, i)
		}
	}
	// The shape guard runs before the profile check so a role can never hide
	// in a group that has no recording profile.
	if err := validateMemberMetadataShape(group); err != nil {
		return err
	}
	if group.isWideProfile() {
		if err := normalizeWideProfilePolicy(group); err != nil {
			return err
		}
	}
	if isMixedWriteGroup(group) {
		if group.Destination.StorageStrategy != WriteGroupStorageStrategyManaged {
			return fmt.Errorf("%w: mixed records require managed storage", ErrWriteGroupValidation)
		}
		if group.RowPolicy.IntervalSeconds != 0 {
			return fmt.Errorf("%w: mixed records cannot carry a group interval", ErrWriteGroupValidation)
		}
		if group.StorageMode != "" && group.StorageMode != WriteGroupStorageModePlain {
			return fmt.Errorf("%w: mixed storage mode is not supported", ErrWriteGroupValidation)
		}
		_, err := mixedRecordProfileKey(group)
		return err
	}
	profile := group.RecordingProfile
	if profile == nil {
		return nil
	}
	if err := profile.validate(); err != nil {
		return err
	}
	if group.Destination.StorageStrategy != WriteGroupStorageStrategyManaged {
		return fmt.Errorf("%w: a recording profile requires a managed destination", ErrWriteGroupValidation)
	}
	if _, err := circuitProfileKeyOf(group); err != nil {
		return err
	}
	// The profile is the only authority for the execution interval, so a
	// caller-supplied value can never contradict it.
	group.RowPolicy.IntervalSeconds = profile.IntervalSeconds()
	return nil
}

func supportedLayoutVersion(version WriteGroupLayoutVersion) bool {
	return version == "" || version == WriteGroupLayoutManagedWideV1 ||
		version == WriteGroupLayoutManagedLongV1 || version == WriteGroupLayoutManagedWideProfileV1
}

// isLongLayout reports whether the group uses the shared typed long table.
func (g *WriteGroup) isLongLayout() bool {
	return g.EffectiveLayoutVersion() == WriteGroupLayoutManagedLongV1
}

func validCircuitText(value string) bool {
	return utf8.RuneCountInString(value) <= maxCircuitIdentityLength &&
		strings.IndexFunc(value, unicode.IsControl) < 0
}

func sameRecordingProfile(left, right *WriteGroupRecordingProfile) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	// The reduction compares through its effective value: a profile that omits
	// the field and one that explicitly pins latest are the same intent, so a
	// no-op resend replays the saved group instead of conflicting.
	return left.ResultShape == right.ResultShape && left.Kind == right.Kind && left.Interval == right.Interval && left.Unit == right.Unit &&
		left.EffectiveSampleReduction() == right.EffectiveSampleReduction()
}

// normalizeWriteGroupStorageMode keeps the long layout and its mode together:
// the long table only serves canonical profile groups with explicitly selected
// storage, and a wide group never carries a mode. Connector validation checks
// the dialect-specific modes separately.
func normalizeWriteGroupStorageMode(group *WriteGroup) error {
	if !group.isLongLayout() {
		if group.StorageMode != "" {
			return fmt.Errorf("%w: a storage mode requires the long layout", ErrWriteGroupValidation)
		}
		return nil
	}
	if group.RecordingProfile == nil {
		return fmt.Errorf("%w: the long layout serves recording profile groups only", ErrWriteGroupValidation)
	}
	switch group.StorageMode {
	case "":
		group.StorageMode = WriteGroupStorageModePlain
	case WriteGroupStorageModePlain, WriteGroupStorageModeTimescale, WriteGroupStorageModeNativeMonthly:
	default:
		return fmt.Errorf("%w: unsupported storage mode %q", ErrWriteGroupValidation, group.StorageMode)
	}
	return nil
}

// managedLongGroup reports whether the group is stored in the shared managed
// long table, the only storage whose lost table can be restored from a
// confirmed layout.
func managedLongGroup(group *WriteGroup) bool {
	return group.isLongLayout() && group.Destination.StorageStrategy == WriteGroupStorageStrategyManaged
}

// managedCutoverTarget reports whether the group is a layout a reviewed cutover may
// start: the long layout or the known managed wide profile layout. It only names the
// layouts; a target still has to pass the full readiness and runtime proof of its own
// layout, and any other layout stays refused.
func managedCutoverTarget(group *WriteGroup) bool {
	return managedLongGroup(group) || (group.isWideProfile() && group.Destination.StorageStrategy == WriteGroupStorageStrategyManaged)
}

// defaultRecreateMissing saves the design default, recovery on, as an explicit
// value on a new managed long group that did not choose. It runs only when a
// group is created: a saved group that never confirmed a policy is not given
// one later, so an absent policy is never read as consent.
func defaultRecreateMissing(group *WriteGroup) {
	if managedLongGroup(group) && group.WritePolicy.RecreateMissing == nil {
		on := true
		group.WritePolicy.RecreateMissing = &on
	}
}

// sameRecreateMissing compares the saved policy tri-state exactly.
func sameRecreateMissing(left, right *bool) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}
