package modbus

import (
	"context"
	"encoding/binary"
	"net"
	"testing"
	"time"
)

// 測試 TCPTransport 在對端斷線後的自動自癒重連機制
func TestTCPTransport_SelfHealing_Reconnect(t *testing.T) {
	// 建立本機 Modbus TCP 模擬服務端
	lc := net.ListenConfig{}
	ln, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("無法建立測試監聽: %v", err)
	}
	defer ln.Close()

	addr := ln.Addr().(*net.TCPAddr)

	// 模擬服務端協程
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				buf := make([]byte, 256)
				for {
					_ = c.SetReadDeadline(time.Now().Add(2 * time.Second))
					n, err := c.Read(buf)
					if err != nil {
						return
					}
					if n >= 12 {
						// 正常回應 Modbus 保持暫存器讀取
						txID := binary.BigEndian.Uint16(buf[0:2])
						unitID := buf[6]
						resp := BuildTCPFrame(txID, unitID, FuncReadHoldingRegisters, []byte{2, 0x12, 0x34})
						_, _ = c.Write(resp)
						// 回應一次後主動斷開連線，模擬 PLC 閒置超時或強制切斷
						return
					}
				}
			}(conn)
		}
	}()

	transport := NewTCPTransport(addr.IP.String(), addr.Port)
	transport.Timeout = 1 * time.Second
	defer transport.Close()

	if err := transport.Connect(); err != nil {
		t.Fatalf("第一次連線失敗: %v", err)
	}
	if !transport.IsConnected() {
		t.Fatal("連線後 IsConnected 應為 true")
	}

	req := BuildTCPFrame(1, 1, FuncReadHoldingRegisters, []byte{0, 0, 0, 1})

	// 第一次請求（連線成功，服務端回應後會立即主動關閉連線）
	resp1, err := transport.SendReceive(req)
	if err != nil {
		t.Fatalf("第一次請求失敗: %v", err)
	}
	if len(resp1) < 9 {
		t.Fatalf("第一次回應長度過短: %d", len(resp1))
	}

	// 稍作等待確保對端 FIN 到達本地 socket
	time.Sleep(50 * time.Millisecond)

	// 第二次請求：此時底層 socket 已被對端斷開，SendReceive 應能自動自癒重新建立連線並成功取得回應！
	req2 := BuildTCPFrame(2, 1, FuncReadHoldingRegisters, []byte{0, 0, 0, 1})
	resp2, err := transport.SendReceive(req2)
	if err != nil {
		t.Fatalf("第二次請求自癒失敗: %v", err)
	}
	if len(resp2) < 9 {
		t.Fatalf("第二次回應長度過短: %d", len(resp2))
	}
}

func TestTCPTransport_IsConnected(t *testing.T) {
	transport := NewTCPTransport("127.0.0.1", 502)
	if transport.IsConnected() {
		t.Fatal("未 Connect 前 IsConnected 應為 false")
	}
}
