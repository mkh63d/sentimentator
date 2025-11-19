<template>
  <div class="min-h-screen bg-secondary">
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
      </Sidebar>

      <div class="px-4 sm:px-8 py-8 max-w-[1400px] mx-auto flex-1">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between mb-6 gap-3 sm:gap-0">
          <h1 class="m-0 text-2xl sm:text-3xl font-extrabold text-tertiary min-w-0 truncate">Analysis Library</h1>
          <div class="flex gap-2 sm:gap-3 items-center">
        <button
          @click="fetchAnalyses"
          class="flex items-center gap-2 px-2 py-2 sm:px-4 sm:py-3 text-sm sm:text-base font-semibold transition rounded-lg cursor-pointer bg-primary text-background hover:opacity-80"
          title="Refresh analyses"
        >
          <RotateCcw class="w-4 h-4 sm:w-5 sm:h-5" />
          <span class="hidden sm:inline">Refresh</span>
        </button>
        <button
          @click="createNew"
          class="flex items-center justify-center px-2 py-2 sm:px-6 sm:py-3 text-sm sm:text-base font-semibold transition rounded-lg cursor-pointer bg-primary text-background hover:opacity-80"
          title="New Analysis"
        >
          <Plus class="pr-1 w-4 h-4 sm:w-5 sm:h-5" />
          <span class="hidden sm:inline">New Analysis</span>
        </button>
          </div>
        </div>

        <div class="flex gap-3 mb-6">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search analyses..."
            class="flex-1 px-3 py-2 text-sm sm:text-base border rounded-lg border-background"
          />
          <select
            v-model="filterStatus"
            class="px-3 py-3 text-base border rounded-lg cursor-pointer border-background"
          >
            <option value="">All Statuses</option>
            <option value="completed">Completed</option>
            <option value="processing">Processing</option>
            <option value="failed">Failed</option>
          </select>
        </div>

        <div v-if="loading" class="py-12 text-center">
          <span
            class="inline-block w-5 h-5 border-4 border-gray-200 rounded-full border-t-blue-500 animate-spin"
          ></span>
          Loading library...
        </div>

        <div v-else-if="error" class="py-12 text-center">
          <p>
            Error fetching library: <br /><code>{{ error }}</code>
          </p>
          <button
            @click="fetchAnalyses"
            class="px-4 py-2 transition bg-gray-200 rounded hover:bg-gray-300"
          >
            Retry
          </button>
        </div>

        <div v-else-if="filteredAnalyses.length === 0" class="py-12 text-center">
          <p>No analyses found in your library.</p>
          <button
            @click="createNew"
            class="px-4 py-2 transition bg-gray-100 rounded hover:bg-gray-200"
          >
            Create your first analysis
          </button>
        </div>

        <!-- Desktop / Tablet: regular table -->
        <div class="hidden md:block overflow-hidden rounded-lg shadow-md bg-background/80">
          <div class="overflow-x-auto"> <!-- allow horizontal scroll if needed -->
            <table class="w-full border-collapse">
              <thead class="border-b border-gray-100 bg-background/90">
                <tr>
                  <th class="p-4 font-semibold text-left text-tertiary">Title</th>
                  <th class="p-4 font-semibold text-left text-tertiary">Created</th>
                  <th class="p-4 font-semibold text-left text-tertiary">Updated</th>
                  <th class="p-4 font-semibold text-left text-tertiary">Status</th>
                  <th class="p-4 font-semibold text-left text-tertiary">Actions</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="analysis in filteredAnalyses.filter(a => !filterStatus || ((a.summary === null ? 'processing' : 'completed') === filterStatus))"
                  :key="analysis.id"
                  @click="analysis.transcription && analysis.id && openAnalysis(analysis.id)"
                  :class="[
                    analysis && analysis.transcription && analysis.id
                      ? 'cursor-pointer transition hover:bg-background/10'
                      : '',
                    'bg-background/90 border-b border-gray-100',
                  ]"
                >
                  <td class="flex flex-col p-4 text-tertiary">
                    <strong class="block truncate max-w-[40ch]">{{ analysis.url }}</strong>
                  </td>
                  <td class="p-4 text-tertiary">{{ analysis.created_at ? formatDate(analysis.created_at) : '' }}</td>
                  <td class="p-4 text-tertiary">{{ analysis.updated_at ? formatDate(analysis.updated_at) : '' }}</td>
                  <td class="p-4">
                    <span
                      :class="[
                        'px-3 py-1 rounded-full text-sm font-semibold',
                        analysis.summary !== null && 'bg-green-100 text-green-800',
                        analysis.summary === null && 'bg-yellow-100 text-yellow-800',
                      ]"
                    >
                      {{ analysis.summary === null ? 'Processing' : 'Completed' }}
                    </span>
                  </td>
                  <td class="flex items-center gap-2 p-4" @click.stop>
                    <button
                      @click="analysis.transcription && analysis.id && openAnalysis(analysis.id)"
                      :class="
                        analysis.summary === null
                          ? 'bg-transparent border-none text-xl p-1 opacity-50 cursor-not-allowed pointer-events-none'
                          : 'bg-transparent border-none cursor-pointer text-xl p-1'
                      "
                      :disabled="analysis.summary === null"
                      title="View"
                    >
                      <Eye class="inline w-5 h-5" color="var(--color-tertiary)" />
                    </button>
                    <button
                      @click="analysis.id && eraseAnalysis(analysis.id)"
                      class="p-1 text-xl bg-transparent border-none cursor-pointer hover:opacity-70"
                      title="Delete"
                    >
                      <Trash2 class="inline w-5 h-5" color="var(--color-tertiary)" />
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Mobile: stacked cards -->
        <div class="md:hidden space-y-3">
          <div
            v-for="analysis in filteredAnalyses"
            :key="analysis.id"
            class="border rounded-lg p-3 flex flex-col gap-2"
            :class="[
              // color the whole card based on status on small screens so the status badge can be hidden
              analysis.summary !== null ? 'bg-green-50 border-green-200' : 'bg-yellow-50 border-yellow-200',
              analysis && analysis.transcription && analysis.id ? 'hover:bg-background/10' : '',
            ]"
          >
            <div class="flex justify-between items-start gap-3">
              <div class="min-w-0">
                <div
                  class="text-sm font-semibold text-tertiary truncate max-w-[60%] sm:max-w-[70%]"
                  :title="analysis.url"
                >
                  {{ analysis.url }}
                </div>
                <div class="text-xs text-tertiary/70 mt-1">
                  <div>Created: {{ analysis.created_at ? formatDate(analysis.created_at) : '' }}</div>
                  <div>Updated: {{ analysis.updated_at ? formatDate(analysis.updated_at) : '' }}</div>
                </div>
              </div>

              <div class="flex flex-col items-end gap-2">
                <!-- hide the compact badge on very small screens to save space; color indicates status -->
                <span
                  class="hidden sm:inline-block px-2 py-1 rounded-full text-xs font-semibold whitespace-nowrap"
                  :class="analysis.summary !== null ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'"
                >
                  {{ analysis.summary === null ? 'Processing' : 'Completed' }}
                </span>

                <div class="flex items-center gap-2">
                  <button
                    @click.stop="analysis.transcription && analysis.id && openAnalysis(analysis.id)"
                    :disabled="analysis.summary === null"
                    :class="[
                      'p-1 rounded shrink-0 inline-flex items-center justify-center',
                      analysis.summary === null ? 'opacity-50 cursor-not-allowed' : '',
                    ]"
                    title="View"
                    aria-label="View analysis"
                  >
                    <Eye class="w-5 h-5" color="var(--color-tertiary)" />
                  </button>

                  <button
                    @click.stop="analysis.id && eraseAnalysis(analysis.id)"
                    class="p-1 rounded hover:opacity-75 shrink-0 inline-flex items-center justify-center"
                    title="Delete"
                    aria-label="Delete analysis"
                  >
                    <Trash2 class="w-5 h-5" color="var(--color-tertiary)" />
                  </button>
                </div>
              </div>
            </div>

            <div v-if="analysis.summary === null" class="text-xs text-tertiary/60">
              Summary is not ready yet.
            </div>
          </div>
        </div>
      </div>
    </div>
    <SearchModal v-model="searchStore.visible" @analyze="fetchAnalyses" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { Plus, Eye, Trash2, RotateCcw } from 'lucide-vue-next';
