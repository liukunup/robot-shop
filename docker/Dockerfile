# 默认镜像仓库
ARG REGISTRY=docker.io

# 国内镜像源加速
ARG USE_CHINA_MIRROR=true

# 构建镜像的版本
ARG NODE_VERSION=22-alpine
ARG GO_VERSION=1.23-alpine
ARG NGINX_VERSION=1.29-alpine


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
    mkdir -p output/bin output/config && \
    go mod download && \
    go build -ldflags="-s -w" -o ./output/bin/server ./cmd/server && \
    go build -ldflags="-s -w" -o ./output/bin/migration ./cmd/migration && \
    go build -ldflags="-s -w" -o ./output/bin/task ./cmd/task && \
    cp ./config/prod.yml ./output/config/prod.yml


# 第三阶段:构建运行镜像
FROM ${REGISTRY}/nginx:${NGINX_VERSION}

RUN apk add --no-cache ca-certificates openssl && \
    mkdir -p /etc/nginx/ssl && \
    openssl req -x509 -nodes -newkey rsa:2048 \
        -keyout /etc/nginx/ssl/server.key \
        -out /etc/nginx/ssl/server.crt \
        -days 3650 \
        -subj "/C=CN/ST=Beijing/L=Beijing/O=/OU=/CN=RobotShop" \
        -addext "subjectAltName=DNS:localhost" && \
    openssl dhparam -out /etc/nginx/dhparam.pem 2048 && \
    cp /etc/nginx/ssl/server.crt /etc/nginx/ssl/fe.local.crt && \
    cp /etc/nginx/ssl/server.key /etc/nginx/ssl/fe.local.key && \
    apk del openssl

WORKDIR /app
RUN mkdir -p storage/logs config /var/www/fe.local/public
COPY --from=frontend-builder /app/frontend/dist /var/www/fe.local/public
COPY --from=backend-builder /app/backend/output/ ./
COPY docker/nginxconfig/ /etc/nginx/

RUN rm -f /etc/nginx/conf.d/default.conf

COPY docker/scripts/entrypoint.sh ./entrypoint.sh

RUN chmod +x ./entrypoint.sh

EXPOSE 80 443

ENTRYPOINT [ "./entrypoint.sh" ]
