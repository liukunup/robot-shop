#!/usr/bin/env bash

###############################################################################
# 生产环境部署脚本
# 用途：一键部署 Robot Shop 到生产环境
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
COMPOSE_FILE="${PROJECT_ROOT}/deploy/docker-compose/docker-compose.yml"
ENV_FILE="${PROJECT_ROOT}/.env"

# 显示横幅
show_banner() {
    echo ""
    echo -e "${BLUE}╔════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║                                            ║${NC}"
    echo -e "${BLUE}║     Robot Shop - 生产环境部署              ║${NC}"
    echo -e "${BLUE}║                                            ║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════════╝${NC}"
    echo ""
}

# 检查环境文件
check_env_file() {
    log_info "检查环境配置文件..."
    
    if [[ ! -f "${ENV_FILE}" ]]; then
        log_warning ".env 文件不存在"
        log_info "从 .env.example 创建 .env 文件..."
        
        if [[ -f "${PROJECT_ROOT}/.env.example" ]]; then
            cp "${PROJECT_ROOT}/.env.example" "${ENV_FILE}"
            log_warning "请编辑 .env 文件并配置生产环境参数"
            log_info "编辑文件: ${ENV_FILE}"
            
            read -p "是否现在编辑 .env 文件? (y/N): " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                ${EDITOR:-vim} "${ENV_FILE}"
            else
                log_error "请先配置 .env 文件后再运行此脚本"
                exit 1
            fi
        else
            log_error ".env.example 文件不存在"
            exit 1
        fi
    fi
    
    log_success "环境配置文件检查通过"
}

# 验证关键配置
validate_config() {
    log_info "验证关键配置..."
    
    # 加载环境变量
    set -a
    source "${ENV_FILE}"
    set +a
    
    local has_error=false
    
    # 检查 JWT 密钥
    if [[ "${JWT_SECRET_KEY:-}" == "CHANGE_THIS_TO_RANDOM_SECRET_KEY_IN_PRODUCTION" ]] || \
       [[ "${JWT_SECRET_KEY:-}" == "dev_secret_key_not_for_production" ]]; then
        log_error "JWT_SECRET_KEY 未设置或使用默认值，请设置强随机密钥"
        has_error=true
    fi
    
    # 检查数据库密码
    if [[ "${MYSQL_ROOT_PASSWORD:-}" == "root" ]] || \
       [[ "${MYSQL_PASSWORD:-}" == "robotshop" ]] || \
       [[ "${MYSQL_PASSWORD:-}" == "ChangeMe123!" ]]; then
        log_warning "MySQL 密码使用默认值，建议修改为强密码"
    fi
    
    # 检查应用环境
    if [[ "${APP_ENV:-}" != "prod" ]] && [[ "${APP_ENV:-}" != "production" ]]; then
        log_warning "APP_ENV 未设置为 prod，当前值: ${APP_ENV:-未设置}"
    fi
    
    if [[ "${has_error}" == "true" ]]; then
        log_error "配置验证失败，请修复上述问题"
        exit 1
    fi
    
    log_success "配置验证通过"
}

# 检查 Docker
check_docker() {
    log_info "检查 Docker 环境..."
    
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装"
        exit 1
    fi
    
    if ! command -v docker compose &> /dev/null; then
        log_error "Docker Compose 未安装"
        exit 1
    fi
    
    if ! docker info &> /dev/null; then
        log_error "Docker 未运行"
        exit 1
    fi
    
    log_success "Docker 环境检查通过"
}

# 拉取镜像
pull_images() {
    log_info "拉取最新镜像..."
    
    cd "${PROJECT_ROOT}/deploy/docker-compose"
    
    if docker compose --env-file "${ENV_FILE}" pull; then
        log_success "镜像拉取成功"
    else
        log_error "镜像拉取失败"
        exit 1
    fi
}

