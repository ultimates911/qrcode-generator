import { createRouter, createWebHistory } from 'vue-router'
import Login from '../pages/login.vue'
import Links from '../pages/links.vue'
import NewLink from '../pages/new-link.vue'
import EditQR from '../pages/edit-qr.vue'

const routes = [
  { path: '/', redirect: '/links' },
  { path: '/login', component: Login },
  { path: '/links', component: Links },
  { path: '/links/new', component: NewLink }
  ,{ path: '/links/:id/edit', component: EditQR }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
