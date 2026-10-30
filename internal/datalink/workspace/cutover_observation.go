package workspace

import (
	"context"
	"encoding/json"
	"reflect"
	"regexp"

	"go-gateway/internal/datalink/dbtarget"
)

// CutoverObservation is an ephemeral confirmed operation reference. It carries
// original source identity, never a current draft, destination or mutation token.
type CutoverObservation struct {
	Available                                              bool
	WorkspaceID, WorkspaceRevision, GroupID, GroupRevision string
	ConnectorID, ConnectorRevision, OperationID            string
	Outcome                                                string // applied, rejected or unavailable; no free-form native error
}

// CutoverObservationConfig is immutable after service construction. Providers
// and callbacks must be nonblocking; they run after the native transaction and
// workspace mutex are released. ScopeRevision must use a pure existing read.
type CutoverObservationConfig struct {
	ScopeRevision            func(context.Context, string, string, string, string, string) (string, bool)
	Observe                  func(context.Context, CutoverObservation)
	RecordLimit              int
	OnDrop                   func(uint64)
	InfoEnabled, WarnEnabled bool
}

var cutoverObservationID = regexp.MustCompile(`^[A-Za-z0-9_-][A-Za-z0-9_.-]{0,127}$`)

const cutoverObservationUnavailable = "unavailable"

func (s *WriteGroupService) observeCutover(ctx context.Context, plan *cutoverPlan, result *WriteGroupCutover, nativeErr error, replayed bool) {
	if s == nil || replayed {
		return
	}
	config := s.cutoverObservation
	if config.Observe == nil {
		return
	}
	success := nativeErr == nil && result != nil && plan != nil
	if (success && !config.InfoEnabled) || (!success && !config.WarnEnabled) {
		return
	}
	// An ApplyCutover has at most one descriptor. There is no second queue or
	// background worker; an exhausted capture limit never changes native return.
	if config.RecordLimit < 1 {
		safeCutoverDrop(config.OnDrop)
		return
	}
	observation := CutoverObservation{Outcome: cutoverObservationUnavailable}
	if success {
		observation = s.proveCutoverObservation(ctx, plan, result, config)
	}
	if !observation.Available && !config.WarnEnabled {
		return
	}
	func() {
		defer func() { recover() }() //nolint:errcheck // Diagnostics must not change the settled cutover.
		config.Observe(ctx, observation)
	}()
}
func safeCutoverDrop(drop func(uint64)) {
	if drop != nil {
		defer func() { recover() }() //nolint:errcheck // Drop reporting must not change the cutover.
		drop(1)
	}
}
func cutoverScopeRevision(ctx context.Context, c CutoverObservationConfig, g *WriteGroup) (revision string, ok bool) {
	defer func() {
		if recover() != nil {
			revision, ok = "", false
		}
	}()
	if c.ScopeRevision == nil || g == nil || ctx == nil || ctx.Err() != nil {
		return "", false
	}
	revision, ok = c.ScopeRevision(ctx, g.WorkspaceID, g.ID, g.Revision, g.Destination.ConnectorID, g.Destination.ConnectorRevision)
	return revision, ok && cutoverObservationID.MatchString(revision)
}

