import { mount } from '@vue/test-utils';
import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import ThemeSwitcher from '../ThemeSwitcher.vue';

describe('ThemeSwitcher.vue', () => {
  beforeEach(() => {
    // Ensure a clean DOM and localStorage for each test
    document.documentElement.removeAttribute('data-theme');
    localStorage.clear();
  });

  afterEach(() => {
    document.documentElement.removeAttribute('data-theme');
    localStorage.clear();
  });

  it('applies theme from localStorage on mount', async () => {
    localStorage.setItem('theme', 'dark');

    const { default: ThemeSwitcher } = await import('../ThemeSwitcher.vue');
    const wrapper = mount(ThemeSwitcher, { props: { changeOnTheme: true } });

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
