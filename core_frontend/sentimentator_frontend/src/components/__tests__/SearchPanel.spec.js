import { mount } from '@vue/test-utils';
import { describe, it, expect, beforeEach } from 'vitest';
import SearchPanel from '@/views/SearchPanel.vue';
import { createPinia, setActivePinia } from 'pinia';

describe('SearchPanel', () => {
  let wrapper;

  beforeEach(() => {
    const pinia = createPinia();
    setActivePinia(pinia);
    wrapper = mount(SearchPanel);
  });

  it('should_render_properly', () => {
    expect(wrapper.exists()).toBe(true);
  });

  it('should_update_link_input_field', () => {
    const linkInput = wrapper.find('input[aria-label="Add link"]');
    linkInput.setValue('https://example.com');
    expect(linkInput.element.value).toBe('https://example.com');
  });

  it('should_add_link_on_enter_key', async () => {
    const linkInput = wrapper.find('input[aria-label="Add link"]');
    await linkInput.setValue('https://example.com');
    await linkInput.trigger('keyup.enter');
    expect(wrapper.emitted()).toHaveProperty('add-link');
  });

  it('should_add_link_on_add_button_press', async () => {
    const linkInput = wrapper.find('input[aria-label="Add link"]');
    await linkInput.setValue('https://example.com');
    const addButton = wrapper.findAll('button').filter((btn) => btn.text() === 'Add')[0];
    await addButton.trigger('click');
    expect(wrapper.emitted()).toHaveProperty('add-link');
  });

  it('should_not_add_invalid_url', async () => {
    const linkInput = wrapper.find('input[aria-label="Add link"]');
    await linkInput.setValue('invalid-url');
    await linkInput.trigger('keyup.enter');
    expect(wrapper.emitted()['add-link']).toBeUndefined();
  });

  it('should_clear_all_links', async () => {
    const linkInput = wrapper.find('input[aria-label="Add link"]');
    await linkInput.setValue('https://example.com');
    await linkInput.trigger('keyup.enter');
    expect(wrapper.emitted()).toHaveProperty('add-link');

    const clearButton = wrapper.findAll('button').filter((btn) => btn.text() === 'Clear')[0];
    await clearButton.trigger('click');
    expect(wrapper.emitted()).toHaveProperty('clear-all');
  });

  it('should_display_error_for_invalid_url', async () => {
    const linkInput = wrapper.find('input[aria-label="Add link"]');
    await linkInput.setValue('invalid-url');
    expect(linkInput.classes()).toContain('border-red-500');
  });

  it('should_not_display_error_for_valid_url', async () => {
    const linkInput = wrapper.find('input[aria-label="Add link"]');
    await linkInput.setValue('https://example.com');
    expect(linkInput.classes()).not.toContain('border-red-500');
  });

  it('should_analyze_links_on_analyze_button_press', async () => {
    const linkInput = wrapper.find('input[aria-label="Add link"]');
    await linkInput.setValue('https://example.com');
    await linkInput.trigger('keyup.enter');
    expect(wrapper.emitted()).toHaveProperty('add-link');

    const analyzeButton = wrapper.findAll('button').filter((btn) => btn.text() === 'Analyze')[0];
    await analyzeButton.trigger('click');
    expect(wrapper.emitted()).toHaveProperty('analyze');
  });

  it('should_not_analyze_links_on_analyze_button_if_empty', async () => {
    const analyzeButton = wrapper.findAll('button').filter((btn) => btn.text() === 'Analyze')[0];
    await analyzeButton.trigger('click');
    expect(wrapper.emitted()['analyze']).toBeUndefined();
  });
});
