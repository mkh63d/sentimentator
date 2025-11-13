<template>
  <div class="bg-linear-to-bl from-primary via-secondary to-fourth">
    <nav class="px-2 py-5">
      <div class="max-w-7xl mx-auto flex justify-between items-center">
        <div class="px-5 text-2xl font-bold text-background" @click="() => router.push('/home')">
          Sentimentator
        </div>
        <ThemeSwitcher class="mx-5" />
      </div>
    </nav>
    <div class="min-h-screen flex items-center justify-center">
      <div class="w-full max-w-md p-5">
        <div class="bg-background p-10 rounded-xl shadow-lg">
          <h1 class="text-center mb-8 text-tertiary text-3xl font-bold">
            {{ isLogin ? 'Login' : 'Register' }}
          </h1>

          <form @submit.prevent="handleSubmit">
            <div class="mb-5">
              <label for="email" class="block mb-1 text-tertiary font-medium">Email</label>
              <input
                id="email"
                v-model="form.email"
                type="email"
                placeholder="Enter your email"
                required
                class="w-full p-3 border border-gray-300 text-textBase rounded-md text-sm focus:outline-none focus:border-fourth transition"
              />
            </div>

            <div class="mb-5">
              <label for="password" class="block mb-1 text-tertiary font-medium">Password</label>
              <input
                id="password"
                v-model="form.password"
                type="password"
                placeholder="Enter your password"
                required
                class="w-full text-textBase p-3 border border-gray-300 rounded-md text-sm focus:outline-none focus:border-fourth transition"
              />
            </div>

            <div v-if="!isLogin" class="mb-5">
              <label for="confirmPassword" class="block mb-1 text-tertiary font-medium"
                >Confirm Password</label
              >
              <input
                id="confirmPassword"
                v-model="form.confirmPassword"
                type="password"
                placeholder="Confirm your password"
                required
                class="w-full p-3 border border-gray-300 text-textBase rounded-md text-sm focus:outline-none focus:border-fourth transition"
              />
            </div>

            <div
              data-testid="error-message"
              v-if="error"
              class="text-red-600 mb-4 p-2 bg-red-100 rounded text-sm"
            >
              {{ error }}
            </div>

            <button
              type="submit"
              class="w-full p-3 bg-fourth text-white rounded-md text-lg font-semibold hover:bg-fourth transition disabled:opacity-60"
              :disabled="!isFormValid"
            >
              {{ isLogin ? 'Login' : 'Register' }}
            </button>
          </form>

          <div class="text-center mt-6">
            <p class="text-gray-600 text-sm">
              {{ isLogin ? "Don't have an account?" : 'Already have an account?' }}
              <a
                @click="toggleMode"
                class="text-accent font-semibold cursor-pointer hover:underline ml-1"
              >
                {{ isLogin ? 'Register' : 'Login' }}
              </a>
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import ThemeSwitcher from '@/components/ThemeSwitcher.vue';
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { FormData } from '@/types/FormData';
import { useAuthStore } from '@/stores/authStore';

const emit = defineEmits(['login', 'register']);
const router = useRouter();
const authStore = useAuthStore();

const isLogin = ref(true);
const form = ref<FormData>({
  email: '',
  password: '',
  confirmPassword: '',
});
const error = ref('');

const isFormValid = computed(() => {
  const emailValid = form.value.email.trim() !== '';
  const passwordValid = form.value.password.trim() !== '';
  if (!isLogin.value) {
    return emailValid && passwordValid && form.value.confirmPassword.trim() !== '';
  }
  return emailValid && passwordValid;
});

function toggleMode() {
  isLogin.value = !isLogin.value;
  error.value = '';
  form.value.confirmPassword = '';
}

async function handleSubmit() {
  error.value = '';

  if (!isLogin.value && form.value.password !== form.value.confirmPassword) {
    error.value = 'Passwords do not match';
    return;
  }

  try {
    if (isLogin.value) {
      console.log('Logging in:', form.value.email);
      try {
        await authStore.loginUser(form.value.email, form.value.password);
      } catch (e) {
        error.value = 'Login failed: ' + (e instanceof Error ? e.message : '');
      }
      emit('login', { email: form.value.email, password: form.value.password });
    } else {
      try {
        await authStore.registerUser(form.value.email, form.value.password);
      } catch (e) {
        error.value = 'Registration failed: ' + (e instanceof Error ? e.message : '');
      }
      emit('register', { email: form.value.email, password: form.value.password });
    }
    router.push({ name: 'list' });
  } catch (err) {
    if (err instanceof Error) {
      error.value = err.message;
    } else if (typeof err === 'string') {
      error.value = err;
    } else {
      error.value = 'An error occurred';
    }
  }
}
</script>
