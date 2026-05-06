# Skill: Lattice-Coding 业务模块全流程交付

## 触发方式

当用户提出以下需求时，按本 Skill 推进：

- “开发 XX 模块”
- “实现 XX 功能”
- “交付 XX”
- “给 Agent 写 XX 模块提示词”
- “按照项目规范生成 XX 后端/前端”
- “把 mock 换成真实 API”
- “完成 XX 模块的全链路验收”

本 Skill 适用于 Lattice-Coding 的所有业务模块：

- provider：模型供应商管理
- agent：Agent 配置管理
- chat：会话与消息
- run：Agent Run 任务执行
- mcp / tool：工具管理与调用
- workflow：简版工作流
- knowledge：知识库与检索
- safety：权限与安全控制
- audit：审计日志

---

## 项目固定上下文

Lattice-Coding 是一个 AI Agent 编码平台，当前阶段采用：

| 层 | 技术 |
|---|---|
| 后端 | Go + Gin + GORM + Eino |
| 数据库 | MySQL |
| 缓存 / 运行态 / 事件流 | Redis |
| 向量检索 | PostgreSQL + pgvector，按需接入 |
| 前端 | Vue 3 + TypeScript + Vite + Element Plus |
| 架构 | 模块化单体 |
| 部署目标 | 实验室服务器 / Docker / 后续 K8s |
| 当前执行模式 | Web + Server + Server Workspace |
| 后续演进 | Server + Web 管理页 + 每用户本地 Client |

后端模块标准结构：

```text
internal/modules/{module}/
├── api/
│   ├── handler.go
│   ├── request.go
│   ├── response.go
│   └── router.go
│
├── application/
│   ├── command.go
│   ├── query.go
│   ├── dto.go
│   ├── assembler.go
│   ├── command_service.go
│   └── query_service.go
│
├── domain/
│   ├── entity.go
│   ├── value_object.go
│   └── repository.go
│
├── infra/
│   └── persistence/
│       ├── model.go
│       ├── converter.go
│       └── repository_impl.go
│
└── module.go
```

---

# 总体原则

## 必须遵守

1. 先咨询，再实现。
2. 先数据模型，再 schema。
3. 先后端，再前端。
4. 先 CRUD，再运行时能力。
5. 每一步必须有明确产出物。
6. 每一步必须有验证方式。
7. 每一步验证通过后，才能进入下一步。
8. 关键设计决策必须等待用户确认。
9. 不得删除既有模块目录结构。
10. 不得通过删除文件解决 import cycle。
11. 不得把业务逻辑塞进 common。
12. 不得让 api 层访问 DB / Redis / Eino。
13. 不得让 application 层依赖 api 或 infra。
14. 不得让 domain 层依赖 Gin / GORM / Redis / Eino。
15. 不得让 infra 层反向依赖 application / api。
16. 跨模块调用只允许通过 application service 或显式接口。
17. 所有 DB / Redis / LLM / Tool / Shell 调用必须传 `context.Context`。
18. 涉及敏感信息必须脱敏。
19. 涉及文件写入、Shell、外部模型、外部网络时必须先确认安全边界。
20. 每次生成代码后必须说明修改了哪些文件、每个文件的职责、如何验证。

## 禁止行为

- 未确认数据模型就生成 Service。
- 未更新 schema 就写 Repository。
- 未后端 curl 验证就写前端页面。
- 前端继续使用 mock 却声称完成模块。
- Handler 中写业务逻辑。
- application 中直接写 SQL。
- domain 中出现 GORM tag。
- infra 中返回 PO 给上层。
- 业务模块直接 import 具体 LLM SDK。
- 业务模块直接 `os.ReadFile` / `os.WriteFile` / `exec.Command`。
- 自动新增第三方依赖。
- 自动修改全局架构、路由前缀、响应结构。
- 自动删除 `api/`、`application/`、`domain/`、`infra/`、`module.go`。

---

# 依赖方向规范

标准依赖方向：

