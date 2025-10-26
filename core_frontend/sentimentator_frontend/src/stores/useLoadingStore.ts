// stores/useLoadingStore.ts
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useLoadingStore = defineStore('loading', () => {
  const visible = ref(false)
  const message = ref<string>('')

  function show(msg?: string) {
    visible.value = true
    message.value = msg ?? 'Loading...'
  }

  function hide() {
    visible.value = false
    message.value = ''
  }

  return { visible, message, show, hide }
})
