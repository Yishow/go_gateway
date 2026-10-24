package grouppipeline

import (
	"testing"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

// TestLongStorageRecoveryIsWiredByNew checks the one place production gets its
// storage recovery from: a pipeline built over a delivery store restores
// confirmed managed tables through the sender it creates, on the operation
// ledger the caller shares with schema confirmation.
func TestLongStorageRecoveryIsWiredByNew(t *testing.T) {
	t.Run("a store with a database installs recovery into the sender", func(t *testing.T) {
		p := newFixture(t).pipeline()
		require.NotNil(t, p.recovery)
		require.Same(t, p.recovery, p.senderConfig().Recovery, "the sender asks the same recovery the view reads")
	})

	t.Run("the caller's ledger is the one recovery coordinates on", func(t *testing.T) {
		f := newFixture(t)
		ledger := recordingplan.NewService(recordingplan.NewSQLRepository(f.store.DB()))
		p := New(Dependencies{Groups: f.groups, Tags: f.tags, Destinations: f.dest, Inspector: f.inspect, Store: f.store, Ledger: ledger},
			Config{NodeID: "node-1", Now: func() time.Time { return f.clock }})
		require.Same(t, ledger, p.recovery.ledger)
	})

	t.Run("a recovery of the caller's own stays in place", func(t *testing.T) {
		f := newFixture(t)
		own := &recordingRecovery{}
		p := New(Dependencies{Groups: f.groups, Tags: f.tags, Destinations: f.dest, Inspector: f.inspect, Store: f.store},
			Config{NodeID: "node-1", Now: func() time.Time { return f.clock }, Sender: groupdelivery.SenderConfig{Recovery: own}})
		require.Same(t, own, p.senderConfig().Recovery)
	})

	t.Run("without a store nothing is wired and nothing panics", func(t *testing.T) {
		f := newFixture(t)
		p := New(Dependencies{Groups: f.groups, Tags: f.tags, Destinations: f.dest, Inspector: f.inspect},
			Config{NodeID: "node-1", Now: func() time.Time { return f.clock }})
		require.Nil(t, p.recovery)
		require.Nil(t, p.senderConfig().Recovery)
	})

	t.Run("a wide group shows no storage view", func(t *testing.T) {
		f := newFixture(t)
		p := f.pipeline()
		view, err := p.Delivery(t.Context(), "group-G")
		require.NoError(t, err)
		require.Nil(t, view.Storage)
	})
}

// recordingRecovery is a StorageRecovery that stands in for a caller's own.
type recordingRecovery struct{ groupdelivery.StorageRecovery }
