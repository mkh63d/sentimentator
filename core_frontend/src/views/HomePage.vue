<template>
  <div class="min-h-screen bg-background">
    <nav
      :class="[
        'fixed top-0 left-0 right-0 z-50 px-4 py-5 bg-white/ backdrop-blur-md transition-transform duration-300',
        showNavbar ? 'translate-y-0' : '-translate-y-full',
      ]"
    >
      <div class="flex flex-wrap md:flex-no-wrap items-center justify-between mx-auto max-w-7xl">
        <div
          class="text-2xl font-bold cursor-pointer text-background"
          @click="
            () => {
              router.push({ name: 'home' });
              showPhoneMenu = false;
            }
          "
        >
          Sentimentator
        </div>
        <ul class="hidden gap-8 md:flex">
          <li class="order-2 md:order-1">
            <router-link
              :to="{ hash: '#features' }"
              class="font-medium transition text-background hover:text-indigo-100"
              >Features</router-link
            >
          </li>
          <li class="order-2 md:order-2">
            <router-link
              :to="{ hash: '#methods' }"
              class="font-medium transition text-background hover:text-indigo-100"
              >Methods</router-link
            >
          </li>
          <li class="order-2 md:order-3">
            <router-link
              :to="{ hash: '#about' }"
              class="font-medium transition text-background hover:text-indigo-100"
              >About</router-link
            >
          </li>
          <li class="order-2 md:order-4">
            <router-link
              :to="{ hash: '#team' }"
              class="font-medium transition text-background hover:text-indigo-100"
              >Team</router-link
            >
          </li>
        </ul>
        <ul class="hidden gap-8 md:flex">
          <li>
            <button
              class="px-6 py-3 font-semibold transition rounded-lg cursor-pointer bg-background text-primary hover:bg-secondary"
              @click="hasTokens ? goToDashboard() : goToLogin()"
            >
              Try for Free
            </button>
          </li>
          <li>
            <button
              class="px-6 py-2 mr-3 font-semibold transition bg-transparent rounded-lg cursor-pointer text-background border-3 border-background"
              @click="goToLogin"
            >
              Login
            </button>
          </li>
          <li>
            <ThemeSwitcher class="mx-5 ml-0" :changeOnTheme="true" />
          </li>
        </ul>
        <button class="md:hidden mr-2" @click="togglePhoneMenu">
          <svg
            class="w-6 h-6 scale-[130%]"
            viewBox="0 0 24 24"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
            v-if="!showPhoneMenu"
          >
            <path
              d="M4 6H20M4 12H20M4 18H20"
              stroke="var(--color-secondary)"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            ></path>
          </svg>
          <svg
            class="w-6 h-6 scale-[130%]"
            viewBox="0 0 24 24"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
            v-else
          >
            <path
              d="M6 18L18 6M6 6L18 18"
              stroke="var(--color-secondary)"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            ></path>
          </svg>
        </button>
        <div v-show="showPhoneMenu" class="md:hidden relative w-full h-full">
          <div class="flex flex-col justify-between items-center text-center h-[90vh] pt-5 z-10">
            <div
              class="p-2 flex-1 font-medium text-3xl w-full transition text-background items-center justify-center flex hover:text-indigo-100 border-b border-gray-900"
            >
              <router-link :to="{ hash: '#features' }" class="z-10 px-4">Features</router-link>
            </div>

            <div
              class="p-2 flex-1 font-medium text-3xl w-full transition text-background items-center justify-center flex hover:text-indigo-100 border-b border-gray-900"
              v-on:click="showPhoneMenu = false"
            >
              <router-link :to="{ hash: '#methods' }" class="z-10 px-4">Methods</router-link>
            </div>

            <div
              class="p-2 flex-1 font-medium text-3xl w-full transition text-background items-center justify-center flex hover:text-indigo-100 border-b border-gray-900"
              v-on:click="showPhoneMenu = false"
            >
              <router-link :to="{ hash: '#about' }" class="z-10 px-4">About</router-link>
            </div>

            <div
              class="p-2 flex-1 font-medium text-3xl w-full transition text-background items-center justify-center flex hover:text-indigo-100 border-b border-gray-900"
              v-on:click="showPhoneMenu = false"
            >
              <router-link :to="{ hash: '#team' }" class="z-10 px-4">Team</router-link>
            </div>
            <div class="flex space-x-5 m-5 justify-center">
              <button
                class="px-4 py-3 font-semibold transition rounded-lg cursor-pointer bg-background text-primary hover:bg-secondary"
                @click="hasTokens ? goToDashboard() : goToLogin()"
              >
                Try for Free
              </button>

              <button
                class="px-6 py-2 font-semibold transition bg-transparent rounded-lg cursor-pointer text-background border-3 border-background"
                @click="goToLogin"
              >
                Login
              </button>
            </div>

            <ThemeSwitcher class="mx-5 mb-10" :changeOnTheme="true" />
          </div>
        </div>
      </div>
    </nav>
    <!-- Hero Section -->
    <section
      id="hero"
      class="relative z-0 bg-linear-to-br from-primary via-secondary to-fourth md:pt-32 md:pb-40"
    >
      <div class="px-4 py-20 max-h-full mx-auto text-center max-w-7xl md:py-32 text-background">
        <div class="absolute md:top-10 md:right-10 hidden md:block">
          <ThemeSwitcher :changeOnTheme="true" />
        </div>
        <h1 class="mb-5 text-5xl font-extrabold md:text-6xl">Analyze Video Reviews in Seconds</h1>
        <p class="mb-20 text-xl opacity-90">
          Transform your favourite video reviews into actionable sentiment insights with ease.
        </p>
        <div class="flex flex-col justify-center gap-5 sm:flex-row">
          <button
            class="px-8 py-4 text-lg font-semibold transition rounded-lg cursor-pointer bg-background text-primary hover:bg-secondary"
            @click="hasTokens ? goToDashboard() : goToLogin()"
          >
            Start Now!
          </button>
          <button
            class="px-8 py-4 text-lg font-semibold transition bg-transparent border-2 rounded-lg cursor-pointer text-background border-background hover:bg-background hover:text-primary"
            @click="goToLogin"
          >
            Sign In
          </button>
        </div>
      </div>
    </section>

    <section id="compatible">
      <div class="px-1 py-1 mx-auto bg-white max-w-12xl">
        <ul class="flex items-center justify-center gap-5 md:gap-15">
          <li>
            <h3 class="text-md md:text-md font-bold p-5 md:px-20 text-[#3d4d5c] border-r-2">
              Compatible with <br />
              Popular Platforms
            </h3>
          </li>
          <li>
            <svg viewBox="0 0 100 70" role="img" aria-label="YouTube" class="w-20 md:w-40">
              <rect width="100" height="70" rx="12" fill="#ff0000" />
              <polygon points="40,20 70,35 40,50" fill="#fff" />
            </svg>
          </li>
          <li>
            <svg viewBox="0 0 100 100" role="img" aria-label="Instagram" class="w-20 md:w-40">
              <defs>
                <linearGradient id="instaGradient" x1="0%" y1="0%" x2="100%" y2="100%">
                  <stop offset="0%" stop-color="#f58529" />
                  <stop offset="50%" stop-color="#dd2a7b" />
                  <stop offset="100%" stop-color="#8134af" />
                </linearGradient>
              </defs>
              <rect x="10" y="10" width="80" height="80" rx="22" fill="url(#instaGradient)" />
              <circle cx="50" cy="50" r="23" fill="none" stroke="#fff" stroke-width="8" />
              <circle cx="74" cy="26" r="6" fill="#fff" />
            </svg>
          </li>
        </ul>
      </div>
    </section>

    <!-- Features Section -->
    <section id="features" class="py-2 bg-tertiary">
      <div
        class="flex flex-wrap justify-center px-0 mx-0 divide-x divide-y max-w-8xl divide-primary"
      >
        <div class="w-[50%] text-center p-12 col-1">
          <h3 class="mb-1 text-xl font-bold text-secondary text-align-left">Available non-stop</h3>
          <h2 class="mb-1 text-2xl font-extrabold text-primary text-align-left">
            24/7 - on demand,
          </h2>
          <h3 class="mb-1 text-xl font-bold text-secondary text-align-left">
            whenever you need it
          </h3>
        </div>
        <div class="text-center w-[30%] p-12 col-1 border-r-0">
          <h2 class="mb-4 text-2xl font-bold text-primary">Accurate Results</h2>
          <p class="text-textBase">
            Our advanced algorithms ensure you get the most accurate sentiment insights.
          </p>
        </div>
        <div class="text-center w-[40%] p-6 col-2">
          <h2 class="mb-4 text-2xl font-bold text-primary">User-Friendly</h2>
          <p class="text-textBase">
            An intuitive interface makes it easy for anyone to analyze video reviews.
          </p>
        </div>
        <div class="text-center p-6 col-2 w-[60%] border-r-0">
          <h2 class="mb-4 text-2xl font-bold text-primary">In One Place</h2>
          <p class="text-textBase">
            All your data for your next better consumer choice in one place.
          </p>
        </div>
        <div class="w-full p-6 text-center col-1">
          <h2 class="mb-4 text-2xl font-bold text-primary">Be a Conscious Customer</h2>
          <p class="text-textBase">
            Make informed decisions based on a comprehensive sentiment analysis.
          </p>
        </div>
      </div>
    </section>

    <section id="methods" class="bg-secondary">
      <div class="px-4 mx-auto text-center max-w-7xl py-18">
        <h2 class="mb-5 text-4xl font-extrabold text-primary">Our Analysis Methods</h2>
        <p class="mb-10 text-xl text-tertiary">
          Try out our cutting-edge techniques to deliver precise sentiment analysis.
        </p>
        <div class="grid grid-cols-1 gap-10 md:grid-cols-3">
          <div
            class="p-6 transition duration-700 border-2 rounded-lg shadow-md bg-secondary hover:shadow-lg opacity-40 hover:opacity-100 border-accent"
          >
            <h3 class="mb-4 text-2xl font-bold text-primary">Natural Language Processing (NLP)</h3>
            <p class="text-tertiary">
              Leveraging advanced NLP algorithms to understand and interpret human language in video
              reviews.
            </p>
          </div>
          <div
            class="p-6 transition duration-700 border-2 rounded-lg shadow-md bg-secondary hover:shadow-lg opacity-40 hover:opacity-100 border-accent"
          >
            <h3 class="mb-4 text-2xl font-bold text-primary">
              Transcript and highlight extraction
            </h3>
            <p class="text-tertiary">
              Extracting key transcripts and highlights from video content for focused sentiment
              analysis from multiple angles with just the link you provide.
            </p>
          </div>
          <div
            class="p-6 transition duration-700 border-2 rounded-lg shadow-md bg-secondary hover:shadow-lg opacity-40 hover:opacity-100 border-accent"
          >
            <h3 class="mb-4 text-2xl font-bold text-primary">Sentiment Scoring</h3>
            <p class="text-tertiary">
              Assigning sentiment scores to reviews to quantify positive, negative, and neutral
              feedback.
            </p>
          </div>
        </div>
      </div>
    </section>

    <section id="about">
      <div class="px-4 pb-12 mx-auto text-center max-w-7xl pt-18">
        <h2 class="mb-5 text-4xl font-extrabold text-primary">About Us</h2>
        <p class="mb-10 text-xl text-tertiary">
          Sentimentator is dedicated to empowering consumers with insights derived from video
          reviews. <br />
          Our mission is to make sentiment analysis accessible and actionable for everyone.
        </p>
      </div>
    </section>
    <section id="team" ref="teamSectionRef">
      <Suspense v-if="shouldRenderTeamCarousel">
        <template #default>
          <TeamCarousel />
        </template>
        <template #fallback>
          <div class="w-[90%] rounded-3xl mx-auto px-4 py-10 bg-slate-900/30">
            <div class="h-60 animate-pulse bg-slate-700 rounded-3xl" />
          </div>
        </template>
      </Suspense>
      <div
        v-else
        class="w-[90%] rounded-3xl mx-auto px-4 py-10 bg-slate-900/10 animate-pulse"
        aria-hidden="true"
      >
        <div class="h-60 bg-slate-800/30 rounded-3xl" />
      </div>
    </section>

    <!-- Footer -->
    <footer class="py-10 text-center bg-background text-primary">
      <div class="px-4 mx-auto max-w-7xl">
        <p>&copy; 2025 Sentimentator. All rights reserved.</p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, computed, defineAsyncComponent, onUnmounted } from 'vue';
