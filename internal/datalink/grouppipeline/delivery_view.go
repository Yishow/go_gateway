package grouppipeline

import (
	"context"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
)

// DeliveryView is the API-safe delivery truth of one group: its intake state,
// where every accepted row stands, and the backlog frozen to each revision. It
// is built from the durable store, never from buffers or ACKs, and contains no
// values, DSNs or driver text.
type DeliveryView struct {
	GroupID              string                        `json:"group_id"`
	Intake               IntakeView                    `json:"intake"`
	Stages               StagesView                    `json:"stages"`
	LastSQLCommittedAt   *time.Time                    `json:"last_sql_committed_at"`
	OldestPendingSeconds float64                       `json:"oldest_pending_seconds"`
	NoDataBuckets        int                           `json:"no_data_buckets"`
	SkippedBuckets       int                           `json:"skipped_buckets"`
	Backlog              []BacklogView                 `json:"backlog"`
	Quota                groupdelivery.QuotaStatusView `json:"quota"`
}

// IntakeView says whether the group is accepting new samples.
type IntakeView struct {
	State  string `json:"state"`
	Reason string `json:"reason,omitempty"`
}

// StagesView counts rows per stage; only SQLCommitted is destination evidence.
type StagesView struct {
	Collecting   int `json:"collecting"`
	Queued       int `json:"queued"`
	Retrying     int `json:"retrying"`
	Blocked      int `json:"blocked"`
	Quarantined  int `json:"quarantined"`
	Unknown      int `json:"unknown"`
	SQLCommitted int `json:"sql_committed"`
	Skipped      int `json:"skipped"`
}

// BacklogView is the undelivered backlog of one group revision and destination.
type BacklogView struct {
	GroupRevision     string   `json:"group_revision"`
	ConnectorID       string   `json:"connector_id"`
	ConnectorRevision string   `json:"connector_revision"`
	TableSchema       string   `json:"table_schema"`
	TableName         string   `json:"table_name"`
	Pending           int      `json:"pending"`
	ErrorCodes        []string `json:"error_codes"`
}

// StateNotRunning is the intake state of a group the pipeline does not know.
const StateNotRunning = "not_running"

// Delivery reads the delivery truth of one group.
func (p *Pipeline) Delivery(ctx context.Context, groupID string) (*DeliveryView, error) {
	status, err := p.deps.Store.GroupStatus(ctx, groupID)
	if err != nil {
		return nil, err
	}
	quota, err := p.deps.Store.QuotaStatus(ctx, groupdelivery.GroupKey{GroupID: groupID})
	if err != nil {
		return nil, err
	}
	view := &DeliveryView{
		GroupID: groupID,
		Intake:  IntakeView{State: StateNotRunning},
		Stages: StagesView{
			Collecting: status.Stages.Collecting, Queued: status.Stages.Queued, Retrying: status.Stages.Retrying,
			Blocked: status.Stages.Blocked, Quarantined: status.Stages.Quarantined, Unknown: status.Stages.Unknown,
			SQLCommitted: status.Stages.SQLCommitted, Skipped: status.Stages.Skipped,
		},
		LastSQLCommittedAt: status.LastSQLCommittedAt, OldestPendingSeconds: status.OldestPendingSeconds,
		NoDataBuckets: status.NoDataBuckets, SkippedBuckets: status.SkippedBuckets,
		Backlog: make([]BacklogView, 0, len(status.Backlog)), Quota: quota.View(),
	}
	for _, entry := range status.Backlog {
		codes := entry.LastErrorCodes
		if codes == nil {
			codes = []string{}
		}
		view.Backlog = append(view.Backlog, BacklogView{
			GroupRevision: entry.GroupRevision, ConnectorID: entry.ConnectorID, ConnectorRevision: entry.ConnectorRevision,
			TableSchema: entry.TableSchema, TableName: entry.TableName, Pending: entry.Pending, ErrorCodes: codes,
		})
	}
	for _, intake := range p.Status() {
		if intake.GroupID != groupID {
			continue
		}
		// An active boundary outranks a blocked newer revision for the same group.
		if view.Intake.State == StateNotRunning || intake.State == StateActive {
			view.Intake = IntakeView{State: intake.State, Reason: intake.Reason}
		}
	}
	return view, nil
}
