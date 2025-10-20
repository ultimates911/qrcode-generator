import { createRouter, createWebHistory } from 'vue-router'
import Login from '../pages/login.vue'

const routes = [
  {path: '/', redirect: '/login'},
  { path: '/login', component: Login }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
