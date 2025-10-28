<template>
  <transition name="fade">
    <div v-if="visible" class="fixed inset-0 flex items-center justify-center bg-black/45 z-9999">
      <div class="text-center text-white">
        <div
          class="w-[55px] h-[55px] border-4 border-gray-300 border-t-[#42b883] rounded-full mx-auto animate-spin"
        ></div>
        <p v-if="message" class="mt-3 text-[1.1rem]">{{ message }}</p>
      </div>
    </div>
  </transition>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia';
import { useLoadingStore } from '@/stores/useLoadingStore';

const store = useLoadingStore();
const { visible, message } = storeToRefs(store);
</script>

<style scoped>
/* Fade transition */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.25s;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
