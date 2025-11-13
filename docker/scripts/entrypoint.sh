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

echo "${YELLOW}[1/2] 执行数据库迁移...${NC}"
# 运行数据库迁移
if [ -f "./migration" ]; then
    ./migration -conf "$CONFIG_FILE"
    echo "${GREEN}✓ 数据库迁移完成${NC}"
else
    echo "${YELLOW}⚠ migration 可执行文件不存在，跳过迁移步骤${NC}"
fi

echo "${YELLOW}[2/2] 启动服务器...${NC}"
# 启动服务器
exec ./server -conf "$CONFIG_FILE"
