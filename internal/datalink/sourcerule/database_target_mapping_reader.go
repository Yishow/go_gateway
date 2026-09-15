package sourcerule

import (
	"context"
	"sync"

	"go-gateway/internal/datalink/schema"
)

var (
	databaseTargetMappingReaders      sync.Map
	databaseTargetConnectorReaders    sync.Map
	databaseTargetConnectorValidators sync.Map
)

type DatabaseTargetMappingReader interface {
	List(ctx context.Context) ([]*schema.DatabaseTargetMapping, error)
}

type DatabaseTargetMappingListFunc func(ctx context.Context) ([]*schema.DatabaseTargetMapping, error)

func (fn DatabaseTargetMappingListFunc) List(ctx context.Context) ([]*schema.DatabaseTargetMapping, error) {
	return fn(ctx)
}

type DatabaseTargetConnectorReader interface {
	GetByID(ctx context.Context, connectorID string) (*schema.DatabaseConnector, error)
}

type DatabaseTargetConnectorGetFunc func(
	ctx context.Context,
	connectorID string,
) (*schema.DatabaseConnector, error)

func (fn DatabaseTargetConnectorGetFunc) GetByID(
	ctx context.Context,
	connectorID string,
) (*schema.DatabaseConnector, error) {
	return fn(ctx, connectorID)
}

type DatabaseTargetConnectorValidation struct {
	Ready  bool
	Issues []DatabaseTargetValidationIssue
}

type DatabaseTargetValidationIssue struct {
	Severity  string
	MappingID string
	Code      string
	Message   string
}

type DatabaseTargetConnectorValidator interface {
	Validate(ctx context.Context, connectorID string) (*DatabaseTargetConnectorValidation, error)
}

type DatabaseTargetConnectorValidatorFunc func(
	ctx context.Context,
	connectorID string,
) (*DatabaseTargetConnectorValidation, error)

func (fn DatabaseTargetConnectorValidatorFunc) Validate(
	ctx context.Context,
	connectorID string,
) (*DatabaseTargetConnectorValidation, error) {
	return fn(ctx, connectorID)
}

func (s *Service) SetDatabaseTargetMappingReader(reader DatabaseTargetMappingReader) {
	if reader == nil {
		databaseTargetMappingReaders.Delete(s)
		return
	}
	databaseTargetMappingReaders.Store(s, reader)
}

func (s *Service) databaseTargetMappingReader() DatabaseTargetMappingReader {
	reader, ok := databaseTargetMappingReaders.Load(s)
	if !ok {
		return nil
	}
	mappingReader, ok := reader.(DatabaseTargetMappingReader)
	if !ok {
		return nil
	}
	return mappingReader
}

func (s *Service) SetDatabaseTargetConnectorReader(reader DatabaseTargetConnectorReader) {
	if reader == nil {
		databaseTargetConnectorReaders.Delete(s)
		return
	}
	databaseTargetConnectorReaders.Store(s, reader)
}

func (s *Service) databaseTargetConnectorReader() DatabaseTargetConnectorReader {
	reader, ok := databaseTargetConnectorReaders.Load(s)
	if !ok {
		return nil
	}
	connectorReader, ok := reader.(DatabaseTargetConnectorReader)
	if !ok {
		return nil
	}
	return connectorReader
}

func (s *Service) SetDatabaseTargetConnectorValidator(validator DatabaseTargetConnectorValidator) {
	if validator == nil {
		databaseTargetConnectorValidators.Delete(s)
		return
	}
	databaseTargetConnectorValidators.Store(s, validator)
}

func (s *Service) databaseTargetConnectorValidator() DatabaseTargetConnectorValidator {
	validator, ok := databaseTargetConnectorValidators.Load(s)
	if !ok {
		return nil
	}
	connectorValidator, ok := validator.(DatabaseTargetConnectorValidator)
	if !ok {
		return nil
	}
	return connectorValidator
}
