<template>
  <div
    class="w-[90%] rounded-3xl mx-auto px-4 py-10 bg-[radial-gradient(circle_at_top_left,var(--color-tertiary),var(--color-primary),transparent_70%),radial-gradient(circle_at_bottom_right,var(--color-primary),var(--color-tertiary),transparent_80%)] bg-blend-overlay"
  >
    <div class="relative w-full translate-y-80 z-50">
      <button
        @click="scrollLeft"
        class="absolute top-0 left-0 z-10 p-2 text-white -translate-y-1/2 bg-gray-800 rounded-full hover:bg-gray-700"
        aria-label="Scroll team carousel left"
      >
        ‹
      </button>
      <button
        @click="scrollRight"
        class="absolute top-0 right-0 z-10 p-2 text-white -translate-y-1/2 bg-gray-800 rounded-full hover:bg-gray-700"
        aria-label="Scroll team carousel right"
      >
        ›
      </button>
    </div>
    <h2 class="pt-5 text-5xl font-bold text-center text-white mb-15">Meet the Team</h2>
    <div ref="scrollContainer" class="px-4 overflow-x-auto scrollbar-hide scroll-smooth">
      <div class="flex py-4 space-x-4">
        <div
          v-for="member in teamMembers"
          :key="member.name"
          class="min-w-[200px] shrink-0 overflow-hidden rounded-2xl shadow-lg relative group"
        >
          <img
            :src="member.photo1x"
            :srcset="`${member.photo1x} 1x, ${member.photo2x} 2x`"
            :alt="member.name"
            class="object-cover object-center w-48 h-80"
            loading="lazy"
            decoding="async"
          />
          <h3
            class="absolute bottom-0 left-0 right-0 py-2 text-2xl font-bold text-center text-white bg-gradient-to-t from-black/70 to-transparent"
          >
            {{ member.name }}
          </h3>
          <div
            class="absolute inset-0 bg-black/70 backdrop-blur-sm opacity-0 group-hover:opacity-100 transition-opacity duration-300 flex flex-col items-center justify-center gap-4"
          >
            <a
              :href="member.linkedin"
              target="_blank"
              rel="noopener noreferrer"
              class="px-6 py-2 text-white bg-blue-600 rounded-lg hover:bg-blue-700 transition"
            >
              LinkedIn
            </a>
            <a
              :href="member.github"
              target="_blank"
              rel="noopener noreferrer"
              class="px-6 py-2 text-white bg-gray-800 rounded-lg hover:bg-gray-900 transition"
            >
              GitHub
            </a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue';

import defaultPortrait from '@/assets/team/default.jpg';
import mhPortrait from '@/assets/team/MH.jpg';
import okPortrait from '@/assets/team/OK.jpg';

const scrollContainer = ref<HTMLDivElement | null>(null);
const doAutoScroll = ref(true);
const scrollAmount = 200;
let autoScroll: ReturnType<typeof setInterval> | undefined;

const teamMembers = [
  {
    name: 'Oskar Kacprzak',
    photo1x: okPortrait,
    photo2x: okPortrait,
    linkedin: 'https://www.linkedin.com/in/oskar-kacprzak',
    github: 'https://github.com/Oskarowski',
  },
  {
    name: 'Michał Hołyński',
    photo1x: mhPortrait,
    photo2x: mhPortrait,
    linkedin: 'https://www.linkedin.com/in/michal-krz-holynski',
    github: 'https://github.com/mkh63d',
  },
  {
    name: 'Wojciech Kapica',
    photo1x: defaultPortrait,
    photo2x: defaultPortrait,
    linkedin: 'https://www.linkedin.com/in/your-linkedin-profile',
    github: 'https://github.com/Soplicow',
  },
  {
    name: 'Błazej Dymek',
    photo1x: defaultPortrait,
    photo2x: defaultPortrait,
    linkedin: 'https://www.linkedin.com/in/your-linkedin-profile',
    github: 'https://github.com/BlazekWasTaken',
  },
  {
    name: 'Natalia Kośmider',
    photo1x: defaultPortrait,
    photo2x: defaultPortrait,
    linkedin: 'https://www.linkedin.com/in/your-linkedin-profile',
    github: 'https://github.com/jumpidolek',
  },
  {
    name: 'Anna Krzempek',
    photo1x: defaultPortrait,
    photo2x: defaultPortrait,
    linkedin: 'https://www.linkedin.com/in/your-linkedin-profile',
    github: 'https://github.com/your-github-profile',
  },
  {
    name: 'Jeremi Cichawa',
    photo1x: defaultPortrait,
    photo2x: defaultPortrait,
    linkedin: 'https://www.linkedin.com/in/your-linkedin-profile',
    github: 'https://github.com/Nekyi4',
  },
];

const scrollBy = (delta: number) => {
  doAutoScroll.value = false;
  scrollContainer.value?.scrollBy({ left: delta, behavior: 'smooth' });
};

const scrollRight = () => scrollBy(scrollAmount);
const scrollLeft = () => scrollBy(-scrollAmount);

const handleWheel = (e: WheelEvent) => {
  if (e.shiftKey || e.deltaX !== 0) {
    doAutoScroll.value = false;
  }
};

onMounted(() => {
  autoScroll = setInterval(() => {
    if (!scrollContainer.value || !doAutoScroll.value) return;
    scrollContainer.value.scrollLeft += 1;
  }, 40);

  scrollContainer.value?.addEventListener('wheel', handleWheel);
});

onBeforeUnmount(() => {
  if (autoScroll) clearInterval(autoScroll);
  scrollContainer.value?.removeEventListener('wheel', handleWheel);
});
</script>
