<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { queueService } from '../services/queue.service'

const queueNumber = ref('00')
const loading = ref(false)
const error = ref('')

onMounted(async () => {
  loading.value = true
  try {
    const result = await queueService.current()
    queueNumber.value = result.queueNumber
  } catch (err) {
    error.value = err instanceof Error
      ? err.message : 'ไม่สามารถโหลดหมายเลขคิวได้'
  } finally {
    loading.value = false
  }
})

async function resetQueue() {
  if (loading.value) return

  loading.value = true
  error.value = ''

  try {
    const result = await queueService.reset()
    queueNumber.value = result.queueNumber
  } catch (err) {
    error.value = err instanceof Error
      ? err.message : 'ไม่สามารถล้างคิวได้'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <main class="screen">
    <header class="screen-header">IT 05-3</header>

    <section class="reset-content">
      <button
        class="action-button reset-action"
        data-testid="reset-button"
        :disabled="loading"
        @click="resetQueue"
      >
        ล้างคิว
      </button>

      <h1 class="queue-label">หมายเลขคิวปัจจุบัน</h1>

      <div
        class="queue-number"
        data-testid="reset-number"
      >
        {{ loading && error ? '-' : queueNumber }}
      </div>

      <router-link
        to="/"
        class="action-button back-button"
        data-testid="back-button"
      >
        กลับไปหน้ารับบัตรคิว
      </router-link>

      <p v-if="error" class="error" role="alert">
        {{ error }}
      </p>
    </section>
  </main>
</template>
