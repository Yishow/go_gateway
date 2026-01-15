package modbus

// Modbus 功能碼 (Function Codes)
const (
	// 讀取功能碼
	FuncReadCoils            = 0x01 // 讀取線圈狀態 (1-2000)
	FuncReadDiscreteInputs   = 0x02 // 讀取離散輸入狀態 (1-2000)
	FuncReadHoldingRegisters = 0x03 // 讀取保持暫存器 (1-125)
	FuncReadInputRegisters   = 0x04 // 讀取輸入暫存器 (1-125)

	// 寫入功能碼
	FuncWriteSingleCoil      = 0x05 // 寫入單個線圈
	FuncWriteSingleRegister  = 0x06 // 寫入單個暫存器
	FuncWriteMultipleCoils   = 0x0F // 寫入多個線圈 (1-1968)
	FuncWriteMultipleRegisters = 0x10 // 寫入多個暫存器 (1-123)

	// 異常功能碼 (功能碼 + 0x80)
	FuncExceptionOffset = 0x80
)

// Modbus 異常碼 (Exception Codes)
const (
	ExceptionIllegalFunction         = 0x01 // 非法功能碼
	ExceptionIllegalDataAddress      = 0x02 // 非法數據地址
	ExceptionIllegalDataValue       = 0x03 // 非法數據值
	ExceptionServerDeviceFailure    = 0x04 // 服務器設備故障
	ExceptionAcknowledge           = 0x05 // 確認
	ExceptionServerDeviceBusy       = 0x06 // 服務器設備忙
	ExceptionMemoryParityError      = 0x08 // 記憶體奇偶校驗錯誤
	ExceptionGatewayPathUnavailable = 0x0A // 網關路徑不可用
	ExceptionGatewayTargetNoResponse = 0x0B // 網關目標設備無響應
)

// Modbus 地址區域定義
const (
	// 線圈 (Coils) - 可讀寫
	CoilStartAddress = 0x0000
	CoilMaxAddress  = 0xFFFF
	CoilMaxQuantity  = 2000

	// 離散輸入 (Discrete Inputs) - 只讀
	DiscreteInputStartAddress = 0x10000
	DiscreteInputMaxAddress   = 0x1FFFF
	DiscreteInputMaxQuantity  = 2000

	// 保持暫存器 (Holding Registers) - 可讀寫
	HoldingRegisterStartAddress = 0x40000
	HoldingRegisterMaxAddress   = 0x4FFFF
	HoldingRegisterMaxQuantity  = 125

	// 輸入暫存器 (Input Registers) - 只讀
	InputRegisterStartAddress = 0x30000
	InputRegisterMaxAddress   = 0x3FFFF
	InputRegisterMaxQuantity = 125
)

// Modbus TCP/UDP MBAP Header 長度
const (
	MBAPHeaderLength = 7 // Transaction ID(2) + Protocol ID(2) + Length(2) + Unit ID(1)
	TCPDefaultPort   = 502
	UDPDefaultPort   = 502
)

// Modbus RTU 相關常數
const (
	RTUFrameMinLength = 4  // Address(1) + Function(1) + CRC(2)
	RTUFrameMaxLength = 256 // 最大封包長度
)

// 預設配置
const (
	DefaultUnitID    = 1
	DefaultTimeout   = 2 // 秒
	DefaultBaudRate  = 9600
	DefaultDataBits  = 8
	DefaultStopBits  = 1
	DefaultParity    = "N" // N=None, E=Even, O=Odd
)

// 異常碼對應的錯誤訊息
var ExceptionMessages = map[byte]string{
	ExceptionIllegalFunction:          "非法功能碼",
	ExceptionIllegalDataAddress:        "非法數據地址",
	ExceptionIllegalDataValue:          "非法數據值",
	ExceptionServerDeviceFailure:       "服務器設備故障",
	ExceptionAcknowledge:              "確認",
	ExceptionServerDeviceBusy:         "服務器設備忙",
	ExceptionMemoryParityError:        "記憶體奇偶校驗錯誤",
	ExceptionGatewayPathUnavailable:    "網關路徑不可用",
	ExceptionGatewayTargetNoResponse:  "網關目標設備無響應",
}