```text
api → application → domain
infra/persistence → domain
module.go → api / application / infra
app/bootstrap → modules
runtime → Eino / LLM / Tool / Event
```

禁止依赖方向：

```text
application → api
application → infra
domain → api
domain → application
domain → infra
infra → application
infra → api
api → infra
api → runtime/llm
api → runtime/eino
```

运行时边界：

```text
provider / agent / chat / run 等业务模块
→ application 编排
→ runtime/llm、runtime/eino、runtime/tool、runtime/event
```

---

# Step 1 — 咨询与设计，不写代码

## 目标

对齐模块边界、业务范围、数据模型、接口范围、运行时能力、安全风险，避免返工。

## 必须输出

1. 模块职责说明。
2. 模块不负责什么。
3. 依赖哪些模块。
4. 被哪些模块依赖。
5. 需要哪些表。
6. 核心字段草案。
7. 状态枚举草案。
8. 接口清单草案。
9. 是否需要 Redis。
10. 是否需要 pgvector。
11. 是否需要 LLM / Eino。
12. 是否需要 Tool / Shell / 文件系统。
13. 是否需要审计。
14. 是否需要权限审批。
15. 候选方案对比表。
16. 推荐方案和理由。
17. 等待用户确认的问题列表。

## 候选方案表格式

| 方案 | 做法 | 优点 | 缺点 | 是否推荐 |
|---|---|---|---|---|

## 等待用户确认

以下情况必须暂停：

- 新增数据库表。
- 新增核心字段。
- 删除字段。
- 改变状态枚举。
- 使用 JSON 字段还是独立表。
- 是否允许软删除。
- 是否允许物理删除。
- 是否接入真实 LLM。
- 是否调用外部网络。
- 是否执行 Shell。
- 是否写文件。
- 是否引入 Redis Stream。
- 是否引入 pgvector。
- 是否新增第三方依赖。
- 是否新增跨模块调用。
- 是否新增定时任务。
- 是否修改统一响应结构。
- 是否修改 `/api` 路由规范。

## 验证方式

用户明确回复类似：

```text
模块边界、数据模型和接口范围确认，可以进入 Step 2。
```

---

# Step 2 — 更新 schema / migration

## 目标

数据库 DDL 与设计对齐。所有后端 Repository、DTO、Service 都必须基于已确认的 schema。

## 产出物

根据项目当前方式生成或修改：

```text
migrations/mysql/{序号}_{module}.sql
```

或当前项目使用的：

```text
schema.sql
```

如涉及 pgvector：

```text
migrations/postgres/{序号}_{module}.sql
```

## 表设计要求

| 项 | 规范 |
|---|---|
| 主键 | `id BIGINT NOT NULL AUTO_INCREMENT` |
| 时间字段 | `created_at DATETIME(3)`、`updated_at DATETIME(3)` |
| 逻辑删除 | 当前项目可使用 `deleted TINYINT` 或统一 BasePO 的软删除字段，必须与已有代码一致 |
| 唯一约束 | 业务唯一字段需加唯一索引，通常包含 deleted |
| 查询索引 | 高频查询字段、状态、类型、父 ID、创建时间必须加索引 |
| JSON 字段 | 只放扩展配置，不放高频查询字段 |
| 敏感字段 | 只存密文，不存明文 |
| 外键 | 默认应用层维护，不强制数据库外键 |
| 状态字段 | 必须有明确枚举含义 |

## 验证方式

执行 SQL：

```bash
mysql -u <user> -p <db> < migrations/mysql/{file}.sql
```

检查表：

```bash
mysql -u <user> -p <db> -e "SHOW TABLES;"
mysql -u <user> -p <db> -e "DESC <table_name>;"
```

检查索引：

```bash
mysql -u <user> -p <db> -e "SHOW INDEX FROM <table_name>;"
```

## 注意事项

