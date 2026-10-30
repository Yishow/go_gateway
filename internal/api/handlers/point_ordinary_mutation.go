package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sqltx"
)

type ordinaryPointRuntime interface {
	CaptureOrdinaryPoint(string) (collector.PointMutationLease, bool)
	ApplyOrdinaryPoint(*schema.Point, collector.PointMutationLease) bool
	InvalidateOrdinaryPoint(collector.PointMutationLease) bool
}

func (h *PointHandler) updateOrdinaryPoint(c *gin.Context, id string, req point.UpdatePointRequest) (*schema.Point, error) {
	runtime, guarded := h.runtimeSync.(ordinaryPointRuntime)
	if !guarded {
		return h.svc.Update(c.Request.Context(), id, req)
	}
	before, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		return nil, err
	}
	proof, ok := runtime.CaptureOrdinaryPoint(id)
	if !ok {
		return nil, point.ErrFormatUnavailable
	}
	next, err := h.svc.UpdateOrdinary(c.Request.Context(), before, req, proof.Lease != 0 && proof.Revision == before.RevisionID)
	if err != nil {
		if errors.Is(err, sqltx.ErrCommitOutcomeUnknown) || errors.Is(err, sqltx.ErrConnectionUnusable) {
			runtime.InvalidateOrdinaryPoint(proof)
			return nil, errors.Join(point.ErrFormatUnavailable, err)
		}
		return nil, err
	}
	if !runtime.ApplyOrdinaryPoint(next, proof) {
		return nil, &ordinaryRuntimePendingError{point: next}
	}
	return next, nil
}

type ordinaryRuntimePendingError struct{ point *schema.Point }

func (e *ordinaryRuntimePendingError) Error() string {
	return "point configuration saved; runtime application pending"
}
func (e *ordinaryRuntimePendingError) Unwrap() error { return point.ErrFormatUnavailable }
