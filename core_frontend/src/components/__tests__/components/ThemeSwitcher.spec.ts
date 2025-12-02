import { mount } from '@vue/test-utils';
import { describe, it, expect, beforeEach, afterEach } from 'vitest';

describe('ThemeSwitcher.vue', () => {
  beforeEach(() => {
    // Ensure a clean DOM for each test
    document.documentElement.removeAttribute('data-theme');
  });

  afterEach(() => {
    document.documentElement.removeAttribute('data-theme');
  });

  // it('applies theme from localStorage on mount', async () => {
  //   // use the global mocked localStorage provided by vitest.setup.js
  //   const getItemSpy = vi.spyOn(globalThis.localStorage, 'getItem');
  //   const setItemSpy = vi.spyOn(globalThis.localStorage, 'setItem');

  //   // simulate theme present in storage before importing the module
  //   globalThis.localStorage.setItem('theme', 'dark');

  //   const { default: ThemeSwitcher } = await import('@/components/ThemeSwitcher.vue');
  //   const wrapper = mount(ThemeSwitcher, { props: { changeOnTheme: true } });
  //   // component onMounted calls applyTheme
  //   expect(document.documentElement.getAttribute('data-theme')).toBe('dark');

  //   // SVG should reflect stroke-black when changeOnTheme=true and theme=dark
  //   const svg = wrapper.find('svg');
  //   expect(svg.classes()).toContain('stroke-black');

  //   getItemSpy.mockRestore();
  //   setItemSpy.mockRestore();
  // });

  it('toggles theme on click and updates localStorage and document attribute', async () => {
    // start with light using the global mocked localStorage
    globalThis.localStorage.setItem('theme', 'light');

    const { default: ThemeSwitcher } = await import('@/components/ThemeSwitcher.vue');
    const wrapper = mount(ThemeSwitcher, { props: { changeOnTheme: true } });

    const btn = wrapper.find('button');
    expect(document.documentElement.getAttribute('data-theme')).toBe('light');
    expect(globalThis.localStorage.getItem('theme')).toBe('light');

    await btn.trigger('click');
    // after first click -> dark
    expect(document.documentElement.getAttribute('data-theme')).toBe('dark');
    expect(globalThis.localStorage.getItem('theme')).toBe('dark');

    await btn.trigger('click');
    // after second click -> light
    expect(document.documentElement.getAttribute('data-theme')).toBe('light');
    expect(globalThis.localStorage.getItem('theme')).toBe('light');

    // spies are not necessary to restore here because vitest.setup resets storage
  });
});
