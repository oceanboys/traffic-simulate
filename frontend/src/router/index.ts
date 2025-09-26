import { createRouter, createWebHistory } from 'vue-router'
import Dashboard from '@/views/Dashboard.vue'
import Roads from '@/views/Roads.vue'
import GPS from '@/views/GPS.vue'
import Alerts from '@/views/Alerts.vue'

const routes = [
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/dashboard',
    name: 'dashboard',
    component: Dashboard
  },
  {
    path: '/roads',
    name: 'roads',
    component: Roads
  },
  {
    path: '/gps',
    name: 'gps',
    component: GPS
  },
  {
    path: '/alerts',
    name: 'alerts',
    component: Alerts
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