func (s *WriteGroupService) proveCutoverObservation(ctx context.Context, plan *cutoverPlan, result *WriteGroupCutover, c CutoverObservationConfig) (observation CutoverObservation) {
	observation = CutoverObservation{Outcome: cutoverObservationUnavailable}
	// A provider panic or proof-read failure cannot change the already-settled
	// operation; neither may it establish a partially attributed success.
	defer func() {
		if recover() != nil {
			observation = CutoverObservation{Outcome: cutoverObservationUnavailable}
		}
	}()
	if ctx == nil || ctx.Err() != nil || s.repo == nil || result == nil || plan == nil || plan.source == nil || result.WorkspaceID != plan.req.WorkspaceID || result.SourceGroupID != plan.source.ID || result.Kind != plan.req.Kind {
		return observation
	}
	var originalRevision string
	switch result.Kind {
	case WriteGroupCutoverForward:
		if plan.full == nil || plan.full.ID != result.SourceGroupID {
			return observation
		}
		originalRevision = plan.full.Revision
		found := false
		for _, end := range result.Ended {
			if end.GroupID == result.SourceGroupID && end.Revision == originalRevision {
				found = true
			}
		}
		if !found {
			return observation
		}
	case WriteGroupCutoverRollback:
		if plan.prior == nil || plan.prior.Kind != WriteGroupCutoverForward || plan.prior.SourceGroupID != result.SourceGroupID || result.PriorCutoverID != plan.prior.ID || len(plan.prior.Ended) != 1 || plan.prior.Ended[0].GroupID != result.SourceGroupID {
			return observation
		}
		originalRevision = plan.prior.Ended[0].Revision
	default:
		return observation
	}
	// Identity comes from the actual original immutable version, not today's saved
	// source or the new rollback Started revision. No bootstrap getters are used.
	original, err := s.repo.versionByRevisionInTx(ctx, s.repo.db, result.WorkspaceID, result.SourceGroupID, originalRevision)
	if err != nil || original == nil || original.Group == nil || original.GroupRevision != originalRevision || original.AppliedRevision != originalRevision || original.Group.AppliedRevision != originalRevision {
		return observation
	}
	group := original.Group
	if result.Kind == WriteGroupCutoverForward && !reflect.DeepEqual(group, plan.full) {
		return observation
	}
	for _, id := range []string{result.ID, result.WorkspaceID, group.ID, group.Revision, group.Destination.ConnectorID, group.Destination.ConnectorRevision} {
		if !cutoverObservationID.MatchString(id) {
			return observation
		}
	}
	before, ok := cutoverScopeRevision(ctx, c, group)
	if !ok {
		return observation
	}
	frozen, ok := s.cutoverSettledObservationProof(ctx, plan, result, original)
	if !ok {
		return observation
	}
	after, ok := cutoverScopeRevision(ctx, c, group)
	if !ok || before != after {
		return observation
	}
	// Scope providers are outside native locks. Revalidate every immutable proof
	// after the final provider returns; unchanged WRevision alone cannot prove
	// that a ledger/version/descriptor was not changed while the read was in flight.
	finalFrozen, ok := s.cutoverSettledObservationProof(ctx, plan, result, original)
	if !ok || frozen != finalFrozen {
		return observation
	}
	return CutoverObservation{Available: true, WorkspaceID: group.WorkspaceID, WorkspaceRevision: after, GroupID: group.ID, GroupRevision: group.Revision, ConnectorID: group.Destination.ConnectorID, ConnectorRevision: group.Destination.ConnectorRevision, OperationID: result.ID, Outcome: "applied"}
}

func (s *WriteGroupService) cutoverSettledObservationProof(ctx context.Context, plan *cutoverPlan, result *WriteGroupCutover, original *WriteGroupAppliedSnapshot) (string, bool) {
	committed, err := s.repo.cutoverInTx(ctx, s.repo.db, result.WorkspaceID, result.ID)
	if err != nil || !reflect.DeepEqual(committed, result) {
		return "", false
	}
	if plan.prior != nil {
		prior, err := s.repo.cutoverInTx(ctx, s.repo.db, result.WorkspaceID, plan.prior.ID)
		if err != nil || !reflect.DeepEqual(prior, plan.prior) {
			return "", false
		}
	}
	group := original.Group
	reread, err := s.repo.versionByRevisionInTx(ctx, s.repo.db, group.WorkspaceID, group.ID, group.Revision)
	if err != nil || !reflect.DeepEqual(original, reread) {
		return "", false
	}
	return s.cutoverFrozenObservationProof(ctx, original)
}

func (s *WriteGroupService) cutoverFrozenObservationProof(ctx context.Context, original *WriteGroupAppliedSnapshot) (string, bool) {
	group := original.Group
	var count int
	if err := s.repo.db.QueryRowContext(ctx, s.repo.query(`SELECT COUNT(*) FROM wg_runtime_versions WHERE group_id=$1 AND group_revision=$2`), group.ID, group.Revision).Scan(&count); err != nil || count != 1 {
		return "", false
	}
	var payload string
	if err := s.repo.db.QueryRowContext(ctx, s.repo.query(`SELECT payload FROM wg_runtime_versions WHERE workspace_id=$1 AND group_id=$2 AND group_revision=$3`), group.WorkspaceID, group.ID, group.Revision).Scan(&payload); err != nil {
		return "", false
	}
	var frozen WriteGroupAppliedSnapshot
	if json.Unmarshal([]byte(payload), &frozen) != nil || frozen.RuntimeLayout == nil || frozen.WorkspaceID != group.WorkspaceID || frozen.GroupID != group.ID || frozen.GroupRevision != group.Revision || frozen.AppliedRevision != group.Revision || !frozen.EffectiveAt.Equal(original.EffectiveAt) || !reflect.DeepEqual(frozen.Group, group) {
		return "", false
	}
	switch frozen.RuntimeLayout.Dialect {
	case dbtarget.SQLDialectSQLite, dbtarget.SQLDialectPostgres:
		return payload, true
	case dbtarget.SQLDialectMySQL:
		return payload, ValidateMySQLRuntimeLayout(group, frozen.RuntimeLayout) == nil
	case dbtarget.SQLDialectSQLServer:
		return payload, ValidateSQLServerRuntimeLayout(group, frozen.RuntimeLayout) == nil
	default:
		return "", false
	}
}
