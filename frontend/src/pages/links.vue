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
                <button class="icon" title="Аналитика" @click="openAnalytics(it)">
                  <span class="i analytics-icon">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                      <line x1="3" y1="3" x2="3" y2="21"></line> <!-- Ось Y -->
                      <line x1="3" y1="21" x2="21" y2="21"></line> <!-- Ось X -->
                      <rect x="6" y="13" width="3" height="8" rx="1"></rect>
                      <rect x="11" y="8" width="3" height="13" rx="1"></rect>
                      <rect x="16" y="3" width="3" height="18" rx="1"></rect>
                    </svg>
                  </span>
                </button>
                <button class="icon" title="Просмотр QR" @click="previewQR(it)">
                  <span class="i qr-icon">

                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
               
                      <path d="M4 4h6v6H4z" stroke-width="2"></path>
                      <path d="M4 4v6h6" stroke-width="2" fill="none"></path>
                      
                      <path d="M14 4h6v6h-6z" stroke-width="2"></path>
                      <path d="M14 4v6h6" stroke-width="2" fill="none"></path>
                      
                      <path d="M4 14h6v6H4z" stroke-width="2"></path>
                      <path d="M4 14v6h6" stroke-width="2" fill="none"></path>
                      
                      <path d="M14 14h6v6h-6z" stroke-width="2"></path>
                      <path d="M14 14v6h6" stroke-width="2" fill="none"></path>
                      
                      <rect x="10" y="10" width="1.5" height="1.5" fill="currentColor"></rect>
                      <rect x="12" y="12" width="1.5" height="1.5" fill="currentColor"></rect>
                      <rect x="10" y="12" width="1.5" height="1.5" fill="currentColor"></rect>
                      <rect x="12" y="10" width="1.5" height="1.5" fill="currentColor"></rect>
                      
                      <rect x="18" y="10" width="1" height="1" fill="currentColor"></rect>
                      <rect x="10" y="18" width="1" height="1" fill="currentColor"></rect>
                      <rect x="18" y="18" width="1" height="1" fill="currentColor"></rect>
                    </svg>
                  </span>
                </button>
                <button type="button" class="icon" title="Настройки QR" @click.stop="editQR(it)">
                  <span class="i settings-icon">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                      <circle cx="12" cy="12" r="3"></circle>
                      <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
                    </svg>
                  </span>
                </button>
                <button type="button" class="icon delete" title="Удалить ссылку" @click.stop="showDeleteModal(it)">
                  <span class="i delete-icon">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                      <polyline points="3 6 5 6 21 6"></polyline>
                      <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                      <line x1="10" y1="11" x2="10" y2="17"></line>
                      <line x1="14" y1="11" x2="14" y2="17"></line>
                    </svg>
                  </span>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal">
        <div class="modal-content">
          <h3>Удаление ссылки</h3>
          <p>Вы точно хотите удалить ссылку "{{ linkToDelete?.name || 'Без имени' }}"?</p>
          <div class="modal-actions">
            <button class="cancel-btn" @click="closeModal">Отмена</button>
            <button class="confirm-delete-btn" @click="confirmDelete">Удалить</button>
          </div>
        </div>
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
const showModal = ref(false)
const linkToDelete = ref(null)

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

function showDeleteModal(link) {
  linkToDelete.value = link
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  linkToDelete.value = null
}

async function confirmDelete() {
  if (!linkToDelete.value) return
  
  try {
    const res = await fetch(`/api/v1/links/${linkToDelete.value.id}`, {
      method: 'DELETE',
      credentials: 'include'
    })

    if (!res.ok) {
      let msg = 'Не удалось удалить ссылку'
      try { 
        const j = await res.json()
        if (j?.error) msg = j.error 
      } catch {}
      throw new Error(msg)
    }

    items.value = items.value.filter(item => item.id !== linkToDelete.value.id)
    closeModal()
  } catch (e) {
    alert(e.message || 'Ошибка при удалении')
    closeModal()
  }
}

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
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  color: var(--text-primary);
}
h2 {
  margin: 0;
  color: var(--text-muted);
  font-weight: 600;
}
.add-btn {
  background: var(--accent-color);
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 10px 18px;
  cursor: pointer;
  transition: transform .15s ease, box-shadow .15s ease, background .2s ease;
}
.add-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 14px var(--shadow-light);
  background: var(--accent-hover);
}
.state { margin: 24px 0;}
.muted { color: var(--text-secondary); }
.error { color: var(--error-color); }
.empty { margin-top: 40px; color: var(--text-muted); }

