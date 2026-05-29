<script setup lang="ts">
import type { HobbyProject } from '@/types'
import SegmentedProgress from '@/components/SegmentedProgress.vue'
import { RouterLink } from 'vue-router'

defineProps<{ project: HobbyProject }>()
</script>

<template>
  <RouterLink
    :to="{ name: 'project-detail', params: { projectId: project.id } }"
    class="project-card"
  >
    <div class="project-card__header">
      <h2 class="project-card__name">{{ project.name }}</h2>
      <div class="project-card__tags">
        <span v-if="project.gameSystem" class="project-card__tag">{{ project.gameSystem }}</span>
        <span v-if="project.faction"    class="project-card__tag">{{ project.faction }}</span>
      </div>
    </div>

    <div v-if="project.progress" class="project-card__progress">
      <SegmentedProgress :progress="project.progress" :compact="true" />
      <div class="project-card__progress-meta">
        <span class="project-card__unit-count">{{ project.progress.totalUnits }} units</span>
        <span class="project-card__complete-count">
          {{ project.progress.byStatus['complete'] ?? 0 }} complete
        </span>
      </div>
    </div>

    <p v-if="project.description" class="project-card__description">
      {{ project.description }}
    </p>

    <div class="project-card__footer">
      <span class="project-card__arrow" aria-hidden="true">VIEW PROJECT →</span>
    </div>
  </RouterLink>
</template>

<style lang="scss" scoped>
@use "@/styles/tokens" as *;
@use "@/styles/mixins" as *;

.project-card {
  @include panel-frame($notch: true);
  display: flex;
  flex-direction: column;
  gap: $spacing-3;
  padding: $spacing-5;
  text-decoration: none;
  color: inherit;
  transition: border-color 0.15s ease, box-shadow 0.15s ease;

  &:hover {
    border-color: $color-accent-amber;
    box-shadow: $shadow-mid;
    text-decoration: none;
  }

  &__header {
    display: flex;
    flex-direction: column;
    gap: $spacing-1;
  }

  &__name {
    font-family: $font-heading;
    font-size: $font-heading-md;
    font-weight: $weight-bold;
    color: $color-text-primary;
    @include truncate;
  }

  &__tags {
    display: flex;
    flex-wrap: wrap;
    gap: $spacing-1;
  }

  &__tag {
    @include section-label;
    font-size: 10px;
    padding: 1px $spacing-2;
    border: 1px solid $color-border-default;
    border-radius: $radius-sm;
  }

  &__progress {
    display: flex;
    flex-direction: column;
    gap: $spacing-1;

    .segmented-progress {
      flex: 1;
    }
  }

  &__progress-meta {
    display: flex;
    justify-content: space-between;
  }

  &__unit-count,
  &__complete-count {
    font-family: $font-mono;
    font-size: 11px;
    color: $color-text-secondary;
  }

  &__complete-count {
    color: $color-status-complete;
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
    margin-top: auto;
    padding-top: $spacing-2;
    border-top: 1px solid $color-border-default;
  }

  &__arrow {
    @include section-label;
    font-size: 10px;
    color: $color-accent-amber;
  }
}
</style>
