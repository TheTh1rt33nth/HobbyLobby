<script setup lang="ts">
import { reactive, watch } from 'vue'
import type { PaintProfile, CreatePaintProfilePayload } from '@/types'
import FormField from '@/components/FormField.vue'
import BaseButton from '@/components/BaseButton.vue'

const props = defineProps<{
  initial?: PaintProfile | null
  loading?: boolean
}>()

const emit = defineEmits<{
  submit: [payload: CreatePaintProfilePayload]
  cancel: []
}>()

const form = reactive<CreatePaintProfilePayload>({
  name: '',
  description: null,
  targetArea: null,
})

const errors = reactive({ name: '' })

watch(
  () => props.initial,
  (val) => {
    if (val) {
      form.name        = val.name
      form.description = val.description
      form.targetArea  = val.targetArea
    }
  },
  { immediate: true },
)

function validate(): boolean {
  errors.name = form.name.trim() ? '' : 'Profile name is required'
  return !errors.name
}

function handleSubmit() {
  if (!validate()) return
  emit('submit', {
    name:        form.name.trim(),
    description: form.description?.trim() || null,
    targetArea:  form.targetArea?.trim()  || null,
  })
}
</script>

<template>
  <form @submit.prevent="handleSubmit" novalidate class="profile-form">
    <FormField id="pp-name" label="Profile Name" :required="true" :error="errors.name">
      <input id="pp-name" v-model="form.name" type="text" placeholder="e.g. Death Guard Green Armour" />
    </FormField>

    <FormField id="pp-area" label="Target Area">
      <input id="pp-area" v-model="form.targetArea" type="text" placeholder="e.g. Power Armour, Trim, Eyes…" />
    </FormField>

    <FormField id="pp-desc" label="Description">
      <textarea id="pp-desc" v-model="form.description" placeholder="Short description of this scheme…" />
    </FormField>

    <div class="d-flex justify-content-end gap-2 pt-2">
      <BaseButton type="button" variant="secondary" @click="emit('cancel')">Cancel</BaseButton>
      <BaseButton type="submit" variant="primary" :loading="loading">
        {{ initial ? 'Save Changes' : 'Create Profile' }}
      </BaseButton>
    </div>
  </form>
</template>

<style lang="scss" scoped>
@use "@/styles/tokens" as *;

.profile-form {
  display: flex;
  flex-direction: column;
  gap: $spacing-4;
}
</style>
