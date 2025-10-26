import { createRouter, createWebHistory } from 'vue-router'
import Login from '../pages/login.vue'
import Links from '../pages/links.vue'
import NewLink from '../pages/new-link.vue'

const routes = [
  { path: '/', redirect: '/links' },
  { path: '/login', component: Login },
  { path: '/links', component: Links },
  { path: '/links/new', component: NewLink }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
