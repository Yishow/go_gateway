package optimizer

import (
	"go-gateway/internal/protocol/fatek"
	"go-gateway/internal/protocol/mcprotocol"
	"testing"
)

func TestAcquisitionRandomFramesMatchSerializers(t *testing.T) {
	mc := randomInput("mc", 0, 1, 10)
	frame := mcprotocol.NewRequestFrame(0, 255, 0)
	// Exact command payload produced by MCClient.RandomRead for D0/D1/D10.
	request := frame.BuildPacket(mcprotocol.CmdRandomRead, mcprotocol.SubCmdWord, []byte{3, 0, 0xA8, 0, 0, 0, 0xA8, 1, 0, 0, 0xA8, 10, 0, 0})
	response := []byte{0xD0, 0, 0, 255, 255, 3, 0, 8, 0, 0, 0, 1, 0, 2, 0, 3, 0}
	length, err := mcprotocol.ParseResponseHeader(response[:9])
	if err != nil || length != len(response)-9 {
		t.Fatal("bad response oracle")
	}
	assertRandomFrameCost(t, mc, request, response)
	f := randomInput("fatek", 0, 1, 10)
	assertRandomFrameCost(t, f, fatek.BuildFrame(1, "48", "03D00000D00001D00010"), fatek.BuildFrame(1, "48", "0000100020003"))
	f.Sources[0].Capability.RandomCrossArea = true
	f.Sources[0].Points[1].Address = "DR1"
	f.Sources[0].Points[1].DataType = "int32"
	f.Sources[0].Points[2].Address = "M10"
	f.Sources[0].Points[2].DataType = "bool"
	assertRandomFrameCost(t, f, fatek.BuildFrame(1, "48", "03D00000DR00001M0010"), fatek.BuildFrame(1, "48", "00001000000021"))
}

func assertRandomFrameCost(t *testing.T, in AcquisitionInput, request, response []byte) {
	t.Helper()
	p := BuildAcquisitionPlan(in)
	r := p.Requests()
	if len(p.Errors()) != 0 || len(r) != 1 || r[0].Kind != "random" || !r[0].FrameCostKnown || r[0].FrameBytes.RequestBytes != len(request) || r[0].FrameBytes.ResponseBytes != len(response) {
		t.Fatalf("serialized random cost: %+v errors%+v want%d/%d", r, p.Errors(), len(request), len(response))
	}
}
