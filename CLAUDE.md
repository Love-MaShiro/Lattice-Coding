# CLAUDE.md

## 项目概述

Lattice-Coding 是一个面向实验室内部使用的 AI Agent 编码平台，目标用户规模约 20 人。

项目定位：

- 内部使用的 AI Agent 编码平台
- 支持模型管理、Agent 配置、对话、工具调用、知识库检索、Run 执行记录
- 后端采用模块化单体架构
- 前端采用 Vue 3 管理控制台
- 不追求企业级 SaaS、完整 Web IDE、插件市场、计费系统和复杂低代码工作流

技术栈：

| 层 | 技术 |
|---|---|
| 后端 | Go + Eino |
| Web 框架 | Gin |
| 主数据库 | MySQL |
| 缓存 / 队列 / 事件流 | Redis |
| 向量数据库 | PostgreSQL + pgvector |
| ORM | GORM |
| 前端 | Vue 3 + TypeScript + Vite + Element Plus |
| 部署 | Docker + K8s |

---

## 核心功能模块

| 模块 | 说明 |
|------|------|
| provider | 模型提供商管理，支持 OpenAI、Claude、DeepSeek、Qwen、Ollama、本地模型、OpenAI-compatible API。 |
| agent | Agent 创建与配置，管理 Prompt、默认模型、工具权限、最大执行步数、默认测试命令。 |
| chat | 对话引擎，负责会话、消息、SSE 流式输出和上下文回填。 |
| run | Agent Run 任务执行，管理任务创建、排队、状态机、中断、恢复、最终结果。 |
| mcp | MCP / Tool 工具管理与调用，负责工具注册、工具 Schema、工具调用分发。 |
| workflow | 简版工作流编排，第一阶段只支持配置化线性流程。 |
| knowledge | 知识库与 RAG，负责文档管理、chunk、embedding、pgvector 检索。 |
| safety | 权限与安全控制，负责工具权限、命令白名单、路径限制、人工确认。 |
| audit | 审计日志，记录模型调用、工具调用、命令执行、文件修改、Run 轨迹。 |
| common | 公共基础能力，只放无业务归属的通用代码。 |

---

## 砍掉的功能

- 完整 Web IDE
- 企业级多租户系统
- 插件市场
- 计费系统
- 模型微调平台
- 复杂可视化工作流拖拽
- 云端并行 Agent 集群
- 自动长期记忆
- 默认联网搜索

---

## 后端代码组织

