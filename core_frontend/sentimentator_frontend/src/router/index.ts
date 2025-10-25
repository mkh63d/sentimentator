import {
  createRouter,
  createWebHistory,
  type RouteRecordRaw,
} from 'vue-router';
import LoginPanel from '../views/LoginPanel.vue';
import SearchPanel from '../views/SearchPanel.vue';
//import Analysis from '../views/Analysis.vue';

const routes: RouteRecordRaw[] = [
  //{ path: '/', name: 'Home', component: Home },
  //{ path: '/about', name: 'About', component: About }
  { path: '/login', name: 'Login', component: LoginPanel },
  { path: '/', redirect: { name: 'Search' } },
  { path: '/search', name: 'Search', component: SearchPanel },
 // { path: '/analysis', name: 'Analysis', component: Analysis },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
