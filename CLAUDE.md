# Lattice-Coding 项目开发规范

## 项目概览

Lattice-Coding 是一个面向实验室内部使用的轻量级 AI Agent 编码平台，主要服务于 10-20 人规模的科研团队。

项目目标不是复刻 Claude Code、Codex 或 Dify 这类完整商业平台，而是构建一个可控、可验证、可复盘的内部 Coding Agent 系统，用于辅助团队完成代码理解、代码修改、测试修复、代码审查、实验脚本维护、项目文档检索等研发任务。

Lattice-Coding 的核心定位是：

> 面向科研团队的内网部署型 AI Agent 编码平台。

## 技术栈

### 后端

- Go
- Eino Agent 框架
- Gin 或 Hertz 作为 HTTP 框架
- MySQL 8.x
- Redis 7.x
- GORM 或 sqlc
- zap / slog 作为结构化日志组件

### 前端

- 轻量 Web 管理控制台
- 支持 SSE 流式输出展示
- 不做完整 Web IDE

### 模型接入

- OpenAI SDK / OpenAI-compatible API
- Anthropic Claude API
- DeepSeek / Qwen / Ollama 等通过统一模型适配层接入

### 部署

- Docker Compose
- 实验室内网服务器部署
- 后端服务、Worker、MySQL、Redis、Nginx 独立容器

---

---

## 核心功能模块

| 模块 | 说明 |
|------|------|
| 多模型提供商管理 | 支持 OpenAI、Claude、DeepSeek、Qwen、Ollama、本地模型和 OpenAI-compatible API。统一模型调用接口，支持模型切换、默认模型、超时配置和调用日志。 |
| Agent 创建与配置 | 支持创建代码修复、代码审查、测试生成、实验复现等 Agent。每个 Agent 可配置模型、系统提示词、工具权限、最大执行步数和默认测试命令。 |
| 对话引擎 | 支持多轮对话、SSE 流式输出、工具调用结果回填、任务中断与继续。对话接口需要稳定支持长连接。 |
| 工具调用系统 | 支持文件读取、代码搜索、文件编辑、Shell 命令执行、Git diff、测试执行、知识库搜索等工具。所有工具调用必须经过权限检查并写入审计日志。 |
| Coding Agent 执行闭环 | 支持“理解需求 → 读取上下文 → 制定计划 → 修改代码 → 执行测试 → 根据错误继续修复 → 输出报告”的完整流程。最终结果必须包含 diff、测试结果、风险说明和执行摘要。 |
| 权限与安全控制 | 控制 Agent 可访问的文件路径、可执行命令和可使用工具。高风险操作必须二次确认或直接禁止。 |
| Git Diff 与回滚 | 每次 Agent 修改代码都必须记录 patch。支持查看 diff、导出 patch、回滚本次修改。 |
| 日志与审计 | 记录用户任务、Agent 计划、模型调用、工具调用、命令执行、文件修改、测试结果和最终报告。所有关键行为必须可追踪、可复盘。 |
| 轻量知识库 | 第一阶段支持 Markdown / TXT / 本地文档检索。优先使用关键词检索、MySQL FULLTEXT 或 pgvector，避免复杂 RAG 流程。 |
| 轻量管理控制台 | 支持管理模型、Agent、项目、任务记录、运行日志和权限策略。控制台以实用为主，不追求复杂可视化。 |

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

## 代码组织规范

### 包结构

