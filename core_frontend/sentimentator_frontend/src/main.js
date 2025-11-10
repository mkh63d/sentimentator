import { createApp } from 'vue';
import { createPinia } from 'pinia';

import App from './App.vue';
import router from './router';

if (import.meta.env.VITE_DEMO_MODE === 'true') {
  const { worker } = await import('./mocks/browser');
  await worker.start();
}

const app = createApp(App);

app.use(createPinia());
app.use(router);

app.mount('#app');
