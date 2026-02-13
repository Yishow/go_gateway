package handlers

import (
"errors"
"fmt"
"net"
)

// categorizeError 將底層錯誤轉換為更易讀的錯誤訊息
func categorizeError(err error) error {
	if err == nil {
		return nil
	}
	
	// 處理 net.Error (Timeout, Connection refused)
	var netErr net.Error
	if errors.As(err, &netErr) {
		if netErr.Timeout() {
			return fmt.Errorf("TIMEOUT: %v", err)
		}
		return fmt.Errorf("NETWORK_ERROR: %v", err)
	}

	// 檢查常見的錯誤字串 (因為部分庫可能返回普通 error)
	s := err.Error()
	if s == "EOF" {
		return fmt.Errorf("CONNECTION_CLOSED: Remote host closed connection")
	}
	// TODO: 可以根據具體協議庫的錯誤類型進行更細緻的分類
	// 例如: CRC Checksum Error, Illegal Function, etc.

	return err
}
