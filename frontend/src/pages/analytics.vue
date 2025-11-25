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
.title { 
  margin: 0; 
  color: var(--text-primary); 
  font-weight: 600; 
}
.subtitle { 
  color: var(--text-secondary); 
  margin: 8px 0 24px; 
}
.state { margin: 24px 0; }
.muted { color: var(--text-muted); }
.error { color: var(--error-color); }
.grid { display: grid; grid-template-columns: 1fr auto; align-items: start; gap: 10px 16px; }
.label { 
  font-weight: 600; 
  color: var(--text-primary); 
}
.total { 
  justify-self: end; 
  color: var(--text-primary); 
}

.list { grid-column: 1 / -1; display: grid; gap: 12px; }
.item { 
  background: var(--bg-tertiary); 
  border: 1px solid var(--border-color); 
  border-radius: 10px; 
  padding: 10px; 
}
.main { display: flex; justify-content: space-between; align-items: center; gap: 12px; }
.text { 
  font-weight: 600; 
  color: var(--text-primary); 
}
.ghost { 
  background: var(--bg-secondary); 
  color: var(--text-primary); 
  border: 1px solid var(--border-color); 
  border-radius: 8px; 
  padding: 8px 12px; 
  cursor: pointer; 
  transition: all 0.2s ease;
}
.ghost:hover { 
  background: var(--bg-tertiary); 
  transform: translateY(-1px); 
  box-shadow: 0 2px 8px var(--shadow-light);
}
.details { 
  margin-top: 10px; 
  padding: 12px; 
  border-left: 4px solid var(--accent-color); 
  background: var(--bg-secondary); 
  border-radius: 8px; 
  color: var(--text-primary); 
}
.ua { word-break: break-all; }
.dots { text-align: center; }
</style>
