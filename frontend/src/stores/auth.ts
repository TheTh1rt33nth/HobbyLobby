import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import type { User } from '@/types'
import { getCurrentUser } from '@/api/auth'

const TOKEN_KEY = 'hl_token'

export const useAuthStore = defineStore('auth', () => {
  // Token is kept in memory; sessionStorage is the persistence backing.
  // sessionStorage is cleared when the tab closes, limiting exposure.
  //TODO: persistent login
  const token = ref<string | null>(sessionStorage.getItem(TOKEN_KEY))
  const user  = ref<User | null>(null)

  const isAuthenticated = computed(() => !!token.value && !!user.value)

  // Called once on app mount. Validates stored token by fetching /api/users/me.
  async function bootstrap(): Promise<void> {
    if (!token.value) return
    try {
      user.value = await getCurrentUser()
    } catch {
      clearSession()
    }
  }

  function setToken(newToken: string): void {
    token.value = newToken
    sessionStorage.setItem(TOKEN_KEY, newToken)
  }

  function clearSession(): void {
    token.value = null
    user.value  = null
    sessionStorage.removeItem(TOKEN_KEY)
  }

  // Sign out is client-side only — no logout endpoint exists on the backend.
  // Token remains valid server-side for the remainder of its 24h window as of now
  // TODO: logout
  function signOut(): void {
    clearSession()
  }

  return {
    token,
    user,
    isAuthenticated,
    bootstrap,
    setToken,
    clearSession,
    signOut,
  }
})
