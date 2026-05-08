<template>
  <PageContainer class="chat-page" title="对话">
    <template #actions>
      <el-select
        v-model="selectedAgentId"
        class="agent-select"
        placeholder="选择 Agent"
        filterable
        :loading="agentLoading"
        @change="handleAgentChange"
      >
        <el-option
          v-for="agent in enabledAgents"
          :key="agent.id"
          :label="agent.name"
          :value="agent.id"
        />
      </el-select>
      <el-tooltip content="刷新" placement="bottom">
        <el-button :icon="Refresh" circle @click="refreshAll" />
      </el-tooltip>
      <el-tooltip content="新对话" placement="bottom">
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
            <span class="session-meta">{{ formatTime(session.updated_at) }}</span>
          </button>
          <el-empty v-if="!sessionLoading && sessions.length === 0" description="暂无会话" :image-size="72" />
        </el-scrollbar>
      </aside>

      <section class="message-pane">
        <div class="message-toolbar">
          <div class="message-title">
            <span>{{ activeSessionTitle }}</span>
            <el-tag v-if="activeSession?.status" size="small" effect="plain">
              {{ activeSession.status }}
            </el-tag>
          </div>
          <el-tooltip v-if="activeSessionId" content="删除会话" placement="bottom">
            <el-button type="danger" :icon="Delete" circle plain @click="deleteActiveSession" />
          </el-tooltip>
        </div>

        <el-scrollbar ref="messageScrollbarRef" class="message-scroll" v-loading="messageLoading">
          <div v-if="messages.length > 0" class="message-list">
            <div
              v-for="message in messages"
              :key="message.id"
              class="message-row"
              :class="message.role"
            >
              <div class="message-bubble">
                <div class="message-role">{{ roleLabel(message.role) }}</div>
                <div class="message-content">
                  <span v-if="message.content">{{ message.content }}</span>
                  <span v-else-if="message.role === 'assistant' && sending" class="typing-dot">正在思考</span>
                </div>
                <div class="message-time">{{ formatTime(message.created_at) }}</div>
              </div>
            </div>
          </div>
          <el-empty v-else-if="!messageLoading" description="暂无消息" :image-size="88" />
        </el-scrollbar>

        <div class="composer">
          <el-input
            v-model="draft"
            type="textarea"
            resize="none"
            :rows="3"
            placeholder="输入消息"
            :disabled="sending"
            @keydown.enter.exact.prevent="sendMessage"
          />
          <el-button type="primary" :icon="Promotion" :loading="sending" @click="sendMessage">
            发送
          </el-button>
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
import { Delete, Plus, Promotion, Refresh } from '@element-plus/icons-vue'

import { agentApi } from '@/api/agent'
import { chatApi, type ChatMessage, type ChatSession } from '@/api/chat'
import PageContainer from '@/components/PageContainer.vue'
import type { Agent } from '@/types/agent'
import { notifyError, notifySuccess, notifyWarning } from '@/utils/notify'

const agentLoading = ref(false)
const sessionLoading = ref(false)
const messageLoading = ref(false)
const sending = ref(false)
const selectedAgentId = ref<number>()
const activeSessionId = ref<number>()
const agents = ref<Agent[]>([])
const sessions = ref<ChatSession[]>([])
const messages = ref<ChatMessage[]>([])
const draft = ref('')
const tempMessageId = ref(-1)
const messageScrollbarRef = ref<ScrollbarInstance>()

const enabledAgents = computed(() => agents.value.filter((agent) => agent.enabled))
const activeSession = computed(() => sessions.value.find((session) => session.id === activeSessionId.value))
const activeSessionTitle = computed(() => activeSession.value?.title || '新对话')

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
}

async function selectSession(id: number) {
  if (sending.value) return
  activeSessionId.value = id
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

    const typewriter = createTypewriter(assistantPreview)
    const result = await chatApi.stream(
      {
        agent_id: selectedAgentId.value,
        session_id: sessionId,
        message: content
      },
      {
        onDelta(delta) {
          typewriter.push(delta)
        }
      }
    )

    await typewriter.drain()
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

async function deleteActiveSession() {
  if (!activeSession.value) return
  await ElMessageBox.confirm(`确认删除会话「${activeSessionTitle.value}」？`, '删除会话', {
    type: 'warning',
    confirmButtonText: '删除',
    cancelButtonText: '取消'
  })
  await chatApi.deleteSession(activeSession.value.id)
  notifySuccess('会话已删除')
  activeSessionId.value = undefined
  messages.value = []
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
  width: 240px;
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

.message-pane {
  min-width: 0;
  min-height: 0;
  display: flex;
  flex-direction: column;
  background: var(--color-bg-card);
}

.message-title {
  min-width: 0;
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--color-text-primary);
  font-weight: 600;
}

.message-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 20px;
}

.message-row {
  display: flex;
  justify-content: flex-start;
}

.message-row.user {
  justify-content: flex-end;
}

.message-bubble {
  max-width: min(720px, 78%);
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px 14px;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  background: var(--color-bg-secondary);
}

.message-row.user .message-bubble {
  color: #fff;
  border-color: var(--el-color-primary);
  background: var(--el-color-primary);
}

.message-row.user .message-time,
.message-row.user .message-role {
  color: rgba(255, 255, 255, 0.78);
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
  line-height: 1.7;
  white-space: pre-wrap;
  overflow-wrap: anywhere;
}

.typing-dot::after {
  display: inline-block;
  width: 20px;
  content: '...';
  animation: typing-pulse 1s infinite steps(3, end);
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

@keyframes typing-pulse {
  0% { opacity: 0.35; }
  50% { opacity: 1; }
  100% { opacity: 0.35; }
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
    max-width: 88%;
  }

  .composer {
    grid-template-columns: 1fr;
  }

  .composer .el-button {
    height: 40px;
  }
}
</style>
