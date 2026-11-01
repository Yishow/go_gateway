package point

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type pointAcquisitionLifecycle struct {
	t               *testing.T
	begun, finished int
	active          bool
	before          func()
}

func (l *pointAcquisitionLifecycle) BeginAcquisitionConfigurationMutation(_ context.Context, deviceID string) func() {
	require.Equal(l.t, "device", deviceID)
	require.False(l.t, l.active)
	l.begun++
	l.active = true
	if l.before != nil {
		l.before()
	}
	return func() { l.finished++; l.active = false }
}

type pointAcquisitionRepository struct {
	*MemoryRepository
	beforeWrite func()
}

func (r *pointAcquisitionRepository) Create(ctx context.Context, p *schema.Point) error {
	r.beforeWrite()
	return r.MemoryRepository.Create(ctx, p)
}
func (r *pointAcquisitionRepository) UpdateRevision(ctx context.Context, p *schema.Point, revision string) error {
	r.beforeWrite()
	return r.MemoryRepository.UpdateRevision(ctx, p, revision)
}
func (r *pointAcquisitionRepository) UpdateOrdinaryRevision(ctx context.Context, p *schema.Point, revision string, installed bool) error {
	r.beforeWrite()
	return r.MemoryRepository.UpdateOrdinaryRevision(ctx, p, revision, installed)
}
func (r *pointAcquisitionRepository) Delete(ctx context.Context, id string) error {
	r.beforeWrite()
	return r.MemoryRepository.Delete(ctx, id)
}
func (r *pointAcquisitionRepository) DeleteRevision(ctx context.Context, id, device, revision string) error {
	r.beforeWrite()
	return r.MemoryRepository.DeleteRevision(ctx, id, device, revision)
}
func (r *pointAcquisitionRepository) Update(ctx context.Context, p *schema.Point) error {
	r.beforeWrite()
	return r.MemoryRepository.Update(ctx, p)
}

func TestAcquisitionMutationPrecedesPointConfigurationWritesNotTelemetry(t *testing.T) {
	for _, operation := range []string{"create", "update", "ordinary", "materialize", "delete", "delete-revision", "restore", "restore-revision", "restore-absent", "telemetry", "batch-telemetry"} {
		t.Run(operation, func(t *testing.T) {
			lifecycle := &pointAcquisitionLifecycle{t: t}
			repo := &pointAcquisitionRepository{MemoryRepository: NewMemoryRepository(), beforeWrite: func() {
				require.True(t, lifecycle.active, "quarantine must precede the actual repository config write")
			}}
			service := NewService(repo, nil)
			service.SetAcquisitionMutationLifecycle(lifecycle)
			original := &schema.Point{ID: "p", DeviceID: "device", Address: "0", DataType: schema.DataTypeUint16, RevisionID: "v1", Enabled: true}
			if operation != "create" && operation != "restore-absent" {
				require.NoError(t, repo.MemoryRepository.Create(t.Context(), original))
			}
			enabled := false
			request := UpdatePointRequest{Enabled: &enabled}
			var err error
			switch operation {
			case "create":
				_, err = service.Create(t.Context(), CreatePointRequest{DeviceID: "device", Address: "0", DataType: schema.DataTypeUint16})
			case "update":
				_, err = service.Update(t.Context(), original.ID, request)
			case "ordinary":
				_, err = service.UpdateOrdinary(t.Context(), original, request, true)
			case "materialize":
				_, err = service.UpdateMaterialized(t.Context(), original, request)
			case "delete":
				err = service.Delete(t.Context(), original.ID)
			case "delete-revision":
				err = service.DeleteRevision(t.Context(), original)
			case "restore":
				err = service.Restore(t.Context(), original)
			case "restore-revision":
				err = service.RestoreRevision(t.Context(), original, original.RevisionID)
			case "restore-absent":
				err = service.RestoreAbsent(t.Context(), original)
			case "telemetry":
				err = service.UpdateReadResult(t.Context(), original.ID, uint16(7), "")
			case "batch-telemetry":
				err = service.BatchUpdateReadResult(t.Context(), []ReadResultUpdate{{PointID: original.ID, Value: uint16(7)}})
			}
			require.NoError(t, err)
			want := 1
			if operation == "telemetry" || operation == "batch-telemetry" {
				want = 0
			}
			require.Equal(t, want, lifecycle.begun)
			require.Equal(t, want, lifecycle.finished)
			require.False(t, lifecycle.active)
		})
	}
}
