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
            class="px-4 py-2 transition bg-gray-100 rounded hover:bg-gray-200"
          >
            {{ t('analysis.createFirst') }}
          </button>
        </div>

        <!-- Use the mobile stacked cards layout on all screen sizes -->
        <div class="space-y-3">
          <div
            v-for="analysis in filteredAnalyses"
            :key="analysis.id"
            class="border rounded-lg p-3 flex flex-col gap-2 bg-background"
            :class="[
              analysis.summary
                ? 'bg-green-50 dark:bg-green-700 border-green-200 cursor-pointer'
                : 'bg-yellow-50 dark:bg-yellow-400 border-yellow-200 cursor-default',
              analysis && analysis.transcription && analysis.id ? 'hover:opacity-50' : '',
            ]"
            @click.stop="analysis.summary !== null && analysis.id && openAnalysis(analysis.id)"
          >
            <div class="sm:flex sm:flex-row flex-col justify-between items-start gap-3">
              <div class="flex items-center gap-3">
                <!-- Thumbnail / video preview area -->
                <div
                  class="w-full aspect-video sm:w-28 sm:h-20 bg-gray-100 rounded overflow-hidden shrink-0 flex items-center justify-center"
                  :class="
                    analysis.summary !== null ? 'cursor-pointer hover:opacity-90' : 'opacity-60'
                  "
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
                  <!-- image if provided by API -->
                  <img
                    v-if="analysis.thumbnail_url"
                    :src="analysis.thumbnail_url"
                    alt="analysis thumbnail"
                    class="w-full h-full object-cover"
                    loading="lazy"
                  />

                  <!-- fallback placeholder -->
                  <template v-else>
                    <div
                      class="flex flex-col items-center justify-center text-xs text-tertiary/60 px-2 text-center"
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="w-6 h-6 mb-1"
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
              </div>

              <div class="hidden sm:flex items-center gap-3">
                <!--Display the logo from which the analysis originated-->
                <div
                  class="w-28 h-20 bg-gray-100 rounded overflow-hidden shrink-0 flex items-center justify-center"
                >
                  <img
                    v-if="analysis.platform"
                    :src="LogoSource[analysis.platform]"
                    :alt="analysis.platform"
                    class="w-full h-full object-cover"
                    loading="lazy"
                  />
                </div>
              </div>

              <div class="min-w-0 flex-1">
                <div
                  class="sm:p-0 pt-2 text-xl sm:text-sm font-semibold text-[#bdccdb] truncate max-w-full sm:max-w-full"
                  :title="analysis.title ?? analysis.url"
                >
                  {{ analysis.title && analysis.title.trim() ? analysis.title : analysis.url }}
                </div>
                <div class="text-xs text-[#bdccdb] mt-1">
                  <div>
                    {{ t('analysis.created') }}
                    {{ analysis.created_at ? formatDate(analysis.created_at) : '' }}
                  </div>
                  <div>
                    {{ t('analysis.updated') }}
                    {{ analysis.updated_at ? formatDate(analysis.updated_at) : '' }}
                  </div>
                </div>
              </div>

              <div class="flex flex-col items-end gap-2">
                <!-- hide the compact badge on very small screens to save space; color indicates status -->
                <span
                  class="hidden sm:inline-block px-2 py-1 rounded-full text-xs font-semibold whitespace-nowrap"
                  :class="
                    analysis.summary !== null
                      ? 'bg-green-100 text-green-800'
                      : 'bg-yellow-100 text-yellow-800'
                  "
                >
                  {{
                    analysis.summary === null ? t('analysis.processing') : t('analysis.completed')
                  }}
                </span>

                <div class="flex items-center gap-2">
                  <button
                    @click.stop="analysis.transcription && analysis.id && openAnalysis(analysis.id)"
                    :disabled="analysis.summary === null"
                    :class="[
                      'p-1 rounded shrink-0 inline-flex items-center justify-center',
                      analysis.summary === null ? 'opacity-50 cursor-not-allowed' : '',
                    ]"
                    :title="t('analysis.view')"
                    :aria-label="t('analysis.view')"
                  >
                    <Eye
                      class="w-7 h-7 sm:w-5 sm:h-5 cursor-pointer"
                      color="var(--color-tertiary)"
                    />
                  </button>

                  <button
                    @click.stop="analysis.id && eraseAnalysis(analysis.id)"
                    class="p-1 rounded hover:opacity-75 shrink-0 inline-flex items-center justify-center"
                    :title="t('analysis.delete')"
                    :aria-label="t('analysis.delete')"
                  >
                    <Trash2 class="w-7 h-7 sm:w-5 sm:h-5" color="var(--color-tertiary)" />
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

import youtubeLogo from '@/assets/logos/youtube-logo.jpg';
import tiktokLogo from '@/assets/logos/tiktok-logo.png';
import instagramLogo from '@/assets/logos/instagram-logo.png';

const LogoSource: Record<string, string> = {
  youtube_platform: youtubeLogo,
  tiktok_platform: tiktokLogo,
  instagram_platform: instagramLogo,
};

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

onMounted(() => {
  fetchAnalyses();
});

const isSidebarOpen = ref(false);
</script>
