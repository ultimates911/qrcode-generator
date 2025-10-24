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

      <button type="submit" class="submit-btn" :disabled="loading">
        {{ loading ? 'Вход...' : 'Войти' }}
      </button>

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
const loading = ref(false)
const router = useRouter()

function getErrorMessage(status, errorData = null) {
  switch (status) {
    case 400:
      if (errorData?.message) {
        return errorData.message
      }
      return 'Неверный формат данных. Проверьте введенные данные.'
    
    case 401:
      return 'Неверный email или пароль'
    
    case 500:
      return 'Внутренняя ошибка сервера. Пожалуйста, попробуйте позже.'
    
    default:
      if (errorData?.message) {
        return errorData.message
      }
      return 'Произошла неизвестная ошибка'
  }
}

async function handleLogin() {
  error.value = null
  loading.value = true

  try {
    const res = await fetch('/api/v1/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({
        email: email.value,
        password: password.value,
      }),
    })

    if (!res.ok) {
      let errorData = null
      
      try {
        errorData = await res.json()
      } catch (e) {
        console.warn('Server returned non-JSON error response')
      }
      
      error.value = getErrorMessage(res.status, errorData)
      return
    }

    const data = await res.json()
    
    if (data.token) {
      document.cookie = `jwt_token=${data.token}; path=/; max-age=${7 * 24 * 60 * 60}`
    }

    router.push('/dashboard')
    
  } catch (err) {
    console.error('Login error:', err)
    if (err.name === 'TypeError' || err.name === 'NetworkError') {
      error.value = 'Ошибка сети. Проверьте подключение к интернету.'
    } else {
      error.value = 'Произошла непредвиденная ошибка'
    }
  } finally {
    loading.value = false
  }
}
</script>