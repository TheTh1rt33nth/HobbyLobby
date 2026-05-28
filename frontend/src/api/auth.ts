import { request } from './client'
import type { User, RegisterPayload, LoginPayload } from '@/types'

export async function register(payload: RegisterPayload): Promise<User> {
  const data = await request<{ user: User }>('/users/register', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
  return data.user
}

export async function login(
  payload: LoginPayload,
): Promise<{ token: string; expiry: string }> {
  const data = await request<{ auth_token: { token: string; expiry: string } }>(
    '/tokens/auth',
    {
      method: 'POST',
      body: JSON.stringify(payload),
    },
  )
  return data.auth_token
}

export async function getCurrentUser(): Promise<User> {
  const data = await request<{ user: User }>('/users/me')
  return data.user
}
