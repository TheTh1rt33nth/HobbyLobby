<script setup lang="ts">
withDefaults(
  defineProps<{
    variant?: 'primary' | 'secondary' | 'danger' | 'ghost'
    size?: 'sm' | 'md' | 'lg'
    type?: 'button' | 'submit' | 'reset'
    loading?: boolean
    disabled?: boolean
  }>(),
  {
    variant: 'primary',
    size: 'md',
    type: 'button',
    loading: false,
    disabled: false,
  },
)
</script>

<template>
  <button
    :type="type"
    :disabled="disabled || loading"
    :class="['btn', `btn--${variant}`, `btn--${size}`, { 'btn--loading': loading }]"
  >
    <span v-if="loading" class="btn__spinner" aria-hidden="true" />
    <slot />
  </button>
</template>

<style lang="scss" scoped>
@use "@/styles/tokens" as *;
@use "sass:color";
@use "@/styles/mixins" as *;

.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: $spacing-2;
  border: 1px solid transparent;
  border-radius: $radius-sm;
  font-family: $font-heading;
  font-weight: $weight-semibold;
  letter-spacing: $tracking-wide;
  text-transform: uppercase;
  cursor: pointer;
  transition: background 0.15s ease, border-color 0.15s ease, color 0.15s ease, opacity 0.15s ease;
  white-space: nowrap;
  user-select: none;

  &:disabled {
    opacity: 0.45;
    cursor: not-allowed;
  }

  // ── Sizes ──────────────────────────────────────────────────────────────────
  &--sm { padding: $spacing-1 $spacing-3; font-size: 12px; }
  &--md { padding: $spacing-2 $spacing-4; font-size: $font-heading-sm; }
  &--lg { padding: $spacing-3 $spacing-6; font-size: 16px; }

  // ── Variants ───────────────────────────────────────────────────────────────
  &--primary {
    background: $color-accent-brass;
    border-color: $color-accent-amber;
    color: $color-text-primary;

    &:hover:not(:disabled) {
      background: $color-accent-amber;
      border-color: $color-accent-amber;
    }
  }

  &--secondary {
    background: $color-bg-elevated;
    border-color: $color-border-default;
    color: $color-text-primary;

    &:hover:not(:disabled) {
      border-color: $color-accent-amber;
      color: $color-accent-amber;
    }
  }

  &--danger {
    background: $color-accent-red;
    border-color: color.adjust($color-accent-red, $lightness: -10%);
    color: $color-text-primary;

    &:hover:not(:disabled) {
      background: color.adjust($color-accent-red, $lightness: 8%);
    }
  }

  &--ghost {
    background: transparent;
    border-color: transparent;
    color: $color-text-secondary;

    &:hover:not(:disabled) {
      border-color: $color-accent-amber;
      color: $color-accent-amber;
    }
  }

  // ── Spinner ────────────────────────────────────────────────────────────────
  &__spinner {
    width: 12px;
    height: 12px;
    border: 2px solid rgba(255, 255, 255, 0.3);
    border-top-color: currentColor;
    border-radius: 50%;
    animation: spin 0.6s linear infinite;
  }
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
