<template>
  <div class="min-h-screen bg-secondary">
    <NavbarComponent />
    <div class="flex min-h-screen">
      <SideBar v-model="sidebarStore.visible">
        <ul class="space-y-2">
          <li>
            <button
              @click="router.push({ name: 'home' })"
              class="font-medium text-left transition cursor-pointer text-background hover:text-indigo-100"
            >
              Home
            </button>
          </li>
          <li>
            <button
              @click="router.push({ name: 'list' })"
              class="font-medium text-left transition cursor-pointer text-background hover:text-indigo-100"
            >
              Analysis Library
            </button>
          </li>
          <li class="self-bottom">
            <button
              @click="router.push({ name: 'settings' })"
              class="font-medium text-left transition cursor-pointer text-background hover:text-indigo-100"
            >
              Settings
            </button>
          </li>
        </ul>
      </SideBar>
      <div class="px-8 py-8 max-w-[1400px] mx-auto flex-1">
        <div class="flex items-center justify-center flex-1 mb-8">
          <h1 class="m-0 text-3xl font-extrabold text-tertiary">Settings</h1>
        </div>
        <div class="flex flex-col mb-8">
          <div class="p-6 mb-6 rounded-lg shadow-md bg-background">
            <h2 class="mb-4 text-xl font-semibold text-tertiary">Account Settings</h2>
            <div class="space-y-4">
              <div>
                <label class="block mb-2 text-sm font-medium text-tertiary">
                  Change Email Address
                </label>
                <input
                  type="email"
                  v-model="email"
                  class="w-full text-tertiary px-3 py-2 my-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#10a20b] focus:border-transparent"
                  placeholder="Enter new email address"
                />
                <button
                  :disabled="!email"
                  class="px-4 py-2 mt-2 text-white transition rounded-md cursor-pointer disabled:cursor-default bg-primary hover:bg-primary/80 disabled:opacity-50"
                >
                  Update Email
                </button>
              </div>
              <div>
                <label class="block mb-2 text-sm font-medium text-tertiary">
                  Change password
                </label>
                <input
                  type="password"
                  class="w-full text-tertiary px-3 py-2 my-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#10a20b] focus:border-transparent"
                  placeholder="Enter current password"
                />
                <input
                  type="text"
                  class="w-full text-tertiary px-3 py-2 my-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#10a20b] focus:border-transparent"
                  placeholder="Enter new password"
                />
                <button
                  class="px-4 py-2 mt-2 text-white transition rounded-md cursor-pointer disabled:cursor-default bg-primary hover:bg-primary/80 disabled:opacity-50"
                  :disabled="!currentPassword || !newPassword || currentPassword === newPassword"
                >
                  Update Password
                </button>
              </div>
            </div>
          </div>

          <div class="p-6 mb-6 rounded-lg shadow-md bg-background">
            <h2 class="mb-4 text-xl font-semibold text-gray-700">Preferences</h2>
            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700">Enable Notifications</span>
                <button
                  class="relative inline-flex items-center h-6 bg-gray-200 rounded-full cursor-pointer w-11"
                >
                  <span
                    class="inline-block w-4 h-4 transition transform translate-x-1 bg-white rounded-full"
                  ></span>
                </button>
              </div>
            </div>
          </div>

          <div class="p-6 rounded-lg shadow-md bg-background">
            <h2 class="mb-4 text-xl font-semibold text-gray-700">API settings</h2>
            <div class="space-y-4">
              <div>
                <label class="block mb-2 text-sm font-medium text-tertiary">
                  Provide custom API key
                </label>
                <input
                  type="text"
                  v-model="apiKey"
                  class="w-full text-tertiary px-3 py-2 my-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#10a20b] focus:border-transparent"
                  placeholder="Enter new API key"
                />
                <button
                  type="button"
                  class="px-4 py-2 mt-2 mr-2 text-white transition bg-green-600 rounded-md cursor-pointer disabled:cursor-default hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-300 disabled:opacity-50"
                  :disabled="!apiKey"
                  @click="updateAPIKey"
                >
                  Update API Key
                </button>
                <button
                  type="button"
                  class="px-4 py-2 mt-2 text-white transition bg-red-600 rounded-md cursor-pointer hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-300 disabled:opacity-50"
                  @click="resetAPIKey"
                >
                  Reset API Key
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import NavbarComponent from '@/components/NavbarComponent.vue';
import SideBar from '@/components/SideBar.vue';
import { useSidebarStore } from '@/stores/sidebarStore';
import { useAuthStore } from '@/stores/authStore';
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { changeAPIKey, resetAPIKey as resetAPIKeyAPI } from '@/api/users';
import { fireTimedAlert } from '@/lib/sweetAlerts';

const sidebarStore = useSidebarStore();
const authStore = useAuthStore();
const router = useRouter();

const email = ref('');
const currentPassword = ref('');
const newPassword = ref('');
const apiKey = ref('');

async function updateAPIKey() {
  if (apiKey.value) {
    try {
      await changeAPIKey(apiKey.value);

      // Refresh the token to get updated permissions
      await authStore.refreshAccountTokens();

      fireTimedAlert('Success', 'API key has been updated.', 2000);
      apiKey.value = ''; // Clear the input
    } catch (error) {
      console.error('Error updating API key:', error);
      fireTimedAlert('Error', 'Failed to update API key', 2000);
    }
  }
}

async function resetAPIKey() {
  try {
    await resetAPIKeyAPI();

    // Refresh the token to get updated permissions
    await authStore.refreshAccountTokens();

    fireTimedAlert('Success', 'API key has been reset to default.', 2000);
    apiKey.value = ''; // Clear the input
  } catch (error) {
    console.error('Error resetting API key:', error);
    fireTimedAlert('Error', 'Failed to reset API key', 2000);
  }
}
</script>
