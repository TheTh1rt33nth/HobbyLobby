<script setup lang="ts">
import { ref, computed } from 'vue'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { fetchPaintProfiles, createPaintProfile } from '@/api/paintProfiles'
import type { CreatePaintProfilePayload } from '@/types'

import PaintProfileCard from '../components/PaintProfileCard.vue'
import PaintProfileForm from '../components/PaintProfileForm.vue'
import BaseButton from '@/components/BaseButton.vue'
import BaseModal from '@/components/BaseModal.vue'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import ErrorBanner from '@/components/ErrorBanner.vue'
import EmptyState from '@/components/EmptyState.vue'

const queryClient = useQueryClient()

const { data: profiles, isLoading, error } = useQuery({
  queryKey: ['paint-profiles'],
  queryFn: fetchPaintProfiles,
})

const showCreateModal = ref(false)
const createError     = ref('')

const createMutation = useMutation({
  mutationFn: createPaintProfile,
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['paint-profiles'] })
    showCreateModal.value = false
    createError.value = ''
  },
  onError: (err) => {
    createError.value = err instanceof Error ? err.message : 'Failed to create profile.'
  },
})

const errorMessage = computed(() =>
  error.value instanceof Error ? error.value.message : 'Failed to load paint profiles.',
)

function handleCreate(payload: CreatePaintProfilePayload) {
  createMutation.mutate(payload)
}
</script>

<template>
  <div class="d-flex flex-column gap-5">
    <header class="page-header">
      <div>
        <h1 class="page-title">PAINT PROFILES</h1>
        <p class="page-subtitle">Reusable painting schemes and step-by-step recipes</p>
      </div>
      <BaseButton variant="primary" @click="showCreateModal = true">+ New Profile</BaseButton>
    </header>

    <LoadingSpinner v-if="isLoading" />
    <ErrorBanner v-else-if="error" :message="errorMessage" />

    <template v-else-if="profiles">
      <EmptyState
        v-if="profiles.length === 0"
        heading="NO SCHEMES LOGGED"
        message="Document your painting process by creating a reusable paint profile."
      >
        <template #action>
          <BaseButton variant="primary" @click="showCreateModal = true">
            Create First Profile
          </BaseButton>
        </template>
      </EmptyState>

      <div v-else class="profiles-view__grid">
        <PaintProfileCard v-for="profile in profiles" :key="profile.id" :profile="profile" />
      </div>
    </template>

    <BaseModal v-if="showCreateModal" title="NEW PAINT PROFILE" @close="showCreateModal = false">
      <ErrorBanner v-if="createError" :message="createError" />
      <PaintProfileForm
        :loading="createMutation.isPending.value"
        @submit="handleCreate"
        @cancel="showCreateModal = false"
      />
    </BaseModal>
  </div>
</template>

<style lang="scss" scoped>
@use "@/styles/tokens" as *;

.profiles-view__grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: $spacing-4;
}
</style>
