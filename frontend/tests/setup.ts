import { config } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, vi } from 'vitest'

// Fresh Pinia before each test
beforeEach(() => {
  setActivePinia(createPinia())
  // Clear sessionStorage between tests
  sessionStorage.clear()
})

// Stub router-link and router-view globally for component tests that don't mount with router
config.global.stubs = {
  RouterLink: { template: '<a><slot /></a>' },
  RouterView: { template: '<div />' },
}
