package sourcerule

import (
	"context"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/schema"
)

func (s *Service) loadDeviceActivationReadiness(ctx context.Context, deviceID string) (*schema.DeviceReadiness, error) {
	readiness, err := s.deviceSvc.CheckReadiness(ctx, deviceID)
	if err != nil {
		return nil, fmt.Errorf("檢查設備就緒狀態失敗: %w", err)
	}
	return readiness, nil
}

func (s *Service) initialRuleEnabled(ctx context.Context, requested bool, deviceID string) (bool, error) {
	if !requested {
		return false, nil
	}
	readiness, err := s.loadDeviceActivationReadiness(ctx, deviceID)
	if err != nil {
		return false, err
	}
	return readiness != nil && readiness.ActivationAllowed, nil
}

func (s *Service) ensureRuleActivationAllowed(ctx context.Context, deviceID string, enabled bool) error {
	if !enabled {
		return nil
	}
	readiness, err := s.loadDeviceActivationReadiness(ctx, deviceID)
	if err != nil {
		return err
	}
	if readiness != nil && readiness.ActivationAllowed {
		return nil
	}
	return ruleActivationBlockedError(readiness)
}

func (s *Service) runtimeRuleEnabled(ctx context.Context, rule *schema.SourceRule) (bool, error) {
	if !rule.Enabled {
		return false, nil
	}
	readiness, err := s.loadDeviceActivationReadiness(ctx, rule.DeviceID)
	if err != nil {
		return false, err
	}
	return readiness != nil && readiness.ActivationAllowed, nil
}

func ruleActivationBlockedError(readiness *schema.DeviceReadiness) error {
	message := "設備尚未達到 activation readiness，不能啟用來源規則"
	if readiness == nil || len(readiness.BlockingReasons) == 0 {
		return validationError(message)
	}
	return validationError(fmt.Sprintf("%s: %s", message, strings.Join(readiness.BlockingReasons, "; ")))
}