```text
lattice-coding/
├── cmd/
│   ├── api/
│   │   └── main.go
│   └── worker/
│       └── main.go
│
├── internal/
│   ├── app/
│   │   ├── bootstrap.go
│   │   ├── dependencies.go
│   │   ├── router.go
│   │   └── modules.go
│   │
│   ├── common/
│   │   ├── config/
│   │   ├── errors/
│   │   ├── response/
│   │   ├── db/
│   │   ├── redis/
│   │   ├── logger/
│   │   ├── middleware/
│   │   └── util/
│   │
│   ├── modules/
│   │   ├── provider/
│   │   ├── agent/
│   │   ├── chat/
│   │   ├── run/
│   │   ├── mcp/
│   │   ├── workflow/
│   │   ├── knowledge/
│   │   ├── safety/
│   │   └── audit/
│   │
│   └── runtime/
│       ├── eino/
│       ├── llm/
│       ├── tool/
│       └── event/
│
├── migrations/
│   ├── mysql/
│   └── postgres/
│
├── configs/
├── deployments/
├── scripts/
└── docs/

---

### 模块内部组织

internal/modules/{module}/
├── api/
│   ├── handler.go
│   ├── request.go
│   ├── response.go
│   └── router.go
│
├── application/
│   ├── command_service.go
│   ├── query_service.go
│   ├── command.go
│   ├── query.go
│   ├── dto.go
│   └── assembler.go
│
├── domain/
│   ├── entity.go
│   ├── value_object.go
│   ├── repository.go
│   ├── service.go
│   └── event.go
│
├── infra/
│   └── persistence/
│       ├── model.go
│       ├── repository_impl.go
│       └── converter.go
│
└── module.go

简单模块可以省略空文件，但不得改变职责边界。
---

### 各层职责边界

| 层 | 职责 | 禁止 |
| ------ | -------------------------------------------------- | ------------------------------------------------------------------ |
| api | 处理 HTTP 请求、参数校验、响应转换、路由注册；只调用本模块 application。 | 不写业务逻辑；不访问 MySQL、Redis、pgvector；不直接调用 Eino、Tool、Mapper；不调用其他模块 infra。 |
| application | 处理用例编排、事务控制、跨模块调用；调用本模块 domain 和 repository 接口。 | 不直接写 SQL；不直接调用其他模块 infra；不返回 PO / Entity 给前端；不依赖 Gin Context。 |
| domain | 定义实体、值对象、领域规则、Repository 接口；维护业务不变量。 | 不依赖 Gin、MySQL、Redis、pgvector、Eino；不出现 SQL；不调用其他模块；不处理 HTTP、配置、日志、外部 SDK。 |
| infra  | 实现 domain.Repository；负责 PO、数据库访问、SQL/GORM、数据库转换。 | 不写业务流程；不被其他模块直接调用；不向外暴露 PO；不跨模块访问其他模块数据库表。 |
| runtime | 封装 Eino、LLM、Tool、Redis Stream 等运行时能力。| 不直接处理业务规则；不直接处理 HTTP；不绕过 application 层被 api 调用；不直接决定权限策略。|
| common| 放通用错误、响应、配置、日志、DB、Redis、中间件等无业务归属能力。| 不放 Agent、Run、Provider、Project 等业务逻辑；不创建业务 Helper。|

### 跨模块调用规则

chat/application      → run/application
run/application       → agent/application
run/application       → provider/application
run/application       → mcp/application
run/application       → safety/application
knowledge/application → provider/application
mcp/application       → safety/application

跨模块只能传递 DTO，不允许传递 Entity、PO、Repository、Mapper。

## Go 设计原则

| 原则             | 说明                                                 |
| -------------- | -------------------------------------------------- |
| Go 优先          | 不把 Go 写成没有注解的 Spring Boot。                         |
| 少抽象            | 一层能解决的不要拆成两层，不引入不必要的 manager/helper/processor。     |
| 接口按需定义         | 只有存在多个实现、需要 mock、隔离外部依赖时才定义 interface。             |
| 显式依赖           | 通过构造函数传递依赖，不使用隐式全局对象。                              |
| 标准库优先          | 能用标准库解决的，不引入第三方依赖。                                 |
| 不滥用泛型          | 泛型只用于 Result、PageResult、Cache 等明显通用结构。             |
| 不 panic 处理业务错误 | 业务错误必须通过 return error 显式返回。                        |
| context 贯穿调用链  | HTTP、DB、Redis、LLM、Shell、Tool 全部传递 context.Context。 |

### 统一响应规范

目录：

internal/common/response/

推荐结构：
```go
type Result[T any] struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Data    T      `json:"data,omitempty"`
}
```
| 项          | 规范                                                        |
| ---------- | --------------------------------------------------------- |
| 成功响应       | 使用 `response.OK(c, data)`                                 |
| 失败响应       | 使用 `response.Fail(c, err)`                                |
| 分页响应       | 使用 `response.Page(c, pageData)` 或 `Result[PageResult[T]]` |
| 成功码        | 统一使用 `"0"`                                                |
| HTTP 状态码   | 普通业务错误可以 HTTP 200；认证、权限、系统错误可使用真实 HTTP 状态码                |
| Data 字段    | 使用 `Data T`，不使用 `Data *T`                                 |
| PageResult | 不继承 Result，作为 data 内容返回                                   |

不要定义实例方法 result.Ok(...)
不要模仿 Java 的 PageResult extends Result
不要在业务层手动拼 JSON 响应

### 错误码与业务错误规范

目录：

internal/common/errors/
├── code.go
├── message.go
├── biz_error.go
└── convert.go

规则：
| 项     | 规范                                                             |
| ----- | -------------------------------------------------------------- |
| 错误类型  | 使用 `BizError`，不使用 `BizException`                               |
| 错误返回  | 通过 `return error` 显式返回                                         |
| 错误码   | 使用五位字符串，成功为 `"0"`                                              |
| 底层错误  | 使用 `Unwrap()` 保留 cause                                         |
| 业务错误  | 使用 `errors.New(code)` / `errors.NewWithMessage(code, message)` |
| 未知错误  | 统一转换为 `InternalError`                                          |
| panic | 只允许 Recovery 兜底捕获，不允许业务逻辑主动 panic                              |


错误码分类：

| 分类                  | 范围            |
| ------------------- | ------------- |
| 通用错误                | `40000-40999` |
| Agent / Run 错误      | `41000-41999` |
| 模型调用错误              | `42000-42999` |
| 工具调用错误              | `43000-43999` |
| 文件 / Git / Shell 错误 | `44000-44999` |
| 知识库错误               | `45000-45999` |
| 系统错误                | `50000+`      |

## LLM 调用规范

目录：

internal/runtime/llm/
├── executor.go
├── client.go
├── router.go
├── retry.go
├── breaker.go
└── types.go

规则：
| 项         | 规范                                      |
| --------- | --------------------------------------- |
| 非流式调用     | 必须走 `llm-pool`                          |
| 流式 SSE 调用 | 必须走 `llm-stream`                        |
| 并发控制      | 使用 goroutine + channel/semaphore        |
| 超时控制      | 所有调用必须支持 context timeout / cancel       |
| 模型 SDK    | 业务代码不得直接调用模型 SDK                        |
| 调用日志      | 必须记录 provider、model、latency、token、error |
| 重试        | 网络抖动、5xx、429 可重试；认证失败、参数错误不重试           |
| 熔断        | 每个 provider/model 独立熔断                  |
| fallback  | 主模型失败后可路由到备用模型                          |
| 流式输出      | 一旦开始输出 token，默认不自动切换 provider           |

推荐配置：
llm:
  pool:
    max_concurrent: 10
    acquire_timeout: 3s
  stream:
    max_concurrent: 8
    acquire_timeout: 3s
  timeout:
    sync_call: 60s
    stream_idle: 120s
    health_check: 10s
  retry:
    max_attempts: 2
    initial_interval: 500ms
    max_interval: 5s
  fallback:
    enabled: true
    max_attempts: 2

## SSE 与 Run 事件规范

SSE 事件必须走：

Worker → Redis Stream → lattice-api → Browser

禁止 Worker 直接向浏览器写数据。

事件类型：

message.delta
tool.call
tool.result
command.output
file.diff
test.output
approval.required
run.completed
run.failed
run.cancelled
ping

## 路由规范

所有后端接口统一挂在 /api 下。

示例：

GET  /api/health
GET  /api/v1/providers
POST /api/v1/agents
POST /api/v1/runs
GET  /api/v1/runs/{run_id}/events

规则：

| 项          | 规范                                    |
| ---------- | ------------------------------------- |
| 后端路由       | 必须统一 `/api` 前缀                        |
| Vite proxy | 只转发 `/api`，不做 rewrite                 |
| 前端 Axios   | `baseURL = "/api"`                    |
| 前端 health  | `get("/health")` 对应后端 `/api/health`   |
| 禁止         | 不要同时保留 `/health` 和 `/api/health` 两套路由 |

## 中间项规范

| 中间件      | 作用                        |
| -------- | ------------------------- |
| TraceID  | 为每个请求生成 trace_id          |
| Logger   | 记录请求路径、状态码、耗时、trace_id    |
| Recovery | 捕获 panic，返回 InternalError |
| CORS     | 支持前端开发环境跨域                |
| Auth 预留  | 当前可放行，但保留 token 解析入口      |


## 部署架构

用户浏览器
  │
  ▼
Ingress Nginx
  │
  ├──▶ lattice-frontend
  │
  ├──▶ lattice-api
  │       ├──▶ MySQL
  │       ├──▶ Redis
  │       └──▶ PostgreSQL + pgvector
  │
  └──▶ lattice-worker
          ├──▶ Redis
          ├──▶ MySQL
          ├──▶ PostgreSQL + pgvector
          └──▶ LLM Provider

## 后端基础组件

业务开发前 P0 基础组件

在实现 provider、agent、chat、run 业务前，必须先补齐：
| 优先级 | 组件                                 | 说明                                         |
| --- | ---------------------------------- | ------------------------------------------ |
| P0  | migrations/mysql                   | MySQL DDL                                  |
| P0  | migrations/postgres                | pgvector DDL                               |
| P0  | internal/app/bootstrap.go          | 初始化 config、logger、DB、Redis、runtime、modules |
| P0  | internal/app/router.go             | 统一 `/api` 路由                               |
| P0  | internal/app/modules.go            | 显式装配所有模块                                   |
| P0  | LLMExecutor                        | llm-pool / llm-stream                      |
| P0  | Redis Stream                       | Run 事件流                                    |
| P0  | Health Check 增强                    | 检查 API、MySQL、Redis、pgvector                |
| P0  | TraceID / Logger / Recovery / CORS | 接口基础中间件                                    |
| P0  | 结构化日志                              | 支撑后续排查                                     |

## 编码规范

| 编号 | 规范                                          |
| -- | ------------------------------------------- |
| 1  | 包名使用小写单词，不使用下划线和驼峰。                         |
| 2  | 文件名使用小写加下划线。                                |
| 3  | 导出类型、函数、常量使用 PascalCase；非导出标识符使用 camelCase。 |
| 4  | 缩略词统一大写，如 ID、URL、HTTP、JSON、SSE。             |
| 5  | 错误必须显式处理，不允许忽略 error。                       |
| 6  | 错误返回必须携带上下文并保留错误链。                          |
| 7  | 不使用 panic 处理业务错误。                           |
| 8  | Handler 只做参数绑定、校验和调用 application。           |
| 9  | application 负责用例编排、事务控制和跨模块调用。              |
| 10 | domain 不依赖外部框架和基础设施。                        |
| 11 | infra 负责数据库、Redis、pgvector、外部 SDK 访问。       |
| 12 | 日志必须结构化输出。                                  |
| 13 | 日志禁止输出密钥、密码、token。                          |
| 14 | 所有外部调用必须传 context.Context。                  |
| 15 | goroutine 必须可退出。                            |
| 16 | 禁止无界 goroutine。                             |
| 17 | channel 必须明确关闭责任。                           |
| 18 | 共享 map 并发读写必须加锁。                            |
| 19 | 新增第三方依赖前必须说明理由。                             |
| 20 | 不引入不必要的设计模式。                                |

AI 生成代码必须遵守
新功能必须先判断归属模块。
不确定归属时先询问，不得放入 common。
不要创建无意义的 manager/helper/processor。
不要创建万能 BaseRepository。
不要在 Handler 中写业务逻辑。
不要在 application 中直接写 SQL。
不要让 domain 依赖 Gin、GORM、Redis、Eino。
不要跨模块调用 infra。
不要跨模块传 Entity / PO。
所有 DB、Redis、LLM、Tool 调用必须传 context.Context。
所有错误必须显式处理。
新增外部依赖必须先说明原因。
生成后必须保证 go test ./... 可通过。
修改前后必须说明改了哪些文件。
如果涉及前后端接口，必须说明最终请求路径。
如果涉及 SSE，必须说明事件类型和关闭条件。
如果涉及 Redis key，必须说明 key 格式和 TTL。
如果涉及数据库，必须同步 migration。
如果涉及 Run 写任务，必须考虑项目写锁。
如果涉及 LLM，必须走 LLMExecutor，不得直接调用 SDK。