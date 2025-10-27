<template>
  <transition name="fade">
    <div v-if="visible" class="overlay">
      <div class="content">
        <div class="spinner"></div>
        <p v-if="message" class="message">{{ message }}</p>
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
.overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.45);
  z-index: 9999;
}

.content {
  text-align: center;
  color: white;
}

.spinner {
  width: 55px;
  height: 55px;
  border: 4px solid #ddd;
  border-top-color: #42b883;
  border-radius: 50%;
  margin: 0 auto;
  animation: spin 1s linear infinite;
}

.message {
  margin-top: 12px;
  font-size: 1.1rem;
}

/* Animation */
@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

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
