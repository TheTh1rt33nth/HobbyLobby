<script setup lang="ts">
defineProps<{
  heading?: string
  message?: string
}>()
</script>

<template>
  <div class="empty-state">
    <div class="empty-state__icon" aria-hidden="true">▭</div>
    <p class="empty-state__heading">{{ heading ?? 'NO RECORDS FOUND' }}</p>
    <p v-if="message" class="empty-state__message">{{ message }}</p>
    <div v-if="$slots.action" class="empty-state__action">
      <slot name="action" />
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use "@/styles/tokens" as *;
@use "@/styles/mixins" as *;

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: $spacing-12 $spacing-8;
  text-align: center;
  gap: $spacing-3;

  &__icon {
    font-size: 40px;
    opacity: 0.12;
    color: $color-text-primary;
    line-height: 1;
  }

  &__heading {
    @include section-label;
    font-size: $font-heading-sm;
  }

  &__message {
    color: $color-text-secondary;
    font-size: $font-body-sm;
    max-width: 360px;
    line-height: 1.6;
  }

  &__action {
    margin-top: $spacing-2;
  }
}
</style>
