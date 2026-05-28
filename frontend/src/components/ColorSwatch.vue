<script setup lang="ts">
withDefaults(
  defineProps<{
    colorHex: string | null
    size?: number
  }>(),
  { size: 24 },
)
</script>

<template>
  <span
    class="swatch"
    :style="{
      width: `${size}px`,
      height: `${size}px`,
      background: colorHex ?? 'transparent',
      borderColor: colorHex ? 'rgba(255,255,255,0.15)' : 'var(--border)',
    }"
    :title="colorHex ?? 'No color'"
    :aria-label="colorHex ? `Color: ${colorHex}` : 'No color assigned'"
  />
</template>

<style lang="scss" scoped>
@use "@/styles/tokens" as *;

.swatch {
  display: inline-block;
  border-radius: 50%;
  border: 1px solid $color-border-default;
  flex-shrink: 0;
  vertical-align: middle;

  --border: #{$color-border-default};

  // Placeholder pattern when no color
  &[style*="background: transparent"] {
    background-image: repeating-linear-gradient(
      -45deg,
      $color-border-default,
      $color-border-default 1px,
      transparent 1px,
      transparent 4px
    ) !important;
  }
}
</style>
