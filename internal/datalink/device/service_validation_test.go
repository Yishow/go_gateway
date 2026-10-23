package device

import (
	"testing"

	_ "go-gateway/internal/datalink/connector/adapters"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// validModbusTCPConfig 回傳通過驗證的 Modbus TCP 配置，供各案例複製後修改
func validModbusTCPConfig() map[string]any {
	return map[string]any{
		"host":     "127.0.0.1",
		"port":     float64(502),
		"slave_id": float64(1),
	}
}

func TestValidateConnectionConfigListsAllMissingFields(t *testing.T) {
	tests := []struct {
		name     string
		protocol schema.ProtocolType
		config   map[string]any
		wantErr  string
	}{
		{
			name:     "Fatek FBs 空 map 一次列齊 mode 與 station_no",
			protocol: schema.ProtocolFatekFBs,
			config:   map[string]any{},
			wantErr:  "缺少必填欄位: mode, station_no",
		},
		{
			name:     "Modbus TCP 空 map 一次列齊 host、port 與 slave_id",
			protocol: schema.ProtocolModbusTCP,
			config:   map[string]any{},
			wantErr:  "缺少必填欄位: host, port, slave_id",
		},
		{
			name:     "Modbus TCP nil map 與空 map 行為一致",
			protocol: schema.ProtocolModbusTCP,
			config:   nil,
			wantErr:  "缺少必填欄位: host, port, slave_id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConnectionConfig(tt.protocol, tt.config)
			require.Error(t, err)
			assert.EqualError(t, err, tt.wantErr)
		})
	}
}

func TestValidateConnectionConfigSingleMissingField(t *testing.T) {
	tests := []struct {
		name     string
		protocol schema.ProtocolType
		config   map[string]any
		wantErr  string
	}{
		{
			name:     "Modbus TCP 只缺 slave_id",
			protocol: schema.ProtocolModbusTCP,
			config:   map[string]any{"host": "127.0.0.1", "port": float64(502)},
			wantErr:  "缺少必填欄位: slave_id",
		},
		{
			name:     "Modbus TCP 只缺 host",
			protocol: schema.ProtocolModbusTCP,
			config:   map[string]any{"port": float64(502), "slave_id": float64(1)},
			wantErr:  "缺少必填欄位: host",
		},
		{
			name:     "Fatek FBs 只缺 station_no",
			protocol: schema.ProtocolFatekFBs,
			config:   map[string]any{"mode": "tcp"},
			wantErr:  "缺少必填欄位: station_no",
		},
		{
			name:     "Fatek FBs 只缺 mode",
			protocol: schema.ProtocolFatekFBs,
			config:   map[string]any{"station_no": float64(1)},
			wantErr:  "缺少必填欄位: mode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConnectionConfig(tt.protocol, tt.config)
			require.Error(t, err)
			assert.EqualError(t, err, tt.wantErr)
		})
	}
}

func TestValidateConnectionConfigEmptyFields(t *testing.T) {
	tests := []struct {
		name     string
		protocol schema.ProtocolType
		config   map[string]any
		wantErr  string
	}{
		{
			name:     "單一欄位為空字串",
			protocol: schema.ProtocolModbusTCP,
			config:   map[string]any{"host": "", "port": float64(502), "slave_id": float64(1)},
			wantErr:  "必填欄位 host 不能為空",
		},
		{
			name:     "單一欄位為 nil",
			protocol: schema.ProtocolModbusTCP,
			config:   map[string]any{"host": "127.0.0.1", "port": nil, "slave_id": float64(1)},
			wantErr:  "必填欄位 port 不能為空",
		},
		{
			name:     "多個欄位為空字串或 nil 時一次列齊",
			protocol: schema.ProtocolModbusTCP,
			config:   map[string]any{"host": "", "port": nil, "slave_id": float64(1)},
			wantErr:  "必填欄位不能為空: host, port",
		},
		{
			name:     "所有必填欄位皆為空時一次列齊",
			protocol: schema.ProtocolModbusTCP,
			config:   map[string]any{"host": "", "port": nil, "slave_id": nil},
			wantErr:  "必填欄位不能為空: host, port, slave_id",
		},
		{
			name:     "Fatek FBs mode 為空字串",
			protocol: schema.ProtocolFatekFBs,
			config:   map[string]any{"mode": "", "station_no": float64(1)},
			wantErr:  "必填欄位 mode 不能為空",
		},
		{
			name:     "Fatek FBs mode 與 station_no 皆為 nil",
			protocol: schema.ProtocolFatekFBs,
			config:   map[string]any{"mode": nil, "station_no": nil},
			wantErr:  "必填欄位不能為空: mode, station_no",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConnectionConfig(tt.protocol, tt.config)
			require.Error(t, err)
			assert.EqualError(t, err, tt.wantErr)
		})
	}
}

