import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router';
//import Analysis from '../views/Analysis.vue';

const routes: RouteRecordRaw[] = [
  //{ path: '/', name: 'Home', component: Home },
  //{ path: '/about', name: 'About', component: About }
  { path: '/login', name: 'login', component: () => import('@/views/LoginPanel.vue') },
  { path: '/', redirect: { name: 'search' } },
  { path: '/search', name: 'search', component: () => import('@/views/SearchPanel.vue') },
  { path: '/list', name: 'list', component: () => import('@/views/ListOfAnalysisView.vue') },
  { path: '/home', name: 'home', component: () => import('@/views/HomePage.vue') },
  // { path: '/analysis', name: 'analysis', component: analysis },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
