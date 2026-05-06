# Skill：Lattice-Coding 新增模型供应商 Adapter

## 适用场景

当项目已经有 Provider / ModelConfig 管理能力，并且需要新增一种模型供应商时使用。

例如：
- "新增 Claude 供应商"
- "新增 Gemini 供应商"
- "新增 OpenRouter 供应商"
- "新增火山方舟 Ark 供应商"
- "新增本地 OpenAI-compatible 网关"
- "新增实验室内部模型服务"

本Skill的目标是：新增供应商时，只新增对应的 Adapter 类，而不需要修改其他代码。
（例如：新增 Claude 供应商时，只需要新增 ClaudeAdapter 类，而不需要修改其他代码）

## 核心设计原则

新增供应商时，禁止让 `provider`、`agent`、`chat`、`run` 模块直接依赖具体 SDK。

正确边界是：
- provider 模块：管理供应商配置
- model_config：管理具体模型配置
- runtime/llm：根据 provider_type 路由到对应 Adapter
- Adapter：封装某个供应商的认证、模型列表、连通性测试、ChatModel 创建
- chat/run/agent：只使用统一 LLMExecutor，不关心具体供应商

```text
internal/runtime/llm/
├── adapter/
│   ├── adapter.go
│   ├── registry.go
│   ├── openai_compatible.go
│   ├── ollama.go
│   ├── claude.go
│   └── ...
│
├── factory.go
├── executor.go
├── health.go
├── model_lister.go
├── types.go
└── errors.go
```

| 文件 | 作用 |
|---|---|
| `adapter.go` | 定义统一 ProviderAdapter 接口。 |
| `registry.go` | 注册和查找 Adapter。 |
| `openai_compatible.go` | OpenAI-compatible 供应商适配。 |
| `ollama.go` | Ollama 适配。 |
| `{provider}.go` | 新供应商适配实现。 |
| `factory.go` | 根据 provider_type 获取 Adapter 并创建 ChatModel。 |
| `health.go` | 连通性测试入口。 |
| `model_lister.go` | 模型列表同步入口。 |

# Step 1：分析目标供应商API
## 目标

在写代码前，先搞清楚目标供应商和现有 Adapter 的差异。

## 必须调研的问题

| 问题                   | 说明                                                           |
| -------------------- | ------------------------------------------------------------ |
| 供应商类型值               | `provider_type` 应该叫什么，例如 `claude`、`gemini`、`ark`。            |
| 是否 OpenAI-compatible | 如果兼容 `/v1/chat/completions`，优先复用 OpenAI-compatible Adapter。  |
| 认证方式                 | Bearer Token、API Key Header、双 Header、无认证、自定义 Header。         |
| 必填 auth_config 字段    | 例如 `api_key`、`api_version`、`anthropic_version`、`deployment`。 |
| base_url 默认值         | 官方默认地址是什么，是否允许用户自定义。                                         |
| Chat 调用路径            | 是 `/v1/chat/completions`，还是供应商自定义路径。                         |
| 列模型接口                | 是否支持拉取模型列表，路径和返回结构是什么。                                       |
| 流式响应格式               | 是否标准 SSE，还是供应商自定义流式格式。                                       |
| Tool Calling         | 是否支持工具调用，字段格式是否兼容 Eino。                                      |
| JSON Mode            | 是否支持结构化输出。                                                   |
| Embedding            | 是否支持 embedding 模型。                                           |
| 错误结构                 | 错误响应字段是什么，如何提取错误码和错误摘要。                                      |
| 限流错误                 | 429 或自定义错误码如何识别。                                             |
| 认证错误                 | 401/403 或自定义错误码如何识别。                                         |

## 产出物

输出一份 API 特征说明：

```text
供应商：xxx
provider_type：xxx
是否 OpenAI-compatible：是 / 否
认证方式：xxx
默认 base_url：xxx
Chat 路径：xxx
模型列表路径：xxx
流式格式：xxx
是否支持 Tool Calling：是 / 否 / 未知
是否支持模型列表同步：是 / 否
是否第一阶段接入：连通性测试 / ChatModel / 模型同步 / 流式
```


## 等待用户确认

必须等待用户确认：

```text
API 特征确认，可以进入 Adapter 实现。
```

## 注意事项

