# Robot Shop - AI 编程助手指导文档

Robot Shop 是一个企业级机器人管理平台，提供多种机器人（钉钉、飞书、企业微信、金山协作）的集成管理功能。

## 前端开发规范

### 技术栈

- Ant Design Pro

### 关键目录

前端代码都在`frontend`目录下，执行任何与前端相关的指令，都应当先`cd frontend`切换到此目录。

- `src/pages/` 业务页面，按功能模块划分
- `src/components/` 可复用的UI组件
- `src/services/` API请求封装，由OpenAPI自动生成
- `src/utils/` 工具函数

### 组件选型优先级

ProComponents > Ant Design Pro > Ant Design > 自定义组件

### 页面布局标准

- **管理页面** 使用 ProTable + ProForm 的标准CRUD布局
- **详情页面** 使用 ProDescriptions 展示详细信息
- **表单页面** 使用 ProForm 提供丰富的表单控件

#### 常用命令

```bash
cd frontend         # 切换到前端目录
npm install         # 安装依赖
npm run start:dev   # 启动开发服务器
npm run start:test  # 启动测试服务器
npm run openapi     # 生成API客户端代码，需要后端先执行 make swag 并启动服务
npm run test        # 运行测试
```

## 后端开发规范

### 技术栈

- go-nunu

### 关键目录

后端代码都在`backend`目录下，执行任何与后端相关的指令，都应当先`cd backend`切换到此目录。

- `api/v1/` API接口定义，包含请求/响应结构
- `internal/handler/` 接收HTTP请求，调用service处理业务
- `internal/service/` 核心业务逻辑实现
- `internal/repository/` 数据访问抽象，与数据库交互
- `internal/model/` 数据库实体模型定义

#### 架构分层规范

api (接口定义) → handler (请求处理) → service (业务逻辑) → repository (数据访问) → model (数据模型)

#### 代码组织规范

- **接口定义** (api/v1/) - 定义请求/响应结构体，版本化管理
- **请求处理** (internal/handler/) - HTTP路由处理，参数验证，调用service
- **业务逻辑** (internal/service/) - 核心业务逻辑，事务控制，业务规则
- **数据访问** (internal/repository/) - 数据库CRUD操作，查询优化
- **数据模型** (internal/model/) - GORM模型定义，数据库表映射

#### 常用命令

```bash
cd backend      # 切换到后端目录
go mod tidy     # 安装依赖
make init       # 安装开发依赖包
make bootstrap  # 启动后端服务
make mock       # 生成mock测试代码
make test       # 运行测试
make swag       # 生成API文档
```

#### 日志查看

```bash
tail -f storage/logs/app.log
```
