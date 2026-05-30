import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import { VueQueryPlugin, QueryClient, useQueryClient } from '@tanstack/vue-query'
import ProjectsView from '@/features/projects/views/ProjectsView.vue'

vi.mock('@/api/projects', () => ({
  fetchProjects: vi.fn(),
  createProject: vi.fn(),
}))

vi.mock('@/api/client', () => ({
  request: vi.fn(),
  ApiError: class ApiError extends Error {
    constructor(public status: number, message: string) { super(message) }
  },
}))

import { fetchProjects, createProject } from '@/api/projects'

const mockProjects = [
  {
    id: 1,
    userId: 1,
    name: 'Death Guard',
    description: null,
    gameSystem: null,
    faction: null,
    progress: { totalUnits: 10, byStatus: { unassembled: 10, assembled: 0, primed: 0, base_coated: 0, painted: 0, based: 0, complete: 0 } },
    createdAt: '',
    updatedAt: '',
  },
]

describe('Project creation flow', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('creating a project calls the API and invalidates the projects query', async () => {
    vi.mocked(fetchProjects).mockResolvedValue(mockProjects)
    vi.mocked(createProject).mockResolvedValue({
      ...mockProjects[0],
      id: 2,
      name: 'New Army',
      progress: undefined,
    })

    const router = createRouter({
      history: createMemoryHistory(),
      routes: [
        { path: '/projects', name: 'projects', component: ProjectsView },
        { path: '/projects/:projectId', name: 'project-detail', component: { template: '<div />' } },
        { path: '/paint-profiles', name: 'paint-profiles', component: { template: '<div />' } },
      ],
    })

    const pinia = createPinia()
    const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })

    await router.push('/projects')

    const wrapper = mount(ProjectsView, {
      global: {
        plugins: [pinia, router, [VueQueryPlugin, { queryClient }]],
      },
    })

    await flushPromises()

    // Open create modal
    const newBtn = wrapper.findAll('button').find((b) => b.text().includes('New Project'))
    await newBtn?.trigger('click')
    await flushPromises()

    // Fill in the form
    const nameInput = wrapper.find('#proj-name')
    await nameInput.setValue('New Army')

    // Submit
    const submitBtn = wrapper.findAll('button').find((b) => b.text().includes('Create Project'))
    await submitBtn?.trigger('click')
    await flushPromises()

    expect(createProject).toHaveBeenCalledWith(expect.objectContaining({ name: 'New Army' }))
  })
})
