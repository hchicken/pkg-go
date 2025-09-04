# Go 环境变量
GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
GOVERSION:=$(shell go version | awk '{print $$3}')
VERSION:=$(shell git describe --tags --always)
COMMIT:=$(shell git rev-parse --short HEAD)
BUILD_TIME:=$(shell date '+%Y-%m-%d %H:%M:%S')

# 项目信息
PROJECT_NAME:=pkg-go
BUILD_DIR:=build
DIST_DIR:=dist

.PHONY: release clean test lint fmt vet deps check build info

# 生成发布版本
release:
	@read -p "Enter package name (or 'all' for all packages): " package; \
	if [ -z "$$package" ]; then \
		echo "Error: Package name cannot be empty"; \
		exit 1; \
	fi; \
	read -p "Enter version number (e.g., v1.0.0): " version; \
	if [ -z "$$version" ]; then \
		echo "Error: Version number cannot be empty"; \
		exit 1; \
	fi; \
	if ! echo "$$version" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+$$'; then \
		echo "Error: Version must follow semantic versioning (e.g., v1.0.0)"; \
		exit 1; \
	fi; \
	echo "Creating release for $$package version $$version"; \
	read -p "Continue? [y/N] " confirm; \
	if [ "$$confirm" != "y" ] && [ "$$confirm" != "Y" ]; then \
		echo "Release cancelled"; \
		exit 1; \
	fi; \
	./release.sh "$$package" "$$version"

# 清理构建文件
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR) $(DIST_DIR)
	@rm -f *.zip *.tar.gz
	@echo "Clean completed"

# 运行测试
test:
	@echo "Running tests..."
	@go test -v ./...

# 运行测试并生成覆盖率报告
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# 代码格式化
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# 代码检查
vet:
	@echo "Running go vet..."
	@go vet ./...

# 代码风格检查 (需要安装 golint)
lint:
	@echo "Running golint..."
	@if command -v golint >/dev/null 2>&1; then \
		golint ./...; \
	else \
		echo "golint not installed. Run: go install golang.org/x/lint/golint@latest"; \
	fi

# 下载依赖
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# 更新依赖
deps-update:
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy

# 综合检查
check: fmt vet lint test
	@echo "All checks completed"

# 构建所有包
build:
	@echo "Building packages..."
	@mkdir -p $(BUILD_DIR)
	@for dir in */; do \
		if [ -f "$$dir/go.mod" ] || [ -f "$$dir/*.go" ]; then \
			echo "Building $$dir..."; \
			(cd "$$dir" && go build -o "../$(BUILD_DIR)/$${dir%/}" .); \
		fi; \
	done

# 显示项目信息
info:
	@echo "Project Information:"
	@echo "  Name: $(PROJECT_NAME)"
	@echo "  Version: $(VERSION)"
	@echo "  Commit: $(COMMIT)"
	@echo "  Build Time: $(BUILD_TIME)"
	@echo "  Go Version: $(GOVERSION)"
	@echo "  Host OS: $(GOHOSTOS)"

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help