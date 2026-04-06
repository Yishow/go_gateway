package sourcerule

import (
	"context"
	"sync"

	"go-gateway/internal/datalink/schema"
)

var (
	databaseTargetMappingReaders      sync.Map
	databaseTargetConnectorValidators sync.Map
)

type DatabaseTargetMappingReader interface {
	List(ctx context.Context) ([]*schema.DatabaseTargetMapping, error)
}

type DatabaseTargetMappingListFunc func(ctx context.Context) ([]*schema.DatabaseTargetMapping, error)

func (fn DatabaseTargetMappingListFunc) List(ctx context.Context) ([]*schema.DatabaseTargetMapping, error) {
	return fn(ctx)
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
	mappingReader, _ := reader.(DatabaseTargetMappingReader)
	return mappingReader
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
	connectorValidator, _ := validator.(DatabaseTargetConnectorValidator)
	return connectorValidator
}