- 不依赖 GORM AutoMigrate 作为正式建表方式。
- schema 未确认前，不进入后端代码。
- JSON 字段不要替代核心关系。
- 密钥、token、密码字段必须考虑加密和脱敏。
- 高频日志表不要设计过多无用字段。
- 大文本日志不要默认塞进 MySQL，优先存摘要和文件路径。
- 需要关联校验的字段要提前设计索引。

---

# Step 3 — domain：实体、值对象、Repository 接口

## 目标

定义业务实体和值对象，明确模块需要的数据访问能力。

## 产出物

```text
internal/modules/{module}/domain/entity.go
internal/modules/{module}/domain/value_object.go
internal/modules/{module}/domain/repository.go
```

### entity.go

定义核心业务实体，例如：

```text
Provider
ModelConfig
Agent
ChatSession
ChatMessage
Run
Tool
KnowledgeDocument
SafetyCheck
AuditLog
```

### value_object.go

定义枚举和值对象，例如：

```text
Status
Type
Mode
Role
Permission
HealthStatus
RunStatus
ToolType
AuthType
ModelType
```

### repository.go

定义本模块的数据访问接口，例如：

```text
Create
Update
FindByID
FindPage
DeleteByID
ExistsByName
UpdateStatus
FindByParentID
```

## 验证方式

```bash
go test ./...
go list ./...
```

检查 domain 层依赖：

```bash
grep -R "gorm.io" internal/modules/{module}/domain || true
grep -R "github.com/gin-gonic/gin" internal/modules/{module}/domain || true
grep -R "github.com/redis" internal/modules/{module}/domain || true
grep -R "github.com/cloudwego/eino" internal/modules/{module}/domain || true
```

以上命令不应输出有效依赖。

## 注意事项

- Entity 不等于 PO。
- Entity 不需要 GORM tag。
- Entity 默认不加 json/form/binding tag。
- domain 不处理 HTTP。
- domain 不写 SQL。
- domain 不记录日志。
- domain 不读配置。
- Repository 接口定义“业务需要什么数据能力”，不暴露 GORM。
- 复杂领域规则可以放 `domain/service.go`，但不要过度抽象。

---

# Step 4 — infra/persistence：PO、Converter、RepositoryImpl

## 目标

实现 MySQL 持久化访问，完成 PO 与 Entity 的转换。

## 产出物

```text
internal/modules/{module}/infra/persistence/model.go
internal/modules/{module}/infra/persistence/converter.go
internal/modules/{module}/infra/persistence/repository_impl.go
```

### model.go

定义 PO：

```text
{Module}PO
```

要求：

- 嵌入项目统一 `db.BasePO` 或符合现有字段规范。
- 包含 GORM tag。
- 定义 `TableName()`。
- JSON 字段使用 GORM 支持的 JSON 类型或统一封装。
- 不包含业务方法。

### converter.go

只负责：

```text
PO -> Entity
Entity -> PO
[]PO -> []Entity
[]Entity -> []PO
```

### repository_impl.go

实现 domain.Repository：

```text
Create
Update
FindByID
FindPage
DeleteByID
Exists...
UpdateStatus...
```

## 验证方式

```bash
go test ./...
go list ./...
```

检查不允许依赖：

```bash
grep -R "internal/modules/{module}/api" internal/modules/{module}/infra || true
grep -R "internal/modules/{module}/application" internal/modules/{module}/infra || true
```

以上命令不应输出有效依赖。

## 注意事项

- infra 可以 import GORM。
- infra 可以 import domain。
- infra 不允许 import application。
- infra 不允许 import api。
- Repository 方法必须接收 `context.Context`。
- GORM 必须使用 `db.WithContext(ctx)`。
- 删除默认逻辑删除。
- 查询默认过滤逻辑删除。
- 排序字段必须白名单，避免 SQL 注入。
- infra 不写业务流程。
- infra 不返回 PO 给 application。

---

# Step 5 — DTO / Request / Response / Command / Query

## 目标

隔离 HTTP、业务输入输出、领域实体和数据库对象，避免层间对象混用。

## 产出物

