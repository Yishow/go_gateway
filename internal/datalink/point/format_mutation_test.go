package point

import (
	"context"
	"go-gateway/internal/datalink/schema"
	"testing"

	"github.com/stretchr/testify/require"
)

type formatStoreSpy struct {
	Repository
	calls int
}

func (s *formatStoreSpy) MutateFormat(context.Context, FormatMutationRequest) (*FormatMutationResult, error) {
	s.calls++
	return &FormatMutationResult{RuntimeApplyStatus: "pending"}, nil
}
func TestFormatMutationRejectsMissingGuardsBeforeStore(t *testing.T) {
	store := &formatStoreSpy{Repository: NewMemoryRepository()}
	svc := NewService(store, nil)
	for _, req := range []FormatMutationRequest{{}, {FormatScope: FormatScope{WorkspaceID: "ws", DeviceID: "dev", PointID: "point", RequestID: "request", ExpectedWorkspaceRevision: "rev"}, FormatMode: schema.PointFormatOverride, ConfiguredFormat: "ABCD"}} {
		result, err := svc.MutateFormat(t.Context(), req)
		require.ErrorIs(t, err, ErrFormatInvalid)
		require.Nil(t, result)
	}
	require.Zero(t, store.calls)
}
