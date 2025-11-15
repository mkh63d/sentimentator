import { mount } from '@vue/test-utils';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { ref, type Component } from 'vue';

// Mock the search store used by the component
const visibleRef = ref(true);
const hideMock = vi.fn();

vi.mock('@/stores/searchStore', () => ({
  useSearchStore: () => ({ visible: visibleRef, hide: hideMock, show: vi.fn() }),
}));

// Mock the createAnalysis API so submit flow doesn't perform network calls
const createAnalysisMock = vi.fn(() => Promise.resolve({}));
vi.mock('@/api/analyses', () => ({ createAnalysis: createAnalysisMock }));

// Stub icons and LoadingOverlay to keep the test focused
const globalStubs = {
  X: true,
  LoadingOverlay: true,
};

describe('SearchModal.vue', () => {
  let SearchModal: Component;

  beforeEach(async () => {
    // reset state before each test
    visibleRef.value = true;
    hideMock.mockClear();
    createAnalysisMock.mockClear();
    // import component lazily to ensure mocks are applied
    // use dynamic ESM import instead of require so TypeScript/ESM environments work
    SearchModal = (await import('../SearchModal.vue')).default;
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('renders when the store visible is true', () => {
    const wrapper = mount(SearchModal, { global: { stubs: globalStubs } });
    expect(wrapper.find('h2').text()).toContain('Add Links for Analysis');
  });

  it('adds a valid link and clears the input', async () => {
    const wrapper = mount(SearchModal, { global: { stubs: globalStubs } });
    const input = wrapper.find('input[aria-label="Add link"]');
    await input.setValue('https://www.youtube.com/watch?v=abc');

    const addButton = wrapper.findAll('button').find((b) => b.text() === 'Add');
    // fallback if find approach didn't work, find by index via buttons
    if (!addButton) {
      const buttons = wrapper.findAll('button');
      await buttons[1].trigger('click');
    } else {
      await addButton.trigger('click');
    }

    // link should appear in the list
    expect(wrapper.findAll('ul li').length).toBeGreaterThan(0);
    expect(wrapper.find('ul li a').text()).toContain('youtube.com');
    // input cleared
    expect((wrapper.find('input[aria-label="Add link"]').element as HTMLInputElement).value).toBe(
      '',
    );
  });

  it('removes a link when remove button clicked', async () => {
    const wrapper = mount(SearchModal, { global: { stubs: globalStubs } });
    const input = wrapper.find('input[aria-label="Add link"]');
    await input.setValue('https://www.youtube.com/watch?v=abc');
    // click Add
    await wrapper.findAll('button')[1].trigger('click');

    expect(wrapper.findAll('ul li').length).toBe(1);
    // click remove button (inside first li)
    await wrapper.find('ul li button[aria-label="Remove link"]').trigger('click');
    expect(wrapper.findAll('ul li').length).toBe(0);
  });

  it('Analyse button is enabled only when there are links and triggers submit flow', async () => {
    const wrapper = mount(SearchModal, { global: { stubs: globalStubs } });

    // Initially disabled
    const analyseButton = wrapper.findAll('button').find((b) => b.text() === 'Analyse');
    expect(analyseButton).toBeDefined();
    expect(analyseButton!.attributes('disabled')).toBeDefined();

    // Add a link
    await wrapper
      .find('input[aria-label="Add link"]')
      .setValue('https://www.youtube.com/watch?v=abc');
    await wrapper.findAll('button')[1].trigger('click');

    // Now enabled
    const analyseBtnAfter = wrapper.findAll('button').find((b) => b.text() === 'Analyse');
    expect(analyseBtnAfter).toBeDefined();
    expect(analyseBtnAfter!.attributes('disabled')).toBeUndefined();

    // Mock timers to advance the submit delay
    vi.useFakeTimers();
    await analyseBtnAfter!.trigger('click');

    // advance the 1s simulated delay inside submit
    vi.advanceTimersByTime(1000);
    // allow pending promises to resolve
    await Promise.resolve();
    await Promise.resolve();

    // createAnalysis should have been called
    expect(createAnalysisMock).toHaveBeenCalled();
    // store.hide should be called
    expect(hideMock).toHaveBeenCalled();
    // component should emit analyze event
    expect(wrapper.emitted()['analyze']).toBeTruthy();
  });
});
