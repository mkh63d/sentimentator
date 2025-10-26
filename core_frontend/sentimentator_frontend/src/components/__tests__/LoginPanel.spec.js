import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import LoginPanel from '@/views/LoginPanel.vue'

describe('LoginPanel', () => {
    let wrapper

    beforeEach(() => {
        wrapper = mount(LoginPanel)
    })

    it('should_render_properly', () => {
        expect(wrapper.exists()).toBe(true)
    })

    it('should_display_username_and_password_input_fields', () => {
        expect(wrapper.find('input[type="email"]').exists()).toBe(true)
        expect(wrapper.find('input[type="password"]').exists()).toBe(true)
    })

    it('should_display_login_button', () => {
        expect(wrapper.find('button').exists()).toBe(true)
    })

    it('should_update_email_when_input_changes', async () => {
        const email = wrapper.find('input[type="email"]')
        await email.setValue('testuser@example.com')
        expect(email.element.value).toBe('testuser@example.com')
    })

    it('should_update_password_when_input_changes', async () => {
        const passwordInput = wrapper.find('input[type="password"]')
        await passwordInput.setValue('testpassword')
        expect(passwordInput.element.value).toBe('testpassword')
    })

    it('should_emit_login_event_when_form_is_submitted', async () => {
        const emailInput = wrapper.find('input[type="email"]')
        const passwordInput = wrapper.find('input[type="password"]')
        
        await emailInput.setValue('testuser@example.com')
        await passwordInput.setValue('testpassword')
        await wrapper.find('form').trigger('submit.prevent')

        expect(wrapper.emitted()).toHaveProperty('login')
    })

    it('should_disable_submit_button_when_fields_are_empty', () => {
        const button = wrapper.find('button[type="submit"]')
        expect(button.attributes('disabled')).toBeDefined()
    })

    it('should_enable_submit_button_when_both_fields_are_filled', async () => {
        await wrapper.find('input[type="email"]').setValue('testuser')
        await wrapper.find('input[type="password"]').setValue('testpassword')
        
        const button = wrapper.find('button[type="submit"]')
        expect(button.attributes('disabled')).toBeUndefined()
    })
})