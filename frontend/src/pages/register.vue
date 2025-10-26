<template>
  <div class="login-container">
    <h2>Регистрация</h2>

    <form @submit.prevent="handleRegister">
      <div class="input-group" :class="{ 'has-error': !!errors.name }">
        <input
          type="text"
          v-model.trim="name"
          placeholder="Имя"
          @blur="touch.name = true; validateName()"
          @input="validateName()"
          aria-describedby="name-hint"
          required
        />
        <p id="name-hint" class="hint" :class="{ 'hint-error': !!errors.name}">Введите имя</p>
      </div>

      <div class="input-group" :class="{ 'has-error': !!errors.email }">
        <input
          type="email"
          v-model.trim="email"
          placeholder="Почта"
          @blur="touch.email = true; validateEmail()"
          @input="validateEmail()"
          required
        />
        <p v-if="errors.email" class="field-error">Неправильно введена почта</p>
      </div>

      <div class="input-group" :class="{ 'has-error': !!errors.password }">
        <div class="control">
          <input
            :type="showPassword ? 'text' : 'password'"
            v-model="password"
            placeholder="Пароль"
            @blur="touch.password = true; validatePassword()"
            @input="validatePassword(); validateConfirm()"
            required
          />
          <button
            type="button"
            class="toggle-password"
            @click="showPassword = !showPassword"
            :class="{ 'visible': showPassword }"
            aria-label="Показать пароль"
          >
            <svg class="eye-icon" width="20" height="20" viewBox="0 0 24 24" fill="none">
              <line x1="4" y1="4" x2="20" y2="20" stroke="currentColor" stroke-width="2" class="strike-through"/>
              <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" stroke="currentColor" stroke-width="2"/>
              <circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="2"/>
            </svg>
          </button>
        </div>
        <p v-if="errors.password" class="field-error">{{ errors.password }}</p>
      </div>

      <div class="input-group" :class="{ 'has-error': !!errors.confirm }">
        <div class="control">
          <input
            :type="showConfirm ? 'text' : 'password'"
            v-model="confirm"
            placeholder="Повторите пароль"
            @blur="touch.confirm = true; validateConfirm()"
            @input="validateConfirm()"
            required
          />
          <button
            type="button"
            class="toggle-password"
            @click="showConfirm = !showConfirm"
            :class="{ 'visible': showConfirm }"
            aria-label="Показать пароль"
          >
            <svg class="eye-icon" width="20" height="20" viewBox="0 0 24 24" fill="none">
              <line x1="4" y1="4" x2="20" y2="20" stroke="currentColor" stroke-width="2" class="strike-through"/>
              <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" stroke="currentColor" stroke-width="2"/>
              <circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="2"/>
            </svg>
          </button>
        </div>
        <p v-if="errors.confirm" class="field-error">Пароли не совпадают</p>
      </div>

      <p v-if="generalError" class="error">{{ generalError }}</p>

      <button type="submit" class="submit-btn" :disabled="loading">
        {{ loading ? 'Отправка...' : 'Отправить' }}
      </button>
    </form>

    <div class="register-section">
      <p class="register-text">Уже есть аккаунт?</p>
      <router-link to="/login" class="create-account-btn">
        Войти
      </router-link>
    </div>
  </div>
</template>


<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const name = ref('')
const email = ref('')
const password = ref('')
const confirm = ref('')

const showPassword = ref(false)
const showConfirm  = ref(false)
const loading = ref(false)
const generalError = ref(null)
const errors = ref({ name: '', email: '', password: '', confirm: '' })
const touch = ref({ name: false, email: false, password: false, confirm: false })

function validateName () {
  errors.value.name = !name.value.trim() && touch.value.name ? 'Введите имя' : ''
}

function validateEmail () {
  const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  errors.value.email = touch.value.email && !re.test(email.value) ? 'Неправильно введена почта' : ''
}

function validatePassword () {
  if (!touch.value.password) return (errors.value.password = '')
  errors.value.password = password.value.length < 8 ? 'Минимум 8 символов' : ''
}

function validateConfirm () {
  if (!touch.value.confirm) return (errors.value.confirm = '')
  errors.value.confirm = password.value !== confirm.value ? 'Пароли не совпадают' : ''
}

