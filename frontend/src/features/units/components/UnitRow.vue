<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Unit, PaintProfile, UpdateUnitPayload } from '@/types'
import StatusChip from '@/components/StatusChip.vue'
import ColorSwatch from '@/components/ColorSwatch.vue'
import BaseButton from '@/components/BaseButton.vue'
import ConfirmModal from '@/components/ConfirmModal.vue'
import BaseModal from '@/components/BaseModal.vue'
import UnitForm from './UnitForm.vue'

const props = defineProps<{
  unit: Unit
  profiles?: PaintProfile[]
  updateLoading?: boolean
  deleteLoading?: boolean
}>()

const emit = defineEmits<{
  update: [unitId: number, payload: UpdateUnitPayload]
  delete: [unitId: number]
}>()

const showEditModal   = ref(false)
const showDeleteModal = ref(false)

const linkedProfile = computed(() =>
  props.profiles?.find((p) => p.id === props.unit.paintProfileId) ?? null,
)

// Find a representative swatch color from linked profile's first step
const swatchColor = computed(() => {
  const firstStep = linkedProfile.value?.steps?.[0]
  return firstStep?.colorHex ?? null
})

function handleUpdate(payload: UpdateUnitPayload) {
  emit('update', props.unit.id, payload)
  showEditModal.value = false
}
</script>

<template>
  <div class="unit-row">
    <div class="unit-row__main">
      <div class="unit-row__identity">
        <span class="unit-row__qty">×{{ unit.quantity }}</span>
        <span class="unit-row__name">{{ unit.name }}</span>
      </div>

      <div class="unit-row__status">
        <StatusChip :status="unit.status" />
      </div>

      <div class="unit-row__profile" v-if="unit.paintProfileId">
        <ColorSwatch :color-hex="swatchColor" :size="16" />
        <span class="unit-row__profile-name">{{ linkedProfile?.name }}</span>
      </div>
    </div>

    <div class="unit-row__actions">
      <BaseButton variant="ghost" size="sm" @click="showEditModal = true">Edit</BaseButton>
      <BaseButton variant="ghost" size="sm" @click="showDeleteModal = true">✕</BaseButton>
    </div>

    <!-- Edit modal -->
    <BaseModal
      v-if="showEditModal"
      :title="`EDIT: ${unit.name}`"
      @close="showEditModal = false"
    >
      <UnitForm
        :initial="unit"
        :profiles="profiles"
        :loading="updateLoading"
        @submit="handleUpdate"
        @cancel="showEditModal = false"
      />
    </BaseModal>

    <!-- Delete confirm -->
    <ConfirmModal
      v-if="showDeleteModal"
      title="REMOVE UNIT"
      :message="`Remove '${unit.name}' (×${unit.quantity}) from this project?`"
      confirm-label="Remove Unit"
      :loading="deleteLoading"
      @confirm="emit('delete', unit.id)"
      @cancel="showDeleteModal = false"
    />
  </div>
</template>

<style lang="scss" scoped>
@use "@/styles/tokens" as *;
@use "@/styles/mixins" as *;

.unit-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: $spacing-4;
  padding: $spacing-3 $spacing-4;
  background: $color-bg-surface;
  border: 1px solid $color-border-default;
  border-radius: $radius-md;
  transition: border-color 0.15s;

  &:hover {
    border-color: $color-border-strong;
  }

  &__main {
    display: flex;
    align-items: center;
    gap: $spacing-4;
    flex: 1;
    min-width: 0;
  }

  &__identity {
    display: flex;
    align-items: baseline;
    gap: $spacing-2;
    min-width: 0;
  }

  &__qty {
    font-family: $font-mono;
    font-size: $font-mono-size;
    color: $color-text-secondary;
    flex-shrink: 0;
  }

  &__name {
    font-family: $font-heading;
    font-size: $font-heading-sm;
    font-weight: $weight-semibold;
    color: $color-text-primary;
    @include truncate;
  }

  &__status {
    flex-shrink: 0;
  }

  &__profile {
    display: flex;
    align-items: center;
    gap: $spacing-1;
    min-width: 0;
  }

  &__profile-name {
    font-size: $font-body-sm;
    color: $color-text-secondary;
    @include truncate;
  }

  &__actions {
    display: flex;
    align-items: center;
    gap: $spacing-1;
    flex-shrink: 0;
  }
}
</style>
