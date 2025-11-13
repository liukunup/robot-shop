#!/usr/bin/env bash

#############################################################################
# Robot Shop 数据库备份脚本
# 用途: 自动备份 MySQL 数据库
# 
# 使用方法:
#   ./backup.sh                  # 执行完整备份
#   ./backup.sh --db robotshop   # 备份指定数据库
#   ./backup.sh --clean          # 清理旧备份
#
# Crontab 示例 (每天凌晨 2 点备份):
#   0 2 * * * /path/to/backup.sh >> /var/log/robot-shop-backup.log 2>&1
#############################################################################

set -e

# 颜色定义
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m' # No Color

# 配置
readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
readonly COMPOSE_DIR="${PROJECT_ROOT}/deploy/docker-compose"
readonly ENV_FILE="${PROJECT_ROOT}/.env"
readonly BACKUP_DIR="${PROJECT_ROOT}/backups"
readonly RETENTION_DAYS=7  # 备份保留天数

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $(date '+%Y-%m-%d %H:%M:%S') $*"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $(date '+%Y-%m-%d %H:%M:%S') $*"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $(date '+%Y-%m-%d %H:%M:%S') $*"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $(date '+%Y-%m-%d %H:%M:%S') $*" >&2
}

# 显示横幅
show_banner() {
    cat << "EOF"
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║         🤖 Robot Shop 数据库备份工具 v1.0.0              ║
║                                                           ║
║         自动备份 MySQL 数据库并管理备份文件               ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
EOF
}

# 检查环境
check_environment() {
    log_info "检查环境..."
    
    # 检查 .env 文件
    if [[ ! -f "${ENV_FILE}" ]]; then
        log_error ".env 文件不存在: ${ENV_FILE}"
        log_info "请先复制 .env.example 并配置环境变量"
        exit 1
    fi
    
    # 加载环境变量
    set -a
    source "${ENV_FILE}"
    set +a
    
    # 检查 Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装或未在 PATH 中"
        exit 1
    fi
    
    # 检查 Docker Compose
    if ! docker compose version &> /dev/null; then
        log_error "Docker Compose 未安装或版本过低"
        exit 1
    fi
    
    # 检查 MySQL 容器是否运行
    if ! docker compose -f "${COMPOSE_DIR}/docker-compose.yml" --env-file "${ENV_FILE}" ps mysql | grep -q "Up"; then
        log_error "MySQL 容器未运行"
        log_info "请先启动 MySQL 容器: docker compose up -d mysql"
        exit 1
    fi
    
    log_success "环境检查通过"
}

# 创建备份目录
create_backup_dir() {
    if [[ ! -d "${BACKUP_DIR}" ]]; then
        log_info "创建备份目录: ${BACKUP_DIR}"
        mkdir -p "${BACKUP_DIR}"
    fi
}

# 备份数据库
backup_database() {
    local db_name="$1"
    local timestamp=$(date '+%Y%m%d_%H%M%S')
    local backup_file="${BACKUP_DIR}/backup_${db_name}_${timestamp}.sql"
    
    log_info "开始备份数据库: ${db_name}"
    
    # 执行备份
    if [[ "${db_name}" == "all" ]]; then
        log_info "备份所有数据库..."
        docker compose -f "${COMPOSE_DIR}/docker-compose.yml" \
            --env-file "${ENV_FILE}" \
            exec -T mysql mysqldump \
            -u root \
            -p"${MYSQL_ROOT_PASSWORD}" \
            --all-databases \
            --single-transaction \
            --quick \
            --lock-tables=false \
            --routines \
            --triggers \
            --events \
            > "${backup_file}" 2>/dev/null || {
                log_error "数据库备份失败"
                rm -f "${backup_file}"
                return 1
            }
    else
        log_info "备份数据库: ${db_name}"
        docker compose -f "${COMPOSE_DIR}/docker-compose.yml" \
            --env-file "${ENV_FILE}" \
            exec -T mysql mysqldump \
            -u root \
            -p"${MYSQL_ROOT_PASSWORD}" \
            "${db_name}" \
            --single-transaction \
            --quick \
            --lock-tables=false \
            --routines \
            --triggers \
            > "${backup_file}" 2>/dev/null || {
                log_error "数据库备份失败"
                rm -f "${backup_file}"
                return 1
            }
    fi
    
    # 压缩备份
    log_info "压缩备份文件..."
    gzip "${backup_file}"
    backup_file="${backup_file}.gz"
    
    # 验证备份文件
    if [[ -f "${backup_file}" ]]; then
        local file_size=$(du -h "${backup_file}" | cut -f1)
        log_success "备份完成: ${backup_file} (${file_size})"
        return 0
    else
        log_error "备份文件不存在"
        return 1
    fi
}

# 清理旧备份
clean_old_backups() {
    log_info "清理 ${RETENTION_DAYS} 天前的备份..."
    
    local deleted_count=0
    while IFS= read -r -d '' file; do
        rm -f "${file}"
        ((deleted_count++))
        log_info "已删除: $(basename "${file}")"
    done < <(find "${BACKUP_DIR}" -name "backup_*.sql.gz" -type f -mtime +${RETENTION_DAYS} -print0 2>/dev/null)
    
    if [[ ${deleted_count} -gt 0 ]]; then
        log_success "已清理 ${deleted_count} 个旧备份文件"
    else
        log_info "没有需要清理的旧备份"
    fi
}