```text
lattice-coding
├── cmd/
│   ├── api/                         # API 服务入口
│   │   └── main.go
│   └── worker/                      # Worker 服务入口
│       └── main.go
│
├── internal/
│   ├── common/
│   │   ├── config/                  # 配置加载、MySQL、Redis、pgvector、Eino 初始化
│   │   ├── errors/                  # BizError、ErrorCode
│   │   ├── response/                # Result、PageResult
│   │   ├── logger/                  # 日志封装
│   │   ├── middleware/              # 鉴权、trace、recover、cors
│   │   └── util/                    # ID、时间、JSON 等无业务工具
│   │
│   ├── modules/
│   │   ├── provider/                # 模型提供商管理
│   │   ├── agent/                   # Agent 配置
│   │   ├── chat/                    # 对话引擎、SSE
│   │   ├── run/                     # Agent Run、任务状态机、队列
│   │   ├── mcp/                     # Tool / MCP 工具管理与调用
│   │   ├── workflow/                # 简版工作流编排
│   │   ├── knowledge/               # 知识库、RAG、pgvector 检索
│   │   ├── safety/                  # 权限、命令拦截、路径拦截
│   │   └── audit/                   # 模型调用、工具调用、命令执行审计
│   │
│   └── runtime/
│       ├── eino/                    # Eino Runtime 封装
│       ├── llm/                     # OpenAI、Claude、Qwen、Ollama 适配
│       ├── tool/                    # 内置 Tool 执行器
│       └── event/                   # Redis Stream 事件封装
│
├── deployments/
│   ├── docker/
│   └── k8s/
│
├── migrations/
│   ├── mysql/
│   └── postgres/
│
├── configs/
│   └── config.yaml
│
└── go.mod

每个模块内部四层结构：

modules/{module}
├── api/                     # Handler、Request、Response、Router
├── application/             # Service、Command、Query、DTO、Assembler
├── domain/                  # Entity、ValueObject、Repository 接口、领域规则
├── infra/                   # Mapper、PO、Repository 实现、外部适配
└── module.go                # 模块初始化、依赖装配、路由注册

### 各层职责边界

| 层 | 职责 | 禁止 |
|----|------|------|
| api | 处理 HTTP 请求、参数校验、响应转换；注册路由；调用本模块 application。 | 不写业务逻辑；不访问 MySQL、Redis、pgvector；不直接调用 Eino、Tool、Mapper；不调用其他模块 infra。 |
| application | 处理用例编排、事务控制、跨模块调用；调用本模块 domain 和 repository 接口；跨模块调用其他模块 application service。 | 不直接写 SQL；不直接调用其他模块 infra；不返回 PO / Entity 给前端；不依赖 HTTP 框架对象。 |
| domain | 定义核心实体、值对象、领域规则、Repository 接口；维护业务不变量。 | 不依赖 Gin、MySQL、Redis、pgvector、Eino；不出现 SQL；不调用其他模块；不处理 HTTP、配置、日志、外部 SDK。 |
| infra | 实现 domain.Repository；负责 Mapper、PO、SQL、数据库转换；访问 MySQL、Redis、pgvector、外部 SDK。 | 不写业务流程；不被其他模块直接调用；不向外暴露 PO；不跨模块访问其他模块数据库表。 |
| runtime | 封装 Eino、LLM、Tool、Redis Stream 等运行时能力；为 run、chat、knowledge 等模块提供执行能力。 | 不直接处理业务规则；不直接处理 HTTP；不绕过 application 层被 api 调用；不直接决定权限策略。 |

### 跨模块调用

禁止
任意模块 api          → 其他模块 api
任意模块 api          → 其他模块 infra
任意模块 application  → 其他模块 infra
任意模块 infra        → 其他模块 infra
任意模块              → 其他模块 mapper
任意模块              → 其他模块 PO / Entity
任意模块              → 直接查询其他模块数据库表

跨模块数据传输只能使用DTO
AgentDTO
ProviderDTO
ModelDTO
RunDTO
ToolDTO
KnowledgeDocDTO

## LLM调用规范

### 基本原则

| 项 | 规范 |
|----|------|
| 调用隔离 | LLM 调用必须和普通 HTTP 请求隔离，不允许在 API Handler 中直接阻塞执行模型请求。 |
| 并发控制 | 所有 LLM 调用必须进入专用并发池：`llm-pool` 或 `llm-stream`。 |
| 超时控制 | 所有 LLM 调用必须传递 `context.Context`，必须支持超时、取消和中断。 |
| 日志审计 | 每次调用必须记录 provider、model、run_id、耗时、状态、错误、token 用量。 |
| 配置化 | provider、model、timeout、retry、breaker、fallback 必须配置化，不允许硬编码。 |

### 并发池配置

使用 goroutine + channel / semaphore 控制并发

```yaml
llm:
  pool:
    max_concurrent: 10
    acquire_timeout: 3s

  stream:
    max_concurrent: 8
    acquire_timeout: 3s
