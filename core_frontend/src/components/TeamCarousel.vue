<template>
  <div
    class="w-[90%] rounded-3xl mx-auto px-4 py-10 bg-[radial-gradient(circle_at_top_left,var(--color-tertiary),var(--color-primary),transparent_70%),radial-gradient(circle_at_bottom_right,var(--color-primary),var(--color-tertiary),transparent_80%)] bg-blend-overlay"
  >
    <div class="relative w-full translate-y-80">
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
          class="min-w-[200px] shrink-0 overflow-hidden rounded-2xl shadow-lg"
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
            class="my-2 text-2xl font-bold text-center text-white text-shadow-xl text-shadow-gray-900"
          >
            {{ member.name }}
          </h3>
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
  },
  {
    name: 'Michał Hołyński',
    photo1x: mhPortrait,
    photo2x: mhPortrait,
  },
  {
    name: 'Wojciech Kapica',
    photo1x: defaultPortrait,
    photo2x: defaultPortrait,
  },
  {
    name: 'Błazej Dymek',
    photo1x: defaultPortrait,
    photo2x: defaultPortrait,
  },
  {
    name: 'Natalia Kośmider',
    photo1x: defaultPortrait,
    photo2x: defaultPortrait,
  },
  {
    name: 'Anna Krzempek',
    photo1x: defaultPortrait,
    photo2x: defaultPortrait,
  },
  {
    name: 'Jeremi Cichawa',
    photo1x: defaultPortrait,
    photo2x: defaultPortrait,
  },
];

const scrollBy = (delta: number) => {
  doAutoScroll.value = false;
  scrollContainer.value?.scrollBy({ left: delta, behavior: 'smooth' });
};

const scrollRight = () => scrollBy(scrollAmount);
const scrollLeft = () => scrollBy(-scrollAmount);

onMounted(() => {
  autoScroll = setInterval(() => {
    if (!scrollContainer.value || !doAutoScroll.value) return;
    scrollContainer.value.scrollLeft += 1;
  }, 40);
});

onBeforeUnmount(() => {
  if (autoScroll) clearInterval(autoScroll);
});
</script>
