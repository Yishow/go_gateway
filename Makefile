.PHONY: build-frontend build-backend build clean test-ui

# 建置前端
build-frontend:
	cd web/test-ui && npm install
	cd web/test-ui && npm run build

# 建置後端（包含前端）
build-backend: build-frontend
	go build -o bin/test-ui.exe ./cmd/test_ui

# 建置測試工具
build: build-backend

# 清理建置檔案
clean:
	rm -rf web/test-ui/dist
	rm -rf web/test-ui/node_modules
	rm -f bin/test-ui.exe

# 執行測試工具
test-ui: build
	./bin/test-ui.exe
