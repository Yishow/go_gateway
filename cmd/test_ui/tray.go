package main

import (
	"log"
	"os/exec"
	"runtime"

	"github.com/getlantern/systray"
)

// SetupTray 設定系統托盤圖示和選單
//
// Args:
//   - serverAddr: 伺服器地址（用於打開瀏覽器）
func SetupTray(serverAddr string) {
	// 設定托盤圖示和選單
	systray.SetIcon(getIconData())
	systray.SetTitle("Go Gateway 測試工具")
	systray.SetTooltip("Go Gateway 測試工具 - 點擊打開瀏覽器")

	// 建立選單項目
	mOpenBrowser := systray.AddMenuItem("打開瀏覽器", "在瀏覽器中打開應用程式")
	mOpenBrowser.SetIcon(getIconData())
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("退出", "退出應用程式")

	// 在背景執行選單事件處理
	go func() {
		for {
			select {
			case <-mOpenBrowser.ClickedCh:
				openBrowser(serverAddr)
			case <-mQuit.ClickedCh:
				log.Println("收到退出請求，正在關閉應用程式...")
				systray.Quit()
				return
			}
		}
	}()
}

// openBrowser 在預設瀏覽器中打開指定 URL
//
// Args:
//   - url: 要打開的 URL
func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default: // linux
		cmd = exec.Command("xdg-open", url)
	}

	if err := cmd.Start(); err != nil {
		log.Printf("無法打開瀏覽器: %v", err)
	}
}

// getIconData 返回系統托盤圖示的位元組資料
// 返回一個簡單的 16x16 藍色方塊圖示作為預設圖示
// 在實際使用中，您可以使用 embed 來嵌入自己的圖示檔案
//
// Returns:
//   - []byte: 圖示的位元組資料（16x16 PNG）
func getIconData() []byte {
	// 這是一個 16x16 的藍色方塊 PNG 圖示
	// 在實際使用中，您應該替換為自己的圖示檔案
	// 可以使用 embed 來嵌入圖示檔案，例如：
	//   //go:embed icon.png
	//   var iconData embed.FS
	//   然後讀取 iconData
	return []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D,
		0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x10,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0xF3, 0xFF, 0x61, 0x00, 0x00, 0x00,
		0x19, 0x49, 0x44, 0x41, 0x54, 0x38, 0x8D, 0x63, 0x64, 0x60, 0x60, 0x60,
		0xF8, 0xCF, 0xC0, 0xC0, 0xF0, 0x9F, 0x81, 0x81, 0xE1, 0x3F, 0x03, 0x03,
		0xC3, 0x7F, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42,
		0x60, 0x82,
	}
}
