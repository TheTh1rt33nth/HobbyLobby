// HobbyLobby — API Client
//

export class ApiError extends Error {
  constructor(
    public readonly status: number,
    message: string,
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

export async function request<T>(
  path: string,
  options: RequestInit = {},
): Promise<T> {
  const { useAuthStore } = await import('@/stores/auth')
  const { getRouter }    = await import('@/router')

  const authStore = useAuthStore()
  const token     = authStore.token

  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string> | undefined),
  }

  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const res = await fetch(`/api${path}`, { ...options, headers })

  if (res.status === 401) {
    authStore.clearSession()
    getRouter().push({ name: 'login' })
    throw new ApiError(401, 'Session expired. Please sign in again.')
  }

  // No-content responses
  if (res.status === 204 || res.headers.get('content-length') === '0') {
    return undefined as T
  }

  const body = await res.json().catch(() => ({}))

  if (!res.ok) {
    throw new ApiError(
      res.status,
      (body as { error?: string }).error ?? `Request failed (${res.status})`,
    )
  }

  return body as T
}
