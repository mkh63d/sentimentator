<template>
    <div class="login-panel">
        <div class="container">
            <div class="form-wrapper">
                <h1>{{ isLogin ? 'Login' : 'Register' }}</h1>
                
                <form @submit.prevent="handleSubmit">
                    <div class="form-group">
                        <label for="email">Email</label>
                        <input
                            id="email"
                            v-model="form.email"
                            type="email"
                            placeholder="Enter your email"
                            required    
                        />
                    </div>

                    <div class="form-group">
                        <label for="password">Password</label>
                        <input
                            id="password"
                            v-model="form.password"
                            type="password"
                            placeholder="Enter your password"
                            required
                        />
                    </div>

                    <div v-if="!isLogin" class="form-group">
                        <label for="confirmPassword">Confirm Password</label>
                        <input
                            id="confirmPassword"
                            v-model="form.confirmPassword"
                            type="password"
                            placeholder="Confirm your password"
                            required
                        />
                    </div>

                    <div v-if="error" class="error-message">{{ error }}</div>

                    <button type="submit" class="btn-primary" :disabled="!isFormValid">
                        {{ isLogin ? 'Login' : 'Register' }}
                    </button>
                </form>

                <div class="toggle-mode">
                    <p>
                        {{ isLogin ? "Don't have an account?" : "Already have an account?" }}
                        <a @click="toggleMode">{{ isLogin ? 'Register' : 'Login' }}</a>
                    </p>
                </div>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import { FormData } from '@/types/FormData'

export default {
    name: 'LoginPanel',
    data() {
        return {
            isLogin: true,
            form: {
                email: '',
                password: '',
                confirmPassword: ''
            } as FormData,
            error: '' as string
        }
    },
    computed: {
        isFormValid(): boolean {
            const emailValid = this.form.email.trim() !== ''
            const passwordValid = this.form.password.trim() !== ''
            if (!this.isLogin) {
                return emailValid && passwordValid && this.form.confirmPassword.trim() !== ''
            }
            return emailValid && passwordValid
        }
    },
    methods: {
        toggleMode(): void {
            this.isLogin = !this.isLogin
            this.error = ''
            this.form.confirmPassword = ''
        },
        async handleSubmit(): Promise<void> {
            this.error = ''

            if (!this.isLogin && this.form.password !== this.form.confirmPassword) {
                this.error = 'Passwords do not match'
                return
            }

            try {
                if (this.isLogin) {
                    // Implement login logic
                    console.log('Logging in:', this.form.email)
                    // await this.$store.dispatch('login', { email: this.form.email, password: this.form.password })
                    this.$emit('login', { email: this.form.email, password: this.form.password })
                } else {
                    // Implement registration logic
                    console.log('Registering:', this.form.email)
                    // await this.$store.dispatch('register', { email: this.form.email, password: this.form.password })
                    this.$emit('register', { email: this.form.email, password: this.form.password })
                }
                this.$router.push('/list')
            } catch (err) {
                this.error = (err as Error).message || 'An error occurred'
            }
        }
           
    }
    
}

</script>

<style scoped>
.login-panel {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.container {
    width: 100%;
    max-width: 400px;
    padding: 20px;
}

.form-wrapper {
    background: white;
    padding: 40px;
    border-radius: 10px;
    box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
}

h1 {
    text-align: center;
    margin-bottom: 30px;
    color: #333;
}

.form-group {
    margin-bottom: 20px;
}

label {
    display: block;
    margin-bottom: 5px;
    color: #555;
    font-weight: 500;
}

input {
    width: 100%;
    padding: 12px;
    border: 1px solid #ddd;
    border-radius: 5px;
    font-size: 14px;
    transition: border-color 0.3s;
    box-sizing: border-box;
}

input:focus {
    outline: none;
    border-color: #667eea;
}

.btn-primary {
    width: 100%;
    padding: 12px;
    background: #667eea;
    color: white;
    border: none;
    border-radius: 5px;
    font-size: 16px;
    font-weight: 600;
    cursor: pointer;
    transition: background 0.3s;
}

.btn-primary:hover {
    background: #5568d3;
}

.error-message {
    color: #e74c3c;
    margin-bottom: 15px;
    padding: 10px;
    background: #fadbd8;
    border-radius: 5px;
    font-size: 14px;
}

.toggle-mode {
    text-align: center;
    margin-top: 20px;
}

.toggle-mode p {
    color: #666;
    font-size: 14px;
}

.toggle-mode a {
    color: #667eea;
    cursor: pointer;
    font-weight: 600;
    text-decoration: none;
}

.toggle-mode a:hover {
    text-decoration: underline;
}
</style>