```text
internal/modules/{module}/api/request.go
internal/modules/{module}/api/response.go
internal/modules/{module}/application/command.go
internal/modules/{module}/application/query.go
internal/modules/{module}/application/dto.go
internal/modules/{module}/application/assembler.go
```

## 职责划分

| 文件 | 职责 |
|---|---|
| api/request.go | HTTP Request；Request -> Command / Query |
| api/response.go | HTTP Response；DTO -> Response |
| application/command.go | 写操作输入 |
| application/query.go | 查询输入 |
| application/dto.go | application 输出和跨模块传输 |
| application/assembler.go | Entity -> DTO；Command -> Entity |
| infra/persistence/converter.go | PO <-> Entity |

## 产出对象

### Request

```text
Create{Module}Request
Update{Module}Request
{Module}PageRequest
```

### Command

```text
Create{Module}Command
Update{Module}Command
ChangeStatusCommand
```

### Query

```text
{Module}PageQuery
{Module}DetailQuery
```

### DTO

```text
{Module}DTO
{Module}DetailDTO
{Module}RuntimeDTO
```

### Response

```text
{Module}Response
{Module}DetailResponse
{Module}PageResponse
```

## 验证方式

```bash
go test ./...
go list ./...
```

## 注意事项

- application/assembler.go 不能 import api。
- api 可以 import application。
- Request 不传入 Service，先转 Command / Query。
- Response 不直接从 PO 生成。
- DTO 不含密文。
- Response 必须脱敏。
- 分页响应使用项目统一 PageResult。
- 不要让 PageResult 继承 Result。
- 如果出现 import cycle，优先检查转换函数是否放错层。

---

# Step 6 — application：业务 Service

## 目标

实现模块核心业务用例。

## 产出物

```text
internal/modules/{module}/application/command_service.go
internal/modules/{module}/application/query_service.go
```

## CommandService

负责写操作：

```text
Create
Update
Delete
Enable
Disable
ChangeStatus
Bind
Unbind
Approve
Reject
SetDefault
Sync
Test
```

## QueryService

负责读操作：

```text
GetByID
ListPage
ListByParentID
GetForRuntime
GetForReference
CheckReference
```

## 必须处理的业务规则

1. 创建唯一性校验。
2. 更新存在性校验。
3. 删除关联校验。
4. 状态流转校验。
5. 启用 / 禁用规则。
6. 默认项唯一性。
7. 敏感字段加密与脱敏。
8. 事务边界。
9. 缓存失效。
10. 审计记录。
11. 跨模块引用检查。
12. 运行时能力调用前的安全边界。

## 验证方式

```bash
go test ./...
```

如 Handler 已完成，再做 curl 验证。

## 注意事项

- application 不依赖 api。
- application 不依赖 infra。
- application 不直接写 SQL。
- application 不直接调用具体 SDK。
- application 不使用 Gin Context。
- 事务放在 application 层。
- 跨模块只调用其他模块 application service 或接口。
- 删除必须做关联校验。
- 返回 DTO，不返回 Entity / PO。
- 复杂运行时能力不要混在 CRUD 里，单独 Service 或 runtime 层封装。

---

# Step 7 — api：Handler / Router

## 目标

暴露 REST API，只做 HTTP 参数和响应。

## 产出物

```text
internal/modules/{module}/api/handler.go
internal/modules/{module}/api/router.go
```

## Handler 职责

只允许：

- 参数绑定
- 参数校验
- Request -> Command / Query
- 调用 application service
- DTO -> Response
- 统一响应返回

禁止：

- 写业务逻辑
- 访问 GORM
- 访问 Redis
- 调用 Eino
- 执行 Shell
- 读写文件
- 跨模块访问 infra
- 返回 PO / Entity

## 路由规范

统一挂载：

```text
/api/v1/{resource}
```

健康检查：

```text
/api/health
```

## 验证方式

启动：

```bash
go run ./cmd/api
```

检查端口：

```bash
ss -lntp | grep -E '8080|8000'
```

curl：

