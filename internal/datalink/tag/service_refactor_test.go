package tag

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceCreateSetsDefaults(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo)
	ctx := context.Background()

	created, err := svc.Create(ctx, CreateTagRequest{
		Key:      "Line1/Temp",
		DataType: schema.DataTypeFloat32,
		Labels: map[string]string{
			"area": "A",
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "Line1/Temp", created.DisplayName)
	assert.Equal(t, "line1/temp", created.KeyLower)
	assert.Equal(t, schema.TagStatusDraft, created.Status)
	assert.Contains(t, created.Labels, "\"area\":\"A\"")
}

func TestServiceUpdateRejectsKeyChange(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo)
	ctx := context.Background()

	created, err := svc.Create(ctx, CreateTagRequest{
		Key:      "tag-1",
		DataType: schema.DataTypeInt16,
	})
	require.NoError(t, err)

	newKey := "tag-2"
	_, err = svc.Update(ctx, created.ID, UpdateTagRequest{Key: &newKey})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "標籤鍵不可修改")
}

func TestServiceActivateRetiredTagFails(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo)
	ctx := context.Background()

	created, err := svc.Create(ctx, CreateTagRequest{
		Key:      "tag-retire",
		DataType: schema.DataTypeString,
	})
	require.NoError(t, err)

	require.NoError(t, svc.Retire(ctx, created.ID))
	err = svc.Activate(ctx, created.ID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "已退役的標籤不能重新啟用")
}
