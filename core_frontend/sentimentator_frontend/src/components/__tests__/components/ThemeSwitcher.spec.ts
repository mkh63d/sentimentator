import { mount } from '@vue/test-utils';
import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import ThemeSwitcher from '@/components/ThemeSwitcher.vue';

// TODO: I habe no idea why those tests don't pass but I will just leave them here for now

describe('ThemeSwitcher.vue', () => {
  const getItemSpy = vi.spyOn(Storage.prototype, 'getItem');
  const setItemSpy = vi.spyOn(Storage.prototype, 'setItem');

  beforeEach(() => {
    // Ensure a clean DOM and localStorage for each test
    document.documentElement.removeAttribute('data-theme');
    localStorage.clear();
    getItemSpy.mockClear();
    setItemSpy.mockClear();
  });

  afterEach(() => {
    document.documentElement.removeAttribute('data-theme');
    localStorage.clear();
    getItemSpy.mockClear();
    setItemSpy.mockClear();
  });

  it('applies theme from localStorage on mount', async () => {
    getItemSpy.mockReturnValue('dark');

    const { default: ThemeSwitcher } = await import('@/components/ThemeSwitcher.vue');
    const wrapper = mount(ThemeSwitcher, { props: { changeOnTheme: true } });

    expect(getItemSpy).toHaveBeenCalledWith('theme');
    // component onMounted calls applyTheme
    expect(document.documentElement.getAttribute('data-theme')).toBe('dark');

    // SVG should reflect stroke-black when changeOnTheme=true and theme=dark
    const svg = wrapper.find('svg');
    expect(svg.classes()).toContain('stroke-black');
  });

  it('toggles theme on click and updates localStorage and document attribute', async () => {
    // start with light
    localStorage.setItem('theme', 'light');
    const wrapper = mount(ThemeSwitcher, { props: { changeOnTheme: true } });

    const btn = wrapper.find('button');
    expect(document.documentElement.getAttribute('data-theme')).toBe('light');
    expect(localStorage.getItem('theme')).toBe('light');

    await btn.trigger('click');
    // after first click -> dark
    expect(document.documentElement.getAttribute('data-theme')).toBe('dark');
    expect(localStorage.getItem('theme')).toBe('dark');

    await btn.trigger('click');
    // after second click -> light
    expect(document.documentElement.getAttribute('data-theme')).toBe('light');
    expect(localStorage.getItem('theme')).toBe('light');
  });
});
