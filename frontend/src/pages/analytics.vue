<template>
  <div class="analytics">
    <h2 class="title">Аналитика</h2>
    <p v-if="link" class="subtitle">{{ link.original_url }}</p>

    <div v-if="loading" class="state muted">Загрузка...</div>
    <p v-else-if="error" class="state error">{{ error }}</p>

    <div v-else class="grid">
      <div class="label">Переходы:</div>
      <div class="total">Всего: <strong>{{ transitions.length }}</strong></div>

      <div class="list">
        <div v-for="t in transitions" :key="t.id" class="item">
          <div class="main">
            <div class="text">{{ formatRow(t) }}</div>
            <button class="ghost" @click="toggle(t)">{{ t._open ? 'Скрыть' : 'Подробнее' }}</button>
          </div>
          <div v-if="t._open" class="details">
            <div>Откуда: <strong>{{ t.referer || '—' }}</strong></div>
            <div>Браузер: <strong>{{ t.browser || '—' }}</strong></div>
            <div>OS: <strong>{{ t.os || '—' }}</strong></div>
            <div>User-agent: <strong class="ua">{{ t.user_agent || '—' }}</strong></div>
          </div>
        </div>
        <div v-if="transitions.length === 0" class="muted dots">...</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const id = Number(route.params.id)

const loading = ref(false)
const error = ref(null)
const link = ref(null)
const transitions = ref([])

function formatRow(t) {
  const country = t.country || 'Unknown'
  const city = t.city || ''
  const dt = new Date(t.created_at)
  const dateStr = dt.toLocaleString()
  return `${country}${city ? ', ' + city : ''} ${dateStr}`
}

function toggle(t) { t._open = !t._open }

async function fetchData() {
  loading.value = true
  error.value = null
  try {
    const [linkRes, tranRes] = await Promise.all([
      fetch(`/api/v1/links/${id}`, { credentials: 'include' }),
      fetch(`/api/v1/links/${id}/transitions`, { credentials: 'include' }),
    ])
    if (!linkRes.ok) {
      if (linkRes.status === 401) { router.push('/login'); return }
      throw new Error('Не удалось получить ссылку')
    }
    if (!tranRes.ok) {
      if (tranRes.status === 401) { router.push('/login'); return }
      throw new Error('Не удалось получить аналитику')
    }
    link.value = await linkRes.json()
    const data = await tranRes.json()
    transitions.value = Array.isArray(data?.transitions) ? data.transitions : []
  } catch (e) {
    error.value = e.message || 'Ошибка'
  } finally {
    loading.value = false
  }
}

onMounted(fetchData)
</script>

<style scoped>
.analytics { max-width: 980px; margin: 60px auto; padding: 0 16px; }
.title { margin: 0; color: #2b3b45; font-weight: 600; }
.subtitle { color: #6a7a86; margin: 8px 0 24px; }
.state { margin: 24px 0; }
.muted { color: #8b98a1; }
.error { color: #c54747; }
.grid { display: grid; grid-template-columns: 1fr auto; align-items: start; gap: 10px 16px; }
.label { font-weight: 600; color: #2b3b45; }
.total { justify-self: end; color: #2b3b45; }

.list { grid-column: 1 / -1; display: grid; gap: 12px; }
.item { background: #f6f7f9; border: 1px solid #eceff3; border-radius: 10px; padding: 10px; }
.main { display: flex; justify-content: space-between; align-items: center; gap: 12px; }
.text { font-weight: 600; color: #2b3b45; }
.ghost { background: #f3f5f7; color: #3a4a56; border: 1px solid #e3e9ee; border-radius: 8px; padding: 8px 12px; cursor: pointer; }
.ghost:hover { background: #eef2f7; transform: translateY(-1px); }
.details { margin-top: 10px; padding: 12px; border-left: 4px solid #6b97a7; background: #f2f4f6; border-radius: 8px; color: #1f2a33; }
.ua { word-break: break-all; }
.dots { text-align: center; }
</style>

