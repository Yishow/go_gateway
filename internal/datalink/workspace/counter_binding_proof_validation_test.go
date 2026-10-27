package workspace

import (
	"encoding/json"
	"math"
	"strings"
	"testing"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestCounterBindingProofScaleAndTypeContract(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mode   CounterScaleMode
		scale  string
		source schema.DataType
		result measurement.ExactType
		signed bool
	}{
		{"unscaled unsigned", CounterScaleUnscaled, "1", schema.DataTypeUint64, measurement.ExactUint64, false},
		{"unscaled signed", CounterScaleUnscaled, "1", schema.DataTypeInt16, measurement.ExactInt64, true},
		{"unscaled float", CounterScaleUnscaled, "1", schema.DataTypeFloat32, measurement.ExactFloat64, false},
		{"scaled unsigned", CounterScaleExact, "0.1", schema.DataTypeUint64, measurement.ExactDecimal, false},
		{"scaled signed", CounterScaleExact, "0.10", schema.DataTypeInt64, measurement.ExactDecimal, true},
		{"explicit scale one", CounterScaleExact, "1", schema.DataTypeInt32, measurement.ExactDecimal, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			facts := counterProofFacts()
			facts.ScaleMode, facts.ExactScale, facts.SourceValueType, facts.ResultValueType = tc.mode, tc.scale, tc.source, tc.result
			facts.CounterPolicy = measurement.CounterPolicy{AllowNegative: tc.signed}
			if tc.signed {
				facts.SemanticKind = measurement.SemanticKindSignedCounter
			}
			proof, err := NewCounterBindingProof(facts)
			require.NoError(t, err)
			wire, err := json.Marshal(proof)
			require.NoError(t, err)
			var decoded CounterBindingProof
			require.NoError(t, json.Unmarshal(wire, &decoded))
			require.Equal(t, facts, decoded.CounterBindingFacts)
			require.Equal(t, tc.scale, decoded.ExactScale, "exact digits, including scale precision, are frozen")
		})
	}
}

func TestCounterBindingProofRefusesIncompleteOrContradictoryFacts(t *testing.T) {
	changes := map[string]func(*CounterBindingFacts){
		"workspace":          func(f *CounterBindingFacts) { f.Member.WorkspaceID = "" },
		"group":              func(f *CounterBindingFacts) { f.Member.GroupID = "" },
		"applied revision":   func(f *CounterBindingFacts) { f.Member.AppliedRevision = "" },
		"device":             func(f *CounterBindingFacts) { f.Member.DeviceID = "" },
		"point":              func(f *CounterBindingFacts) { f.Member.PointID = "" },
		"tag":                func(f *CounterBindingFacts) { f.Member.TagID = "" },
		"measurement":        func(f *CounterBindingFacts) { f.MeasurementID = "" },
		"definition":         func(f *CounterBindingFacts) { f.DefinitionRevision = "" },
		"binding":            func(f *CounterBindingFacts) { f.SourceBindingRevision = "" },
		"epoch":              func(f *CounterBindingFacts) { f.SeriesEpoch = "" },
		"unit":               func(f *CounterBindingFacts) { f.Unit = " " },
		"source revision":    func(f *CounterBindingFacts) { f.SourceRevision = "" },
		"mapping revision":   func(f *CounterBindingFacts) { f.MappingRevision = "" },
		"gauge":              func(f *CounterBindingFacts) { f.SemanticKind = measurement.SemanticKindGauge },
		"negative monotonic": func(f *CounterBindingFacts) { f.CounterPolicy.AllowNegative = true },
		"signed without policy": func(f *CounterBindingFacts) {
			f.SemanticKind = measurement.SemanticKindSignedCounter
			f.CounterPolicy = measurement.CounterPolicy{}
		},
		"rollover":             func(f *CounterBindingFacts) { f.CounterPolicy.RolloverModulus = "0" },
		"rate":                 func(f *CounterBindingFacts) { f.CounterPolicy.MaxRatePerSecond = "-1" },
		"rate without modulus": func(f *CounterBindingFacts) { f.CounterPolicy.RolloverModulus = "" },
		"unknown mode":         func(f *CounterBindingFacts) { f.ScaleMode = "" },
		"missing scale":        func(f *CounterBindingFacts) { f.ExactScale = "" },
		"negative scale":       func(f *CounterBindingFacts) { f.ExactScale = "-1" },
		"non decimal scale":    func(f *CounterBindingFacts) { f.ExactScale = "1e-1" },
		"not really unscaled": func(f *CounterBindingFacts) {
			f.ScaleMode = CounterScaleUnscaled
			f.ResultValueType = measurement.ExactUint64
		},
		"unscaled result mismatch": func(f *CounterBindingFacts) { f.ScaleMode = CounterScaleUnscaled; f.ExactScale = "1" },
		"scaled result mismatch":   func(f *CounterBindingFacts) { f.ResultValueType = measurement.ExactUint64 },
		"unknown source":           func(f *CounterBindingFacts) { f.SourceValueType = "new-type" },
		"non numeric source":       func(f *CounterBindingFacts) { f.SourceValueType = schema.DataTypeBool },
		"float cannot scale": func(f *CounterBindingFacts) {
			f.SourceValueType = schema.DataTypeFloat64
			f.CounterPolicy = measurement.CounterPolicy{}
		},
		"float cannot wrap": func(f *CounterBindingFacts) {
			f.SourceValueType = schema.DataTypeFloat64
			f.ScaleMode = CounterScaleUnscaled
			f.ExactScale = "1"
			f.ResultValueType = measurement.ExactFloat64
		},
		"invalid legacy number": func(f *CounterBindingFacts) { v := math.Inf(1); f.CounterPolicy.MaxRollOver = &v },
	}
	for name, change := range changes {
		t.Run(name, func(t *testing.T) {
			facts := counterProofFacts()
			change(&facts)
			proof, err := NewCounterBindingProof(facts)
			require.ErrorIs(t, err, ErrCounterBindingProofConflict)
			require.Zero(t, proof, "invalid facts never produce a usable sealed proof")
		})
	}
}

