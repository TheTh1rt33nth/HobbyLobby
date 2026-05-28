<script setup lang="ts">
import { ref } from 'vue'
import AppSidebar from './AppSidebar.vue'

const navOpen = ref(false)
</script>

<template>
  <div class="app-layout d-flex min-vh-100">
    <!-- Desktop sidebar — hidden below lg breakpoint -->
    <aside class="app-sidebar d-none d-lg-flex flex-column flex-shrink-0">
      <AppSidebar />
    </aside>

    <!-- Mobile offcanvas -->
    <Teleport to="body">
      <Transition name="offcanvas-fade">
        <div v-if="navOpen" class="offcanvas-overlay" @click="navOpen = false" />
      </Transition>
      <div class="offcanvas-nav" :class="{ 'offcanvas-nav--open': navOpen }" :aria-hidden="!navOpen">
        <AppSidebar :show-close="true" @close="navOpen = false" />
      </div>
    </Teleport>

    <!-- Content column -->
    <div class="d-flex flex-column flex-grow-1 min-w-0">
      <!-- Mobile top bar -->
      <header class="app-topbar d-flex d-lg-none align-items-center px-3">
        <button class="app-topbar__toggle" type="button" aria-label="Open navigation" @click="navOpen = true">
          <span /><span /><span />
        </button>
        <span class="app-topbar__logo">HOBBY LOBBY</span>
      </header>

      <!-- Page content — Bootstrap container centres on ultrawide -->
      <main class="flex-grow-1">
        <div class="container-xxl py-4 px-3 px-lg-4">
          <RouterView />
        </div>
      </main>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use "@/styles/tokens" as *;

.app-layout {
  background: $color-bg-base;
}

.app-sidebar {
  width: $sidebar-width;
  flex-shrink: 0;
}

// Mobile top bar
.app-topbar {
  height: 52px;
  background: $color-bg-surface;
  border-bottom: 1px solid $color-border-default;
  gap: $spacing-3;
  flex-shrink: 0;

  &__toggle {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    width: 22px;
    height: 16px;
    background: none;
    border: none;
    cursor: pointer;
    padding: 0;
    flex-shrink: 0;

    span {
      display: block;
      height: 2px;
      background: $color-text-primary;
      border-radius: 1px;
    }
  }

  &__logo {
    font-family: $font-heading;
    font-size: 18px;
    font-weight: $weight-bold;
    letter-spacing: $tracking-wider;
    color: $color-accent-amber;
  }
}

// Mobile offcanvas
.offcanvas-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  z-index: 200;
}

.offcanvas-nav {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  width: $sidebar-width;
  z-index: 201;
  transform: translateX(-100%);
  transition: transform 0.25s ease;

  &--open {
    transform: translateX(0);
  }
}

// Transition for the overlay
.offcanvas-fade-enter-active,
.offcanvas-fade-leave-active {
  transition: opacity 0.25s ease;
}
.offcanvas-fade-enter-from,
.offcanvas-fade-leave-to {
  opacity: 0;
}
</style>
