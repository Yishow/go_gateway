package handlers

import (
	"fmt"
	"time"

	"go-gateway/internal/protocol/fatek"
	"go-gateway/internal/protocol/mcprotocol"
	"go-gateway/internal/protocol/modbus"
)

// --- Helper Functions ---

// createClientWithDebug 創建帶有數據包記錄功能的客戶端
func (h *TestHandler) createClientWithDebug(protocol string, config map[string]interface{}, connectionID string) (interface{}, error) {
	// Helper to safe cast config values
	getString := func(key string, def string) string {
		if v, ok := config[key].(string); ok {
			return v
		}
		return def
	}
	getInt := func(key string, def int) int {
		if v, ok := config[key].(float64); ok {
			return int(v)
		}
		if v, ok := config[key].(int); ok {
			return v
		}
		return def
	}
	getDuration := func(key string, def time.Duration) time.Duration {
		if v, ok := config[key].(float64); ok {
			return time.Duration(v) * time.Millisecond
		}
		return def
	}

	switch protocol {
	case "modbus_tcp":
		host := getString("host", "localhost")
		port := getInt("port", 502)
		unitID := byte(getInt("unitID", 1))
		timeout := getDuration("timeout", 2000*time.Millisecond)
		transport := modbus.NewTCPTransport(host, port)
		if timeout > 0 {
			transport.Timeout = timeout
		}
		// 包裝 Transport 以記錄數據包（如果 debugHandler 存在）
		if h.debugHandler != nil {
			wrappedTransport := NewWrappedTransport(transport, connectionID, protocol, h.debugHandler)
			return modbus.NewClient(wrappedTransport, unitID), nil
		}
		return modbus.NewClient(transport, unitID), nil

	case "modbus_udp":
		host := getString("host", "localhost")
		port := getInt("port", 502)
		unitID := byte(getInt("unitID", 1))
		timeout := getDuration("timeout", 2000*time.Millisecond)
		transport := modbus.NewUDPTransport(host, port)
		if timeout > 0 {
			transport.Timeout = timeout
		}
		if h.debugHandler != nil {
			wrappedTransport := NewWrappedTransport(transport, connectionID, protocol, h.debugHandler)
			return modbus.NewClient(wrappedTransport, unitID), nil
		}
		return modbus.NewClient(transport, unitID), nil

	case "modbus_rtu":
		port := getString("port", "COM1")
		baudRate := getInt("baudRate", 9600)
		dataBits := getInt("dataBits", 8)
		stopBits := getInt("stopBits", 1)
		parity := getString("parity", "N")
		unitID := byte(getInt("unitID", 1))
		timeout := getDuration("timeout", 2000*time.Millisecond)
		transport := modbus.NewRTUTransport(port, baudRate, dataBits, stopBits, parity, timeout)
		if h.debugHandler != nil {
			wrappedTransport := NewWrappedTransport(transport, connectionID, protocol, h.debugHandler)
			return modbus.NewClient(wrappedTransport, unitID), nil
		}
		return modbus.NewClient(transport, unitID), nil

	case "fatek_tcp":
		host := getString("host", "localhost")
		port := getInt("port", 500)
		station := getInt("station", 1)
		timeout := getDuration("timeout", 2000*time.Millisecond)
		transport := fatek.NewTCPTransport(host, port)
		if timeout > 0 {
			transport.Timeout = timeout
		}
		if h.debugHandler != nil {
			wrappedTransport := NewWrappedTransport(transport, connectionID, protocol, h.debugHandler)
			return fatek.NewClient(wrappedTransport, station), nil
		}
		return fatek.NewClient(transport, station), nil

	case "fatek_serial":
		port := getString("port", "COM1")
		baudRate := getInt("baudRate", 9600)
		dataBits := getInt("dataBits", 7)
		stopBits := getInt("stopBits", 1)
		parity := getString("parity", "E")
		station := getInt("station", 1)
		timeout := getDuration("timeout", 1000*time.Millisecond)
		transport := fatek.NewSerialTransport(port, baudRate, dataBits, stopBits, parity, timeout)
		if h.debugHandler != nil {
			wrappedTransport := NewWrappedTransport(transport, connectionID, protocol, h.debugHandler)
			return fatek.NewClient(wrappedTransport, station), nil
		}
		return fatek.NewClient(transport, station), nil

	case "mc_tcp", "mcprotocol_tcp":
		host := getString("host", "localhost")
		port := getInt("port", 6000)
		timeout := getDuration("timeout", 2000*time.Millisecond)
		transport := mcprotocol.NewTCPTransport(host, port)
		if timeout > 0 {
			transport.Timeout = timeout
		}
		if h.debugHandler != nil {
			wrappedTransport := NewWrappedTransport(transport, connectionID, protocol, h.debugHandler)
			return mcprotocol.NewClientWithTransport(wrappedTransport), nil
		}
		return mcprotocol.NewClientWithTransport(transport), nil

	case "mcprotocol_serial":
		port := getString("port", "COM1")
		baudRate := getInt("baudRate", 9600)
		dataBits := getInt("dataBits", 7)
		stopBits := getInt("stopBits", 2)
		parity := getString("parity", "E")
		timeout := getDuration("timeout", 2000*time.Millisecond)
		transport := mcprotocol.NewSerialTransport(port, baudRate, dataBits, stopBits, parity, timeout)
		if h.debugHandler != nil {
			wrappedTransport := NewWrappedTransport(transport, connectionID, protocol, h.debugHandler)
			return mcprotocol.NewClientWithTransport(wrappedTransport), nil
		}
		return mcprotocol.NewClientWithTransport(transport), nil

	default:
		return nil, fmt.Errorf("不支援的協議: %s", protocol)
	}
}

