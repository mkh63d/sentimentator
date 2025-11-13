<template>
  <aside
    class="h-screen transition-all duration-300 ease-in-out border-r-fourth bg-primary shadow-xl shadow-gray-500 z-50 p-4 shrink-0"
    :class="sidebarStore.visible ? 'w-64' : 'w-0 overflow-hidden p-0'"
  >
    <div v-if="sidebarStore.visible" class="relative h-full">
      <button
        class="absolute mt-1 mr-3 scale-200 top-2 right-2 text-background hover:text-secondary transition"
        @click="sidebarStore.hide()"
      >
        <X size="16" />
      </button>
      <div class="p-4 pt-12">
        <slot></slot>
      </div>
    </div>
  </aside>
</template>

<script setup>
import { X } from 'lucide-vue-next';
import { useSidebarStore } from '@/stores/sidebarStore';

const sidebarStore = useSidebarStore();
</script>