如果目标供应商本质上是 OpenAI-compatible，不要新写完整 Adapter，优先在现有 `openai_compatible` Adapter 中通过 `provider_type` 分支或配置处理。

---

# Step 2 — 确认 Provider 数据模型是否够用

## 目标

判断新增供应商是否需要修改 `provider` 或 `model_config` 表。

当前推荐字段：

```text
provider:
- provider_type
- base_url
- auth_type
- api_key_ciphertext
- auth_config_ciphertext
- config
- enabled
- health_status

model_config:
- provider_id
- name
- model
- model_type
- params
- capabilities
- is_default
- enabled
```

## 判断规则

| 情况 | 做法 |
|---|---|
| 只需要 API Key | 使用 `api_key_ciphertext`。 |
| 需要特殊 Header | 放入 `auth_config_ciphertext`。 |
| 需要 `api_version` | 放入 `auth_config_ciphertext` 或 `config`，敏感则加密。 |
| 需要 `deployment` | 放入 `model_config.params` 或 `auth_config`，视语义决定。 |
| 需要 `region` | 放入 `provider.config`。 |
| 需要模型能力声明 | 放入 `model_config.capabilities`。 |
| 需要高频查询 | 不要放 JSON，应新增字段。 |

## 等待用户确认

如果需要新增字段，必须先说明：

```text
新增字段：
原因：
影响的表：
是否需要迁移旧数据：
```

等待用户确认后再改 schema。

## 产出物

- 是否需要改 schema 的结论
- 如需要，给出 migration SQL
- 如不需要，说明用哪些已有字段表达差异

---

# Step 3 — 定义统一 Adapter 接口

## 目标

保证新增供应商只需要实现统一接口。

## 推荐接口

```go
package adapter

import (
    "context"

    "github.com/cloudwego/eino/components/model"
)

type ProviderAdapter interface {
    SupportedTypes() []string

    BuildChatModel(ctx context.Context, cfg RuntimeModelConfig) (model.ToolCallingChatModel, error)

    TestConnection(ctx context.Context, cfg RuntimeModelConfig) (*ConnectionTestResult, error)

    ListModels(ctx context.Context, cfg RuntimeProviderConfig) ([]ModelInfo, error)

    NormalizeError(err error) *LLMError
}
```

如果当前阶段暂时不支持 Tool Calling，可以先使用更小接口：

```go
type ProviderAdapter interface {
    SupportedTypes() []string
    BuildChatModel(ctx context.Context, cfg RuntimeModelConfig) (model.ChatModel, error)
    TestConnection(ctx context.Context, cfg RuntimeModelConfig) (*ConnectionTestResult, error)
    ListModels(ctx context.Context, cfg RuntimeProviderConfig) ([]ModelInfo, error)
}
```

## 推荐类型

```go
type RuntimeProviderConfig struct {
    ProviderID   uint64
    ProviderType string
    BaseURL      string
    AuthType     string
    APIKey       string
    AuthConfig   map[string]any
    Config       map[string]any
    Enabled      bool
}

type RuntimeModelConfig struct {
    RuntimeProviderConfig

    ModelConfigID uint64
    Model         string
    ModelType     string
    Params        map[string]any
    Capabilities  map[string]any
}

type ModelInfo struct {
    ID           string
    Name         string
    ModelType    string
    Capabilities map[string]any
}

type ConnectionTestResult struct {
    Success      bool
    LatencyMS    int64
    ModelCount   int
    ErrorCode    string
    ErrorMessage string
    HealthStatus string
}

type LLMError struct {
    Code        string
    Message     string
    Retryable   bool
    AuthError   bool
    RateLimited bool
}
```

## 产出物

```text
internal/runtime/llm/adapter/adapter.go
```

包含：

- `ProviderAdapter`
- `RuntimeProviderConfig`
- `RuntimeModelConfig`
- `ModelInfo`
- `ConnectionTestResult`
- `LLMError`

## 验证方式

```bash
go test ./...
```

## 注意事项

- Adapter 接口只放供应商差异能力。
- 不要把 Provider CRUD 放进 Adapter。
- 不要让 Adapter 访问 MySQL。
- Adapter 不负责读取密钥密文，它只接收已解密后的运行时配置。
- Adapter 不负责审计落库，审计由上层 `runtime/llm` 或 `audit` 模块处理。

---

