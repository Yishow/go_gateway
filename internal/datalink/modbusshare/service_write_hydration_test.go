package modbusshare

import (
	"context"
	"errors"
	"testing"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/require"
)

func TestService_WriteProjectedTagValue_HydrationReadErrorFailsClosed(t *testing.T) {
	ctx := context.Background()
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	tagRecord, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "hydration-read-error", DataType: schema.DataTypeInt16})
	require.NoError(t, err)
	svc := NewService(tagSvc, 4096)
	_, err = svc.UpsertMapping(ctx, tagRecord.ID, 0)
	require.NoError(t, err)
	svc.SetHydrationError(errors.New("bootstrap read failed"))

	err = svc.WriteProjectedTagValue(ctx, tagRecord.ID, int16(1))
	var shareErr *Error
	require.ErrorAs(t, err, &shareErr)
	require.Equal(t, ErrCodeHydrationRequired, shareErr.Code)
}
