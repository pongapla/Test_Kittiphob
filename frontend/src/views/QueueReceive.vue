<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { queueService } from '../services/queue.service'

const router = useRouter()
const loading = ref(false)
const error = ref('')

async function receiveTicket() {
  if (loading.value) return
  loading.value = true
  error.value = ''

  try {
    const ticket = await queueService.issue()
    await router.push(`/ticket/${ticket.id}`)
  } catch (err) {
    error.value = err instanceof Error
      ? err.message : 'ไม่สามารถรับบัตรคิวได้'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <main class="screen">
    <header class="screen-header">IT 05-1</header>

    <section class="receive-content">
      <button
        class="receive-button"
        data-testid="receive-button"
        :disabled="loading"
        @click="receiveTicket"
      >
        {{ loading ? 'กำลังรับคิว...' : 'รับบัตรคิว' }}
      </button>

      <router-link
        to="/reset"
        class="action-button reset-navigation"
        data-testid="reset-navigation"
      >
        ล้างคิว
      </router-link>

      <p v-if="error" class="error" role="alert">
        {{ error }}
      </p>
    </section>
  </main>
</template>
