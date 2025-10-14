<template>
  <div class="login-container">
    <h2>Вход</h2>
    <form @submit.prevent="handleLogin">
      <div class="input-group">
        <input
          type="email"
          v-model="email"
          placeholder="Почта"
          required
        />
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
    </form>
    <p class="create-account">
      <router-link to="/register">Создать аккаунт</router-link>
    </p>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const email = ref('')
const password = ref('')
const showPassword = ref(false)
const router = useRouter() 

async function handleLogin() {
  try {
    const res = await fetch('http://localhost:8080/api/v1/login', { 
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: email.value, password: password.value })
    })

    if (!res.ok) {
      const error = await res.json()
      alert('Ошибка: ' + (error.message || res.statusText))
      return
    }

    const data = await res.json()
    console.log('Успешный логин:', data)

    localStorage.setItem('token', data.token)

    router.push('/dashboard')
  } catch (err) {
    console.error(err)
    alert('Ошибка сети')
  }
}
</script>


<style scoped>
.login-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
}

h2 {
  margin-bottom: 20px;
}

.input-group {
  margin-bottom: 15px;
  position: relative;
}

input {
  padding: 8px 12px;
  border: 1px solid #6c8c8c;
  border-radius: 5px;
  width: 250px;
}

.password-group .toggle {
  position: absolute;
  right: 10px;
  top: 8px;
  cursor: pointer;
}

.submit-btn {
  background-color: #6c8c8c;
  border: none;
  color: white;
  padding: 10px 16px;
  border-radius: 6px;
  cursor: pointer;
}

.create-account {
  margin-top: 10px;
  font-size: 14px;
}
</style>
