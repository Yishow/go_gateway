package collector

import "go-gateway/internal/datalink/schema"

// PrepareSelectedAcquisitionGeometry uses the full original SQL definition but
// journals only this operation's actual installs. Retained slots require their
// original matching authorization; a selected install cannot bless unknown slots.
func (s *Scheduler) PrepareSelectedAcquisitionGeometry(mutation AcquisitionGeometryMutation, device schema.Device, definition AcquisitionGeometryDefinition, selected map[string]bool) *AcquisitionGeometryUpdate {
	token := mutation.token
	if token == nil || token.owner != s {
		return nil
	}
	for id := range token.before.slots {
		if selected[id] {
			continue
		}
		before := token.before.publication
		if !token.before.authorized || before == nil || before.workspaceID != definition.WorkspaceID || before.sourceRevision != definition.SourceRevision {
			definition = AcquisitionGeometryDefinition{}
			break
		}
	}
	return s.PrepareAcquisitionGeometry(mutation, device, definition, "")
}
