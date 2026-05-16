<template>
  <PageContainer class="chat-page" title="Chat">
    <template #actions>
      <el-select
        v-model="selectedAgentId"
        class="agent-select"
        placeholder="选择 Agent"
        filterable
        :loading="agentLoading"
        @change="handleAgentChange"
      >
        <el-option v-for="agent in enabledAgents" :key="agent.id" :label="agent.name" :value="agent.id" />
      </el-select>
      <el-select v-model="executionMode" class="mode-select" :disabled="sending">
        <el-option label="Direct" value="direct_chat" />
        <el-option label="Workflow" value="fixed_workflow" />
        <el-option label="PlanGraph" value="plan_graph" />
        <el-option label="ReAct" value="pure_react" />
      </el-select>
      <el-tooltip content="刷新" placement="bottom">
        <el-button :icon="Refresh" circle @click="refreshAll" />
      </el-tooltip>
      <el-tooltip content="新建会话" placement="bottom">
        <el-button type="primary" :icon="Plus" circle @click="startNewChat" />
      </el-tooltip>
    </template>

    <div class="chat-shell">
      <aside class="session-pane">
        <div class="session-header">
          <span>会话</span>
          <el-tag size="small" type="info" effect="plain">{{ sessions.length }}</el-tag>
        </div>

        <el-scrollbar class="session-scroll" v-loading="sessionLoading">
          <button
            v-for="session in sessions"
            :key="session.id"
            class="session-item"
            :class="{ active: session.id === activeSessionId }"
            @click="selectSession(session.id)"
          >
            <span class="session-title">{{ session.title || `Session #${session.id}` }}</span>
            <span class="session-meta">
              <span>{{ formatTime(session.updated_at) }}</span>
              <span v-if="session.summary" class="summary-dot">已压缩</span>
            </span>
          </button>
          <el-empty v-if="!sessionLoading && sessions.length === 0" description="暂无会话" :image-size="72" />
        </el-scrollbar>
      </aside>

      <section class="message-pane">
        <div class="message-toolbar">
          <div class="message-title">
            <span>{{ activeSessionTitle }}</span>
            <el-tag v-if="activeSession?.status" size="small" effect="plain">{{ activeSession.status }}</el-tag>
            <el-tag size="small" type="success" effect="plain">{{ modeLabel }}</el-tag>
          </div>
          <div class="message-actions">
            <el-tooltip content="压缩上下文" placement="bottom">
              <el-button :icon="Files" circle plain :disabled="!activeSessionId || sending" @click="compactActiveSession" />
            </el-tooltip>
            <el-tooltip content="清空当前会话" placement="bottom">
              <el-button type="warning" :icon="Delete" circle plain :disabled="!activeSessionId || sending" @click="clearActiveSession" />
            </el-tooltip>
          </div>
        </div>

        <el-alert
          v-if="activeSession?.summary"
          class="summary-alert"
          type="info"
          :closable="false"
          show-icon
        >
          <template #title>当前会话已有压缩摘要，后续请求会携带摘要上下文。</template>
        </el-alert>

        <el-scrollbar ref="messageScrollbarRef" class="message-scroll" v-loading="messageLoading">
          <div v-if="messages.length > 0" class="message-list">
            <div v-for="message in messages" :key="message.id" class="message-row" :class="message.role">
              <div class="message-bubble" :class="{ pending: message.role === 'assistant' && !message.content && sending }">
                <div class="message-role">{{ roleLabel(message.role) }}</div>
                <div class="message-content">
                  <div v-if="message.content" class="markdown-body" v-html="renderMarkdown(message.content)" />
                  <span v-else-if="message.role === 'assistant' && sending" class="typing-indicator" aria-label="正在思考">
                    <span />
                    <span />
                    <span />
                  </span>
                </div>
                <details v-if="visibleTrace(message).length > 0" class="trace-collapse">
                  <summary>
                    <span>思考与工具调用</span>
                    <el-tag size="small" effect="plain">{{ visibleTrace(message).length }}</el-tag>
                  </summary>
                  <div class="trace-list">
                    <div v-for="event in visibleTrace(message)" :key="event.id" class="trace-item" :class="event.level">
                      <span class="trace-type">{{ event.title }}</span>
                      <span class="trace-content">{{ event.content }}</span>
                    </div>
                  </div>
                </details>
                <div class="message-time">{{ formatTime(message.created_at) }}</div>
              </div>
            </div>
          </div>
          <el-empty v-else-if="!messageLoading" description="暂无消息" :image-size="88" />
        </el-scrollbar>

        <div v-if="executionMode === 'pure_react' && traceEvents.length > 0" class="react-status">
          <span>ReAct 执行轨迹</span>
          <el-tag v-for="event in traceEvents.slice(-3)" :key="event.id" size="small" effect="plain">
            {{ event.title }}
          </el-tag>
        </div>

        <div class="composer">
          <el-input
            v-model="draft"
            type="textarea"
            resize="none"
            :rows="3"
            placeholder="输入消息。ReAct 模式会在消息上方展示运行事件和工具调用占位信息。"
            :disabled="sending"
            @keydown.enter.exact.prevent="sendMessage"
          />
          <el-button type="primary" :icon="Promotion" :loading="sending" @click="sendMessage">发送</el-button>
        </div>
      </section>
    </div>
  </PageContainer>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, ref } from 'vue'
