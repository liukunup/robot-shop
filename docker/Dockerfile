# 默认镜像仓库
ARG REGISTRY=docker.io

# 国内镜像源加速
ARG USE_CHINA_MIRROR=true

# 构建镜像的版本
ARG NODE_VERSION=22-alpine
ARG GO_VERSION=1.23-alpine
ARG NGINX_VERSION=alpine


# 第一阶段：构建前端
FROM ${REGISTRY}/node:${NODE_VERSION} AS frontend-builder

# 国内加速
RUN if [ "$USE_CHINA_MIRROR" = "true" ]; then \
      # 设置 Alpine 镜像源
      sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
      # 设置 NPM 镜像
      npm config set registry https://registry.npmmirror.com; \
    fi

# 构建代码
WORKDIR /app/frontend
COPY frontend/ ./
RUN npm install && \
    npm run build


# 第二阶段：构建后端
FROM ${REGISTRY}/golang:${GO_VERSION} AS backend-builder

# 国内加速
ARG USE_CHINA_MIRROR
RUN if [ "$USE_CHINA_MIRROR" = "true" ]; then \
      # 设置 Alpine 镜像源
      sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
      # 设置 Go 代理
      go env -w GOPROXY=https://goproxy.cn,direct; \
    fi

# 构建代码
WORKDIR /app/backend
COPY backend/ ./
RUN export CGO_ENABLED=0 && \
    export GOOS=linux && \
    export GOARCH=${TARGETARCH:-amd64} && \
    mkdir -p bin && \
    go mod download && \
    go build -ldflags="-s -w" -o ./bin/server ./cmd/server && \
    go build -ldflags="-s -w" -o ./bin/migration ./cmd/migration && \
    go build -ldflags="-s -w" -o ./bin/task ./cmd/task


# 第三阶段：构建运行镜像
FROM ${REGISTRY}/nginx:${NGINX_VERSION}

# 安装时区数据和必要的依赖（alpine 版本）
RUN apk add --no-cache tzdata ca-certificates && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata

WORKDIR /app
RUN mkdir -p storage/logs config /run/nginx /usr/share/nginx/html
COPY --from=backend-builder /app/backend/bin/server ./server
COPY --from=backend-builder /app/backend/bin/migration ./migration
COPY --from=backend-builder /app/backend/bin/task ./task
COPY --from=backend-builder /app/backend/config/prod.yml ./config/prod.yml
COPY --from=frontend-builder /app/frontend/dist /usr/share/nginx/html
COPY docker/nginxconfig/ /etc/nginx/
COPY docker/scripts/entrypoint.sh ./entrypoint.sh

RUN chmod +x ./entrypoint.sh

EXPOSE 80

ENTRYPOINT [ "./entrypoint.sh" ]
