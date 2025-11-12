# ============================================
# 初始化开发环境
# ============================================
.PHONY: init
init:
	# 安装Wire依赖注入工具 - 用于自动生成依赖注入代码
	go install github.com/google/wire/cmd/wire@latest
	# 安装MockGen工具 - 用于生成测试Mock对象
	go install github.com/golang/mock/mockgen@latest
	# 安装Swag工具 - 用于生成Swagger API文档
	go install github.com/swaggo/swag/cmd/swag@latest
	# 安装Delve调试器 - Go语言调试工具
	go install github.com/go-delve/delve/cmd/dlv@latest

# ============================================
# 启动项目完整流程
# ============================================
.PHONY: bootstrap
bootstrap:
	# 启动Docker Compose中的依赖服务(数据库、Redis等)
	cd ./deploy/docker-compose && docker compose up -d && cd ../../
	# 执行数据库迁移,初始化数据库表结构
	go run ./cmd/migration
	# 使用nunu工具启动HTTP服务器(支持热重载)
	nunu run ./cmd/server

# ============================================
# 生成测试Mock对象
# ============================================
.PHONY: mock
mock:
	# 为Service层生成Mock对象 - 用于Handler层测试
	mockgen -source=internal/service/user.go -destination test/mocks/service/user.go
	mockgen -source=internal/service/robot.go -destination test/mocks/service/robot.go
	mockgen -source=internal/service/role.go -destination test/mocks/service/role.go
	mockgen -source=internal/service/menu.go -destination test/mocks/service/menu.go
	mockgen -source=internal/service/api.go -destination test/mocks/service/api.go
	# 为Repository层生成Mock对象 - 用于Service层测试
	mockgen -source=internal/repository/user.go -destination test/mocks/repository/user.go
	mockgen -source=internal/repository/robot.go -destination test/mocks/repository/robot.go
	mockgen -source=internal/repository/role.go -destination test/mocks/repository/role.go
	mockgen -source=internal/repository/menu.go -destination test/mocks/repository/menu.go
	mockgen -source=internal/repository/api.go -destination test/mocks/repository/api.go
	# 为Transaction等基础Repository生成Mock对象
	mockgen -source=internal/repository/repository.go -destination test/mocks/repository/repository.go

# ============================================
# 运行测试并生成覆盖率报告
# ============================================
.PHONY: test
test:
	# 运行所有测试用例并收集覆盖率数据
	# -coverpkg: 指定要统计覆盖率的包(handler/service/repository)
	# -coverprofile: 输出覆盖率数据到文件
	go test -coverpkg=./internal/handler,./internal/service,./internal/repository -coverprofile=./coverage.out ./test/server/...
	# 将覆盖率数据生成HTML可视化报告
	go tool cover -html=./coverage.out -o coverage.html

# ============================================
# 编译生产环境可执行文件
# ============================================
.PHONY: build
build:
	# 编译Go程序,生成优化后的二进制文件
	# -ldflags="-s -w": 去除调试信息和符号表,减小文件体积
	# -o: 指定输出文件路径
	go build -ldflags="-s -w" -o ./bin/server ./cmd/server

# ============================================
# 构建和测试Docker镜像
# ============================================
.PHONY: docker
docker:
	# 构建Docker镜像
	# -f: 指定Dockerfile路径
	# --build-arg: 传递构建参数(应用相对路径)
	# -t: 指定镜像名称和标签
	docker build -f deploy/build/Dockerfile --build-arg APP_RELATIVE_PATH=./cmd/task -t 1.1.1.1:5000/demo-task:v1 .
	# 运行Docker容器进行测试
	# --rm: 容器退出后自动删除
	# -i: 交互模式
	docker run --rm -i 1.1.1.1:5000/demo-task:v1

# ============================================
# 生成Swagger API文档
# ============================================
.PHONY: swag
swag:
	# 扫描代码注释生成Swagger文档
	# -g: 指定主入口文件
	# -o: 指定文档输出目录
	# --parseDependency: 解析依赖包中的注释
	swag init -g cmd/server/main.go -o docs --parseDependency
