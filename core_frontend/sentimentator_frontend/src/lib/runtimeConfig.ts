interface RuntimeConfig {
  VITE_BACKEND_API: string;
  VITE_DEMO_MODE: string;
}

let configCache: RuntimeConfig | null = null;
let configPromise: Promise<RuntimeConfig> | null = null;

/**
 * Loads runtime configuration from /runtime-config.json
 * Cached after first successful load
 */
export async function getRuntimeConfig(): Promise<RuntimeConfig> {
  if (configCache) {
    return configCache;
  }

  if (configPromise) {
    return configPromise;
  }

  configPromise = (async () => {
    try {
      const response = await fetch('/runtime-config.json');
      if (!response.ok) {
        throw new Error(`Failed to load runtime config: ${response.status}`);
      }
      const config = await response.json();
      configCache = config;
      return config;
    } catch (error) {
      console.error('Failed to load runtime configuration, using defaults:', error);
      // Fallback to build-time env or defaults
      configCache = {
        VITE_BACKEND_API: import.meta.env.VITE_BACKEND_API || 'http://localhost:8081/api',
        VITE_DEMO_MODE: import.meta.env.VITE_DEMO_MODE || 'false',
      };
      return configCache;
    }
  })();

  return configPromise;
}

/**
 * Get config synchronously (only after getRuntimeConfig has been called)
 * Use this after config has been loaded during app bootstrap
 */
export function getConfig(): RuntimeConfig {
  if (!configCache) {
    throw new Error('Runtime config not loaded. Call getRuntimeConfig() first.');
  }
  return configCache;
}
