.PHONY: build-frontend sync-frontend-static build-backend build clean test-ui gen-docs gate-smoke gate-final gate-soak gatev11 gatev12 points-precheck-up points-precheck-down points-migrate-up points-migrate-down longtask-smoke cross-platform-loop check-lines

GATE_WORKDIR ?= $(CURDIR)
GATE_SMOKE_DURATION ?= 60s
GATE_FINAL_DURATION ?= 10m
GATE_SOAK_DURATION ?= 30m

POINTS_DB_TYPE ?= postgres
POINTS_DSN ?=
POINTS_SQLITE_FILE ?=
POINTS_DOWN_STRATEGY ?= keep-latest


# 建置前端
build-frontend:
	cd frontend && npm install
	cd frontend && npm run build

# 同步前端 build 到 embed 目錄
sync-frontend-static: build-frontend
	mkdir -p cmd/test_ui/static
	cp -R frontend/dist/. cmd/test_ui/static/

# 建置後端（包含前端）
build-backend: sync-frontend-static
	go build -o bin/test-ui.exe ./cmd/test_ui

# 建置測試工具
build: build-backend

# 清理建置檔案
clean:
	rm -rf frontend/dist
	rm -rf frontend/node_modules
	rm -rf bin/

# 執行測試工具
test-ui: build
	./bin/test-ui.exe

# 產生 Swagger API 文檔
gen-docs:
	swag init -g cmd/test_ui/main.go -o docs/swagger

# 僅執行 smoke gate
gate-smoke:
	./scripts/run_modbus_soak.sh --workdir "$(GATE_WORKDIR)" --smoke-duration "$(GATE_SMOKE_DURATION)" --final-duration "$(GATE_FINAL_DURATION)" --soak-duration "$(GATE_SOAK_DURATION)" --skip-final --skip-soak

# 僅執行 final gate
gate-final:
	./scripts/run_modbus_soak.sh --workdir "$(GATE_WORKDIR)" --smoke-duration "$(GATE_SMOKE_DURATION)" --final-duration "$(GATE_FINAL_DURATION)" --soak-duration "$(GATE_SOAK_DURATION)" --skip-smoke --skip-soak

# 僅執行 soak gate
gate-soak:
	./scripts/run_modbus_soak.sh --workdir "$(GATE_WORKDIR)" --smoke-duration "$(GATE_SMOKE_DURATION)" --final-duration "$(GATE_FINAL_DURATION)" --soak-duration "$(GATE_SOAK_DURATION)" --skip-smoke --skip-final

# Gate v1.1 全流程（smoke + final + soak）
gatev11:
	./scripts/run_modbus_soak.sh --workdir "$(GATE_WORKDIR)" --smoke-duration "$(GATE_SMOKE_DURATION)" --final-duration "$(GATE_FINAL_DURATION)" --soak-duration "$(GATE_SOAK_DURATION)"

# Gate v1.2 預設驗證（smoke + final）
gatev12:
	./scripts/run_modbus_soak.sh --workdir "$(GATE_WORKDIR)" --smoke-duration "$(GATE_SMOKE_DURATION)" --final-duration "$(GATE_FINAL_DURATION)" --soak-duration "$(GATE_SOAK_DURATION)" --skip-soak

# Points unique migration precheck (up)
points-precheck-up:
	@if [ "$(POINTS_DB_TYPE)" = "postgres" ]; then \
		./scripts/check_points_unique_conflicts.sh --db postgres --stage up --dsn "$(POINTS_DSN)"; \
	else \
		./scripts/check_points_unique_conflicts.sh --db sqlite --stage up --sqlite-file "$(POINTS_SQLITE_FILE)"; \
	fi

# Points unique migration precheck (down)
points-precheck-down:
	@if [ "$(POINTS_DB_TYPE)" = "postgres" ]; then \
		./scripts/check_points_unique_conflicts.sh --db postgres --stage down --dsn "$(POINTS_DSN)"; \
	else \
		./scripts/check_points_unique_conflicts.sh --db sqlite --stage down --sqlite-file "$(POINTS_SQLITE_FILE)"; \
	fi

# Points unique migration apply up
points-migrate-up:
	@if [ "$(POINTS_DB_TYPE)" = "postgres" ]; then \
		./scripts/run_points_unique_migration.sh --db postgres --action up --dsn "$(POINTS_DSN)" --execute; \
	else \
		./scripts/run_points_unique_migration.sh --db sqlite --action up --sqlite-file "$(POINTS_SQLITE_FILE)" --execute; \
	fi

# Points unique migration apply down (requires explicit confirm)
points-migrate-down:
	@if [ "$(POINTS_DB_TYPE)" = "postgres" ]; then \
		./scripts/run_points_unique_migration.sh --db postgres --action down --dsn "$(POINTS_DSN)" --execute --confirm-down --down-conflict-strategy "$(POINTS_DOWN_STRATEGY)"; \
	else \
		./scripts/run_points_unique_migration.sh --db sqlite --action down --sqlite-file "$(POINTS_SQLITE_FILE)" --execute --confirm-down --down-conflict-strategy "$(POINTS_DOWN_STRATEGY)"; \
	fi


# 長任務提醒機制 smoke test（controller start/finish + reconcile）
longtask-smoke:
	bash ./scripts/longtask_smoke.sh

# 跨平台測試契約驗證序列（Vite x2 → handler → row-group -count=20 → 全域）
cross-platform-loop:
	bash ./scripts/run_cross_platform_contract_loop.sh

# 檔案行數規範檢查（預設檢查本次變更）
check-lines:
	bash ./scripts/check_file_lines.sh
