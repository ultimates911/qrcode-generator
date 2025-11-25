<template>
  <div class="edit-qr">
    <h2 class="title">Настройка <span class="accent">qr</span>-кода</h2>

    <div class="layout">
      <div class="preview">
        <div class="canvas">
          <div v-if="qrLoading || !qrSrc" class="spinner"></div>
          <img v-else :src="qrSrc" alt="QR preview" />
        </div>
      </div>

      <div class="panel">
        <div class="field">
          <label>Ссылка</label>
          <input v-model.trim="url" type="url" placeholder="https://example.com/page/1" @input="queueGenerate" />
        </div>

        <div class="field">
          <label>Цвет</label>
          <div class="swatches">
            <button v-for="c in fgPalette" :key="c" class="swatch" :style="{ background: '#'+c }" :class="{ active: color===c }" @click="setColor(c)"></button>
          </div>
        </div>

        <div class="field">
          <label>Фон</label>
          <div class="swatches">
            <button v-for="c in bgPalette" :key="c" class="swatch" :style="{ background: '#'+c }" :class="{ active: background===c }" @click="setBackground(c)"></button>
          </div>
        </div>

        <div class="field">
          <label>Уровень сглаживания</label>
          <div class="swatches smoothing">
            <button v-for="s in smoothingPreset" :key="s" class="pill" :class="{ active: smoothing===s }" @click="setSmoothing(s)">{{ s.toFixed(2) }}</button>
          </div>
        </div>

        <p v-if="error" class="error">{{ error }}</p>

        <div class="actions">
          <button class="ghost" @click="router.push('/links')">Отмена</button>
          <button class="primary" :disabled="saving" @click="save">{{ saving ? 'Сохранение...' : 'Сохранить' }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const linkId = Number(route.params.id)

const url = ref('')
const hash = ref('')
const color = ref('000000')
const background = ref('FFFFFF')
const smoothing = ref(0.0)

const fgPalette = ['0EA5E9','2563EB','22C55E','EF4444','F59E0B','C8A696','837DA2','000000']
const bgPalette = ['FFFFFF','F8FAFC','E2E8F0','F1F5F9','FFF7ED','FDF2F8','ECFDF5','FAFAFA']
const smoothingPreset = [0.00, 0.05, 0.10, 0.20, 0.30, 0.40, 0.50]

const error = ref(null)
const qrSrc = ref('')
const qrLoading = ref(false)
const saving = ref(false)
let pendingTimer = null
let lastObjectUrl = null

function queueGenerate() {
  if (pendingTimer) clearTimeout(pendingTimer)
  pendingTimer = setTimeout(generateQR, 250)
}

function setColor(c) { color.value = c; queueGenerate() }
function setBackground(c) { background.value = c; queueGenerate() }
function setSmoothing(s) { smoothing.value = s; queueGenerate() }

async function fetchLink() {
  error.value = null
  try {
    const res = await fetch(`/api/v1/links/${linkId}`, { credentials: 'include' })
    if (!res.ok) {
      if (res.status === 401) { router.push('/login'); return }
      let msg = 'Не удалось получить данные ссылки'
      try { const j = await res.json(); if (j?.error) msg = j.error } catch {}
      throw new Error(msg)
    }
    const data = await res.json()
    url.value = data.original_url
    hash.value = data.hash || ''
    color.value = data.color || color.value
    background.value = data.background || background.value
    smoothing.value = typeof data.smoothing === 'number' ? data.smoothing : 0.0
    await generateQR()
  } catch (e) {
    error.value = e.message || 'Ошибка'
  }
}

async function generateQR() {
  if (!hash.value) return
  qrLoading.value = true
  error.value = null
  try {
    const redirectUrl = `${window.location.protocol}//${window.location.hostname}:8080/redirect/${hash.value}`
    const res = await fetch('/api/v1/qrcode', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ url: redirectUrl, color: color.value, background: background.value, smoothing: smoothing.value })
    })
    if (!res.ok) {
      if (res.status === 401) { router.push('/login'); return }
      throw new Error('Не удалось сгенерировать QR-код')
    }
    const blob = await res.blob()
    if (lastObjectUrl) URL.revokeObjectURL(lastObjectUrl)
    lastObjectUrl = URL.createObjectURL(blob)
    qrSrc.value = lastObjectUrl
  } catch (e) {
    error.value = e.message || 'Ошибка'
  } finally {
    qrLoading.value = false
  }
}