```

并发池实现案例
```go
package llm

import (
	"context"
	"errors"
	"time"
)

var ErrPoolFull = errors.New("llm pool is full")

type SemaphorePool struct {
	sem chan struct{}
}

func NewSemaphorePool(maxConcurrent int) *SemaphorePool {
	return &SemaphorePool{
		sem: make(chan struct{}, maxConcurrent),
	}
}

func (p *SemaphorePool) Acquire(ctx context.Context, acquireTimeout time.Duration) error {
	timer := time.NewTimer(acquireTimeout)
	defer timer.Stop()

	select {
	case p.sem <- struct{}{}:
		return nil
	case <-timer.C:
		return ErrPoolFull
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *SemaphorePool) Release() {
	select {
	case <-p.sem:
	default:
	}
}
```

LLM Pool定义
```go
package llm

import "time"

type PoolConfig struct {
	MaxConcurrent int
	AcquireTimeout time.Duration
}

type LLMExecutor struct {
	Pool       *SemaphorePool // 非流式调用
	StreamPool *SemaphorePool // 流式调用

	PoolConfig       PoolConfig
	StreamPoolConfig PoolConfig
}

func NewLLMExecutor(poolCfg, streamCfg PoolConfig) *LLMExecutor {
	return &LLMExecutor{
		Pool:       NewSemaphorePool(poolCfg.MaxConcurrent),
		StreamPool: NewSemaphorePool(streamCfg.MaxConcurrent),

		PoolConfig:       poolCfg,
		StreamPoolConfig: streamCfg,
	}
}
```

```go
\\llm-pool:非流式调用(用于阻塞等待完整响应的模型调用)
package llm

import (
	"context"
	"time"
)

type ChatRequest struct {
	Provider string
	Model    string
	Messages []Message
}

type ChatResponse struct {
	Content string
}

type Message struct {
	Role    string
	Content string
}

type ChatClient interface {
	Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)
}

func (e *LLMExecutor) Chat(
	ctx context.Context,
	client ChatClient,
	req ChatRequest,
) (*ChatResponse, error) {
	if err := e.Pool.Acquire(ctx, e.PoolConfig.AcquireTimeout); err != nil {
		return nil, err
	}
	defer e.Pool.Release()

	callCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	return client.Chat(callCtx, req)
}
```
```go
\\llm-stream:流式SSE调用(长连接)
package llm

import (
	"context"
	"time"
)

type StreamChunk struct {
	Content string
	Done    bool
	Err     error
}

type StreamClient interface {
	Stream(ctx context.Context, req ChatRequest) (<-chan StreamChunk, error)
}

