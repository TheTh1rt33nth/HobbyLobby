<script setup lang="ts">
import BaseModal from './BaseModal.vue'
import BaseButton from './BaseButton.vue'

withDefaults(
  defineProps<{
    title?: string
    message: string
    confirmLabel?: string
    cancelLabel?: string
    loading?: boolean
  }>(),
  {
    title: 'Confirm Action',
    confirmLabel: 'Confirm',
    cancelLabel: 'Cancel',
    loading: false,
  },
)

const emit = defineEmits<{
  confirm: []
  cancel: []
}>()
</script>

<template>
  <BaseModal :title="title" max-width="440px" @close="emit('cancel')">
    <div class="confirm-modal">
      <div class="confirm-modal__danger-zone">
        <p class="confirm-modal__message">{{ message }}</p>
      </div>
    </div>

    <template #footer>
      <BaseButton variant="secondary" @click="emit('cancel')" :disabled="loading">
        {{ cancelLabel }}
      </BaseButton>
      <BaseButton variant="danger" @click="emit('confirm')" :loading="loading">
        {{ confirmLabel }}
      </BaseButton>
    </template>
  </BaseModal>
</template>

<style lang="scss" scoped>
@use "@/styles/tokens" as *;
@use "@/styles/mixins" as *;

.confirm-modal {
  &__danger-zone {
    @include hazard-stripe($color-accent-red, 0.06);
    border: 1px solid rgba($color-accent-red, 0.25);
    border-radius: $radius-md;
    padding: $spacing-4;
  }

  &__message {
    color: $color-text-primary;
    line-height: 1.6;
  }
}
</style>
