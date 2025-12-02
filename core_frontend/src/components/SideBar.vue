<template>
  <aside
    class="h-[85vh] mb-2 transition-all duration-300 ease-in-out border-r-fourth bg-primary shadow-xl dark:shadow-black/40 shadow-gray-500 z-50 p-4 shrink-0"
    :class="sidebarStore.visible ? 'w-64' : 'w-0 overflow-hidden p-0'"
  >
    <div v-if="sidebarStore.visible" class="relative">
      <!-- close button remains available immediately -->
      <button
        class="cursor-pointer absolute mt-1 mr-3 scale-200 top-2 right-2 text-background hover:text-secondary transition"
        @click="sidebarStore.hide()"
      >
        <X size="16" />
      </button>

      <!-- content is shown with a delay when sidebar opens -->
      <div v-if="showContent" class="p-4 pt-12">
        <slot></slot>
      </div>

      <div v-if="showContent" class="flex flex-col gap-3">
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
          {{ t('nav.logout') }}
        </button>
        <LanguageSwitcher class="sm:hidden block self-center" />
      </div>
    </div>
  </aside>
</template>

<script setup>
defineOptions({ name: 'SideBar' });

import { X } from 'lucide-vue-next';
import LanguageSwitcher from '@/components/LanguageSwitcher.vue';
import { useSidebarStore } from '@/stores/sidebarStore';
import { useAuthStore } from '@/stores/authStore';
import { useRouter } from 'vue-router';
import { ref, watch, onUnmounted } from 'vue';
import { useI18n } from 'vue-i18n';

const router = useRouter();
const authStore = useAuthStore();
const sidebarStore = useSidebarStore();
const { t } = useI18n();

const showContent = ref(false);
const OPEN_DELAY = 300; // match CSS duration-300
let openTimer = null;

watch(
  () => sidebarStore.visible,
  (visible) => {
    if (visible) {
      // delay showing inner content until after the width transition
      openTimer = setTimeout(() => {
        showContent.value = true;
        openTimer = null;
      }, OPEN_DELAY);
    } else {
      // hide immediately when closing and clear any pending timer
      if (openTimer) {
        clearTimeout(openTimer);
        openTimer = null;
      }
      showContent.value = false;
    }
  },
  { immediate: true },
);

onUnmounted(() => {
  if (openTimer) clearTimeout(openTimer);
});
</script>
