import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/provider'
    },
    {
      path: '/provider',
      name: 'Provider',
      component: () => import('@/views/provider/ProviderList.vue'),
      meta: { title: '模型管理' }
    },
    {
      path: '/agent',
      name: 'Agent',
      component: () => import('@/views/agent/AgentList.vue'),
      meta: { title: 'Agent 管理' }
    },
    {
      path: '/chat',
      name: 'Chat',
      component: () => import('@/views/chat/ChatPage.vue'),
      meta: { title: '对话' }
    }
  ]
})

export default router
