package main

import (
	"context"
	"time"

	"go-gateway/internal/apphost"
	"go-gateway/internal/desktop"
	"go-gateway/internal/diagnostics"
)

// reportAsync keeps a visible native error without delaying resource cleanup
// until the operator acknowledges its dialog. The caller waits after cleanup.
func (h *gatewayHost) reportAsync(err error) <-chan struct{} {
	done := make(chan struct{})
	if err == nil {
		close(done)
		return done
	}
	code := apphost.Code(err)
	h.emit(code, nil)
	saved := false
	if h.logs != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		saved = h.logs.Flush(ctx) == nil
		cancel()
	}
	if h.launch.Mode != apphost.Desktop {
		close(done)
		return done
	}
	info := desktop.ErrorInfo{Code: code, Message: "Gateway 無法完成這項操作。", NextStep: nativeNextStep(code)}
	if saved && h.sink != nil {
		result := make(chan diagnostics.FileStatus, 1)
		go func() { result <- h.sink.Status() }()
		select {
		case status := <-result:
			info.DiagnosticSaved = status.Saved
			info.DiagnosticPath = status.Path
		case <-time.After(time.Second):
		}
	}
	go func() { desktop.ShowError(info); close(done) }()
	return done
}
func nativeNextStep(code string) string {
	switch code {
	case "startup.path_ambiguous":
		return "工作目錄與程式目錄可能有不同資料來源。請以絕對 GATEWAY_DB_PATH（相容 DB_PATH、SQLITE_PATH）指定預期資料庫後再啟動；未搬移資料。"
	case "startup.database_linked":
		return "資料庫有多個 hardlink 名稱，無法確認同一 WAL/SHM 復原來源。請保留原始 DB 與 journal/WAL/SHM，由資料管理者辨識原始路徑及完整復原狀態；不要搬移、刪除 sidecar 或切換 alias 規避檢查。"
	case "startup.path_unusable":
		return "請確認資料庫及父目錄的存取權限，使用可信的絕對資料庫路徑；不會改用另一個資料库。"
	case "startup.owner_busy":
		return "相同資料庫已有Gateway owner。請使用已開啟的程序；headless或無法驗證的owner不會開啟網頁。"
	case "startup.bind_failed":
		return "請檢查HOST / PORT是否可用。Gateway不會終止其他程序或自動換連接埠。"
	case "startup.assets_missing":
		return "執行檔缺少完整網頁資產。請以完整前端與後端build重新建置。"
	case "startup.log_unavailable":
		return "請檢查LOG_FILE與資料庫/其他檔案是否衝突，以及診斷目錄寫入權限。"
	case "shutdown.failed":
		return "關閉未被確認為完整成功。保留現有資料與診斷，檢查未完成或結果未知的工作。"
	default:
		return "請檢查設定、資料目錄權限及診斷紀錄，修正後再啟動。"
	}
}
