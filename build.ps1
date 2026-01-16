# 建置腳本
Write-Host "建置前端..." -ForegroundColor Green
Set-Location web/test-ui
npm run build
Set-Location ../..

Write-Host "複製前端檔案到 embed 目錄..." -ForegroundColor Green
if (Test-Path cmd/test_ui/static) {
    Remove-Item -Recurse -Force cmd/test_ui/static
}
Copy-Item -Recurse web/test-ui/dist cmd/test_ui/static

Write-Host "建置後端..." -ForegroundColor Green
# -s: 移除符號表，-w: 移除 DWARF 除錯資訊，-trimpath: 移除檔案路徑資訊
go build -ldflags "-H=windowsgui -s -w" -trimpath -o bin/test-ui.exe ./cmd/test_ui

Write-Host "建置完成！" -ForegroundColor Green
