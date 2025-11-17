<template>
  <div class="min-h-screen bg-secondary">
    <NavbarComponent />

    <div class="flex min-h-screen">
      <SideBar v-model="sidebarStore.visible">
        <ul class="space-y-2">
          <li>
            <button
              @click="router.push({ name: 'home' })"
              class="font-medium text-left transition cursor-pointer text-background hover:text-indigo-100"
            >
              Home
            </button>
          </li>
          <li>
            <button
              @click="router.push({ name: 'list' })"
              class="font-medium text-left transition cursor-pointer text-background hover:text-indigo-100"
            >
              Analysis Library
            </button>
          </li>
          <li class="self-bottom">
            <button
              @click="router.push({ name: 'settings' })"
              class="font-medium text-left transition cursor-pointer text-background hover:text-indigo-100"
            >
              Settings
            </button>
          </li>
        </ul>
      </SideBar>
      <div class="px-6 py-8 mx-auto max-w-7xl">
        <!-- Loading State -->
        <div v-if="loading" class="flex items-center justify-center min-h-screen">
          <div class="text-2xl text-tertiary">Loading analysis...</div>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="p-6 bg-red-100 rounded-lg dark:bg-red-900">
          <p class="text-red-800 dark:text-red-200">{{ error }}</p>
        </div>

        <!-- Data Display -->
        <div v-else-if="analysisData">
          <!-- Page Header -->
          <div class="mb-8">
            <h1 class="m-0 text-3xl font-extrabold text-tertiary">Sentiment Analysis Results</h1>
            <a :href="analysisData.title" class="mt-1 text-xl text-tertiary/90 hover:font-semibold">
              {{ analysisData.title || 'Smartphone review sentiment breakdown by features' }}
            </a>
          </div>

          <!-- Overall Sentiment Summary -->
          <div class="p-6 mb-6 shadow-lg bg-fourth rounded-xl">
            <h2 class="text-2xl font-semibold text-[#bdccdb] mb-4">Overall Sentiment</h2>
            <div
              class="flex flex-col items-center justify-center space-y-6 md:flex-row md:space-x-8 md:space-y-0"
            >
              <div class="text-center">
                <div class="mb-2 text-6xl font-bold text-primary">{{ dominantSentiment }}</div>
                <div class="text-lg text-[#bdccdb]">Dominant Sentiment</div>
              </div>
              <div class="w-full md:w-auto">
                <SentimentChart
                  :positive="analysisData.overallSentiment.positive"
                  :neutral="analysisData.overallSentiment.neutral"
                  :negative="analysisData.overallSentiment.negative"
                />
              </div>
              <div class="flex space-x-4 md:flex-col md:space-x-0 md:space-y-4">
                <div class="p-4 text-center bg-green-100 rounded-lg dark:bg-green-900">
                  <div class="mb-1 text-2xl font-bold text-green-600">
                    {{ analysisData.overallSentiment.positive }}%
                  </div>
                  <div class="text-sm text-green-600">Positive</div>
                </div>
                <div class="p-4 text-center bg-gray-100 rounded-lg dark:bg-gray-700">
                  <div class="mb-1 text-2xl font-bold text-gray-400">
                    {{ analysisData.overallSentiment.neutral }}%
                  </div>
                  <div class="text-sm text-gray-400">Neutral</div>
                </div>
                <div class="p-4 text-center bg-red-100 rounded-lg dark:bg-red-950">
                  <div class="mb-1 text-2xl font-bold text-red-600">
                    {{ analysisData.overallSentiment.negative }}%
                  </div>
                  <div class="text-sm text-red-600">Negative</div>
                </div>
              </div>
            </div>
          </div>

          <!-- Feature-Based Sentiment Analysis -->
          <div class="p-6 mb-6 shadow-lg bg-fourth rounded-xl">
            <h2 class="text-xl font-semibold text-[#bdccdb] mb-4">Sentiment by Feature</h2>
            <div class="space-y-4">
              <div
                v-for="feature in analysisData.features"
                :key="feature.name"
                class="p-4 border border-gray-200 rounded-lg dark:border-gray-700"
              >
                <div class="flex items-center justify-between mb-3">
                  <h3 class="text-lg font-semibold text-[#bdccdb]">{{ feature.name }}</h3>
                  <span
                    :class="getSentimentBadgeClass(feature.sentiment)"
                    class="px-3 py-1 text-sm font-semibold rounded-full"
                  >
                    {{ feature.sentiment }}
                  </span>
                </div>
                <div class="flex mb-2 space-x-2">
                  <div
                    class="h-3 bg-green-500 rounded"
                    :style="{ width: `${feature.positive}%` }"
                  ></div>
                  <div
                    class="h-3 bg-gray-300 rounded"
                    :style="{ width: `${feature.neutral}%` }"
                  ></div>
                  <div
                    class="h-3 bg-red-500 rounded"
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
          <div class="grid grid-cols-1 gap-6 mb-6 md:grid-cols-3">
            <div class="p-6 bg-white shadow-lg dark:bg-gray-800 rounded-xl">
              <h3 class="text-lg font-semibold text-[#bdccdb] mb-3">Features Detected</h3>
              <div class="mb-2 text-4xl font-bold text-primary">
                {{ analysisData.features.length }}
              </div>
              <p class="text-sm text-[#bdccdb]/80">Features identified</p>
            </div>

            <div class="p-6 bg-white shadow-lg dark:bg-gray-800 rounded-xl">
              <h3 class="text-lg font-semibold text-[#bdccdb] mb-3">Sentences Analyzed</h3>
              <div class="text-4xl font-bold text-[#bdccdb] mb-2">
                {{ analysisData.statistics.sentencesAnalyzed ?? '—' }}
              </div>
              <p class="text-sm text-[#bdccdb]/80">Total sentences processed</p>
            </div>

            <div class="p-6 bg-white shadow-lg dark:bg-gray-800 rounded-xl">
              <h3 class="text-lg font-semibold text-[#bdccdb] mb-3">Confidence</h3>
              <div class="mb-2 text-4xl font-bold text-primary">
                {{ analysisData.statistics.confidence }}%
              </div>
              <p class="text-sm text-[#bdccdb]/80">Average model confidence</p>
            </div>
          </div>

          <!-- Extracted Keywords -->
          <div class="p-6 mb-6 bg-white shadow-lg dark:bg-gray-800 rounded-xl">
            <h2 class="text-xl font-semibold text-[#bdccdb] mb-4">Extracted Keywords</h2>
            <div class="flex flex-wrap gap-2">
              <span
                v-for="keyword in analysisData.keywords"
                :key="keyword"
                class="px-3 py-1 text-sm text-white rounded-full bg-fourth"
              >
                {{ keyword }}
              </span>
            </div>
          </div>

          <!-- Action Buttons -->
          <div class="flex justify-center space-x-4">
            <button
              @click="() => router.push({ name: 'list' })"
              class="px-6 py-3 font-semibold transition rounded-lg cursor-pointer bg-tertiary text-background hover:bg-opacity-90"
            >
              Back to List
            </button>
            <button
              @click="exportResults"
              class="px-6 py-3 font-semibold text-white transition rounded-lg cursor-pointer bg-fourth hover:bg-opacity-90"
            >
              Export Results
            </button>
            <button
              @click="searchStore.show()"
              class="px-6 py-3 font-semibold transition rounded-lg cursor-pointer bg-accent text-background hover:bg-opacity-90"
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
  score: number;
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
    sentencesAnalyzed: number | null;
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

