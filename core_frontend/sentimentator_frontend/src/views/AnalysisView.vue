<template>
  <div class="min-h-screen bg-secondary">
    <NavbarComponent
      :isSidebarOpen="isSidebarOpen"
      @update:isSidebarOpen="isSidebarOpen = $event"
    />

    <div class="flex min-h-screen">
      <Sidebar v-model="isSidebarOpen" @update:modelValue="isSidebarOpen = $event">
        <ul class="space-y-2">
          <li>
            <button
              @click="router.push('/home')"
              class="text-background font-medium hover:text-indigo-100 transition text-left"
            >
              Home
            </button>
          </li>
          <li>
            <button
              @click="router.push('/list')"
              class="text-background font-medium hover:text-indigo-100 transition text-left"
            >
              Analysis Library
            </button>
          </li>
          <li>
            <button
              @click="router.push('/search')"
              class="text-background font-medium hover:text-indigo-100 transition text-left"
            >
              New Analysis
            </button>
          </li>
          <li class="self-bottom">
            <button
              @click="router.push('/settings')"
              class="text-background font-medium hover:text-indigo-100 transition text-left"
            >
              Settings
            </button>
          </li>
        </ul>
      </Sidebar>
      <div class="max-w-7xl mx-auto px-6 py-8">
        <!-- Page Header -->
        <div class="mb-8">
          <h1 class="m-0 text-3xl font-extrabold text-tertiary">Sentiment Analysis Results</h1>
          <p class="text-textBase">Smartphone review sentiment breakdown by features</p>
        </div>

        <!-- Overall Sentiment Summary -->
        <div class="bg-fourthary rounded-xl shadow-lg p-6 mb-6">
          <h2 class="text-2xl font-semibold text-secondary mb-4">Overall Sentiment</h2>
          <div
            class="flex flex-col md:flex-row items-center justify-center md:space-x-8 space-y-6 md:space-y-0"
          >
            <div class="text-center">
              <div class="text-6xl font-bold text-fourthary mb-2">Positive</div>
              <div class="text-lg text-secondary">Dominant Sentiment</div>
            </div>
            <div class="w-full md:w-auto">
              <SentimentChart :positive="65" :neutral="25" :negative="10" />
            </div>
            <div class="flex md:flex-col space-x-4 md:space-x-0 md:space-y-4">
              <div class="text-center p-4 bg-green-100 dark:bg-green-900 rounded-lg">
                <div class="text-2xl font-bold text-green-600 mb-1">65%</div>
                <div class="text-sm text-textBase">Positive</div>
              </div>
              <div class="text-center p-4 bg-gray-100 dark:bg-gray-700 rounded-lg">
                <div class="text-2xl font-bold text-gray-600 mb-1">25%</div>
                <div class="text-sm text-textBase">Neutral</div>
              </div>
              <div class="text-center p-4 bg-red-100 dark:bg-red-900 rounded-lg">
                <div class="text-2xl font-bold text-red-600 mb-1">10%</div>
                <div class="text-sm text-textBase">Negative</div>
              </div>
            </div>
          </div>
        </div>

        <!-- Feature-Based Sentiment Analysis -->
        <div class="bg-fourthary rounded-xl shadow-lg p-6 mb-6">
          <h2 class="text-xl font-semibold text-secondary mb-4">Sentiment by Smartphone Feature</h2>
          <div class="space-y-4">
            <!-- Camera Feature -->
            <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
              <div class="flex justify-between items-center mb-3">
                <h3 class="text-lg font-semibold text-secondary">📷 Camera</h3>
                <span
                  class="px-3 py-1 bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200 rounded-full text-sm font-semibold"
                >
                  Positive
                </span>
              </div>
              <div class="flex space-x-2 mb-2">
                <div class="flex-1 bg-green-500 h-3 rounded" style="width: 70%"></div>
                <div class="flex-1 bg-gray-300 h-3 rounded" style="width: 20%"></div>
                <div class="flex-1 bg-red-500 h-3 rounded" style="width: 10%"></div>
              </div>
              <div class="flex justify-between text-sm text-tertiary">
                <span>Positive: 70%</span>
                <span>Neutral: 20%</span>
                <span>Negative: 10%</span>
              </div>
            </div>

            <!-- Battery Feature -->
            <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
              <div class="flex justify-between items-center mb-3">
                <h3 class="text-lg font-semibold text-secondary">🔋 Battery</h3>
                <span
                  class="px-3 py-1 bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200 rounded-full text-sm font-semibold"
                >
                  Positive
                </span>
              </div>
              <div class="flex space-x-2 mb-2">
                <div class="flex-1 bg-green-500 h-3 rounded" style="width: 80%"></div>
                <div class="flex-1 bg-gray-300 h-3 rounded" style="width: 15%"></div>
                <div class="flex-1 bg-red-500 h-3 rounded" style="width: 5%"></div>
              </div>
              <div class="flex justify-between text-sm text-tertiary">
                <span>Positive: 80%</span>
                <span>Neutral: 15%</span>
                <span>Negative: 5%</span>
              </div>
            </div>

            <!-- Display Feature -->
            <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
              <div class="flex justify-between items-center mb-3">
                <h3 class="text-lg font-semibold text-secondary">📱 Display</h3>
                <span
                  class="px-3 py-1 bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200 rounded-full text-sm font-semibold"
                >
                  Positive
                </span>
              </div>
              <div class="flex space-x-2 mb-2">
                <div class="flex-1 bg-green-500 h-3 rounded" style="width: 75%"></div>
                <div class="flex-1 bg-gray-300 h-3 rounded" style="width: 20%"></div>
                <div class="flex-1 bg-red-500 h-3 rounded" style="width: 5%"></div>
              </div>
              <div class="flex justify-between text-sm text-tertiary">
                <span>Positive: 75%</span>
                <span>Neutral: 20%</span>
                <span>Negative: 5%</span>
              </div>
            </div>

            <!-- Performance Feature -->
            <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
              <div class="flex justify-between items-center mb-3">
                <h3 class="text-lg font-semibold text-secondary">⚡ Performance</h3>
                <span
                  class="px-3 py-1 bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200 rounded-full text-sm font-semibold"
                >
                  Neutral
                </span>
              </div>
              <div class="flex space-x-2 mb-2">
                <div class="flex-1 bg-green-500 h-3 rounded" style="width: 50%"></div>
                <div class="flex-1 bg-gray-300 h-3 rounded" style="width: 40%"></div>
                <div class="flex-1 bg-red-500 h-3 rounded" style="width: 10%"></div>
              </div>
              <div class="flex justify-between text-sm text-tertiary">
                <span>Positive: 50%</span>
                <span>Neutral: 40%</span>
                <span>Negative: 10%</span>
              </div>
            </div>

            <!-- Price Feature -->
            <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
              <div class="flex justify-between items-center mb-3">
                <h3 class="text-lg font-semibold text-secondary">💰 Price</h3>
                <span
                  class="px-3 py-1 bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200 rounded-full text-sm font-semibold"
                >
                  Negative
                </span>
              </div>
              <div class="flex space-x-2 mb-2">
                <div class="flex-1 bg-green-500 h-3 rounded" style="width: 30%"></div>
                <div class="flex-1 bg-gray-300 h-3 rounded" style="width: 20%"></div>
                <div class="flex-1 bg-red-500 h-3 rounded" style="width: 50%"></div>
              </div>
              <div class="flex justify-between text-sm text-tertiary">
                <span>Positive: 30%</span>
                <span>Neutral: 20%</span>
                <span>Negative: 50%</span>
              </div>
            </div>

            <!-- Design Feature -->
            <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
              <div class="flex justify-between items-center mb-3">
                <h3 class="text-lg font-semibold text-secondary">🎨 Design</h3>
                <span
                  class="px-3 py-1 bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200 rounded-full text-sm font-semibold"
                >
                  Positive
                </span>
              </div>
              <div class="flex space-x-2 mb-2">
                <div class="flex-1 bg-green-500 h-3 rounded" style="width: 85%"></div>
                <div class="flex-1 bg-gray-300 h-3 rounded" style="width: 10%"></div>
                <div class="flex-1 bg-red-500 h-3 rounded" style="width: 5%"></div>
              </div>
              <div class="flex justify-between text-sm text-tertiary">
                <span>Positive: 85%</span>
                <span>Neutral: 10%</span>
                <span>Negative: 5%</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Feature Summary Statistics -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
          <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6">
            <h3 class="text-lg font-semibold text-secondary mb-3">Features Detected</h3>
            <div class="text-4xl font-bold text-primary mb-2">6</div>
            <p class="text-sm text-tertiary">Smartphone features identified</p>
          </div>

          <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6">
            <h3 class="text-lg font-semibold text-secondary mb-3">Sentences Analyzed</h3>
            <div class="text-4xl font-bold text-secondary mb-2">47</div>
            <p class="text-sm text-tertiary">Total sentences processed</p>
          </div>

          <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6">
            <h3 class="text-lg font-semibold text-secondary mb-3">Confidence</h3>
            <div class="text-4xl font-bold text-primary mb-2">89%</div>
            <p class="text-sm text-tertiary">Average model confidence</p>
          </div>
        </div>

        <!-- Extracted Keywords -->
        <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6 mb-6">
          <h2 class="text-xl font-semibold text-tertiary mb-4">Extracted Keywords</h2>
          <div class="flex flex-wrap gap-2">
            <span class="px-3 py-1 bg-fourthary text-white rounded-full text-sm"
              >camera quality</span
            >
            <span class="px-3 py-1 bg-fourthary text-white rounded-full text-sm">battery life</span>
            <span class="px-3 py-1 bg-fourthary text-white rounded-full text-sm"
              >display brightness</span
            >
            <span class="px-3 py-1 bg-fourthary text-white rounded-full text-sm"
              >fast charging</span
            >
            <span class="px-3 py-1 bg-fourthary text-white rounded-full text-sm">screen size</span>
            <span class="px-3 py-1 bg-fourthary text-white rounded-full text-sm">expensive</span>
            <span class="px-3 py-1 bg-fourthary text-white rounded-full text-sm"
              >premium design</span
            >
            <span class="px-3 py-1 bg-fourthary text-white rounded-full text-sm"
              >smooth performance</span
            >
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="flex justify-center space-x-4">
          <button
            @click="() => router.push('/list')"
            class="px-6 py-3 bg-tertiary text-white rounded-lg font-semibold hover:bg-opacity-90 transition"
          >
            Back to List
          </button>
          <button
            @click="exportResults"
            class="px-6 py-3 bg-fourthary text-white rounded-lg font-semibold hover:bg-opacity-90 transition"
          >
            Export Results
          </button>
          <button
            @click="analyzeNew"
            class="px-6 py-3 bg-accent text-white rounded-lg font-semibold hover:bg-opacity-90 transition"
          >
            Analyze New Video
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import SentimentChart from '@/components/SentimentChart.vue';
import { useRouter } from 'vue-router';
import NavbarComponent from '@/components/NavbarComponent.vue';
import Sidebar from '@/components/Sidebar.vue';
import { ref } from 'vue';

const router = useRouter();

const isSidebarOpen = ref(false);
</script>
