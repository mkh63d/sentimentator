import { mount } from '@vue/test-utils';
import { describe, it, expect } from 'vitest';
import LoadingOverlay from '../LoadingOverlay.vue';

describe('LoadingOverlay.vue', () => {
  it('renders when `visible` is true and shows the message', () => {
    const wrapper = mount(LoadingOverlay, {
      props: { visible: true, message: 'Loading…' },
    });

    // overlay should exist and contain the message
    expect(wrapper.find('.fixed').exists()).toBe(true);
    expect(wrapper.text()).toContain('Loading…');
  });

  it('does not render when `visible` is false', () => {
    const wrapper = mount(LoadingOverlay, {
      props: { visible: false },
    });

    expect(wrapper.find('.fixed').exists()).toBe(false);
  });

  it('emits update:visible when the internal computed setter is used', async () => {
    const wrapper = mount(LoadingOverlay, {
      props: { visible: true },
    });

    // Call the computed setter on the VM to simulate child updating the value
    type Exposed = { visibleLocal: boolean; messageLocal: string };
    const vm = wrapper.vm as unknown as Exposed;
    vm.visibleLocal = false;
    await wrapper.vm.$nextTick();

    const emitted = wrapper.emitted()['update:visible'];
    expect(emitted).toBeTruthy();
    expect(emitted![0]).toEqual([false]);
  });

  it('emits update:message when the message computed setter is used', async () => {
    const wrapper = mount(LoadingOverlay, {
      props: { visible: true, message: 'Old' },
    });

    const vm2 = wrapper.vm as unknown as { messageLocal: string };
    vm2.messageLocal = 'New message';
    await wrapper.vm.$nextTick();

    const emitted = wrapper.emitted()['update:message'];
    expect(emitted).toBeTruthy();
    expect(emitted![0]).toEqual(['New message']);
  });
});
