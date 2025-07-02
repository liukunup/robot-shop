# 第一阶段：构建前端
ARG NODE_VERSION=latest
FROM node:${NODE_VERSION} AS frontend-builder

WORKDIR /app
# 注意路径调整为从项目根目录开始的相对路径
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci

COPY frontend .
RUN npm run build

# 第二阶段：构建后端
ARG GO_VERSION=latest
FROM golang:${GO_VERSION} AS backend-builder

WORKDIR /app
# 注意路径调整为从项目根目录开始的相对路径
COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend .
RUN CGO_ENABLED=0 GOOS=linux go build -o ../deploy/server

# 最终镜像
FROM alpine:latest

WORKDIR /app

# 安装时区数据和必要的依赖（可选）
RUN apk add --no-cache tzdata ca-certificates

# 从构建阶段复制文件
COPY --from=frontend-builder /app/dist ./static
COPY --from=backend-builder /app/deploy/server .

# 应用配置文件和静态资源（如果有）
# COPY deploy/config.yaml ./

EXPOSE 8080
CMD ["./server"]
