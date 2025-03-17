# FFlow 流程平台

## 项目简介

FFlow 是一个使用 Go 语言编写的高性能流程编排工具，专注于 AI Agent 任务调度与通用工作流程管理。项目采用领域驱动设计（DDD）架构思想，并使用类似大仓（Monorepo）的方式组织多个微服务，主要包括：
- 基于事件驱动的流程引擎，支持 AI Agent 任务编排
- 多语言脚本执行引擎，兼容多种 Agent 实现
- 灵活的服务集成能力，支持 Agent 服务调用
- AI 任务编排与资源调度

项目采用 DDD 分层架构，将复杂的业务逻辑解耦为清晰的领域模型，同时通过大仓方式管理多个微服务，实现代码共享和版本一致性

## 核心特性

- **AI Agent 任务编排**：
  - 基于事件驱动的 Agent 工作流编排
  - Agent 任务状态追踪与管理
  - Agent 间协作流程编排
  - 灵活的任务调度策略
- **多语言支持**：内置多语言脚本执行引擎，支持多种 Agent 实现
- **服务编排**：支持多种服务调用方式
  - Event：事件消费服务
  - RPC：远程过程调用服务
  - Web：HTTP/HTTPS 服务
  - Timer：定时任务服务

## 项目结构

项目采用 DDD 分层架构和大仓（Monorepo）管理方式组织代码：

```
├── api/            # 接口定义层，包含各服务的 API 定义
│   ├── foundation/ # 基础设施服务接口
│   └── workflow-app/ # 工作流应用服务接口
├── service/        # 领域实现层，包含多个微服务
│   ├── cmd/       # 应用层：服务入口和启动配置
│   ├── internal/  # 领域层：领域模型和业务逻辑
│   ├── pkg/       # 基础设施层：技术组件和工具
│   └── test/      # 测试代码
└── deployer/      # 部署配置层：服务部署和运维配置
```

## 快速开始

### 环境要求

- Go 1.16+

### 本地开发

1. 克隆项目
```bash
git clone https://github.com/fflow-tech/fflow.git
cd fflow
```

2. 安装依赖
```bash
go mod download
```

3. 运行服务
```bash
cd service/cmd/demo-app/blank-demo
go run main.go
```

## 服务说明

### 工作流应用服务 (workflow-app)

工作流应用服务是 FFlow 的核心服务，提供：
- 流程定义和管理
- 事件处理和路由
- 任务调度和执行

### 基础服务 (foundation)

提供基础设施支持：
- 认证授权 (auth)
- 函数计算 (faas)
- 定时任务 (timer)

## 调试说明

### Nocalhost 调试

如果生成的调试 pod 没有流量进入，可以尝试修改 pod 的 labels 让服务匹配上。

## 贡献指南

1. Fork 本仓库
2. 创建您的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交您的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开一个 Pull Request

## 许可证

本项目采用 [Apache License 2.0](LICENSE) 许可证。