package device

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
)

type readProbeTarget struct {
	enabled bool
	req     connector.ReadRequest
	details string
}

// ProbeAndActivate 在啟用設備前先完成 connect + read probe 驗證。
func (s *Service) ProbeAndActivate(ctx context.Context, id string) error {
	device, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("取得設備失敗: %w", err)
	}

	if device.Status == schema.DeviceStatusActive {
		return nil
	}

	result, err := s.TestConnectionWithResult(ctx, id)
	if err != nil {
		return fmt.Errorf("連線測試失敗，無法啟用設備: %w", err)
	}
	if result == nil {
		return fmt.Errorf("連線測試失敗，無法啟用設備")
	}
	if !result.CanActivate {
		if result.Probe.Status == TestConnectionStageFailed && result.Probe.Error != "" {
			return fmt.Errorf("讀取探測失敗，無法啟用設備: %s", result.Probe.Error)
		}
		if result.Error != "" {
			return fmt.Errorf("連線測試失敗，無法啟用設備: %s", result.Error)
		}
		return fmt.Errorf("連線測試失敗，無法啟用設備")
	}

	if err := s.repo.UpdateStatus(ctx, id, schema.DeviceStatusActive); err != nil {
		return fmt.Errorf("更新設備狀態失敗: %w", err)
	}

	return nil
}

func probeReadWithProtocol(
	ctx context.Context,
	protocol connector.Protocol,
	device *schema.Device,
) (bool, error) {
	target, err := buildReadProbeTarget(device)
	if err != nil {
		return false, err
	}
	if !target.enabled {
		return false, nil
	}

	result, err := protocol.Read(ctx, target.req)
	if err != nil {
		return true, fmt.Errorf("%s, read 失敗: %w", target.details, err)
	}
	if result.Error != "" {
		return true, fmt.Errorf("%s, read 錯誤: %s", target.details, result.Error)
	}
	if result.Quality == schema.QualityBad {
		return true, fmt.Errorf("%s, read 品質為 bad", target.details)
	}

	return true, nil
}

func buildReadProbeTarget(device *schema.Device) (readProbeTarget, error) {
	overrideAddr, overrideFunc := extractProbeOverride(device.ConnectionConfig)

	switch device.Protocol {
	case schema.ProtocolModbusTCP:
		var cfg schema.ConnectionConfigModbusTCP
		if err := json.Unmarshal([]byte(device.ConnectionConfig), &cfg); err != nil {
			return readProbeTarget{}, fmt.Errorf("解析連線設定失敗: %w", err)
		}
		addr := pickProbeValue(overrideAddr, "40001")
		function := pickProbeValue(overrideFunc, "03")
		return readProbeTarget{
			enabled: true,
			req: connector.ReadRequest{
				Address:  addr,
				Function: function,
				DataType: schema.DataTypeUint16,
				Count:    1,
			},
			details: fmt.Sprintf(
				"host=%s port=%d slave=%d function=%s address=%s",
				cfg.Host, cfg.Port, cfg.SlaveID, function, addr,
			),
		}, nil
	case schema.ProtocolModbusUDP:
		var cfg schema.ConnectionConfigModbusUDP
		if err := json.Unmarshal([]byte(device.ConnectionConfig), &cfg); err != nil {
			return readProbeTarget{}, fmt.Errorf("解析連線設定失敗: %w", err)
		}
		addr := pickProbeValue(overrideAddr, "40001")
		function := pickProbeValue(overrideFunc, "03")
		return readProbeTarget{
			enabled: true,
			req: connector.ReadRequest{
				Address:  addr,
				Function: function,
				DataType: schema.DataTypeUint16,
				Count:    1,
			},
			details: fmt.Sprintf(
				"host=%s port=%d slave=%d function=%s address=%s",
				cfg.Host, cfg.Port, cfg.SlaveID, function, addr,
			),
		}, nil
	case schema.ProtocolModbusRTU:
		var cfg schema.ConnectionConfigModbusRTU
		if err := json.Unmarshal([]byte(device.ConnectionConfig), &cfg); err != nil {
			return readProbeTarget{}, fmt.Errorf("解析連線設定失敗: %w", err)
		}
		addr := pickProbeValue(overrideAddr, "40001")
		function := pickProbeValue(overrideFunc, "03")
		return readProbeTarget{
			enabled: true,
			req: connector.ReadRequest{
				Address:  addr,
				Function: function,
				DataType: schema.DataTypeUint16,
				Count:    1,
			},
			details: fmt.Sprintf(
				"serial=%s baud=%d slave=%d function=%s address=%s",
				cfg.SerialPort, cfg.BaudRate, cfg.SlaveID, function, addr,
			),
		}, nil
	case schema.ProtocolFatekFBs:
		var cfg schema.ConnectionConfigFatekFBs
		if err := json.Unmarshal([]byte(device.ConnectionConfig), &cfg); err != nil {
			return readProbeTarget{}, fmt.Errorf("解析連線設定失敗: %w", err)
		}
		addr := pickProbeValue(overrideAddr, "D0")
		return readProbeTarget{
			enabled: true,
			req: connector.ReadRequest{
				Address:  addr,
				Function: pickProbeValue(overrideFunc, ""),
				DataType: schema.DataTypeInt16,
				Count:    1,
			},
			details: fmt.Sprintf(
				"mode=%s host=%s port=%d station=%d address=%s",
				cfg.Mode, cfg.Host, cfg.Port, cfg.StationNo, addr,
			),
		}, nil
	case schema.ProtocolMC3E:
		var cfg schema.ConnectionConfigMC3E
		if err := json.Unmarshal([]byte(device.ConnectionConfig), &cfg); err != nil {
			return readProbeTarget{}, fmt.Errorf("解析連線設定失敗: %w", err)
		}
		addr := pickProbeValue(overrideAddr, "D0")
		return readProbeTarget{
			enabled: true,
			req: connector.ReadRequest{
				Address:  addr,
				Function: pickProbeValue(overrideFunc, ""),
				DataType: schema.DataTypeInt16,
				Count:    1,
			},
			details: fmt.Sprintf(
				"host=%s port=%d station=%d address=%s",
				cfg.Host, cfg.Port, cfg.StationNo, addr,
			),
		}, nil
	default:
		return readProbeTarget{enabled: false}, nil
	}
}

func extractProbeOverride(config string) (address, function string) {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(config), &raw); err != nil {
		return "", ""
	}
	address = extractProbeString(raw["probe_address"])
	function = strings.TrimSpace(extractProbeString(raw["probe_function"]))
	return address, function
}

func extractProbeString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v)
	case float64:
		return strconv.FormatInt(int64(v), 10)
	default:
		return ""
	}
}

func pickProbeValue(override, fallback string) string {
	if strings.TrimSpace(override) == "" {
		return fallback
	}
	return strings.TrimSpace(override)
}