.empty a {
  color: var(--accent-color);
  text-decoration: none;
  transition: color 0.2s ease;
}

.empty a:hover {
  color: var(--accent-hover); 
  text-decoration: underline;
}

.table { width: 100%; border-collapse: collapse; }
.table th, .table td { padding: 10px 12px; border-bottom: 1px solid var(--border-color); text-align: left; }
.table thead th { font-weight: 600; color: var(--text-primary); background: var(--bg-secondary); }
.table tbody tr:hover { background: var(--bg-secondary); }
.sortable { cursor: pointer; user-select: none; }
.sort { margin-left: 6px; color: var(--text-muted); }
.name { font-weight: 600; color: var(--text-primary); max-width: 280px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.cnt { width: 140px; }
.url {
  width: 100%;
  padding: 12px 14px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  color: var(--text-primary);
  transition: border-color .15s ease, box-shadow .15s ease;
}
.table tbody tr:hover .url {
  border-color: var(--border-hover);
  box-shadow: inset 0 0 0 1px var(--border-hover);
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
  background: var(--icon-bg);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  cursor: pointer;
  transition: transform .12s ease, background .15s ease, box-shadow .15s ease, border-color .15s ease;
  font-size: 16px; 
  color: var(--icon-color);
}

.icon:hover {
  transform: translateY(-1px);
  background: var(--icon-hover-bg);
  border-color: var(--icon-border-hover);
  box-shadow: 0 4px 10px var(--shadow-light);
}

.icon.delete:hover {
  background: linear-gradient(135deg, #ffeaea 0%, #ffdbdb 100%);
  border-color: #f3c5c5;
  box-shadow: 0 4px 10px rgba(149, 58, 58, 0.25);
  color: #c54747; 
}

:root.dark-theme .icon.delete:hover {
  background: #c54747;
  border-color: #c54747;
  color: white; 
}

.i { 
  display: inline-flex; 
  align-items: center;
  justify-content: center;
}

.i svg { 
  width: 20px; 
  height: 20px; 
  color: var(--icon-color); 
  transition: color 0.15s ease;
}

.analytics-icon svg {
  width: 25px; 
  height: 25px; 
}

.qr-icon svg {
  width: 25px; 
  height: 25px; 
  stroke-width: 1%; 
}

.settings-icon svg {
  width: 20px;
  height: 20px;
}

.delete-icon svg {
  width: 20px;
  height: 20px;
}

.icon:hover .i svg {
  color: var(--icon-hover-color);
}


.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--modal-overlay);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: var(--bg-secondary);
  border-radius: 12px;
  padding: 0;
  width: 450px;
  max-width: 90vw;
  box-shadow: 0 20px 40px var(--shadow-dark);
  border: 1px solid var(--border-color);
}

.modal-content {
  padding: 32px;
}

.modal-content h3 {
  margin: 0 0 16px 0;
  color: var(--text-primary);
  font-weight: 600;
  font-size: 20px;
}

.modal-content p {
  margin: 0 0 32px 0;
  color: var(--text-secondary);
  line-height: 1.5;
  font-size: 15px;
}

.modal-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

.cancel-btn, .confirm-delete-btn {
  padding: 12px 24px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: transform .15s ease, box-shadow .15s ease, background .2s ease;
  min-width: 100px;
}

.cancel-btn {
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  color: var(--text-primary);
}

.cancel-btn:hover {
  transform: translateY(-1px);
  background: linear-gradient(135deg, var(--bg-secondary) 0%, var(--bg-tertiary) 100%);
  border-color: var(--border-hover);
  box-shadow: 0 4px 10px var(--shadow-light);
}

.confirm-delete-btn {
  background: transparent;
  color: var(--accent-color);
  border: 1px solid var(--accent-color);
}

.confirm-delete-btn:hover {
  transform: translateY(-1px);
  background: var(--error-color);
  color: white;
  border-color: var(--error-color);
  box-shadow: 0 4px 10px rgba(149, 58, 58, 0.25);
}
</style>