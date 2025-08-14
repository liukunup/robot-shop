# 开发指导手册

技术栈：`Go` + `Ant Design Pro`

推荐使用`VS Code`作为集成开发环境，打开`*.code-workspace`工作空间文件，一键开启开发之旅吧~

如果你的项目托管在`GitHub`，推荐使用`Codespaces`进行开发~

## 前端

### 前置准备工作

- 安装 Node.js 环境

- 运行 npm install 安装依赖

## 后端

### 前置准备工作

- 安装 Golang 环境

### 文档

### 调试

1. 安装`dlv`

```shell
go install github.com/go-delve/delve/cmd/dlv@latest
```

2. 在有需要的地方打断点（就是行号旁边鼠标点一下，标记小红点）

3. 左侧侧边栏找到调试按钮，点击`Debug`运行

4. 通过调用接口触发程序执行，进入到断点位置，查看变量、堆栈等信息

## 构建

### 本地构建

- 安装 Docker 环境

### [GitHub Actions](https://docs.github.com/zh/actions)

## 部署

- Docker

- Docker Compose

- Kubernetes

- Helm
