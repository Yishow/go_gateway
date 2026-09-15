package schema

import "time"

// =============================================================================
// 設備模型 (Device)
// =============================================================================

// Device 設備實體，代表可連接的 PLC 或資料來源
type Device struct {
	// ID 主鍵，UUID 格式
	ID string `json:"id" db:"id"`

	// Name 設備顯示名稱
	Name string `json:"name" db:"name"`

	// Description 設備描述
	Description string `json:"description,omitempty" db:"description"`

	// Protocol 通訊協議類型
	Protocol ProtocolType `json:"protocol" db:"protocol"`

	// Status 設備狀態
	Status DeviceStatus `json:"status" db:"status"`

	// ConnectionConfig 協議特定的連線配置 (JSON 格式)
	// 不同協議有不同的欄位結構，詳見 ConnectionConfigXXX 類型
	ConnectionConfig string `json:"connection_config" db:"connection_config"`

	// LastTestAt 最後連線測試時間
	LastTestAt *time.Time `json:"last_test_at,omitempty" db:"last_test_at"`

	// LastTestSuccess 最後連線測試是否成功
	LastTestSuccess *bool `json:"last_test_success,omitempty" db:"last_test_success"`

	// LastTestError 最後連線測試錯誤訊息
	LastTestError string `json:"last_test_error,omitempty" db:"last_test_error"`

	// LastCollectedAt 最後資料收集時間 (新增欄位)
	LastCollectedAt *time.Time `json:"last_collected_at,omitempty" db:"last_collected_at"`

	// CollectionCount 總收集次數 (新增欄位)
	CollectionCount int64 `json:"collection_count" db:"collection_count"`

	// ErrorCount 總錯誤次數 (新增欄位)
	ErrorCount int64 `json:"error_count" db:"error_count"`

	// ReadinessStatus 設備就緒狀態 (新增欄位)
	// JSON 儲存多維 readiness contract，至少含 connect/probe 狀態、planning/activation/apply eligibility 與 blocking reasons。
	ReadinessStatus string `json:"readiness_status,omitempty" db:"readiness_status"`

	// CreatedAt 建立時間
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	// UpdatedAt 更新時間
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ReadinessCheck 單項檢查結果
type ReadinessCheck struct {
	Name    string `json:"name"`
	Pass    bool   `json:"pass"`
	Message string `json:"message"`
}

// ReadinessStageStatus describes one diagnostics stage outcome inside the readiness contract.
type ReadinessStageStatus string

const (
	// ReadinessStageStatusUnknown means the diagnostics stage has not produced a conclusive result yet.
	ReadinessStageStatusUnknown ReadinessStageStatus = "unknown"
	// ReadinessStageStatusSuccess means the diagnostics stage completed successfully.
	ReadinessStageStatusSuccess ReadinessStageStatus = "success"
	// ReadinessStageStatusFailed means the diagnostics stage failed.
	ReadinessStageStatusFailed ReadinessStageStatus = "failed"
	// ReadinessStageStatusSkipped means the diagnostics stage is unavailable or intentionally skipped.
	ReadinessStageStatusSkipped ReadinessStageStatus = "skipped"
)

// DeviceReadiness 設備就緒狀態詳情
type DeviceReadiness struct {
	DeviceID           string               `json:"device_id"`
	Status             string               `json:"status"` // ready, warning, error
	AvailabilityStatus string               `json:"availability_status,omitempty"`
	AvailabilityReason string               `json:"availability_reason,omitempty"`
	Checks             []ReadinessCheck     `json:"checks"`
	ConnectStatus      ReadinessStageStatus `json:"connect_status"`
	ProbeStatus        ReadinessStageStatus `json:"probe_status"`
	PlanningAllowed    bool                 `json:"planning_allowed"`
	ActivationAllowed  bool                 `json:"activation_allowed"`
	ApplyAllowed       bool                 `json:"apply_allowed"`
	BlockingReasons    []string             `json:"blocking_reasons"`
}

// =============================================================================
// 連線配置結構 (協議特定)
// =============================================================================

// ConnectionConfigModbusTCP Modbus TCP 連線配置
type ConnectionConfigModbusTCP struct {
	// Host 目標主機位址
	Host string `json:"host"`
	// Port TCP 埠號 (預設 502)
	Port int `json:"port"`
	// SlaveID 從站 ID (1-247)
	SlaveID byte `json:"slave_id"`
	// Timeout 連線逾時秒數
	Timeout int `json:"timeout"`
	// DataFormat 字節序格式 (ABCD, BADC, CDAB, DCBA)，預設為 ABCD
	DataFormat string `json:"data_format,omitempty"`
}

// ConnectionConfigModbusRTU Modbus RTU 連線配置
type ConnectionConfigModbusRTU struct {
	// SerialPort 串列埠名稱 (如 COM1, /dev/ttyUSB0)
	SerialPort string `json:"serial_port"`
	// BaudRate 鮑率 (9600, 19200, 38400, 57600, 115200)
	BaudRate int `json:"baud_rate"`
	// DataBits 資料位元 (7 或 8)
	DataBits int `json:"data_bits"`
	// StopBits 停止位元 (1 或 2)
	StopBits int `json:"stop_bits"`
	// Parity 同位檢查 (none, odd, even)
	Parity string `json:"parity"`
	// SlaveID 從站 ID (1-247)
	SlaveID byte `json:"slave_id"`
	// Timeout 通訊逾時秒數
	Timeout int `json:"timeout"`
	// DataFormat 字節序格式 (ABCD, BADC, CDAB, DCBA)，預設為 ABCD
	DataFormat string `json:"data_format,omitempty"`
}

// ConnectionConfigModbusUDP Modbus UDP 連線配置
type ConnectionConfigModbusUDP struct {
	// Host 目標主機位址
	Host string `json:"host"`
	// Port UDP 埠號
	Port int `json:"port"`
	// SlaveID 從站 ID
	SlaveID byte `json:"slave_id"`
	// Timeout 通訊逾時秒數
	Timeout int `json:"timeout"`
	// DataFormat 字節序格式 (ABCD, BADC, CDAB, DCBA)，預設為 ABCD
	DataFormat string `json:"data_format,omitempty"`
}

// ConnectionConfigFatekFBs FATEK FBs 連線配置
type ConnectionConfigFatekFBs struct {
	// Mode 連線模式 (tcp 或 serial)
	Mode string `json:"mode"`
	// Host TCP 模式的主機位址
	Host string `json:"host,omitempty"`
	// Port TCP 模式的埠號 (預設 500)
	Port int `json:"port,omitempty"`
	// SerialPort 串列模式的埠名稱
	SerialPort string `json:"serial_port,omitempty"`
	// BaudRate 串列模式的鮑率
	BaudRate int `json:"baud_rate,omitempty"`
	// DataBits 串列模式資料位元（7 或 8；預設 7，符合 FATEK ASCII 常見設定）
	DataBits int `json:"data_bits,omitempty"`
	// StopBits 串列模式停止位元（1 或 2；預設 1）
	StopBits int `json:"stop_bits,omitempty"`
	// Parity 串列模式同位：none / odd / even（預設 even）
	Parity string `json:"parity,omitempty"`
	// StationNo 站號 (0-255)
	StationNo byte `json:"station_no"`
	// Timeout 通訊逾時秒數
	Timeout int `json:"timeout"`
}

// ConnectionConfigMC3E Mitsubishi MC 3E 連線配置
type ConnectionConfigMC3E struct {
	// Host 目標主機位址
	Host string `json:"host"`
	// Port TCP 埠號 (預設依 PLC 型號)
	Port int `json:"port"`
	// NetworkNo 網路編號
	NetworkNo byte `json:"network_no"`
	// PCNo PC 編號
	PCNo byte `json:"pc_no"`
	// IONo I/O 編號
	IONo uint16 `json:"io_no"`
	// StationNo 站號
	StationNo byte `json:"station_no"`
	// Timeout 通訊逾時秒數
	Timeout int `json:"timeout"`
	// DataFormat 字節序格式 (ABCD, BADC, CDAB, DCBA)，預設為 CDAB (三菱標準)
	DataFormat string `json:"data_format,omitempty"`
}

// ConnectionConfigMQTT MQTT 接收連線配置
type ConnectionConfigMQTT struct {
	// BrokerURL MQTT Broker 位址 (如 tcp://localhost:1883)
	BrokerURL string `json:"broker_url"`
	// ClientID 客戶端識別碼
	ClientID string `json:"client_id"`
	// Username 使用者名稱 (選用)
	Username string `json:"username,omitempty"`
	// Password 密碼 (選用)
	Password string `json:"password,omitempty"`
	// UseTLS 是否使用 TLS
	UseTLS bool `json:"use_tls"`
	// Topics 訂閱的主題清單
	Topics []string `json:"topics"`
	// QoS 服務品質等級 (0, 1, 2)
	QoS byte `json:"qos"`
}