func TestValidateConnectionConfigMixedRequiredFields(t *testing.T) {
	tests := []struct {
		name     string
		protocol schema.ProtocolType
		config   map[string]any
		wantErr  string
	}{
		{
			name:     "Modbus 同時回報缺少 port 與空 host",
			protocol: schema.ProtocolModbusTCP,
			config:   map[string]any{"host": "", "slave_id": float64(1)},
			wantErr:  "缺少必填欄位: port\n必填欄位 host 不能為空",
		},
		{
			name:     "Fatek 同時回報缺少 station_no 與空 mode",
			protocol: schema.ProtocolFatekFBs,
			config:   map[string]any{"mode": ""},
			wantErr:  "缺少必填欄位: station_no\n必填欄位 mode 不能為空",
		},
		{
			name:     "多個缺少欄位依 schema 順序排在 null 欄位之前",
			protocol: schema.ProtocolModbusTCP,
			config:   map[string]any{"port": nil},
			wantErr:  "缺少必填欄位: host, slave_id\n必填欄位 port 不能為空",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConnectionConfig(tt.protocol, tt.config)
			require.Error(t, err)
			assert.EqualError(t, err, tt.wantErr)
		})
	}
}

func TestValidateConnectionConfigErrorPriority(t *testing.T) {
	t.Run("缺少與空值同時回報且缺少在前", func(t *testing.T) {
		err := validateConnectionConfig(schema.ProtocolModbusTCP, map[string]any{"host": ""})
		require.Error(t, err)
		assert.EqualError(t, err, "缺少必填欄位: port, slave_id\n必填欄位 host 不能為空")
	})

	t.Run("缺少欄位優先於數值範圍錯誤", func(t *testing.T) {
		err := validateConnectionConfig(schema.ProtocolModbusTCP, map[string]any{
			"host": "127.0.0.1",
			"port": float64(0),
		})
		require.Error(t, err)
		assert.EqualError(t, err, "缺少必填欄位: slave_id")
	})

	t.Run("空值欄位優先於數值範圍錯誤", func(t *testing.T) {
		err := validateConnectionConfig(schema.ProtocolModbusTCP, map[string]any{
			"host":     "",
			"port":     float64(0),
			"slave_id": float64(1),
		})
		require.Error(t, err)
		assert.EqualError(t, err, "必填欄位 host 不能為空")
	})
}

func TestValidateConnectionConfigKeepsRangeAndEnumChecks(t *testing.T) {
	tests := []struct {
		name     string
		protocol schema.ProtocolType
		config   map[string]any
		wantErr  string
	}{
		{
			name:     "數值小於最小值",
			protocol: schema.ProtocolModbusTCP,
			config:   map[string]any{"host": "127.0.0.1", "port": float64(0), "slave_id": float64(1)},
			wantErr:  "欄位 port 的值 0 小於最小值 1",
		},
		{
			name:     "數值大於最大值",
			protocol: schema.ProtocolModbusTCP,
			config:   map[string]any{"host": "127.0.0.1", "port": float64(70000), "slave_id": float64(1)},
			wantErr:  "欄位 port 的值 70000 大於最大值 65535",
		},
		{
			name:     "多個範圍錯誤只回報第一個",
			protocol: schema.ProtocolModbusTCP,
			config:   map[string]any{"host": "127.0.0.1", "port": float64(0), "slave_id": float64(300)},
			wantErr:  "欄位 port 的值 0 小於最小值 1",
		},
		{
			name:     "字串不在 enum 內",
			protocol: schema.ProtocolFatekFBs,
			config:   map[string]any{"mode": "udp", "station_no": float64(1)},
			wantErr:  "欄位 mode 的值 udp 不在允許的選項中",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConnectionConfig(tt.protocol, tt.config)
			require.Error(t, err)
			assert.EqualError(t, err, tt.wantErr)
		})
	}
}

func TestValidateConnectionConfigAcceptsValidConfig(t *testing.T) {
	t.Run("Modbus TCP", func(t *testing.T) {
		assert.NoError(t, validateConnectionConfig(schema.ProtocolModbusTCP, validModbusTCPConfig()))
	})

	t.Run("Fatek FBs", func(t *testing.T) {
		config := map[string]any{"mode": "tcp", "station_no": float64(1)}
		assert.NoError(t, validateConnectionConfig(schema.ProtocolFatekFBs, config))
	})
}

func TestIsConditionallyOptional(t *testing.T) {
	tests := []struct {
		name      string
		protocol  schema.ProtocolType
		fieldName string
		config    map[string]any
		want      bool
	}{
		{"serial 模式不需要 host", schema.ProtocolFatekFBs, "host", map[string]any{"mode": "serial"}, true},
		{"serial 模式不需要 port", schema.ProtocolFatekFBs, "port", map[string]any{"mode": "serial"}, true},
		{"tcp 模式需要 host", schema.ProtocolFatekFBs, "host", map[string]any{"mode": "tcp"}, false},
		{"tcp 模式不需要 serial_port", schema.ProtocolFatekFBs, "serial_port", map[string]any{"mode": "tcp"}, true},
		{"tcp 模式不需要 baud_rate", schema.ProtocolFatekFBs, "baud_rate", map[string]any{"mode": "tcp"}, true},
		{"serial 模式需要 serial_port", schema.ProtocolFatekFBs, "serial_port", map[string]any{"mode": "serial"}, false},
		{"mode 缺少時不略過任何欄位", schema.ProtocolFatekFBs, "host", map[string]any{}, false},
		{"mode 非字串時不略過任何欄位", schema.ProtocolFatekFBs, "serial_port", map[string]any{"mode": nil}, false},
		{"非 Fatek 協議不套用條件", schema.ProtocolModbusTCP, "host", map[string]any{"mode": "serial"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isConditionallyOptional(tt.protocol, tt.fieldName, tt.config)
			assert.Equal(t, tt.want, got)
		})
	}
}
