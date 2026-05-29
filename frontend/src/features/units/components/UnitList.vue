<script setup lang="ts">
import { ref, computed } from 'vue'
import type { PaintProfile, UpdateUnitPayload, CreateUnitPayload } from '@/types'
import { useUnits } from '../composables/useUnits'
import UnitRow from './UnitRow.vue'
import UnitForm from './UnitForm.vue'
import BaseButton from '@/components/BaseButton.vue'
import BaseModal from '@/components/BaseModal.vue'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import ErrorBanner from '@/components/ErrorBanner.vue'
import EmptyState from '@/components/EmptyState.vue'

const props = defineProps<{
  projectId: number
  assignedProfileIds?: number[]
  profiles?: PaintProfile[]
}>()

const { unitsQuery, createMutation, updateMutation, deleteMutation } = useUnits(props.projectId)

const showAddModal   = ref(false)
const mutationError  = ref('')

const assignedProfiles = computed(() =>
  (props.profiles ?? []).filter((p) => props.assignedProfileIds?.includes(p.id)),
)

function handleCreate(payload: CreateUnitPayload) {
  mutationError.value = ''
  createMutation.mutate(payload, {
    onSuccess: () => { showAddModal.value = false },
    onError: (err) => { mutationError.value = err instanceof Error ? err.message : 'Failed to add unit.' },
  })
}

function handleUpdate(unitId: number, payload: UpdateUnitPayload) {
  mutationError.value = ''
  updateMutation.mutate({ unitId, payload }, {
    onError: (err) => { mutationError.value = err instanceof Error ? err.message : 'Failed to update unit.' },
  })
}

function handleDelete(unitId: number) {
  mutationError.value = ''
  deleteMutation.mutate(unitId, {
    onError: (err) => { mutationError.value = err instanceof Error ? err.message : 'Failed to remove unit.' },
  })
}

const errorMessage = computed(() =>
  unitsQuery.error.value instanceof Error ? unitsQuery.error.value.message : 'Failed to load units.',
)
</script>

<template>
  <div class="unit-list">
    <div class="unit-list__toolbar">
      <span class="unit-list__count" v-if="unitsQuery.data.value">
        {{ unitsQuery.data.value.length }} unit{{ unitsQuery.data.value.length !== 1 ? 's' : '' }}
      </span>
      <BaseButton variant="ghost" size="sm" @click="showAddModal = true">+ Add Unit</BaseButton>
    </div>

    <ErrorBanner v-if="mutationError" :message="mutationError" />

    <LoadingSpinner v-if="unitsQuery.isLoading.value" />
    <ErrorBanner v-else-if="unitsQuery.error.value" :message="errorMessage" />

    <EmptyState
      v-else-if="!unitsQuery.data.value?.length"
      heading="NO UNITS FILED"
      message="Add your first unit to start tracking progress."
    >
      <template #action>
        <BaseButton variant="primary" @click="showAddModal = true">Add Unit</BaseButton>
      </template>
    </EmptyState>

    <div v-else class="unit-list__rows">
      <UnitRow
        v-for="unit in unitsQuery.data.value"
        :key="unit.id"
        :unit="unit"
        :profiles="assignedProfiles"
        :update-loading="updateMutation.isPending.value"
        :delete-loading="deleteMutation.isPending.value"
        @update="handleUpdate"
        @delete="handleDelete"
      />
    </div>

    <BaseModal
      v-if="showAddModal"
      title="ADD UNIT"
      @close="showAddModal = false"
    >
      <UnitForm
        :profiles="assignedProfiles"
        :loading="createMutation.isPending.value"
        @submit="handleCreate"
        @cancel="showAddModal = false"
      />
    </BaseModal>
  </div>
</template>

<style lang="scss" scoped>
@use "@/styles/tokens" as *;
@use "@/styles/mixins" as *;

.unit-list {
  display: flex;
  flex-direction: column;
  gap: $spacing-3;

  &__toolbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  &__count {
    font-family: $font-mono;
    font-size: $font-mono-size;
    color: $color-text-secondary;
  }

  &__rows {
    display: flex;
    flex-direction: column;
    gap: $spacing-2;
  }
}
</style>
