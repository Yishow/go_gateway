package modbusshare

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
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

func TestService_Status_AfterStartAndStop(t *testing.T) {
	ctx := context.Background()
	tagSvc := setupTagSvc(t)
	svc := NewService(tagSvc, 4096)

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
}