// createClient 保留原方法以向後兼容（如果其他地方有調用）
// 如果沒有 debugHandler，使用原始的工廠方法創建客戶端
func (h *TestHandler) createClient(protocol string, config map[string]interface{}) (interface{}, error) {
	// Helper to safe cast config values
	getString := func(key string, def string) string {
		if v, ok := config[key].(string); ok {
			return v
		}
		return def
	}
	getInt := func(key string, def int) int {
		if v, ok := config[key].(float64); ok {
			return int(v)
		}
		if v, ok := config[key].(int); ok {
			return v
		}
		return def
	}
	getDuration := func(key string, def time.Duration) time.Duration {
		if v, ok := config[key].(float64); ok {
			return time.Duration(v) * time.Millisecond
		}
		return def
	}

	switch protocol {
	case "modbus_tcp":
		host := getString("host", "localhost")
		port := getInt("port", 502)
		unitID := byte(getInt("unitID", 1))
		timeout := getDuration("timeout", 2000*time.Millisecond)
		return modbus.CreateTCPClient(host, port, unitID, timeout), nil

	case "modbus_udp":
		host := getString("host", "localhost")
		port := getInt("port", 502)
		unitID := byte(getInt("unitID", 1))
		timeout := getDuration("timeout", 2000*time.Millisecond)
		return modbus.CreateUDPClient(host, port, unitID, timeout), nil

	case "modbus_rtu":
		port := getString("port", "COM1")
		baudRate := getInt("baudRate", 9600)
		dataBits := getInt("dataBits", 8)
		stopBits := getInt("stopBits", 1)
		parity := getString("parity", "N")
		unitID := byte(getInt("unitID", 1))
		timeout := getDuration("timeout", 2000*time.Millisecond)
		return modbus.CreateRTUClient(port, baudRate, dataBits, stopBits, parity, timeout, unitID), nil

	case "fatek_tcp":
		host := getString("host", "localhost")
		port := getInt("port", 500)
		station := getInt("station", 1)
		timeout := getDuration("timeout", 2000*time.Millisecond)
		return fatek.CreateTCPClient(host, port, station, timeout), nil

	case "fatek_serial":
		port := getString("port", "COM1")
		baudRate := getInt("baudRate", 9600)
		dataBits := getInt("dataBits", 7)
		stopBits := getInt("stopBits", 1)
		parity := getString("parity", "E")
		station := getInt("station", 1)
		timeout := getDuration("timeout", 1000*time.Millisecond)
		return fatek.CreateSerialClient(port, station, baudRate, dataBits, stopBits, parity, timeout), nil

	case "mc_tcp", "mcprotocol_tcp":
		host := getString("host", "localhost")
		port := getInt("port", 6000)
		timeout := getDuration("timeout", 2000*time.Millisecond)
		return mcprotocol.CreateTCPClient(host, port, timeout), nil

	case "mcprotocol_serial":
		port := getString("port", "COM1")
		baudRate := getInt("baudRate", 9600)
		dataBits := getInt("dataBits", 7)
		stopBits := getInt("stopBits", 2)
		parity := getString("parity", "E")
		timeout := getDuration("timeout", 2000*time.Millisecond)
		return mcprotocol.CreateSerialClient(port, baudRate, dataBits, stopBits, parity, timeout), nil

	default:
		return nil, fmt.Errorf("unsupported protocol: %s", protocol)
	}
}
