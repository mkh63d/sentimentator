<template>
  <div class="min-h-screen bg-secondary">
    <NavbarComponent />

    <div class="flex min-h-screen">
      <SideBar v-model="sidebarStore.visible">
        <ul class="space-y-2">
          <li>
            <button
              @click="router.push({ name: 'home' })"
              class="cursor-pointer text-background font-medium hover:text-indigo-100 transition text-left"
            >
              Home
            </button>
          </li>
          <li>
            <button
              @click="router.push({ name: 'list' })"
              class="cursor-pointer text-background font-medium hover:text-indigo-100 transition text-left"
            >
              Analysis Library
            </button>
          </li>
          <li class="self-bottom">
            <button
              @click="router.push({ name: 'settings' })"
              class="cursor-pointer text-background font-medium hover:text-indigo-100 transition text-left"
            >
              Settings
            </button>
          </li>
        </ul>
      </SideBar>
      <div class="max-w-7xl mx-auto px-6 py-8">
        <!-- Loading State -->
        <div v-if="loading" class="flex justify-center items-center min-h-screen">
          <div class="text-2xl text-tertiary">Loading analysis...</div>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="bg-red-100 dark:bg-red-900 p-6 rounded-lg">
          <p class="text-red-800 dark:text-red-200">{{ error }}</p>
        </div>

        <!-- Data Display -->
        <div v-else-if="analysisData">
          <!-- Page Header -->
          <div class="mb-8">
            <h1 class="m-0 text-3xl font-extrabold text-tertiary">Sentiment Analysis Results</h1>
            <a :href="analysisData.title" class="text-tertiary/90 hover:font-semibold text-xl mt-1">
              {{ analysisData.title || 'Smartphone review sentiment breakdown by features' }}
            </a>
          </div>

          <!-- Overall Sentiment Summary -->
          <div class="bg-fourth rounded-xl shadow-lg p-6 mb-6">
            <h2 class="text-2xl font-semibold text-[#bdccdb] mb-4">Overall Sentiment</h2>
            <div
              class="flex flex-col md:flex-row items-center justify-center md:space-x-8 space-y-6 md:space-y-0"
            >
              <div class="text-center">
                <div class="text-6xl font-bold text-primary mb-2">{{ dominantSentiment }}</div>
                <div class="text-lg text-[#bdccdb]">Dominant Sentiment</div>
              </div>
              <div class="w-full md:w-auto">
                <SentimentChart
                  :positive="analysisData.overallSentiment.positive"
                  :neutral="analysisData.overallSentiment.neutral"
                  :negative="analysisData.overallSentiment.negative"
                />
              </div>
              <div class="flex md:flex-col space-x-4 md:space-x-0 md:space-y-4">
                <div class="text-center p-4 bg-green-100 dark:bg-green-900 rounded-lg">
                  <div class="text-2xl font-bold text-green-600 mb-1">
                    {{ analysisData.overallSentiment.positive }}%
                  </div>
                  <div class="text-sm text-green-600">Positive</div>
                </div>
                <div class="text-center p-4 bg-gray-100 dark:bg-gray-700 rounded-lg">
                  <div class="text-2xl font-bold text-gray-400 mb-1">
                    {{ analysisData.overallSentiment.neutral }}%
                  </div>
                  <div class="text-sm text-gray-400">Neutral</div>
                </div>
                <div class="text-center p-4 bg-red-100 dark:bg-red-950 rounded-lg">
                  <div class="text-2xl font-bold text-red-600 mb-1">
                    {{ analysisData.overallSentiment.negative }}%
                  </div>
                  <div class="text-sm text-red-600">Negative</div>
                </div>
              </div>
            </div>
          </div>

          <!-- Feature-Based Sentiment Analysis -->
          <div class="bg-fourth rounded-xl shadow-lg p-6 mb-6">
            <h2 class="text-xl font-semibold text-[#bdccdb] mb-4">Sentiment by Feature</h2>
            <div class="space-y-4">
              <div
                v-for="feature in analysisData.features"
                :key="feature.name"
                class="border border-gray-200 dark:border-gray-700 rounded-lg p-4"
              >
                <div class="flex justify-between items-center mb-3">
                  <h3 class="text-lg font-semibold text-[#bdccdb]">{{ feature.name }}</h3>
                  <span
                    :class="getSentimentBadgeClass(feature.sentiment)"
                    class="px-3 py-1 rounded-full text-sm font-semibold"
                  >
                    {{ feature.sentiment }}
                  </span>
                </div>
                <div class="flex space-x-2 mb-2">
                  <div
                    class="bg-green-500 h-3 rounded"
                    :style="{ width: `${feature.positive}%` }"
                  ></div>
                  <div
                    class="bg-gray-300 h-3 rounded"
                    :style="{ width: `${feature.neutral}%` }"
                  ></div>
                  <div
                    class="bg-red-500 h-3 rounded"
                    :style="{ width: `${feature.negative}%` }"
                  ></div>
                </div>
                <div class="flex justify-between text-sm text-[#bdccdb] opacity-85">
                  <span>Positive: {{ feature.positive }}%</span>
                  <span>Neutral: {{ feature.neutral }}%</span>
                  <span>Negative: {{ feature.negative }}%</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Feature Summary Statistics -->
          <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
            <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6">
              <h3 class="text-lg font-semibold text-[#bdccdb] mb-3">Features Detected</h3>
              <div class="text-4xl font-bold text-primary mb-2">
                {{ analysisData.features.length }}
              </div>
              <p class="text-sm text-[#bdccdb]/80">Features identified</p>
            </div>

            <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6">
              <h3 class="text-lg font-semibold text-[#bdccdb] mb-3">Sentences Analyzed</h3>
              <div class="text-4xl font-bold text-[#bdccdb] mb-2">
                {{ analysisData.statistics.sentencesAnalyzed }}
              </div>
              <p class="text-sm text-[#bdccdb]/80">Total sentences processed</p>
            </div>

            <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6">
              <h3 class="text-lg font-semibold text-[#bdccdb] mb-3">Confidence</h3>
              <div class="text-4xl font-bold text-primary mb-2">
                {{ analysisData.statistics.confidence }}%
              </div>
              <p class="text-sm text-[#bdccdb]/80">Average model confidence</p>
            </div>
          </div>

          <!-- Extracted Keywords -->
          <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6 mb-6">
            <h2 class="text-xl font-semibold text-[#bdccdb] mb-4">Extracted Keywords</h2>
            <div class="flex flex-wrap gap-2">
              <span
                v-for="keyword in analysisData.keywords"
                :key="keyword"
                class="px-3 py-1 bg-fourth text-white rounded-full text-sm"
              >
                {{ keyword }}
              </span>
            </div>
          </div>

          <!-- Action Buttons -->
          <div class="flex justify-center space-x-4">
            <button
              @click="() => router.push({ name: 'list' })"
              class="cursor-pointer px-6 py-3 bg-tertiary text-background rounded-lg font-semibold hover:bg-opacity-90 transition"
            >
              Back to List
            </button>
            <button
              @click="exportResults"
              class="cursor-pointer px-6 py-3 bg-fourth text-white rounded-lg font-semibold hover:bg-opacity-90 transition"
            >
              Export Results
            </button>
            <button
              @click="searchStore.show()"
              class="cursor-pointer px-6 py-3 bg-accent text-background rounded-lg font-semibold hover:bg-opacity-90 transition"
            >
              Analyze New Video
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
  <SearchModal v-model="searchStore.visible" />
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue';
import SentimentChart from '@/components/SentimentChart.vue';
import { useRouter, useRoute } from 'vue-router';
import NavbarComponent from '@/components/NavbarComponent.vue';
import SideBar from '@/components/SideBar.vue';
import { useSidebarStore } from '@/stores/sidebarStore';
import { useSearchStore } from '@/stores/searchStore';
import { getAnalysisById } from '@/api/analyses';
import SearchModal from '@/components/SearchModal.vue';

