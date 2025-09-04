#!/bin/bash

# 脚本配置
set -euo pipefail  # 严格模式：遇到错误立即退出，未定义变量报错，管道错误传播

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

# 帮助文档
usage() {
    cat << EOF
Usage: $0 <package_name|all> <version> [options]

Arguments:
  package_name    Package name to release, or 'all' for all packages
  version         Version number in semantic versioning format (e.g., v1.0.0)

Options:
  -h, --help      Show this help message
  -d, --dry-run   Show what would be done without actually doing it
  -f, --force     Force release even if tag already exists
  --no-archive    Skip creating archive files

Examples:
  $0 all v1.0.0
  $0 cache v1.2.3
  $0 all v2.0.0 --dry-run
  $0 cache v1.0.1 --force

EOF
    exit 1
}

# 参数解析
DRY_RUN=false
FORCE=false
NO_ARCHIVE=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            usage
            ;;
        -d|--dry-run)
            DRY_RUN=true
            shift
            ;;
        -f|--force)
            FORCE=true
            shift
            ;;
        --no-archive)
            NO_ARCHIVE=true
            shift
            ;;
        -*)
            log_error "Unknown option: $1"
            usage
            ;;
        *)
            if [[ -z "${release_pkg:-}" ]]; then
                release_pkg=$1
            elif [[ -z "${version:-}" ]]; then
                version=$1
            else
                log_error "Too many arguments"
                usage
            fi
            shift
            ;;
    esac
done

# 检查必需参数
if [[ -z "${release_pkg:-}" ]] || [[ -z "${version:-}" ]]; then
    log_error "Missing required arguments"
    usage
fi

log_info "Release package: $release_pkg"
log_info "Version: $version"
[[ "$DRY_RUN" == true ]] && log_warning "DRY RUN MODE - No actual changes will be made"

# 版本格式验证
if [[ -z "$version" ]]; then
    log_error "Version cannot be empty"
    exit 1
elif ! [[ "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    log_error "Version must match format v1.0.0"
    exit 1
fi

# 检查Git仓库状态
check_git_status() {
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        log_error "Not in a git repository"
        exit 1
    fi
    
    if [[ $(git status --porcelain) ]]; then
        log_warning "Working directory is not clean"
        git status --short
        read -p "Continue anyway? [y/N] " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "Release cancelled"
            exit 1
        fi
    fi
}

# 检查标签是否已存在
check_tag_exists() {
    local tag=$1
    if git rev-parse "$tag" >/dev/null 2>&1; then
        if [[ "$FORCE" == false ]]; then
            log_error "Tag '$tag' already exists. Use --force to overwrite"
            return 1
        else
            log_warning "Tag '$tag' already exists, will be overwritten"
        fi
    fi
    return 0
}

# 创建归档文件
create_archive() {
    local file_name=$1
    local version=$2
    local tag=$3
    local archive_name="pkg-go-${file_name}-${version}"
    
    if [[ "$NO_ARCHIVE" == true ]]; then
        log_info "Skipping archive creation (--no-archive specified)"
        return 0
    fi
    
    log_info "Creating archive: ${archive_name}.zip"
    if [[ "$DRY_RUN" == false ]]; then
        if ! git archive --format=zip --output="${archive_name}.zip" "$tag"; then
            log_error "Failed to create archive"
            return 1
        fi
        
        # 创建tar.gz格式的归档
        if ! git archive --format=tar.gz --output="${archive_name}.tar.gz" "$tag"; then
            log_error "Failed to create tar.gz archive"
            return 1
        fi
        
        log_success "Archives created: ${archive_name}.zip, ${archive_name}.tar.gz"
    fi
}

# 发布所有包
release_all_version() {
    local version=$1
    local failed_packages=()
    local successful_packages=()
    
    log_info "Releasing all packages with version $version"
    
    # 使用更兼容的方式查找目录
    for dir in */; do
        # 移除末尾的斜杠
        local file_name="${dir%/}"
        
        # 跳过隐藏目录和非目录文件
        [[ "$file_name" =~ ^\. ]] && continue
        [[ ! -d "$file_name" ]] && continue
        
        log_info "Processing package: $file_name"
        
        if release_single_version "$file_name" "$version"; then
            successful_packages+=("$file_name")
        else
            failed_packages+=("$file_name")
            log_error "Failed to release package: $file_name"
        fi
    done
    
    # 输出总结
    echo
    log_info "Release Summary:"
    if [[ ${#successful_packages[@]} -gt 0 ]]; then
        log_success "Successfully released: ${successful_packages[*]}"
    fi
    if [[ ${#failed_packages[@]} -gt 0 ]]; then
        log_error "Failed to release: ${failed_packages[*]}"
        return 1
    fi
    
    return 0
}

# 发布单个包
release_single_version() {
    local file_name=$1
    local version=$2
    local tag="${file_name}/${version}"
    
    log_info "Releasing package '$file_name' version '$version'"
    
    # 检查包目录是否存在
    if [[ ! -d "$file_name" ]]; then
        log_error "Package directory '$file_name' does not exist"
        return 1
    fi
    
    # 检查标签是否已存在
    if ! check_tag_exists "$tag"; then
        return 1
    fi
    
    # 创建标签
    log_info "Creating tag: $tag"
    if [[ "$DRY_RUN" == false ]]; then
        if [[ "$FORCE" == true ]] && git rev-parse "$tag" >/dev/null 2>&1; then
            log_info "Deleting existing tag: $tag"
            git tag -d "$tag" || {
                log_error "Failed to delete existing tag"
                return 1
            }
            git push origin ":refs/tags/$tag" 2>/dev/null || true
        fi
        
        if ! git tag -a "$tag" -m "Release $file_name $version"; then
            log_error "Failed to create tag"
            return 1
        fi
        log_success "Tag created: $tag"
    fi
    
    # 推送标签
    log_info "Pushing tag to origin"
    if [[ "$DRY_RUN" == false ]]; then
        if ! git push origin "$tag"; then
            log_error "Failed to push tag"
            return 1
        fi
        log_success "Tag pushed: $tag"
    fi
    
    # 创建归档
    if ! create_archive "$file_name" "$version" "$tag"; then
        return 1
    fi
    
    log_success "Successfully released $file_name version $version"
    return 0
}

# 主逻辑
main() {
    # 检查Git状态
    check_git_status
    
    if [[ "$release_pkg" == "all" ]]; then
        release_all_version "$version"
    else
        if [[ -d "$release_pkg" ]]; then
            release_single_version "$release_pkg" "$version"
        else
            log_error "Unknown package name: $release_pkg"
            log_info "Available packages:"
            for dir in */; do
                [[ -d "$dir" ]] && echo "  - ${dir%/}"
            done
            exit 1
        fi
    fi
}

# 执行主函数
main "$@"