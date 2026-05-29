<script setup lang="ts">
import { reactive, watch } from 'vue'
import type { HobbyProject, CreateProjectPayload } from '@/types'
import FormField from '@/components/FormField.vue'
import BaseButton from '@/components/BaseButton.vue'

const props = defineProps<{
  initial?: HobbyProject | null
  loading?: boolean
}>()

const emit = defineEmits<{
  submit: [payload: CreateProjectPayload]
  cancel: []
}>()

const form = reactive<CreateProjectPayload>({
  name: '',
  description: null,
  gameSystem: null,
  faction: null,
})

const errors = reactive({ name: '' })

watch(
  () => props.initial,
  (val) => {
    if (val) {
      form.name        = val.name
      form.description = val.description
      form.gameSystem  = val.gameSystem
      form.faction     = val.faction
    }
  },
  { immediate: true },
)

function validate(): boolean {
  errors.name = form.name.trim() ? '' : 'Project name is required'
  return !errors.name
}

function handleSubmit() {
  if (!validate()) return
  emit('submit', {
    name:        form.name.trim(),
    description: form.description?.trim() || null,
    gameSystem:  form.gameSystem?.trim()  || null,
    faction:     form.faction?.trim()     || null,
  })
}
</script>

<template>
  <form @submit.prevent="handleSubmit" novalidate class="project-form">
    <FormField id="proj-name" label="Project Name" :required="true" :error="errors.name">
      <input id="proj-name" v-model="form.name" type="text" placeholder="e.g. Death Guard Battleforce" />
    </FormField>

    <div class="row g-3">
      <div class="col-6">
        <FormField id="proj-game-system" label="Game System">
          <input id="proj-game-system" v-model="form.gameSystem" type="text" placeholder="e.g. Warhammer 40,000" />
        </FormField>
      </div>
      <div class="col-6">
        <FormField id="proj-faction" label="Faction">
          <input id="proj-faction" v-model="form.faction" type="text" placeholder="e.g. Chaos Space Marines" />
        </FormField>
      </div>
    </div>

    <FormField id="proj-description" label="Description">
      <textarea id="proj-description" v-model="form.description" placeholder="Campaign notes, painting goals…" />
    </FormField>

    <div class="d-flex justify-content-end gap-2 pt-2">
      <BaseButton type="button" variant="secondary" @click="emit('cancel')">Cancel</BaseButton>
      <BaseButton type="submit" variant="primary" :loading="loading">
        {{ initial ? 'Save Changes' : 'Create Project' }}
      </BaseButton>
    </div>
  </form>
</template>

<style lang="scss" scoped>
@use "@/styles/tokens" as *;

.project-form {
  display: flex;
  flex-direction: column;
  gap: $spacing-4;
}
</style>
