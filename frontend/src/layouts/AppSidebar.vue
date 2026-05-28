<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

defineProps<{ showClose?: boolean }>()
const emit = defineEmits<{ close: [] }>()

const authStore = useAuthStore()
const router    = useRouter()

function handleSignOut() {
  authStore.signOut()
  router.push({ name: 'login' })
}
</script>

<template>
  <nav class="sidebar">
    <div class="sidebar__header">
      <span class="sidebar__logo">HOBBY<br />LOBBY</span>
      <span class="sidebar__tagline">OPERATIONS CONSOLE</span>
      <button v-if="showClose" class="sidebar__close" type="button" aria-label="Close navigation" @click="emit('close')">✕</button>
    </div>

    <ul class="sidebar__nav" role="list">
      <li>
        <RouterLink :to="{ name: 'dashboard' }" class="sidebar__link" active-class="sidebar__link--active" @click="emit('close')">
          <span class="sidebar__link-icon" aria-hidden="true">▸</span>
          Dashboard
        </RouterLink>
      </li>
      <li>
        <RouterLink :to="{ name: 'projects' }" class="sidebar__link" active-class="sidebar__link--active" @click="emit('close')">
          <span class="sidebar__link-icon" aria-hidden="true">▸</span>
          Projects
        </RouterLink>
      </li>
      <li>
        <RouterLink :to="{ name: 'paint-profiles' }" class="sidebar__link" active-class="sidebar__link--active" @click="emit('close')">
          <span class="sidebar__link-icon" aria-hidden="true">▸</span>
          Paint Profiles
        </RouterLink>
      </li>
    </ul>

    <div class="sidebar__footer">
      <div v-if="authStore.user" class="sidebar__user">
        <span class="sidebar__user-label">OPERATIVE</span>
        <span class="sidebar__user-name">{{ authStore.user.username }}</span>
      </div>
      <button class="sidebar__signout" @click="handleSignOut" type="button">
        SIGN OUT
      </button>
    </div>
  </nav>
</template>

<style lang="scss" scoped>
@use "@/styles/tokens" as *;

.sidebar {
  background: $color-bg-surface;
  border-right: 1px solid $color-border-default;
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 100vh;
  overflow-y: auto;

  &__header {
    position: relative;
    padding: $spacing-6 $spacing-4 $spacing-4;
    border-bottom: 1px solid $color-border-default;
    background: linear-gradient(180deg, $color-bg-elevated 0%, $color-bg-surface 100%);
  }

  &__close {
    position: absolute;
    top: $spacing-3;
    right: $spacing-3;
    background: none;
    border: none;
    color: $color-text-secondary;
    font-size: 14px;
    cursor: pointer;
    padding: $spacing-1;
    line-height: 1;

    &:hover {
      color: $color-text-primary;
    }
  }

  &__logo {
    display: block;
    font-family: $font-heading;
    font-size: 22px;
    font-weight: $weight-bold;
    letter-spacing: $tracking-wider;
    color: $color-accent-amber;
    line-height: 1.1;
  }

  &__tagline {
    display: block;
    font-family: $font-mono;
    font-size: 9px;
    letter-spacing: $tracking-wider;
    color: $color-text-disabled;
    margin-top: $spacing-1;
  }

  &__nav {
    list-style: none;
    flex: 1;
    padding: $spacing-4 0;
  }

  &__link {
    display: flex;
    align-items: center;
    gap: $spacing-2;
    padding: $spacing-3 $spacing-4;
    font-family: $font-heading;
    font-size: $font-heading-sm;
    font-weight: $weight-semibold;
    letter-spacing: $tracking-wide;
    text-transform: uppercase;
    color: $color-text-secondary;
    text-decoration: none;
    transition: color 0.15s ease, background 0.15s ease;
    border-left: 2px solid transparent;

    &:hover {
      color: $color-text-primary;
      background: rgba($color-accent-amber, 0.05);
      border-left-color: $color-accent-amber;
      text-decoration: none;
    }

    &--active {
      color: $color-accent-amber;
      border-left-color: $color-accent-amber;
    }
  }

  &__link-icon {
    font-size: 10px;
    opacity: 0.5;
  }

  &__footer {
    padding: $spacing-4;
    border-top: 1px solid $color-border-default;
  }

  &__user {
    display: flex;
    flex-direction: column;
    margin-bottom: $spacing-3;
  }

  &__user-label {
    font-family: $font-mono;
    font-size: 9px;
    letter-spacing: $tracking-wider;
    text-transform: uppercase;
    color: $color-text-disabled;
  }

  &__user-name {
    font-family: $font-mono;
    font-size: $font-mono-size;
    color: $color-text-primary;
  }

  &__signout {
    width: 100%;
    padding: $spacing-2 $spacing-3;
    background: transparent;
    border: 1px solid $color-border-default;
    border-radius: $radius-sm;
    color: $color-text-secondary;
    font-family: $font-heading;
    font-size: $font-heading-sm;
    font-weight: $weight-semibold;
    letter-spacing: $tracking-wide;
    text-transform: uppercase;
    cursor: pointer;
    transition: border-color 0.15s ease, color 0.15s ease;

    &:hover {
      border-color: $color-accent-red;
      color: $color-accent-red;
    }
  }
}
</style>
