import { setActivePinia, createPinia } from 'pinia';
import { describe, it, expect, beforeEach } from 'vitest';
import { useSidebarStore } from '@/stores/sidebarStore';

describe('sidebarStore (component tests)', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it('defaults to hidden', () => {
    const store = useSidebarStore();
    expect(store.visible).toBe(false);
  });

  it('show sets visible to true', () => {
    const store = useSidebarStore();
    store.show();
    expect(store.visible).toBe(true);
  });

  it('hide sets visible to false', () => {
    const store = useSidebarStore();
    store.show();
    expect(store.visible).toBe(true);
    store.hide();
    expect(store.visible).toBe(false);
  });
});
