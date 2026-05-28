import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { VueQueryPlugin, QueryClient } from '@tanstack/vue-query'
import 'bootstrap/dist/css/bootstrap.min.css'

import App from './App.vue'
import router from './router'
import { useAuthStore } from './stores/auth'
import './styles/main.scss'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30_000,       // 30s — avoids refetch storms on rapid navigation
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
})

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(VueQueryPlugin, { queryClient })

;(async () => {
  const authStore = useAuthStore()
  await authStore.bootstrap()
  app.mount('#app')
})()
