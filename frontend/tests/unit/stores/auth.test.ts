import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '@/stores/auth'

// Mock the API module
vi.mock('@/api/auth', () => ({
  getCurrentUser: vi.fn(),
}))

import { getCurrentUser } from '@/api/auth'

describe('auth store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    sessionStorage.clear()
    vi.clearAllMocks()
  })

  it('bootstrap: valid token → populates user', async () => {
    const mockUser = { id: 1, username: 'testUser', email: 'test@example.com', createdAt: '', updatedAt: '' }
    vi.mocked(getCurrentUser).mockResolvedValue(mockUser)

    sessionStorage.setItem('hl_token', 'valid-token-123')
    const store = useAuthStore()

    await store.bootstrap()

    expect(store.user).toEqual(mockUser)
    expect(store.isAuthenticated).toBe(true)
  })

  it('bootstrap: no token → does not call API', async () => {
    const store = useAuthStore()
    await store.bootstrap()

    expect(getCurrentUser).not.toHaveBeenCalled()
    expect(store.isAuthenticated).toBe(false)
  })

  it('bootstrap: 401 → clears session', async () => {
    const { ApiError } = await import('@/api/client')
    vi.mocked(getCurrentUser).mockRejectedValue(new ApiError(401, 'Unauthorized'))

    sessionStorage.setItem('hl_token', 'expired-token')
    const store = useAuthStore()

    await store.bootstrap()

    expect(store.user).toBeNull()
    expect(store.token).toBeNull()
    expect(store.isAuthenticated).toBe(false)
    expect(sessionStorage.getItem('hl_token')).toBeNull()
  })

  it('signOut: clears all session state', () => {
    const store = useAuthStore()
    store.setToken('some-token')
    store.user = { id: 2, username: 'op', email: 'op@x.com', createdAt: '', updatedAt: '' }

    store.signOut()

    expect(store.token).toBeNull()
    expect(store.user).toBeNull()
    expect(store.isAuthenticated).toBe(false)
    expect(sessionStorage.getItem('hl_token')).toBeNull()
  })
})
