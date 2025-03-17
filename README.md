<div align="center">

# 🚀 FFlow - 智能流程编排平台

**强大的 AI Agent 任务调度与工作流管理系统**

[![Go Version](https://img.shields.io/badge/Go-1.16+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-Apache%202.0-green.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

</div>

## ✨ 为什么选择 FFlow？

FFlow 不仅仅是一个流程编排工具，更是一个**智能任务调度生态系统**。它能够：

- 🤖 **无缝连接 AI 与业务**：专为 AI Agent 设计的工作流编排，让智能任务与业务流程完美结合
- 🔄 **事件驱动架构**：基于高性能事件机制，支持复杂业务场景下的实时响应
- 🛠️ **灵活的扩展性**：多语言脚本支持，兼容各种 Agent 实现和第三方服务集成
- 🏗️ **企业级架构**：基于 DDD 领域驱动设计，清晰的代码组织便于扩展和维护

## 🌟 产品亮点

### 多语言脚本执行

> 💡 **内置多语言执行引擎**，一套流程定义，多种语言实现

- 支持 Python, JavaScript, Go 等多种语言
- 安全的脚本隔离与资源控制
- 丰富的内置函数库

### 服务编排与集成

| 服务类型 | 描述 | 应用场景 |
|---------|------|----------|
| **Event** | 事件驱动的消费服务 | 异步任务处理、数据流转 |
| **RPC** | 远程过程调用服务 | 跨服务通信、微服务调用 |
| **Web** | HTTP/HTTPS RESTful服务 | API集成、前端交互 |
| **Timer** | 定时任务服务 | 周期性任务、延时执行 |

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

## 🚀 快速开始

### 环境要求

- Go 1.16+
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

4. **访问API文档**

浏览器打开 `http://localhost:8080/swagger/index.html` 查看API文档。

## 🔍 使用场景

### 业务流程自动化

- 客户服务流程自动化
- 基于协同工具的跨部门审批流程
- 数据处理与分析管道
- IoT 设备控制与监控

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

## 📄 许可证

本项目采用 [Apache License 2.0](LICENSE) 许可证。

---

<div align="center">

**FFlow** ⚡ **让工作流程更智能，让开发更高效**

[GitHub](https://github.com/fflow-tech/fflow) · [编排指南](https://github.com/fflow-tech/fflow/blob/main/docs/user-guide.md)

</div>
