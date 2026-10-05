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
	Attention              []groupdelivery.AttentionItem  `json:"attention"`
	GroupID                string                         `json:"group_id"`
	Intake                 IntakeView                     `json:"intake"`
	Stages                 StagesView                     `json:"stages"`
	LastSQLCommittedAt     *time.Time                     `json:"last_sql_committed_at"`
	LastSQLCommittedEffect *groupdelivery.CommittedEffect `json:"last_sql_committed_effect"`
	RevisionStages         *RevisionStagesView            `json:"revision_stages"`
	OldestPendingSeconds   float64                        `json:"oldest_pending_seconds"`
	NoDataBuckets          int                            `json:"no_data_buckets"`
	SkippedBuckets         int                            `json:"skipped_buckets"`
	RecentBucketIssues     []groupdelivery.BucketIssue    `json:"recent_bucket_issues"`
	Backlog                []BacklogView                  `json:"backlog"`
	Quota                  groupdelivery.QuotaStatusView  `json:"quota"`
}

// RevisionStagesView reports acceptance and delivery for one applied revision.
type RevisionStagesView struct {
	GroupRevision string     `json:"group_revision"`
	Stages        StagesView `json:"stages"`
}

func deliveryStages(status groupdelivery.StageCounts) StagesView {
	return StagesView{
		Collecting: status.Collecting, Queued: status.Queued, Retrying: status.Retrying,
		Blocked: status.Blocked, Quarantined: status.Quarantined, Unknown: status.Unknown,
		SQLCommitted: status.SQLCommitted, Skipped: status.Skipped,
	}
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
	attention, err := p.deps.Store.Attention(ctx, groupID)
	if err != nil {
		return nil, err
	}
	view := &DeliveryView{
		Attention:          attention,
		GroupID:            groupID,
		Intake:             IntakeView{State: StateNotRunning},
		Stages:             deliveryStages(status.Stages),
		LastSQLCommittedAt: status.LastSQLCommittedAt, OldestPendingSeconds: status.OldestPendingSeconds,
		NoDataBuckets: status.NoDataBuckets, SkippedBuckets: status.SkippedBuckets,
		RecentBucketIssues: status.RecentBucketIssues,
		Backlog:            make([]BacklogView, 0, len(status.Backlog)), Quota: quota.View(),
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

// DeliveryForRevision adds actual receipt identity for the current applied
// revision while retaining historical counts and backlog as delivery truth.
func (p *Pipeline) DeliveryForRevision(ctx context.Context, groupID, appliedRevision string) (*DeliveryView, error) {
	// Read the effect before totals so a concurrently committed first row
	// cannot produce identity with an older zero-commit count.
	effect, err := p.deps.Store.LastCommittedEffect(ctx, groupID, appliedRevision)
	if err != nil {
		return nil, err
	}
	stages, err := p.deps.Store.RevisionStageCounts(ctx, groupID, appliedRevision)
	if err != nil {
		return nil, err
	}
	view, err := p.Delivery(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if appliedRevision != "" {
		view.RevisionStages = &RevisionStagesView{GroupRevision: appliedRevision, Stages: deliveryStages(stages)}
	}
	if stages.SQLCommitted > 0 {
		view.LastSQLCommittedEffect = effect
	}
	return view, nil
}
