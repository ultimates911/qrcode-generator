<template>
  <div class="links-page">
    <div class="header">
      <h2>Мои ссылки</h2>
      <div class="actions-header">
        <input class="search" type="text" v-model="search" placeholder="Поиск по имени" />
        <button class="add-btn" @click="onAdd">Добавить</button>
      </div>
    </div>

    <div v-if="loading" class="state muted">Загрузка...</div>
    <div v-else-if="error" class="state error">{{ error }}</div>

    <div v-else>
      <div v-if="items.length === 0" class="empty">
        <p class="muted">Здесь пока ничe нет. <a href="#" @click.prevent="onAdd">Хотите начать?</a></p>
      </div>

      <div v-else>
        <table class="table">
          <thead>
            <tr>
              <th>Название</th>
              <th @click="toggleSort('transitions')" class="sortable">
                Количество переходов
                <span class="sort" v-if="sortBy==='transitions'">{{ sortOrder==='asc' ? '▲' : '▼' }}</span>
              </th>
              <th>Ссылка</th>
              <th @click="toggleSort('created_at')" class="sortable">
                Дата и время создания
                <span class="sort" v-if="sortBy==='created_at'">{{ sortOrder==='asc' ? '▲' : '▼' }}</span>
              </th>
              <th class="_w"></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="it in items" :key="it.id">
              <td class="name">{{ it.name || 'Без имени' }}</td>
              <td class="cnt">{{ it.transitions_count ?? it.transitions ?? 0 }}</td>
              <td class="url-cell"><input class="url" :value="it.original_url" readonly /></td>
              <td class="date">{{ formatDate(it.created_at) }}</td>
              <td class="actions">
                <button class="icon" title="Аналитика" @click="openAnalytics(it)"><span class="i">A</span></button>
                <button class="icon" title="Просмотр QR" @click="previewQR(it)"><span class="i">QR</span></button>
                <button type="button" class="icon" title="Настройки QR" @click.stop="editQR(it)"><span class="i">⚙</span></button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

const items = ref([])
const loading = ref(false)
const error = ref(null)
const router = useRouter()
const search = ref('')
const sortBy = ref('created_at')
const sortOrder = ref('desc')

async function fetchLinks() {
  loading.value = true
  error.value = null
  try {
    const params = new URLSearchParams()
    if (search.value) params.set('search', search.value)
    if (sortBy.value) params.set('sort_by', sortBy.value)
    if (sortOrder.value) params.set('order', sortOrder.value)
    const qs = params.toString() ? `?${params.toString()}` : ''
    const res = await fetch(`/api/v1/links${qs}`, { credentials: 'include' })
    if (!res.ok) {
      if (res.status === 401) {
        router.push('/login')
        return
      }
      let msg = 'Не удалось загрузить ссылки'
      try { const j = await res.json(); if (j?.error) msg = j.error } catch {}
      throw new Error(msg)
    }
    const data = await res.json()
    items.value = Array.isArray(data?.links) ? data.links : []
  } catch (e) {
    error.value = e.message || 'Ошибка'
  } finally {
    loading.value = false
  }
}

function onAdd() { router.push('/links/new') }

function openAnalytics(it) { router.push(`/links/${it.id}/analytics`) }

function previewQR(it) { router.push(`/links/${it.id}/download`) }

function editQR(it) { router.push(`/links/${it.id}/edit`) }

onMounted(fetchLinks)

let debounceId = null
watch(search, () => {
  clearTimeout(debounceId)
  debounceId = setTimeout(() => {
    fetchLinks()
  }, 300)
})

function toggleSort(by) {
  if (sortBy.value === by) {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortBy.value = by
    sortOrder.value = by === 'created_at' ? 'desc' : 'desc'
  }
  fetchLinks()
}

function formatDate(s) {
  try {
    const d = new Date(s)
    if (Number.isNaN(d.getTime())) return '—'
    return d.toLocaleString()
  } catch { return '—' }
}
</script>

<style scoped>
.links-page {
  max-width: 900px;
  margin: 60px auto;
  padding: 0 16px;
}
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}
.actions-header { display: flex; gap: 10px; align-items: center; }
.search {
  width: 240px;
  padding: 10px 12px;
  background: #f3f5f7;
  border: 1px solid #edf0f2;
  border-radius: 8px;
  color: #2e2e2e;
}
h2 {
  margin: 0;
  color: #4b6576;
  font-weight: 600;
}
.add-btn {
  background: #6c8c8c;
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 10px 18px;
  cursor: pointer;
  transition: transform .15s ease, box-shadow .15s ease, background .2s ease;
}
.add-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 14px rgba(75, 121, 161, .25);
  background: linear-gradient(135deg, #6c8c8c 0%, #6c8c8c 100%);
}
.state { margin: 24px 0;}
.muted { color: #666; }
.error { color: #c54747; }
.empty { margin-top: 40px; color: #6c8c8c; }

.empty a {
  color: #7c81a2;
  text-decoration: aquamarine;
  transition: color 0.2s ease;
}

.empty a:hover {
  color: #5a7a7a; 
  text-decoration: underline;
}

.table { width: 100%; border-collapse: collapse; }
.table th, .table td { padding: 10px 12px; border-bottom: 1px solid #edf0f2; text-align: left; }
.table thead th { font-weight: 600; color: #2e3a44; background: #f8fafb; }
.table tbody tr:hover { background: #f9fbfd; }
.sortable { cursor: pointer; user-select: none; }
.sort { margin-left: 6px; color: #3a6a95; }
.name { font-weight: 600; color: #2e3a44; max-width: 280px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.cnt { width: 140px; }
.url {
  width: 100%;
  padding: 12px 14px;
  background: #f3f5f7;
  border: 1px solid #edf0f2;
  border-radius: 8px;
  color: #2e2e2e;
  transition: border-color .15s ease, box-shadow .15s ease;
}
.table tbody tr:hover .url {
  border-color: #dfe6eb;
  box-shadow: inset 0 0 0 1px #e6edf2;
}
.actions { display: flex; gap: 8px; justify-content: flex-end; }
.url-cell { min-width: 300px; }
th._w, td._w { width: 140px; }
.icon {
  width: 40px;
  height: 40px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: #f3f5f7;
  border: 1px solid #edf0f2;
  border-radius: 8px;
  cursor: pointer;
  transition: transform .12s ease, background .15s ease, box-shadow .15s ease, border-color .15s ease;
}
.icon:hover {
  transform: translateY(-1px);
  background: linear-gradient(135deg, #f7faff 0%, #eef5ff 100%);
  border-color: #d5e2f3;
  box-shadow: 0 4px 10px rgba(58,106,149,.15);
}
.i { display: inline-flex; }
.i svg { width: 20px; height: 20px; color: #3a6a95; }
</style>

