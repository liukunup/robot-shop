#!/bin/bash

# 设置环境变量
COMPOSE_PROJECT_NAME=${COMPOSE_PROJECT_NAME:-robotshop}
REGISTRY=${REGISTRY:-docker.io}
MYSQL_VERSION=${MYSQL_VERSION:-8}
MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD:-ChangeMe123!}
MYSQL_DATABASE=${MYSQL_DATABASE:-robotshop}
MYSQL_USER=${MYSQL_USER:-robotshop}
MYSQL_PASSWORD=${MYSQL_PASSWORD:-ChangeMe123!}
MYSQL_PORT=${MYSQL_PORT:-3306}
REDIS_VERSION=${REDIS_VERSION:-7}
REDIS_PASSWORD=${REDIS_PASSWORD:-ChangeMe123!}
REDIS_PORT=${REDIS_PORT:-6379}
TZ=${TZ:-Asia/Shanghai}

# 创建网络
docker network create \
  --driver bridge \
  --attachable \
  ${COMPOSE_PROJECT_NAME}-network

# 创建MySQL数据卷
docker volume create \
  --driver local \
  --name ${COMPOSE_PROJECT_NAME}-mysql-data

# 创建Redis数据卷
docker volume create \
  --driver local \
  --name ${COMPOSE_PROJECT_NAME}-redis-data

# 启动MySQL容器
docker run -d \
  --name ${COMPOSE_PROJECT_NAME}-mysql \
  --network ${COMPOSE_PROJECT_NAME}-network \
  -e MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD} \
  -e MYSQL_DATABASE=${MYSQL_DATABASE} \
  -e MYSQL_USER=${MYSQL_USER} \
  -e MYSQL_PASSWORD=${MYSQL_PASSWORD} \
  -e TZ=${TZ} \
  -p ${MYSQL_PORT}:3306 \
  -v ${COMPOSE_PROJECT_NAME}-mysql-data:/var/lib/mysql \
  -v $(pwd)/mysql/conf.d:/etc/mysql/conf.d \
  -v $(pwd)/mysql/initdb:/docker-entrypoint-initdb.d:ro \
  --health-cmd="mysqladmin ping -p${MYSQL_ROOT_PASSWORD}" \
  --health-interval=10s \
  --health-timeout=5s \
  --health-retries=3 \
  --restart unless-stopped \
  ${REGISTRY}/mysql:${MYSQL_VERSION}

# 启动Redis容器
docker run -d \
  --name ${COMPOSE_PROJECT_NAME}-redis \
  --network ${COMPOSE_PROJECT_NAME}-network \
  -e TZ=${TZ} \
  -p ${REDIS_PORT}:6379 \
  -v ${COMPOSE_PROJECT_NAME}-redis-data:/data \
  -v $(pwd)/redis/redis.conf:/usr/local/etc/redis/redis.conf:ro \
  --health-cmd="redis-cli -a ${REDIS_PASSWORD} ping" \
  --health-interval=10s \
  --health-timeout=3s \
  --health-retries=3 \
  --restart unless-stopped \
  ${REGISTRY}/redis:${REDIS_VERSION} \
  redis-server /usr/local/etc/redis/redis.conf --requirepass ${REDIS_PASSWORD}

echo "部署完成!"
echo "MySQL 运行在端口: ${MYSQL_PORT}"
echo "Redis 运行在端口: ${REDIS_PORT}"