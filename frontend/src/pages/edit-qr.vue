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
            <button v-for="c in fgPalette" :key="c" 
                    class="swatch" 
                    :style="{ background: '#'+c }" 
                    :class="{ active: qrColor===c }" 
                    @click="setColor(c)">
            </button>
            <button class="swatch custom-color" 
                    :class="{ active: isCustomFgColor }"
                    @click="openColorPicker('color')"
                    :style="{ 
                      background: isCustomFgColor ? '#' + qrColor : 'var(--bg-tertiary)',
                      borderColor: isCustomFgColor ? '#' + qrColor : 'var(--border-color)'
                    }">
              <span v-if="!isCustomFgColor">+</span>
              <span v-else>✓</span>
            </button>
          </div>
          
          <div v-if="showColorPicker === 'color'" class="color-picker">
            <input type="color" 
                   v-model="customFgColorValue" 
                   @input="onCustomColorChange('color')" />
            <div class="color-inputs">
              <input type="text" 
                     v-model="customFgColorHex" 
                     placeholder="#000000"
                     @input="onHexInput('color')" />
              <button class="apply-btn" @click="applyCustomColor('color')">Применить</button>
              <button class="cancel-btn" @click="closeColorPicker">Отмена</button>
            </div>
          </div>
        </div>

        <div class="field">
          <label>Фон</label>
          <div class="swatches">
            <button v-for="c in bgPalette" :key="c" 
                    class="swatch" 
                    :style="{ background: '#'+c }" 
                    :class="{ active: qrBackground===c }" 
                    @click="setBackground(c)">
            </button>
            <button class="swatch custom-color" 
                    :class="{ active: isCustomBgColor }"
                    @click="openColorPicker('background')"
                    :style="{ 
                      background: isCustomBgColor ? '#' + qrBackground : 'var(--bg-tertiary)',
                      borderColor: isCustomBgColor ? '#' + qrBackground : 'var(--border-color)'
                    }">
              <span v-if="!isCustomBgColor">+</span>
              <span v-else>✓</span>
            </button>
          </div>
          
          <div v-if="showColorPicker === 'background'" class="color-picker">
            <input type="color" 
                   v-model="customBgColorValue" 
                   @input="onCustomColorChange('background')" />
            <div class="color-inputs">
              <input type="text" 
                     v-model="customBgColorHex" 
                     placeholder="#FFFFFF"
                     @input="onHexInput('background')" />
              <button class="apply-btn" @click="applyCustomColor('background')">Применить</button>
              <button class="cancel-btn" @click="closeColorPicker">Отмена</button>
            </div>
          </div>
        </div>

        <div class="field">
          <label>Уровень сглаживания</label>
          <div class="swatches smoothing">
            <button v-for="s in smoothingPreset" 
                    :key="s" 
                    class="pill" 
                    :class="{ active: smoothing===s }" 
                    @click="setSmoothing(s)">
              {{ s.toFixed(2) }}
            </button>
          </div>
        </div>

        <p v-if="error" class="error">{{ error }}</p>

        <div class="actions">
          <button class="ghost" @click="router.push('/links')">Отмена</button>
          <button class="primary" :disabled="saving" @click="save">
            {{ saving ? 'Сохранение...' : 'Сохранить' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const linkId = Number(route.params.id)

const url = ref('')
const hash = ref('')
const qrColor = ref('000000')
const qrBackground = ref('FFFFFF')
const smoothing = ref(0.0)

const fgPalette = ['0EA5E9','2563EB','22C55E','EF4444','F59E0B','C8A696','837DA2','000000']
const bgPalette = ['FFFFFF','F8FAFC','E2E8F0','F1F5F9','FFF7ED','FDF2F8','ECFDF5','FAFAFA']

const showColorPicker = ref(null)

// Переменные для палитры цвета
const customFgColorValue = ref('#000000')
const customFgColorHex = ref('')
// Переменные для палитры фона
const customBgColorValue = ref('#FFFFFF')
const customBgColorHex = ref('')

const lastCustomFgColor = ref('000000')
const lastCustomBgColor = ref('FFFFFF')

const isCustomFgColor = computed(() => {
  return !fgPalette.includes(qrColor.value)
})

const isCustomBgColor = computed(() => {
  return !bgPalette.includes(qrBackground.value)
})

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

function setColor(c) { 
  qrColor.value = c
  queueGenerate() 
}

function setBackground(c) { 
  qrBackground.value = c
  queueGenerate() 
}

function setSmoothing(s) { smoothing.value = s; queueGenerate() }

function openColorPicker(type) {
  console.log('Opening color picker for:', type) // Для отладки
  showColorPicker.value = type
  
  if (type === 'color') {
    console.log('Current qrColor:', qrColor.value)
    console.log('isCustomFgColor:', isCustomFgColor.value)
    console.log('lastCustomFgColor:', lastCustomFgColor.value)
    
    if (isCustomFgColor.value) {
      customFgColorValue.value = '#' + qrColor.value
      customFgColorHex.value = '#' + qrColor.value
    } else {
      customFgColorValue.value = '#' + lastCustomFgColor.value
      customFgColorHex.value = '#' + lastCustomFgColor.value
    }
    console.log('Set customFgColorValue to:', customFgColorValue.value)
  } else {
    console.log('Current qrBackground:', qrBackground.value)
    console.log('isCustomBgColor:', isCustomBgColor.value)
    console.log('lastCustomBgColor:', lastCustomBgColor.value)
    
    if (isCustomBgColor.value) {
      customBgColorValue.value = '#' + qrBackground.value
      customBgColorHex.value = '#' + qrBackground.value
    } else {
      customBgColorValue.value = '#' + lastCustomBgColor.value
      customBgColorHex.value = '#' + lastCustomBgColor.value
    }
    console.log('Set customBgColorValue to:', customBgColorValue.value)
  }
}

function closeColorPicker() {
  showColorPicker.value = null
}

function onCustomColorChange(type) {
  console.log('onCustomColorChange:', type)
  if (type === 'color') {
    customFgColorHex.value = customFgColorValue.value
    console.log('Updated customFgColorHex:', customFgColorHex.value)
  } else {
    customBgColorHex.value = customBgColorValue.value
    console.log('Updated customBgColorHex:', customBgColorHex.value)
  }
}

function onHexInput(type) {
  console.log('onHexInput:', type)
  if (type === 'color') {
    const hex = customFgColorHex.value.replace('#', '')
    console.log('Parsed hex:', hex)
    if (/^[0-9A-Fa-f]{6}$/.test(hex)) {
      customFgColorValue.value = '#' + hex
      console.log('Updated customFgColorValue:', customFgColorValue.value)
    }
  } else {
    const hex = customBgColorHex.value.replace('#', '')
    console.log('Parsed hex:', hex)
    if (/^[0-9A-Fa-f]{6}$/.test(hex)) {
      customBgColorValue.value = '#' + hex
      console.log('Updated customBgColorValue:', customBgColorValue.value)
    }
  }
}

function applyCustomColor(type) {
  console.log('applyCustomColor:', type)
  if (type === 'color') {
    const hex = customFgColorHex.value.replace('#', '')
    console.log('Applying color hex:', hex)
    if (/^[0-9A-Fa-f]{6}$/.test(hex)) {
      qrColor.value = hex
      lastCustomFgColor.value = hex
      console.log('Updated qrColor to:', qrColor.value)
      queueGenerate()
    } else {
      console.error('Invalid hex for color:', hex)
    }
  } else {
    const hex = customBgColorHex.value.replace('#', '')
    console.log('Applying background hex:', hex)
    if (/^[0-9A-Fa-f]{6}$/.test(hex)) {
      qrBackground.value = hex
      lastCustomBgColor.value = hex
      console.log('Updated qrBackground to:', qrBackground.value)
      queueGenerate()
    } else {
      console.error('Invalid hex for background:', hex)
    }
  }
  
  closeColorPicker()
}

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
    console.log('Fetched link data:', data)
    url.value = data.original_url
    hash.value = data.hash || ''
    
    if (data.color) {
      qrColor.value = data.color
      console.log('Set qrColor from DB:', data.color)
      if (!fgPalette.includes(data.color)) {
        lastCustomFgColor.value = data.color
        console.log('Saved as lastCustomFgColor:', data.color)
      }
    }
    
    if (data.background) {
      qrBackground.value = data.background
      console.log('Set qrBackground from DB:', data.background)
      if (!bgPalette.includes(data.background)) {
        lastCustomBgColor.value = data.background
        console.log('Saved as lastCustomBgColor:', data.background)
      }
    }
    
    smoothing.value = typeof data.smoothing === 'number' ? data.smoothing : 0.0
    await generateQR()
  } catch (e) {
    error.value = e.message || 'Ошибка'
    console.error('Error fetching link:', e)
  }
}