```bash
curl -i http://127.0.0.1:8080/api/health
curl -i http://127.0.0.1:8080/api/v1/{resource}
```

## 注意事项

- 如果 curl 连接失败，先检查服务是否启动。
- 如果 404，检查路由是否注册。
- 如果返回 HTML，可能请求打到了前端服务。
- 如果 500，看后端日志。
- 如果路径出现 `/api/api`，检查前端 API 文件。
- Handler 不执行长任务；长任务走 run / worker。

---

# Step 8 — module.go：模块装配

## 目标

显式装配 Repository、Service、Handler，并注册路由。

## 产出物

```text
internal/modules/{module}/module.go
```

## 必须包含

- `Module` struct
- `NewModule(deps ...)`
- `RegisterRoutes(r gin.IRouter)`
- 对外暴露必要的 QueryService / CommandService
- 注入跨模块依赖接口

## 典型装配

```text
DB
→ RepositoryImpl
→ CommandService / QueryService
→ Handler
→ RegisterRoutes
```

## 验证方式

```bash
go test ./...
go run ./cmd/api
curl -i http://127.0.0.1:8080/api/v1/{resource}
```

## 注意事项

- module.go 可以 import api / application / infra。
- module.go 不写业务逻辑。
- module.go 不处理 HTTP 细节。
- 所有模块还必须在 internal/app/modules.go 中统一加载。
- 如果忘记 RegisterRoutes，接口会 404。
- 如果忘记初始化 Service，可能 nil pointer。
- 如果跨模块依赖未注入，运行时会失败。

---

# Step 9 — 后端 curl 验收

## 目标

确认 HTTP → Handler → Service → Repository → MySQL 的完整链路可用。

## 验收命令

```bash
go test ./...
go run ./cmd/api
```

新开终端：

```bash
ss -lntp | grep -E '8080|8000'
curl -i http://127.0.0.1:8080/api/health
curl -i http://127.0.0.1:8080/api/v1/{resource}
```

创建：

```bash
curl -i -X POST http://127.0.0.1:8080/api/v1/{resource} \
  -H "Content-Type: application/json" \
  -d '{...}'
```

详情：

```bash
curl -i http://127.0.0.1:8080/api/v1/{resource}/{id}
```

更新：

```bash
curl -i -X PUT http://127.0.0.1:8080/api/v1/{resource}/{id} \
  -H "Content-Type: application/json" \
  -d '{...}'
```

删除：

```bash
curl -i -X DELETE http://127.0.0.1:8080/api/v1/{resource}/{id}
```

数据库：

```sql
SELECT * FROM <table> ORDER BY id DESC LIMIT 5;
```

## 验收标准

- `go test ./...` 通过。
- 后端服务正常监听。
- health 返回 200。
- 列表返回统一 Result。
- 创建成功。
- 查询详情成功。
- 更新成功。
- 删除成功。
- 错误响应符合统一格式。
- 敏感字段不泄露。
- 数据库记录符合预期。

---

# Step 10 — 前端 API 文件

## 目标

创建模块真实 API 封装，替换 mock 的基础。

## 产出物

```text
lattice-coding-web/src/api/{module}.ts
```

## 必须包含

- 类型定义
- list
- get
- create
- update
- delete
- enable / disable，按需
- test / sync / health，按需
- PageResult 类型复用

## 规范

- 使用 `src/api/request.ts`。
- `baseURL` 已经是 `/api`。
- API 方法路径从 `/v1/...` 开始。
- 不写 `/api/v1/...`。
- 不写死后端 IP 和端口。
- 不绕过统一 request。

## 验证方式

```bash
cd lattice-coding-web
npm run dev
```

浏览器 F12 Network：

- Request URL 是 `/api/v1/...`
- 不是 `/api/api/v1/...`
- 没有直接请求 `127.0.0.1:8080`
- 状态码 200 或合理错误码
- 响应结构符合 `code/message/data`

## 注意事项

