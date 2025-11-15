import { setActivePinia, createPinia } from 'pinia';
import { describe, it, expect, beforeEach } from 'vitest';
import { useSearchStore } from '@/stores/searchStore';

describe('searchStore (component tests)', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it('defaults to hidden', () => {
    const store = useSearchStore();
    expect(store.visible).toBe(false);
  });

  it('show sets visible to true', () => {
    const store = useSearchStore();
    store.show();
    expect(store.visible).toBe(true);
  });

  it('hide sets visible to false', () => {
    const store = useSearchStore();
    store.show();
    expect(store.visible).toBe(true);
    store.hide();
    expect(store.visible).toBe(false);
  });
});
