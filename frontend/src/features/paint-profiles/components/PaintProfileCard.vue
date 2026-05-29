<script setup lang="ts">
import type { PaintProfile } from '@/types'
import ColorSwatch from '@/components/ColorSwatch.vue'
import { RouterLink } from 'vue-router'

defineProps<{ profile: PaintProfile }>()
</script>

<template>
  <RouterLink
    :to="{ name: 'paint-profile-detail', params: { profileId: profile.id } }"
    class="profile-card"
  >
    <div class="profile-card__header">
      <h2 class="profile-card__name">{{ profile.name }}</h2>
      <span v-if="profile.targetArea" class="profile-card__area">{{ profile.targetArea }}</span>
    </div>

    <p v-if="profile.description" class="profile-card__description">{{ profile.description }}</p>

    <div class="profile-card__footer">
      <div class="profile-card__swatches">
        <ColorSwatch
          v-for="step in (profile.steps ?? []).slice(0, 6)"
          :key="step.id"
          :color-hex="step.colorHex"
          :size="16"
        />
      </div>
      <span class="profile-card__step-count">
        {{ (profile.steps ?? []).length }} step{{ (profile.steps ?? []).length !== 1 ? 's' : '' }}
      </span>
    </div>
  </RouterLink>
</template>

<style lang="scss" scoped>
@use "@/styles/tokens" as *;
@use "@/styles/mixins" as *;

.profile-card {
  @include panel-frame($notch: true);
  display: flex;
  flex-direction: column;
  gap: $spacing-3;
  padding: $spacing-5;
  text-decoration: none;
  color: inherit;
  transition: border-color 0.15s, box-shadow 0.15s;

  &:hover {
    border-color: $color-accent-amber;
    box-shadow: $shadow-mid;
    text-decoration: none;
  }

  &__header {
    display: flex;
    align-items: baseline;
    justify-content: space-between;
    gap: $spacing-3;
  }

  &__name {
    font-family: $font-heading;
    font-size: $font-heading-md;
    font-weight: $weight-bold;
    color: $color-text-primary;
    @include truncate;
  }

  &__area {
    @include section-label;
    font-size: 10px;
    flex-shrink: 0;
  }

  &__description {
    font-size: $font-body-sm;
    color: $color-text-secondary;
    line-height: 1.5;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  &__footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-top: auto;
    padding-top: $spacing-2;
    border-top: 1px solid $color-border-default;
  }

  &__swatches {
    display: flex;
    gap: $spacing-1;
    align-items: center;
  }

  &__step-count {
    font-family: $font-mono;
    font-size: 11px;
    color: $color-text-secondary;
  }
}
</style>
