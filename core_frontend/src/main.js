import { createApp } from 'vue';
import { createPinia } from 'pinia';

import App from './App.vue';
import router from './router';
import i18n from './i18n';
import { getRuntimeConfig } from './lib/runtimeConfig';
import { initializeApiClient } from './lib/client';

// Load runtime configuration before starting the app
const config = await getRuntimeConfig();

initializeApiClient();

if (config.VITE_DEMO_MODE === 'true') {
  const { worker } = await import('./mocks/browser');
  await worker.start({ onUnhandledRequest: 'warn' });
}

const app = createApp(App);
const pinia = createPinia();

app.use(pinia);
app.use(router);
app.use(i18n);

app.mount('#app');
