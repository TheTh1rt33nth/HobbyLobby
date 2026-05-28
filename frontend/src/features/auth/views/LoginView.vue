<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuth } from '../composables/useAuth'
import BaseButton from '@/components/BaseButton.vue'
import FormField from '@/components/FormField.vue'
import ErrorBanner from '@/components/ErrorBanner.vue'

const router = useRouter()
const route  = useRoute()
const { loginMutation } = useAuth()

const form = reactive({ username: '', password: '' })
const errors = reactive({ username: '', password: '' })
const serverError = ref('')

function validate(): boolean {
  errors.username = form.username.trim() ? '' : 'Username is required'
  errors.password = form.password       ? '' : 'Password is required'
  return !errors.username && !errors.password
}

async function handleSubmit() {
  serverError.value = ''
  if (!validate()) return

  try {
    await loginMutation.mutateAsync({ username: form.username.trim(), password: form.password })
    const redirect = route.query['redirect'] as string | undefined
    router.push(redirect ?? { name: 'dashboard' })
  } catch (err: unknown) {
    serverError.value = err instanceof Error ? err.message : 'Login failed. Check your credentials.'
  }
}
</script>

<template>
  <div class="auth-page">
    <div class="auth-card">
      <div class="auth-card__header">
        <span class="auth-card__logo">HOBBYLOBBY</span>
        <span class="auth-card__tagline">OPERATIONS CONSOLE — AUTHENTICATION</span>
      </div>

      <ErrorBanner v-if="serverError" :message="serverError" />

      <form @submit.prevent="handleSubmit" novalidate class="auth-form">
        <FormField id="username" label="Username" :required="true" :error="errors.username">
          <input
            id="username"
            v-model="form.username"
            type="text"
            autocomplete="username"
            placeholder="operative handle"
          />
        </FormField>

        <FormField id="password" label="Password" :required="true" :error="errors.password">
          <input
            id="password"
            v-model="form.password"
            type="password"
            autocomplete="current-password"
            placeholder="••••••••"
          />
        </FormField>

        <BaseButton
          type="submit"
          variant="primary"
          size="lg"
          :loading="loginMutation.isPending.value"
          class="w-100"
        >
          Authenticate
        </BaseButton>
      </form>

      <p class="auth-card__footer-text">
        New operative?
        <RouterLink :to="{ name: 'register' }">Register here</RouterLink>
      </p>
    </div>
  </div>
</template>

