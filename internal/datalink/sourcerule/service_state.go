package sourcerule

import (
	"context"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/tag"
)

// Service manages persisted source rules and their derived runtime projections.
type Service struct {
	repo             Repository
	deviceSvc        *device.Service
	pointSvc         *point.Service
	tagSvc           *tag.Service
	mappingSvc       *mapping.Service
	runtimeSync      RuntimeSyncer
	shareRuntimeSync ShareRuntimeReconciler
	candidateScope   candidateScope
}

type candidateScope struct {
	workspace CandidateWorkspaceReader
	hydration CandidateHydrationReader
}

// CandidateWorkspaceReader exposes the durable workspace/device membership
// needed to authorize candidate review and apply operations.
type CandidateWorkspaceReader interface {
	ShareOwnershipSnapshot(context.Context) (string, []string, error)
}

// CandidateHydrationReader exposes the backend Share hydration barrier.
type CandidateHydrationReader interface {
	CheckHydration(context.Context) (modbusshare.HydrationState, error)
}
