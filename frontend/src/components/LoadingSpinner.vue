<template>
  <div class="spinner" role="status" aria-label="Loading">
    <div class="spinner__ring" />
  </div>
</template>

<style lang="scss" scoped>
@use "@/styles/tokens" as *;

.spinner {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: $spacing-8;

  &__ring {
    width: 32px;
    height: 32px;
    border: 3px solid $color-border-default;
    border-top-color: $color-accent-amber;
    border-radius: 50%;
    animation: spin 0.7s linear infinite;
  }
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
