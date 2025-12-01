<template>
  <div class="min-h-screen bg-secondary">
    <NavbarComponent />
    <div class="flex min-h-screen">
      <SideBar v-model="sidebarStore.visible">
        <ul class="space-y-2">
          <li>
            <button
              @click="router.push({ name: 'home' })"
              class="cursor-pointer text-background font-medium hover:text-indigo-100 transition text-left"
            >
              Home
            </button>
          </li>
          <li>
            <button
              @click="router.push({ name: 'list' })"
              class="cursor-pointer text-background font-medium hover:text-indigo-100 transition text-left"
            >
              Analysis Library
            </button>
          </li>
          <li class="self-bottom">
            <button
              @click="router.push({ name: 'settings' })"
              class="cursor-pointer text-background font-medium hover:text-indigo-100 transition text-left"
            >
              Settings
            </button>
          </li>
        </ul>
      </SideBar>
      <div class="px-8 py-8 max-w-[1400px] mx-auto flex-1">
        <div class="flex flex-1 justify-center items-center mb-8">
          <h1 class="m-0 text-3xl font-extrabold text-tertiary">Settings</h1>
        </div>
        <div class="flex flex-col mb-8">
          <div class="bg-background rounded-lg shadow-md p-6 mb-6">
            <h2 class="text-xl font-semibold text-tertiary mb-4">Account Settings</h2>
            <div class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-tertiary mb-2">
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
                  class="cursor-pointer disabled:cursor-default mt-2 px-4 py-2 bg-primary text-white rounded-md hover:bg-primary/80 transition disabled:opacity-50"
                >
                  Update Email
                </button>
              </div>
              <div>
                <label class="block text-sm font-medium text-tertiary mb-2">
                  Change password
                </label>
                <input
                  type="password"
                  class="w-full text-tertiary px-3 py-2 my-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#10a20b] focus:border-transparent"
                  placeholder="Enter current password"
                />
                <input
                  type="password"
                  class="w-full text-tertiary px-3 py-2 my-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#10a20b] focus:border-transparent"
                  placeholder="Enter new password"
                />
                <button
                  class="cursor-pointer disabled:cursor-default mt-2 px-4 py-2 bg-primary text-white rounded-md hover:bg-primary/80 disabled:opacity-50 disabled:cursor-default transition"
                  :disabled="!currentPassword || !newPassword || currentPassword === newPassword"
                >
                  Update Password
                </button>
              </div>
            </div>
          </div>

          <div class="bg-background rounded-lg shadow-md p-6 mb-6">
            <h2 class="text-xl font-semibold text-gray-700 mb-4">Preferences</h2>
            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700">Enable Notifications</span>
                <button
                  class="cursor-pointer bg-gray-200 relative inline-flex h-6 w-11 items-center rounded-full"
                >
                  <span
                    class="translate-x-1 inline-block h-4 w-4 transform rounded-full bg-white transition"
                  ></span>
                </button>
              </div>
            </div>
          </div>

          <div class="bg-background rounded-lg shadow-md p-6">
            <h2 class="text-xl font-semibold text-gray-700 mb-4">API settings</h2>
            <div class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-tertiary mb-2">
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
                  class="cursor-pointer disabled:cursor-default mt-2 px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-300 disabled:opacity-50 transition mr-2"
                  :disabled="!apiKey"
                  @click="updateAPIKey"
                >
                  Update API Key
                </button>
                <button
                  type="button"
                  class="cursor-pointer mt-2 px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-300 disabled:opacity-50 transition"
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
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { changeAPIKey } from '@/api/users';
import { fireTimedAlert } from '@/lib/sweetAlerts';

const sidebarStore = useSidebarStore();
const router = useRouter();

const email = ref('');
const currentPassword = ref('');
const newPassword = ref('');
const apiKey = ref('');

async function updateAPIKey() {
  if (apiKey.value) {
    await changeAPIKey(apiKey.value)
      .then(() => {
        fireTimedAlert('Success', 'API key has been updated.', 2000);
      })
      .catch((error) => {
        console.error('Error updating API key:', error);
        fireTimedAlert('Error', 'Failed to update API key', 2000);
      });
  }
}

async function resetAPIKey() {
  await changeAPIKey('a')
    .then(() => {
      fireTimedAlert('Success', 'API key has been reset to default.', 2000);
    })
    .catch((error) => {
      console.error('Error resetting API key:', error);
      fireTimedAlert('Error', 'Failed to reset API key', 2000);
    });
}
</script>
