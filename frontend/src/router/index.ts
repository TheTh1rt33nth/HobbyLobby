import { createRouter, createWebHistory } from 'vue-router'

// Lazy-loaded views

const LoginView          = () => import('@/features/auth/views/LoginView.vue')
const RegisterView       = () => import('@/features/auth/views/RegisterView.vue')
const DashboardView      = () => import('@/features/dashboard/views/DashboardView.vue')
const ProjectsView       = () => import('@/features/projects/views/ProjectsView.vue')
const ProjectDetailView  = () => import('@/features/projects/views/ProjectDetailView.vue')
const PaintProfilesView  = () => import('@/features/paint-profiles/views/PaintProfilesView.vue')
const PaintProfileDetail = () => import('@/features/paint-profiles/views/PaintProfileDetailView.vue')
const AppLayout          = () => import('@/layouts/AppLayout.vue')

// Route table

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login',    name: 'login',    component: LoginView },
    { path: '/register', name: 'register', component: RegisterView },
    {
      path: '/',
      redirect: { name: 'dashboard' },
      component: AppLayout,
      meta: { requiresAuth: true },
      children: [
        {
          path: 'dashboard',
          name: 'dashboard',
          component: DashboardView,
        },
        {
          path: 'projects',
          name: 'projects',
          component: ProjectsView,
        },
        {
          path: 'projects/:projectId',
          name: 'project-detail',
          component: ProjectDetailView,
          props: (route) => ({ projectId: Number(route.params['projectId']) }),
        },
        {
          path: 'paint-profiles',
          name: 'paint-profiles',
          component: PaintProfilesView,
        },
        {
          path: 'paint-profiles/:profileId',
          name: 'paint-profile-detail',
          component: PaintProfileDetail,
          props: (route) => ({ profileId: Number(route.params['profileId']) }),
        },
      ],
    },
    { path: '/:pathMatch(.*)*', redirect: '/dashboard' },
  ],
})


router.beforeEach(async (to) => {
  if (to.meta['requiresAuth']) {
    const { useAuthStore } = await import('@/stores/auth')
    const authStore = useAuthStore()

    if (!authStore.isAuthenticated) {
      return { name: 'login', query: { redirect: to.fullPath } }
    }
  }
})


export function getRouter() {
  return router
}

export default router
