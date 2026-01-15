package fatek

// 協定常數定義
const (
	// STX 起始字元，固定為 0x02
	STX = byte(0x02)
	// ETX 結束字元，固定為 0x03
	ETX = byte(0x03)
	// DefaultStation 預設站號
	DefaultStation = 1
	// DefaultTCPPort 預設 TCP 埠號
	DefaultTCPPort = 500
	// DefaultTimeout 預設逾時時間（秒）
	DefaultTimeout = 2.0
	// DefaultBaudRate 預設串列埠波特率
	DefaultBaudRate = 9600
	// DefaultDataBits 預設資料位元數（FATEK ASCII 標準為 7）
	DefaultDataBits = 7
	// DefaultParity 預設同位檢查（Even）
	DefaultParity = "E"
	// DefaultStopBits 預設停止位元數
	DefaultStopBits = 1
	// MaxDiscreteCount 單點讀寫最大數量
	MaxDiscreteCount = 255
	// MaxRegisterCount 暫存器讀寫最大數量
	MaxRegisterCount = 64
)
