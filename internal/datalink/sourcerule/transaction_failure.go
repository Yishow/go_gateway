package sourcerule

import (
	"context"
	"errors"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/sqltx"
)

func (s *Service) noteTransactionFailure(ctx context.Context, err error) bool {
	var dirty *modbusshare.Error
	typedDirty := errors.As(err, &dirty) && dirty.Code == modbusshare.ErrCodeDirtyUnknown
	if !typedDirty && !errors.Is(err, sqltx.ErrCommitOutcomeUnknown) && !errors.Is(err, sqltx.ErrConnectionUnusable) {
		return false
	}
	if j := s.pointJournal(ctx); j != nil && j.storageFailure == nil {
		j.storageFailure = err
	}
	return true
}
func (s *Service) storageCompensationError(ctx context.Context) error {
	if j := s.pointJournal(ctx); j != nil {
		return j.storageFailure
	}
	return nil
}

type sourceTransactionUncertainError struct{ state, cause error }

func (e *sourceTransactionUncertainError) Error() string   { return e.state.Error() }
func (e *sourceTransactionUncertainError) Unwrap() []error { return []error{e.state, e.cause} }
