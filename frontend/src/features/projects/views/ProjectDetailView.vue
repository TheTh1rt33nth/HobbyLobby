<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import {
  fetchProject,
  updateProject,
  deleteProject,
  fetchProjectPaintProfiles,
  assignPaintProfile,
  unassignPaintProfile,
} from '@/api/projects'
import { fetchPaintProfiles } from '@/api/paintProfiles'
import type { UpdateProjectPayload } from '@/types'

import LoadingSpinner from '@/components/LoadingSpinner.vue'
import ErrorBanner from '@/components/ErrorBanner.vue'
import BaseButton from '@/components/BaseButton.vue'
import BaseModal from '@/components/BaseModal.vue'
import ConfirmModal from '@/components/ConfirmModal.vue'
import SegmentedProgress from '@/components/SegmentedProgress.vue'
import EmptyState from '@/components/EmptyState.vue'
import ProjectForm from '../components/ProjectForm.vue'
import UnitList from '@/features/units/components/UnitList.vue'

const props = defineProps<{ projectId: number }>()

const router       = useRouter()
const queryClient  = useQueryClient()

// ── Queries ──────────────────────────────────────────────────────────────────

const { data: project, isLoading, error } = useQuery({
  queryKey: computed(() => ['project', props.projectId]),
  queryFn: () => fetchProject(props.projectId),
})

const { data: assignedProfiles } = useQuery({
  queryKey: computed(() => ['project-paint-profiles', props.projectId]),
  queryFn: () => fetchProjectPaintProfiles(props.projectId),
})

const { data: allProfiles } = useQuery({
  queryKey: ['paint-profiles'],
  queryFn: fetchPaintProfiles,
})

// ── Edit project mutation ─────────────────────────────────────────────────────

const showEditModal  = ref(false)
const editError      = ref('')

const editMutation = useMutation({
  mutationFn: (payload: UpdateProjectPayload) =>
    updateProject(props.projectId, payload),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['project', props.projectId] })
    queryClient.invalidateQueries({ queryKey: ['projects'] })
    showEditModal.value = false
    editError.value = ''
  },
  onError: (err) => {
    editError.value = err instanceof Error ? err.message : 'Failed to update project.'
  },
})

// ── Delete project mutation ───────────────────────────────────────────────────

const showDeleteModal  = ref(false)
const deleteError      = ref('')

const deleteMutation = useMutation({
  mutationFn: () => deleteProject(props.projectId),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['projects'] })
    router.push({ name: 'projects' })
  },
  onError: (err) => {
    deleteError.value = err instanceof Error ? err.message : 'Failed to delete project.'
  },
})

// ── Assign / unassign paint profile mutations ─────────────────────────────────

const assignMutation = useMutation({
  mutationFn: (profileId: number) => assignPaintProfile(props.projectId, profileId),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['project-paint-profiles', props.projectId] })
  },
})

const unassignMutation = useMutation({
  mutationFn: (profileId: number) => unassignPaintProfile(props.projectId, profileId),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['project-paint-profiles', props.projectId] })
  },
})

// ── Helpers ───────────────────────────────────────────────────────────────────

const assignedIds = computed(() =>
  new Set((assignedProfiles.value ?? []).map((p) => p.id)),
)

const unassignedProfiles = computed(() =>
  (allProfiles.value ?? []).filter((p) => !assignedIds.value.has(p.id)),
)

const errorMessage = computed(() =>
  error.value instanceof Error ? error.value.message : 'Failed to load project.',
)

const showAssignModal = ref(false)
</script>