async function generateQR() {
  if (!hash.value) return
  qrLoading.value = true
  error.value = null
  try {
    const redirectUrl = `${window.location.protocol}//${window.location.hostname}:8080/redirect/${hash.value}`
    console.log('Generating QR with:', { 
      color: qrColor.value, 
      background: qrBackground.value 
    })
    const res = await fetch('/api/v1/qrcode', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ 
        url: redirectUrl, 
        color: qrColor.value,
        background: qrBackground.value,
        smoothing: smoothing.value 
      })
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
    console.error('Error generating QR:', e)
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
      body: JSON.stringify({ 
        original_url: url.value, 
        color: qrColor.value,
        background: qrBackground.value,
        smoothing: smoothing.value 
      })
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
    console.error('Error saving:', e)
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

.swatches { display: grid; grid-template-columns: repeat(7, 36px); gap: 10px; }
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

.swatch.custom-color {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: bold;
  color: var(--text-secondary);
  background: var(--bg-tertiary);
  border: 2px dashed var(--border-color);
}

.swatch.custom-color:hover {
  border-color: var(--accent-color);
  color: var(--accent-color);
}

.swatch.custom-color.active {
  border-color: var(--accent-color);
  border-style: solid;
  color: var(--accent-color);
}

.color-picker {
  margin-top: 12px;
  padding: 16px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  box-shadow: 0 8px 20px var(--shadow-dark);
}

.color-picker input[type="color"] {
  width: 100%;
  height: 60px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  margin-bottom: 12px;
  cursor: pointer;
  background: var(--bg-tertiary);
}

.color-inputs {
  display: flex;
  gap: 10px;
  align-items: center;
}

.color-inputs input[type="text"] {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.apply-btn, .cancel-btn {
  padding: 8px 16px;
  border-radius: 6px;
  border: 1px solid var(--border-color);
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s ease;
}

.apply-btn {
  background: var(--accent-color);
  color: white;
  border-color: var(--accent-color);
}

.apply-btn:hover {
  background: var(--accent-hover);
}

.cancel-btn {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.cancel-btn:hover {
  background: var(--bg-secondary);
}
</style>