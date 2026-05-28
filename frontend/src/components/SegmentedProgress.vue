<script setup lang="ts">
import { computed } from 'vue'
import type { ProjectProgress } from '@/types'
import { HOBBY_STATUS_ORDER, HOBBY_STATUS_LABELS } from '@/types'

const props = defineProps<{
  progress: ProjectProgress
  compact?: boolean
}>()

const segments = computed(() =>
  HOBBY_STATUS_ORDER.map((status) => ({
    status,
    label: HOBBY_STATUS_LABELS[status],
    count: props.progress.byStatus[status] ?? 0,
    width: props.progress.totalUnits > 0
      ? ((props.progress.byStatus[status] ?? 0) / props.progress.totalUnits) * 100
      : 0,
  })).filter((s) => s.count > 0),
)

const isEmpty = computed(() => props.progress.totalUnits === 0)
</script>

<template>
  <div class="segmented-progress" :class="{ 'segmented-progress--compact': compact }">
    <div class="segmented-progress__bar" role="img" :aria-label="`Progress: ${progress.totalUnits} total units`">
      <div
        v-if="isEmpty"
        class="segmented-progress__empty"
      />
      <div
        v-for="seg in segments"
        :key="seg.status"
        :class="['segmented-progress__segment', `segmented-progress__segment--${seg.status}`]"
        :style="{ width: `${seg.width}%` }"
        :title="`${seg.label}: ${seg.count}`"
      />
    </div>
    <div v-if="!compact" class="segmented-progress__legend">
      <span
        v-for="seg in segments"
        :key="seg.status"
        class="segmented-progress__legend-item"
      >
        <span
          :class="['segmented-progress__legend-dot', `segmented-progress__segment--${seg.status}`]"
        />
        {{ seg.label }} · {{ seg.count }}
      </span>
      <span v-if="isEmpty" class="segmented-progress__legend-item segmented-progress__legend-item--empty">
        NO UNITS
      </span>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use "@/styles/tokens" as *;

.segmented-progress {
  &__bar {
    height: 8px;
    border-radius: $radius-sm;
    background: $color-bg-elevated;
    overflow: hidden;
    display: flex;
    border: 1px solid $color-border-default;
  }

  &__empty {
    width: 100%;
    background: $color-bg-elevated;
  }

  &__segment {
    height: 100%;
    transition: width 0.3s ease;

    &--unassembled { background: $color-status-unassembled; }
    &--assembled   { background: $color-status-assembled; }
    &--primed      { background: $color-status-primed; }
    &--base_coated { background: $color-status-base-coated; }
    &--painted     { background: $color-status-painted; }
    &--based       { background: $color-status-based; }
    &--complete    { background: $color-status-complete; }
  }

  &__legend {
    display: flex;
    flex-wrap: wrap;
    gap: $spacing-2 $spacing-4;
    margin-top: $spacing-2;
  }

  &__legend-item {
    display: flex;
    align-items: center;
    gap: $spacing-1;
    font-family: $font-mono;
    font-size: 11px;
    color: $color-text-secondary;

    &--empty {
      font-family: $font-heading;
      font-size: $font-heading-sm;
      letter-spacing: $tracking-wider;
    }
  }

  &__legend-dot {
    display: inline-block;
    width: 8px;
    height: 8px;
    border-radius: 50%;
  }

  &--compact &__bar {
    height: 5px;
  }
}
</style>
