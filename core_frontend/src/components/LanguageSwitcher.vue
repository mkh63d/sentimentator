<template>
  <div class="language-switcher">
    <button
      @click="toggleLanguage"
      class="px-3 py-2 rounded-lg transition bg-primary text-background hover:bg-secondary flex items-center gap-2"
      :title="currentLocale === 'en' ? 'Switch to Polish' : 'Przełącz na angielski'"
    >
      <span class="text-sm font-medium">{{ currentLocale.toUpperCase() }}</span>
      <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path
          d="M12 22C17.5228 22 22 17.5228 22 12C22 6.47715 17.5228 2 12 2C6.47715 2 2 6.47715 2 12C2 17.5228 6.47715 22 12 22Z"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        />
        <path
          d="M2 12H22"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        />
        <path
          d="M12 2C14.5013 4.73835 15.9228 8.29203 16 12C15.9228 15.708 14.5013 19.2616 12 22C9.49872 19.2616 8.07725 15.708 8 12C8.07725 8.29203 9.49872 4.73835 12 2Z"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        />
      </svg>
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

const { locale } = useI18n();

const currentLocale = computed(() => locale.value);

const toggleLanguage = () => {
  const newLocale = locale.value === 'en' ? 'pl' : 'en';
  locale.value = newLocale;
  localStorage.setItem('locale', newLocale);
};
</script>