import dayjs from 'dayjs'
import { ElMessageBox } from 'element-plus'
import type { ScrollbarInstance } from 'element-plus'
import { Delete, Files, Plus, Promotion, Refresh } from '@element-plus/icons-vue'

import { agentApi } from '@/api/agent'
import { chatApi, type ChatExecutionMode, type ChatMessage, type ChatSession, type ChatStreamEvent } from '@/api/chat'
import PageContainer from '@/components/PageContainer.vue'
import type { Agent } from '@/types/agent'
import { renderMarkdown } from '@/utils/markdown'
import { notifyError, notifySuccess, notifyWarning } from '@/utils/notify'

interface TraceEvent {
  id: number
  title: string
  content: string
  level: 'info' | 'success' | 'danger'
}

interface MessageMeta {
  trace?: TraceEvent[]
}

const agentLoading = ref(false)
const sessionLoading = ref(false)
const messageLoading = ref(false)
const sending = ref(false)
const selectedAgentId = ref<number>()
const activeSessionId = ref<number>()
const executionMode = ref<ChatExecutionMode>('direct_chat')
const agents = ref<Agent[]>([])
const sessions = ref<ChatSession[]>([])
const messages = ref<ChatMessage[]>([])
const draft = ref('')
const tempMessageId = ref(-1)
const traceId = ref(1)
const activeAssistantId = ref<number>()
const traceEvents = ref<TraceEvent[]>([])
const messageScrollbarRef = ref<ScrollbarInstance>()

const enabledAgents = computed(() => agents.value.filter((agent) => agent.enabled))
const activeSession = computed(() => sessions.value.find((session) => session.id === activeSessionId.value))
const activeSessionTitle = computed(() => activeSession.value?.title || '新会话')
const modeLabel = computed(() => {
  const labels: Record<ChatExecutionMode, string> = {
    direct_chat: 'Direct',
    fixed_workflow: 'Workflow',
    plan_graph: 'PlanGraph',
    pure_react: 'ReAct'
  }
  return labels[executionMode.value]
})

async function loadAgents() {
  agentLoading.value = true
  try {
    const result = await agentApi.list({ page: 1, page_size: 100 })
    agents.value = result.items || []
    if (!selectedAgentId.value && enabledAgents.value.length > 0) {
      selectedAgentId.value = enabledAgents.value[0].id
    }
  } finally {
    agentLoading.value = false
  }
}