func (e *LLMExecutor) Stream(
	ctx context.Context,
	client StreamClient,
	req ChatRequest,
) (<-chan StreamChunk, error) {
	if err := e.StreamPool.Acquire(ctx, e.StreamPoolConfig.AcquireTimeout); err != nil {
		return nil, err
	}

	out := make(chan StreamChunk, 32)

	go func() {
		defer e.StreamPool.Release()
		defer close(out)

		streamCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		ch, err := client.Stream(streamCtx, req)
		if err != nil {
			out <- StreamChunk{Err: err}
			return
		}

		idleTimer := time.NewTimer(120 * time.Second)
		defer idleTimer.Stop()

		for {
			select {
			case chunk, ok := <-ch:
				if !ok {
					out <- StreamChunk{Done: true}
					return
				}

				if !idleTimer.Stop() {
					select {
					case <-idleTimer.C:
					default:
					}
				}
				idleTimer.Reset(120 * time.Second)

				select {
				case out <- chunk:
				case <-ctx.Done():
					out <- StreamChunk{Err: ctx.Err()}
					return
				}

			case <-idleTimer.C:
				out <- StreamChunk{Err: context.DeadlineExceeded}
				return

			case <-ctx.Done():
				out <- StreamChunk{Err: ctx.Err()}
				return
			}
		}
	}()

	return out, nil
}
```

### Go HTTP Client 配置

取代Java的OkHttpClient，Go 中由 `http.Client` + `http.Transport` 承担。

| 配置项 | 建议值 |
|----|----|
| MaxIdleConns | 100 |
| MaxIdleConnsPerHost | 20 |
| MaxConnsPerHost | 50 |
| IdleConnTimeout | 90s |
| TLSHandshakeTimeout | 10s |
| ResponseHeaderTimeout | 60s |
| ExpectContinueTimeout | 1s |
| DialTimeout | 10s |
| KeepAlive | 30s |

示例配置：

```go
transport := &http.Transport{
    MaxIdleConns:        100,
    MaxIdleConnsPerHost: 20,
    MaxConnsPerHost:     50,
    IdleConnTimeout:     90 * time.Second,
    TLSHandshakeTimeout: 10 * time.Second,
    ResponseHeaderTimeout: 60 * time.Second,
    ExpectContinueTimeout: 1 * time.Second,
    DialContext: (&net.Dialer{
        Timeout:   10 * time.Second,
        KeepAlive: 30 * time.Second,
    }).DialContext,
}

client := &http.Client{
    Transport: transport,
    Timeout:   0, // 不在 Client 层设置总超时，统一由 context 控制
}

### 超时层次

| 层级             | 用途             | 建议值   |
| -------------- | -------------- | ----- |
| API 请求超时       | 普通 HTTP 请求     | 30s   |
| 同步 LLM 超时      | 非流式模型调用        | 60s   |
| SSE 空闲超时       | 流式无输出超时        | 120s  |
| Agent Run 总超时  | 一次 Agent 任务总时长 | 20min |
| Provider 连通性测试 | health check   | 10s   |
| Embedding 调用   | 文档向量化          | 30s   |
| Rerank 调用      | 检索结果重排         | 30s   |

### 重试策略

| 异常类型    | 是否重试 | 策略                  |
| ------- | ---- | ------------------- |
| 网络抖动    | 是    | 最多 2 次，指数退避         |
| 连接超时    | 是    | 最多 2 次，指数退避         |
| 5xx     | 是    | 最多 2 次，指数退避         |
| 429 限流  | 是    | 指数退避，尊重 Retry-After |
| 认证失败    | 否    | 直接失败                |
| 参数错误    | 否    | 直接失败                |
| 模型不存在   | 否    | 直接失败                |
| 上下文超限   | 否    | 直接失败，交给上层裁剪上下文      |
| 用户主动取消  | 否    | 立即中断                |
| 流式输出已开始 | 谨慎   | 默认不重试，避免重复输出        |

### 熔断配置
```yaml
llm:
  circuit_breaker:
    enabled: true
    window: 60s
    min_requests: 10
    failure_rate_threshold: 0.5
    slow_call_threshold: 30s
    slow_call_rate_threshold: 0.6
    open_duration: 30s
    half_open_max_requests: 3
```
### Fallback 路由

```yaml
llm:
  routing:
    default:
      primary: claude-sonnet
      fallback:
        - gpt-4.1
        - deepseek-chat
        - qwen-plus
        - ollama-qwen-coder
```

## 部署架构

## 部署架构

用户浏览器  
  │  
  ▼  
