<template>
  <div class="min-h-screen bg-secondary overflow-hidden">
    <NavbarComponent
      :isSidebarOpen="isSidebarOpen"
      @update:isSidebarOpen="isSidebarOpen = $event"
    />
    <div class="flex min-h-screen">
      <Sidebar v-model="isSidebarOpen">
        <ul class="space-y-2">
          <li>
            <button
              @click="router.push({ name: 'home' })"
              class="font-medium text-left transition cursor-pointer text-background hover:text-indigo-100"
            >
              {{ t('home.title') }}
            </button>
          </li>
          <li>
            <button
              @click="router.push({ name: 'list' })"
              class="font-medium text-left transition cursor-pointer text-background hover:text-indigo-100"
            >
              {{ t('analysis.library') }}
            </button>
          </li>
          <li class="self-bottom">
            <button
              @click="router.push({ name: 'settings' })"
              class="font-medium text-left transition cursor-pointer text-background hover:text-indigo-100"
            >
              {{ t('settings.title') }}
            </button>
          </li>
        </ul>
      </Sidebar>

      <div class="px-4 sm:px-8 py-8 max-w-[1400px] mx-auto flex-1">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between mb-6 gap-3 sm:gap-0">
          <h1 class="m-0 text-2xl sm:text-3xl font-extrabold text-tertiary min-w-0 truncate">
            {{ t('analysis.library') }}
          </h1>
          <div class="flex gap-2 sm:gap-3 items-center">
            <button
              @click="fetchAnalyses"
              class="flex items-center gap-2 px-2 py-2 sm:px-4 sm:py-3 text-sm sm:text-base font-semibold transition rounded-lg cursor-pointer bg-primary text-background hover:opacity-80"
              :title="t('analysis.refresh')"
            >
              <RotateCcw class="w-4 h-4 sm:w-5 sm:h-5" />
              <span class="hidden sm:inline">{{ t('analysis.refresh') }}</span>
            </button>
            <button
              @click="createNew"
              class="flex items-center gap-2 px-2 py-2 sm:px-4 sm:py-3 text-sm sm:text-base font-semibold transition rounded-lg cursor-pointer bg-primary text-background hover:opacity-80"
              :title="t('analysis.newAnalysis')"
            >
              <Plus class="w-4 h-4 sm:w-5 sm:h-5" />
              <span class="hidden sm:inline">{{ t('analysis.newAnalysis') }}</span>
            </button>
            <button
              @click="mobileFilterWindow = true"
              class="flex sm:hidden items-center gap-2 px-2 py-2 sm:px-4 sm:py-3 text-sm sm:text-base font-semibold transition rounded-lg cursor-pointer bg-primary text-background hover:opacity-80"
              :title="t('analysis.filterAnalyses')"
            >
              <Filter class="w-4 h-4 sm:w-5 sm:h-5" />
            </button>
          </div>
        </div>

        <div class="gap-3 mb-6 hidden sm:flex">
          <input
            v-model="searchQuery"
            type="text"
            :placeholder="t('analysis.searchAnalyses')"
            class="flex-1 px-3 py-2 text-sm sm:text-base border rounded-lg border-background"
          />
          <select
            v-model="filterStatus"
            class="px-3 py-3 text-base border rounded-lg cursor-pointer border-background"
          >
            <option value="">{{ t('analysis.allStatuses') }}</option>
            <option value="completed">{{ t('analysis.completed') }}</option>
            <option value="processing">{{ t('analysis.processing') }}</option>
            <option value="failed">{{ t('analysis.failed') }}</option>
          </select>
        </div>

        <div v-if="loading" class="py-12 text-center">
          <span
            class="inline-block w-5 h-5 border-4 border-gray-200 rounded-full border-t-blue-500 animate-spin"
          ></span>
          {{ t('analysis.loadingLibrary') }}
        </div>

        <div v-else-if="error" class="py-12 text-center">
          <p>
            {{ t('analysis.errorFetching') }} <br /><code>{{ error }}</code>
          </p>
          <button
            @click="fetchAnalyses"
            class="px-4 py-2 transition bg-gray-200 rounded hover:bg-gray-300"
          >
            {{ t('analysis.retry') }}
          </button>
        </div>

        <div v-else-if="filteredAnalyses.length === 0" class="py-12 text-center">
          <p>{{ t('analysis.noAnalyses') }}</p>
          <button
            @click="createNew"
            class="px-4 py-2 transition bg-gray-100 rounded hover:bg-gray-200 cursor-pointer"
          >
            {{ t('analysis.createFirst') }}
          </button>
        </div>

        <!-- Redesigned analysis cards -->
        <div class="space-y-4">
          <div
            v-for="analysis in filteredAnalyses"
            :key="analysis.id"
            class="border rounded-xl p-4 sm:p-5 flex flex-col gap-4 shadow-sm hover:shadow-md transition-all duration-200"
            :class="[
              analysis.summary
                ? 'bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700 cursor-pointer hover:border-primary'
                : 'bg-gray-50 dark:bg-gray-800/50 border-gray-300 dark:border-gray-600 cursor-default opacity-90',
            ]"
            @click.stop="analysis.summary !== null && analysis.id && openAnalysis(analysis.id)"
          >
            <div class="flex flex-col sm:flex-row gap-4">
              <!-- Thumbnail Section with Platform Border -->
              <div class="relative shrink-0">
                <div
                  class="w-full sm:w-48 h-full min-h-[120px] bg-gray-200 dark:bg-gray-700 rounded-lg overflow-hidden flex items-center justify-center group border-4 transition-all"
                  :class="[
                    analysis.summary !== null ? 'cursor-pointer' : 'opacity-60',
                    getPlatformBorderClass(analysis.platform),
                  ]"
                  @click.stop="
                    analysis.summary !== null && analysis.id && openAnalysis(analysis.id)
                  "
                  role="button"
                  :aria-disabled="analysis.summary === null"
                  :title="
                    analysis.summary !== null
                      ? t('analysis.openAnalysis')
                      : t('analysis.summaryNotReady')
                  "
                >
                  <!-- Thumbnail image -->
                  <img
                    v-if="analysis.thumbnail_url"
                    :src="analysis.thumbnail_url"
                    alt="analysis thumbnail"
                    class="w-full h-full object-cover transition-transform duration-200 group-hover:scale-105"
                    loading="lazy"
                  />

                  <!-- Fallback placeholder -->
                  <template v-else>
                    <div
                      class="flex flex-col items-center justify-center text-gray-400 dark:text-gray-500 px-2 text-center"
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="w-8 h-8 mb-1"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="1.5"
                          d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M4 6h11a2 2 0 012 2v8a2 2 0 01-2 2H4a2 2 0 01-2-2V8a2 2 0 012-2z"
                        />
                      </svg>
                      <div class="text-xs">{{ t('analysis.noThumbnail') }}</div>
                    </div>
                  </template>
                </div>
                
                <!-- Platform indicator pill -->
                <div
                  v-if="analysis.platform"
                  class="absolute -bottom-2 left-1/2 -translate-x-1/2 px-3 py-1 rounded-full text-xs font-bold shadow-lg"
                  :class="getPlatformPillClass(analysis.platform)"
                  :title="analysis.platform"
                >
                  {{ getPlatformName(analysis.platform) }}
                </div>
              </div>

              <!-- Content Section -->
              <div class="flex-1 min-w-0 flex flex-col justify-between">
                <div>
                  <div class="flex items-start justify-between gap-3 mb-2">
                    <h3
                      class="text-base sm:text-lg font-semibold text-gray-900 dark:text-gray-100 line-clamp-2"
                      :title="analysis.title ?? analysis.url"
                    >
                      {{ analysis.title && analysis.title.trim() ? analysis.title : analysis.url }}
                    </h3>
                    
                    <!-- Status badge -->
                    <span
                      class="px-3 py-1 rounded-full text-xs font-semibold whitespace-nowrap shrink-0"
                      :class="
                        analysis.summary !== null
                          ? 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300'
                          : 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900 dark:text-yellow-300'
                      "
                    >
                      {{
                        analysis.summary === null
                          ? t('analysis.processing')
                          : t('analysis.completed')
                      }}
                    </span>
                  </div>
                  
                  <!-- Metadata -->
                  <div class="text-xs text-gray-600 dark:text-gray-400 space-y-1">
                    <div class="flex items-center gap-2">
                      <span class="font-medium">{{ t('analysis.created') }}</span>
                      <span>{{ analysis.created_at ? formatDate(analysis.created_at) : '' }}</span>
                    </div>
                    <div class="flex items-center gap-2">
                      <span class="font-medium">{{ t('analysis.updated') }}</span>
                      <span>{{ analysis.updated_at ? formatDate(analysis.updated_at) : '' }}</span>
                    </div>
                  </div>
                </div>

                <!-- Action buttons -->
                <div class="flex items-center justify-between gap-2 mt-3">
                  <button
                    @click.stop="analysis.transcription && analysis.id && openAnalysis(analysis.id)"
                    :disabled="analysis.summary === null"
                    :class="[
                      'flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all shadow-sm',
                      analysis.summary === null
                        ? 'bg-gray-200 text-gray-400 cursor-not-allowed dark:bg-gray-700 dark:text-gray-600'
                        : 'bg-gradient-to-r from-blue-500 to-indigo-600 text-white hover:from-blue-600 hover:to-indigo-700 hover:shadow-md cursor-pointer',
                    ]"
                    :title="t('analysis.view')"
                    :aria-label="t('analysis.view')"
                  >
                    <Eye class="w-4 h-4" />
                    <span class="hidden sm:inline">{{ t('analysis.view') }}</span>
                  </button>

                  <button
                    @click.stop="analysis.id && eraseAnalysis(analysis.id)"
                    class="flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all shadow-sm bg-gray-100 text-gray-700 hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-300 dark:hover:bg-gray-600 cursor-pointer"
                    :title="t('analysis.delete')"
                    :aria-label="t('analysis.delete')"
                  >
                    <Trash2 class="w-4 h-4" />
                    <span class="hidden sm:inline">{{ t('analysis.delete') }}</span>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <SearchModal v-model="searchStore.visible" @analyze="fetchAnalyses" />
    <MobileFilterModal
      v-model="mobileFilterWindow"
      v-model:searchQuery="searchQuery"
      v-model:filterStatus="filterStatus"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { Plus, Eye, Trash2, RotateCcw, Filter } from 'lucide-vue-next';
