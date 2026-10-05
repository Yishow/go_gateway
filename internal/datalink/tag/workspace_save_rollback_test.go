package tag

import (
	"go-gateway/internal/datalink/schema"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMemoryWorkspaceTagRollbackPreservesDirectEditAndRestoresExactOwnedFields(t *testing.T) {
	repository := NewMemoryRepository()
	service := NewService(repository)
	old := &schema.Tag{ID: "tag", Key: "owned", KeyLower: "owned", DisplayName: "Old", DataType: schema.DataTypeInt16, Unit: "A", Status: schema.TagStatusActive, Labels: `{"source":"source-rule"}`, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	require.NoError(t, repository.Create(t.Context(), old))
	authored := *old
	authored.DisplayName = "Workspace"
	authored.Unit = "B"
	authored.DataType = schema.DataTypeFloat64
	require.NoError(t, repository.Update(t.Context(), &authored))
	direct := authored
	direct.DisplayName = "Direct"
	require.NoError(t, repository.Update(t.Context(), &direct))
	require.ErrorIs(t, service.RestoreWorkspaceSave(t.Context(), &authored, old), ErrWorkspaceRollbackConflict)
	current, err := service.GetByID(t.Context(), old.ID)
	require.NoError(t, err)
	require.Equal(t, "Direct", current.DisplayName)
	require.NoError(t, repository.Update(t.Context(), &authored))
	require.NoError(t, service.RestoreWorkspaceSave(t.Context(), &authored, old))
	current, err = service.GetByID(t.Context(), old.ID)
	require.NoError(t, err)
	require.Equal(t, old, current)
}
