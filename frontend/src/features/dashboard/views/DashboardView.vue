<script setup lang="ts">
import { computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { fetchProjects } from '@/api/projects'
import { useAuthStore } from '@/stores/auth'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import ErrorBanner from '@/components/ErrorBanner.vue'
import EmptyState from '@/components/EmptyState.vue'
import ProjectCard from '@/features/projects/components/ProjectCard.vue'
import BaseButton from '@/components/BaseButton.vue'
import { RouterLink } from 'vue-router'

const authStore = useAuthStore()

const { data: projects, isLoading, error } = useQuery({
  queryKey: ['projects'],
  queryFn: fetchProjects,
})

const stats = computed(() => {
  if (!projects.value) return null

  const totalProjects = projects.value.length
  let totalUnits = 0
  let completeUnits = 0

  for (const p of projects.value) {
    if (p.progress) {
      totalUnits += p.progress.totalUnits
      completeUnits += p.progress.byStatus['complete'] ?? 0
    }
  }

  const completionPct = totalUnits > 0 ? Math.round((completeUnits / totalUnits) * 100) : 0

  return { totalProjects, totalUnits, completeUnits, completionPct }
})

const errorMessage = computed(() =>
  error.value instanceof Error ? error.value.message : 'Failed to load projects.',
)
</script>

<template>
  <div class="dashboard">
    <header class="page-header">
      <div>
        <h1 class="page-title">OPERATIONS DASHBOARD</h1>
        <p class="page-subtitle">
          Welcome back, <span class="dashboard__operative">{{ authStore.user?.username }}</span>
        </p>
      </div>
    </header>

    <LoadingSpinner v-if="isLoading" />
    <ErrorBanner v-else-if="error" :message="errorMessage" />

    <template v-else-if="projects">
      <!-- Stat cards -->
      <div class="row row-cols-2 row-cols-lg-4 g-3">
        <div class="col"><div class="stat-card">
          <span class="stat-card__label">ACTIVE PROJECTS</span>
          <span class="stat-card__value">{{ stats?.totalProjects ?? 0 }}</span>
        </div></div>
        <div class="col"><div class="stat-card">
          <span class="stat-card__label">TOTAL UNITS</span>
          <span class="stat-card__value">{{ stats?.totalUnits ?? 0 }}</span>
        </div></div>
        <div class="col"><div class="stat-card stat-card--accent">
          <span class="stat-card__label">COMPLETION</span>
          <span class="stat-card__value">{{ stats?.completionPct ?? 0 }}%</span>
        </div></div>
        <div class="col"><div class="stat-card">
          <span class="stat-card__label">COMPLETE UNITS</span>
          <span class="stat-card__value">{{ stats?.completeUnits ?? 0 }}</span>
        </div></div>
      </div>

      <!-- Project cards -->
      <section class="d-flex flex-column gap-3">
        <h2 class="section-title">CAMPAIGN ROSTER</h2>

        <EmptyState
          v-if="projects.length === 0"
          heading="NO PROJECTS FILED"
          message="Begin your campaign by creating your first hobby project."
        >
          <template #action>
            <RouterLink :to="{ name: 'projects' }">
              <BaseButton variant="primary">Go to Projects</BaseButton>
            </RouterLink>
          </template>
        </EmptyState>

        <div v-else class="dashboard__project-grid">
          <ProjectCard v-for="project in projects" :key="project.id" :project="project" />
        </div>
      </section>
    </template>
  </div>
</template>

<style lang="scss" scoped>
@use "@/styles/tokens" as *;
@use "@/styles/mixins" as *;

.dashboard {
  display: flex;
  flex-direction: column;
  gap: $spacing-8;

  &__operative {
    color: $color-accent-amber;
    font-family: $font-mono;
  }

  &__project-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: $spacing-4;
  }
}

// ── Stat cards ──────────────────────────────────────────────────────────────

.stat-card {
  @include panel-frame($notch: true);
  padding: $spacing-5;
  display: flex;
  flex-direction: column;
  gap: $spacing-2;
  background: linear-gradient(135deg, $color-bg-elevated 0%, $color-bg-surface 100%);

  &__label {
    @include section-label;
    font-size: 11px;
  }

  &__value {
    font-family: $font-mono;
    font-size: 28px;
    font-weight: $weight-bold;
    color: $color-text-primary;
    line-height: 1;
  }

  &--accent &__value {
    color: $color-accent-amber;
  }
}
</style>