import Swal from 'sweetalert2';
import Sidebar from '@/components/SideBar.vue';
import NavbarComponent from '@/components/NavbarComponent.vue';
import SearchModal from '@/components/SearchModal.vue';
import MobileFilterModal from '@/components/MobileFilterModal.vue';
import { useSearchStore } from '@/stores/searchStore';
import { getAnalyses, deleteAnalysis } from '@/api/analyses';
import { components } from '@/generated/schema';

interface Sentiment {
  positive: number;
  negative: number;
  neutral: number;
}

interface Analysis {
  id: string;
  title: string;
  description?: string;
  createdAt: string;
  status: 'completed' | 'processing' | 'failed';
  sentiment?: Sentiment;
  transcription?: string;
}

const router = useRouter();
const { t } = useI18n();
const analyses = ref<Analysis[]>([]);
const realAnalyses = ref<components['schemas']['Analysis'][]>([]);
const loading = ref(true);
const error = ref<string | null>(null);
const searchQuery = ref('');
const filterStatus = ref('');
const mobileFilterWindow = ref(false);

const filteredAnalyses = computed(() => {
  const query = searchQuery.value.trim().toLowerCase();

  function inferStatus(
    analysis: components['schemas']['Analysis'],
  ): 'completed' | 'processing' | 'failed' {
    const s = analysis?.summary;
    if (s === null || typeof s === 'undefined') return 'processing';
    // treat an explicit empty string or empty object as failed
    if (s === '' || (typeof s === 'object' && Object.keys(s).length === 0)) return 'failed';
    return 'completed';
  }

  return realAnalyses.value.filter((analysis) => {
    const title = String(analysis.title ?? '').toLowerCase();
    const matchesSearch = !query || title.includes(query);

    const status = inferStatus(analysis);
    const matchesStatus = !filterStatus.value || filterStatus.value === status;

    return matchesSearch && matchesStatus;
  });
});

