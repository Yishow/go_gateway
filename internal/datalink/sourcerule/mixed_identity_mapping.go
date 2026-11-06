package sourcerule

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-gateway/internal/datalink/schema"
)

// identitySyncLinks retains the previously-owned address while synchronizing
// existing mappings. Their address-based candidate identities move afterwards.
func identitySyncLinks(previous, current []*schema.SourceRuleLink) []*schema.SourceRuleLink {
	old := make(map[string]*schema.SourceRuleLink, len(previous))
	for _, link := range previous {
		old[link.ID] = link
	}
	links := cloneSourceRuleLinks(current)
	for _, link := range links {
		if prior := old[link.ID]; prior != nil && prior.PointID == link.PointID {
			link.Address = prior.Address
		}
	}
	return links
}

// rekeyMixedIdentities moves only metadata still owned by the same point and
// item. Compare-and-swap protects an intervening operator edit.
func (s *Service) rekeyMixedIdentities(ctx context.Context, rule *schema.SourceRule, previous, current []*schema.SourceRuleLink) (tagMappingSyncResult, error) {
	var result tagMappingSyncResult
	old := make(map[string]*schema.SourceRuleLink, len(previous))
	for _, link := range previous {
		old[link.ID] = link
	}
	seenTags, seenMaps := map[string]bool{}, map[string]bool{}
	for _, link := range current {
		prior := old[link.ID]
		if prior == nil || prior.PointID != link.PointID || normalizeAddressKey(prior.Address) == normalizeAddressKey(link.Address) {
			continue
		}
		if link.TagID != nil && *link.TagID != "" {
			if s.tagSvc == nil || seenTags[*link.TagID] {
				return result, fmt.Errorf("mixed item tag ownership unavailable")
			}
			seenTags[*link.TagID] = true
			tagRecord, err := s.tagSvc.GetByID(ctx, *link.TagID)
			if err != nil {
				return result, err
			}
			if isRuleManagedTag(tagRecord, rule.ID, prior.Address) {
				var labels map[string]string
				if err := json.Unmarshal([]byte(tagRecord.Labels), &labels); err != nil {
					return result, err
				}
				labels[ruleManagedTagLabelAddress] = normalizeAddressKey(link.Address)
				encoded, err := json.Marshal(labels)
				if err != nil {
					return result, err
				}
				after := *tagRecord
				after.Labels, after.UpdatedAt = string(encoded), time.Now()
				if err := s.tagSvc.RestoreWorkspaceSave(ctx, tagRecord, &after); err != nil {
					return result, &identityMoveUncertainError{cause: err}
				}
				result.tagIdentityMoves = append(result.tagIdentityMoves, tagIdentityMove{tagRecord, &after})
			}
		}
		if link.MappingID == nil || *link.MappingID == "" {
			continue
		}
		if s.mappingSvc == nil || seenMaps[*link.MappingID] {
			return result, fmt.Errorf("mixed item mapping ownership unavailable")
		}
		seenMaps[*link.MappingID] = true
		mappingRecord, err := s.mappingSvc.GetByID(ctx, *link.MappingID)
		if err != nil {
			return result, err
		}
		if mappingRecord.PointID != link.PointID || (link.TagID != nil && mappingRecord.TagID != *link.TagID) {
			return result, fmt.Errorf("mixed item mapping ownership changed")
		}
		if mappingRecord.RuleCandidateID == "" {
			continue
		}
		oldID, err := candidateID(buildMappingCandidateIdentity(rule.ID, prior.Address, ""))
		if err != nil {
			return result, err
		}
		if mappingRecord.RuleCandidateID != oldID {
			return result, fmt.Errorf("mixed item mapping candidate ownership changed")
		}
		newID, err := candidateID(buildMappingCandidateIdentity(rule.ID, link.Address, ""))
		if err != nil {
			return result, err
		}
		after := *mappingRecord
		after.RuleCandidateID, after.UpdatedAt = newID, time.Now()
		if err := s.mappingSvc.RestoreWorkspaceSave(ctx, mappingRecord, &after); err != nil {
			return result, &identityMoveUncertainError{cause: err}
		}
		result.mapIdentityMoves = append(result.mapIdentityMoves, mappingIdentityMove{mappingRecord, &after})
	}
	return result, nil
}

type identityMoveUncertainError struct{ cause error }

func (e *identityMoveUncertainError) Error() string {
	return "mixed item ownership could not be updated safely: " + e.cause.Error()
}
func (e *identityMoveUncertainError) Unwrap() error { return e.cause }

func (s *Service) rollbackIdentityMoves(ctx context.Context, result tagMappingSyncResult) error {
	for i := len(result.mapIdentityMoves) - 1; i >= 0; i-- {
		move := result.mapIdentityMoves[i]
		if err := s.mappingSvc.RestoreWorkspaceSave(ctx, move.after, move.before); err != nil {
			return err
		}
	}
	for i := len(result.tagIdentityMoves) - 1; i >= 0; i-- {
		move := result.tagIdentityMoves[i]
		if err := s.tagSvc.RestoreWorkspaceSave(ctx, move.after, move.before); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) noteRestoredIdentityMoves(ctx context.Context, result tagMappingSyncResult) {
	if len(result.tagIdentityMoves)+len(result.mapIdentityMoves) == 0 {
		return
	}
	journal := s.pointJournal(ctx)
	if journal == nil {
		return
	}
	if journal.restoredIdentityTags == nil {
		journal.restoredIdentityTags = map[string]bool{}
		journal.restoredIdentityMaps = map[string]bool{}
	}
	for _, move := range result.tagIdentityMoves {
		journal.restoredIdentityTags[move.before.ID] = true
	}
	for _, move := range result.mapIdentityMoves {
		journal.restoredIdentityMaps[move.before.ID] = true
	}
	for id := range result.updatedTags {
		journal.restoredIdentityTags[id] = true
	}
	for id := range result.updatedMappings {
		journal.restoredIdentityMaps[id] = true
	}
}
