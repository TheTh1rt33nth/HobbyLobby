<script setup lang="ts">
import { ref, computed } from 'vue'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { fetchProjects, createProject } from '@/api/projects'
import type { CreateProjectPayload } from '@/types'
import ProjectCard from '../components/ProjectCard.vue'
import ProjectForm from '../components/ProjectForm.vue'
import BaseButton from '@/components/BaseButton.vue'
import BaseModal from '@/components/BaseModal.vue'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import ErrorBanner from '@/components/ErrorBanner.vue'
import EmptyState from '@/components/EmptyState.vue'

const queryClient = useQueryClient()

const { data: projects, isLoading, error } = useQuery({
  queryKey: ['projects'],
  queryFn: fetchProjects,
})

const showCreateModal = ref(false)
const createError     = ref('')

const createMutation = useMutation({
  mutationFn: createProject,
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['projects'] })
    showCreateModal.value = false
    createError.value = ''
  },
  onError: (err) => {
    createError.value = err instanceof Error ? err.message : 'Failed to create project.'
  },
})

const errorMessage = computed(() =>
  error.value instanceof Error ? error.value.message : 'Failed to load projects.',
)

function handleCreate(payload: CreateProjectPayload) {
  createMutation.mutate(payload)
}
</script>

<template>
  <div class="d-flex flex-column gap-5">
    <header class="page-header">
      <div>
        <h1 class="page-title">PROJECTS</h1>
        <p class="page-subtitle">Your hobby campaigns and battle forces</p>
      </div>
      <BaseButton variant="primary" @click="showCreateModal = true">
        + New Project
      </BaseButton>
    </header>

    <LoadingSpinner v-if="isLoading" />
    <ErrorBanner v-else-if="error" :message="errorMessage" />

    <template v-else-if="projects">
      <EmptyState
        v-if="projects.length === 0"
        heading="NO PROJECTS FILED"
        message="Begin your campaign by creating your first hobby project."
      >
        <template #action>
          <BaseButton variant="primary" @click="showCreateModal = true">
            Create First Project
          </BaseButton>
        </template>
      </EmptyState>

      <div v-else class="projects-view__grid">
        <ProjectCard v-for="project in projects" :key="project.id" :project="project" />
      </div>
    </template>

    <!-- Create modal -->
    <BaseModal
      v-if="showCreateModal"
      title="NEW PROJECT"
      @close="showCreateModal = false"
    >
      <ErrorBanner v-if="createError" :message="createError" />
      <ProjectForm
        :loading="createMutation.isPending.value"
        @submit="handleCreate"
        @cancel="showCreateModal = false"
      />
    </BaseModal>
  </div>
</template>

<style lang="scss" scoped>
@use "@/styles/tokens" as *;

.projects-view__grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: $spacing-4;
}
</style>
