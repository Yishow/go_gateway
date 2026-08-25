package modbusshare

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"

	datalinksettings "go-gateway/internal/datalink/settings"
	virtualmodbus "go-gateway/internal/virtual/server/modbus"
)

// ApplySettingsCAS persists settings only when expectedRevision matches the
// server's hydrated revision. A successful save always gets a server revision.
func (s *Service) ApplySettingsCAS(ctx context.Context, settings Settings, expectedRevision string) error {
	s.lifecycleMu.Lock()
	defer s.lifecycleMu.Unlock()
	s.mu.RLock()
	current := s.settings
	authoritative := s.authoritativeSettings
	s.mu.RUnlock()
	if authoritative && expectedRevision == "" {
		return &Error{Code: ErrCodeRevisionConflict, Message: "settings revision is required", Retryable: true, SettingsRevision: current.SettingsRevision}
	}
	if authoritative && expectedRevision != current.SettingsRevision {
		return &Error{Code: ErrCodeRevisionConflict, Message: settingsRevisionConflictMessage, Retryable: true, SettingsRevision: current.SettingsRevision}
	}
	if err := s.validateSettings(settings); err != nil {
		return err
	}
	settings.SettingsRevision = uuid.NewString()
	if repo, ok := s.settingsRepo.(datalinksettings.RevisionCASRepository); ok {
		casExpected := expectedRevision
		if !authoritative {
			casExpected = ""
		}
		if err := repo.SetIfRevision(ctx, settingsKey, casExpected, settings); err != nil {
			if errors.Is(err, datalinksettings.ErrRevisionConflict) {
				return &Error{Code: ErrCodeRevisionConflict, Message: "settings revision conflict", Retryable: true, SettingsRevision: current.SettingsRevision}
			}
			return &Error{Code: ErrCodeSettingsUpdateFailed, Message: "settings storage update failed", Retryable: true, SettingsRevision: current.SettingsRevision}
		}
		return s.applySettingsLocked(ctx, settings, false)
	}
	return s.applySettingsLocked(ctx, settings, true)
}

// ApplySettings applies new global Modbus Share settings and manages listener lifecycle accordingly.
func (s *Service) ApplySettings(ctx context.Context, settings Settings) error {
	s.lifecycleMu.Lock()
	defer s.lifecycleMu.Unlock()
	return s.applySettingsLocked(ctx, settings, true)
}

// StartConfiguredListener binds the persisted listener only after hydration
// has produced a ready runtime projection.
func (s *Service) StartConfiguredListener(ctx context.Context) error {
	s.lifecycleMu.Lock()
	defer s.lifecycleMu.Unlock()
	s.mu.RLock()
	settings := s.settings
	s.mu.RUnlock()
	if !settings.Enabled {
		return nil
	}
	hydration, err := s.CheckHydration(ctx)
	if err != nil || hydration.State != HydrationStateReady || !hydration.Readiness {
		s.mu.Lock()
		s.state = HydrationStateFailed
		s.mu.Unlock()
		return &Error{Code: ErrCodeHydrationRequired, Message: "share workspace hydration is not ready", Retryable: true, Action: "restore the workspace projection before starting the listener"}
	}
	return s.applySettingsLocked(ctx, settings, false)
}

// StopListener stops only the process-local listener. It intentionally keeps
// the durable enabled setting unchanged so a normal process shutdown restores
// the listener on the next startup.
func (s *Service) StopListener() error {
	s.lifecycleMu.Lock()
	defer s.lifecycleMu.Unlock()
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.server.Stop(); err != nil {
		return fmt.Errorf("stop Modbus Share listener: %w", err)
	}
	if s.settings.Enabled {
		s.state = shareStateStopped
	} else {
		s.state = shareStateDisabled
	}
	return nil
}

// CloseRuntime releases process-local Share resources without changing the
// durable operator setting.
func (s *Service) CloseRuntime() error { return s.StopListener() }