import ThemeSwitcher from '@/components/ThemeSwitcher.vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/authStore';

const TeamCarousel = defineAsyncComponent(() => import('@/components/TeamCarousel.vue'));

const router = useRouter();
const authStore = useAuthStore();

const teamSectionRef = ref<HTMLElement | null>(null);
const shouldRenderTeamCarousel = ref(false);
let teamObserver: IntersectionObserver | null = null;
const showNavbar = ref(false);
const lastScroll = ref(0);
const showPhoneMenu = ref(false);

onMounted(() => {
  if (!('IntersectionObserver' in window)) {
    shouldRenderTeamCarousel.value = true;
    return;
  }

  teamObserver = new IntersectionObserver(
    (entries) => {
      if (entries.some((entry) => entry.isIntersecting)) {
        shouldRenderTeamCarousel.value = true;
        teamObserver?.disconnect();
        teamObserver = null;
      }
    },
    {
      rootMargin: '0px 0px -15% 0px',
      threshold: 0.2,
    },
  );

  if (teamSectionRef.value) {
    teamObserver.observe(teamSectionRef.value);
  }

  const hero = document.getElementById('hero');

  if (!hero) return;

  const heroBottom = hero.offsetTop + hero.offsetHeight;
  const methodsElements = document.querySelectorAll('#methods .opacity-40');

  // Show navbar only if below the hero section
  //
  const onScroll = () => {
    const current = window.scrollY;

    if (current > heroBottom) {
      showNavbar.value = true;
    } else {
      showNavbar.value = false;
      showPhoneMenu.value = false;
    }

    for (let i = 0; i < methodsElements.length; i++) {
      let method = methodsElements[i] as HTMLElement;
      if (
        current >
        method.getBoundingClientRect().top + window.scrollY - window.innerHeight + 100
      ) {
        method.classList.remove('opacity-40');
        method.classList.add('opacity-100');
      } else {
        method.classList.remove('opacity-100');
        method.classList.add('opacity-40');
      }
    }

    lastScroll.value = current;
  };

  window.addEventListener('scroll', onScroll);

  onUnmounted(() => {
    window.removeEventListener('scroll', onScroll);
  });
});

onBeforeUnmount(() => {
  teamObserver?.disconnect();
  teamObserver = null;
});

const hasTokens = computed(() => !!authStore.accessToken);

const goToLogin = () => {
  router.push({ name: 'login' });
};

const goToDashboard = () => {
  router.push({ name: 'list' });
};

const togglePhoneMenu = () => {
  showPhoneMenu.value = !showPhoneMenu.value;
};
</script>