async function loadSessions(keepActive = true) {
  sessionLoading.value = true
  try {
    const result = await chatApi.listSessions({ page: 1, page_size: 100 })
    sessions.value = result.items || []
    if (!keepActive || !sessions.value.some((session) => session.id === activeSessionId.value)) {
      activeSessionId.value = sessions.value[0]?.id
    }
  } finally {
    sessionLoading.value = false
  }
}

async function loadMessages(sessionId?: number) {
  if (!sessionId) {
    messages.value = []
    return
  }
  messageLoading.value = true
  try {
    messages.value = await chatApi.listMessages(sessionId, { limit: 100 })
    await scrollToBottom()
  } finally {
    messageLoading.value = false
  }
}

async function refreshAll() {
  await Promise.all([loadAgents(), loadSessions(true)])
  await loadMessages(activeSessionId.value)
}

function handleAgentChange() {
  startNewChat()
}

function startNewChat() {
  activeSessionId.value = undefined
  messages.value = []
  draft.value = ''
  traceEvents.value = []
  activeAssistantId.value = undefined
}

async function selectSession(id: number) {
  if (sending.value) return
  activeSessionId.value = id
  traceEvents.value = []
  activeAssistantId.value = undefined
  await loadMessages(id)
}

async function sendMessage() {
  const content = draft.value.trim()
  if (!content) return
  if (!activeSessionId.value && !selectedAgentId.value) {
    notifyWarning('请先选择 Agent')
    return
  }

  const now = new Date().toISOString()
  const userPreview = createPreviewMessage('user', content, now)
  const assistantPreview = createPreviewMessage('assistant', '', now)

  sending.value = true
  draft.value = ''
  traceEvents.value = []
  activeAssistantId.value = assistantPreview.id
  messages.value = [...messages.value, userPreview, assistantPreview]
  await nextTick()
  await scrollToBottom()

  let sessionId = activeSessionId.value
  try {
    if (!sessionId) {
      const session = await chatApi.createSession({
        agent_id: selectedAgentId.value!,
        title: buildTitle(content)
      })
      sessionId = session.id
      activeSessionId.value = session.id
      userPreview.session_id = session.id
      assistantPreview.session_id = session.id
      sessions.value = [session, ...sessions.value]
      messages.value = [...messages.value]
    }

    const result =
      executionMode.value === 'direct_chat'
        ? await sendStreamMessage(sessionId, content, assistantPreview)
        : await sendSingleRequestMessage(sessionId, content, assistantPreview)

    activeSessionId.value = result.session_id
    await loadSessions(true)
    await loadMessages(result.session_id)
  } catch (err) {
    messages.value = messages.value.filter((message) => message.id !== assistantPreview.id)
    notifyError(err instanceof Error ? err.message : '流式对话失败')
  } finally {
    sending.value = false
  }
}

async function sendStreamMessage(sessionId: number, content: string, assistantPreview: ChatMessage) {
  const typewriter = createTypewriter(assistantPreview)
  const result = await chatApi.stream(
    {
      agent_id: selectedAgentId.value,
      session_id: sessionId,
      message: content,
      mode: executionMode.value
    },
    {
      onDelta(delta) {
        typewriter.push(delta)
      },
      onEvent(event) {
        appendTraceEvent(event)
      }
    }
  )
  await typewriter.drain()
  return result
}

async function sendSingleRequestMessage(sessionId: number, content: string, assistantPreview: ChatMessage) {
  appendTraceEvent({ type: 'run.started', content: modeLabel.value })
  try {
    const result = await chatApi.complete({
      agent_id: selectedAgentId.value,
      session_id: sessionId,
      message: content,
      mode: executionMode.value
    })
    assistantPreview.id = result.message.id
    assistantPreview.session_id = result.message.session_id
    assistantPreview.content = result.message.content || result.content
    assistantPreview.token_count = result.message.token_count
    assistantPreview.meta = result.message.meta || '{}'
    assistantPreview.created_at = result.message.created_at
    assistantPreview.updated_at = result.message.updated_at
    activeAssistantId.value = undefined
    traceEvents.value = []
    messages.value = [...messages.value]
    return result
  } catch (err) {
    appendTraceEvent({
      type: 'run.error',
      message: err instanceof Error ? err.message : '请求失败'
    })
    throw err
  }
}

