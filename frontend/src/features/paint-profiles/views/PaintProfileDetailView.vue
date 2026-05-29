<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import {
  fetchPaintProfile,
  updatePaintProfile,
  deletePaintProfile,
  createPaintStep,
  updatePaintStep,
  deletePaintStep,
} from '@/api/paintProfiles'
import type { PaintStep, CreatePaintStepPayload, UpdatePaintStepPayload } from '@/types'

import LoadingSpinner from '@/components/LoadingSpinner.vue'
import ErrorBanner from '@/components/ErrorBanner.vue'
import BaseButton from '@/components/BaseButton.vue'
import BaseModal from '@/components/BaseModal.vue'
import ConfirmModal from '@/components/ConfirmModal.vue'
import PaintStepList from '../components/PaintStepList.vue'
import PaintStepForm from '../components/PaintStepForm.vue'
import PaintProfileForm from '../components/PaintProfileForm.vue'
import type { UpdatePaintProfilePayload } from '@/types'

const props = defineProps<{ profileId: number }>()

const router      = useRouter()
const queryClient = useQueryClient()

const profileKey = computed(() => ['paint-profile', props.profileId])

const { data: profile, isLoading, error } = useQuery({
  queryKey: profileKey,
  queryFn: () => fetchPaintProfile(props.profileId),
})

// ── Edit profile ──────────────────────────────────────────────────────────────

const showEditModal = ref(false)
const editError     = ref('')

const editMutation = useMutation({
  mutationFn: (payload: UpdatePaintProfilePayload) =>
    updatePaintProfile(props.profileId, payload),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: profileKey.value })
    queryClient.invalidateQueries({ queryKey: ['paint-profiles'] })
    showEditModal.value = false
    editError.value = ''
  },
  onError: (err) => {
    editError.value = err instanceof Error ? err.message : 'Failed to update profile.'
  },
})

// ── Delete profile ────────────────────────────────────────────────────────────

const showDeleteModal = ref(false)
const deleteError     = ref('')

const deleteMutation = useMutation({
  mutationFn: () => deletePaintProfile(props.profileId),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['paint-profiles'] })
    router.push({ name: 'paint-profiles' })
  },
  onError: (err) => {
    deleteError.value = err instanceof Error ? err.message : 'Failed to delete profile.'
  },
})

// ── Steps ─────────────────────────────────────────────────────────────────────

const showAddStepModal  = ref(false)
const editingStep       = ref<PaintStep | null>(null)
const stepMutationError = ref('')

const nextStepOrder = computed(() => {
  const steps = profile.value?.steps ?? []
  return steps.length > 0 ? Math.max(...steps.map((s) => s.stepOrder)) + 1 : 1
})

const addStepMutation = useMutation({
  mutationFn: (payload: CreatePaintStepPayload) =>
    createPaintStep(props.profileId, payload),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: profileKey.value })
    showAddStepModal.value = false
    stepMutationError.value = ''
  },
  onError: (err) => {
    stepMutationError.value = err instanceof Error ? err.message : 'Failed to add step.'
  },
})

const editStepMutation = useMutation({
  mutationFn: ({ stepId, payload }: { stepId: number; payload: UpdatePaintStepPayload }) =>
    updatePaintStep(props.profileId, stepId, payload),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: profileKey.value })
    editingStep.value = null
    stepMutationError.value = ''
  },
  onError: (err) => {
    stepMutationError.value = err instanceof Error ? err.message : 'Failed to update step.'
  },
})

const deleteStepMutation = useMutation({
  mutationFn: (stepId: number) => deletePaintStep(props.profileId, stepId),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: profileKey.value })
  },
})

const errorMessage = computed(() =>
  error.value instanceof Error ? error.value.message : 'Failed to load profile.',
)

function handleEditStep(step: PaintStep) {
  editingStep.value = step
}

function handleSubmitEditStep(payload: UpdatePaintStepPayload) {
  if (!editingStep.value) return
  editStepMutation.mutate({ stepId: editingStep.value.id, payload })
}
</script>

