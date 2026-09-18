package workspace

import (
	"context"
	"fmt"
	"strings"
)

func (s *Service) BindDatabaseConnector(ctx context.Context, connectorID string) (*Record, error) {
	connectorID = strings.TrimSpace(connectorID)
	if connectorID == "" {
		return nil, workspaceValidationError("database connector id 不能為空")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	record, err := s.getOrCreate(ctx)
	if err != nil {
		return nil, err
	}

	record.DatabaseConnectorID = connectorID
	record.DatabaseSetupRevision = s.newID()
	record.UpdatedAt = s.now()
	if err := s.repo.Save(ctx, record); err != nil {
		return nil, fmt.Errorf("bind workspace database connector: %w", err)
	}

	return cloneRecord(record), nil
}
