package modbusshare

import (
	"context"
	"math"
	"net"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/protocol/modbus"
)

func setupTagSvc(t *testing.T) *tag.Service {
	t.Helper()
	repo := tag.NewMemoryRepository()
	svc := tag.NewService(repo)
	_, err := svc.Create(context.Background(), tag.CreateTagRequest{
		Key:      "test.temp.float32",
		DataType: schema.DataTypeFloat32,
	})
	if err != nil {
		t.Fatalf("create float32 tag failed: %v", err)
	}
	_, err = svc.Create(context.Background(), tag.CreateTagRequest{
		Key:      "test.count.int16",
		DataType: schema.DataTypeInt16,
	})
	if err != nil {
		t.Fatalf("create int16 tag failed: %v", err)
	}
	return svc
}

func TestService_WriteTagValue_Float32AndInt16(t *testing.T) {
	ctx := context.Background()
	tagSvc := setupTagSvc(t)
	svc := NewService(tagSvc, 8192)
	svc.SetHydrationState(HydrationState{State: HydrationStateReady, Readiness: true})

	floatTag, err := tagSvc.GetByKey(ctx, "test.temp.float32")
	if err != nil {
		t.Fatalf("get float tag failed: %v", err)
	}
	intTag, err := tagSvc.GetByKey(ctx, "test.count.int16")
	if err != nil {
		t.Fatalf("get int tag failed: %v", err)
	}

	if _, err := svc.UpsertMapping(ctx, floatTag.ID, 100); err != nil {
		t.Fatalf("upsert float mapping failed: %v", err)
	}
	if _, err := svc.UpsertMapping(ctx, intTag.ID, 200); err != nil {
		t.Fatalf("upsert int mapping failed: %v", err)
	}

	if err := svc.WriteTagValue(ctx, floatTag.ID, 12.5); err != nil {
		t.Fatalf("write float value failed: %v", err)
	}
	if err := svc.WriteTagValue(ctx, intTag.ID, 321); err != nil {
		t.Fatalf("write int value failed: %v", err)
	}

	words, err := svc.ReadHoldingWords(100, 2)
	if err != nil {
		t.Fatalf("read float words failed: %v", err)
	}
	if len(words) != 2 {
		t.Fatalf("unexpected float words len: %d", len(words))
	}

	intWords, err := svc.ReadHoldingWords(200, 1)
	if err != nil {
		t.Fatalf("read int words failed: %v", err)
	}
	if intWords[0] != 321 {
		t.Fatalf("expected 321, got %d", intWords[0])
	}
}

func TestService_WriteTagValue_RejectsDatatypeIdentityDrift(t *testing.T) {
	ctx := context.Background()
	tagSvc := setupTagSvc(t)
	svc := NewService(tagSvc, 8192)
	svc.SetHydrationState(HydrationState{State: HydrationStateReady, Readiness: true})
	intTag, err := tagSvc.GetByKey(ctx, "test.count.int16")
	if err != nil {
		t.Fatalf("get tag: %v", err)
	}
	svc.ReplaceMappings(map[string]TagMirrorMapping{intTag.ID: {
		TagID: intTag.ID, Register: 100, ZeroBasedRegister: 100, ShareStartRegister: 40101,
		SpanRegisters: 2, StrideRegisters: 2, CapacityRegisters: 4096, DataType: schema.DataTypeFloat32,
	}})
	if err := svc.WriteTagValue(ctx, intTag.ID, int16(7)); err == nil {
		t.Fatal("expected datatype identity drift to fail closed")
	}
}

func TestService_ListMappingsForWorkspaceRequiresMappingIdentityOwnership(t *testing.T) {
	tagSvc := setupTagSvc(t)
	svc := NewService(tagSvc, 4096)
	svc.SetHydrationState(HydrationState{State: HydrationStateReady, Readiness: true})
	ctx := context.Background()
	floatTag, err := tagSvc.GetByKey(ctx, "test.temp.float32")
	if err != nil {
		t.Fatalf("get float tag: %v", err)
	}
	svc.ReplaceMappings(map[string]TagMirrorMapping{floatTag.ID: {
		WorkspaceID: "ws-1", SourceRuleID: "rule-foreign", SourceRuleRevision: "r1", TagID: floatTag.ID,
		Register: 10, ShareStartRegister: 40011, ZeroBasedRegister: 10, SpanRegisters: 2, StrideRegisters: 2, CapacityRegisters: 2048, DataType: schema.DataTypeFloat32,
	}})
	svc.SetDesiredMappingOwnershipChecker(func(_ context.Context, desired DesiredMapping) error {
		if desired.SourceRuleID != "rule-owned" {
			return NewError(ErrCodeWorkspaceScope, "mapping identity is not owned", false)
		}
		return nil
	})
	mappings, err := svc.ListMappingsForWorkspace(ctx, "ws-1")
	if err != nil {
		t.Fatalf("list mappings: %v", err)
	}
	if len(mappings) != 0 {
		t.Fatalf("expected no mappings, got %v", mappings)
	}
}

