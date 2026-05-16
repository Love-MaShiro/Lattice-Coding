# Lattice Coding

Lattice Coding 是一个 Go + Vue 实现的 AI Coding Agent 平台原型。当前重点是把 Web Chat、Agent Runtime、QueryEngine、LLM Runtime、Tool Runtime 和 Provider 管理链路跑通。

## 当前实现效果

- Provider / ModelConfig 管理：支持在页面或 API 中配置模型供应商、模型配置、启用/禁用、健康检查、同步模型。
- Agent 管理：Agent 绑定一个默认 `model_config_id`，Chat 请求会通过 Agent 找到模型配置。
- Chat 页面：
  - 新建会话、清空当前会话、会话列表、历史消息加载。
  - 执行模式下拉框：Direct / Workflow / PlanGraph / ReAct。
  - Direct 模式走流式输出。
  - ReAct 模式走非流式 QueryEngine，并在最终回答气泡里展示可折叠的“思考与工具调用”轨迹。
  - 压缩上下文按钮已接入占位/基础会话摘要链路。
- QueryEngine：
  - 普通 Chat completion 已迁入 `runtime/query.Engine.Run`。
  - Direct stream 已迁入 `runtime/query.Engine.Stream`。
  - Pure ReAct 已接入 `runtime/agent.Runtime`。
- LLM Runtime：
  - `runtime/llm` 已通过 `ModelConfigResolver` 与 `modules/provider` 解耦。
  - Provider 模块负责读取 Provider / ModelConfig、解密 API Key，并映射为 `llm.ResolvedModelConfig`。
- Tool Runtime：
  - 当前静态开放给 Agent 的工具包括 `file.list`、`file.read`、`code.grep`、`git.diff`、`file.edit`、`shell.run`。
  - 工具执行统一经过 `ToolExecutor`，权限、安全检查和审计/调用记录在工具层处理。
- Run 页面：
  - 可以查看 Run 列表和详情。
  - 预留 token / latency / cost 展示区域，当前取决于后端 Run/LLM 记录是否完整写入。
- 测试工作区：
  - Docker API 容器会把项目根目录下的 `test_workspace` 挂载到容器内 `/workspace`。
  - Agent 默认工作目录暂定为 `/workspace`，用于前端进行简单开发测试。

## 目录说明

```text
cmd/
  api/                 Web API 入口
  worker/              后台 Worker 入口

internal/
  app/                 启动、模块装配、依赖容器
  common/              config/db/redis/logger/http 等基础设施
  modules/             业务管理平面：provider、agent、chat、workflow、knowledge、run 等
  runtime/             Agent 执行平面：query、agent、llm、tool、prompt、context 等

lattice-coding-web/    Vue 前端
deployments/docker/    Docker Compose 开发环境
migrations/            MySQL/Postgres 初始化脚本
test_workspace/        当前 Agent 测试工作区
```

## Docker 配置需求

本地需要安装：

- Docker
- Docker Compose

默认开发端口：

| 服务 | 端口 | 说明 |
| --- | --- | --- |
| API | `8080` | Go 后端 |
| Web | `3000` | Vite 前端开发服务 |
| MySQL | `3306` | 业务数据 |
| Redis | `6379` | 缓存/上下文 |
| Postgres | `5432` | pgvector / 知识库预留 |

Docker Compose 文件位于：

```bash
deployments/docker/docker-compose.yml
```

API 容器关键环境变量：

```yaml
MYSQL_HOST=mysql
REDIS_HOST=redis
POSTGRES_HOST=postgres
LATTICE_WORKSPACE_DIR=/workspace
```

工作区挂载：

```yaml
../../test_workspace:/workspace
```

这意味着 Agent 在前端执行文件相关工具时，默认看到的是 `test_workspace` 目录，而不是整个项目根目录。

## 启动方式

启动后端依赖和 API：

```bash
docker compose -f deployments/docker/docker-compose.yml up -d --build
```

查看容器状态：

```bash
docker ps
```

验证 API：

```bash
curl http://localhost:8080/api/health
```

启动前端开发服务：

```bash
docker run --rm -it ^
  -p 3000:3000 ^
  -v "%cd%/lattice-coding-web:/app" ^
  -w /app ^
  -e VITE_API_PROXY_TARGET=http://host.docker.internal:8080 ^
  node:20-alpine sh -c "npm install && npm run dev -- --host 0.0.0.0"
```

如果不用 Docker 启动前端，也可以：

```bash
cd lattice-coding-web
npm install
npm run dev
```

打开页面：

```text
http://localhost:3000
```

## 供应商和模型管理

页面入口通常是 Provider / ModelConfig 管理页面。

推荐配置顺序：

1. 创建 Provider。
   - 选择 Provider Type，例如 OpenAI compatible、Ollama、Claude 等当前后端支持的类型。
   - 配置 Base URL 和 API Key。
2. 启用 Provider。
3. 创建 ModelConfig。
   - 选择 Provider。
   - 填写模型名称，例如 `gpt-4.1`、`deepseek-chat`、`qwen-plus`、本地 Ollama 模型名等。
   - 设置 temperature、top_p、max_tokens。
4. 启用 ModelConfig。
5. 可选：设置默认 ModelConfig。
6. 创建或更新 Agent，让 Agent 绑定该 `model_config_id`。

当前后端链路：

