import { useMutation } from '@tanstack/vue-query'
import { useAuthStore } from '@/stores/auth'
import { login as apiLogin, register as apiRegister, getCurrentUser } from '@/api/auth'
import type { LoginPayload, RegisterPayload } from '@/types'

export function useAuth() {
  const authStore = useAuthStore()

  const loginMutation = useMutation({
    mutationFn: async (payload: LoginPayload) => {
      const { token } = await apiLogin(payload)
      authStore.setToken(token)
      // Populate user after setting token
      authStore.user = await getCurrentUser()
    },
  })

  const registerMutation = useMutation({
    mutationFn: (payload: RegisterPayload) => apiRegister(payload),
  })

  return { loginMutation, registerMutation }
}