const fetchAnalyses = async () => {
  try {
    loading.value = true;
    error.value = null;
    const response = await getAnalyses();
    console.log('Fetched analyses response:', response);

    realAnalyses.value = response;
    console.log('Real analyses set:', realAnalyses.value);
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'An error occurred';
  } finally {
    loading.value = false;
  }
};

const openAnalysis = (id: string) => {
  router.push({ name: 'analysis', params: { id } });
};

const searchStore = useSearchStore();

const createNew = () => {
  searchStore.show();
};

const eraseAnalysis = async (id: string) => {
  // Use SweetAlert2 for confirmation
  // Make sure to install sweetalert2 and import it
  await Swal.fire({
    title: t('analysis.deleteConfirmTitle'),
    text: t('analysis.deleteConfirmText'),
    icon: 'warning',
    showCancelButton: true,
    confirmButtonColor: '#d33',
    cancelButtonColor: '#3085d6',
    confirmButtonText: t('analysis.deleteConfirmButton'),
    cancelButtonText: t('common.cancel'),
  }).then(async (result) => {
    if (result.isConfirmed) {
      try {
        await deleteAnalysis(id);
        analyses.value = analyses.value.filter((a) => a.id !== id);
        Swal.fire({
          title: t('analysis.deletedTitle'),
          text: t('analysis.deletedText'),
          icon: 'success',
          timer: 2000,
          showConfirmButton: false,
        });
      } catch (e) {
        Swal.fire({
          title: t('analysis.deleteErrorTitle'),
          text: t('analysis.deleteErrorText') + ' ' + (e instanceof Error ? e.message : ''),
          icon: 'error',
          timer: 2000,
          showConfirmButton: false,
        });
      }
    }
  });
};

