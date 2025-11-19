<template>
  <aside
    class="h-auto mb-2 transition-all duration-300 ease-in-out border-r-fourth bg-primary shadow-xl dark:shadow-black/40 shadow-gray-500 z-50 p-4 shrink-0"
    :class="sidebarStore.visible ? 'w-64' : 'w-0 overflow-hidden p-0'"
  >
    <div v-if="sidebarStore.visible" class="relative h-full">
      <button
        class="cursor-pointer absolute mt-1 mr-3 scale-200 top-2 right-2 text-background hover:text-secondary transition"
        @click="sidebarStore.hide()"
      >
        <X size="16" />
      </button>
      <div class="p-4 pt-12">
        <slot></slot>
      </div>
      <div>
        <button
          class="cursor-pointer w-full mt-4 py-2 px-4 bg-secondary text-tertiary font-bold rounded hover:bg-secondary/80 transition"
          @click="
            () => {
              authStore.logoutUser();
              sidebarStore.hide();
              router.push({ name: 'login' });
            }
          "
        >
          Logout
        </button>
      </div>
    </div>
  </aside>
</template>

<script setup>
defineOptions({ name: 'SideBar' });

import { X } from 'lucide-vue-next';
import { useSidebarStore } from '@/stores/sidebarStore';
import { useAuthStore } from '@/stores/authStore';
import { useRouter } from 'vue-router';

const router = useRouter();
const authStore = useAuthStore();
const sidebarStore = useSidebarStore();
</script>
