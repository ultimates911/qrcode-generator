import { createRouter, createWebHistory } from 'vue-router'
import Login from '../pages/login.vue'
import Links from '../pages/links.vue'

const routes = [
  { path: '/', redirect: '/links' },
  { path: '/login', component: Login },
  { path: '/links', component: Links }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