function anyErrors () {
  return Boolean(errors.value.name || errors.value.email || errors.value.password || errors.value.confirm)
}

function getErrorMessage(status, errorData = null) {
  switch (status) {
    case 400:
      return errorData?.message || 'Неверный формат данных. Проверьте введенные данные.'
    case 409:
      return errorData?.message || 'Пользователь с такой почтой уже существует.'
    case 500:
      return 'Внутренняя ошибка сервера. Пожалуйста, попробуйте позже.'
    default:
      return errorData?.message || 'Произошла неизвестная ошибка'
  }
}

async function handleRegister () {
  touch.value = { name: true, email: true, password: true, confirm: true }
  validateName(); validateEmail(); validatePassword(); validateConfirm()

  if (anyErrors()) return

  generalError.value = null
  loading.value = true

  try {
    const res = await fetch('/api/v1/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({
        name: name.value.trim(),
        email: email.value.trim(),
        password: password.value,
        second_password: confirm.value
      })
    })

    if (!res.ok) {
      let errorData = null
      try { errorData = await res.json() } catch {}
      generalError.value = getErrorMessage(res.status, errorData)
      return
    }

    router.push('/login')
  } catch (err) {
    if (err.name === 'TypeError' || err.name === 'NetworkError') {
      generalError.value = 'Ошибка сети. Проверьте подключение к интернету.'
    } else {
      generalError.value = 'Произошла непредвиденная ошибка'
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
body {
  margin: 0;
  font-family: system-ui, sans-serif;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
}

.login-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
}

h2 {
  margin-bottom: 20px;
  color: #2a3c3c;
}

.input-group {
  margin-bottom: 15px;
  position: relative;
}

input {
  padding: 12px 14px;
  border: 1px solid #6c8c8c;
  border-radius: 5px;
  width: 300px;
  font-size: 15.5px;
  height: 40px;
  box-sizing: border-box;
  outline: none;
  transition: border-color .15s ease;
}

.control {
  position: relative;
  height: 40px; 
}

.control > input {
  height: 100%;
  width: 100%;
  padding-right: 45px;  
  box-sizing: border-box;
  line-height: 40px;     
  padding-top: 0;
  padding-bottom: 0;
}

.control > .toggle-password {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  cursor: pointer;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #666;
  border-radius: 4px;
  transition: all 0.2s ease;
}

.toggle-password:hover {
  color: #333;
  background-color: #f5f5f5;
}

.toggle-password:active {
  transform: translateY(-50%) scale(0.95);
}

.eye-icon { transition: all 0.2s ease; }
.strike-through { opacity: 1; transition: opacity 0.2s ease; }
.toggle-password.visible .strike-through { opacity: 0; }

.hint {
  margin: 6px 0 0;
  font-size: 12.5px;
  color: #9aa7a7;
}

.hint-error {
  color: #d43c3c;
}

.field-error {
  color: #d43c3c;
  margin: 6px 0 0;
  font-size: 13px;
}

.has-error input { border-color: #d43c3c; }

.error {
  color: red;
  margin-top: -6px;
  margin-bottom: 10px;
  font-size: 14px;
  width: 300px;
}

.submit-btn {
  background-color: #6c8c8c;
  border: none;
  color: white;
  height: 35px;
  padding: 10px 16px;
  border-radius: 6px;
  cursor: pointer;
}

.submit-btn:disabled {
  opacity: .8;
  cursor: default;
}

.register-section {
  text-align: center;
  margin-top: 20px;
}

.register-text {
  margin-top: -10px;
  margin-bottom: 10px;
  color: #666;
  font-size: 14px;
}

.create-account-btn {
  display: inline-block;
  background-color: transparent;
  border: 2px solid #6c8c8c;
  color: #6c8c8c;
  height: 35px;
  padding: 0 16px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  width: 150px;
  text-align: center;
  line-height: 31px;
  text-decoration: none;
  transition: all 0.2s ease;
  box-sizing: border-box;
}

.create-account-btn:hover {
  background-color: #6c8c8c;
  color: white;
}


input::-ms-reveal,
input::-ms-clear { display: none; }
</style>