Ingress Nginx（L7 负载均衡 + SSL 终止 + SSE 支持）  
  │  
  ├──▶ lattice-frontend（Vue / React SPA，Nginx 静态文件服务，2 副本）  
  │  
  ├──▶ lattice-api（Go API 服务，HTTP 接口 + SSE 输出，2 副本）  
  │       │  
  │       ├──▶ MySQL 8.x（主业务数据存储）  
  │       ├──▶ Redis 7.x（Run 状态 / 队列 / SSE Stream / 缓存 / 限流）  
  │       └──▶ PostgreSQL + pgvector（知识库向量存储）  
  │  
  └──▶ lattice-worker（Go Worker，Agent Run / Tool 调用 / Shell 执行，1-2 副本）  
          │  
          ├──▶ Redis 7.x（消费 Run 队列，写入 SSE Stream）  
          ├──▶ MySQL 8.x（写入 Run 结果与审计日志）  
          ├──▶ PostgreSQL + pgvector（知识库检索与索引）  
          └──▶ LLM Provider（OpenAI / Claude / DeepSeek / Qwen / Ollama）


**Ingress 关键配置（SSE 必须）**:
```yaml
nginx.ingress.kubernetes.io/proxy-read-timeout: "300"
nginx.ingress.kubernetes.io/proxy-send-timeout: "300"
nginx.ingress.kubernetes.io/proxy-buffering: "off"
nginx.ingress.kubernetes.io/proxy-request-buffering: "off"
nginx.ingress.kubernetes.io/limit-rps: "20"
```
---

## 编码规范（基于字节跳动Go开发手册）

| 编号 | 规范 |
|------|------|
| 1 | 包名使用小写单词，不使用下划线和驼峰；包名应短小、明确，例如 `agent`、`run`、`safety`。 |
| 2 | 文件名使用小写加下划线，例如 `agent_service.go`、`run_handler.go`。 |
| 3 | 导出类型、函数、常量使用 PascalCase；非导出标识符使用 camelCase。 |
| 4 | 缩略词保持统一大写，例如 `ID`、`URL`、`HTTP`、`JSON`、`SSE`，不要写成 `Id`、`Url`、`Http`。 |
| 5 | 布尔变量使用肯定语义，例如 `isEnabled`、`hasPermission`、`canExecute`，避免 `isNotValid` 这类否定命名。 |
| 6 | 错误必须显式处理，不允许使用 `_` 忽略 `error`。 |
| 7 | 错误返回必须携带上下文，使用 `fmt.Errorf("xxx failed: %w", err)` 保留错误链。 |
| 8 | 不使用 `panic` 处理业务错误；`panic` 只允许用于启动阶段不可恢复错误或明确的程序员错误。 |
| 9 | Handler 层只做参数绑定、校验、调用 application，不写业务逻辑。 |
| 10 | application 层负责用例编排、事务控制、跨模块调用，不直接写 SQL。 |
| 11 | domain 层不依赖 Gin、MySQL、Redis、pgvector、Eino，不出现 SQL 和外部 SDK。 |
| 12 | infra 层负责数据库、Redis、pgvector、外部 SDK 访问，不向 api/application 暴露 PO。 |
| 13 | 日志必须结构化输出，至少包含 `trace_id`、`run_id`、`session_id`、`user_id`、`error`。 |
| 14 | 日志中禁止输出 API Key、Token、密码、Cookie、完整 `.env` 内容。 |
| 15 | 外部调用、数据库访问、Redis 操作、Shell 执行、LLM 调用必须传递 `context.Context`。 |
| 16 | goroutine 必须可退出；启动 goroutine 时必须考虑 `context.Done()`、超时或关闭信号。 |
| 17 | 禁止无界 goroutine；LLM、Shell、Embedding、SSE 等长耗时任务必须通过 semaphore / worker pool 控制并发。 |
| 18 | channel 必须明确关闭责任；发送方负责关闭，接收方不得随意关闭 channel。 |
| 19 | 共享 map 不允许并发读写；并发场景必须使用 `sync.Mutex`、`sync.RWMutex`、`sync.Map` 或单 goroutine 管理。 |
| 20 | 不引入不必要的设计模式和第三方依赖；能用标准库解决的优先使用标准库，新增依赖前必须说明理由。 |
| 21 | Go 项目更适合“公共能力轻量封装 + 各模块显式 Repository 模式“ ，不要想着替代JAVA中的某些能力，代码需要尊重Go项目风格 |

