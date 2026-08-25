package modbusshare

import (
	"context"
	"encoding/json"
	"fmt"
)

// LoadPersistedSettings hydrates settings from system_settings without binding
// the listener before workspace projection restore has completed.
func (s *Service) LoadPersistedSettings(ctx context.Context) error {
	s.mu.RLock()
	repo := s.settingsRepo
	s.mu.RUnlock()
	if repo == nil {
		return fmt.Errorf("modbus share settings repository is not configured")
	}
	item, err := repo.Get(ctx, settingsKey)
	if err != nil {
		return fmt.Errorf("load modbus share settings: %w", err)
	}
	settings, err := decodeSettings(item.Value)
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.settings = settings
	if settings.Enabled {
		s.state = shareStateStopped
	} else {
		s.state = shareStateDisabled
	}
	s.authoritativeSettings = true
	s.mu.Unlock()
	return nil
}

func decodeSettings(value interface{}) (Settings, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return Settings{}, fmt.Errorf("encode persisted modbus share settings: %w", err)
	}
	var settings Settings
	if err := json.Unmarshal(payload, &settings); err != nil {
		return Settings{}, fmt.Errorf("decode persisted modbus share settings: %w", err)
	}
	if settings.BindAddress == "" {
		return Settings{}, fmt.Errorf("persisted modbus share bind address is required")
	}
	return settings, nil
}
