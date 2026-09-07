package adapters

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"testing"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/modbus"
)

// TestModbusTCPConnector_SelfHealing_Read 測試 Modbus TCP 連接器在設備中斷連線後的自動自癒讀取
func TestModbusTCPConnector_SelfHealing_Read(t *testing.T) {
	lc := net.ListenConfig{}
	ln, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("無法建立測試監聽: %v", err)
	}
	defer ln.Close()

	addr := ln.Addr().(*net.TCPAddr)

	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			go func(conn net.Conn) {
				defer conn.Close()
				buf := make([]byte, 256)
				for {
					_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
					n, err := conn.Read(buf)
					if err != nil {
						return
					}
					if n >= 12 {
						txID := binary.BigEndian.Uint16(buf[0:2])
						unitID := buf[6]
						resp := modbus.BuildTCPFrame(txID, unitID, modbus.FuncReadHoldingRegisters, []byte{2, 0x12, 0x34})
						_, _ = conn.Write(resp)
						// 模擬 PLC idle timeout 主動切斷連線
						return
					}
				}
			}(c)
		}
	}()

	conn := NewModbusTCPConnector().(*ModbusTCPConnector)
	cfg := fmt.Sprintf(`{"host":%q,"port":%d,"slave_id":1,"timeout":1}`, addr.IP.String(), addr.Port)

	ctx := context.Background()
	if err := conn.Connect(ctx, cfg); err != nil {
		t.Fatalf("Connect 失敗: %v", err)
	}
	defer conn.Close()

	req := connector.ReadRequest{
		Address:  "40001",
		Function: "03",
		DataType: schema.DataTypeUint16,
		Count:    1,
	}

	// 第一次讀取（正常建立並讀取，服務端隨後切斷）
	res1, err := conn.Read(ctx, req)
	if err != nil {
		t.Fatalf("第一次讀取失敗: %v", err)
	}
	if res1.Quality != schema.QualityGood {
		t.Fatalf("第一次讀取品質異常: %v", res1.Quality)
	}

	// 等待對端 FIN 抵達
	time.Sleep(50 * time.Millisecond)

	// 第二次讀取：此時對端連線已切斷，Read 應能自動自癒重新連線並成功讀取
	res2, err := conn.Read(ctx, req)
	if err != nil {
		t.Fatalf("第二次讀取自癒失敗: %v", err)
	}
	if res2.Quality != schema.QualityGood {
		t.Fatalf("第二次讀取品質異常: %v", res2.Quality)
	}
}

// TestMC3EConnector_SelfHealing_EnsureConnection 測試 MC3E 連接器在連線斷開後的自癒
func TestMC3EConnector_SelfHealing_EnsureConnection(t *testing.T) {
	conn := &MC3EConnector{
		connected:      true,
		persistentMode: true,
	}

	// client 為 nil 時 ensureConnection 不應崩潰且返回 nil（既有行為）
	if err := conn.ensureConnection(); err != nil {
		t.Fatalf("client 為 nil 時不應報錯: %v", err)
	}
}