新增第三方库前必须说明：

解决什么问题
为什么标准库或现有依赖不能满足
是否会增加部署、维护或安全成本
是否与 Go + Eino + MySQL 8.x + Redis 7.x 的技术栈边界一致


## 数据库规范

### MySQL 通用字段约定每张表必须包含以下字段：
```sql
id BIGINT NOT NULL AUTO_INCREMENT, -- 主键，禁用 UUID
created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
deleted TINYINT(1) NOT NULL DEFAULT 0, -- 逻辑删除标志PRIMARY KEY (id)
```
- 字符集：`utf8mb4`，排序规则：`utf8mb4_unicode_ci`- 禁用 `VARCHAR` 无长度约束，text 类 content 字段用 `MEDIUMTEXT`- 金额用 `DECIMAL(19,4)`，禁用 `FLOAT/DOUBLE`- 布尔用 `TINYINT(1)`，不用 `BIT`

### 索引设计原则

1. **区分度低的字段不单独建索引**（如 deleted、status 枚举），必须与高区分度字段组合

2. **组合索引遵循最左前缀**：等值查询字段在左，范围查询字段在右

3. **查询条件中含 `deleted`**，必须将 `deleted` 纳入索引

4. **每表索引不超过 5 个**（含主键），写多读少的表控制在 3 个以内

5. **禁止在 `TEXT/BLOB` 类型字段上建普通索引**，需要时建前缀索引或全文索引```sql-- 正确示例：conversation_id 高区分度在左，deleted 次之，created_at 范围在右INDEX idx_conv_created (conversation_id, deleted, created_at)```

### 大表处理策略

判断为大表的阈值：行数 > 500 万 或 数据量 > 2GB
| 场景 | 策略 |
|------|------|
| t_message | 按 conversation_id 分区，或按月归档冷数据 |
| 知识库向量表 | ivfflat 索引，lists = sqrt(行数) |
| 日志类表 | 只保留 90 天，定期 DELETE + OPTIMIZE TABLE |

### 分页查询规范- 
**禁止** `LIMIT offset, size` 深分页（offset > 1000 全表扫描）
- 对话记录类使用**游标分页**：
```sql
SELECT id, role, content, created_at FROM t_message
WHERE conversation_id = ? 
    AND deleted = 0 
    AND (created_at < ? OR (created_at = ? AND id < ?))
ORDER BY created_at DESC, id DESC
LIMIT 20;
```
- 管理后台必须分页时，用 `WHERE id > lastId LIMIT size` 替代 offset

### pgvector 索引规范
```sql
-- 余弦相似度索引，lists 值 = sqrt(总行数)，行数 <10 万时 lists=100
CREATE INDEX idx_embedding_ivfflat ON knowledge_embedding
USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);
-- 查询时设置 probes，精度和速度平衡
SET ivfflat.probes = 10;
SELECT * FROM knowledge_embeddingORDER BY embedding <=> '[...]'::vector LIMIT 5;
```

### 索引检测措施

**开发阶段**：启用 p6spy，拦截执行 >10ms 的查询自动 EXPLAIN，type=ALL 时打印警告日志。

**CI 阶段**：关键查询写 `IndexCoverageTest`，EXPLAIN 结果中 type=ALL 则测试失败，阻断合并。

**生产阶段**：定期查询 `performance_schema.events_statements_summary_by_digest`，找出 `sum_no_index_used > 0` 的 SQL。

```sql
SELECT digest_text, count_star AS 执行次数, sum_no_index_used AS 未用索引次数
FROM performance_schema.events_statements_summary_by_digest
WHERE sum_no_index_used > 0
ORDER BY sum_no_index_used DESC LIMIT 20;
```
---

## 性能瓶颈优先级