// Types
interface Feature {
  name: string;
  sentiment: 'Positive' | 'Neutral' | 'Negative';
  positive: number;
  neutral: number;
  negative: number;
}

interface AnalysisData {
  id: string;
  title?: string;
  overallSentiment: {
    positive: number;
    neutral: number;
    negative: number;
  };
  transcription: string;
  features: Feature[];
  statistics: {
    sentencesAnalyzed: number;
    confidence: number;
  };
  keywords: string[];
}

const router = useRouter();
const route = useRoute();
const searchStore = useSearchStore();
const sidebarStore = useSidebarStore();

// State
const analysisData = ref<AnalysisData | null>(null);
const loading = ref(true);
const error = ref<string | null>(null);

// Computed
const dominantSentiment = computed(() => {
  if (!analysisData.value) return 'Unknown';
  const { positive, neutral, negative } = analysisData.value.overallSentiment;
  const max = Math.max(positive, neutral, negative);
  if (max === positive) return 'Positive';
  if (max === neutral) return 'Neutral';
  return 'Negative';
});

// Methods
function getSentimentBadgeClass(sentiment: string): string {
  switch (sentiment) {
    case 'Positive':
      return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200';
    case 'Negative':
      return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200';
    default:
      return 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200';
  }
}

/**
 * Generates sentiment values that sum to 100 based on the primary sentiment. Fully randomized.
 * @deprecated Use actual sentiment analysis data instead.
 * @param sentiment
 */
