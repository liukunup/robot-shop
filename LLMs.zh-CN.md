# LLMs (beta)

本文档旨在为大模型辅助编程提供标准化的项目信息，使其能够更准确、高效的理解用户指令，并严格按照用户需求生成或修改代码。

## 推理规范

1. [意图识别] 当接收到用户指令时，你需要先思考用户的意图；
2. [任务拆解] 你需要根据用户的意图拆解子任务，然后根据子任务的要求去生成或修改代码；
3. [路径规划] 当存在多个子任务时，你需要规划任务的执行顺序，确保每个子任务都能正确执行；
4. [严格模式] 你需要严格按照用户需求生成或修改代码，不得修改任何非代码文件。
5. [迭代开发] 你需要先检查项目中是否已经存在类似代码，并在此基础上进行迭代开发。

## 目录树
```
<projectRoot>
├── .github                # GitHub 配置文件
├── backend                        # 后端项目的目录
│   ├── README.md                  # 说明文档
│   ├── README_zh.md               # 说明文档（中文）
│   ├── LICENSE                    # 许可文件
│   ├── .dockerignore              # Docker 忽略文件
│   ├── .gitignore                 # Git 忽略文件
│   ├── .gitattributes             # Git 属性文件
│   ├── go.mod                     # Go Module 依赖管理文件
│   ├── go.sum                     # Go Module 依赖管理文件
│   ├── Makefile                   # 提供了一些开发或调试相关的操作命令
│   ├── api                        # 接口定义，包括请求和响应的结构体
│   │   └── v1                     # v1版本
│   ├── cmd                        # 应用的入口点，根据不同的命令执行不同的操作 
│   │   ├── migration              # 用于执行数据库迁移操作
│   │   │   ├── main.go            # 应用程序入口
│   │   │   └── wire               # 依赖注入配置
│   │   │       ├── wire.go        # 依赖注入配置文件
│   │   │       └── wire_gen.go    # 该文件是由 Wire 工具自动生成的，请勿手动修改
│   │   ├── server                 # 用于启动后端服务器
│   │   │   ├── main.go            # 应用程序入口
│   │   │   └── wire               # 依赖注入配置
│   │   │       ├── wire.go        # 依赖注入配置文件
│   │   │       └── wire_gen.go    # 该文件是由 Wire 工具自动生成的，请勿手动修改
│   │   └── task                   # 用于执行后台任务
│   │       ├── main.go            # 应用程序入口
│   │       └── wire               # 依赖注入配置
│   │           ├── wire.go        # 依赖注入配置文件
│   │           └── wire_gen.go    # 该文件是由 Wire 工具自动生成的，请勿手动修改
│   ├── config                     # 应用配置文件，根据不同的环境提供不同的配置
│   │   ├── local.yml              # 本地环境，用于开发、测试
│   │   ├── gray.yml               # 灰度环境，用于上线前的测试
│   │   └── prod.yml               # 生产环境，用于生产
│   ├── deploy                     # 构建和部署相关的配置文件
│   │   ├── build                  # 构建
│   │   │   └── Dockerfile         # 镜像构建文件
│   │   └── docker-compose         # 部署
│   │       └── docker-compose.yml # 镜像部署文件
│   ├── docs                       # 接口文档，使用 Swagger 生成
│   │   ├── docs.go                # 该文件是由 Swag 工具自动生成的，请勿手动修改
│   │   ├── swagger.json           # 生成 Swagger 文档的 JSON 配置文件
│   │   └── swagger.yaml           # 生成 Swagger 文档的 YAML 配置文件
│   ├── internal                   # 应用的核心模块，包含了各种业务逻辑的实现
│   │   ├── handler                # 处理 HTTP 请求的实现，负责接收请求并调用相应的服务进行处理
│   │   │   ├── handler.go         # Handler 类
│   │   │   └── user.go            # User Handler
│   │   ├── job                    # 包含了异步事件的实现
│   │   │   ├── job.go             # Job 类
│   │   │   └── user.go            # User 异步事件
│   │   ├── middleware             # 包含了中间件的实现，用于处理请求的预处理和后处理
│   │   │   ├── cors.go            # 处理跨域请求
│   │   │   ├── jwt.go             # 处理 JWT 认证
│   │   │   ├── log.go             # 处理日志记录
│   │   │   └── sign.go            # 处理签名验证
│   │   ├── model                  # 包含了数据模型的定义
│   │   │   └── user.go            # User 数据库表
│   │   ├── repository             # 包含了数据访问层的实现，负责与数据库进行交互
│   │   │   ├── repository.go      # Repository 类
│   │   │   └── user.go            # User DAO层
│   │   ├── server                 # 包含了服务端的实现，用于接收请求并处理
│   │   │   ├── http.go            # 处理 HTTP 请求
│   │   │   ├── job.go             # 处理后台任务
│   │   │   ├── migration.go       # 处理数据库迁移
│   │   │   └── task.go            # 处理后台任务
│   │   ├── service                # 包含了业务逻辑的实现，负责处理具体的业务操作
│   │   │   ├── service.go         # Service 类
│   │   │   └── user.go            # User 业务逻辑类
│   │   └── task                   # 包含了定时任务的实现
│   │       ├── task.go            # Task 类
│   │       └── user.go            # User 定时任务
│   ├── pkg                        # 包含了一些通用的功能和工具
│   │   ├── app                    # 包含了应用程序的配置和初始化
│   │   │   └── app.go
│   │   ├── config                 # 包含了应用程序的配置文件
│   │   │   └── config.go
│   │   ├── jwt                    # 包含了 JWT 认证的实现
│   │   │   └── jwt.go             # JWT Token 生成和解析
│   │   ├── log                    # 包含了日志记录的实现
│   │   │   └── log.go
│   │   ├── server                 # 包含了服务器的实现
│   │   │   ├── grpc               # 包含了 gRPC 服务器的实现
│   │   │   │   └── grpc.go
│   │   │   ├── http               # 包含了 HTTP 服务器的实现
│   │   │   │   └── http.go
│   │   │   └── server.go
│   │   ├── sid                    # 包含了 SID 生成器的实现
│   │   │   ├── convert.go
│   │   │   └── sid.go
│   │   └── zapgorm2               # 包含了 ZapGorm2 的实现
│   │       └── zapgorm2.go
│   ├── scripts                    # 包含了一些脚本文件，用于项目的构建、测试和部署等操作
│   │   └── README.md
│   ├── storage                    # 该模块用于存储文件或其他静态资源。
│   │   └── nunu-test.db           # 测试数据库文件
│   ├── test                       # 该模块包含了各个模块的单元测试，按照模块划分子目录
│   │   ├── mocks                  # 模拟数据
│   │   │   ├── repository         # 模拟数据库 DAO 层
│   │   │   │   ├── repository.go
│   │   │   │   └── user.go        # 模拟 User DAO 层
│   │   │   └── service            # 模拟业务逻辑
│   │   │       └── user.go
│   │   └── server                 # 服务测试
│   │       ├── handler            # 测试 Handler
│   │       │   ├── main_test.go
│   │       │   └── user_test.go
│   │       ├── repository         # 测试 Repository
│   │       │   └── user_test.go
│   │       └── service            # 测试 Service
│   │           └── user_test.go
│   └── web
│       └── index.html
├── frontend                           # 前端项目的目录
│   ├── README.md                      # 说明文档
│   ├── mock                           # 模拟数据，用于模拟后端接口的返回数据，便于前端开发和测试
│   │   └── userAPI.ts                 # 这是一个代码样例
│   ├── package.json                   # 依赖管理文件，通常不需要手动修改
│   ├── pnpm-lock.yaml                 # 依赖锁定文件，通常不需要手动修改
│   ├── src                            # 源代码目录
│   │   ├── access.ts                  # 权限
│   │   ├── app.ts                     # 运行时配置
│   │   ├── assets                     # 静态资源
│   │   ├── components                 # 全局组件
│   │   │   └── Guide                  # 这是一个示例组件
│   │   │       ├── Guide.less         # 样式
│   │   │       ├── Guide.tsx          # 组件
│   │   │       └── index.ts
│   │   ├── constants                  # 常量
│   │   │   └── index.ts               # 全局常量
│   │   ├── models                     # 模型
│   │   │   └── global.ts              # 全局共享数据
│   │   ├── pages                      # 页面目录
│   │   │   ├── Access                 # 访问控制
│   │   │   │   └── index.tsx
│   │   │   ├── Home                   # 默认主页
│   │   │   │   ├── index.less
│   │   │   │   └── index.tsx
│   │   │   └── Table                  # 表格页面
│   │   │       ├── components         # 表格页面的构成组件，如创建、更新的表单等，这是一个代码样例
│   │   │       │   ├── CreateForm.tsx # 创建表单
│   │   │       │   └── UpdateForm.tsx # 更新表单
│   │   │       └── index.tsx          # 表格页
│   │   ├── services                   # 服务目录，用于封装网络请求等，均由 OpenAPI 生成，请勿手动修改
│   │   │   └── demo                   # 忽略此目录
│   │   │       ├── UserController.ts  # User Controller，用于处理用户的增删改查等操作，这是一个代码样例
│   │   │       ├── index.ts           # 服务入口
│   │   │       └── typings.d.ts       # 类型定义文件
│   │   └── utils                      # 工具
│   │       └── format.ts              # 格式化工具
│   ├── tsconfig.json                  # TypeScript 配置文件
│   └── typings.d.ts                   # 类型定义文件
├── .gitignore             # Git忽略文件
├── openhrm.code-workspace # VSCode工作区文件
├── LLMs.zh-CN.md          # 本文档
├── LICENSE                # 许可协议
└── README.md              # 项目简介
```

