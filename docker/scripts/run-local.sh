#!/usr/bin/env bash

###############################################################################
# Docker 容器运行测试脚本
# 用途：测试构建的 Docker 镜像是否能正常运行
# 作者：robot-shop team
# 日期：2025-11-12
###############################################################################

set -e
set -u

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

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
CONTAINER_NAME="${CONTAINER_NAME:-robot-shop-test}"
HOST_PORT="${HOST_PORT:-8000}"
CONTAINER_PORT="${CONTAINER_PORT:-8000}"

# 测试超时时间（秒）
STARTUP_TIMEOUT=30
HEALTH_CHECK_TIMEOUT=60

# 显示配置
show_config() {
    log_info "===== 运行配置 ====="
    echo "  镜像: ${IMAGE_NAME}:${IMAGE_TAG}"
    echo "  容器名: ${CONTAINER_NAME}"
    echo "  端口映射: ${HOST_PORT}:${CONTAINER_PORT}"
    echo "  启动超时: ${STARTUP_TIMEOUT}s"
    echo "  健康检查超时: ${HEALTH_CHECK_TIMEOUT}s"
    echo ""
}

# 检查镜像是否存在
check_image() {
    log_info "检查镜像是否存在..."
    
    if ! docker images "${IMAGE_NAME}:${IMAGE_TAG}" | grep -q "${IMAGE_TAG}"; then
        log_error "镜像 ${IMAGE_NAME}:${IMAGE_TAG} 不存在"
        log_info "请先运行构建脚本: ./deploy/scripts/build-local.sh"
        exit 1
    fi
    
    log_success "镜像检查通过"
}

# 清理已存在的容器
cleanup_container() {
    log_info "清理已存在的容器..."
    
    if docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
        log_warning "发现已存在的容器: ${CONTAINER_NAME}"
        docker stop "${CONTAINER_NAME}" >/dev/null 2>&1 || true
        docker rm "${CONTAINER_NAME}" >/dev/null 2>&1 || true
        log_success "已清理旧容器"
    else
        log_info "没有需要清理的容器"
    fi
}

# 启动容器
start_container() {
    log_info "启动容器..."
    
    # 创建临时配置文件（如果需要）
    local temp_config="${PROJECT_ROOT}/backend/config/test.yml"
    
    docker run -d \
        --name "${CONTAINER_NAME}" \
        -p "${HOST_PORT}:${CONTAINER_PORT}" \
        -e APP_ENV=test \
        "${IMAGE_NAME}:${IMAGE_TAG}"
    
    log_success "容器已启动: ${CONTAINER_NAME}"
}

# 等待容器就绪
wait_for_container() {
    log_info "等待容器就绪..."
    
    local count=0
    local max_attempts=$((STARTUP_TIMEOUT / 2))
    
    while [[ $count -lt $max_attempts ]]; do
        if docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
            log_success "容器运行中"
            return 0
        fi
        
        count=$((count + 1))
        sleep 2
    done
    
    log_error "容器启动超时"
    show_container_logs
    exit 1
}

# 健康检查
health_check() {
    log_info "执行健康检查..."
    
    local count=0
    local max_attempts=$((HEALTH_CHECK_TIMEOUT / 5))
    local url="http://localhost:${HOST_PORT}"
    
    while [[ $count -lt $max_attempts ]]; do
        log_info "尝试访问: ${url} (${count}/${max_attempts})"
        
        # 检查根路径
        if curl -sf "${url}/" >/dev/null 2>&1; then
            log_success "服务响应正常"
            return 0
        fi
        
        # 检查 Swagger 文档
        if curl -sf "${url}/swagger/index.html" >/dev/null 2>&1; then
            log_success "Swagger 文档可访问"
            return 0
        fi
        
        count=$((count + 1))
        sleep 5
    done
    
    log_error "健康检查超时"
    return 1
}