const formatDate = (dateString: string): string => {
  const locale = navigator.language || navigator.languages?.[0] || 'en-US';
  return new Date(dateString).toLocaleDateString(locale, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
};

const getPlatformBorderClass = (platform?: string): string => {
  if (!platform) return 'border-gray-300 dark:border-gray-600';
  
  const platformMap: Record<string, string> = {
    youtube_platform: 'border-red-500 dark:border-red-400',
    tiktok_platform: 'border-black dark:border-white',
    instagram_platform: 'border-pink-500 dark:border-pink-400',
  };
  
  return platformMap[platform] || 'border-gray-300 dark:border-gray-600';
};

const getPlatformPillClass = (platform?: string): string => {
  if (!platform) return 'bg-gray-500 text-white';
  
  const platformMap: Record<string, string> = {
    youtube_platform: 'bg-red-600 text-white',
    tiktok_platform: 'bg-black text-white dark:bg-white dark:text-black',
    instagram_platform: 'bg-gradient-to-r from-purple-500 via-pink-500 to-orange-500 text-white',
  };
  
  return platformMap[platform] || 'bg-gray-500 text-white';
};

const getPlatformName = (platform?: string): string => {
  if (!platform) return '';
  
  const platformMap: Record<string, string> = {
    youtube_platform: 'YOUTUBE',
    tiktok_platform: 'TIKTOK',
    instagram_platform: 'INSTAGRAM',
  };
  
  return platformMap[platform] || platform;
};

onMounted(() => {
  fetchAnalyses();
});

const isSidebarOpen = ref(false);
</script>
