package point

import (
	"context"
	"go-gateway/internal/datalink/schema"
)

type pointRevisionUpdater interface {
	UpdateRevision(context.Context, *schema.Point, string) error
}

func (s *Service) updateRevision(ctx context.Context, p *schema.Point, expected string) error {
	if repo, ok := s.repo.(pointRevisionUpdater); ok {
		return repo.UpdateRevision(ctx, p, expected)
	}
	return s.repo.Update(ctx, p)
}
