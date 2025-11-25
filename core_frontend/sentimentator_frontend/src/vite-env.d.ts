/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_BACKEND_API: string;
  readonly VITE_DEMO_MODE: boolean;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
