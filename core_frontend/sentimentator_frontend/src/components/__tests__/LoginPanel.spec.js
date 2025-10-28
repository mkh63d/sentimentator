import { describe, it, expect, beforeEach } from 'vitest';
import { mount } from '@vue/test-utils';
import LoginPanel from '@/views/LoginPanel.vue';

describe('LoginPanel', () => {
  let wrapper;

  beforeEach(() => {
    wrapper = mount(LoginPanel);
  });

  it('should_render_properly', () => {
    expect(wrapper.exists()).toBe(true);
  });

  it('should_display_username_and_password_input_fields', () => {
    expect(wrapper.find('input[type="email"]').exists()).toBe(true);
    expect(wrapper.find('input[type="password"]').exists()).toBe(true);
  });

  it('should_display_login_button', () => {
    expect(wrapper.find('button').exists()).toBe(true);
  });

  it('should_update_email_when_input_changes', async () => {
    const email = wrapper.find('input[type="email"]');
    await email.setValue('testuser@example.com');
    expect(email.element.value).toBe('testuser@example.com');
  });

  it('should_update_password_when_input_changes', async () => {
    const passwordInput = wrapper.find('input[type="password"]');
    await passwordInput.setValue('testpassword');
    expect(passwordInput.element.value).toBe('testpassword');
  });

  it('should_emit_login_event_when_form_is_submitted', async () => {
    const emailInput = wrapper.find('input[type="email"]');
    const passwordInput = wrapper.find('input[type="password"]');

    await emailInput.setValue('testuser@example.com');
    await passwordInput.setValue('testpassword');
    await wrapper.find('form').trigger('submit.prevent');

    expect(wrapper.emitted()).toHaveProperty('login');
  });

  it('should_disable_submit_button_when_fields_are_empty', () => {
    const button = wrapper.find('button[type="submit"]');
    expect(button.attributes('disabled')).toBeDefined();
  });

  it('should_enable_submit_button_when_both_fields_are_filled', async () => {
    await wrapper.find('input[type="email"]').setValue('testuser@example.com');
    await wrapper.find('input[type="password"]').setValue('testpassword');

    const button = wrapper.find('button[type="submit"]');
    expect(button.attributes('disabled')).toBeUndefined();
  });

  it('should_change_header_text_based_on_mode', async () => {
    const header = wrapper.find('h1').text();
    expect(header).toBe('Login');

    // Switch to register mode
    await wrapper.find('a').trigger('click');

    expect(wrapper.find('h1').text()).toBe('Register');
  });

  it('should_change_button_text_based_on_mode', async () => {
    const button = wrapper.find('button[type="submit"]');
    expect(button.text()).toBe('Login');

    // Switch to register mode
    await wrapper.find('a').trigger('click');

    expect(wrapper.find('button[type="submit"]').text()).toBe('Register');
  });

  it('should_toggle_mode_when_toggle_link_is_clicked', async () => {
    const toggleLink = wrapper
      .findAll('a')
      .filter((node) => node.text().includes('Register') || node.text().includes('Login'))[0];
    expect(wrapper.vm.isLogin).toBe(true);

    await toggleLink.trigger('click');

    expect(wrapper.vm.isLogin).toBe(false);
  });

  it('should_display_error_message_when_different_passwords_are_entered_during_registration', async () => {
    // Switch to register mode
    await wrapper.find('a').trigger('click');

    await wrapper.find('input[type="email"]').setValue('testuser@example.com');
    await wrapper.find('input[type="password"]').setValue('password1');
    await wrapper.find('input[id="confirmPassword"]').setValue('password2');
    await wrapper.find('form').trigger('submit.prevent');

    expect(wrapper.find('[data-testid="error-message"]').exists()).toBe(true);
  });

  it('should_emit_register_event_when_registration_form_is_submitted', async () => {
    // Switch to register mode
    await wrapper.find('a').trigger('click');

    await wrapper.find('input[type="email"]').setValue('testuser@example.com');
    await wrapper.find('input[type="password"]').setValue('password1');
    await wrapper.find('input[id="confirmPassword"]').setValue('password1');
    await wrapper.find('form').trigger('submit.prevent');

    expect(wrapper.emitted()).toHaveProperty('register');
  });
});
