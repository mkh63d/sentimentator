<template>
  <nav class="px-2 pt-2 shadow-2xl shadow-gray-500">
    <div class="max-w-full mb-3 mx-auto flex justify-between items-center bg-primary p-3">
      <div class="flex justify-center items-center">
        <div
          class="pr-1 ml-5 mr-4 cursor-pointer"
          @click="sidebarStore.visible = !sidebarStore.visible"
        >
          <svg class="w-6 h-6" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path
              d="M4 6H20M4 12H20M4 18H20"
              stroke="var(--color-secondary)"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            ></path>
          </svg>
        </div>
        <div
          class="cursor-pointer px-5 text-2xl font-bold text-background"
          @click="() => router.push('/home')"
        >
          Sentimentator
        </div>
      </div>
      <ThemeSwitcher class="mx-5" :changeOnTheme="true" />
    </div>
  </nav>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
import ThemeSwitcher from './ThemeSwitcher.vue';
import { useSidebarStore } from '@/stores/sidebarStore';

const router = useRouter();
const sidebarStore = useSidebarStore();
</script>
