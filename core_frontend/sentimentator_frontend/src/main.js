import { createApp } from 'vue';
import { createPinia } from 'pinia';

import App from './App.vue';
import router from './router';
import { getRuntimeConfig } from './lib/runtimeConfig';

// Load runtime configuration before starting the app
const config = await getRuntimeConfig();

if (config.VITE_DEMO_MODE === 'true') {
  const { worker } = await import('./mocks/browser');
  await worker.start({ onUnhandledRequest: 'warn' });
}

const app = createApp(App);
const pinia = createPinia();

app.use(pinia);
app.use(router);

app.mount('#app');