# 备份数据库（如果存在）
backup_database() {
    log_info "检查是否需要备份..."
    
    # 检查是否有运行中的 MySQL 容器
    local mysql_container="${COMPOSE_PROJECT_NAME:-robotshop}-mysql"
    
    if docker ps --format '{{.Names}}' | grep -q "^${mysql_container}$"; then
        log_warning "发现运行中的数据库，建议先备份"
        
        read -p "是否创建数据库备份? (Y/n): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Nn]$ ]]; then
            local backup_dir="${PROJECT_ROOT}/backups"
            local backup_file="${backup_dir}/backup_$(date +%Y%m%d_%H%M%S).sql"
            
            mkdir -p "${backup_dir}"
            
            log_info "创建数据库备份: ${backup_file}"
            
            if docker exec "${mysql_container}" mysqldump \
                -u root \
                -p"${MYSQL_ROOT_PASSWORD:-ChangeMe123!}" \
                --all-databases \
                --single-transaction \
                --quick \
                --lock-tables=false \
                > "${backup_file}"; then
                
                log_success "数据库备份成功: ${backup_file}"
                
                # 压缩备份
                gzip "${backup_file}"
                log_success "备份已压缩: ${backup_file}.gz"
            else
                log_error "数据库备份失败"
                exit 1
            fi
        fi
    else
        log_info "没有需要备份的数据"
    fi
}

# 停止旧容器
stop_old_containers() {
    log_info "停止旧容器..."
    
    cd "${PROJECT_ROOT}/deploy/docker-compose"
    
    if docker compose --env-file "${ENV_FILE}" ps -q | grep -q .; then
        docker compose --env-file "${ENV_FILE}" down
        log_success "旧容器已停止"
    else
        log_info "没有运行中的容器"
    fi
}

# 启动服务
start_services() {
    log_info "启动服务..."
    
    cd "${PROJECT_ROOT}/deploy/docker-compose"
    
    if docker compose --env-file "${ENV_FILE}" up -d; then
        log_success "服务启动成功"
    else
        log_error "服务启动失败"
        exit 1
    fi
}

# 等待服务就绪
wait_for_services() {
    log_info "等待服务就绪..."
    
    local max_attempts=60
    local attempt=0
    
    while [[ $attempt -lt $max_attempts ]]; do
        if docker compose --env-file "${ENV_FILE}" ps | grep -q "healthy"; then
            log_success "服务已就绪"
            return 0
        fi
        
        attempt=$((attempt + 1))
        echo -n "."
        sleep 2
    done
    
    echo ""
    log_error "服务启动超时"
    return 1
}

# 显示服务状态
show_status() {
    log_info "服务状态:"
    echo ""
    
    cd "${PROJECT_ROOT}/deploy/docker-compose"
    docker compose --env-file "${ENV_FILE}" ps
    
    echo ""
    log_info "服务日志 (最后 20 行):"
    docker compose --env-file "${ENV_FILE}" logs --tail=20
}

# 运行健康检查
health_check() {
    log_info "运行健康检查..."
    
    local app_url="http://localhost:${HTTP_PORT:-8000}"
    
    # 检查应用
    if curl -sf "${app_url}/" > /dev/null; then
        log_success "应用服务健康检查通过"
    else
        log_error "应用服务健康检查失败"
        return 1
    fi
    
    # 检查 Swagger
    if curl -sf "${app_url}/swagger/index.html" > /dev/null; then
        log_success "Swagger 文档可访问"
    else
        log_warning "Swagger 文档不可访问"
    fi
    
    return 0
}

# 显示访问信息
show_access_info() {
    echo ""
    log_success "===== 部署完成 ====="
    echo ""
    log_info "访问信息:"
    echo "  应用地址: http://localhost:${HTTP_PORT:-8000}"
    echo "  Swagger 文档: http://localhost:${HTTP_PORT:-8000}/swagger/index.html"
    echo "  MySQL: localhost:${MYSQL_HOST_PORT:-3306}"
    echo "  Redis: localhost:${REDIS_HOST_PORT:-6379}"
    echo ""
    log_info "管理命令:"
    echo "  查看日志: cd deploy/docker-compose && docker compose --env-file ../../.env logs -f"
    echo "  停止服务: cd deploy/docker-compose && docker compose --env-file ../../.env down"
    echo "  重启服务: cd deploy/docker-compose && docker compose --env-file ../../.env restart"
    echo "  查看状态: cd deploy/docker-compose && docker compose --env-file ../../.env ps"
    echo ""
}

# 主函数
main() {
    show_banner
    
    check_env_file
    validate_config
    check_docker
    
    # 备份（如果需要）
    backup_database || true
    
    # 拉取镜像
    pull_images
    
    # 停止旧服务
    stop_old_containers
    
    # 启动新服务
    start_services
    
    # 等待服务就绪
    if wait_for_services; then
        show_status
        
        # 健康检查
        sleep 5
        if health_check; then
            show_access_info
            exit 0
        else
            log_error "健康检查失败，请查看日志"
            exit 1
        fi
    else
        log_error "服务启动失败"
        show_status
        exit 1
    fi
}

# 执行主函数
main "$@"
