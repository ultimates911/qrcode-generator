<template>
  <div class="new-link">
    <h2 class="title">Добавление ссылки</h2>

    <form class="form" @submit.prevent="submit">
      <input class="input" v-model.trim="url" type="url" required placeholder="Ссылка" />

      <p v-if="error" class="error">{{ error }}</p>

      <button type="submit" class="submit" :disabled="loading">
        {{ loading ? 'Отправка...' : 'Отправить' }}
      </button>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const url = ref('')
const error = ref(null)
const loading = ref(false)

async function submit() {
  error.value = null
  if (!url.value) return
  loading.value = true
  try {
    const res = await fetch('/api/v1/links/create', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ original_url: url.value }),
    })
    if (!res.ok) {
      if (res.status === 401) {
        router.push('/login')
        return
      }
      let msg = 'Не удалось создать ссылку'
      try { const j = await res.json(); if (j?.error) msg = j.error } catch {}
      throw new Error(msg)
    }
    router.push('/links')
  } catch (e) {
    error.value = e.message || 'Ошибка'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.new-link {
  max-width: 920px;
  margin: 80px auto;
  padding: 0 16px;
  text-align: center;
}
.title { margin: 0 0 40px; color: #4b6576; font-weight: 600; }

.form { display: grid; justify-items: center; gap: 18px; }
.input {
  width: min(600px, 100%);
  padding: 12px 14px;
  border: 1px solid #cddfe8;
  border-radius: 8px;
  background: #ffffff;
  outline: none;
  transition: border-color .15s ease, box-shadow .15s ease;
}
.input:focus { border-color: #6ea3c4; box-shadow: 0 0 0 3px rgba(85,143,179,.2); }
.error { color: #c54747; margin: 4px 0 0; }

.submit {
  background: #4b79a1;
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 10px 18px;
  cursor: pointer;
  transition: transform .12s ease, box-shadow .15s ease, background .2s ease;
}
.submit:hover { transform: translateY(-1px); box-shadow: 0 6px 14px rgba(75,121,161,.25); }
</style>

