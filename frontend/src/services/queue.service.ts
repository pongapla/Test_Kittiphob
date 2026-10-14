export interface Ticket {
  id: string
  cycle: number
  queueNumber: string
  issuedAt: string
}

const BASE_URL = '/api/queues'

async function request<T>(
  path: string,
  options?: RequestInit
): Promise<T> {
  const response = await fetch(`${BASE_URL}${path}`, options)
  const data = await response.json()

  if (!response.ok) {
    throw new Error(data.error || 'API request failed')
  }

  return data as T
}

export const queueService = {
  issue: () =>
    request<Ticket>('/', { method: 'POST' }),

  current: () =>
    request<{ queueNumber: string }>('/current'),

  ticket: (id: string) =>
    request<Ticket>(`/${encodeURIComponent(id)}`),

  reset: () =>
    request<{ queueNumber: string; message: string }>(
      '/reset',
      { method: 'POST' }
    ),
}
