package optimizer

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionRandomProtocolLimitOnlyTightensRandom(t *testing.T) {
	in := randomInput("mc", 0, 10, 20)
	in.Sources[0].Capability.RandomProtocolLimit = 2
	plan := BuildAcquisitionPlan(in)
	require.Empty(t, plan.Errors())
	require.Len(t, plan.Requests(), 2)
	for _, request := range plan.Requests() {
		require.LessOrEqual(t, request.Length, 2)
	}
	// The same random-only cap must not split a safe contiguous command.
	dense := randomInput("mc", 0, 1, 2)
	dense.Sources[0].Capability.RandomProtocolLimit = 1
	plan = BuildAcquisitionPlan(dense)
	require.Empty(t, plan.Errors())
	require.Len(t, plan.Requests(), 1)
	require.Equal(t, "contiguous", plan.Requests()[0].Kind)
	require.Equal(t, 3, plan.Requests()[0].Length)
	// A larger random-only cap cannot loosen an existing common cap.
	in.Sources[0].Capability.ProtocolLimit = 1
	in.Sources[0].Capability.RandomProtocolLimit = 999
	plan = BuildAcquisitionPlan(in)
	require.Empty(t, plan.Errors())
	require.Len(t, plan.Requests(), 3)
	for _, request := range plan.Requests() {
		require.Equal(t, 1, request.Length)
	}
	// Nor can it loosen the actual random serializer's client ceiling.
	offsets := make([]int, 256)
	for i := range offsets {
		offsets[i] = i * 10
	}
	in = randomInput("mc", offsets...)
	in.Sources[0].Capability.RandomProtocolLimit = 999
	plan = BuildAcquisitionPlan(in)
	require.Empty(t, plan.Errors())
	require.Len(t, plan.Requests(), 2)
	require.Equal(t, 255, plan.Requests()[0].Length)
}

func TestAcquisitionRandomProtocolLimitFingerprintAndValidation(t *testing.T) {
	in := randomInput("mc", 0, 10, 20)
	before := BuildAcquisitionPlan(in).Fingerprint()
	// Omitting the new optional field retains the previous canonical JSON.
	legacy := `{"Revision":"cap1","RandomCrossArea":false,"RandomSupport":"supported","ProtocolLimit":1000,"ClientLimit":0,"DeviceLimit":0,"ProtocolFrameBudget":{"RequestBytes":0,"ResponseBytes":0},"ClientFrameBudget":{"RequestBytes":0,"ResponseBytes":0},"DeviceFrameBudget":{"RequestBytes":0,"ResponseBytes":0},"ProtocolAddressBounds":null,"ClientAddressBounds":null,"DeviceAddressBounds":null}`
	require.Equal(t, legacy, canonicalJSON(in.Sources[0].Capability))
	in.Sources[0].Capability.RandomProtocolLimit = 192
	require.True(t, strings.Contains(canonicalJSON(in.Sources[0].Capability), `"RandomProtocolLimit":192`))
	require.NotEqual(t, before, BuildAcquisitionPlan(in).Fingerprint())
	in.Sources[0].Capability.RandomProtocolLimit = -1
	plan := BuildAcquisitionPlan(in)
	require.Empty(t, plan.Requests())
	require.Len(t, plan.Errors(), 3)
	for _, err := range plan.Errors() {
		require.Equal(t, "invalid-capability", err.Code)
	}
}
