import { setActivePinia, createPinia } from 'pinia';
import { useLoadingStore } from '../../../stores/useLoadingStore';
import { describe, it, expect, beforeEach } from 'vitest';

describe('Loading Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it('should have default values', () => {
    const loadingStore = useLoadingStore();
    expect(loadingStore.visible).toBe(false);
    expect(loadingStore.message).toBe('');
  });

  it('should show loading with default message', () => {
    const loadingStore = useLoadingStore();
    loadingStore.show();
    expect(loadingStore.visible).toBe(true);
    expect(loadingStore.message).toBe('Loading...');
  });

  it('should show loading with custom message', () => {
    const loadingStore = useLoadingStore();
    const customMessage = 'Please wait...';
    loadingStore.show(customMessage);
    expect(loadingStore.visible).toBe(true);
    expect(loadingStore.message).toBe(customMessage);
  });

  it('should hide loading', () => {
    const loadingStore = useLoadingStore();
    loadingStore.show('Loading data...');
    loadingStore.hide();
    expect(loadingStore.visible).toBe(false);
    expect(loadingStore.message).toBe('');
  });
});