func TestService_ListMappingsForWorkspacePassesCompleteIdentityToOwnership(t *testing.T) {
	tagSvc := setupTagSvc(t)
	svc := NewService(tagSvc, 4096)
	svc.SetHydrationState(HydrationState{State: HydrationStateReady, Readiness: true})
	ctx := context.Background()
	tagRecord, err := tagSvc.GetByKey(ctx, "test.temp.float32")
	if err != nil {
		t.Fatalf("get tag: %v", err)
	}
	mapping := TagMirrorMapping{WorkspaceID: "ws-1", SourceRuleID: "rule-1", SourceRuleRevision: "rev-1", TagID: tagRecord.ID, MappingID: "mapping-1", Register: 10, ShareStartRegister: 40011, ZeroBasedRegister: 10, SpanRegisters: 2, StrideRegisters: 2, CapacityRegisters: 2048, DataType: schema.DataTypeFloat32, TagKey: tagRecord.Key, DisplayName: tagRecord.DisplayName}
	svc.ReplaceMappings(map[string]TagMirrorMapping{tagRecord.ID: mapping})
	svc.SetDesiredMappingOwnershipChecker(func(_ context.Context, desired DesiredMapping) error {
		if desired.WorkspaceID != mapping.WorkspaceID || desired.SourceRuleID != mapping.SourceRuleID || desired.SourceRuleRevision != mapping.SourceRuleRevision || desired.TagID != mapping.TagID || desired.MappingID != mapping.MappingID || desired.StrideRegisters != mapping.StrideRegisters || desired.CapacityRegisters != mapping.CapacityRegisters || desired.TagKey != mapping.TagKey || desired.DisplayName != mapping.DisplayName {
			return NewError(ErrCodeWorkspaceScope, "mapping identity is incomplete", false)
		}
		return nil
	})

	mappings, err := svc.ListMappingsForWorkspace(ctx, "ws-1")
	if err != nil {
		t.Fatalf("list mappings: %v", err)
	}
	if len(mappings) != 1 || mappings[0].MappingID != mapping.MappingID {
		t.Fatalf("expected complete owned mapping, got %v", mappings)
	}
}

func TestService_Status_AfterStartAndStop(t *testing.T) {
	ctx := context.Background()
	tagSvc := setupTagSvc(t)
	svc := NewService(tagSvc, 4096)
	svc.SetHydrationState(HydrationState{State: HydrationStateReady, Readiness: true})

	if err := svc.Start(0); err == nil {
		t.Fatal("expected invalid port error")
	}

	if err := svc.Start(5020); err != nil {
		t.Fatalf("start server failed: %v", err)
	}
	st := svc.Status()
	if !st.Enabled {
		t.Fatal("expected enabled status")
	}
	if st.Port != 5020 {
		t.Fatalf("expected port 5020, got %d", st.Port)
	}

	floatTag, _ := tagSvc.GetByKey(ctx, "test.temp.float32")
	if _, err := svc.UpsertMapping(ctx, floatTag.ID, 10); err != nil {
		t.Fatalf("upsert mapping failed: %v", err)
	}
	st = svc.Status()
	if st.MappingCount != 1 {
		t.Fatalf("expected 1 mapping, got %d", st.MappingCount)
	}

	if err := svc.Stop(); err != nil {
		t.Fatalf("stop server failed: %v", err)
	}
	st = svc.Status()
	if st.Enabled {
		t.Fatal("expected disabled status after stop")
	}
	if st.Port != 0 {
		t.Fatalf("expected port 0 after stop, got %d", st.Port)
	}
	if st.BindState != "disabled" && st.BindState != "fail" {
		t.Fatalf("expected bind_state disabled or fail after stop, got %s", st.BindState)
	}
}