async function save() {
  saving.value = true
  error.value = null
  try {
    const res = await fetch(`/api/v1/links/${linkId}`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ original_url: url.value, color: color.value, background: background.value, smoothing: smoothing.value })
    })
    if (!res.ok) {
      if (res.status === 401) { router.push('/login'); return }
      let msg = 'Не удалось сохранить изменения'
      try { const j = await res.json(); if (j?.error) msg = j.error } catch {}
      throw new Error(msg)
    }
    router.push('/links')
  } catch (e) {
    error.value = e.message || 'Ошибка'
  } finally {
    saving.value = false
  }
}

onMounted(fetchLink)
</script>

<style scoped>
.edit-qr { max-width: 1100px; margin: 60px auto; padding: 0 16px; }
.title { 
  margin: 0 0 24px; 
  color: var(--text-primary); 
  font-weight: 600; 
}
.accent { color: var(--accent-color); }
.layout { display: grid; grid-template-columns: 1fr 1fr; gap: 28px; align-items: start; }

.preview .canvas { 
  width: 360px; 
  height: 360px; 
  border-radius: 22px; 
  border: 1px solid var(--border-color); 
  box-shadow: 0 12px 30px var(--shadow-dark); 
  display: flex; 
  align-items: center; 
  justify-content: center; 
  background: #fff; 
}
:root.dark-theme .preview .canvas {
  background: #000; 
}
.preview img { 
  width: 320px; 
  height: 320px; 
  border-radius: 12px; 
  box-shadow: 0 4px 16px var(--shadow-dark); 
}

.spinner { 
  width: 28px; 
  height: 28px; 
  border: 3px solid var(--border-color); 
  border-top-color: var(--accent-color); 
  border-radius: 50%; 
  animation: spin 1s linear infinite; 
}
@keyframes spin { to { transform: rotate(360deg) } }

.panel { display: grid; gap: 18px; }
.field { display: grid; gap: 8px; }
label { 
  color: var(--text-secondary); 
  font-size: 14px; 
}
input { 
  padding: 10px 12px; 
  border: 1px solid var(--border-color); 
  border-radius: 8px; 
  outline: none;
  background: var(--bg-tertiary);
  color: var(--text-primary);
}
input:focus { 
  border-color: var(--accent-color); 
  box-shadow: 0 0 0 3px var(--shadow-light); 
}

.swatches { display: grid; grid-template-columns: repeat(6, 36px); gap: 10px; }
.swatch { 
  width: 36px; 
  height: 36px; 
  border-radius: 8px; 
  border: 2px solid var(--border-color); 
  cursor: pointer; 
  transition: transform .12s ease, box-shadow .15s ease, border-color .15s ease; 
}
.swatch:hover { 
  transform: translateY(-1px); 
  box-shadow: 0 4px 10px var(--shadow-light); 
}
.swatch.active { 
  border-color: var(--accent-color); 
  box-shadow: 0 0 0 3px var(--shadow-light); 
}

.smoothing { grid-template-columns: repeat(7, auto); }
.pill { 
  padding: 6px 10px; 
  border-radius: 8px; 
  border: 1px solid var(--border-color); 
  background: var(--bg-tertiary); 
  color: var(--text-primary);
  cursor: pointer; 
  transition: transform .12s ease, box-shadow .15s ease, background .2s ease; 
}
.pill.active, .pill:hover { 
  background: var(--accent-color); 
  color: white;
  box-shadow: 0 4px 10px var(--shadow-light); 
}

.actions { display: flex; gap: 10px; }
.primary, .ghost { 
  border: 1px solid transparent; 
  border-radius: 8px; 
  padding: 10px 16px; 
  cursor: pointer; 
  transition: all 0.2s ease;
}
.primary { 
  background: var(--accent-color); 
  color: #fff; 
}
.primary:hover { 
  transform: translateY(-1px); 
  box-shadow: 0 6px 14px var(--shadow-light); 
  background: var(--accent-hover);
}
.ghost { 
  background: var(--bg-tertiary); 
  color: var(--text-primary); 
  border-color: var(--border-color); 
}
.ghost:hover { 
  background: var(--bg-secondary); 
  transform: translateY(-1px); 
  box-shadow: 0 4px 10px var(--shadow-light);
}
.error { color: var(--error-color); }
</style>