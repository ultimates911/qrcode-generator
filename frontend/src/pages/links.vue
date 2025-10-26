<template>
  <div class="links-page">
    <div class="header">
      <h2>Мои ссылки</h2>
      <button class="add-btn" @click="onAdd">Добавить</button>
    </div>

    <div v-if="loading" class="state muted">Загрузка...</div>
    <div v-else-if="error" class="state error">{{ error }}</div>

    <div v-else>
      <div v-if="items.length === 0" class="empty">
        <p class="muted">Здесь пока ничe нет. <a href="#" @click.prevent="onAdd">Хотите начать?</a></p>
      </div>

      <div v-else class="list">
        <div v-for="it in items" :key="it.id" class="row">
          <input class="url" :value="it.original_url" readonly />
          <div class="actions">
            <button class="icon" title="Аналитика" @click="openAnalytics(it)">
              <span class="i" aria-hidden="true">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="3 17 9 11 13 15 21 7"></polyline>
                  <polyline points="21 7 15 7 15 13"></polyline>
                </svg>
              </span>
            </button>
            <button class="icon" title="Просмотр QR" @click="previewQR(it)">
              <span class="i" aria-hidden="true">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <rect x="3" y="3" width="7" height="7"></rect>
                  <rect x="14" y="3" width="7" height="7"></rect>
                  <rect x="14" y="14" width="7" height="7"></rect>
                  <path d="M3 14h7v7H3z"></path>
                  <path d="M7 7h0"></path>
                </svg>
              </span>
            </button>
            <button class="icon" title="Настройки QR" @click="editQR(it)">
              <span class="i" aria-hidden="true">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="12" cy="12" r="3"></circle>
                  <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 1 1-4 0v-.09a1.65 1.65 0 0 0-1-1.51 1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 1 1 0-4h.09a1.65 1.65 0 0 0 1.51-1 1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 1 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9c0 .66.39 1.25 1 1.51.2.09.43.14.66.14H21a2 2 0 1 1 0 4h-.09c-.23 0-.46.05-.66.14-.61.26-1 .85-1 1.51z"></path>
                </svg>
              </span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const items = ref([])
const loading = ref(false)
const error = ref(null)
const router = useRouter()

async function fetchLinks() {
  loading.value = true
  error.value = null
  try {
    const res = await fetch('/api/v1/links', { credentials: 'include' })
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

function onAdd() {
  // Заглушка: позже можно открыть модал для создания ссылки
  // Пока просто навигируем на форму если появится
  // alert('Добавление ссылки в разработке')
}

function openAnalytics(it) {
  // Навигация/заглушка под страницу аналитики
  // console.log('analytics', it)
}

function previewQR(it) {
  // Открыть скачивание PNG текущего QR как предпросмотр
  window.open(`/api/v1/links/${it.id}/download?type=png`, '_blank')
}

function editQR(it) {
  // Навигация/заглушка под страницу настройки
  // console.log('edit qr', it)
}

onMounted(fetchLinks)
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
h2 {
  margin: 0;
  color: #4b6576;
  font-weight: 600;
}
.add-btn {
  background: #4b79a1;
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
  background: linear-gradient(135deg, #4b79a1 0%, #3a6a95 100%);
}
.state { margin: 24px 0; }
.muted { color: #8b98a1; }
.error { color: #c54747; }
.empty { margin-top: 40px; }

.list { display: grid; gap: 14px; }
.row {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 12px;
  align-items: center;
  padding: 6px;
  border-radius: 10px;
  transition: background .15s ease, box-shadow .15s ease, transform .15s ease;
}
.row:hover {
  background: #f8fafb;
  box-shadow: 0 3px 10px rgba(0,0,0,.05);
}
.url {
  width: 100%;
  padding: 12px 14px;
  background: #f3f5f7;
  border: 1px solid #edf0f2;
  border-radius: 8px;
  color: #2e2e2e;
  transition: border-color .15s ease, box-shadow .15s ease;
}
.row:hover .url {
  border-color: #dfe6eb;
  box-shadow: inset 0 0 0 1px #e6edf2;
}
.actions { display: flex; gap: 8px; }
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

