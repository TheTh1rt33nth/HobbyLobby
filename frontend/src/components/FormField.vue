<script setup lang="ts">
defineProps<{
  label: string
  required?: boolean
  error?: string
  hint?: string
  id: string
}>()
</script>

<template>
  <div class="form-field" :class="{ 'form-field--error': error }">
    <label :for="id" class="form-field__label">
      {{ label }}
      <span v-if="required" class="form-field__required" aria-hidden="true">*</span>
    </label>
    <slot />
    <p v-if="hint && !error" class="form-field__hint">{{ hint }}</p>
    <p v-if="error" class="form-field__error" role="alert">{{ error }}</p>
  </div>
</template>

<style lang="scss" scoped>
@use "@/styles/tokens" as *;
@use "@/styles/mixins" as *;

.form-field {
  display: flex;
  flex-direction: column;
  gap: $spacing-1;

  &__label {
    @include section-label;
    font-size: 12px;
    display: flex;
    align-items: center;
    gap: $spacing-1;
  }

  &__required {
    color: $color-accent-red;
  }

  &__hint {
    font-size: $font-body-sm;
    color: $color-text-disabled;
  }

  &__error {
    font-size: $font-body-sm;
    color: $color-accent-red;
  }

  &--error {
    :deep(input),
    :deep(select),
    :deep(textarea) {
      border-color: $color-accent-red;
    }
  }
}
</style>
