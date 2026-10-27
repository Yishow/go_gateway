package workspace

import (
	"context"
	"slices"
	"time"
)

// ReviewCutover computes, without changing anything, what a cutover would do:
// which sources move, which stay, what ends and what starts at the boundary,
// and every reason it cannot be applied. The result is bound to everything it
// read by a digest that Apply must present.
func (s *WriteGroupService) ReviewCutover(ctx context.Context, req CutoverRequest) (*CutoverReview, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	req, err := normalizeCutoverRequest(req)
	if err != nil {
		return nil, err
	}
	plan, err := s.planCutover(ctx, req)
	if err != nil {
		return nil, err
	}
	return plan.review, nil
}

// cutoverPlan is a review together with what apply needs to carry it out.
type cutoverPlan struct {
	req    CutoverRequest
	review *CutoverReview
	// source is the saved group; full is the applied revision that ends, and
	// successor the saved draft that keeps the sources that are not moved.
	source, full, successor *WriteGroup
	ended                   []CutoverEnd
	started                 []CutoverStart
	targets                 []*WriteGroup
	readiness               map[string]*WriteGroupReadiness
	prior                   *WriteGroupCutover
	moved, remaining        []CutoverMemberRef
	facts                   []string
}

func (s *WriteGroupService) planCutover(ctx context.Context, req CutoverRequest) (*cutoverPlan, error) {
	if req.Kind == WriteGroupCutoverRollback {
		return s.planRollback(ctx, req)
	}
	return s.planForward(ctx, req)
}

// planBoundary fills the proposed boundary and reports the reasons an instant
// cannot be used. earliest is the start of the revision that would be cut.
func (s *WriteGroupService) planBoundary(req CutoverRequest, earliest time.Time) (boundary time.Time, issues []string) {
	now := s.clockNow()
	boundary = req.Boundary
	if boundary.IsZero() {
		// An unaligned instant on purpose: the cut windows are real, not avoided.
		boundary = now.Add(cutoverLead).Truncate(time.Millisecond)
	}
	if !boundary.After(now) {
		issues = append(issues, CutoverIssueBoundaryNotFuture)
	}
	if !earliest.IsZero() && !boundary.After(earliest) {
		issues = append(issues, CutoverIssueBoundaryBeforeActive)
	}
	return boundary, issues
}

func uniqueSorted(values []string) []string {
	out := slices.Clone(values)
	slices.Sort(out)
	return slices.Compact(out)
}

// finish turns the gathered facts and issues into the review and its digest.
func (p *cutoverPlan) finish(boundary time.Time, issues []string, source CutoverGroupView, targets []CutoverGroupView, unselected []CutoverMemberRef) (*cutoverPlan, error) {
	req := p.req
	req.Boundary = boundary
	facts := slices.Clone(p.facts)
	for _, ready := range p.readiness {
		for _, proof := range ready.counterBindings {
			facts = append(facts, "counter-proof:"+proof.Digest)
		}
	}
	digest, err := cutoverDigest(req, p.moved, p.remaining, facts, boundary)
	if err != nil {
		return nil, err
	}
	if issues == nil {
		issues = []string{}
	}
	p.req = req
	p.review = &CutoverReview{
		Kind: req.Kind, Boundary: boundary, Digest: digest, Applicable: len(issues) == 0, Issues: issues,
		Source: source, Targets: targets, Moved: append([]CutoverMemberRef{}, p.moved...),
		Unselected: append([]CutoverMemberRef{}, unselected...), Remaining: append([]CutoverMemberRef{}, p.remaining...),
	}
	return p, nil
}
