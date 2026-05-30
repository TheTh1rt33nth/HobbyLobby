import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import { createPinia } from 'pinia'
import { VueQueryPlugin, QueryClient } from '@tanstack/vue-query'
import LoginView from '@/features/auth/views/LoginView.vue'

// Mock the auth API
vi.mock('@/api/auth', () => ({
  login: vi.fn(),
  getCurrentUser: vi.fn(),
  register: vi.fn(),
}))

vi.mock('@/api/client', () => ({
  request: vi.fn(),
  ApiError: class ApiError extends Error {
    constructor(public status: number, message: string) { super(message) }
  },
}))

import { login, getCurrentUser } from '@/api/auth'

const mockUser = { id: 1, username: 'testOp', email: 'op@test.com', createdAt: '', updatedAt: '' }

function createTestApp() {
  const router = createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/login',     name: 'login',     component: LoginView },
      { path: '/dashboard', name: 'dashboard', component: { template: '<div>Dashboard</div>' } },
    ],
  })

  const pinia = createPinia()
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })

  return { router, pinia, queryClient }
}

describe('Login flow', () => {
  beforeEach(() => {
    sessionStorage.clear()
    vi.clearAllMocks()
  })

  it('successful login navigates to dashboard', async () => {
    vi.mocked(login).mockResolvedValue({ token: 'abc123', expiry: '2099-01-01' })
    vi.mocked(getCurrentUser).mockResolvedValue(mockUser)

    const { router, pinia, queryClient } = createTestApp()
    await router.push('/login')

    const wrapper = mount(LoginView, {
      global: {
        plugins: [pinia, router, [VueQueryPlugin, { queryClient }]],
      },
    })

    await wrapper.find('#username').setValue('testOp')
    await wrapper.find('#password').setValue('password123')
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    expect(router.currentRoute.value.name).toBe('dashboard')
  })

  it('failed login shows error banner', async () => {
    const { ApiError } = await import('@/api/client')
    vi.mocked(login).mockRejectedValue(new ApiError(401, 'Invalid credentials'))

    const { router, pinia, queryClient } = createTestApp()
    await router.push('/login')

    const wrapper = mount(LoginView, {
      global: {
        plugins: [pinia, router, [VueQueryPlugin, { queryClient }]],
      },
    })

    await wrapper.find('#username').setValue('badOp')
    await wrapper.find('#password').setValue('wrongPass')
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    expect(wrapper.find('.error-banner').exists()).toBe(true)
    expect(wrapper.find('.error-banner__message').text()).toContain('Invalid credentials')
  })
})
