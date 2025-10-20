<template>
  <div class="login-container">
    <h2>Вход</h2>
    <form @submit.prevent="handleLogin">
      <div class="input-group">
        <input type="email" v-model="email" placeholder="Почта" required />
      </div>
      <div class="input-group password-group">
        <input
          :type="showPassword ? 'text' : 'password'"
          v-model="password"
          placeholder="Пароль"
          required
        />
        <span class="toggle" @click="showPassword = !showPassword">👁</span>
      </div>

      <button type="submit" class="submit-btn">Отправить</button>

      <p v-if="error" class="error">{{ error }}</p>
    </form>

    <p class="create-account">
      <router-link to="/register">Создать аккаунт</router-link>
    </p>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import '../assets/main.css'

const email = ref('')
const password = ref('')
const showPassword = ref(false)
const error = ref(null)
const router = useRouter()

async function handleLogin() {
  error.value = null

  try {
    const res = await fetch('http://localhost:8080/api/v1/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include', // 
      body: JSON.stringify({
        email: email.value,
        password: password.value,
      }),
    })

    if (!res.ok) {
      let errText = res.statusText
      try {
        const errData = await res.json()
        errText = errData.message || errText
      } catch (e) { /* noop */ }
      error.value = `Ошибка: ${errText}`
      return
    }

    const data = await res.json()
    if (data.token) {
      document.cookie = `jwt_token=${data.token}; path=/; max-age=${7 * 24 * 60 * 60}`
    }

    router.push('/dashboard')
  } catch (err) {
    console.error(err)
    error.value = 'Ошибка сети. Проверьте подключение.'
  }
}
</script>