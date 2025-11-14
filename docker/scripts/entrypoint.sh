#!/bin/sh
set -e

echo "========================================"
echo "Robot Shop 容器启动中..."
echo "========================================"

# 配置文件路径
CONFIG_FILE="${CONFIG_FILE:-./config/prod.yml}"

# 检查配置文件是否存在
if [ ! -f "$CONFIG_FILE" ]; then
    echo "错误: 配置文件 $CONFIG_FILE 不存在"
    exit 1
fi

echo "[1/3] 执行数据库迁移..."
if [ -x "./bin/migration" ]; then
    if ./bin/migration -conf "$CONFIG_FILE"; then
        echo "✓ 数据库迁移完成"
    else
        echo "⚠ 数据库迁移失败，继续启动（请检查日志）"
    fi
else
    echo "⚠ migration 可执行文件不存在，跳过迁移步骤"
fi

echo "[2/3] 启动后端服务..."
if [ -x "./bin/server" ]; then
    ./bin/server -conf "$CONFIG_FILE" &
    BACK_PID=$!
    echo "✓ 后端进程 PID: $BACK_PID"
else
    echo "错误: server 可执行文件不存在"
    exit 1
fi

ATTEMPT=0
MAX_ATTEMPTS=60
until wget -qO- http://127.0.0.1:8000/healthz >/dev/null 2>&1; do
    ATTEMPT=$((ATTEMPT + 1))
    if [ "$ATTEMPT" -ge "$MAX_ATTEMPTS" ]; then
        echo "⚠ 后端健康检查超时，继续启动 nginx"
        break
    fi
    echo "等待后端就绪($ATTEMPT/$MAX_ATTEMPTS)..."
    sleep 1
done

echo "[3/3] 启动 nginx..."
exec nginx -g 'daemon off;'
