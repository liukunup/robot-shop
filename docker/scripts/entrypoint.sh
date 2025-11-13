#!/bin/sh
set -e

# 颜色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "${GREEN}========================================${NC}"
echo "${GREEN}Robot Shop 容器启动中...${NC}"
echo "${GREEN}========================================${NC}"

# 配置文件路径
CONFIG_FILE="${CONFIG_FILE:-./config/prod.yml}"

# 检查配置文件是否存在
if [ ! -f "$CONFIG_FILE" ]; then
    echo "${RED}错误: 配置文件 $CONFIG_FILE 不存在${NC}"
    exit 1
fi

echo "${YELLOW}[1/3] 执行数据库迁移...${NC}"
if [ -x "./migration" ]; then
    if ./migration -conf "$CONFIG_FILE"; then
        echo "${GREEN}✓ 数据库迁移完成${NC}"
    else
        echo "${YELLOW}⚠ 数据库迁移失败，继续启动（请检查日志）${NC}"
    fi
else
    echo "${YELLOW}⚠ migration 可执行文件不存在，跳过迁移步骤${NC}"
fi

echo "${YELLOW}[2/3] 启动后端服务...${NC}"
if [ -x "./server" ]; then
    ./server -conf "$CONFIG_FILE" &
    BACK_PID=$!
    echo "${GREEN}✓ 后端进程 PID: $BACK_PID${NC}"
else
    echo "${RED}错误: server 可执行文件不存在${NC}"
    exit 1
fi

ATTEMPT=0
MAX_ATTEMPTS=20
until wget -qO- http://127.0.0.1:8000/healthz >/dev/null 2>&1; do
    ATTEMPT=$((ATTEMPT + 1))
    if [ "$ATTEMPT" -ge "$MAX_ATTEMPTS" ]; then
        echo "${YELLOW}⚠ 后端健康检查超时，继续启动 nginx${NC}"
        break
    fi
    echo "${YELLOW}等待后端就绪($ATTEMPT/$MAX_ATTEMPTS)...${NC}"
    sleep 1
done

echo "${YELLOW}[3/3] 启动 nginx...${NC}"
exec nginx -g 'daemon off;'