| 级别 | 瓶颈 | 一期处理方式 |
|------|------|-------------|
| P0 | LLM API 延迟高（3-30s） | 线程隔离 + 熔断 + Fallback（已设计） |
| P0 | 向量检索无索引全表扫描 | 建 ivfflat 索引（建表时必须创建） |


## 当前版本定位

目前 Lattice-Coding 先实现的是：

Web 前端 + 服务端 API + 服务端 Worker + 服务器 workspace

也就是说，用户通过网页管理 Agent、模型和任务，后端 Worker 在服务器上预配置的 workspace 中读取代码、执行测试、生成报告和 diff。

这个阶段的目标不是直接替代 Cursor 或 Claude Code，而是先做一个：

服务器侧 AI 代码审计 / 测试 / Patch 生成平台

它适合统一管理实验室服务器上的共享项目、实验代码和测试任务。

### 需要注意的问题

| 问题           | 当前风险                               | 需要注意的设计                                                                     |
| ------------ | ---------------------------------- | --------------------------------------------------------------------------- |
| 项目权限         | 内网用户可能通过网页操作不属于自己的项目。              | 必须做项目级权限控制，区分查看、执行、写入、审批权限。                                                 |
| workspace 安全 | Worker 如果直接操作服务器目录，可能越权访问敏感文件。     | 项目路径必须限制在白名单 workspace 下，禁止访问 `.env`、`.ssh`、`/etc`、`/root`、Docker socket 等。 |
| 共享目录污染       | Agent 直接改共享项目，可能影响其他用户。            | 不直接改主目录，优先使用独立 worktree / 临时副本，最终生成 patch。                                  |
| 多人并发冲突       | 多个用户同时让 Agent 修改同一个项目。             | 同一项目写任务必须加锁，同一时间只允许一个写任务。                                                   |
| 人工 CR 体验弱    | 用户只能在网页看报告和 diff，不能结合本地 IDE 上下文审查。 | 当前结果定位为 AI 初审和建议，不作为最终合入依据。                                                 |
| 服务器环境差异      | Worker 跑的是服务器环境，不一定等于用户本地开发环境。     | Run 记录必须保存 commit、branch、命令、环境摘要、模型版本和执行日志。                                 |
| Shell 执行风险   | Agent 可能被诱导执行危险命令。                 | 命令白名单、路径白名单、人工审批、超时和审计必须存在。                                                 |
| 结果可信度        | Agent 可能误报、漏报或幻觉。                  | 最终报告必须包含依据、修改文件、diff、测试结果和风险说明。                                             |
| 日志膨胀         | 测试日志、命令输出、SSE 事件可能很大。              | Redis 只放运行态，MySQL 只存摘要，大日志写文件并设置清理策略。                                       |

### 未来目标

服务端 + Web 管理页面 + 每个用户本地客户端

因此现在开发服务器 workspace 版本时，最重要的是不要把“文件系统操作”写死在业务逻辑里。

应该提前抽象：

| 能力       | 当前实现                                   | 未来实现                            |
| -------- | -------------------------------------- | ------------------------------- |
| 读文件      | ServerWorkspaceExecutor 读服务器 workspace | LocalClientExecutor 请求用户本地客户端读取 |
| 写文件      | 服务器临时 workspace 写 patch                | 用户客户端在本地生成或应用 patch             |
| 搜索代码     | 服务器执行 grep / rg                        | 用户客户端在本地执行搜索                    |
| 运行测试     | 服务器 Worker 执行 shell                    | 用户客户端在本地 shell 环境执行             |
| Git diff | 服务器 workspace 执行 git diff              | 用户客户端在本地仓库执行 git diff           |
| 工具调用事件   | Worker 写 Redis Stream                  | 客户端回传事件，Server 写 Redis Stream   |

建议抽象为：`WorkspaceExecutor`
业务层只依赖这个接口，不直接依赖 os.ReadFile、os.WriteFile、exec.Command。

不得删除既有模块目录结构。除非我明确要求，否则不要删除 api、application、domain、infra、module.go 中任何一层。