<script setup lang="ts">
import { ref } from 'vue'

defineProps<{ message: string }>()

const dismissed = ref(false)
</script>

<template>
  <div v-if="!dismissed" class="error-banner" role="alert">
    <span class="error-banner__icon" aria-hidden="true">⚠</span>
    <span class="error-banner__message">{{ message }}</span>
    <button class="error-banner__dismiss" @click="dismissed = true" aria-label="Dismiss error" type="button">✕</button>
  </div>
</template>

<style lang="scss" scoped>
@use "@/styles/tokens" as *;
@use "@/styles/mixins" as *;

.error-banner {
  @include hazard-stripe($color-accent-red, 0.07);
  display: flex;
  align-items: center;
  gap: $spacing-3;
  padding: $spacing-3 $spacing-4;
  border: 1px solid rgba($color-accent-red, 0.3);
  border-radius: $radius-md;
  color: $color-text-primary;
  margin-bottom: $spacing-4;

  &__icon {
    color: $color-accent-red;
    font-size: 16px;
    flex-shrink: 0;
  }

  &__message {
    flex: 1;
    font-size: $font-body-sm;
    line-height: 1.5;
  }

  &__dismiss {
    background: none;
    border: none;
    color: $color-text-secondary;
    cursor: pointer;
    font-size: 12px;
    padding: $spacing-1;
    flex-shrink: 0;
    transition: color 0.15s;

    &:hover { color: $color-text-primary; }
  }
}
</style>
