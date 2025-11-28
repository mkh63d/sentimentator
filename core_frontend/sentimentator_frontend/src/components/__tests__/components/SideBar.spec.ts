import { setActivePinia, createPinia } from 'pinia';
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { mount } from '@vue/test-utils';

// mock router at module scope so vitest hoists the mock and the component's
// module import sees the mocked router
let pushMock = vi.fn();
vi.mock('vue-router', () => ({ useRouter: () => ({ push: pushMock }) }));

describe('SideBar.vue', () => {
  beforeEach(() => {
    vi.resetModules();
    vi.resetAllMocks();
    pushMock = vi.fn();
    setActivePinia(createPinia());
  });

  it('is collapsed when sidebarStore.visible is false', async () => {
    const { default: SideBar } = await import('@/components/SideBar.vue');
    const wrapper = mount(SideBar, { global: { stubs: { X: true } } });

    const aside = wrapper.find('aside');
    expect(aside.exists()).toBe(true);
    // when hidden, the bound class includes 'w-0'
    expect(aside.attributes('class')).toContain('w-0');
  });

  it('shows content when sidebarStore.visible is true and close button hides it', async () => {
    const { useSidebarStore } = await import('@/stores/sidebarStore');
    const store = useSidebarStore();
    store.show();

    const { default: SideBar } = await import('@/components/SideBar.vue');
    const wrapper = mount(SideBar, { global: { stubs: { X: true } } });

    // Logout button should be present when visible - find by text
    const logoutByText = wrapper.findAll('button').find((b) => b.text().includes('Logout'));
    expect(logoutByText).toBeTruthy();

    // Close button is the first absolute button; trigger it
    const closeBtn = wrapper.find('button');
    await closeBtn.trigger('click');
    expect(store.visible).toBe(false);
  });

  it('calls authStore.logoutUser and navigates to /login on Logout click', async () => {
    const { useSidebarStore } = await import('@/stores/sidebarStore');
    const sidebar = useSidebarStore();
    sidebar.show();

    const { useAuthStore } = await import('@/stores/authStore');
    const auth = useAuthStore();
    // replace logoutUser with a spy
    const logoutSpy = vi.fn().mockResolvedValue(undefined);
    auth.logoutUser = logoutSpy;

    const { default: SideBar } = await import('@/components/SideBar.vue');
    const wrapper = mount(SideBar, { global: { stubs: { X: true } } });

    const logoutBtn = wrapper.findAll('button').find((b) => b.text().includes('Logout'));
    expect(logoutBtn).toBeTruthy();
    await logoutBtn!.trigger('click');

    expect(logoutSpy).toHaveBeenCalled();
    expect(pushMock).toHaveBeenCalledWith('/login');
    expect(sidebar.visible).toBe(false);
  });
});
