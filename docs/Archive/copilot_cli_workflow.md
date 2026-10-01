# Copilot CLI Workflow（固定模板 + Watchdog）

## 1) 先準備任務 prompt

複製模板：

```bash
cp docs/templates/copilot_task_prompt.template.md /tmp/task.prompt.md
```

把 `{{...}}` 區塊換成你的任務內容。

---

## 2) 用 watchdog 執行（避免長任務卡死）

```bash
scripts/copilot_watchdog.sh \
  --workdir /home/yishow/github/go_gateway \
  --prompt-file /tmp/task.prompt.md \
  --model gpt-5.3-codex \
  --idle-minutes 8 \
  --max-restarts 2 \
  --poll-seconds 15 \
  --log-file /tmp/copilot_go_gateway.log
```

### 參數說明
- `--idle-minutes`: 超過 N 分鐘無新 log，視為卡住並重啟
- `--max-restarts`: 最多自動重啟次數
- `--poll-seconds`: 監控 log 的輪詢秒數
- `--log-file`: 執行日誌輸出位置

---

## 3) 建議規範（避免只讀不跑）

在 prompt 內固定要求：
1. 明確「目標 / 交付物 / 驗證命令」
2. 一定要跑測試與壓測（不可只靜態分析）
3. 輸出固定格式（檔案、diff、命令、結果、風險）
4. 完成時執行 `openclaw system event --mode now`

---

## 4) Go 無法直接執行時

在 prompt 內要求使用 Docker `golang:1.25`，並掛載 module/build cache：

```bash
docker run --rm \
  -v "$PWD":/workspace \
  -v /home/yishow/.cache/go-mod:/go/pkg/mod \
  -v /home/yishow/.cache/go-build:/root/.cache/go-build \
  -w /workspace \
  -e GOMODCACHE=/go/pkg/mod \
  -e GOPROXY=off \
  golang:1.25 /usr/local/go/bin/go test ./...
```


---

## 5) Soak 測試自動化（含 gate + 失敗通知）

可直接使用：

```bash
scripts/run_modbus_soak.sh \
  --workdir /home/yishow/github/go_gateway \
  --report-dir /tmp/modbus_reports_latest
```

細節見：`docs/modbus_soak_workflow.md`