function scoreToSentiment(score: number): Feature['sentiment'] {
  if (score >= 0.6) return 'Positive';
  if (score <= 0.4) return 'Negative';
  return 'Neutral';
}

function toPercent(score: number): number {
  return Math.round(Math.min(Math.max(score, 0), 1) * 100);
}

function formatFeatureName(name: string): string {
  return name.replace(/_/g, ' ').replace(/\b\w/g, (char) => char.toUpperCase());
}

function normalizeSummary(rawSummary: unknown): Record<string, unknown> {
  if (!rawSummary) return {};
  if (typeof rawSummary === 'string') {
    try {
      return JSON.parse(rawSummary);
    } catch (parseError) {
      console.warn('Unable to parse summary JSON', parseError);
      return {};
    }
  }
  if (typeof rawSummary === 'object') {
    return rawSummary as Record<string, unknown>;
  }
  return {};
}

async function fetchAnalysisData(id: string) {
  loading.value = true;
  error.value = null;

  try {
    const response = await getAnalysisById(id);
    const rawSummary = (response as any)?.summary ?? response;
    const parsedSummary = normalizeSummary(rawSummary) as {
      overall?: number;
      features?: Record<string, number>;
      confidence?: number;
      statistics?: { sentencesAnalyzed?: number };
      sentencesAnalyzed?: number;
    };

    const summaryFeatures = parsedSummary.features ?? {};
    const featureEntries = Object.entries(summaryFeatures) as [string, number][];
    const features: Feature[] = featureEntries.map(([name, value]) => {
      const score = typeof value === 'number' ? value : Number(value) || 0;
      const positive = toPercent(score);
      const negative = Math.max(0, 100 - positive);
      return {
        name: formatFeatureName(name),
        sentiment: scoreToSentiment(score),
        positive,
        neutral: 0,
        negative,
        score,
      };
    });

    const overallPositive = toPercent(parsedSummary.overall ?? 0);
    const overallSentiment = {
      positive: overallPositive,
      neutral: 0,
      negative: Math.max(0, 100 - overallPositive),
    };

    analysisData.value = {
      id: response.id || 'mock-id',
      title: response.url,
      overallSentiment,
      transcription: response.transcription || '',
      features,
      statistics: {
        sentencesAnalyzed:
          parsedSummary?.statistics?.sentencesAnalyzed ?? parsedSummary?.sentencesAnalyzed ?? null,
        confidence: toPercent(parsedSummary.confidence ?? 0),
      },
      keywords: features.map((feature) => feature.name),
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
- Sentences Analyzed: ${data.statistics.sentencesAnalyzed ?? 'N/A'}
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
