<script setup lang="ts">
import { computed } from 'vue'
import type { PaintStep } from '@/types'
import ColorSwatch from '@/components/ColorSwatch.vue'
import BaseButton from '@/components/BaseButton.vue'

const props = defineProps<{
  steps: PaintStep[]
  editable?: boolean
}>()

const emit = defineEmits<{
  edit: [step: PaintStep]
  delete: [stepId: number]
}>()

const sortedSteps = computed(() =>
  [...props.steps].sort((a, b) => a.stepOrder - b.stepOrder),
)
</script>

<template>
  <div class="step-list">
    <div v-if="steps.length === 0" class="step-list__empty">
      No steps added yet.
    </div>
    <div
      v-for="step in sortedSteps"
      :key="step.id"
      class="step-row"
    >
      <span class="step-row__order">{{ step.stepOrder }}</span>
      <ColorSwatch :color-hex="step.colorHex" :size="20" />
      <div class="step-row__details">
        <span class="step-row__paint-name">{{ step.paintName }}</span>
        <div class="step-row__meta">
          <span v-if="step.brand"             class="step-row__meta-item">{{ step.brand }}</span>
          <span v-if="step.paintType"         class="step-row__meta-item">{{ step.paintType }}</span>
          <span v-if="step.applicationMethod" class="step-row__meta-item">{{ step.applicationMethod }}</span>
        </div>
        <p v-if="step.notes" class="step-row__notes">{{ step.notes }}</p>
      </div>
      <div v-if="editable" class="step-row__actions">
        <BaseButton variant="ghost" size="sm" @click="emit('edit', step)">Edit</BaseButton>
        <BaseButton variant="ghost" size="sm" @click="emit('delete', step.id)">✕</BaseButton>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use "@/styles/tokens" as *;
@use "@/styles/mixins" as *;

.step-list {
  display: flex;
  flex-direction: column;
  gap: $spacing-2;
}

.step-row {
  display: flex;
  align-items: flex-start;
  gap: $spacing-3;
  padding: $spacing-3 $spacing-4;
  background: $color-bg-elevated;
  border: 1px solid $color-border-default;
  border-radius: $radius-md;
  border-left: 3px solid $color-border-strong;

  &__order {
    font-family: $font-mono;
    font-size: $font-mono-size;
    color: $color-text-secondary;
    min-width: 20px;
    flex-shrink: 0;
    padding-top: 2px;
  }

  &__details {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: $spacing-1;
  }

  &__paint-name {
    font-family: $font-heading;
    font-size: $font-heading-sm;
    font-weight: $weight-semibold;
    color: $color-text-primary;
  }

  &__meta {
    display: flex;
    flex-wrap: wrap;
    gap: $spacing-1;
  }

  &__meta-item {
    @include section-label;
    font-size: 10px;
    padding: 1px $spacing-2;
    border: 1px solid $color-border-default;
    border-radius: $radius-sm;
  }

  &__notes {
    font-size: $font-body-sm;
    color: $color-text-secondary;
    line-height: 1.4;
  }

  &__actions {
    display: flex;
    gap: $spacing-1;
    flex-shrink: 0;
  }
}

.step-list__empty {
  color: $color-text-secondary;
  font-size: $font-body-sm;
  text-align: center;
  padding: $spacing-4;
}
</style>
