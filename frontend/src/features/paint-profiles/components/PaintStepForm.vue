<script setup lang="ts">
import { reactive, watch, computed } from 'vue'
import type { PaintStep, CreatePaintStepPayload } from '@/types'
import FormField from '@/components/FormField.vue'
import BaseButton from '@/components/BaseButton.vue'
import ColorSwatch from '@/components/ColorSwatch.vue'

const props = defineProps<{
  initial?: PaintStep | null
  nextOrder?: number
  loading?: boolean
}>()

const emit = defineEmits<{
  submit: [payload: CreatePaintStepPayload]
  cancel: []
}>()

const form = reactive<CreatePaintStepPayload>({
  stepOrder:         1,
  paintName:         '',
  brand:             null,
  paintType:         null,
  applicationMethod: null,
  colorHex:          null,
  notes:             null,
})

const errors = reactive({ paintName: '', stepOrder: '' })

watch(
  () => props.initial,
  (val) => {
    if (val) {
      form.stepOrder         = val.stepOrder
      form.paintName         = val.paintName
      form.brand             = val.brand
      form.paintType         = val.paintType
      form.applicationMethod = val.applicationMethod
      form.colorHex          = val.colorHex
      form.notes             = val.notes
    } else if (props.nextOrder != null) {
      form.stepOrder = props.nextOrder
    }
  },
  { immediate: true },
)

// Live preview of entered hex value
const previewColor = computed(() => {
  const v = form.colorHex?.trim() ?? ''
  return /^#[0-9a-fA-F]{3,6}$/.test(v) ? v : null
})

function validate(): boolean {
  errors.paintName = form.paintName.trim()   ? '' : 'Paint name is required'
  errors.stepOrder = form.stepOrder > 0      ? '' : 'Step order must be positive'
  return !errors.paintName && !errors.stepOrder
}

function handleSubmit() {
  if (!validate()) return
  emit('submit', {
    stepOrder:         Number(form.stepOrder),
    paintName:         form.paintName.trim(),
    brand:             form.brand?.trim()             || null,
    paintType:         form.paintType?.trim()         || null,
    applicationMethod: form.applicationMethod?.trim() || null,
    colorHex:          previewColor.value             || null,
    notes:             form.notes?.trim()             || null,
  })
}
</script>

<template>
  <form @submit.prevent="handleSubmit" novalidate class="step-form">
    <div class="row g-3">
      <div class="col-auto" style="min-width: 110px">
        <FormField id="step-order" label="Step #" :required="true" :error="errors.stepOrder">
          <input id="step-order" v-model.number="form.stepOrder" type="number" min="1" />
        </FormField>
      </div>
      <div class="col">
        <FormField id="step-paint" label="Paint Name" :required="true" :error="errors.paintName">
          <input id="step-paint" v-model="form.paintName" type="text" placeholder="e.g. Abaddon Black" />
        </FormField>
      </div>
    </div>

    <div class="row g-3">
      <div class="col-12 col-sm-4">
        <FormField id="step-brand" label="Brand">
          <input id="step-brand" v-model="form.brand" type="text" placeholder="e.g. Citadel" />
        </FormField>
      </div>
      <div class="col-12 col-sm-4">
        <FormField id="step-type" label="Paint Type">
          <input id="step-type" v-model="form.paintType" type="text" placeholder="e.g. Base, Layer, Shade…" />
        </FormField>
      </div>
      <div class="col-12 col-sm-4">
        <FormField id="step-method" label="Application">
          <input id="step-method" v-model="form.applicationMethod" type="text" placeholder="e.g. Drybrushing" />
        </FormField>
      </div>
    </div>

    <FormField id="step-color" label="Color Hex" hint="e.g. #2b2b2b">
      <div class="d-flex align-items-center gap-2">
        <input id="step-color" v-model="form.colorHex" type="text" placeholder="#rrggbb" />
        <ColorSwatch :color-hex="previewColor" :size="32" />
      </div>
    </FormField>

    <FormField id="step-notes" label="Notes">
      <textarea id="step-notes" v-model="form.notes" placeholder="Thin paint to milk consistency, use wet palette…" />
    </FormField>

    <div class="d-flex justify-content-end gap-2 pt-2">
      <BaseButton type="button" variant="secondary" @click="emit('cancel')">Cancel</BaseButton>
      <BaseButton type="submit" variant="primary" :loading="loading">
        {{ initial ? 'Save Step' : 'Add Step' }}
      </BaseButton>
    </div>
  </form>
</template>

<style lang="scss" scoped>
@use "@/styles/tokens" as *;

.step-form {
  display: flex;
  flex-direction: column;
  gap: $spacing-4;
}
</style>