<template>
  <div class="profile-detail">
    <LoadingSpinner v-if="isLoading" />
    <ErrorBanner v-else-if="error" :message="errorMessage" />

    <template v-else-if="profile">
      <!-- Header -->
      <header class="profile-detail__header">
        <div class="profile-detail__breadcrumb">
          <RouterLink :to="{ name: 'paint-profiles' }" class="profile-detail__back">← PAINT PROFILES</RouterLink>
        </div>

        <div class="profile-detail__title-row">
          <div>
            <h1 class="page-title">{{ profile.name }}</h1>
            <span v-if="profile.targetArea" class="profile-detail__area">{{ profile.targetArea }}</span>
          </div>
          <div class="profile-detail__actions">
            <BaseButton variant="secondary" size="sm" @click="showEditModal = true">Edit</BaseButton>
            <BaseButton variant="danger"    size="sm" @click="showDeleteModal = true">Delete</BaseButton>
          </div>
        </div>

        <p v-if="profile.description" class="profile-detail__description">{{ profile.description }}</p>
      </header>

      <!-- Steps section -->
      <section class="d-flex flex-column gap-3">
        <div class="section-header">
          <h2 class="section-title">PAINT STEPS</h2>
          <BaseButton variant="ghost" size="sm" @click="showAddStepModal = true">+ Add Step</BaseButton>
        </div>

        <ErrorBanner v-if="stepMutationError" :message="stepMutationError" />

        <PaintStepList
          :steps="profile.steps ?? []"
          :editable="true"
          @edit="handleEditStep"
          @delete="deleteStepMutation.mutate"
        />
      </section>
    </template>

    <!-- Edit Profile Modal -->
    <BaseModal v-if="showEditModal" title="EDIT PROFILE" @close="showEditModal = false">
      <ErrorBanner v-if="editError" :message="editError" />
      <PaintProfileForm
        :initial="profile"
        :loading="editMutation.isPending.value"
        @submit="editMutation.mutate"
        @cancel="showEditModal = false"
      />
    </BaseModal>

    <!-- Delete Profile Confirm -->
    <ConfirmModal
      v-if="showDeleteModal"
      title="DELETE PAINT PROFILE"
      :message="`Permanently delete '${profile?.name}'? This will remove it from all assigned projects.`"
      confirm-label="Delete Profile"
      :loading="deleteMutation.isPending.value"
      @confirm="deleteMutation.mutate"
      @cancel="showDeleteModal = false"
    />

    <!-- Add Step Modal -->
    <BaseModal v-if="showAddStepModal" title="ADD PAINT STEP" @close="showAddStepModal = false">
      <PaintStepForm
        :next-order="nextStepOrder"
        :loading="addStepMutation.isPending.value"
        @submit="addStepMutation.mutate"
        @cancel="showAddStepModal = false"
      />
    </BaseModal>

    <!-- Edit Step Modal -->
    <BaseModal
      v-if="editingStep"
      title="EDIT PAINT STEP"
      @close="editingStep = null"
    >
      <PaintStepForm
        :initial="editingStep"
        :loading="editStepMutation.isPending.value"
        @submit="handleSubmitEditStep"
        @cancel="editingStep = null"
      />
    </BaseModal>

    <ErrorBanner v-if="deleteError" :message="deleteError" />
  </div>
</template>

<style lang="scss" scoped>
@use "@/styles/tokens" as *;
@use "@/styles/mixins" as *;

.profile-detail {
  display: flex;
  flex-direction: column;
  gap: $spacing-8;

  &__header {
    display: flex;
    flex-direction: column;
    gap: $spacing-4;
    padding-bottom: $spacing-6;
    border-bottom: 1px solid $color-border-default;
  }

  &__back {
    @include section-label;
    font-size: 11px;
    color: $color-accent-amber;
    text-decoration: none;

    &:hover { color: $color-text-primary; text-decoration: none; }
  }

  &__title-row {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: $spacing-4;
  }

  &__area {
    @include section-label;
    font-size: 11px;
    margin-top: $spacing-1;
    display: block;
  }

  &__actions {
    display: flex;
    gap: $spacing-2;
    flex-shrink: 0;
  }

  &__description {
    color: $color-text-secondary;
    font-size: $font-body-sm;
    line-height: 1.6;
  }
}
</style>
