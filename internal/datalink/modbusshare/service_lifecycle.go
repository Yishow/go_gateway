package modbusshare

import (
	"context"
	"net"
	"strconv"
)

// StartCAS starts the listener while checking the caller's durable settings revision.
func (s *Service) StartCAS(ctx context.Context, port int, expectedSettingsRevision string) error {
	s.mu.RLock()
	running := s.state == shareStateRunning && s.server.Port() == port
	authoritative := s.authoritativeSettings
	s.mu.RUnlock()
	if running && !authoritative {
		return &Error{Code: ErrCodeListenerBindFailed, Message: "modbus share listener is already running", Retryable: false, Action: "stop the listener before starting it again"}
	}
	settings := s.Settings()
	if settings.BindAddress == "" {
		settings.BindAddress = defaultBindAddress
	}
	if settings.SlaveID == 0 {
		settings.SlaveID = 1
	}
	if settings.CapacityRegisters == 0 {
		settings.CapacityRegisters = s.bankCapacityRegisters()
	}
	settings.Enabled = true
	settings.Port = port
	return s.ApplySettingsCAS(ctx, settings, expectedSettingsRevision)
}

// StopCAS stops the listener while checking the caller's durable settings revision.
func (s *Service) StopCAS(ctx context.Context, expectedSettingsRevision string) error {
	settings := s.Settings()
	settings.Enabled = false
	return s.ApplySettingsCAS(ctx, settings, expectedSettingsRevision)
}

// Start starts local Modbus TCP server on specified port.
func (s *Service) Start(port int) error {
	s.mu.RLock()
	settings := s.settings
	running := s.state == shareStateRunning && s.server.Port() == port
	s.mu.RUnlock()
	if running {
		return &Error{Code: ErrCodeListenerBindFailed, Message: "modbus share listener is already running", Retryable: false, Action: "stop the listener before starting it again"}
	}
	if settings.BindAddress == "" {
		settings.BindAddress = "127.0.0.1"
	}
	if settings.SlaveID == 0 {
		settings.SlaveID = 1
	}
	if settings.CapacityRegisters == 0 {
		settings.CapacityRegisters = s.bankCapacityRegisters()
	}
	settings.Enabled = true
	settings.Port = port
	return s.ApplySettings(context.Background(), settings)
}

// Stop stops the local Modbus server.
func (s *Service) Stop() error {
	settings := s.Settings()
	settings.Enabled = false
	return s.ApplySettings(context.Background(), settings)
}

// Status describes runtime state of the local Modbus share server.
type Status struct {
	WorkspaceID              string              `json:"workspace_id,omitempty"`
	Enabled                  bool                `json:"enabled"`
	ConfiguredEnabled        bool                `json:"configured_enabled"`
	Port                     int                 `json:"port"`
	Address                  string              `json:"address"`
	BindAddress              string              `json:"bind_address"`
	BindState                string              `json:"bind_state"`
	LifecycleState           string              `json:"lifecycle_state"`
	ListenerState            string              `json:"listener_state"`
	MappingCount             int                 `json:"mapping_count"`
	SlaveID                  uint8               `json:"slave_id,omitempty"`
	CapacityRegisters        int                 `json:"capacity_registers,omitempty"`
	SettingsRevision         string              `json:"settings_revision,omitempty"`
	WorkspaceRevision        string              `json:"workspace_revision,omitempty"`
	HydrationState           string              `json:"hydration_state,omitempty"`
	Readiness                bool                `json:"readiness"`
	ReadinessToken           string              `json:"readiness_token,omitempty"`
	DirtyState               string              `json:"dirty_state,omitempty"`
	Recovery                 *RecoveryStatus     `json:"recovery,omitempty"`
	LastFailure              *Diagnostic         `json:"last_failure,omitempty"`
	CanonicalPlan            *CanonicalSharePlan `json:"canonical_plan,omitempty"`
	CanonicalDesiredMappings []DesiredMapping    `json:"canonical_desired_mappings,omitempty"`
}

// StatusForWorkspace scopes the in-memory mapping count to the requested
// workspace while retaining the lifecycle state from the same service.
func (s *Service) StatusForWorkspace(workspaceID string) Status {
	status := s.Status()
	status.WorkspaceID = workspaceID
	status.MappingCount = len(s.ListMappingsByWorkspace(workspaceID))
	return status
}

// Status returns runtime status.
func (s *Service) Status() Status {
	s.mu.RLock()
	defer s.mu.RUnlock()

	port := s.server.Port()
	if port == 0 && s.settings.Enabled && s.state == HydrationStateFailed {
		port = s.settings.Port
	}
	configuredEnabled := s.settings.Enabled
	enabled := configuredEnabled && s.server.Port() > 0
	hydrationReady := s.hydration.State == HydrationStateReady && s.hydration.Readiness
	addr := ""
	if enabled {
		addr = net.JoinHostPort(s.settings.BindAddress, strconv.Itoa(port))
	} else if s.settings.Enabled && port > 0 {
		addr = net.JoinHostPort(s.settings.BindAddress, strconv.Itoa(port))
	}
	bindState := shareStateDisabled
	//nolint:gocritic // Bind-state precedence mirrors the externally visible lifecycle states.
	if s.state == shareStateRunning && port > 0 {
		bindState = "pass"
	} else if s.state == HydrationStateFailed {
		bindState = "fail"
	} else if s.state == shareStateStopped {
		bindState = shareStateStopped
	} else if !s.settings.Enabled {
		bindState = "disabled"
	} else if port == 0 {
		bindState = "fail"
	}

	status := Status{
		WorkspaceID:       s.hydration.WorkspaceID,
		Enabled:           enabled,
		ConfiguredEnabled: configuredEnabled,
		Port:              port,
		Address:           addr,
		BindAddress:       s.settings.BindAddress,
		BindState:         bindState,
		LifecycleState:    s.state,
		ListenerState:     s.state,
		MappingCount:      len(s.mappings),
		SlaveID:           s.settings.SlaveID,
		CapacityRegisters: s.settings.CapacityRegisters,
		SettingsRevision:  s.settings.SettingsRevision,
		WorkspaceRevision: s.hydration.WorkspaceRevision,
		HydrationState:    s.hydration.State,
		Readiness:         s.hydration.Readiness,
		ReadinessToken:    s.hydration.ReadinessToken,
		LastFailure:       cloneDiagnostic(s.lastFailure),
	}
	if !hydrationReady {
		status.Port = 0
		status.Address = ""
		status.BindAddress = ""
		status.SlaveID = 0
		status.CapacityRegisters = 0
		status.SettingsRevision = ""
		status.WorkspaceRevision = ""
		status.ReadinessToken = ""
	}
	return status
}