# Step 4 — 实现目标供应商 Adapter

## 目标

新增一个 `{provider}.go`，实现 `ProviderAdapter`。

## 文件位置

```text
internal/runtime/llm/adapter/{provider}.go
```

## 必须实现

```go
type XxxAdapter struct {
    // 可注入 logger、http client、配置
}

func NewXxxAdapter() *XxxAdapter

func (a *XxxAdapter) SupportedTypes() []string

func (a *XxxAdapter) BuildChatModel(ctx context.Context, cfg RuntimeModelConfig) (model.ChatModel, error)

func (a *XxxAdapter) TestConnection(ctx context.Context, cfg RuntimeModelConfig) (*ConnectionTestResult, error)

func (a *XxxAdapter) ListModels(ctx context.Context, cfg RuntimeProviderConfig) ([]ModelInfo, error)

func (a *XxxAdapter) NormalizeError(err error) *LLMError
```

## 实现要求

| 方法 | 要求 |
|---|---|
| `SupportedTypes` | 返回支持的 `provider_type`，可多个。 |
| `BuildChatModel` | 根据 cfg 创建 Eino ChatModel 或 eino-ext 模型。 |
| `TestConnection` | 发送轻量请求，例如“请只回复 OK”。 |
| `ListModels` | 调供应商模型列表接口；不支持则返回 unsupported。 |
| `NormalizeError` | 把供应商错误转成统一错误。 |

## 注意事项

- 不要在 Adapter 中打印 API Key。
- 不要返回完整 `auth_config`。
- 不要在 Adapter 中访问数据库。
- 不要在 Adapter 中处理 Provider 启用禁用逻辑。
- 不要在 Adapter 中决定 fallback。
- 不要在 Adapter 中写审计日志。
- 不要在 Adapter 中直接读配置文件。
- 供应商不支持的能力要返回明确错误，不要写假实现。
- 如果目标供应商兼容 OpenAI API，优先复用 OpenAI-compatible 实现。

## 验证方式

```bash
go test ./...
```

如果有真实配置，再做：

```bash
curl -i -X POST http://127.0.0.1:8080/api/v1/model-configs/{id}/test
```

---

# Step 5 — 注册 Adapter

## 目标

让 Factory 可以根据 `provider_type` 找到对应 Adapter。

## 推荐 Registry

```go
type Registry struct {
    adapters map[string]ProviderAdapter
}

func NewRegistry(adapters ...ProviderAdapter) *Registry

func (r *Registry) Register(adapter ProviderAdapter)

func (r *Registry) Get(providerType string) (ProviderAdapter, bool)
```

## 注册位置

```text
internal/runtime/llm/adapter/registry.go
```

或：

```text
internal/runtime/llm/factory.go
```

推荐：

```go
registry := adapter.NewRegistry(
    adapter.NewOpenAICompatibleAdapter(),
    adapter.NewOllamaAdapter(),
    adapter.NewXxxAdapter(),
)
```

## 产出物

- Adapter 注册到 Registry
- Factory 能通过 `provider_type` 获取 Adapter

## 验证方式

```bash
go test ./...
```

新增单元测试：

```text
provider_type = xxx
registry.Get("xxx") 能返回 XxxAdapter
```

## 注意事项

- 新增供应商时，只允许改注册位置和新增 Adapter。
- 不要修改已有 Adapter 的逻辑，除非是在修复通用 bug。
- 不要在业务模块中 `switch provider_type`。
- `provider_type` 路由只应出现在 `runtime/llm` 的 Adapter Registry 或 Factory 中。

---

# Step 6 — 更新 ProviderType 枚举和校验

## 目标

让后端允许创建新 `provider_type` 的 Provider。

## 可能需要修改

```text
internal/modules/provider/domain/value_object.go
internal/modules/provider/api/request.go
internal/modules/provider/application/command_service.go
```

## 产出物

- 新增 `provider_type` 枚举值
- 创建 Provider 时校验通过
- 前端选择项可展示该供应商

## 验证方式

创建 Provider：

```bash
curl -i -X POST http://127.0.0.1:8080/api/v1/providers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "xxx-provider",
    "provider_type": "xxx",
    "base_url": "https://xxx",
    "auth_type": "bearer",
    "api_key": "test",
    "enabled": true
  }'
```

## 注意事项

