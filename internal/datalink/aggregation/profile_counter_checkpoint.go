package aggregation

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"
)

// counterCheckpointVersion 2 added the rollover evidence to the binding; a
// version 1 baseline was taken without it and is refused, never reinterpreted.
const counterCheckpointVersion = 2

// CounterCheckpoint is the serializable baseline of one counter binding: the
// last reliable reading, and the end of the window it was taken at. It is the
// only state the cumulative path carries between windows, so it round-trips
// through JSON with every exact value intact.
type CounterCheckpoint struct {
	Binding   CounterBinding
	WindowEnd time.Time
	Reading   CounterReading
}

type counterCheckpointWire struct {
	Version   int            `json:"version"`
	Binding   CounterBinding `json:"binding"`
	WindowEnd time.Time      `json:"window_end"`
	Reading   CounterReading `json:"reading"`
}

// EncodeCounterCheckpoint serializes a checkpoint. Times are UTC and exact
// values keep their typed digits, so a reading above 2^53 never passes through
// a JSON number.
func EncodeCounterCheckpoint(checkpoint CounterCheckpoint) ([]byte, error) {
	if err := checkpoint.validate(); err != nil {
		return nil, err
	}
	reading := checkpoint.Reading
	reading.ObservedAt = reading.ObservedAt.UTC()
	if !reading.ReceivedAt.IsZero() {
		reading.ReceivedAt = reading.ReceivedAt.UTC()
	}
	return json.Marshal(counterCheckpointWire{
		Version: counterCheckpointVersion, Binding: checkpoint.Binding,
		WindowEnd: checkpoint.WindowEnd.UTC(), Reading: reading,
	})
}

// DecodeCounterCheckpoint restores a checkpoint, refusing unknown versions,
// unknown fields and incomplete content rather than guessing a baseline.
func DecodeCounterCheckpoint(data []byte) (CounterCheckpoint, error) {
	invalid := func(reason string) (CounterCheckpoint, error) {
		return CounterCheckpoint{}, fmt.Errorf("%w: %s", ErrCounterCheckpointInvalid, reason)
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var wire counterCheckpointWire
	if err := decoder.Decode(&wire); err != nil {
		return invalid("malformed")
	}
	// Only the end of input may follow the object. Any other result, a second
	// value or a syntax error in damaged bytes, is a baseline we refuse to guess.
	if _, err := decoder.Token(); !errors.Is(err, io.EOF) {
		return invalid("trailing data")
	}
	if wire.Version == 0 {
		return invalid("version is missing")
	}
	if wire.Version != counterCheckpointVersion {
		return CounterCheckpoint{}, fmt.Errorf("%w: %w", ErrCounterCheckpointInvalid, ErrCounterCheckpointVersion)
	}
	checkpoint := CounterCheckpoint{Binding: wire.Binding, WindowEnd: wire.WindowEnd.UTC(), Reading: wire.Reading}
	checkpoint.Reading.ObservedAt = checkpoint.Reading.ObservedAt.UTC()
	if !checkpoint.Reading.ReceivedAt.IsZero() {
		checkpoint.Reading.ReceivedAt = checkpoint.Reading.ReceivedAt.UTC()
	}
	if err := checkpoint.validate(); err != nil {
		return CounterCheckpoint{}, err
	}
	return checkpoint, nil
}

func (c CounterCheckpoint) validate() error {
	if err := c.Binding.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrCounterCheckpointInvalid, err)
	}
	reading := c.Reading
	if reading.SampleID == "" || reading.ObservedAt.IsZero() || c.WindowEnd.IsZero() || !c.WindowEnd.After(reading.ObservedAt) {
		return fmt.Errorf("%w: baseline reading and window end are required", ErrCounterCheckpointInvalid)
	}
	if _, numeric := reading.Value.Sign(); !numeric {
		return fmt.Errorf("%w: baseline reading is not numeric", ErrCounterCheckpointInvalid)
	}
	return nil
}
