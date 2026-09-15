package hsllogic

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// =============================================================================
// 內部方法
// =============================================================================

// log 記錄報文
func (pl *PacketLogger) log(
	direction PacketDirection,
	protocol, deviceID string,
	data []byte,
	desc string,
	duration time.Duration,
) {
	if !pl.options.Enabled {
		return
	}

	// 從物件池取得 PacketLog
	entry := acquirePacketLog()
	entry.Timestamp = time.Now()
	entry.Direction = direction
	entry.Protocol = protocol
	entry.DeviceID = deviceID
	entry.RawData = data
	// HexString 延遲生成，僅在實際需要時（String() 被調用時）才格式化
	entry.Description = desc
	entry.Duration = duration

	// 控制台輸出
	if pl.options.LogToConsole && pl.options.ConsoleWriter != nil {
		fmt.Fprintln(pl.options.ConsoleWriter, entry.String())
	}

	// 寫入緩衝區
	if pl.options.LogToFile {
		select {
		case pl.buffer <- entry:
		default:
			// 緩衝區已滿，丟棄最舊的條目
			select {
			case <-pl.buffer:
				pl.buffer <- entry
			default:
			}
		}
	} else {
		// 如果不寫入檔案，立即歸還物件至池
		releasePacketLog(entry)
	}
}

// writeLoop 寫入迴圈
func (pl *PacketLogger) writeLoop() {
	defer pl.wg.Done()

	ticker := time.NewTicker(pl.options.FlushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-pl.stopChan:
			// 處理剩餘的緩衝區
			for len(pl.buffer) > 0 {
				entry := <-pl.buffer
				pl.writeEntry(entry)
			}
			return
		case entry := <-pl.buffer:
			pl.writeEntry(entry)
		case <-ticker.C:
			pl.Flush()
		}
	}
}

// writeEntry 寫入單筆條目
func (pl *PacketLogger) writeEntry(entry *PacketLog) {
	// 確保在函數結束時歸還物件至池
	defer releasePacketLog(entry)

	pl.mu.Lock()
	defer pl.mu.Unlock()

	// 確保檔案已開啟
	if err := pl.ensureFile(); err != nil {
		return
	}

	// 寫入日誌
	line := entry.String() + "\n"
	n, err := pl.currentFile.WriteString(line)
	if err != nil {
		return
	}

	pl.currentSize += int64(n)

	// 檢查是否需要輪替
	if pl.currentSize >= pl.options.MaxFileSize {
		pl.rotateFile()
	}
}

// ensureFile 確保日誌檔案已開啟
func (pl *PacketLogger) ensureFile() error {
	if pl.currentFile != nil {
		return nil
	}

	// 建立目錄
	if err := os.MkdirAll(pl.options.LogDir, 0o755); err != nil {
		return err
	}

	// 開啟新檔案
	filename := filepath.Join(
		pl.options.LogDir,
		fmt.Sprintf("packet_%s.log", time.Now().Format("20060102_150405")),
	)

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	pl.currentFile = f
	pl.currentSize = 0
	return nil
}

// rotateFile 輪替檔案
func (pl *PacketLogger) rotateFile() {
	if pl.currentFile != nil {
		pl.currentFile.Close()
		pl.currentFile = nil
	}

	pl.fileIndex++
	pl.cleanOldFiles()
}

// cleanOldFiles 清理舊檔案
func (pl *PacketLogger) cleanOldFiles() {
	files, err := filepath.Glob(filepath.Join(pl.options.LogDir, "packet_*.log"))
	if err != nil {
		return
	}

	if len(files) <= pl.options.MaxFiles {
		return
	}

	// 刪除最舊的檔案
	for i := 0; i < len(files)-pl.options.MaxFiles; i++ {
		os.Remove(files[i])
	}
}
