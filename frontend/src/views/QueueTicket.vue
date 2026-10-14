<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { queueService, type Ticket } from '../services/queue.service'

const route = useRoute()
const ticket = ref<Ticket | null>(null)
const loading = ref(true)
const error = ref('')


const issuedDateTime = computed(() => {
  if (!ticket.value) return ''

  const date = new Date(ticket.value.issuedAt)

  const dateText = new Intl.DateTimeFormat('en-GB', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    timeZone: 'Asia/Bangkok',
  }).format(date)

  const timeText = new Intl.DateTimeFormat('en-GB', {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
    timeZone: 'Asia/Bangkok',
  }).format(date)

  return `วันที่ : ${dateText} เวลา ${timeText} น.`
})

onMounted(async () => {
  try {
    ticket.value = await queueService.ticket(
      String(route.params.id)
    )
  } catch (err) {
    error.value = err instanceof Error
      ? err.message
      : 'ไม่พบข้อมูลบัตรคิว'
  } finally {
    loading.value = false
  }

})

</script>

<template>
  <main class="screen">
    <header class="screen-header">IT 05-2</header>

    <section class="ticket-content">
      <h1 class="queue-label">หมายเลขคิว</h1>

      <p v-if="loading">กำลังโหลดข้อมูล...</p>

      <template v-else-if="ticket">
        <div class="queue-number" data-testid="queue-number">
          {{ ticket.queueNumber }}
        </div>

        <p class="ticket-datetime" data-testid="ticket-date">
          {{ issuedDateTime }}
        </p>


      </template>

      <p v-else class="error" role="alert">
        {{ error }}
      </p>

      <router-link
        to="/"
        class="action-button back-button"
        data-testid="back-button"
      >
        กลับไปหน้ารับบัตรคิว
      </router-link>
    </section>
  </main>
</template>