async function compactActiveSession() {
  if (!activeSessionId.value) return
  try {
    const session = await chatApi.compactSession(activeSessionId.value)
    sessions.value = sessions.value.map((item) => (item.id === session.id ? session : item))
    notifySuccess('上下文已压缩')
  } catch (err) {
    notifyError(err instanceof Error ? err.message : '压缩上下文失败')
  }
}

async function clearActiveSession() {
  if (!activeSession.value) return
  await ElMessageBox.confirm(`清空当前会话「${activeSessionTitle.value}」？该操作会删除这个会话。`, '清空会话', {
    type: 'warning',
    confirmButtonText: '清空',
    cancelButtonText: '取消'
  })
  await chatApi.deleteSession(activeSession.value.id)
  notifySuccess('会话已清空')
  startNewChat()
  await loadSessions(false)
  await loadMessages(activeSessionId.value)
}

function createPreviewMessage(role: 'user' | 'assistant', content: string, createdAt: string): ChatMessage {
  return {
    id: tempMessageId.value--,
    session_id: activeSessionId.value || 0,
    role,
    content,
    token_count: 0,
    meta: '{}',
    created_at: createdAt,
    updated_at: createdAt
  }
}

function createTypewriter(target: ChatMessage) {
  let buffer = ''
  let timer: number | undefined
  let idleResolve: (() => void) | undefined

  const tick = () => {
    const size = buffer.length > 24 ? 3 : 1
    const next = buffer.slice(0, size)
    buffer = buffer.slice(size)
    target.content += next
    messages.value = [...messages.value]
    scrollToBottom()

    if (buffer.length > 0) {
      timer = window.setTimeout(tick, 18)
      return
    }

    timer = undefined
    idleResolve?.()
    idleResolve = undefined
  }

  return {
    push(delta: string) {
      buffer += delta
      if (!timer) {
        timer = window.setTimeout(tick, 18)
      }
    },
    drain() {
      if (buffer.length === 0 && !timer) {
        return Promise.resolve()
      }
      return new Promise<void>((resolve) => {
        idleResolve = resolve
      })
    }
  }
}

function appendTraceEvent(event: ChatStreamEvent) {
  const type = event.type || ''
  if (type === 'llm.delta') return
  const map: Record<string, { title: string; level: TraceEvent['level'] }> = {
    'run.started': { title: '开始运行', level: 'info' },
    'llm.done': { title: '模型完成', level: 'success' },
    'run.finished': { title: '运行完成', level: 'success' },
    'run.error': { title: '运行失败', level: 'danger' },
    'tool.call': { title: '调用工具', level: 'info' },
    'tool.result': { title: '工具结果', level: 'success' },
    thought: { title: '思考中', level: 'info' }
  }
  const meta = map[type] || { title: type || '事件', level: 'info' }
  traceEvents.value.push({
    id: traceId.value++,
    title: meta.title,
    level: meta.level,
    content: event.message || event.content || event.run_id || modeLabel.value
  })
}

function visibleTrace(message: ChatMessage) {
  if (message.role !== 'assistant') return []
  const savedTrace = messageTrace(message)
  if (savedTrace.length > 0) {
    return savedTrace
  }
  if (message.id === activeAssistantId.value && traceEvents.value.length > 0) {
    return traceEvents.value
  }
  return []
}

