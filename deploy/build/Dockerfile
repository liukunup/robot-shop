# 默认镜像仓库
ARG REGISTRY=docker.io

# 是否使用国内加速
ARG USE_CHINA_MIRROR=true

# 第一阶段：构建前端
ARG NODE_VERSION=22.17.0
FROM ${REGISTRY}/node:${NODE_VERSION}-alpine AS frontend-builder

# 国内加速
RUN if [ "$USE_CHINA_MIRROR" = "true" ]; then \
      # 设置 Alpine 镜像源
      set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
      # 设置 NPM 镜像
      npm config set registry https://registry.npmmirror.com; \
    fi

WORKDIR /app
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci

COPY frontend .
RUN npm run build

# 第二阶段：构建后端
ARG GO_VERSION=1.24.4
FROM ${REGISTRY}/golang:${GO_VERSION}-alpine AS backend-builder

# 国内加速
RUN if [ "$USE_CHINA_MIRROR" = "true" ]; then \
      # 设置 Alpine 镜像源
      set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
      # 设置 Go 代理
      go env -w GOPROXY=https://goproxy.cn,direct; \
    fi

WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ../deploy/server

# 最终镜像
FROM alpine:3.22

WORKDIR /app

# 安装时区数据和必要的依赖
RUN apk add --no-cache tzdata ca-certificates

# 从构建阶段复制文件
COPY --from=frontend-builder /app/dist ./static
COPY --from=backend-builder /app/deploy/server .

# 应用配置文件和静态资源（如果有）
# COPY deploy/config.yaml ./

EXPOSE 8080
CMD ["./server"]
