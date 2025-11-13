import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router';
import { useAuthStore } from '@/stores/authStore';

//TODO: Add route guards for authenticated routes
// and redirect to login if not authenticated.
const routes: RouteRecordRaw[] = [
  { path: '/login', name: 'login', component: () => import('@/views/LoginPanel.vue') },
  { path: '/', redirect: { name: 'home' } },
  { path: '/list', name: 'list', component: () => import('@/views/ListOfAnalysisView.vue') },
  { path: '/home', name: 'home', component: () => import('@/views/HomePage.vue') },
  { path: '/analysis/:id', name: 'analysis', component: () => import('@/views/AnalysisView.vue') },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// Route guard for authentication
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();
  const isAuthenticated = checkAuthStatus(authStore); // You'll need to implement this function
  if (!(to.name === 'home')) {
    if (to.meta.requiresAuth && !isAuthenticated) {
      next({ name: 'login' });
    } else if (to.name === 'login' && isAuthenticated) {
      next({ name: 'home' });
    } else {
      next();
    }
  }
});

// Helper function to check authentication status
function checkAuthStatus(authStore: any): boolean {
  // Implement your authentication check here
  const token = authStore.accessToken;
  return token !== null && token !== '';
}

export default router;
