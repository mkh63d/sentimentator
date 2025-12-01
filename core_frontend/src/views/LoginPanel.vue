<template>
  <div class="bg-linear-to-bl from-primary via-secondary to-fourth overflow-visible min-h-screen">
    <nav class="px-2 pt-5 pb-0">
      <div class="max-w-full w-full mx-auto flex justify-between items-center">
        <div
          class="px-5 text-2xl cursor-pointer font-bold text-textBase-light"
          @click="() => router.push({ name: 'home' })"
        >
          Sentimentator
        </div>
        <ThemeSwitcher class="mx-5" />
      </div>
    </nav>
    <div class="flex items-center justify-center">
      <div class="w-full max-w-md px-5 pt-7">
        <div v-if="mode === 'login'" class="bg-background p-9 rounded-xl shadow-lg">
          <h1 class="text-center mb-5 text-tertiary text-3xl font-bold">
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
                class="w-full p-3 border border-gray-300 text-textBase-dark rounded-md text-sm focus:outline-none focus:border-fourth transition"
              />
            </div>

            <div class="mb-2">
              <label for="password" class="block mb-1 text-tertiary font-medium">Password</label>
              <input
                id="password"
                v-model="form.password"
                type="password"
                placeholder="Enter your password"
                required
                class="w-full text-textBase-dark p-3 border border-gray-300 rounded-md text-sm focus:outline-none focus:border-fourth transition"
              />
            </div>

            <div v-if="isLogin" class="mb-2">
              <a
                @click="
                  () => {
                    mode = 'forgot';
                  }
                "
                class="text-sm text-accent font-medium hover:underline cursor-pointer"
                >Forgot Password?</a
              >
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

            <div class="flex items-center mb-4">
              <input
                id="rememberMe"
                type="checkbox"
                v-model="useCookies"
                class="h-4 w-4 text-fourth focus:ring-fourth border-gray-300 rounded"
              />
              <label for="rememberMe" class="ml-2 text-textBase-dark">Remember me</label>
            </div>

            <button
              type="submit"
              class="w-full cursor-pointer p-3 bg-fourth text-white rounded-md text-lg font-semibold hover:bg-fourth transition disabled:opacity-60"
              :disabled="!isFormValid"
            >
              {{ isLogin ? 'Login' : 'Register' }}
            </button>
          </form>

          <div class="text-center mt-6">
            <p class="text-textBase-dark border-gray-300 text-sm">
              {{ isLogin ? "Don't have an account?" : 'Already have an account?' }}
              <a
                @click="toggleMode"
                data-testid="mode-toggle"
                class="text-accent font-semibold cursor-pointer hover:underline ml-1"
              >
                {{ isLogin ? 'Register' : 'Login' }}
              </a>
            </p>
          </div>
        </div>

        <div
          v-else-if="mode === 'forgot'"
          class="bg-background p-10 rounded-xl shadow-lg flex flex-col items-center"
        >
          <h1 class="text-center mb-4 text-tertiary text-3xl font-bold">Forgot Password</h1>
          <p class="block mb-5 text-center text-tertiary font-medium">
            Please enter your email address to receive password reset instructions.
          </p>
          <div class="mb-5 w-full">
            <input
              id="email"
              v-model="form.email"
              type="email"
              placeholder="Enter your email"
              required
              class="w-full p-3 border border-gray-300 text-textBase rounded-md text-sm focus:outline-none focus:border-fourth transition"
            />
          </div>
          <div
            v-if="isForgotSent"
            class="text-primary mb-4 p-2 bg-primary/10 rounded text-center text-sm"
          >
            If an account with that email exists, a password reset link has been sent.
          </div>
          <button
            @click="forgotPassword()"
            :disabled="!form.email.trim() || isForgotSent"
            class="w-full cursor-pointer px-2 py-3 bg-fourth text-white rounded-md text-lg font-semibold hover:bg-fourth transition disabled:opacity-60 disabled:cursor-default"
          >
            Reset Password
          </button>
          <div class="text-center mt-6">
            <a
              @click="
                () => {
                  mode = 'login';
                  isForgotSent = false;
                }
              "
              class="text-primary font-semibold cursor-pointer hover:underline"
            >
              Back to Login
            </a>
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
import { forgotPassword as apiForgotPassword } from '@/api/account';
import { fireTimedAlert } from '@/lib/sweetAlerts';

const emit = defineEmits(['login', 'register']);
const router = useRouter();
const authStore = useAuthStore();

const useCookies = ref(false);
const useSessionCookies = ref(false);

const isLogin = ref(true);
const isForgotSent = ref(false);
const mode = ref<'login' | 'forgot' | '2fa'>('login');
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
    fireTimedAlert('Registration Error', error.value, 1000);
    error.value = '';
    return;
  }

  try {
    let pushToList = false;
    if (isLogin.value) {
      console.log('Logging in:', form.value.email);
      try {
        await authStore.loginUser(
          form.value.email,
          form.value.password,
          useCookies.value,
          useSessionCookies.value,
        );
        pushToList = true;
      } catch (e) {
        error.value = 'Login failed: ' + (e instanceof Error ? e.message : '');
        const errTitle = 'Login failed';
        const errMessage = e instanceof Error ? e.message : 'An unknown error occurred';
        fireTimedAlert(errTitle, errMessage, 1000);
        return;
      }
      emit('login', { email: form.value.email, password: form.value.password });
    } else {
      try {
        await authStore.registerUser(form.value.email, form.value.password);
        fireTimedAlert(
          'Success',
          'Registration successful! Confirm your email before logging in.',
          3000,
        );
      } catch (e) {
        error.value = 'Registration failed: ' + (e instanceof Error ? e.message : '');
        const errTitle = 'Registration failed';
        const errMessage = e instanceof Error ? e.message : 'An unknown error occurred';
        fireTimedAlert(errTitle, errMessage, 1000);
        return;
      }
      emit('register', { email: form.value.email, password: form.value.password });
    }
    if (pushToList) {
      await router.push({ name: 'list' });
    }
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

function forgotPassword() {
  try {
    apiForgotPassword(form.value.email);
  } catch (err) {
    const errTitle = 'Password Reset Failed';
    const errMessage = err instanceof Error ? err.message : 'An unknown error occurred';
    fireTimedAlert(errTitle, errMessage, 1000);
    return;
    // alert('An error occurred while trying to reset the password.');
    // console.error(err);
  }
  isForgotSent.value = true;
}
</script>