```text
Chat API
  -> modules/chat/application.CommandService
  -> runtime/query.Engine
  -> runtime/query.Strategy
  -> runtime/agent.Runtime 或 DirectChat
  -> runtime/llm.Executor
  -> runtime/llm.Router / Resolver
  -> modules/provider/application.ModelConfigResolver
  -> runtime/llm.Client
```

`runtime/llm` 不再直接 import `modules/provider/domain`。Provider 模块只负责 Provider / ModelConfig 管理和模型配置解析。

## Agent 执行方式如何选择

Chat 页面右上角有执行模式下拉框。

### Direct

适合普通聊天、简单问答、代码解释。

链路：

```text
Chat -> QueryEngine.Stream -> DirectChatStrategy.Stream -> llm.Executor.Stream
```

特点：

- 当前支持流式输出。
- 不调用工具。
- 不展示 ReAct 工具轨迹。

### ReAct

适合需要让模型决定是否调用工具的任务，例如查看工作区文件、读取文件、搜索代码、尝试执行安全命令。

链路：

```text
Chat -> QueryEngine.Run -> PureReActStrategy -> AgentRuntime -> LLM -> ToolExecutor -> LLM -> final
```

模型必须输出规范 JSON：

```json
{"type":"tool_call","reason":"short reason","tool":"tool.name","args":{}}
```

或：

```json
{"type":"final","answer":"final answer"}
```

页面效果：

- 最终回答气泡下方会出现“思考与工具调用”折叠区域。
- 可以看到每一步工具名、简短 reason、输入参数摘要、观察结果或错误。
- 不展示隐藏推理链，只展示后端 ReAct step 记录。

当前限制：

- ReAct 暂时是非流式。
- 工具列表是静态注入，还没有从 Agent 工具绑定表动态加载。
- 新建文件缺少专门的 `file.write` 工具。
- `shell.run` 的写文件重定向通常会被安全检查判定为需要审批，所以让 Agent “创建文件”可能失败，但失败轨迹会显示在气泡里。

### Workflow

当前作为执行模式入口保留。

第一阶段目标是后续接入 Workflow Engine / Workflow Node 执行，不建议现在用于真实任务。

### PlanGraph

当前作为规划模式入口保留。

第一阶段目标是让 LLM 生成 WorkflowSpec DAG，再交给 Workflow 执行。当前不建议用于真实任务。

## 在测试工作区里使用 ReAct

当前 Agent 的工作目录是 Docker 容器内：

```text
/workspace
```

对应宿主机目录：

```text
test_workspace/
```

可以在 Chat 页面选择 ReAct，然后询问：

```text
当前工作区有哪些文件？
```

期望 Agent 调用 `file.list`，并返回 `test_workspace` 下的文件列表。

也可以尝试：

```text
读取 hello.txt 的内容
```

期望 Agent 调用 `file.read`。

创建新文件目前不稳定，因为缺少 `file.write` 工具。比如：

```text
帮我在工作区内写一个 hello_world.c，能够在终端输出 hello world!
```

当前大概率会失败，原因是：

- `file.edit` 只能替换已有文件内容，不能用空 `old_string` 创建文件。
- `shell.run` 的 `echo ... > file` 会触发安全审批。

失败时页面应该仍展示 ReAct 轨迹，说明调用了哪些工具以及为什么失败。

## 常用 API

Provider：

```text
GET    /api/v1/providers
POST   /api/v1/providers
POST   /api/v1/providers/:id/test
POST   /api/v1/providers/:id/sync-models
```

ModelConfig：

```text
GET    /api/v1/model-configs
POST   /api/v1/model-configs
POST   /api/v1/model-configs/:id/test
POST   /api/v1/model-configs/:id/default
```

Agent：

```text
GET    /api/v1/agents
POST   /api/v1/agents
POST   /api/v1/agents/:id/enable
POST   /api/v1/agents/:id/disable
```

Chat：

```text
GET    /api/v1/chat/sessions
POST   /api/v1/chat/sessions
GET    /api/v1/chat/sessions/:id/messages
POST   /api/v1/chat/completions
POST   /api/v1/chat/stream
POST   /api/v1/chat/sessions/:id/compact
DELETE /api/v1/chat/sessions/:id
```

Run：

```text
GET /api/v1/runs
GET /api/v1/runs/:id
GET /api/v1/runs/:id/tool-invocations
```

## 构建和测试

后端测试：

```bash
go test ./...
```

Docker 方式运行 Go 测试：

```bash
docker run --rm -v "%cd%:/app" -w /app golang:1.22.0-alpine go test ./...
```

前端构建：

```bash
cd lattice-coding-web
npm run build
```

Docker 方式前端构建：

```bash
docker run --rm -v "%cd%/lattice-coding-web:/app" -w /app node:20-alpine npm run build
```

## 当前技术债和后续计划

- 实现安全的 `file.write` 工具，让 Agent 可以在 `test_workspace` 内创建文件。
- Agent 工具列表从静态列表改为读取 Agent 工具绑定配置。
- ReAct 支持流式事件，实时显示 tool call / observation。
- Workflow / PlanGraph 接入真实执行引擎。
- Run token / latency / cost 的完整统计和展示。
- 更严格的工作区边界和 shell 命令安全策略。
- Chat 的 `/clear`、`/compact`、`/cost`、`/plan` 等网页动作进一步产品化。
