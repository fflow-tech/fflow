<div align="center">

# 🚀 FFlow - 智能流程编排平台

**强大的 AI Agent 任务调度与工作流管理系统**

[![Go Version](https://img.shields.io/badge/Go-1.16+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-Apache%202.0-green.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

[English](README.md) | [简体中文](README.zh-CN.md)

</div>

## ✨ 为什么选择 FFlow？

FFlow 不仅仅是一个流程编排工具，更是一个**智能任务调度生态系统**。它能够：

- 🤖 **无缝连接 AI 与业务**：专为 AI Agent 设计的工作流编排，让智能任务与业务流程完美结合
- 🔄 **事件驱动架构**：基于高性能事件机制，支持复杂业务场景下的实时响应
- 🛠️ **灵活的扩展性**：多语言脚本支持，兼容各种 Agent 实现和第三方服务集成
- 🏗️ **企业级架构**：基于 DDD 领域驱动设计，清晰的代码组织便于扩展和维护

## 🌟 产品亮点

### 多语言脚本执行

> 💡 **内置多语言执行引擎**

- 支持 Python, JavaScript, Go 等多种语言
- 安全的脚本隔离与资源控制
- 丰富的内置函数库

### 统一流程定义

> 💡 **一次定义，随处运行**

- 同一套工作流定义同时支持分布式服务端和本地运行
- 无缝切换开发环境和生产环境
- 标准化的工作流定义格式

## 🏗️ 架构设计

FFlow 采用**领域驱动设计(DDD)**和**大仓(Monorepo)管理**方式组织代码，实现高内聚低耦合的系统架构。

### 项目结构

```
├── api/            # 接口定义层：服务API协议
│   ├── foundation/ # 基础设施服务接口
│   └── workflow-app/ # 工作流应用服务接口
├── service/        # 实现层：服务核心逻辑
│   ├── cmd/       # 应用层：服务入口
│   │   ├── foundation/  # 基础服务
│   │   └── workflow-app/ # 工作流服务
│   ├── internal/  # 领域层：业务模型
│   ├── pkg/       # 基础层：通用组件
│   └── test/      # 测试代码
└── deployer/      # 部署层：运维配置
```

### 命令行执行

FFlow 提供了便捷的命令行工具 `fflow-cli` 用于本地执行工作流:

```bash
# 执行工作流
fflow-cli -f <工作流定义文件> -i <输入参数文件>

# 示例
fflow-cli -f examples/example-http.json -i examples/example-http-input.json
fflow-cli -f examples/example-http.yaml -i examples/example-http-input.json
```

CLI 工具会自动在当前目录下创建 `.fflow` 文件夹用于存储工作流定义和实例数据。

#### 主要参数说明

- `-f`: 工作流定义文件路径
- `-i`: 工作流输入参数文件路径
- `-config.path`: 配置文件路径，默认为 `.fflow/`
- `-def.path`: 工作流定义目录，默认为 `.fflow/definitions`
- `-inst.path`: 工作流实例目录，默认为 `.fflow/instances`

## 🚀 快速开始

### 一键安装

使用以下命令可以快速安装 FFlow：

```bash
curl -sSL https://raw.githubusercontent.com/fflow-tech/fflow/main/install.sh | bash
```

### 环境要求

- Go 1.23+
- Docker (可选，用于容器化部署)
- Kubernetes (可选，用于集群部署)

### 本地开发

1. **获取代码**

```bash
git clone https://github.com/fflow-tech/fflow.git
cd fflow
```

2. **安装依赖**

```bash
go mod download
```

3. **运行示例服务**

```bash
cd service/cmd/demo-app/blank-demo
go run main.go
```

## 🔍 使用场景

### 业务流程自动化

- 客户服务流程自动化
- 基于协同工具的跨部门审批流程
- 数据处理与分析管道
- IoT 设备控制与监控
- AI Agent 的调度和执行

## 📊 性能指标

- **高吞吐量**：单节点支持每秒数千任务调度
- **低延迟**：任务调度平均延迟<10ms
- **高可用**：支持多节点部署，无单点故障
- **水平扩展**：线性扩展性能，按需增加节点

## 🔨 调试与开发

### Nocalhost调试

> 如果生成的调试pod没有流量进入，可以尝试修改pod的labels让服务匹配上。

## 🤝 贡献指南

我们欢迎各种形式的贡献，无论是新功能、文档改进还是bug修复！

1. Fork 本仓库
2. 创建您的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交您的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开一个Pull Request

## 📋 路线图

- [ ] WebUI工作流可视化设计器
- [ ] 更多AI模型集成支持
- [ ] 分布式任务调度优化
- [ ] 更完善的监控与报警系统
- [ ] 支持 MCP(Model Context Protocol) 工具的调度

## 📄 许可证

本项目采用 [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0) 许可证。

---

<div align="center">

**FFlow** ⚡ **让工作流程更智能，让开发更高效**

[GitHub](https://github.com/fflow-tech/fflow) · [编排指南](https://github.com/fflow-tech/fflow/blob/main/docs/user-guide.md)

</div> 