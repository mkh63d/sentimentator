<template>
  <transition name="fade">
    <div
      v-if="visibleLocal"
      class="fixed inset-0 flex items-center justify-center bg-black/45 z-9999"
    >
      <div class="text-center text-white">
        <div
          class="w-[55px] h-[55px] border-4 border-gray-300 border-t-[#42b883] rounded-full mx-auto animate-spin"
        ></div>
        <p v-if="messageLocal" class="mt-3 text-[1.1rem]">{{ messageLocal }}</p>
      </div>
    </div>
  </transition>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps({
  visible: {
    type: Boolean,
    default: false,
  },
  message: {
    type: String,
    default: '',
  },
});

const emit = defineEmits(['update:visible', 'update:message']);

const visibleLocal = computed({
  get: () => props.visible,
  set(value) {
    emit('update:visible', value);
  },
});

const messageLocal = computed({
  get: () => props.message,
  set(value) {
    emit('update:message', value);
  },
});
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