func TestService_Start_PortConflict_ReturnsActionableError(t *testing.T) {
	tagSvc := setupTagSvc(t)
	svc := NewService(tagSvc, 4096)

	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("listen failed: %v", err)
	}
	defer ln.Close()

	port := ln.Addr().(*net.TCPAddr).Port
	err = svc.Start(port)
	if err == nil {
		t.Fatal("expected port conflict error")
	}
	msg := strings.ToLower(err.Error())
	if !strings.Contains(msg, "already in use") {
		t.Fatalf("expected actionable conflict message, got: %v", err)
	}
}

func TestService_ModbusClientRead_MirroredValues(t *testing.T) {
	ctx := context.Background()
	tagSvc := setupTagSvc(t)
	svc := NewService(tagSvc, 8192)
	svc.SetHydrationState(HydrationState{State: HydrationStateReady, Readiness: true})

	floatTag, err := tagSvc.GetByKey(ctx, "test.temp.float32")
	if err != nil {
		t.Fatalf("get float tag failed: %v", err)
	}

	if _, err := svc.UpsertMapping(ctx, floatTag.ID, 120); err != nil {
		t.Fatalf("upsert mapping failed: %v", err)
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve random port failed: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()

	if err := svc.Start(port); err != nil {
		t.Fatalf("start server failed: %v", err)
	}
	defer svc.Stop()

	if err := svc.WriteTagValue(ctx, floatTag.ID, 12.5); err != nil {
		t.Fatalf("write mirrored value failed: %v", err)
	}

	transport := modbus.NewTCPTransport("127.0.0.1", port)
	transport.Timeout = 3 * time.Second
	client := modbus.NewClient(transport, 1)
	if err := client.Connect(); err != nil {
		t.Fatalf("modbus client connect failed: %v", err)
	}
	defer client.Close()

	registers, err := client.ReadHoldingRegisters(120, 2)
	if err != nil {
		t.Fatalf("read holding registers failed: %v", err)
	}
	if len(registers) != 2 {
		t.Fatalf("expected 2 registers, got %d", len(registers))
	}
	if registers[0] != 16712 || registers[1] != 0 {
		t.Fatalf("unexpected mirrored registers: %#v", registers)
	}
}

func TestService_RemoveMappingAndHasMapping(t *testing.T) {
	ctx := context.Background()
	tagSvc := setupTagSvc(t)
	svc := NewService(tagSvc, 4096)

	floatTag, err := tagSvc.GetByKey(ctx, "test.temp.float32")
	if err != nil {
		t.Fatalf("get float tag failed: %v", err)
	}

	if _, err := svc.UpsertMapping(ctx, floatTag.ID, 50); err != nil {
		t.Fatalf("upsert mapping failed: %v", err)
	}
	if !svc.HasMapping(floatTag.ID) {
		t.Fatal("expected mapping to exist")
	}

	svc.RemoveMapping(floatTag.ID)
	if svc.HasMapping(floatTag.ID) {
		t.Fatal("expected mapping to be removed")
	}
}

func TestToUint64_NegativeSignedValue_ReturnsError(t *testing.T) {
	if _, err := toUint64(int(-1)); err == nil {
		t.Fatal("expected error for negative signed value")
	}
}

func TestToInt64_Uint64Overflow_ReturnsError(t *testing.T) {
	if _, err := toInt64(uint64(math.MaxUint64)); err == nil {
		t.Fatal("expected overflow error for uint64 to int64 conversion")
	}
}

func TestService_WriteTagValue_RegisterOverflow_ReturnsError(t *testing.T) {
	ctx := context.Background()
	tagSvc := setupTagSvc(t)
	svc := NewService(tagSvc, 200000)

	floatTag, err := tagSvc.GetByKey(ctx, "test.temp.float32")
	if err != nil {
		t.Fatalf("get float tag failed: %v", err)
	}

	// Upserting float32 (span=2) at 65535 overflows max register capacity (65536)
	_, err = svc.UpsertMapping(ctx, floatTag.ID, math.MaxUint16)
	if err == nil {
		if err := svc.WriteTagValue(ctx, floatTag.ID, 1.25); err == nil {
			t.Fatal("expected overflow error when mapping or writing multi-word value at last register")
		}
	}
}
