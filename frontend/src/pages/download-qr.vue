<template>
  <div class="download">
    <h2 class="title">Скачать <span class="accent">qr</span>-код</h2>
    <p v-if="link" class="subtitle">{{ link.original_url }}</p>

    <div v-if="loading" class="state muted">Загрузка...</div>
    <p v-else-if="error" class="state error">{{ error }}</p>

    <div v-else class="center">
      <div class="canvas">
        <div v-if="qrLoading || !qrSrc" class="spinner"></div>
        <img v-else :src="qrSrc" alt="QR preview" />
      </div>

      <div class="controls">
        <div class="label">Скачать</div>
        <div class="buttons">
          <button class="ghost" @click="download('svg')">SVG</button>
          <button class="primary" @click="download('png')">PNG</button>
          <button class="ghost" @click="download('pdf')">PDF</button>
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
const id = Number(route.params.id)

const loading = ref(false)
const error = ref(null)
const link = ref(null)
const qrSrc = ref('')
const qrLoading = ref(false)
let lastObjectUrl = null

async function fetchLink() {
  loading.value = true
  error.value = null
  try {
    const res = await fetch(`/api/v1/links/${id}`, { credentials: 'include' })
    if (!res.ok) {
      if (res.status === 401) { router.push('/login'); return }
      throw new Error('Не удалось получить ссылку')
    }
    link.value = await res.json()
    await generatePreview()
  } catch (e) {
    error.value = e.message || 'Ошибка'
  } finally {
    loading.value = false
  }
}

async function generatePreview() {
  if (!link.value) return
  qrLoading.value = true
  try {
    const res = await fetch('/api/v1/qrcode', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({
        url: link.value.original_url,
        color: link.value.color,
        background: link.value.background,
        smoothing: link.value.smoothing ?? 0,
      }),
    })
    if (!res.ok) {
      if (res.status === 401) { router.push('/login'); return }
      throw new Error('Не удалось сгенерировать превью')
    }
    const blob = await res.blob()
    if (lastObjectUrl) URL.revokeObjectURL(lastObjectUrl)
    lastObjectUrl = URL.createObjectURL(blob)
    qrSrc.value = lastObjectUrl
  } finally {
    qrLoading.value = false
  }
}

function download(type) {
  const t = String(type || 'png').toLowerCase()
  window.open(`/api/v1/links/${id}/download?type=${encodeURIComponent(t)}`, '_blank')
}

onMounted(fetchLink)
</script>

<style scoped>
.download { max-width: 980px; margin: 60px auto; padding: 0 16px; text-align: center; }
.title { 
  margin: 0; 
  color: var(--text-primary); 
  font-weight: 600; 
}
.accent { color: var(--accent-color); }
.subtitle { 
  color: var(--text-secondary); 
  margin: 8px 0 24px; 
}
.state { margin: 24px 0; }
.muted { color: var(--text-muted); }
.error { color: var(--error-color); }

.center { display: grid; justify-items: center; gap: 26px; }
.canvas { 
  width: 360px; 
  height: 360px; 
  border-radius: 28px; 
  border: 1px solid black; 
  box-shadow: 0 18px 44px var(--shadow-dark); 
  display: flex; 
  align-items: center; 
  justify-content: center; 
  background: #fff; 
}
:root.dark-theme .canvas {
  background: #000; 
}
.canvas img { 
  width: 320px; 
  height: 320px; 
  border-radius: 14px; 
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

.controls { display: grid; gap: 10px; }
.label { 
  color: var(--text-primary); 
  font-weight: 600; 
}
.buttons { display: flex; gap: 14px; }
.primary, .ghost { 
  border: 1px solid transparent; 
  border-radius: 10px; 
  padding: 10px 22px; 
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
</style>