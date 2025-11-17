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

      <div class="px-8 py-8 max-w-[1400px] mx-auto flex-1">
        <div class="flex items-center justify-between mb-8">
          <h1 class="m-0 text-3xl font-extrabold text-tertiary">Analysis Library</h1>
          <div class="flex gap-3">
            <button
              @click="fetchAnalyses"
              class="flex items-center gap-2 px-4 py-3 font-semibold transition border-none rounded-lg cursor-pointer bg-primary text-background hover:opacity-80"
              title="Refresh analyses"
            >
              <RotateCcw class="w-5 h-5" />
              Refresh
            </button>
            <button
              @click="createNew"
              class="flex items-center justify-center px-6 py-3 font-semibold transition border-none rounded-lg cursor-pointer bg-primary text-background hover:opacity-80"
            >
              <Plus class="pr-2" />
              <p>New Analysis</p>
            </button>
          </div>
        </div>

        <div class="flex gap-4 mb-8">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search analyses..."
            class="flex-1 px-3 py-3 text-base border rounded-lg border-background"
          />
          <!-- <select
            v-model="filterStatus"
            class="px-3 py-3 text-base border rounded-lg cursor-pointer border-background"
          >
            <option value="">All Statuses</option>
            <option value="completed">Completed</option>
            <option value="processing">Processing</option>
            <option value="failed">Failed</option>
          </select> -->
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

        <div v-else class="overflow-hidden rounded-lg shadow-md bg-background/80">
          <table class="w-full border-collapse">
            <thead class="border-b border-gray-100 bg-background/90">
              <tr>
                <th class="p-4 font-semibold text-left text-tertiary">Title</th>
                <!-- <th class="p-4 font-semibold text-left text-tertiary">Created</th> -->
                <th class="p-4 font-semibold text-left text-tertiary">Status</th>
                <!-- <th class="p-4 font-semibold text-left text-tertiary">Results</th> -->
                <th class="p-4 font-semibold text-left text-tertiary">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="analysis in filteredAnalyses"
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
                  <strong>{{ analysis.url }}</strong>
                </td>
                <!-- <td class="p-4 text-tertiary">{{ formatDate(analysis.createdAt) }}</td> -->
                <td class="p-4">
                  <span
                    :class="[
                      'px-3 py-1 rounded-full text-sm font-semibold',
                      analysis.summary !== null && 'bg-green-100 text-green-800',
                      analysis.summary === null && 'bg-yellow-100 text-yellow-800',
                      // analysis.status === 'failed' && 'bg-red-100 text-red-800',
                    ]"
                  >
                    {{ analysis.summary === null ? 'Processing' : 'Completed' }}
                  </span>
                </td>
                <!-- <td class="p-4 border-t border-gray-100">
                  <div v-if="analysis.sentiment" class="flex h-5 overflow-hidden rounded">
                    <span
                      class="h-full bg-green-600"
                      :style="`width: ${analysis.sentiment.positive}%`"
                    ></span>
                    <span
                      class="h-full bg-red-500"
                      :style="`width: ${analysis.sentiment.negative}%`"
                    ></span>
                    <span
                      class="h-full bg-gray-400"
                      :style="`width: ${analysis.sentiment.neutral}%`"
                    ></span>
                  </div>
                  <span v-else class="text-gray-400">—</span>
                </td> -->
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
