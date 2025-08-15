# Makefile for building installation for multiple platforms

# 變數定義
BINARY_NAME=installation
VERSION?=v1.0.0
BUILD_DIR=build
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

# Go 相關變數
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod

# 預設目標
.PHONY: all
all: clean deps build-all

# 初始化和安裝依賴
.PHONY: deps
deps:
	$(GOMOD) tidy

# 清理構建文件
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# 創建構建目錄
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# 構建所有平台
.PHONY: build-all
build-all: build-darwin build-linux build-windows

# macOS 構建
.PHONY: build-darwin
build-darwin: build-darwin-amd64 build-darwin-arm64

.PHONY: build-darwin-amd64
build-darwin-amd64: $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .

.PHONY: build-darwin-arm64
build-darwin-arm64: $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .

# Linux 構建
.PHONY: build-linux
build-linux: build-linux-amd64 build-linux-arm64

.PHONY: build-linux-amd64
build-linux-amd64: $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .

.PHONY: build-linux-arm64
build-linux-arm64: $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .

# Windows 構建
.PHONY: build-windows
build-windows: build-windows-amd64

.PHONY: build-windows-amd64
build-windows-amd64: $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .

# 本地構建（當前平台）
.PHONY: local
local: $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

# 測試
.PHONY: test
test:
	$(GOTEST) -v ./...

# 格式化代碼
.PHONY: fmt
fmt:
	go fmt ./...

# 檢查代碼
.PHONY: vet
vet:
	go vet ./...

# 壓縮二進制文件
.PHONY: compress
compress: build-all
	@echo "壓縮二進制文件..."
	@for file in $(BUILD_DIR)/*; do \
		if [ -f "$$file" ]; then \
			echo "壓縮 $$file"; \
			gzip -f "$$file"; \
		fi \
	done

# 創建發布包
.PHONY: release
release: build-all
	@echo "創建發布包..."
	@cd $(BUILD_DIR) && \
	for file in *; do \
		if [ -f "$$file" ]; then \
			echo "打包 $$file"; \
			tar -czf "$$file.tar.gz" "$$file"; \
		fi \
	done

# 顯示構建信息
.PHONY: info
info:
	@echo "二進制名稱: $(BINARY_NAME)"
	@echo "版本: $(VERSION)"
	@echo "構建目錄: $(BUILD_DIR)"

# 安裝到本地 GOPATH/bin
.PHONY: install
install:
	$(GOCMD) install $(LDFLAGS) .

# 快速構建和測試
.PHONY: quick
quick: fmt vet local test

# 幫助信息
.PHONY: help
help:
	@echo "可用的 make 目標:"
	@echo "  all          - 清理、安裝依賴並構建所有平台"
	@echo "  local        - 構建當前平台的二進制文件"
	@echo "  build-all    - 構建所有平台的二進制文件"
	@echo "  build-darwin - 構建 macOS 平台 (amd64, arm64)"
	@echo "  build-linux  - 構建 Linux 平台 (amd64, arm64)"
	@echo "  build-windows- 構建 Windows 平台 (amd64)"
	@echo "  clean        - 清理構建文件"
	@echo "  deps         - 安裝 Go 依賴"
	@echo "  test         - 運行測試"
	@echo "  fmt          - 格式化代碼"
	@echo "  vet          - 檢查代碼"
	@echo "  compress     - 壓縮二進制文件"
	@echo "  release      - 創建發布包"
	@echo "  install      - 安裝到本地"
	@echo "  quick        - 快速構建和測試"
	@echo "  info         - 顯示構建信息"
	@echo "  help         - 顯示此幫助信息"