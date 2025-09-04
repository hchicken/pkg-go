# 项目信息
PROJECT_NAME:=pkg-go
VERSION:=$(shell git describe --tags --always)

.PHONY: release clean get-version help

# 获取包的最新版本
get-version:
	@read -p "Enter package name: " package; \
	if [ -z "$$package" ]; then \
		echo "Error: Package name cannot be empty"; \
		exit 1; \
	fi; \
	echo "Getting latest version for package: $$package"; \
	latest_tag=$$(git tag -l "$$package/*" --sort=-version:refname | head -1); \
	if [ -z "$$latest_tag" ]; then \
		echo "No version found for package: $$package"; \
		echo "Available packages:"; \
		git tag -l | cut -d'/' -f1 | sort -u | sed 's/^/  - /'; \
	else \
		echo "Latest version: $$latest_tag"; \
	fi

# 打包发布
release:
	@echo "Starting release process..."
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
	@rm -rf build dist *.zip *.tar.gz
	@echo "Clean completed"

# 显示帮助信息
help:
	@echo "$(PROJECT_NAME) Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  release      - Create and publish a package release"
	@echo "  get-version  - Get the latest version of a package"
	@echo "  clean        - Clean build artifacts"
	@echo "  help         - Show this help message"
	@echo ""
	@echo "Usage:"
	@echo "  make release      # Interactive release process"
	@echo "  make get-version  # Get latest version of a package"
	@echo "  make clean        # Clean up build files"

.DEFAULT_GOAL := help