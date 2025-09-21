<script setup>
import { ref, onMounted } from 'vue'

const message = ref('...')
onMounted(async () => {
  const r = await fetch('/api/hello')
  const data = await r.json()
  message.value = data.message
})
</script>

<template>
  <main style="font: 16px/1.4 system-ui, sans-serif; padding: 24px">
    <h1>QR Code Gen (Vue)</h1>
    <p>Ответ бэкенда: <b>{{ message }}</b></p>
  </main>
</template>

<style>
body { margin: 0; }
</style>
