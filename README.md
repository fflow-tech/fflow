<div align="center">

# 🚀 FFlow - Intelligent Flow Orchestration Platform

**Powerful AI Agent Task Scheduling and Workflow Management System**

[![Go Version](https://img.shields.io/badge/Go-1.16+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-Apache%202.0-green.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

[English](README.md) | [简体中文](README.zh-CN.md)

</div>

## ✨ Why FFlow?

FFlow is not just a workflow orchestration tool, but an **intelligent task scheduling ecosystem**. It can:

- 🤖 **Seamless AI Integration**: Workflow orchestration designed for AI Agents, perfectly combining intelligent tasks with business processes
- 🔄 **Event-Driven Architecture**: High-performance event mechanism supporting real-time response in complex business scenarios
- 🛠️ **Flexible Extensibility**: Multi-language script support, compatible with various LLM API implementations and third-party service integrations
- 🏗️ **Enterprise Architecture**: Based on DDD (Domain-Driven Design) for clear code organization and maintainability

## 🌟 Key Features

### Multi-Language Script Execution

> 💡 **Built-in Multi-Language Engine**

- Supports Python, JavaScript, Go, and more
- Secure script isolation and resource control
- Rich built-in function library

### Unified Workflow Definition

> 💡 **Define Once, Run Anywhere**

- Same workflow definition supports both distributed server and local execution
- Seamless switching between development and production environments
- Standardized workflow definition format

## 🏗️ Architecture Design

FFlow adopts **Domain-Driven Design (DDD)** and **Monorepo Management** to achieve a highly cohesive and loosely coupled system architecture.

### Project Structure

```
├── api/            # Interface Layer: Service API Protocols
│   ├── foundation/ # Infrastructure Service Interfaces
│   └── workflow-app/ # Workflow Application Service Interfaces
├── service/        # Implementation Layer: Core Service Logic
│   ├── cmd/       # Application Layer: Service Entry Points
│   │   ├── foundation/  # Foundation Services
│   │   └── workflow-app/ # Workflow Services
│   ├── internal/  # Domain Layer: Business Models
│   ├── pkg/       # Foundation Layer: Common Components
│   └── test/      # Test Code
└── deployer/      # Deployment Layer: Operations Configuration
```

### Command Line Execution

FFlow provides a convenient CLI tool `fflow-cli` for local workflow execution:

```bash
# Execute workflow
fflow-cli -f <workflow-definition-file> -i <input-parameters-file>

# Example
fflow-cli -f examples/example-http.json -i examples/example-http-input.json
fflow-cli -f examples/example-http.yaml -i examples/example-http-input.json
fflow-cli -f examples/example-openai.yaml -i examples/example-openai-input.json
```

The CLI tool automatically creates a `.fflow` folder in the current directory for storing workflow definitions and instance data.

#### Main Parameters

- `-f`: Workflow definition file path
- `-i`: Workflow input parameters file path
- `-config.path`: Configuration file path, defaults to `.fflow/`
- `-def.path`: Workflow definition directory, defaults to `.fflow/definitions`
- `-inst.path`: Workflow instance directory, defaults to `.fflow/instances`

## 🚀 Quick Start

### One-Click Installation

Install FFlow quickly using the following command:

```bash
curl -sSL https://raw.githubusercontent.com/fflow-tech/fflow/main/install.sh | bash
```

### Requirements

- Go 1.23+
- Docker (optional, for containerized deployment)
- Kubernetes (optional, for cluster deployment)

### Local Development

1. **Get the Code**

```bash
git clone https://github.com/fflow-tech/fflow.git
cd fflow
```

2. **Install Dependencies**

```bash
go mod download
```

3. **Run Example Service**

```bash
cd service/cmd/demo-app/blank-demo
go run main.go
```

## 🔍 Use Cases

### Business Process Automation

- Customer service process automation
- Cross-department approval processes based on collaboration tools
- Data processing and analysis pipelines
- IoT device control and monitoring
- AI Agent scheduling and execution

## 📊 Performance Metrics

- **High Throughput**: Thousands of tasks per second per node
- **Low Latency**: Average task scheduling latency <10ms
- **High Availability**: Multi-node deployment support, no single point of failure
- **Horizontal Scaling**: Linear performance scaling, add nodes as needed

## 🔨 Debugging and Development

### Nocalhost Debugging

> If the generated debug pod has no incoming traffic, try modifying the pod's labels to match the service.

## 🤝 Contributing

We welcome all forms of contributions, whether they're new features, documentation improvements, or bug fixes!

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📋 Roadmap

- [ ] WebUI workflow visual designer
- [ ] More AI model integration support
- [ ] Distributed task scheduling optimization
- [ ] Enhanced monitoring and alerting system
- [ ] Support for MCP (Model Context Protocol) tool scheduling

## 📄 License

This project is licensed under the [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0).

---

<div align="center">

**FFlow** ⚡ **Making Workflows Smarter, Development More Efficient**

[GitHub](https://github.com/fflow-tech/fflow) · [Orchestration Guide](https://github.com/fflow-tech/fflow/blob/main/docs/user-guide.md)

</div>
