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
      component: () => import('@/views/provider/ProviderList.vue')
    },
    {
      path: '/agent',
      name: 'Agent',
      component: () => import('@/views/agent/AgentList.vue')
    },
    {
      path: '/chat',
      name: 'Chat',
      component: () => import('@/views/chat/ChatPage.vue')
    }
  ]
})

export default router