- 前端页面不要继续使用 mock。
- 删除 mock 时确认空状态能显示。
- 如果页面一直 loading，先看 Network。
- 如果 Network 没请求，检查页面是否调用 API。
- 如果请求 404，检查后端路由和 Vite proxy。
- 如果连接失败，检查后端是否启动。

---

# Step 11 — 前端页面对接

## 目标

完成真实 API 数据驱动的页面。

## 产出物

```text
lattice-coding-web/src/views/{module}/
```

通常包含：

- 列表页
- 新增弹窗
- 编辑弹窗
- 删除确认
- 状态操作
- 详情展示
- 分页
- 空状态
- loading
- error 处理

## 推荐组件

- LatticeTable
- LatticeFormDialog
- useRequest
- useConfirm
- notify
- Element Plus 表格 / 表单 / 弹窗 / Tag / Switch

## 验证方式

```bash
cd lattice-coding-web
npm run dev
```

浏览器访问：

```text
http://服务器IP:3000/{module}
```

检查：

- 页面正常加载。
- 列表真实请求后端。
- 无 mock 数据。
- 新增成功。
- 编辑成功。
- 删除前确认。
- 删除后刷新。
- 分页正常。
- 错误提示正常。
- 敏感字段不显示。
- F12 无 CORS、404、500、Network Error。

## 注意事项

- 前端不要做业务唯一性校验的最终判断，最终以服务端为准。
- 前端不要展示密钥明文。
- 前端不要写死状态枚举之外的字符串。
- 前端不要重复拼 `/api`。
- 前端不要直接访问后端 IP。
- 表单提交后必须刷新列表。
- 删除操作必须确认。

---

# Step 12 — 运行时能力接入

## 触发条件

只有在以下条件满足后，才进入运行时能力：

- 后端 CRUD 验收通过。
- 前端页面对接完成。
- 用户确认需要接入运行时能力。
- 安全边界确认。

## 运行时能力类型

| 能力 | 所属层 |
|---|---|
| LLM 调用 | runtime/llm |
| Eino Agent | runtime/eino |
| Tool 调用 | runtime/tool |
| Redis Stream | runtime/event |
| SSE | api SSE handler |
| 文件读写 | WorkspaceExecutor |
| Shell 执行 | WorkspaceExecutor |
| Git diff | WorkspaceExecutor |
| 权限检查 | safety/application |
| 审计记录 | audit/application |
| 向量检索 | knowledge + pgvector |

## 等待用户确认

- 是否接真实 LLM。
- 是否只用 Ollama。
- 是否允许外部模型。
- 是否允许外网。
- 是否允许 Shell。
- 是否允许写文件。
- 是否只生成 patch。
- 是否需要人工审批。
- 是否写审计日志。
- 是否需要 Redis Stream。
- 是否需要限流。
- 是否需要熔断。
- 是否需要 fallback。

## 验证方式

按能力分别验证：

- LLM：轻量 prompt 返回 OK。
- Tool：只读工具先跑通。
- Shell：白名单命令测试。
- Redis Stream：写入和读取事件。
- SSE：浏览器接收事件。
- 文件写入：生成 patch，不直接改主目录。
- 审计：数据库有记录。

## 注意事项

- 不要在业务模块直接 import eino-ext。
- 不要在 Handler 中执行 Agent。
- 长任务必须走 run / worker。
- Worker 事件通过 Redis Stream 给 API SSE。
- 文件路径必须限制在 workspace 内。
- 写操作优先 patch。
- Shell 必须超时和白名单。
- 高风险操作必须审批。
- 所有事件必须可追踪。

---

# Step 13 — 完整验收

## 后端

```bash
go test ./...
go run ./cmd/api
curl -i http://127.0.0.1:8080/api/health
curl -i http://127.0.0.1:8080/api/v1/{resource}
```

## 前端

```bash
cd lattice-coding-web
npm run dev
```

浏览器：

```text
http://服务器IP:3000/{module}
```

## 数据库

```sql
SELECT * FROM <table> ORDER BY id DESC LIMIT 5;
```

## 浏览器 F12

检查：

