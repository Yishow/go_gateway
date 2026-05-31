package workspace

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestService_ConcurrentAttachDeviceNeverLosesUpdate 驗證多個 goroutine 並行呼叫
// AttachDevice 時，每個 deviceID 都會正確出現在最終工作區，不會因為
// read-modify-write race 而被覆寫（lost update）。
func TestService_ConcurrentAttachDeviceNeverLosesUpdate(t *testing.T) {
	const n = 20

	ctx := context.Background()
	svc := NewService(NewMemoryRepository())

	var wg sync.WaitGroup
	wg.Add(n)
	for i := range n {
		go func(i int) {
			defer wg.Done()
			deviceID := fmt.Sprintf("dev-%d", i)
			_, err := svc.AttachDevice(ctx, deviceID)
			require.NoError(t, err)
		}(i)
	}
	wg.Wait()

	record, err := svc.GetOrCreate(ctx)
	require.NoError(t, err)
	require.Len(t, record.OrderedDeviceIDs, n, "所有並行寫入的 deviceID 都應存在，不應有 lost update")

	idSet := make(map[string]struct{}, len(record.OrderedDeviceIDs))
	for _, id := range record.OrderedDeviceIDs {
		idSet[id] = struct{}{}
	}
	for i := range n {
		require.Contains(t, idSet, fmt.Sprintf("dev-%d", i),
			"dev-%d 應出現在工作區，代表沒有 lost update", i)
	}
}

// TestService_ConcurrentAttachAndDetachResultsConsistent 驗證並行 Attach/Detach
// 操作不會讓工作區進入不一致狀態（例如重複 ID 或非法 ID）。
func TestService_ConcurrentAttachAndDetachResultsConsistent(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryRepository())

	// 先建立 5 個設備
	for i := range 5 {
		_, err := svc.AttachDevice(ctx, fmt.Sprintf("dev-%d", i))
		require.NoError(t, err)
	}

	var wg sync.WaitGroup
	// 並行 detach dev-0..dev-2 並 attach dev-5..dev-7
	for i := range 3 {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			_, _ = svc.DetachDevice(ctx, fmt.Sprintf("dev-%d", i))
		}(i)
		go func(i int) {
			defer wg.Done()
			_, _ = svc.AttachDevice(ctx, fmt.Sprintf("dev-%d", i+5))
		}(i)
	}
	wg.Wait()

	record, err := svc.GetOrCreate(ctx)
	require.NoError(t, err)

	// 無重複 ID
	seen := make(map[string]struct{}, len(record.OrderedDeviceIDs))
	for _, id := range record.OrderedDeviceIDs {
		require.NotContains(t, seen, id, "OrderedDeviceIDs 不應包含重複 ID")
		seen[id] = struct{}{}
	}

	// dev-3、dev-4 未被任何 goroutine detach，必須仍存在
	require.Contains(t, seen, "dev-3", "dev-3 未被 detach，應仍存在於工作區")
	require.Contains(t, seen, "dev-4", "dev-4 未被 detach，應仍存在於工作區")
}
