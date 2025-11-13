#!/usr/bin/env bash

###############################################################################
# Docker 本地构建测试脚本
# 用途：在本地测试 Docker 镜像构建流程
# 作者：robot-shop team
# 日期：2025-11-12
###############################################################################

set -e  # 遇到错误立即退出
set -u  # 使用未定义变量时报错

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
    echo -e "${RED}[ERROR]${NC} $1"
}

# 配置变量
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
IMAGE_NAME="${IMAGE_NAME:-robot-shop}"
IMAGE_TAG="${IMAGE_TAG:-local-test}"
DOCKERFILE_PATH="${PROJECT_ROOT}/deploy/build/Dockerfile"
BUILD_CONTEXT="${PROJECT_ROOT}"

# 架构选择
PLATFORM="${PLATFORM:-linux/amd64}"  # 默认构建当前平台，可选：linux/amd64,linux/arm64
USE_CHINA_MIRROR="${USE_CHINA_MIRROR:-false}"  # 是否使用国内镜像加速

# 显示配置信息
show_config() {
    log_info "===== 构建配置 ====="
    echo "  项目根目录: ${PROJECT_ROOT}"
    echo "  Dockerfile: ${DOCKERFILE_PATH}"
    echo "  构建上下文: ${BUILD_CONTEXT}"
    echo "  镜像名称: ${IMAGE_NAME}:${IMAGE_TAG}"
    echo "  目标平台: ${PLATFORM}"
    echo "  国内加速: ${USE_CHINA_MIRROR}"
    echo ""
}

# 检查依赖
check_dependencies() {
    log_info "检查依赖环境..."
    
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi
    
    # 检查 Docker 是否运行
    if ! docker info &> /dev/null; then
        log_error "Docker 未运行，请启动 Docker"
        exit 1
    fi
    
    log_success "依赖检查通过"
}

# 检查必要文件
check_files() {
    log_info "检查必要文件..."
    
    local files=(
        "${DOCKERFILE_PATH}"
        "${PROJECT_ROOT}/frontend/package.json"
        "${PROJECT_ROOT}/backend/go.mod"
    )
    
    for file in "${files[@]}"; do
        if [[ ! -f "$file" ]]; then
            log_error "文件不存在: ${file}"
            exit 1
        fi
    done
    
    log_success "文件检查通过"
}

# 清理旧镜像（可选）
clean_old_images() {
    log_info "清理旧的测试镜像..."
    
    if docker images "${IMAGE_NAME}:${IMAGE_TAG}" | grep -q "${IMAGE_TAG}"; then
        docker rmi "${IMAGE_NAME}:${IMAGE_TAG}" || true
        log_success "已清理旧镜像"
    else
        log_info "没有需要清理的旧镜像"
    fi
}

# 构建镜像
build_image() {
    log_info "开始构建 Docker 镜像..."
    log_info "这可能需要几分钟时间，请耐心等待..."
    
    local build_start=$(date +%s)
    
    # 构建命令
    local build_cmd="docker build \
        --platform ${PLATFORM} \
        --build-arg USE_CHINA_MIRROR=${USE_CHINA_MIRROR} \
        -t ${IMAGE_NAME}:${IMAGE_TAG} \
        -f ${DOCKERFILE_PATH} \
        ${BUILD_CONTEXT}"
    
    log_info "执行命令: ${build_cmd}"
    echo ""
    
    if eval "${build_cmd}"; then
        local build_end=$(date +%s)
        local build_duration=$((build_end - build_start))
        log_success "镜像构建成功！耗时: ${build_duration} 秒"
    else
        log_error "镜像构建失败"
        exit 1
    fi
}

# 检查镜像
verify_image() {
    log_info "验证镜像..."
    
    if docker images "${IMAGE_NAME}:${IMAGE_TAG}" | grep -q "${IMAGE_TAG}"; then
        local image_info=$(docker images "${IMAGE_NAME}:${IMAGE_TAG}" --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}")
        echo ""
        echo "${image_info}"
        echo ""
        log_success "镜像验证通过"
    else
        log_error "镜像不存在"
        exit 1
    fi
}

# 显示镜像详情
show_image_details() {
    log_info "镜像详细信息..."
    docker inspect "${IMAGE_NAME}:${IMAGE_TAG}" --format='
镜像ID: {{.Id}}
创建时间: {{.Created}}
架构: {{.Architecture}}
操作系统: {{.Os}}
大小: {{.Size}} bytes
层数: {{len .RootFS.Layers}}
'
}

# 主函数
main() {
    echo ""
    log_info "===== Docker 本地构建测试 ====="
    echo ""
    
    show_config
    check_dependencies
    check_files
    
    # 询问是否清理旧镜像
    if [[ "${CLEAN:-}" == "true" ]]; then
        clean_old_images
    fi
    
    build_image
    verify_image
    show_image_details
    
    echo ""
    log_success "===== 构建测试完成 ====="
    echo ""
    log_info "下一步："
    log_info "  1. 运行容器测试: ./deploy/scripts/run-local.sh"
    log_info "  2. 或手动运行: docker run -p 8000:8000 ${IMAGE_NAME}:${IMAGE_TAG}"
    echo ""
}

# 执行主函数
main "$@"
