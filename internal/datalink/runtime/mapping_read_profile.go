package runtime

import (
	"context"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

func mappingProjectionReadProfile(projection *workspace.RuntimeProjection, p *schema.Point, contexts ...context.Context) (connector.ReadProfile, error) {
	if !modernPoint(p) {
		return connector.ReadProfile{}, nil
	}
	if len(contexts) > 0 {
		if profiles, ok := contexts[0].Value(readProfilesKey{}).(map[string]connector.ReadProfile); ok {
			if profile, found := profiles[p.ID]; found {
				return profile, connector.ValidatePointReadProfile(p, profile)
			}
			// Disabled/draft mapping metadata may exist without a scheduled install.
			// It cannot establish proof, but still needs the same immutable projection.
		}
	}
	return projectionReadProfile(projection, p, true)
}
