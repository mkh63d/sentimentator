<template>
  <nav class="px-2 pt-2 shadow-2xl shadow-gray-500 dark:shadow-black/40">
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
        <router-link to="/home" class="px-5 text-2xl font-bold text-background">
          Sentimentator
        </router-link>
      </div>
      <div class="flex items-center gap-3">
        <LanguageSwitcher class="hidden sm:block mx-0" />
        <ThemeSwitcher class="mx-5 sm:mx-3" :changeOnTheme="true" />
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import ThemeSwitcher from '@/components/ThemeSwitcher.vue';
import LanguageSwitcher from '@/components/LanguageSwitcher.vue';
import { useSidebarStore } from '@/stores/sidebarStore';

const sidebarStore = useSidebarStore();
</script>