- 请求路径正确。
- 状态码正确。
- 响应结构正确。
- 页面渲染正确。
- 无 CORS。
- 无 `/api/api`。
- 无 mock 数据。
- 无敏感字段。
- 无一直 loading。
- 无前端 JS 报错。

## 模块交付文档

每个模块完成后，必须输出：

1. 模块职责。
2. 不负责什么。
3. 表结构。
4. 核心 Entity / PO / DTO。
5. API 列表。
6. 后端调用链路。
7. 前端调用链路。
8. 跨模块依赖。
9. 安全与脱敏点。
10. 验收命令。
11. 已知限制。
12. 后续扩展点。
13. 面试讲解版本。

---

# 常见坑速查

| 现象 | 原因 | 修复 |
|---|---|---|
| `curl: Failed to connect` | 后端没启动或端口错 | `go run ./cmd/api`，`ss -lntp` |
| 404 | 路由没注册或路径不一致 | 检查 `RegisterRoutes` 和 `/api/v1` |
| `/api/api/...` | 前端 API 文件重复写 `/api` | API 方法从 `/v1/...` 开始 |
| 页面一直转圈 | 请求 pending / 后端没启动 / loading 没关闭 | 先看 F12 Network |
| import cycle | 转换函数放错层 | 按 Request/Response/Assembler/Converter 拆分 |
| application import api | 分层错误 | 移动 DTO->Response 到 api 包 |
| application import infra | 分层错误 | 依赖 domain.Repository 接口 |
| domain import GORM | 分层错误 | GORM 只能在 infra |
| api 直接查 DB | 分层错误 | 通过 application service |
| 删除目录 | Agent 误以为未使用 | 恢复标准目录结构 |
| API Key 出现在响应 | 脱敏遗漏 | 只返回 `api_key_set` 或类似布尔值 |
| 前端请求真实 IP | 绕过 Vite proxy | 使用相对路径 `/api` |
| 健康接口 404 | `/api` 前缀不一致 | 后端统一 `/api`，proxy 不 rewrite |
| 长任务卡死 API | Handler 直接执行 | 改成 run + worker |
| Redis 存业务事实 | 存储边界错误 | MySQL 存事实，Redis 存运行态 |
| 代码写死服务器文件系统 | 后续本地客户端难扩展 | 抽象 WorkspaceExecutor |

---

# Agent 通用执行模板

当用户要求开发模块时，先输出计划，不直接改代码：

```text
我将按 Lattice-Coding 业务模块全流程交付 Skill 推进 {module} 模块。

本次先进入 Step 1：咨询与设计，不写代码。

我会输出：
1. 模块职责
2. 模块不负责什么
3. 数据模型草案
4. API 草案
5. 跨模块依赖
6. 风险点
7. 需要你确认的问题

请确认后，我再进入 Step 2 更新 schema。
```

当用户确认后，再按步骤执行：

```text
请实现 {module} 模块 Step {n}。

要求：
1. 只做本步骤内容。
2. 不修改无关模块。
3. 不删除标准目录结构。
4. 不新增第三方依赖。
5. 遵守 api/application/domain/infra 分层。
6. 完成后运行 go test ./...。
7. 输出修改文件、结构体、方法、验证方式。
```

---

## 2. 跨模块依赖处理通用范式

### 2.1 核心原则

遵循 **依赖倒置原则（DIP）** 和 **接口隔离原则（ISP）**：

1. **依赖抽象而非具体实现**：定义接口而非直接依赖模块
2. **接口定义在被依赖模块**：由被依赖方定义接口，依赖方实现
3. **避免循环依赖**：通过接口解耦，禁止模块间双向引用
4. **依赖注入装配**：在应用入口统一装配跨模块依赖

### 2.2 范式一：单向依赖（A 依赖 B）

**场景**：Agent 创建时需要校验 Provider 是否存在

**实现步骤**：

