# 🚀 FFlow 工作流系统使用指南

## 👋 欢迎使用 FFlow

FFlow 是一个强大的工作流引擎系统，让你能够轻松定义、部署和管理各种复杂场景下的自动化流程。本指南将帮助你快速入门并掌握 FFlow 的核心功能。

## 📑 目录

- 🔍 [快速参考](#快速参考)
- 📝 [流程定义基础](#流程定义基础)
- 🔌 [节点类型详解](#节点类型详解)
- ⏱️ [触发器配置](#触发器配置)
- 🔣 [表达式编写指南](#表达式编写指南)

## 🔍 快速参考

### 📊 数据引用快捷方式

| 缩写 | 含义 | 示例用法 |
|-----|-----|--------|
| `w.i` | 流程输入 | `${w.i.operator}` |
| `w.v` | 流程全局变量 | `${w.v.projectName}` |
| `this.o` | 当前节点的输出 | `${this.o.status}` |
| `t1.o` | 特定节点(t1)的输出 | `${t1.o.data.id}` |
| `t1.po` | 特定节点(t1)的轮询结果 | `${t1.po.data.completed}` |

### 🔗 节点级别支持的特殊引用

| 引用 | 含义 |
|-----|-----|
| `this.nodeInstID` | 当前节点的实例ID |
| `this.nodeRefName` | 当前节点的引用名称 |
| `this.owner` | 当前节点的所有者 |

## 📝 流程定义基础

流程是 FFlow 的核心概念，每个流程包含基本信息、输入参数、节点定义和触发器设置等。

### 📋 流程结构示例

```yaml
name: 简单审批流程
desc: 一个简单的审批流程示例
timeout:
  duration: 1d
  policy: TIME_OUT_WF

input:
- operator:
    options:
    - workflow
    - conductor
    default: workflow
    required: true

owner:
  wechat: user1;user2
  chatGroup: workflow-group

variables:
  projectName: 测试项目

nodes:
- start:
    type: ASSIGN
    assign:
    - variables:
        status: pending
    next: approval
    
- approval:
    type: SERVICE
    name: 审批节点
    args:
      protocol: HTTP
      method: POST
      url: https://api.example.com/approval
      body:
        requester: "${w.i.operator}"
        project: "${w.v.projectName}"
    next: end
```

### 📌 基本字段说明

| 字段 | 说明 |
|-----|-----|
| `name` | 流程名称 (必填) |
| `desc` | 流程描述 (必填) |
| `timeout` | 流程超时设置，包含 `duration` 和 `policy` |
| `input` | 流程输入参数定义 |
| `variables` | 全局变量定义 |
| `owner` | 流程所有者和通知接收方 |
| `nodes` | 流程节点列表 (必填) |
| `triggers` | 流程触发器定义 |
| `webhooks` | 流程事件webhook地址列表 |
| `subworkflows` | 子流程定义 |

## 🔌 节点类型详解

FFlow 支持多种类型的节点，每种类型具有特定的功能和配置方式。

### 🌐 SERVICE 节点

服务节点用于调用外部服务，支持 TRPC、HTTP、FAAS 等多种协议。

```yaml
name: HTTP调用示例
type: SERVICE
args:
  protocol: HTTP
  method: POST
  url: https://api.example.com
  body:
    params:
      width: 100
      height: 100
  headers:
    Content-Type: application/json
```

#### 🔄 轮询功能

可通过 `pollArgs` 配置轮询功能，用于监控异步任务的执行状态：

```yaml
pollArgs:
  protocol: HTTP
  method: GET
  url: https://api.example.com/status/${this.o.taskId}
  timeoutDuration: 30m
  pollCondition: ${this.po.status == "completed"}
  successCondition: ${this.po.result == "success"}
```

### 🔀 SWITCH 节点

条件分支节点，根据条件决定流程走向：

```yaml
type: SWITCH
switch:
- condition: ${w.i.type == "urgent"}
  next: urgentProcess
- condition: ${w.i.type == "normal"}
  next: normalProcess
next: defaultProcess  # 默认路径
```

### ⑂ FORK & JOIN 节点

并行执行多个分支，并在 JOIN 处汇合：

```yaml
# 并行开始
type: FORK
fork:
- branch1
- branch2

# 其他节点定义...

# 并行结束
type: JOIN
```

### 📊 更多节点类型

- 🔄 **TRANSFORM**: 数据转换节点
- 📝 **ASSIGN**: 变量赋值节点
- 🔗 **REF**: 引用节点
- 🔀 **EXCLUSIVE_JOIN**: 互斥汇合节点
- 📦 **SUB_WORKFLOW**: 子流程节点
- ⏱️ **WAIT**: 等待节点

## ⏱️ 触发器配置

触发器用于自动启动流程或执行特定操作。

### 🕒 定时触发器

```yaml
type: timer
expr: 0 15 10 ? * MON-FRI  # 工作日上午10:15触发
actions:
  sw:
    action: START_WORKFLOW
    allowDays: WORKDAY
    args:
      name: 每日构建
      input:
        operator: timer
```

### 📣 事件触发器

```yaml
type: event
event: ProjectUpdateEvent
actions:
- a1:
    action: RERUN_NODE
    condition: ${event.type == w.v.type}
    args:
      node: buildNode
```

### ⏰ 常用定时表达式

| 含义 | 表达式 |
|-----|-------|
| 每小时的10分30秒 | `30 10 * * * ?` |
| 每天1点10分30秒 | `30 10 1 * * ?` |
| 每周一到周五的10点15分 | `0 15 10 ? * MON-FRI` |
| 每天5-15点整点 | `0 0 5-15 * * ?` |

## 🔣 表达式编写指南

FFlow 使用表达式语言来访问和处理数据。表达式使用 `${...}` 格式。

### 📝 基本用法

```
${w.i.user == "admin"}
${t1.o.count > 0 && w.v.status == "active"}
```

### 🛠️ 内置函数

#### 📝 sprintf

用于格式化字符串：

```
${sprintf("%s项目-%s版本发布", w.v.project, w.v.version)}
```

#### 🕒 curtimeformat

格式化当前时间：

```
${curtimeformat("0102")}  # 输出当天日期，如0927
```

---

> 💡 **提示**：需要更详细的信息或有任何问题，请参考完整文档或联系技术支持团队。
