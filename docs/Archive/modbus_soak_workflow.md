# Modbus Soak Test Workflow（60s + 10m + 30m）

## 一鍵執行

```bash
scripts/run_modbus_soak.sh \
  --workdir /home/yishow/github/go_gateway \
  --report-dir /tmp/modbus_reports_latest
```

預設會依序跑：
1. smoke 60s
2. final 10m
3. soak 30m

也可用 Makefile 統一入口（Gate v1.2）：

```bash
make gate-smoke   # 只跑 smoke
make gate-final   # 只跑 final
make gate-soak    # 只跑 soak
make gatev11      # smoke + final + soak 全流程
```

每一階段都會做 gate check：
- `pass=true`
- `mismatch_count=0`
- `bad_rows=0`
- `total_rows>=1`

任一階段失敗會立即停止並（預設）送出 OpenClaw system event 通知。

---

## 可重複基準資料（deterministic）

腳本會把以下參數固定傳給 `cmd/loadtest_modbus`：
- `-value-base`（預設 1000）
- `-tag-key-prefix`（預設 `bench/server`）

可透過參數覆寫：

```bash
scripts/run_modbus_soak.sh \
  --value-base 1200 \
  --tag-key-prefix bench/stable \
  --report-dir /tmp/modbus_reports_stable
```

---

## 只跑特定階段

```bash
# 跳過 30m soak
scripts/run_modbus_soak.sh --skip-soak

# 只跑 30m soak（跳過 smoke + final）
scripts/run_modbus_soak.sh --skip-smoke --skip-final
```

---

## Dry-run 檢查命令

```bash
scripts/run_modbus_soak.sh --dry-run --skip-soak
```

---

## 報告檢查工具

手動檢查單一報告：

```bash
scripts/check_modbus_report.py /tmp/modbus_10m_report.json
```