- 如果当前 `provider_type` 没有强校验，可以不改。
- 如果有枚举校验，必须新增合法值。
- 不要为了新增供应商删除校验。
- 不要把供应商差异逻辑写进 provider application。

---

# Step 7 — 实现模型连通性测试

## 目标

验证新供应商的 ChatModel 可以被创建并调用。

## 调用链

```text
POST /api/v1/model-configs/{id}/test
→ provider/application 查询 Provider + ModelConfig
→ 解密 API Key
→ 组装 RuntimeModelConfig
→ runtime/llm.HealthChecker
→ Factory.GetAdapter(provider_type)
→ Adapter.BuildChatModel
→ ChatModel.Generate
→ 更新 provider_health
```

## 验证方式

```bash
curl -i -X POST http://127.0.0.1:8080/api/v1/model-configs/{id}/test
```

检查数据库：

```sql
SELECT health_status, last_checked_at, last_error
FROM provider
WHERE id = <provider_id>;

SELECT *
FROM provider_health
WHERE provider_id = <provider_id>
ORDER BY id DESC
LIMIT 5;
```

## 注意事项

- prompt 固定为“请只回复 OK”或更短 ping prompt。
- 超时 10s。
- 失败要更新健康状态。
- 错误只返回摘要。
- 日志禁止输出密钥。
- 不要把连通性测试写在 Handler 里。
- 不要在 provider 模块直接 import 具体供应商 SDK。

---

# Step 8 — 实现模型列表同步

## 目标

支持从新供应商拉取模型列表，写入 `model_config`。

## 调用链

```text
POST /api/v1/providers/{id}/sync-models
→ provider/application.SyncModels
→ runtime/llm.ModelLister
→ Adapter.ListModels
→ model_config 去重插入
```

## 验证方式

```bash
curl -i -X POST http://127.0.0.1:8080/api/v1/providers/{id}/sync-models
```

检查数据库：

```sql
SELECT provider_id, name, model, model_type, enabled
FROM model_config
WHERE provider_id = <provider_id>;
```

## 注意事项

- 如果供应商不支持列模型，返回 unsupported。
- 不要写假模型列表。
- 已存在模型不要重复插入。
- 不要覆盖用户手动配置的 `params`。
- `enabled` 默认值必须等待用户确认。
- 同步过程建议使用事务。
- 同步失败不能影响已有 `model_config`。

---

# Step 9 — 前端 Provider 页面更新

## 目标

让前端可以选择新供应商类型，并填写该供应商需要的鉴权字段。

## 修改位置

```text
lattice-coding-web/src/views/provider/
lattice-coding-web/src/api/provider.ts
```

## 产出物

- Provider 类型下拉框新增供应商
- 根据 `provider_type` 显示不同 auth 表单字段
- `base_url` 默认值自动填充或提示
- `auth_type` 默认值合理
- `capabilities` 或 `params` 可选配置

## 表单建议

| 字段 | 说明 |
|---|---|
| `provider_type` | 新供应商类型 |
| `base_url` | 默认值可预填，允许修改 |
| `auth_type` | bearer / api_key / custom_header / none |
| `api_key` | 普通密钥 |
| `auth_config` | 复杂鉴权 JSON，可先用文本框 |
| `config` | 非敏感扩展配置 |

## 验证方式

浏览器访问：

```text
http://服务器IP:3000/provider
```

检查：

- 可以创建新供应商
- 可以创建新模型
- 可以点击测试
- 错误能提示
- 不展示 API Key 明文
- F12 Network 路径正确

## 注意事项

- 前端不要写死后端 IP。
- API 方法从 `/v1/...` 开始，不写 `/api/v1/...`。
- 不要把 API Key 回显到编辑表单。
- 只显示 `api_key_set`。
- `auth_config` 如包含敏感字段，编辑时不要原样回显密文。

---

# Step 10 — 完整验收

## 后端验收

```bash
go test ./...
go run ./cmd/api
```

新开终端：

```bash
curl -i http://127.0.0.1:8080/api/health
curl -i http://127.0.0.1:8080/api/v1/providers
```

创建供应商：

```bash
curl -i -X POST http://127.0.0.1:8080/api/v1/providers \
  -H "Content-Type: application/json" \
  -d '{...}'
```

创建模型：

