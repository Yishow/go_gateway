package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"go-gateway/internal/protocol/fatek"
)

func main() {
	var (
		mode       = flag.String("mode", "", "連接模式: tcp 或 serial (必填)")
		host       = flag.String("host", "", "TCP 模式: PLC IP 位址")
		port       = flag.Int("port", 500, "TCP 模式: TCP 埠號 (預設 500)")
		serialPort = flag.String("serial", "", "Serial 模式: 串列埠名稱 (例如 COM3 或 /dev/ttyUSB0)")
		baudRate   = flag.Int("baud", 9600, "Serial 模式: 波特率 (預設 9600)")
		dataBits   = flag.Int("databits", 7, "Serial 模式: 資料位元數 (預設 7)")
		stopBits   = flag.Int("stopbits", 1, "Serial 模式: 停止位元數 (預設 1)")
		parity     = flag.String("parity", "E", "Serial 模式: 同位檢查 N/E/O (預設 E)")
		station    = flag.Int("station", 1, "PLC 站號 (預設 1)")
		timeout    = flag.Duration("timeout", 2*time.Second, "逾時時間 (預設 2s)")
		action     = flag.String("action", "", "操作類型: read, write, random (必填)")
		symbol     = flag.String("symbol", "", "組件符號: X, Y, M, S, T, C, R, D, RT, RC, F, DR 等")
		addr       = flag.Int("addr", 0, "起始位址")
		count      = flag.Int("count", 1, "讀取數量")
		values     = flag.String("values", "", "寫入值 (逗號分隔，例如: 1,2,3 或 true,false,true)")
		random     = flag.String("random", "", "隨機讀取項目 (格式: SYMBOL1:ADDR1,SYMBOL2:ADDR2)")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Fatek PLC 測試工具\n\n")
		fmt.Fprintf(os.Stderr, "用法:\n")
		fmt.Fprintf(os.Stderr, "  %s [選項]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "TCP 模式範例:\n")
		fmt.Fprintf(os.Stderr, "  %s -mode=tcp -host=\"192.168.1.5\" -port=500 -action=read -symbol=D -addr=0 -count=10\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -mode=tcp -host=\"192.168.1.5\" -action=write -symbol=Y -addr=0 -values=true,false,true\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nSerial 模式範例:\n")
		fmt.Fprintf(os.Stderr, "  %s -mode=serial -serial=COM3 -action=read -symbol=X -addr=0 -count=5\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -mode=serial -serial=/dev/ttyUSB0 -action=write -symbol=D -addr=0 -values=100,200,300\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n隨機讀取範例:\n")
		fmt.Fprintf(os.Stderr, "  %s -mode=tcp -host=\"192.168.1.5\" -action=random -random=\"X:0,D:0,R:10\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n注意：在 Windows PowerShell 中，IP 位址和包含逗號的值需用引號包裹\n")
		fmt.Fprintf(os.Stderr, "\n選項:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *mode == "" {
		fmt.Fprintf(os.Stderr, "錯誤: 必須指定 -mode (tcp 或 serial)\n\n")
		flag.Usage()
		os.Exit(1)
	}

	if *action == "" {
		fmt.Fprintf(os.Stderr, "錯誤: 必須指定 -action (read, write, random)\n\n")
		flag.Usage()
		os.Exit(1)
	}

	var client *fatek.FatekClient

	// 建立客戶端
	switch *mode {
	case "tcp":
		if *host == "" {
			fmt.Fprintf(os.Stderr, "錯誤: TCP 模式必須指定 -host\n\n")
			flag.Usage()
			os.Exit(1)
		}
		client = fatek.CreateTCPClient(*host, *port, *station, *timeout)
		fmt.Printf("建立 TCP 客戶端: %s:%d (站號: %d, 逾時: %v)\n", *host, *port, *station, *timeout)

	case "serial":
		if *serialPort == "" {
			fmt.Fprintf(os.Stderr, "錯誤: Serial 模式必須指定 -serial\n\n")
			flag.Usage()
			os.Exit(1)
		}
		client = fatek.CreateSerialClient(*serialPort, *station, *baudRate, *dataBits, *stopBits, *parity, *timeout)
		fmt.Printf("建立 Serial 客戶端: %s (站號: %d, 波特率: %d, 資料位元: %d, 停止位元: %d, 同位: %s, 逾時: %v)\n",
			*serialPort, *station, *baudRate, *dataBits, *stopBits, *parity, *timeout)

	default:
		fmt.Fprintf(os.Stderr, "錯誤: 不支援的模式 '%s'，請使用 'tcp' 或 'serial'\n", *mode)
		os.Exit(1)
	}

	// 連線
	fmt.Println("正在連線...")
	if err := client.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "連線失敗: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()
	fmt.Println("連線成功！")

	// 執行操作
	switch *action {
	case "read":
		if err := handleRead(client, *symbol, *addr, *count); err != nil {
			fmt.Fprintf(os.Stderr, "讀取失敗: %v\n", err)
			os.Exit(1)
		}

	case "write":
		if err := handleWrite(client, *symbol, *addr, *values); err != nil {
			fmt.Fprintf(os.Stderr, "寫入失敗: %v\n", err)
			os.Exit(1)
		}

	case "random":
		if err := handleRandom(client, *random); err != nil {
			fmt.Fprintf(os.Stderr, "隨機讀取失敗: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "錯誤: 不支援的操作 '%s'，請使用 'read', 'write' 或 'random'\n", *action)
		os.Exit(1)
	}
}

// handleRead 處理讀取操作
func handleRead(client *fatek.FatekClient, symbol string, addr, count int) error {
	if symbol == "" {
		return fmt.Errorf("必須指定 -symbol")
	}

	comp, err := fatek.GetComponentType(symbol)
	if err != nil {
		return fmt.Errorf("無效的組件符號 '%s': %w", symbol, err)
	}

	if comp.IsDiscrete {
		// 讀取離散狀態
		fmt.Printf("讀取離散狀態: %s%d-%d\n", symbol, addr, addr+count-1)
		values, err := client.ReadStatus(symbol, addr, count)
		if err != nil {
			return err
		}
		fmt.Printf("結果: [")
		for i, v := range values {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%v", v)
		}
		fmt.Println("]")
	} else {
		// 讀取暫存器
		fmt.Printf("讀取暫存器: %s%d-%d\n", symbol, addr, addr+count-1)
		values, err := client.ReadRegisters(symbol, addr, count)
		if err != nil {
			return err
		}
		fmt.Printf("結果: [")
		for i, v := range values {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%d", v)
		}
		fmt.Println("]")
	}

	return nil
}

// handleWrite 處理寫入操作
func handleWrite(client *fatek.FatekClient, symbol string, addr int, valuesStr string) error {
	if symbol == "" {
		return fmt.Errorf("必須指定 -symbol")
	}
	if valuesStr == "" {
		return fmt.Errorf("必須指定 -values")
	}

	comp, err := fatek.GetComponentType(symbol)
	if err != nil {
		return fmt.Errorf("無效的組件符號 '%s': %w", symbol, err)
	}

	valuesList := strings.Split(valuesStr, ",")
	if len(valuesList) == 0 {
		return fmt.Errorf("無效的值列表")
	}

	if comp.IsDiscrete {
		// 寫入離散狀態
		values := make([]bool, 0, len(valuesList))
		for _, v := range valuesList {
			v = strings.TrimSpace(v)
			b, err := strconv.ParseBool(v)
			if err != nil {
				return fmt.Errorf("無法解析布林值 '%s': %w", v, err)
			}
			values = append(values, b)
		}

		fmt.Printf("寫入離散狀態: %s%d-%d = %v\n", symbol, addr, addr+len(values)-1, values)
		if err := client.WriteStatus(symbol, addr, values); err != nil {
			return err
		}
		fmt.Println("寫入成功！")
	} else {
		// 寫入暫存器
		values := make([]int, 0, len(valuesList))
		for _, v := range valuesList {
			v = strings.TrimSpace(v)
			i, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("無法解析整數值 '%s': %w", v, err)
			}
			values = append(values, i)
		}

		fmt.Printf("寫入暫存器: %s%d-%d = %v\n", symbol, addr, addr+len(values)-1, values)
		if err := client.WriteRegisters(symbol, addr, values); err != nil {
			return err
		}
		fmt.Println("寫入成功！")
	}

	return nil
}

// handleRandom 處理隨機讀取操作
func handleRandom(client *fatek.FatekClient, randomStr string) error {
	if randomStr == "" {
		return fmt.Errorf("必須指定 -random")
	}

	items := make([]fatek.RandomReadItem, 0)
	parts := strings.Split(randomStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		idx := strings.LastIndex(part, ":")
		if idx < 0 {
			return fmt.Errorf("無效的隨機讀取格式 '%s'，應為 SYMBOL:ADDR", part)
		}

		symbol := part[:idx]
		addrStr := part[idx+1:]
		addr, err := strconv.Atoi(addrStr)
		if err != nil {
			return fmt.Errorf("無法解析位址 '%s': %w", addrStr, err)
		}

		items = append(items, fatek.RandomReadItem{
			Symbol: symbol,
			Addr:   addr,
		})
	}

	if len(items) == 0 {
		return fmt.Errorf("沒有指定任何讀取項目")
	}

	fmt.Printf("隨機讀取 %d 個項目:\n", len(items))
	for _, item := range items {
		fmt.Printf("  - %s%d\n", item.Symbol, item.Addr)
	}

	results, err := client.ReadRandom(items)
	if err != nil {
		return err
	}

	fmt.Println("\n結果:")
	for key, value := range results {
		fmt.Printf("  %s = %v\n", key, value)
	}

	return nil
}
