<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'

const props = withDefaults(
  defineProps<{
    title?: string
    maxWidth?: string
  }>(),
  { maxWidth: '520px' },
)

const emit = defineEmits<{ close: [] }>()

const dialogEl = ref<HTMLDivElement | null>(null)

// Track whether the mousedown started on the backdrop itself (not inside the modal).
// This prevents drag-selecting text inside the modal from closing it when the
// mouse is released over the backdrop.
let mousedownOnBackdrop = false

function handleBackdropMousedown(e: MouseEvent) {
  mousedownOnBackdrop = e.target === e.currentTarget
}

function handleBackdropClick(e: MouseEvent) {
  if (mousedownOnBackdrop && e.target === e.currentTarget) emit('close')
  mousedownOnBackdrop = false
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') emit('close')
}

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
  // Trap focus: focus the first focusable element inside the modal
  dialogEl.value?.querySelector<HTMLElement>('input, button, [tabindex]')?.focus()
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})
</script>

<template>
  <Teleport to="body">
    <Transition name="fade">
        <div class="modal-backdrop" @mousedown="handleBackdropMousedown" @click="handleBackdropClick">
        <div
          ref="dialogEl"
          class="modal-panel"
          role="dialog"
          aria-modal="true"
          :aria-label="title"
          :style="{ maxWidth: maxWidth }"
          @click.stop
        >
          <div v-if="title || $slots.header" class="modal-panel__header">
            <slot name="header">
              <h2 class="modal-panel__title">{{ title }}</h2>
            </slot>
            <button class="modal-panel__close" @click="emit('close')" aria-label="Close" type="button">✕</button>
          </div>
          <div class="modal-panel__body">
            <slot />
          </div>
          <div v-if="$slots.footer" class="modal-panel__footer">
            <slot name="footer" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style lang="scss" scoped>
@use "@/styles/tokens" as *;
@use "@/styles/mixins" as *;

.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.65);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: $spacing-4;
  z-index: 100;
}

.modal-panel {
  @include panel-frame($bg: $color-bg-elevated, $notch: true);
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;

  &__header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: $spacing-4 $spacing-5;
    border-bottom: 1px solid $color-border-default;
    @include panel-header-gradient;
  }

  &__title {
    font-family: $font-heading;
    font-size: $font-heading-md;
    font-weight: $weight-bold;
    letter-spacing: $tracking-wide;
    color: $color-text-primary;
  }

  &__close {
    background: none;
    border: none;
    color: $color-text-secondary;
    font-size: 16px;
    cursor: pointer;
    padding: $spacing-1;
    line-height: 1;
    transition: color 0.15s;

    &:hover { color: $color-text-primary; }
  }

  &__body {
    padding: $spacing-5;
  }

  &__footer {
    padding: $spacing-4 $spacing-5;
    border-top: 1px solid $color-border-default;
    display: flex;
    justify-content: flex-end;
    gap: $spacing-2;
  }
}
</style>
