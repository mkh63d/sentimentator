import { setupWorker } from 'msw/browser';
import { handlers } from './handlers';

// Create a service worker with all the mock handlers
export const worker = setupWorker(...handlers);