function messageTrace(message: ChatMessage): TraceEvent[] {
  if (!message.meta || message.meta === '{}') return []
  try {
    const meta = JSON.parse(message.meta) as MessageMeta
    if (!Array.isArray(meta.trace)) return []
    return meta.trace
      .filter((item) => item && item.title)
      .map((item, index) => ({
        id: Number(item.id || index + 1),
        title: String(item.title || '事件'),
        content: String(item.content || ''),
        level: normalizeTraceLevel(item.level)
      }))
  } catch {
    return []
  }
}

function normalizeTraceLevel(level: string): TraceEvent['level'] {
  if (level === 'success' || level === 'danger') return level
  return 'info'
}

async function scrollToBottom() {
  await nextTick()
  messageScrollbarRef.value?.setScrollTop(100000)
}

function roleLabel(role: string) {
  switch (role) {
    case 'user':
      return 'User'
    case 'assistant':
      return 'Assistant'
    case 'system':
      return 'System'
    case 'tool':
      return 'Tool'
    default:
      return role
  }
}

function formatTime(value?: string) {
  return value ? dayjs(value).format('MM-DD HH:mm') : '-'
}

function buildTitle(message: string) {
  return message.length > 40 ? `${message.slice(0, 40)}...` : message
}

onMounted(async () => {
  await refreshAll()
})
</script>

<style scoped>
:deep(.page-content) {
  min-height: calc(100vh - 172px);
  padding: 0;
  overflow: hidden;
}

.agent-select {
  width: 220px;
}

.mode-select {
  width: 150px;
}

.chat-shell {
  flex: 1;
  min-height: 0;
  display: grid;
  grid-template-columns: 280px minmax(0, 1fr);
}

.session-pane {
  min-height: 0;
  display: flex;
  flex-direction: column;
  border-right: 1px solid var(--color-border);
  background: var(--color-bg-secondary);
}

.session-header,
.message-toolbar {
  height: 56px;
  flex: 0 0 56px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 0 16px;
  border-bottom: 1px solid var(--color-border);
}

.session-header {
  font-weight: 600;
  color: var(--color-text-primary);
}

.session-scroll,
.message-scroll {
  flex: 1;
  min-height: 0;
}

.session-item {
  width: 100%;
  min-height: 68px;
  border: 0;
  border-bottom: 1px solid var(--color-border);
  background: transparent;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  justify-content: center;
  gap: 6px;
  padding: 12px 16px;
  text-align: left;
  cursor: pointer;
}

.session-item:hover,
.session-item.active {
  background: var(--color-bg-card);
}

.session-item.active {
  box-shadow: inset 3px 0 0 var(--el-color-primary);
}

