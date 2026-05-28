<script setup lang="ts">
import type { HobbyStatus } from '@/types'
import { HOBBY_STATUS_LABELS } from '@/types'

defineProps<{ status: HobbyStatus }>()
</script>

<template>
  <span :class="['chip', `chip--${status}`]">{{ HOBBY_STATUS_LABELS[status] }}</span>
</template>

<style lang="scss" scoped>
@use "@/styles/tokens" as *;

.chip {
  display: inline-block;
  padding: 2px $spacing-2;
  border-radius: $radius-sm;
  font-family: $font-heading;
  font-size: 11px;
  font-weight: $weight-semibold;
  letter-spacing: $tracking-wider;
  text-transform: uppercase;
  color: $color-bg-base;
  white-space: nowrap;

  &--unassembled { background: $color-status-unassembled; color: $color-text-primary; }
  &--assembled   { background: $color-status-assembled;   color: $color-text-primary; }
  &--primed      { background: $color-status-primed;      color: $color-text-primary; }
  &--base_coated { background: $color-status-base-coated; color: $color-bg-base; }
  &--painted     { background: $color-status-painted;     color: $color-text-primary; }
  &--based       { background: $color-status-based;       color: $color-bg-base; }
  &--complete    { background: $color-status-complete;    color: $color-bg-base; }
}
</style>
