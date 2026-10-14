import { createRouter, createWebHistory } from 'vue-router'

import QueueReceive from '../views/QueueReceive.vue'
import QueueTicket from '../views/QueueTicket.vue'
import QueueReset from '../views/QueueReset.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'queue-receive',
      component: QueueReceive,
    },
    {
      path: '/ticket/:id',
      name: 'queue-ticket',
      component: QueueTicket,
    },
    {
      path: '/reset',
      name: 'queue-reset',
      component: QueueReset,
    },
  ],
})

export default router
