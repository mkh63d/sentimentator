import { setActivePinia, createPinia } from 'pinia';
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { mount, shallowMount } from '@vue/test-utils';

// router push mock needs to be visible to the mocked module factory (vitest hoists mocks),
// so declare it at module scope and reset it in beforeEach.
let pushMock = vi.fn();
vi.mock('vue-router', () => ({ useRouter: () => ({ push: pushMock }) }));

describe('NavbarComponent', () => {
  beforeEach(() => {
    vi.resetModules();
    vi.resetAllMocks();
    // reset the push mock used by the mocked router
    pushMock = vi.fn();
    setActivePinia(createPinia());
  });

  it('toggles sidebar visible when hamburger is clicked', async () => {
    const { useSidebarStore } = await import('@/stores/sidebarStore');
    const store = useSidebarStore();
    expect(store.visible).toBe(false);

    const { default: Navbar } = await import('@/components/NavbarComponent.vue');
    const wrapper = mount(Navbar, { global: { stubs: { ThemeSwitcher: true } } });

    const btn = wrapper.find('.cursor-pointer');
    expect(btn.exists()).toBe(true);
    await btn.trigger('click');

    expect(store.visible).toBe(true);
  });

  it('navigates to /home when title is clicked', async () => {
    const { default: Navbar } = await import('@/components/NavbarComponent.vue');
    const wrapper = shallowMount(Navbar, { global: { stubs: { ThemeSwitcher: true } } });

    const title = wrapper.find('.text-2xl');
    expect(title.exists()).toBe(true);
    await title.trigger('click');

    expect(pushMock).toHaveBeenCalledWith('/home');
  });

  it('renders ThemeSwitcher with changeOnTheme prop', async () => {
    const { default: Navbar } = await import('@/components/NavbarComponent.vue');
    const wrapper = shallowMount(Navbar, { global: { stubs: ['ThemeSwitcher'] } });

    const themeStub = wrapper.findComponent({ name: 'ThemeSwitcher' });
    expect(themeStub.exists()).toBe(true);
    const hasProp =
      themeStub.attributes()['changeontheme'] !== undefined ||
      themeStub.attributes()['change-on-theme'] !== undefined;
    expect(hasProp).toBe(true);
  });
});
