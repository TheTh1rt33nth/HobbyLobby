<script setup lang="ts">
import { reactive, watch } from 'vue'
import type { Unit, CreateUnitPayload, PaintProfile } from '@/types'
import { HOBBY_STATUS_ORDER, HOBBY_STATUS_LABELS } from '@/types'
import FormField from '@/components/FormField.vue'
import BaseButton from '@/components/BaseButton.vue'

const props = defineProps<{
  initial?: Unit | null
  loading?: boolean
  profiles?: PaintProfile[]
}>()

const emit = defineEmits<{
  submit: [payload: CreateUnitPayload]
  cancel: []
}>()

const form = reactive<CreateUnitPayload>({
  name: '',
  quantity: 1,
  status: 'unassembled',
  notes: null,
  paintProfileId: null,
})

const errors = reactive({ name: '', quantity: '' })

watch(
  () => props.initial,
  (val) => {
    if (val) {
      form.name           = val.name
      form.quantity       = val.quantity
      form.status         = val.status
      form.notes          = val.notes
      form.paintProfileId = val.paintProfileId
    }
  },
  { immediate: true },
)

function validate(): boolean {
  errors.name     = form.name.trim()  ? '' : 'Unit name is required'
  errors.quantity = form.quantity > 0 ? '' : 'Quantity must be at least 1'
  return !errors.name && !errors.quantity
}

function handleSubmit() {
  if (!validate()) return
  emit('submit', {
    name:           form.name.trim(),
    quantity:       Number(form.quantity),
    status:         form.status,
    notes:          form.notes?.trim() || null,
    paintProfileId: form.paintProfileId ? Number(form.paintProfileId) : null,
  })
}
</script>

<template>
  <form @submit.prevent="handleSubmit" novalidate class="unit-form">
    <FormField id="unit-name" label="Unit Name" :required="true" :error="errors.name">
      <input id="unit-name" v-model="form.name" type="text" placeholder="e.g. Chaos Space Marines Squad" />
    </FormField>

    <div class="row g-3">
      <div class="col-6">
        <FormField id="unit-qty" label="Quantity" :required="true" :error="errors.quantity">
          <input id="unit-qty" v-model.number="form.quantity" type="number" min="1" max="999" />
        </FormField>
      </div>
      <div class="col-6">
        <FormField id="unit-status" label="Status" :required="true">
          <select id="unit-status" v-model="form.status">
            <option v-for="s in HOBBY_STATUS_ORDER" :key="s" :value="s">
              {{ HOBBY_STATUS_LABELS[s] }}
            </option>
          </select>
        </FormField>
      </div>
    </div>

    <FormField v-if="profiles?.length" id="unit-profile" label="Paint Profile">
      <select id="unit-profile" v-model="form.paintProfileId">
        <option :value="null">— None —</option>
        <option v-for="p in profiles" :key="p.id" :value="p.id">{{ p.name }}</option>
      </select>
    </FormField>

    <FormField id="unit-notes" label="Notes">
      <textarea id="unit-notes" v-model="form.notes" placeholder="Assembly notes, special instructions…" />
    </FormField>

    <div class="d-flex justify-content-end gap-2 pt-2">
      <BaseButton type="button" variant="secondary" @click="emit('cancel')">Cancel</BaseButton>
      <BaseButton type="submit" variant="primary" :loading="loading">
        {{ initial ? 'Save Changes' : 'Add Unit' }}
      </BaseButton>
    </div>
  </form>
</template>

<style lang="scss" scoped>
@use "@/styles/tokens" as *;

.unit-form {
  display: flex;
  flex-direction: column;
  gap: $spacing-4;
}
</style>