func (s *Service) applySettingsLocked(ctx context.Context, settings Settings, persist bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lifecycleEpoch++
	if !settings.Enabled {
		// A disable request may carry only the enabled flag. Keep the last
		// persisted endpoint/geometry instead of inventing zero values.
		if settings.BindAddress == "" {
			settings.BindAddress = s.settings.BindAddress
		}
		if settings.Port == 0 {
			settings.Port = s.settings.Port
		}
		if settings.SlaveID == 0 {
			settings.SlaveID = s.settings.SlaveID
		}
		if settings.CapacityRegisters == 0 {
			settings.CapacityRegisters = s.settings.CapacityRegisters
		}
		if settings.SettingsRevision == "" {
			settings.SettingsRevision = s.settings.SettingsRevision
		}
	}
	if err := s.validateSettings(settings); err != nil {
		return err
	}
	if settings.CapacityRegisters == 0 {
		settings.CapacityRegisters = s.bankCapacityRegisters()
	}
	if settings.SettingsRevision == "" {
		settings.SettingsRevision = fmt.Sprintf("settings-%d", time.Now().UTC().UnixNano())
	}
	if persist && s.settingsRepo != nil {
		if err := s.settingsRepo.Set(ctx, settingsKey, settings); err != nil {
			return fmt.Errorf("persist modbus share settings: %w", err)
		}
	}

	s.settings = settings
	s.authoritativeSettings = s.settingsRepo != nil || s.authoritativeSettings
	if settings.Enabled && s.authoritativeSettings && (s.hydration.State != HydrationStateReady || !s.hydration.Readiness) {
		// Durable settings may be saved before workspace hydration completes,
		// but binding is a runtime side effect and must wait for the restored
		// projection to become ready.
		if err := s.server.Stop(); err != nil {
			s.state = HydrationStateFailed
			return fmt.Errorf("stop Modbus Share listener: %w", err)
		}
		s.state = shareStateStopped
		return nil
	}
	if !settings.Enabled {
		if err := s.server.Stop(); err != nil {
			s.state = HydrationStateFailed
			return fmt.Errorf("stop Modbus Share listener: %w", err)
		}
		s.state = "disabled"
		s.lastFailure = nil
		return nil
	}

	port := settings.Port

	// If already running with the same durable endpoint and geometry, keep it.
	if s.server.Port() == port && s.server.BindAddress() == settings.BindAddress &&
		s.server.SlaveID() == settings.SlaveID && s.server.CapacityRegisters() == settings.CapacityRegisters && s.state == shareStateRunning {
		s.lastFailure = nil
		return nil
	}

	// Stop previous listener if any
	if s.server.Port() > 0 {
		s.state = "stopping"
	}
	if err := s.server.Stop(); err != nil {
		s.state = HydrationStateFailed
		return fmt.Errorf("stop Modbus Share listener: %w", err)
	}
	s.state = "starting"

	if err := preflightPortAvailable(ctx, settings.BindAddress, port); err != nil {
		s.state = HydrationStateFailed
		s.lastFailure = listenerBindFailureDiagnostic(port, "port is already in use or unavailable")
		return &Error{
			Code:      ErrCodeListenerBindFailed,
			Message:   fmt.Sprintf("port %d is already in use or unavailable", port),
			Retryable: true,
			Action:    "release port or configure a different port in settings",
		}
	}

	if err := s.server.StartWithConfig(ctx, virtualServerConfig(settings)); err != nil {
		s.state = HydrationStateFailed
		s.lastFailure = listenerBindFailureDiagnostic(port, "listener could not be bound")
		return &Error{
			Code:      ErrCodeListenerBindFailed,
			Message:   fmt.Sprintf("failed to bind Modbus Share listener on port %d", port),
			Retryable: true,
			Action:    "check network permissions and configure a free port",
		}
	}

	s.state = shareStateRunning
	s.lastFailure = nil
	return nil
}

func listenerBindFailureDiagnostic(port int, reason string) *Diagnostic {
	return &Diagnostic{
		Code:      ErrCodeListenerBindFailed,
		Severity:  "error",
		Message:   fmt.Sprintf("Modbus Share listener could not bind port %d: %s", port, reason),
		Retryable: true,
		Action:    "release the port or configure a different port in settings",
	}
}

func (s *Service) validateSettings(settings Settings) error {
	if settings.CapacityRegisters > s.bankCapacityRegisters() {
		return NewError(ErrCodeCapacityExceeded, fmt.Sprintf("capacity must be between 1 and %d registers", s.bankCapacityRegisters()), false)
	}
	if !settings.Enabled {
		return nil
	}
	if settings.BindAddress == "" {
		return NewError(ErrCodeListenerBindFailed, "bind address is required", false)
	}
	if net.ParseIP(settings.BindAddress) == nil {
		return NewError(ErrCodeListenerBindFailed, "bind address is invalid", false)
	}
	if settings.Port < 1 || settings.Port > 65535 {
		return NewError(ErrCodeListenerBindFailed, "port must be between 1 and 65535", false)
	}
	if settings.SlaveID < 1 || settings.SlaveID > 247 {
		return NewError(ErrCodeListenerBindFailed, "slave id must be between 1 and 247", false)
	}
	if settings.CapacityRegisters < 1 {
		return NewError(ErrCodeCapacityExceeded, fmt.Sprintf("capacity must be between 1 and %d registers", s.bankCapacityRegisters()), false)
	}
	return nil
}

func virtualServerConfig(settings Settings) virtualmodbus.Config {
	return virtualmodbus.Config{
		BindAddress:       settings.BindAddress,
		Port:              settings.Port,
		SlaveID:           settings.SlaveID,
		CapacityRegisters: settings.CapacityRegisters,
	}
}