# 列出备份文件
list_backups() {
    log_info "备份文件列表:"
    echo ""
    
    if [[ ! -d "${BACKUP_DIR}" ]] || [[ -z "$(ls -A "${BACKUP_DIR}"/*.sql.gz 2>/dev/null)" ]]; then
        log_warning "没有找到备份文件"
        return 0
    fi
    
    printf "%-40s %10s %20s\n" "文件名" "大小" "创建时间"
    printf "%-40s %10s %20s\n" "----------------------------------------" "----------" "--------------------"
    
    while IFS= read -r file; do
        if [[ -f "${file}" ]]; then
            local filename=$(basename "${file}")
            local size=$(du -h "${file}" | cut -f1)
            local mtime=$(stat -f "%Sm" -t "%Y-%m-%d %H:%M:%S" "${file}" 2>/dev/null || stat -c "%y" "${file}" 2>/dev/null | cut -d'.' -f1)
            printf "%-40s %10s %20s\n" "${filename}" "${size}" "${mtime}"
        fi
    done < <(find "${BACKUP_DIR}" -name "backup_*.sql.gz" -type f | sort -r)
    
    echo ""
}

# 恢复数据库
restore_database() {
    local backup_file="$1"
    
    if [[ ! -f "${backup_file}" ]]; then
        log_error "备份文件不存在: ${backup_file}"
        return 1
    fi
    
    log_warning "⚠️  即将恢复数据库,这将覆盖现有数据!"
    read -p "确认要继续吗? (yes/no): " confirm
    
    if [[ "${confirm}" != "yes" ]]; then
        log_info "操作已取消"
        return 0
    fi
    
    log_info "开始恢复数据库..."
    
    # 解压并恢复
    if [[ "${backup_file}" == *.gz ]]; then
        gunzip < "${backup_file}" | \
            docker compose -f "${COMPOSE_DIR}/docker-compose.yml" \
                --env-file "${ENV_FILE}" \
                exec -T mysql mysql \
                -u root \
                -p"${MYSQL_ROOT_PASSWORD}" 2>/dev/null || {
                    log_error "数据库恢复失败"
                    return 1
                }
    else
        docker compose -f "${COMPOSE_DIR}/docker-compose.yml" \
            --env-file "${ENV_FILE}" \
            exec -T mysql mysql \
            -u root \
            -p"${MYSQL_ROOT_PASSWORD}" < "${backup_file}" 2>/dev/null || {
                log_error "数据库恢复失败"
                return 1
            }
    fi
    
    log_success "数据库恢复完成"
}

# 显示帮助
show_help() {
    cat << EOF
Robot Shop 数据库备份工具

用法:
    $0 [选项]

选项:
    --db <name>         备份指定数据库 (默认: all)
    --clean             清理旧备份 (保留 ${RETENTION_DAYS} 天)
    --list              列出所有备份文件
    --restore <file>    从备份文件恢复数据库
    --retention <days>  设置备份保留天数 (默认: ${RETENTION_DAYS})
    -h, --help          显示帮助信息

示例:
    # 备份所有数据库
    $0
    
    # 备份指定数据库
    $0 --db robotshop
    
    # 清理旧备份
    $0 --clean
    
    # 列出备份文件
    $0 --list
    
    # 恢复数据库
    $0 --restore backups/backup_all_20251112_120000.sql.gz
    
    # 设置备份保留 30 天
    $0 --retention 30

Crontab 示例:
    # 每天凌晨 2 点备份
    0 2 * * * /path/to/backup.sh >> /var/log/robot-shop-backup.log 2>&1
    
    # 每周日凌晨 3 点清理旧备份
    0 3 * * 0 /path/to/backup.sh --clean >> /var/log/robot-shop-backup.log 2>&1

EOF
}

# 主函数
main() {
    local db_name="all"
    local do_clean=false
    local do_list=false
    local restore_file=""
    local custom_retention=false
    
    # 解析参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            --db)
                db_name="$2"
                shift 2
                ;;
            --clean)
                do_clean=true
                shift
                ;;
            --list)
                do_list=true
                shift
                ;;
            --restore)
                restore_file="$2"
                shift 2
                ;;
            --retention)
                RETENTION_DAYS="$2"
                custom_retention=true
                shift 2
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                log_error "未知参数: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    show_banner
    echo ""
    
    # 列出备份
    if [[ "${do_list}" == true ]]; then
        check_environment
        create_backup_dir
        list_backups
        exit 0
    fi
    
    # 恢复数据库
    if [[ -n "${restore_file}" ]]; then
        check_environment
        restore_database "${restore_file}"
        exit $?
    fi
    
    # 执行备份
    check_environment
    create_backup_dir
    
    if ! backup_database "${db_name}"; then
        log_error "备份失败"
        exit 1
    fi
    
    # 清理旧备份
    if [[ "${do_clean}" == true ]]; then
        clean_old_backups
    fi
    
    # 显示备份列表
    list_backups
    
    log_success "所有操作完成"
}

# 运行主函数
main "$@"
