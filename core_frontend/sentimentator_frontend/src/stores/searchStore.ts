import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useSearchStore = defineStore('search', () => {
  const visible = ref(false);

  function show() {
    visible.value = true;
  }

  function hide() {
    visible.value = false;
  }

  return { visible, show, hide };
});
