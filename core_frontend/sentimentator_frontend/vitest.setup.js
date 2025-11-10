import { vi } from 'vitest';

// Prevent Pinia from initializing devtools
vi.stubGlobal('__VUE_DEVTOOLS_GLOBAL_HOOK__', {
  on: vi.fn(),
  emit: vi.fn(),
});

// Provide a working localStorage for libraries that expect it
globalThis.localStorage = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
};
