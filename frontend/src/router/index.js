import { createRouter, createWebHistory } from 'vue-router'
import Login from '../pages/login.vue'
import Register from '../pages/register.vue'
import Links from '../pages/links.vue'
import NewLink from '../pages/new-link.vue'
import EditQR from '../pages/edit-qr.vue'
import Analytics from '../pages/analytics.vue'
import DownloadQR from '../pages/download-qr.vue'

const routes = [
  { path: '/', redirect: '/links' },
  { path: '/login', component: Login },
  { path: '/register', component: Register },
  { path: '/links', component: Links },
  { path: '/links/new', component: NewLink }
  ,{ path: '/links/:id/edit', component: EditQR }
  ,{ path: '/links/:id/analytics', component: Analytics }
  ,{ path: '/links/:id/download', component: DownloadQR }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