func TestCounterBindingProofStrictWireContract(t *testing.T) {
	for name, wire := range map[string]string{
		"null":                    `null`,
		"array":                   `[]`,
		"non integer version":     strings.Replace(counterProofJSON, `"version":1`, `"version":"1"`, 1),
		"unknown field":           strings.Replace(counterProofJSON, `"version":1`, `"version":1,"future":true`, 1),
		"duplicate version":       strings.Replace(counterProofJSON, `"version":1`, `"version":2,"version":1`, 1),
		"duplicate unit":          strings.Replace(counterProofJSON, `"unit":"kWh"`, `"unit":"MWh","unit":"kWh"`, 1),
		"duplicate nested member": strings.Replace(counterProofJSON, `"tag_id":"tag-1"`, `"tag_id":"other","tag_id":"tag-1"`, 1),
		"case alias duplicate":    strings.Replace(counterProofJSON, `"unit":"kWh"`, `"unit":"MWh","Unit":"kWh"`, 1),
		"case alias":              strings.Replace(counterProofJSON, `"unit":"kWh"`, `"Unit":"kWh"`, 1),
		"unicode case alias":      strings.Replace(counterProofJSON, `"source_revision":"source-1"`, `"source_revision":"source-2","ſource_revision":"source-1"`, 1),
		"null version":            strings.Replace(counterProofJSON, `"version":1`, `"version":null`, 1),
		"missing digest":          counterProofBody,
		"numeric exact policy":    strings.Replace(counterProofJSON, `"9007199254740993"`, `9007199254740993`, 1),
	} {
		t.Run(name, func(t *testing.T) {
			var proof CounterBindingProof
			err := json.Unmarshal([]byte(wire), &proof)
			require.ErrorIs(t, err, ErrCounterBindingProofConflict)
			require.Zero(t, proof)
		})
	}
	proof, err := NewCounterBindingProof(counterProofFacts())
	require.NoError(t, err)
	proof.Version = 2
	_, err = json.Marshal(proof)
	require.ErrorIs(t, err, ErrCounterBindingProofUnsupported)
}
