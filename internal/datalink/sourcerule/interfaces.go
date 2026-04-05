package sourcerule

import (
	"context"

	"go-gateway/internal/datalink/schema"
)

type Repository interface {
	Create(ctx context.Context, rule *schema.SourceRule) error
	Update(ctx context.Context, rule *schema.SourceRule) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*schema.SourceRule, error)
	List(ctx context.Context, filter ListFilter) ([]*schema.SourceRule, error)
	CreateLinks(ctx context.Context, links []*schema.SourceRuleLink) error
	ListLinks(ctx context.Context, ruleID string) ([]*schema.SourceRuleLink, error)
	DeleteLinks(ctx context.Context, ruleID string) error
	ReplaceCandidateSnapshots(ctx context.Context, snapshots []*schema.SourceRuleCandidateSnapshot) error
	ListCandidateSnapshots(ctx context.Context, ruleID, revisionID string) ([]*schema.SourceRuleCandidateSnapshot, error)
	DeleteCandidateSnapshots(ctx context.Context, ruleID, revisionID string) error
}

type RuntimeSyncer interface {
	UpsertPoint(point *schema.Point)
	RemovePoint(pointID string)
}

type ListFilter struct {
	DeviceID *string
	Enabled  *bool
}