.session-title {
  color: var(--color-text-primary);
  font-size: 14px;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.session-meta,
.message-time {
  color: var(--color-text-secondary);
  font-size: 12px;
  font-variant-numeric: tabular-nums;
}

.session-meta {
  display: flex;
  justify-content: space-between;
  gap: 8px;
}

.summary-dot {
  color: var(--el-color-success);
}

.message-pane {
  min-width: 0;
  min-height: 0;
  display: flex;
  flex-direction: column;
  background: linear-gradient(180deg, #ffffff 0%, #f8fafc 100%);
}

.message-title,
.message-actions {
  min-width: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.message-title {
  color: var(--color-text-primary);
  font-weight: 600;
}

.summary-alert {
  border-radius: 0;
}

.message-list {
  display: flex;
  flex-direction: column;
  gap: 18px;
  padding: 24px 28px 28px;
}

.message-row {
  display: flex;
  justify-content: flex-start;
  padding-right: 18%;
}

.message-row.user {
  justify-content: flex-end;
  padding-right: 0;
  padding-left: 18%;
}

.message-bubble {
  max-width: min(760px, 100%);
  min-width: 96px;
  display: flex;
  flex-direction: column;
  gap: 9px;
  padding: 13px 16px 11px;
  border: 1px solid #e2e8f0;
  border-radius: 8px 8px 8px 4px;
  background: #ffffff;
  box-shadow: 0 10px 24px rgb(15 23 42 / 0.06);
}

.message-row.assistant .message-bubble {
  background: #f8fafc;
  border-color: #dbe3ef;
}

.message-row.user .message-bubble {
  color: #fff;
  border-color: var(--el-color-primary);
  border-radius: 8px 8px 4px 8px;
  background: var(--el-color-primary);
  box-shadow: 0 12px 26px rgb(64 158 255 / 0.22);
}

.message-bubble.pending {
  min-width: 112px;
  min-height: 72px;
  justify-content: center;
}

.message-row.user .message-time,
.message-row.user .message-role {
  color: rgba(255, 255, 255, 0.72);
}

.message-role {
  color: var(--color-text-secondary);
  font-size: 12px;
  font-weight: 600;
}

.message-content {
  min-height: 24px;
  color: inherit;
  font-size: 14px;
  line-height: 1.72;
  overflow-wrap: anywhere;
}

.markdown-body :deep(*) {
  margin-top: 0;
}

.markdown-body :deep(*:last-child) {
  margin-bottom: 0;
}

.markdown-body :deep(p) {
  margin-bottom: 10px;
  white-space: normal;
}

.trace-list,
.react-status {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding-top: 4px;
}

.trace-collapse {
  padding: 8px;
  border: 1px solid #dbe3ef;
  border-radius: 8px;
  background: #ffffff;
}

.trace-collapse summary {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  color: #334155;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
  list-style: none;
}

.trace-collapse summary::-webkit-details-marker {
  display: none;
}

.trace-item {
  display: grid;
  grid-template-columns: 72px minmax(0, 1fr);
  gap: 8px;
  padding: 6px 8px;
  border-radius: 6px;
  background: #eef2ff;
  color: #334155;
  font-size: 12px;
}

.trace-item.success {
  background: #ecfdf5;
}

.trace-item.danger {
  background: #fef2f2;
  color: #991b1b;
}

.trace-type {
  font-weight: 700;
}

.trace-content {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.react-status {
  flex: 0 0 auto;
  flex-direction: row;
  align-items: center;
  padding: 8px 16px;
  border-top: 1px solid var(--color-border);
  color: var(--color-text-secondary);
  font-size: 12px;
  background: #f8fafc;
}

.typing-indicator {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  height: 24px;
}

.typing-indicator span {
  width: 7px;
  height: 7px;
  border-radius: 999px;
  background: var(--color-text-tertiary);
  animation: typing-bounce 1.05s infinite ease-in-out;
}

.typing-indicator span:nth-child(2) {
  animation-delay: 0.14s;
}

.typing-indicator span:nth-child(3) {
  animation-delay: 0.28s;
}

.composer {
  flex: 0 0 auto;
  display: grid;
  grid-template-columns: minmax(0, 1fr) 96px;
  align-items: stretch;
  gap: 12px;
  padding: 16px;
  border-top: 1px solid var(--color-border);
  background: var(--color-bg-card);
}

.composer :deep(.el-textarea__inner) {
  min-height: 78px !important;
}

.composer .el-button {
  height: 78px;
}

@keyframes typing-bounce {
  0%, 80%, 100% {
    opacity: 0.35;
    transform: translateY(0);
  }

  40% {
    opacity: 1;
    transform: translateY(-4px);
  }
}

@media (max-width: 900px) {
  .chat-shell {
    grid-template-columns: 1fr;
  }

  .session-pane {
    max-height: 220px;
    border-right: 0;
    border-bottom: 1px solid var(--color-border);
  }

  .message-bubble {
    max-width: 100%;
  }

  .message-list {
    padding: 18px 14px 20px;
  }

  .message-row {
    padding-right: 8%;
  }

  .message-row.user {
    padding-left: 8%;
  }

  .composer {
    grid-template-columns: 1fr;
  }

  .composer .el-button {
    height: 40px;
  }
}
</style>