import Swal from 'sweetalert2';
import Sidebar from '@/components/SideBar.vue';
import NavbarComponent from '@/components/NavbarComponent.vue';
import SearchModal from '@/components/SearchModal.vue';
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
const analyses = ref<Analysis[]>([]);
const realAnalyses = ref<components['schemas']['Analysis'][]>([]);
const loading = ref(true);
const error = ref<string | null>(null);
const searchQuery = ref('');
const filterStatus = ref('');

const filteredAnalyses = computed(() => {
  return realAnalyses.value.filter((analysis) => {
    const matchesSearch = analysis.url.toLowerCase().includes(searchQuery.value.toLowerCase());
    // const matchesStatus = !filterStatus.value || analysis.status === filterStatus.value;
    return matchesSearch;
  });
  // return analyses.value.filter((analysis) => {
  //   const matchesSearch = analysis.title.toLowerCase().includes(searchQuery.value.toLowerCase());
  //   const matchesStatus = !filterStatus.value || analysis.status === filterStatus.value;
  //   return matchesSearch && matchesStatus;
  // });
});

const fetchAnalyses = async () => {
  try {
    loading.value = true;
    error.value = null;
    const response = await getAnalyses();
    console.log('Fetched analyses response:', response);

    // TODO: Adjust mapping for description, createdAt, status, sentiment based on real response structure
    // Waiting for backend API to provide these fields
    analyses.value = await response.map((a) => ({
      id: a.id,
      title: a.url,
      description: '',
      createdAt: new Date().toISOString(),
      status: 'processing',
      sentiment: { positive: 33, negative: 33, neutral: 34 },
      transcription: a.transcription,
    }));
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
  const result = await Swal.fire({
    title: 'Are you sure?',
    text: 'Do you really want to delete this analysis?',
    icon: 'warning',
    showCancelButton: true,
    confirmButtonColor: '#d33',
    cancelButtonColor: '#3085d6',
    confirmButtonText: 'Yes, delete it!',
    cancelButtonText: 'Cancel',
  }).then(async (result) => {
    if (result.isConfirmed) {
      try {
        const response = await deleteAnalysis(id);
        analyses.value = analyses.value.filter((a) => a.id !== id);
        Swal.fire({
          title: 'Deleted!',
          text: 'This analysis has been deleted.',
          icon: 'success',
          timer: 2000,
          showConfirmButton: false,
        });
      } catch (err) {
        Swal.fire({
          title: 'Error!',
          text: 'Failed to delete analysis.',
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