1. **在依赖方（Agent）定义接口**（避免循环引用）
    ```go
    // agent/application/command_service.go
    type ModelConfigGetter interface {
        GetModelConfig(ctx context.Context, id uint64) (uint64, error)  // 返回 provider_id
        GetProvider(ctx context.Context, id uint64) error
    }
    ```

2. **在依赖方的 infra 层实现接口**（直接访问被依赖方的数据库表）
    ```go
    // agent/infra/persistence/provider_getter.go
    type ProviderGetter struct {
        db *gorm.DB
    }
    
    func (g *ProviderGetter) GetModelConfig(ctx context.Context, id uint64) (uint64, error) {
        // 直接查询 model_configs 表
    }
    
    func (g *ProviderGetter) GetProvider(ctx context.Context, id uint64) error {
        // 直接查询 providers 表
    }
    ```

3. **在依赖方的 application 层使用接口**
    ```go
    type CommandService struct {
        agentRepo         domain.AgentRepository
        modelConfigGetter ModelConfigGetter  // 注入接口
    }
    ```

4. **模块装配时注入实现**
    ```go
    // agent/module.go
    func NewModule(p *ModuleProvider) *Module {
        providerGetter := persistence.NewProviderGetter(p.DB)
        cmdSvc := application.NewCommandService(agentRepo, providerGetter)
    }
    ```

### 2.3 范式二：反向依赖检查（B 需要检查 A 的引用）

**场景**：删除 Provider/ModelConfig 前检查是否被 Agent 引用

**实现步骤**：

1. **在被检查方（Agent）定义接口**
    ```go
    // agent/domain/repository.go
    type AgentReferenceChecker interface {
        HasModelConfigReferences(ctx context.Context, modelConfigID uint64) (bool, error)
    }
    ```

2. **在被检查方的 infra 层实现接口**
    ```go
    // agent/infra/persistence/agent_ref_counter.go
    type AgentRefCounter struct {
        db *gorm.DB
    }
    
    func (r *AgentRefCounter) HasModelConfigReferences(ctx context.Context, modelConfigID uint64) (bool, error) {
        var count int64
        // 查询 agents 表统计引用数
        return count > 0, nil
    }
    ```

3. **在检查方（Provider）的 application 层依赖接口**
    ```go
    // provider/application/command_service.go
    type CommandService struct {
        providerRepo    domain.ProviderRepository
        agentRefChecker domain.AgentReferenceChecker  // 来自 agent/domain
    }
    ```

4. **应用入口统一装配**
    ```go
    // app/modules.go
    func InitModules(d *Dependencies) *Modules {
        // 先创建 Agent 模块（无依赖）
        agentModule := agent.NewModule(&agent.ModuleProvider{DB: d.MySQL})
        
        // 创建 AgentRefCounter 供 Provider 使用
        agentRefCounter := persistence.NewAgentRefCounter(d.MySQL)
        
        // 创建 Provider 模块（依赖 AgentRefCounter）
        providerModule := provider.NewModule(&provider.ModuleProvider{
            DB:           d.MySQL,
            AgentChecker: agentRefCounter,
        })
    }
    ```

### 2.4 范式三：双向依赖解耦

**场景**：Agent 和 Provider 互相需要对方的信息

**解决方案**：避免双向依赖，通过以下方式处理：

1. **避免双向引用**：不允许 A 引用 B 的 domain，同时 B 引用 A 的 domain
2. **使用数据层直接访问**：在 infra 层直接访问对方的数据库表（如范式一所示）
3. **事件驱动**：通过事件总线解耦（复杂场景）

## 3. 总结

跨模块依赖处理的核心是**通过接口解耦，依赖倒置**：

1. **接口定义位置**：依赖方需要的数据访问接口定义在依赖方内部（避免循环）
2. **接口实现位置**：依赖方的 infra 层直接访问数据库实现接口
3. **反向依赖处理**：被检查方定义接口，检查方依赖接口，应用入口注入实现
4. **装配顺序**：按依赖关系从无到有依次创建模块

这种设计保证了模块间的松散耦合，符合 Clean Architecture 和 DDD 的分层原则。