## 技术栈

- 前端
  - [Ant Design Pro](https://pro.ant.design/zh-CN/docs/overview) 基于 Ant Design 和 UmiJS 封装的一整套企业级中后台前端/设计解决方案
    - [UmiJS](https://umijs.org/docs/guides/getting-started) 可扩展的企业级前端应用框架，用来生成前端工程

- 后端
  - [Go](https://go.dev) 作为后端编程语言
  - [Nunu](https://github.com/go-nunu/nunu/blob/main/README_zh.md) 应用脚手架，用来生成后端工程
    - [Gin](https://github.com/gin-gonic/gin) Web 框架
    - [Gorm](https://github.com/go-gorm/gorm) ORM 框架
    - [Wire](https://github.com/google/wire) 依赖注入
    - [Viper](https://github.com/spf13/viper) 配置管理
    - [Zap](https://github.com/uber-go/zap) 日志
    - [Golang-jwt](https://github.com/golang-jwt/jwt) JWT
    - [Go-redis](https://github.com/go-redis/redis) Redis
    - [Testify](https://github.com/stretchr/testify) 测试
    - [Sonyflake](https://github.com/sony/sonyflake) 分布式 ID
    - [Gocron](https://github.com/go-co-op/gocron) 定时任务
    - [Go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) SQL Mock
    - [Gomock](https://github.com/golang/mock) Mock
    - [Swaggo](https://github.com/swaggo/swag) Swagger
    - [Pitaya](https://github.com/topfreegames/pitaya) RPC
    - [Casbin](https://github.com/casbin/casbin) RBAC

- 构建
  - GitHub Actions
  - Docker

- 部署
  - Docker Compose
  - Kubernetes

## 设计规范

当用户要求生成前端页面或组件，同时又不提供设计参考时，你应当遵循以下设计规范：

基于「自然」、「确定性」、「意义感」、「生长性」四大设计价值观，通过模块化解决方案，降低冗余的生产成本，让设计专注于更好的用户体验。

> 摘抄自 https://ant-design.antgroup.com/docs/spec/values-cn

## 编码规范

- 在输出任何代码前，先学习技术栈如何正确使用，避免出现不符合框架规范的结果。

- 实现用户需求时，应当检查前后端是否都需要编码，完整的实现用户功能。

- 当你尝试新增或修改某个功能时，需要考虑以下业务实现的关联影响：

   - 前端：Page -> Component -> API 影响点位于哪个页面、哪个组件、哪个API
   - 后端：Handler -> Service -> Repository -> Model 影响点位于哪个Handler、哪个Service、哪个Repository、哪个Model

- 开发前端页面时，你应当考虑优先使用 ProComponents > Ant Design Pro > Ant Design > 自定义编写 的组件，避免重复造轮子。

- 开发后端接口时，你应当考虑遵循RESTful风格，同时考虑性能开销，避免出现性能瓶颈或不必要的计算。

- 当发生代码变更时，你应当考虑撰写测试用例进行验证，以确保代码的正确性和稳定性。

- 你应当考虑使用设计模式，灵活处理不同业务场景带来的挑战。

- 你应当考虑代码的可读性、可维护性，避免出现冗余代码或复杂的逻辑。

- 你应当考虑代码的兼容性，避免出现代码兼容问题。

- 你应当考虑代码的安全性，避免出现安全漏洞或数据泄露。

- 你应当考虑代码的国际化，避免出现语言差异或翻译错误。