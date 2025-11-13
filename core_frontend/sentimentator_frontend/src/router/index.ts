import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  { path: '/login', name: 'login', component: () => import('@/views/LoginPanel.vue') },
  { path: '/', redirect: { name: 'home' } },
  { path: '/list', name: 'list', component: () => import('@/views/ListOfAnalysisView.vue') },
  { path: '/home', name: 'home', component: () => import('@/views/HomePage.vue') },
  { path: '/analysis', name: 'analysis', component: () => import('@/views/AnalysisView.vue') },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
