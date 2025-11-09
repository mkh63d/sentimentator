<template>
  <button
    @click="toggleTheme"
    class="group scale-400 relative flex items-center justify-center h-10 rounded-full transition-colors duration-300 mr-8"
  >
    <!-- SVG container -->
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 40 40"
      class="w-6 h-6 stroke-1 transition-transform duration-500"
      :class="changeOnTheme && currentTheme === 'dark' ? 'stroke-black' : 'stroke-white'"
    >
      <!-- Sun -->
      <circle
        cx="20"
        cy="20"
        r="6"
        fill="none"
        class="transition-all duration-500"
        :class="currentTheme === 'dark' ? 'opacity-35' : 'opacity-100'"
      />

      <!-- Moon (croissant) -->
      <g
        class="transition-all duration-500"
        :class="currentTheme === 'dark' ? 'opacity-100' : 'opacity-35'"
      >
        <path
          fill="transparent"
          style="stroke-width: 2px"
          d="M29.14,28.77c-0.06-0.3-0.3-0.54-0.61-0.59c-5.84-1.1-10.41-5.56-11.65-11.37 c-1.24-5.81,1.1-11.75,5.97-15.15c0.25-0.18,0.38-0.49,0.31-0.79c-0.06-0.3-0.3-0.54-0.61-0.59 c-2.11-0.4-4.29-0.37-6.41,0.08c-4.24,0.91-7.87,3.41-10.23,7.05c-2.36,3.64-3.15,7.98-2.25,12.22 c1.59,7.44,8.26,12.84,15.86,12.84c0,0,0,0,0,0c1.14,0,2.29-0.12,3.41-0.36c2.13-0.46,4.1-1.31,5.88-2.55 C29.08,29.38,29.2,29.07,29.14,28.77z"
          transform="translate(21, 13.5) scale(0.45, 0.4)"
        />
      </g>
    </svg>
  </button>
</template>

<script setup>
import { ref, onMounted, watchEffect } from 'vue';

const props = defineProps({
  changeOnTheme: {
    type: Boolean,
    default: false
  }
});

const themes = ['light', 'dark'];
const currentTheme = ref(localStorage.getItem('theme') || 'light');

function applyTheme(theme) {
  document.documentElement.setAttribute('data-theme', theme);
  localStorage.setItem('theme', theme);
}

function toggleTheme() {
  const nextIndex = (themes.indexOf(currentTheme.value) + 1) % themes.length;
  const nextTheme = themes[nextIndex];
  changeTheme(nextTheme);
}

function changeTheme(theme) {
  currentTheme.value = theme;
  applyTheme(theme);
}

// Make sure it's applied on first mount
onMounted(() => {
  applyTheme(currentTheme.value);
});

// Optional: react to external changes (e.g., system or other components)
watchEffect(() => {
  if (document.documentElement.getAttribute('data-theme') !== currentTheme.value) {
    currentTheme.value = document.documentElement.getAttribute('data-theme') || 'light';
  }
});
</script>