<template>
  <div class="project-detail">
    <LoadingSpinner v-if="isLoading" />
    <ErrorBanner v-else-if="error" :message="errorMessage" />

    <template v-else-if="project">
      <!-- Header -->
      <header class="project-detail__header">
        <div class="project-detail__breadcrumb">
          <RouterLink :to="{ name: 'projects' }" class="project-detail__back">← PROJECTS</RouterLink>
        </div>

        <div class="project-detail__title-row">
          <div>
            <h1 class="page-title">{{ project.name }}</h1>
            <div class="project-detail__meta">
              <span v-if="project.gameSystem" class="project-detail__tag">{{ project.gameSystem }}</span>
              <span v-if="project.faction"    class="project-detail__tag">{{ project.faction }}</span>
            </div>
          </div>
          <div class="project-detail__actions">
            <BaseButton variant="secondary" size="sm" @click="showEditModal = true">Edit</BaseButton>
            <BaseButton variant="danger"    size="sm" @click="showDeleteModal = true">Delete</BaseButton>
          </div>
        </div>

        <p v-if="project.description" class="project-detail__description">{{ project.description }}</p>

        <div v-if="project.progress" class="project-detail__progress">
          <SegmentedProgress :progress="project.progress" />
        </div>
      </header>

      <!-- Units section -->
      <section class="d-flex flex-column gap-3">
        <h2 class="section-title">UNIT ROSTER</h2>
        <UnitList
          :project-id="projectId"
          :assigned-profile-ids="[...assignedIds]"
          :profiles="allProfiles ?? []"
        />
      </section>

      <!-- Paint profiles section -->
      <section class="d-flex flex-column gap-3">
        <div class="section-header">
          <h2 class="section-title">ASSIGNED PAINT SCHEMES</h2>
          <BaseButton
            variant="ghost"
            size="sm"
            @click="showAssignModal = true"
            :disabled="unassignedProfiles.length === 0"
          >
            + Assign Profile
          </BaseButton>
        </div>

        <EmptyState
          v-if="!assignedProfiles?.length"
          heading="NO SCHEMES ASSIGNED"
          message="Assign a paint profile to track the painting scheme for this project."
        />

        <div v-else class="project-detail__profiles">
          <div
            v-for="profile in assignedProfiles"
            :key="profile.id"
            class="profile-row"
          >
            <div class="profile-row__info">
              <RouterLink
                :to="{ name: 'paint-profile-detail', params: { profileId: profile.id } }"
                class="profile-row__name"
              >{{ profile.name }}</RouterLink>
              <span v-if="profile.targetArea" class="profile-row__area">{{ profile.targetArea }}</span>
            </div>
            <BaseButton
              variant="ghost"
              size="sm"
              @click="unassignMutation.mutate(profile.id)"
              :loading="unassignMutation.isPending.value"
            >
              Remove
            </BaseButton>
          </div>
        </div>
      </section>
    </template>

    <!-- Edit Modal -->
    <BaseModal v-if="showEditModal" title="EDIT PROJECT" @close="showEditModal = false">
      <ErrorBanner v-if="editError" :message="editError" />
      <ProjectForm
        :initial="project"
        :loading="editMutation.isPending.value"
        @submit="editMutation.mutate"
        @cancel="showEditModal = false"
      />
    </BaseModal>

    <!-- Delete Confirm -->
    <ConfirmModal
      v-if="showDeleteModal"
      title="DELETE PROJECT"
      :message="`Permanently delete '${project?.name}'? This will also delete all units in this project.`"
      confirm-label="Delete Project"
      :loading="deleteMutation.isPending.value"
      @confirm="deleteMutation.mutate"
      @cancel="showDeleteModal = false"
    />

    <!-- Assign Profile Modal -->
    <BaseModal
      v-if="showAssignModal"
      title="ASSIGN PAINT PROFILE"
      max-width="400px"
      @close="showAssignModal = false"
    >
      <div class="assign-list">
        <p v-if="unassignedProfiles.length === 0" class="assign-list__empty">
          All profiles are already assigned to this project.
        </p>
        <button
          v-for="p in unassignedProfiles"
          :key="p.id"
          class="assign-list__item"
          @click="assignMutation.mutate(p.id); showAssignModal = false"
          type="button"
        >
          <span class="assign-list__name">{{ p.name }}</span>
          <span v-if="p.targetArea" class="assign-list__area">{{ p.targetArea }}</span>
        </button>
      </div>
    </BaseModal>

    <ErrorBanner v-if="deleteError" :message="deleteError" />
  </div>
</template>

<style lang="scss" scoped>
@use "@/styles/tokens" as *;
@use "@/styles/mixins" as *;

.project-detail {
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

  &__meta {
    display: flex;
    flex-wrap: wrap;
    gap: $spacing-1;
    margin-top: $spacing-1;
  }

  &__tag {
    @include section-label;
    font-size: 10px;
    padding: 1px $spacing-2;
    border: 1px solid $color-border-default;
    border-radius: $radius-sm;
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

  &__progress {
    max-width: 600px;
  }

  &__profiles {
    display: flex;
    flex-direction: column;
    gap: $spacing-2;
  }
}

.profile-row {
  @include panel-frame;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: $spacing-3 $spacing-4;

  &__info {
    display: flex;
    align-items: center;
    gap: $spacing-3;
  }

  &__name {
    font-family: $font-heading;
    font-size: $font-heading-sm;
    font-weight: $weight-semibold;
    color: $color-text-primary;
    text-decoration: none;

    &:hover { color: $color-accent-amber; text-decoration: none; }
  }

  &__area {
    @include section-label;
    font-size: 10px;
    color: $color-text-secondary;
  }
}

.assign-list {
  display: flex;
  flex-direction: column;
  gap: $spacing-2;

  &__empty {
    color: $color-text-secondary;
    font-size: $font-body-sm;
    text-align: center;
    padding: $spacing-4 0;
  }

  &__item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: $spacing-3 $spacing-4;
    background: $color-bg-elevated;
    border: 1px solid $color-border-default;
    border-radius: $radius-md;
    cursor: pointer;
    text-align: left;
    transition: border-color 0.15s;

    &:hover { border-color: $color-accent-amber; }
  }

  &__name {
    font-family: $font-heading;
    font-size: $font-heading-sm;
    font-weight: $weight-semibold;
    color: $color-text-primary;
  }

  &__area {
    @include section-label;
    font-size: 10px;
    color: $color-text-secondary;
  }
}
</style>
