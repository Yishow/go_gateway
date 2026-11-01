package connector

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"time"

	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
)

const (
	acquisitionCatalogRevision = "acquisition-profile-catalog-v1"
	acquisitionProfileSafe     = "safe-v1"
	acquisitionProfileFatek    = "fatek-fbs-ascii-v1"
	acquisitionProfileFX5CPU   = "mitsubishi-fx5-cpu-3e-binary-v1"
)

// ErrAcquisitionConfig is intentionally free of connection settings or secrets.
var ErrAcquisitionConfig = errors.New("invalid acquisition declaration")

// AcquisitionConfig is a detached execution declaration derived only from the
// saved device configuration. It does not probe or infer device capabilities.
type AcquisitionConfig struct {
	identity          optimizer.AcquisitionIdentity
	revision, profile string
	declared          bool
	readTimeout       time.Duration
	nativeLimit       int
	frameLimit        optimizer.ApplicationFrameBudget
	coverage          optimizer.ReadableCoverage
}

// ResolveAcquisitionConfig resolves the existing saved configuration without I/O.
// Missing acquisition declarations retain unknown random support, including on
// named protocol subtypes. A selected profile is an operator declaration, not a
// successful device observation or proof that gaps are readable.
func ResolveAcquisitionConfig(protocol schema.ProtocolType, exactConfig string) (AcquisitionConfig, error) {
	fields, err := acquisitionConfigFields([]byte(exactConfig))
	if err != nil {
		return AcquisitionConfig{}, ErrAcquisitionConfig
	}
	identity, err := acquisitionTransportIdentity(protocol, exactConfig)
	if err != nil {
		return AcquisitionConfig{}, err
	}
	var timeout struct {
		Seconds int64 `json:"timeout"`
	}
	if json.Unmarshal([]byte(exactConfig), &timeout) != nil || timeout.Seconds > math.MaxInt64/int64(time.Second) {
		return AcquisitionConfig{}, ErrAcquisitionConfig
	}
	if timeout.Seconds <= 0 {
		timeout.Seconds = 5
	}
	result := AcquisitionConfig{identity: identity, readTimeout: time.Duration(timeout.Seconds) * time.Second}
	if raw, exists := fields["acquisition"]; exists {
		declaration, err := decodeAcquisitionDeclaration(raw)
		if err != nil || !acquisitionProfileCompatible(protocol, declaration.Profile) {
			return AcquisitionConfig{}, ErrAcquisitionConfig
		}
		if declaration.Profile == acquisitionProfileFX5CPU && !ownFX5CPURoute(exactConfig) {
			return AcquisitionConfig{}, ErrAcquisitionConfig
		}
		result.declared, result.profile = true, declaration.Profile
		if declaration.Limits != nil {
			result.nativeLimit = positiveAcquisitionLimit(declaration.Limits.NativeUnits)
			result.frameLimit = optimizer.ApplicationFrameBudget{RequestBytes: positiveAcquisitionLimit(declaration.Limits.RequestBytes), ResponseBytes: positiveAcquisitionLimit(declaration.Limits.ResponseBytes)}
		}
		result.coverage = declaration.coverage()
		if !optimizer.ValidateAcquisitionCoverage(identity.Protocol, result.coverage) {
			return AcquisitionConfig{}, ErrAcquisitionConfig
		}
	}
	digest := sha256.Sum256([]byte(acquisitionCatalogRevision + "\n" + string(protocol) + "\n" + exactConfig))
	result.revision = "acquisition-config-v1:" + hex.EncodeToString(digest[:])
	return result, nil
}

// Identity returns the transport-only identity. The installed caller adds the
// original polling group and cadence; neither comes from recording selectors.
func (c AcquisitionConfig) Identity(deviceID string) optimizer.AcquisitionIdentity {
	result := c.identity
	result.DeviceID = deviceID
	return result
}

// Revision identifies exact saved input and catalog semantics without exposing it.
func (c AcquisitionConfig) Revision() string { return c.revision }

// ReadTimeout preserves the existing adapter timeout and its five-second default.
func (c AcquisitionConfig) ReadTimeout() time.Duration { return c.readTimeout }

// ProfileID is empty for a historical configuration without a declaration.
func (c AcquisitionConfig) ProfileID() string { return c.profile }

// Declared distinguishes explicit safe policy from unknown historical policy.
func (c AcquisitionConfig) Declared() bool { return c.declared }