function generateSentimentValues(sentiment: string): {
  positive: number;
  neutral: number;
  negative: number;
} {
  let primary: number;
  let secondary: number;
  let tertiary: number;

  if (sentiment === 'positive') {
    primary = Math.floor(Math.random() * 81) + 15; // 15-95 (leaving room for min 5 for others)
    const remaining = 100 - primary;
    secondary = Math.floor(Math.random() * (remaining - 10 + 1)) + 5; // At least 5
    tertiary = 100 - primary - secondary;
    return { positive: primary, neutral: secondary, negative: tertiary };
  } else if (sentiment === 'negative') {
    primary = Math.floor(Math.random() * 81) + 15; // 15-95
    const remaining = 100 - primary;
    secondary = Math.floor(Math.random() * (remaining - 10 + 1)) + 5; // At least 5
    tertiary = 100 - primary - secondary;
    return { positive: tertiary, neutral: secondary, negative: primary };
  } else {
    // neutral
    primary = Math.floor(Math.random() * 81) + 15; // 15-95
    const remaining = 100 - primary;
    secondary = Math.floor(Math.random() * (remaining - 10 + 1)) + 5; // At least 5
    tertiary = 100 - primary - secondary;
    return { positive: secondary, neutral: primary, negative: tertiary };
  }
}

function calculateOverall(features: Feature[]): {
  positive: number;
  neutral: number;
  negative: number;
} {
  // Calculate overall sentiment based on feature sentiments
  const sentimentCounts = { positive: 0, neutral: 0, negative: 0 };
  features.forEach((feature) => {
    const sentiment = feature.sentiment.toLowerCase();
    if (sentiment === 'positive') sentimentCounts.positive++;
    else if (sentiment === 'negative') sentimentCounts.negative++;
    else sentimentCounts.neutral++;
  });

  const totalFeatures = features.length || 1; // Avoid division by zero
  const overallSentiment = {
    positive: Math.round((sentimentCounts.positive / totalFeatures) * 100),
    neutral: Math.round((sentimentCounts.neutral / totalFeatures) * 100),
    negative: Math.round((sentimentCounts.negative / totalFeatures) * 100),
  };

  return overallSentiment;
}

//TODO: Replace with actual API call for individual sentiment scores fetching
async function fetchAnalysisData(id: string) {
  loading.value = true;
  error.value = null;

  try {
    const response = await getAnalysisById(id);
    const parsedSummary = response.summary as any;
    console.log('Parsed Summary:', parsedSummary);
    const features = Object.entries(parsedSummary)
      .filter(([key]) => key !== 'statistics')
      .map(([key, value]: [string, any]) => {
        const sentiment = typeof value === 'string' ? value.toLowerCase() : 'neutral';
        const sentimentValues = generateSentimentValues(sentiment);
        return {
          name: key.replace(/_/g, ' ').replace(/\b\w/g, (c) => c.toUpperCase()),
          sentiment: value.charAt(0).toUpperCase() + value.slice(1),
          ...sentimentValues,
        };
      });

    //TODO: add statistics fetching
    analysisData.value = {
      id: response.id || 'mock-id',
      title: response.url,
      overallSentiment: {
        ...calculateOverall(features),
      },
      transcription: response.transcription || '',
      features: features,
      statistics: parsedSummary.statistics || { sentencesAnalyzed: 156, confidence: 90 },
      keywords: Object.keys(parsedSummary).filter((key) => key !== 'statistics'),
    };
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'An error occurred';
    console.error('Error fetching analysis:', err);
  } finally {
    loading.value = false;
  }
}

function exportResults() {
  if (!analysisData.value) return;

  const data = analysisData.value;
  const pdfContent = `
# Sentiment Analysis Results
${data.title ? `## ${data.title}` : ''}

### Overall Sentiment
- **Dominant Sentiment:** ${dominantSentiment.value}
- Positive: ${data.overallSentiment.positive}%
- Neutral: ${data.overallSentiment.neutral}%
- Negative: ${data.overallSentiment.negative}%

### Sentiment by Feature

${data.features
  .map(
    (f) => `#### ${f.name} - ${f.sentiment}
- Positive: ${f.positive}% | Neutral: ${f.neutral}% | Negative: ${f.negative}%`,
  )
  .join('\n\n')}

### Statistics
- Features Detected: ${data.features.length}
- Sentences Analyzed: ${data.statistics.sentencesAnalyzed}
- Average Confidence: ${data.statistics.confidence}%

### Extracted Keywords
${data.keywords.join(', ')}
  `.trim();

  const blob = new Blob([pdfContent], { type: 'text/plain' });
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = `sentiment-analysis-${data.id}-${new Date().getTime()}.txt`;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  window.URL.revokeObjectURL(url);
}

// Lifecycle
onMounted(() => {
  const analysisId = route.params.id as string;
  if (analysisId) {
    fetchAnalysisData(analysisId);
  } else {
    error.value = 'No analysis ID provided';
    loading.value = false;
  }
});
</script>
