<template>
  <div
    class="min-h-screen flex items-center justify-center bg-linear-to-br from-indigo-500 to-purple-700"
  >
    <div class="w-full max-w-md p-5">
      <div class="bg-white p-10 rounded-xl shadow-lg">
        <h1 class="text-center mb-8 text-gray-800 text-3xl font-bold">
          {{ isLogin ? 'Login' : 'Register' }}
        </h1>

        <form @submit.prevent="handleSubmit">
          <div class="mb-5">
            <label for="email" class="block mb-1 text-gray-700 font-medium">Email</label>
            <input
              id="email"
              v-model="form.email"
              type="email"
              placeholder="Enter your email"
              required
              class="w-full p-3 border border-gray-300 rounded-md text-sm focus:outline-none focus:border-indigo-500 transition"
            />
          </div>

          <div class="mb-5">
            <label for="password" class="block mb-1 text-gray-700 font-medium">Password</label>
            <input
              id="password"
              v-model="form.password"
              type="password"
              placeholder="Enter your password"
              required
              class="w-full p-3 border border-gray-300 rounded-md text-sm focus:outline-none focus:border-indigo-500 transition"
            />
          </div>

          <div v-if="!isLogin" class="mb-5">
            <label for="confirmPassword" class="block mb-1 text-gray-700 font-medium"
              >Confirm Password</label
            >
            <input
              id="confirmPassword"
              v-model="form.confirmPassword"
              type="password"
              placeholder="Confirm your password"
              required
              class="w-full p-3 border border-gray-300 rounded-md text-sm focus:outline-none focus:border-indigo-500 transition"
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
            class="w-full p-3 bg-indigo-500 text-white rounded-md text-lg font-semibold hover:bg-indigo-600 transition disabled:opacity-60"
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
              class="text-indigo-500 font-semibold cursor-pointer hover:underline ml-1"
            >
              {{ isLogin ? 'Register' : 'Login' }}
            </a>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { FormData } from '@/types/FormData';

const emit = defineEmits(['login', 'register']);
const router = useRouter();

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
      // Implement login logic
      console.log('Logging in:', form.value.email);
      emit('login', { email: form.value.email, password: form.value.password });
    } else {
      // Implement registration logic
      console.log('Registering:', form.value.email);
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
