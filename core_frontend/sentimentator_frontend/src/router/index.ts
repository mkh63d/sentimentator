import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router';
import { useAuthStore } from '@/stores/authStore';

const routes: RouteRecordRaw[] = [
  { path: '/login', name: 'login', component: () => import('@/views/LoginPanel.vue') },
  { path: '/', redirect: { name: 'home' } },
  { path: '/list', name: 'list', component: () => import('@/views/ListOfAnalysisView.vue') },
  { path: '/home', name: 'home', component: () => import('@/views/HomePage.vue') },
  { path: '/analysis/:id', name: 'analysis', component: () => import('@/views/AnalysisView.vue') },
  { path: '/settings', name: 'settings', component: () => import('@/views/SettingsView.vue') },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) return savedPosition;

    if (to.hash) {
      return { el: to.hash, behavior: 'smooth' };
      // Smooth scrolling doesn't work ):<
      // I'll leave it here as a testament to my efforts
    }

    return { top: 0 };
  },
});

// Route guard for authentication
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();
  const isAuthenticated = checkAuthStatus(authStore); // You'll need to implement this function

  const publicPages = ['home', 'login'];
  if (!publicPages.includes(to.name as string) && !isAuthenticated) {
    // Trying to access a protected page without authentication
    return next({ name: 'login' });
  }

  // Allow access
  next();
});

// Helper function to check authentication status
function checkAuthStatus(authStore: ReturnType<typeof useAuthStore>): boolean {
  // Implement your authentication check here
  const token = authStore.accessToken;
  return token !== null && token !== undefined && token !== '';
}

export default router;