# API 测试
test_apis() {
    log_info "测试 API 端点..."
    
    local base_url="http://localhost:${HOST_PORT}"
    local test_passed=0
    local test_failed=0
    
    # 测试根路径
    echo -n "  测试 GET / ... "
    if curl -sf "${base_url}/" >/dev/null 2>&1; then
        echo -e "${GREEN}✓${NC}"
        test_passed=$((test_passed + 1))
    else
        echo -e "${RED}✗${NC}"
        test_failed=$((test_failed + 1))
    fi
    
    # 测试 Swagger 文档
    echo -n "  测试 GET /swagger/index.html ... "
    if curl -sf "${base_url}/swagger/index.html" >/dev/null 2>&1; then
        echo -e "${GREEN}✓${NC}"
        test_passed=$((test_passed + 1))
    else
        echo -e "${RED}✗${NC}"
        test_failed=$((test_failed + 1))
    fi
    
    # 测试 Swagger JSON
    echo -n "  测试 GET /swagger/doc.json ... "
    if curl -sf "${base_url}/swagger/doc.json" >/dev/null 2>&1; then
        echo -e "${GREEN}✓${NC}"
        test_passed=$((test_passed + 1))
    else
        echo -e "${RED}✗${NC}"
        test_failed=$((test_failed + 1))
    fi
    
    echo ""
    log_info "API 测试结果: ${test_passed} 通过, ${test_failed} 失败"
    
    return $test_failed
}

# 显示容器日志
show_container_logs() {
    log_info "容器日志（最后50行）:"
    echo "----------------------------------------"
    docker logs --tail 50 "${CONTAINER_NAME}" 2>&1 || true
    echo "----------------------------------------"
}

# 显示容器信息
show_container_info() {
    log_info "容器信息:"
    docker inspect "${CONTAINER_NAME}" --format='
容器ID: {{.Id}}
状态: {{.State.Status}}
运行时间: {{.State.StartedAt}}
重启次数: {{.RestartCount}}
端口映射: {{range $p, $conf := .NetworkSettings.Ports}}{{$p}} -> {{(index $conf 0).HostPort}} {{end}}
'
}

# 显示资源使用情况
show_resource_usage() {
    log_info "资源使用情况:"
    docker stats "${CONTAINER_NAME}" --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}\t{{.BlockIO}}"
}

# 停止并清理容器
cleanup() {
    if [[ "${KEEP_RUNNING:-}" != "true" ]]; then
        log_info "停止并清理容器..."
        docker stop "${CONTAINER_NAME}" >/dev/null 2>&1 || true
        docker rm "${CONTAINER_NAME}" >/dev/null 2>&1 || true
        log_success "容器已清理"
    else
        log_info "容器保持运行状态"
        log_info "访问应用: http://localhost:${HOST_PORT}"
        log_info "访问 Swagger: http://localhost:${HOST_PORT}/swagger/index.html"
        log_info "停止容器: docker stop ${CONTAINER_NAME}"
        log_info "删除容器: docker rm ${CONTAINER_NAME}"
    fi
}

# 主函数
main() {
    echo ""
    log_info "===== Docker 容器运行测试 ====="
    echo ""
    
    show_config
    check_image
    cleanup_container
    
    # 启动容器
    start_container
    wait_for_container
    
    # 显示初始日志
    echo ""
    show_container_logs
    echo ""
    
    # 健康检查
    sleep 3  # 给服务一些启动时间
    
    if health_check; then
        log_success "健康检查通过"
        
        # 执行 API 测试
        echo ""
        if test_apis; then
            log_success "API 测试通过"
        else
            log_warning "部分 API 测试失败"
        fi
        
        # 显示容器信息
        echo ""
        show_container_info
        show_resource_usage
        
        echo ""
        log_success "===== 运行测试完成 ====="
        echo ""
        
        cleanup
        exit 0
    else
        log_error "健康检查失败"
        echo ""
        show_container_logs
        echo ""
        cleanup
        exit 1
    fi
}

# 捕获退出信号，确保清理
trap cleanup EXIT INT TERM

# 执行主函数
main "$@"
