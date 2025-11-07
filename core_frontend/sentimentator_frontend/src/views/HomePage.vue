<template>
  <div class="min-h-screen bg-background">
    <!-- Hero Section -->
    <section class="bg-linear-to-br from-primary via-secondary to-[#011936]">
      <nav class="px-4 py-5">
        <div class="max-w-7xl mx-auto flex justify-between items-center">
          <div class="text-2xl font-bold text-textBase">Sentimentator</div>
          <ul class="hidden md:flex gap-8">
            <li>
              <a href="#features" class="text-textBase font-medium hover:text-indigo-100 transition"
                >Features</a
              >
            </li>
            <li>
              <a href="#methods" class="text-textBase font-medium hover:text-indigo-100 transition"
                >Methods</a
              >
            </li>
            <li>
              <a href="#about" class="text-textBase font-medium hover:text-indigo-100 transition"
                >About</a
              >
            </li>
            <li>
              <a href="#team" class="text-textBase font-medium hover:text-indigo-100 transition"
                >Team</a
              >
            </li>
          </ul>
          <ul class="hidden md:flex gap-8">
            <li>
              <button
                class="px-6 py-3 bg-background text-primary rounded-lg font-semibold hover:bg-secondary transition"
              >
                Try for Free
              </button>
            </li>
            <li>
              <button
                class="px-6 py-2 bg-transparent text-textBase rounded-lg font-semibold transition border-3 border-background"
              >
                Login
              </button>
            </li>
            <li>
              <ThemeSwitcher class="mx-5 ml-0" />
            </li>
          </ul>
        </div>
      </nav>

      <div class="max-w-7xl mx-auto px-4 py-20 md:py-32 text-center text-textBase">
        <h1 class="text-5xl md:text-6xl font-extrabold mb-5">Analyze Video Reviews in Seconds</h1>
        <p class="text-xl mb-10 opacity-90">
          Transform your favourite video reviews into actionable sentiment insights with ease.
        </p>
        <div class="flex flex-col sm:flex-row gap-5 justify-center">
          <button
            class="px-8 py-4 bg-background text-primary rounded-lg text-lg font-semibold hover:bg-secondary transition"
          >
            Start Now!
          </button>
          <button
            class="px-8 py-4 bg-transparent text-background border-2 border-background rounded-lg text-lg font-semibold hover:bg-background hover:text-primary transition"
          >
            Sign In
          </button>
        </div>
      </div>
    </section>

    <section id="compatible">
      <div class="max-w-12xl mx-auto px-1 py-1 bg-white">
        <ul class="flex flex-wrap justify-center items-center gap-15">
          <li>
            <h3 class="text-md md:text-md font-bold px-20 text-[#3d4d5c] border-r-2">
              Compatible with <br />
              Popular Platforms
            </h3>
          </li>
          <li></li>
          <li>
            <img
              src="@/assets/logos/youtube-logo.jpg"
              alt="YouTube"
              class="w-30 md:w-40 object-contain"
            />
          </li>
          <li>
            <img
              src="@/assets/logos/instagram-logo.png"
              alt="Instagram"
              class="w-30 md:w-40 object-contain"
            />
          </li>
          <li>
            <img
              src="@/assets/logos/tiktok-logo.png"
              alt="TikTok"
              class="h-8 md:h-8 object-contain"
            />
          </li>
        </ul>
      </div>
    </section>

    <!-- Features Section -->
    <section id="features" class="py-2 bg-tertiary">
      <div
        class="max-w-8xl mx-0 px-0 flex flex-wrap justify-center divide-y divide-x divide-primary"
      >
        <div class="w-[50%] text-center p-12 col-1">
          <h3 class="text-xl font-bold mb-1 text-secondary text-align-left">Available non-stop</h3>
          <h2 class="text-2xl font-extrabold mb-1 text-primary text-align-left">
            24/7 - on demand,
          </h2>
          <h3 class="text-xl font-bold mb-1 text-secondary text-align-left">
            whenever you need it
          </h3>
        </div>
        <div class="text-center w-[30%] p-12 col-1 border-r-0">
          <h2 class="text-2xl font-bold mb-4 text-primary">Accurate Results</h2>
          <p class="text-textBase">
            Our advanced algorithms ensure you get the most accurate sentiment insights.
          </p>
        </div>
        <div class="text-center w-[40%] p-6 col-2">
          <h2 class="text-2xl font-bold mb-4 text-primary">User-Friendly</h2>
          <p class="text-textBase">
            An intuitive interface makes it easy for anyone to analyze video reviews.
          </p>
        </div>
        <div class="text-center p-6 col-2 w-[60%] border-r-0">
          <h2 class="text-2xl font-bold mb-4 text-primary">In One Place</h2>
          <p class="text-textBase">
            All your data for your next better consumer choice in one place.
          </p>
        </div>
        <div class="text-center p-6 col-1 w-full">
          <h2 class="text-2xl font-bold mb-4 text-primary">Be a Conscious Customer</h2>
          <p class="text-textBase">
            Make informed decisions based on a comprehensive sentiment analysis.
          </p>
        </div>
      </div>
    </section>

    <section id="methods" class="bg-secondary">
      <div class="max-w-7xl mx-auto px-4 py-18 text-center">
        <h2 class="text-4xl font-extrabold mb-5 text-primary">Our Analysis Methods</h2>
        <p class="text-xl mb-10 text-tertiary">
          Try out our cutting-edge techniques to deliver precise sentiment analysis.
        </p>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-10">
          <div
            class="bg-secondary p-6 rounded-lg shadow-md hover:shadow-lg transition duration-700 opacity-40 hover:opacity-100 border-accent border-2"
          >
            <h3 class="text-2xl font-bold mb-4 text-primary">Natural Language Processing (NLP)</h3>
            <p class="text-tertiary">
              Leveraging advanced NLP algorithms to understand and interpret human language in video
              reviews.
            </p>
          </div>
          <div
            class="bg-secondary p-6 rounded-lg shadow-md hover:shadow-lg transition duration-700 opacity-40 hover:opacity-100 border-accent border-2"
          >
            <h3 class="text-2xl font-bold mb-4 text-primary">
              Transcript and highlight extraction
            </h3>
            <p class="text-tertiary">
              Extracting key transcripts and highlights from video content for focused sentiment
              analysis from multiple angles with just the link you provide.
            </p>
          </div>
          <div
            class="bg-secondary p-6 rounded-lg shadow-md hover:shadow-lg transition duration-700 opacity-40 hover:opacity-100 border-accent border-2"
          >
            <h3 class="text-2xl font-bold mb-4 text-primary">Sentiment Scoring</h3>
            <p class="text-tertiary">
              Assigning sentiment scores to reviews to quantify positive, negative, and neutral
              feedback.
            </p>
          </div>
        </div>
      </div>
    </section>

    <section id="about">
      <div class="max-w-7xl mx-auto px-4 pt-18 pb-12 text-center">
        <h2 class="text-4xl font-extrabold mb-5 text-primary">About Us</h2>
        <p class="text-xl mb-10 text-tertiary">
          Sentimentator is dedicated to empowering consumers with insights derived from video
          reviews. <br />
          Our mission is to make sentiment analysis accessible and actionable for everyone.
        </p>
      </div>
    </section>
    <section id="team">
      <div
        class="w-[90%] rounded-3xl mx-auto px-4 py-10 bg-[radial-gradient(circle_at_top_left,var(--color-tertiary),var(--color-primary),transparent_70%),radial-gradient(circle_at_bottom_right,var(--color-primary),var(--color-tertiary),transparent_80%)] bg-blend-overlay"
      >
        <div class="relative translate-y-80 w-full">
          <!-- Left Arrow -->
          <button
            @click="scrollLeft"
            class="absolute left-0 top-0 -translate-y-1/2 bg-gray-800 text-white p-2 rounded-full z-10 hover:bg-gray-700"
          >
            ‹
          </button>
          <button
            @click="scrollRight"
            class="absolute right-0 top-0 -translate-y-1/2 bg-gray-800 text-white p-2 rounded-full z-10 hover:bg-gray-700"
          >
            ›
          </button>
        </div>
        <h2 class="text-5xl pt-5 font-bold mb-15 text-white text-center">Meet the Team</h2>
        <div ref="scrollContainer" class="overflow-x-auto scrollbar-hide px-4 scroll-smooth">
          <div class="flex space-x-4 py-4">
            <div class="min-w-[200px] shrink-0 overflow-hidden rounded-2xl shadow-lg">
              <img
                src="@/assets/team/OK.jpg"
                alt="Oskar Kacprzak"
                class="w-48 h-80 object-cover object-center"
              />
              <h3
                class="text-2xl font-bold my-2 text-white text-shadow-xl text-shadow-gray-900 text-center"
              >
                Oskar Kacprzak
              </h3>
            </div>
            <div class="min-w-[200px] shrink-0 overflow-hidden rounded-2xl shadow-lg">
              <img
                src="@/assets/team/MH.jpg"
                alt="Michał Hołyński"
                class="w-48 h-80 object-cover object-center"
              />
              <h3
                class="text-2xl font-bold my-2 text-white text-shadow-xl text-shadow-gray-900 text-center"
              >
                Michał Hołyński
              </h3>
            </div>
            <div class="min-w-[200px] shrink-0 overflow-hidden rounded-2xl shadow-lg">
              <img
                src="@/assets/team/default.jpg"
                alt="Wojciech Kapica"
                class="w-48 h-80 object-cover object-center"
              />
              <h3
                class="text-2xl font-bold my-2 text-white text-shadow-xl text-shadow-gray-900 text-center"
              >
                Wojciech Kapica
              </h3>
            </div>
            <div class="min-w-[200px] shrink-0 overflow-hidden rounded-2xl shadow-lg">
              <img
                src="@/assets/team/default.jpg"
                alt="Błazej Dymek"
                class="w-48 h-80 object-cover object-center"
              />
              <h3
                class="text-2xl font-bold my-2 text-white text-shadow-xl text-shadow-gray-900 text-center"
              >
                Błazej Dymek
              </h3>
            </div>
            <div class="min-w-[200px] shrink-0 overflow-hidden rounded-2xl shadow-lg">
              <img
                src="@/assets/team/default.jpg"
                alt="Natalia Kośmider"
                class="w-48 h-80 object-cover object-center"
              />
              <h3
                class="text-2xl font-bold my-2 text-white text-shadow-xl text-shadow-gray-900 text-center"
              >
                Natalia Kośmider
              </h3>
            </div>
            <div class="min-w-[200px] shrink-0 overflow-hidden rounded-2xl shadow-lg">
              <img
                src="@/assets/team/default.jpg"
                alt="Anna Krzempek"
                class="w-48 h-80 object-cover object-center"
              />
              <h3
                class="text-2xl font-bold my-2 text-white text-shadow-xl text-shadow-gray-900 text-center"
              >
                Anna Krzempek
              </h3>
            </div>
            <div class="min-w-[200px] shrink-0 overflow-hidden rounded-2xl shadow-lg">
              <img
                src="@/assets/team/default.jpg"
                alt="Jeremi Cichawa"
                class="w-48 h-80 object-cover object-center"
              />
              <h3
                class="text-2xl font-bold my-2 text-white text-shadow-xl text-shadow-gray-900 text-center"
              >
                Jeremi Cichawa
              </h3>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="bg-background text-primary py-10 text-center">
      <div class="max-w-7xl mx-auto px-4">
        <p>&copy; 2025 Sentimentator. All rights reserved.</p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue';
import ThemeSwitcher from '@/components/ThemeSwitcher.vue';

const scrollContainer = ref(null);

const doAutoScroll = ref(true);
const scrollAmount = 200;

const scrollRight = () => {
  doAutoScroll.value = false;
  scrollContainer.value.scrollBy({ left: scrollAmount, behavior: 'smooth' });
};

const scrollLeft = () => {
  doAutoScroll.value = false;
  scrollContainer.value.scrollBy({ left: -scrollAmount, behavior: 'smooth' });
};

let autoScroll;
onMounted(() => {
  autoScroll = setInterval(() => {
    if (!scrollContainer.value) return;
    if (!doAutoScroll.value) return;

    const el = scrollContainer.value;
    // Scroll right slowly
    el.scrollLeft += 1;
  }, 40); // smaller = faster (try 10–50)
});

onBeforeUnmount(() => {
  clearInterval(autoScroll);
});

const goToLogin = () => {
  // Logic to navigate to the login page
};
</script>
