<script setup lang="ts">
withDefaults(
  defineProps<{
    notch?: boolean
    elevated?: boolean
    padding?: 'none' | 'sm' | 'md' | 'lg'
  }>(),
  { notch: false, elevated: false, padding: 'md' },
)
</script>

<template>
  <div :class="['card', { 'card--notch': notch, 'card--elevated': elevated }, `card--pad-${padding}`]">
    <slot />
  </div>
</template>

<style lang="scss" scoped>
@use "@/styles/tokens" as *;
@use "@/styles/mixins" as *;

.card {
  @include panel-frame;

  &--elevated {
    @include panel-frame($bg: $color-bg-elevated);
  }

  &--notch {
    @include panel-frame($notch: true);
  }

  &--pad-none  { padding: 0; }
  &--pad-sm    { padding: $spacing-3; }
  &--pad-md    { padding: $spacing-4 $spacing-5; }
  &--pad-lg    { padding: $spacing-6 $spacing-8; }
}
</style>