```bash
curl -i -X POST http://127.0.0.1:8080/api/v1/model-configs \
  -H "Content-Type: application/json" \
  -d '{...}'
```

测试模型：

```bash
curl -i -X POST http://127.0.0.1:8080/api/v1/model-configs/{id}/test
```

同步模型：

```bash
curl -i -X POST http://127.0.0.1:8080/api/v1/providers/{id}/sync-models
```

## 前端验收

```bash
cd lattice-coding-web
npm run dev
```

浏览器：

```text
http://服务器IP:3000/provider
```

检查：

- 新供应商可选
- 可保存 Provider
- 可保存 ModelConfig
- 可测试连通性
- 可同步模型列表
- API Key 不泄露
- 错误提示合理
- Network 无 `/api/api`

---

# 常见坑

## 1. 供应商其实是 OpenAI-compatible，却新写 Adapter

如果 Chat API、鉴权、流式格式都兼容 OpenAI，应复用 OpenAI-compatible Adapter。

新增供应商只需要：

```text
provider_type = xxx
base_url = xxx
auth_type = bearer
model = xxx
```

---

## 2. 把供应商差异写进 provider Service

错误：

```text
provider/application 里 switch provider_type 调 SDK
```

正确：

```text
provider/application 查询配置
runtime/llm Factory 按 provider_type 找 Adapter
Adapter 处理供应商差异
```

---

## 3. 在 Handler 中测试模型

错误：

```text
api.Handler 里直接 new ChatModel
```

正确：

```text
Handler → application → runtime/llm HealthChecker → Adapter
```

---

## 4. auth_config 明文泄露

禁止：

- 返回 `api_key`
- 返回 `api_key_ciphertext`
- 返回 `auth_config_ciphertext`
- 日志打印 auth_config
- 错误中包含 token

允许：

- `api_key_set: true`
- `auth_config_set: true`
- 错误摘要

---

## 5. 模型同步写假数据

如果供应商不支持 list models，必须返回 unsupported，不允许为了页面好看写死模型。

---

## 6. 新增供应商改动太多已有代码

理想改动范围：

```text
新增 internal/runtime/llm/adapter/{provider}.go
修改 adapter registry 注册
必要时新增 provider_type 枚举
必要时更新前端供应商选项
必要时新增 schema 字段，但要先确认
```

不应该大面积改：

```text
provider CRUD
agent
chat
run
api response
数据库基础结构
统一 request.ts
```

---

## 7. import cycle

Adapter 应在 `runtime/llm` 下，不要让 `provider` 模块反向依赖 runtime 后又被 runtime 依赖 provider。

推荐：

```text
provider/application 输出 RuntimeProviderDTO / RuntimeModelConfigDTO
runtime/llm 定义自己的 RuntimeModelConfig
由 application 或组装层做转换
```

不要：

```text
runtime/llm import provider/application
provider/application import runtime/llm
```

如有循环，抽出中立 DTO 或在上层做转换。

---

# 给 Agent 的通用提示词

## 先分析，不写代码

```text
请为 Lattice-Coding 新增一个模型供应商 Adapter：{provider_name}。

先不要写代码，请先完成 Step 1 API 特征分析，回答：
1. provider_type 应该是什么？
2. 是否 OpenAI-compatible？
3. 认证方式是什么？
4. 默认 base_url 是什么？
5. Chat 调用路径是什么？
6. 是否支持流式？
7. 是否支持 Tool Calling？
8. 是否支持模型列表接口？
9. 需要哪些 auth_config 字段？
10. 是否需要修改 provider 或 model_config 表？

请输出推荐接入方案，并标注哪些地方需要我确认。
```

## 确认后实现

```text
请实现 {provider_name} Adapter。

要求：
1. 只新增或修改 runtime/llm/adapter 相关代码。
2. 不修改 provider CRUD 逻辑。
3. 不修改 agent/chat/run 模块。
4. 不让 provider/api 或 provider/application 直接依赖具体 SDK。
5. Adapter 实现统一 ProviderAdapter 接口。
6. 在 Registry 中注册该 Adapter。
7. 支持 TestConnection。
8. 如果供应商支持模型列表，实现 ListModels；不支持则返回 unsupported。
9. 禁止打印或返回 API Key。
10. 完成后运行 go test ./...。
11. 输出修改文件、结构体、方法、验证命令。
```

---
