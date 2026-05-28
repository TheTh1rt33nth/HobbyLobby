<script setup lang="ts">
import { ref, reactive, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '../composables/useAuth'
import BaseButton from '@/components/BaseButton.vue'
import FormField from '@/components/FormField.vue'
import ErrorBanner from '@/components/ErrorBanner.vue'

const router = useRouter()
const { registerMutation } = useAuth()

const form = reactive({ username: '', email: '', password: '' })
const errors = reactive({ username: '', email: '', password: '' })
const serverError = ref('')
const success = ref(false)

let redirectTimer: ReturnType<typeof setTimeout> | null = null
onUnmounted(() => { if (redirectTimer !== null) clearTimeout(redirectTimer) })

function validate(): boolean {
  errors.username = form.username.trim() ? '' : 'Username is required'
  errors.email    = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email) ? '' : 'Valid email required'
  errors.password = form.password.length >= 8 ? '' : 'Password must be at least 8 characters'
  return !errors.username && !errors.email && !errors.password
}

async function handleSubmit() {
  serverError.value = ''
  if (!validate()) return

  try {
    await registerMutation.mutateAsync({
      username: form.username.trim(),
      email: form.email.trim(),
      password: form.password,
    })
    success.value = true
    redirectTimer = setTimeout(() => router.push({ name: 'login' }), 1500)
  } catch (err: unknown) {
    serverError.value = err instanceof Error ? err.message : 'Registration failed.'
  }
}
</script>

<template>
  <div class="auth-page">
    <div class="auth-card">
      <div class="auth-card__header">
        <span class="auth-card__logo">HOBBYLOBBY</span>
        <span class="auth-card__tagline">OPERATIVE ENLISTMENT</span>
      </div>

      <div v-if="success" class="auth-card__success">
        Enlistment confirmed. Redirecting to authentication…
      </div>

      <template v-else>
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

          <FormField id="email" label="Email" :required="true" :error="errors.email">
            <input
              id="email"
              v-model="form.email"
              type="email"
              autocomplete="email"
              placeholder="operative@sector.mil"
            />
          </FormField>

          <FormField id="password" label="Password" :required="true" :error="errors.password" hint="Minimum 8 characters">
            <input
              id="password"
              v-model="form.password"
              type="password"
              autocomplete="new-password"
              placeholder="••••••••"
            />
          </FormField>

          <BaseButton
            type="submit"
            variant="primary"
            size="lg"
            :loading="registerMutation.isPending.value"
            class="w-100"
          >
            Enlist
          </BaseButton>
        </form>
      </template>

      <p class="auth-card__footer-text">
        Already enlisted?
        <RouterLink :to="{ name: 'login' }">Sign in</RouterLink>
      </p>
    </div>
  </div>
</template>

