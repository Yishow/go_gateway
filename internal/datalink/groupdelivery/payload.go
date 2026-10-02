package groupdelivery

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go-gateway/internal/datalink/dbtarget"
)

const rowPayloadVersion = 1

// Cell kinds in an outbox payload. Integers and floats travel as strings so no
// JSON number can round a value; NULL is explicit.
const (
	cellNull    = "null"
	cellBool    = "bool"
	cellInt64   = "int64"
	cellFloat64 = "float64"
	cellString  = "string"
	cellTime    = "time"
)

type payloadCell struct {
	Column string `json:"column"`
	Kind   string `json:"kind"`
	Value  string `json:"value,omitempty"`
}

type rowPayload struct {
	Version     int           `json:"version"`
	RecordID    string        `json:"record_id"`
	EffectKey   string        `json:"effect_key"`
	EntityKey   string        `json:"entity_key,omitempty"`
	BucketStart string        `json:"bucket_start"`
	Partial     bool          `json:"partial,omitempty"`
	Cells       []payloadCell `json:"cells"`
}

// RowPayloadDigest returns the digest delivery stores and sends with a row, so
// other writers of the same row codec (explicit test writes) derive the same
// destination receipt content.
func RowPayloadDigest(row dbtarget.EncodedRow) (string, error) {
	_, digest, err := encodeRowPayload(row)
	return digest, err
}

func encodeRowPayload(row dbtarget.EncodedRow) (payload []byte, digest string, err error) {
	wire := rowPayload{
		Version: rowPayloadVersion, RecordID: row.RecordID, EffectKey: row.EffectKey, EntityKey: row.EntityKey,
		BucketStart: timestamp(row.BucketStart), Partial: row.Partial, Cells: make([]payloadCell, 0, len(row.Cells)),
	}
	for _, cell := range row.Cells {
		encoded, err := encodeCell(cell)
		if err != nil {
			return nil, "", err
		}
		wire.Cells = append(wire.Cells, encoded)
	}
	payload, err = json.Marshal(wire)
	if err != nil {
		return nil, "", fmt.Errorf("%w: encode row payload", ErrInvalidClosure)
	}
	sum := sha256.Sum256(payload)
	return payload, hex.EncodeToString(sum[:]), nil
}

func encodeCell(cell dbtarget.EncodedCell) (payloadCell, error) {
	out := payloadCell{Column: cell.Column}
	switch v := cell.Value.(type) {
	case nil:
		out.Kind = cellNull
	case bool:
		out.Kind, out.Value = cellBool, strconv.FormatBool(v)
	case int64:
		out.Kind, out.Value = cellInt64, strconv.FormatInt(v, 10)
	case float64:
		out.Kind, out.Value = cellFloat64, strconv.FormatFloat(v, 'g', -1, 64)
	case string:
		out.Kind, out.Value = cellString, v
	case time.Time:
		out.Kind, out.Value = cellTime, timestamp(v)
	default:
		return payloadCell{}, fmt.Errorf("%w: unsupported cell type for column %s", ErrInvalidClosure, cell.Column)
	}
	return out, nil
}

// DecodeRowPayload rebuilds an encoded row from an outbox payload with the
// original Go types, so a sender binds exactly what was closed.
func DecodeRowPayload(payload []byte) (dbtarget.EncodedRow, error) {
	var wire rowPayload
	if err := json.Unmarshal(payload, &wire); err != nil {
		return dbtarget.EncodedRow{}, fmt.Errorf("decode row payload: %w", err)
	}
	if wire.Version != rowPayloadVersion {
		return dbtarget.EncodedRow{}, fmt.Errorf("decode row payload: unsupported version %d", wire.Version)
	}
	start, err := time.Parse(time.RFC3339Nano, wire.BucketStart)
	if err != nil {
		return dbtarget.EncodedRow{}, fmt.Errorf("decode row payload time: %w", err)
	}
	row := dbtarget.EncodedRow{
		RecordID: wire.RecordID, EffectKey: wire.EffectKey, EntityKey: wire.EntityKey,
		BucketStart: start.UTC(), Partial: wire.Partial,
	}
	for _, cell := range wire.Cells {
		value, err := decodeCell(cell)
		if err != nil {
			return dbtarget.EncodedRow{}, err
		}
		row.Cells = append(row.Cells, dbtarget.EncodedCell{Column: cell.Column, Value: value})
	}
	return row, nil
}

func decodeCell(cell payloadCell) (any, error) {
	bad := fmt.Errorf("decode row payload: malformed %s cell", cell.Kind)
	switch cell.Kind {
	case cellNull:
		return nil, nil
	case cellBool:
		v, err := strconv.ParseBool(cell.Value)
		if err != nil {
			return nil, bad
		}
		return v, nil
	case cellInt64:
		v, err := strconv.ParseInt(cell.Value, 10, 64)
		if err != nil {
			return nil, bad
		}
		return v, nil
	case cellFloat64:
		v, err := strconv.ParseFloat(cell.Value, 64)
		if err != nil {
			return nil, bad
		}
		return v, nil
	case cellString:
		return cell.Value, nil
	case cellTime:
		v, err := time.Parse(time.RFC3339Nano, cell.Value)
		if err != nil {
			return nil, bad
		}
		return v.UTC(), nil
	}
	return nil, fmt.Errorf("decode row payload: unknown cell kind %q", cell.Kind)
}
