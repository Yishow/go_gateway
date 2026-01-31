.PHONY: build-frontend build-backend build clean test-ui gen-docs

# 建置前端
build-frontend:
	cd frontend && npm install
	cd frontend && npm run build

# 建置後端（包含前端）
build-backend: build-frontend
	go build -o bin/test-ui.exe ./cmd/test_ui

# 建置測試工具
build: build-backend

# 清理建置檔案
clean:
	rm -rf frontend/dist
	rm -rf frontend/node_modules
	rm -f bin/test-ui.exe

# 執行測試工具
test-ui: build
	./bin/test-ui.exe

# 產生 Swagger API 文檔
gen-docs:
	swag init -g cmd/test_ui/main.go -o docs/swagger