// Capability returns detached caps for one native family. The optimizer applies
// its verified command-specific client ceilings and complete-scalar boundaries.
func (c AcquisitionConfig) Capability(nativeBits int) optimizer.AcquisitionCapability {
	capability := optimizer.AcquisitionCapability{Revision: c.revision, RandomSupport: "unknown", DeviceLimit: c.nativeLimit, DeviceFrameBudget: c.frameLimit}
	switch c.identity.Protocol {
	case "modbus":
		capability.ClientLimit = 125
		if nativeBits == 1 {
			capability.ClientLimit = 2000
		}
	case "mc":
		capability.ClientLimit = 255 // Random serializer cap; contiguous remains clamped to 125.
		if nativeBits == 1 {
			capability.ClientLimit = 2000
		}
	case "fatek":
		// 255 discrete components; the planner independently clamps word commands to
		// 64/32 native components and command 48 to its 64-word weighted budget.
		capability.ClientLimit = 255
	}
	switch c.profile {
	case acquisitionProfileSafe:
		capability.RandomSupport = "unsupported"
	// FBs manual Appendix 1-3/4/5 command 48: 64 weighted words; a
	// 32-bit component costs two. The planner enforces that separate budget.
	// https://www.esea.cz/support/fatek/FBs_Manual/Manual_2/Appendix1.pdf
	case acquisitionProfileFatek:
		capability.RandomSupport, capability.RandomCrossArea = "supported", true
	// FX5 CPU built-in Ethernet only: JY997D56001K, printed pp.57,151.
	// 0403/0000 BIN: 192 word accesses; 0401/0001: 3584 bit accesses.
	// https://www.mitsubishielectric.com/dl/fa/document/manual/plcf/jy997d56001/jy997d56001k.pdf
	case acquisitionProfileFX5CPU:
		capability.RandomSupport, capability.RandomCrossArea = "supported", true
		capability.ProtocolLimit = 960       // 0401/0000 BIN word access; client remains 125.
		capability.RandomProtocolLimit = 192 // 0403/0000 BIN, independently of client 255.
		if nativeBits == 1 {
			capability.ProtocolLimit = 3584
			capability.RandomProtocolLimit = 0
			capability.RandomSupport = "unsupported"
		}
	}
	return capability
}

// Coverage returns an owned copy. A profile alone never grants readable holes.
func (c AcquisitionConfig) Coverage() optimizer.ReadableCoverage {
	result := c.coverage
	result.Readable = append([]optimizer.NativeRange(nil), result.Readable...)
	result.Forbidden = append([]optimizer.NativeRange(nil), result.Forbidden...)
	result.Excluded = append([]optimizer.NativeRange(nil), result.Excluded...)
	return result
}

func acquisitionProfileCompatible(protocol schema.ProtocolType, profile string) bool {
	for _, candidate := range acquisitionProfiles(protocol) {
		if candidate == profile {
			return true
		}
	}
	return false
}

// Both validation and the existing registry schema project this one catalog.
func acquisitionProfiles(protocol schema.ProtocolType) []string {
	if !isAcquisitionProtocol(protocol) {
		return nil
	}
	profiles := []string{acquisitionProfileSafe}
	switch protocol {
	case schema.ProtocolMC3E:
		profiles = append(profiles, acquisitionProfileFX5CPU)
	case schema.ProtocolFatekFBs:
		profiles = append(profiles, acquisitionProfileFatek)
	}
	return profiles
}
func isAcquisitionProtocol(protocol schema.ProtocolType) bool {
	switch protocol {
	case schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP, schema.ProtocolMC3E, schema.ProtocolFatekFBs:
		return true
	default:
		return false
	}
}
func acquisitionTransportIdentity(protocol schema.ProtocolType, config string) (optimizer.AcquisitionIdentity, error) {
	var route struct {
		Station byte `json:"station_no"`
		Slave   byte `json:"slave_id"`
	}
	if json.Unmarshal([]byte(config), &route) != nil {
		return optimizer.AcquisitionIdentity{}, ErrAcquisitionConfig
	}
	identity := optimizer.AcquisitionIdentity{ConnectionRevision: ReadConnectionRevision(config), ProtocolRevision: acquisitionCatalogRevision}
	switch protocol {
	case schema.ProtocolMC3E:
		identity.Protocol, identity.Frame, identity.Station = "mc", "3e-binary", strconv.Itoa(int(route.Station))
	case schema.ProtocolFatekFBs:
		if route.Station == 0 {
			route.Station = 1
		}
		identity.Protocol, identity.Frame, identity.Station = "fatek", "ascii", strconv.Itoa(int(route.Station))
	case schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP:
		if route.Slave == 0 {
			route.Slave = 1
		}
		identity.Protocol, identity.Frame, identity.Station = "modbus", "tcp", strconv.Itoa(int(route.Slave))
		if protocol == schema.ProtocolModbusRTU {
			identity.Frame = "rtu"
		}
	default:
		return optimizer.AcquisitionIdentity{}, ErrAcquisitionConfig
	}
	return identity, nil
}
func ownFX5CPURoute(config string) bool {
	var route schema.ConnectionConfigMC3E
	if json.Unmarshal([]byte(config), &route) != nil {
		return false
	}
	return route.NetworkNo == 0 && (route.PCNo == 0 || route.PCNo == 255) && (route.IONo == 0 || route.IONo == 0x03ff) && route.StationNo == 0
}
