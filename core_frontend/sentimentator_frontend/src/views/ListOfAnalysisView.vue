<template>
  <div class="bg-secondary min-h-screen">
    <navbar class="px-2 py-5 shadow-2xl shadow-gray-500">
      <div class="max-w-full mx-auto flex justify-between items-center bg-primary p-3">
        <div class="flex justify-center items-center">
          <div class="px-1 ml-8 cursor-pointer" @click="isSidebarOpen = !isSidebarOpen">
            <svg class="w-6 h-6" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path
                d="M4 6H20M4 12H20M4 18H20"
                stroke="var(--color-secondary)"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              ></path>
            </svg>
          </div>
          <div class="px-5 text-2xl font-bold text-background" @click="() => router.push('/home')">
            Sentimentator
          </div>
        </div>
        <ThemeSwitcher class="mx-5" :changeOnTheme="true"/>
      </div>
    </navbar>
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

      <div class="px-8 py-8 max-w-[1400px] mx-auto flex-1">
        <div class="flex justify-between items-center mb-8">
          <h1 class="m-0 text-3xl font-extrabold text-tertiary">Analysis Library</h1>
          <button
            @click="createNew"
            class="bg-primary text-background border-none px-6 py-3 rounded-lg cursor-pointer font-semibold hover:bg-blue-700 transition"
          >
            + New Analysis
          </button>
        </div>

        <div class="flex gap-4 mb-8">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search analyses..."
            class="flex-1 px-3 py-3 border border-background rounded-lg text-base"
          />
          <select
            v-model="filterStatus"
            class="px-3 py-3 border border-background rounded-lg text-base"
          >
            <option value="">All Statuses</option>
            <option value="completed">Completed</option>
            <option value="processing">Processing</option>
            <option value="failed">Failed</option>
          </select>
        </div>

        <div v-if="loading" class="text-center py-12">
          <span
            class="inline-block w-5 h-5 border-4 border-gray-200 border-t-blue-500 rounded-full animate-spin"
          ></span>
          Loading library...
        </div>

        <div v-else-if="error" class="text-center py-12">
          <p>{{ error }}</p>
          <button
            @click="fetchAnalyses"
            class="bg-gray-200 px-4 py-2 rounded hover:bg-gray-300 transition"
          >
            Retry
          </button>
        </div>

        <div v-else-if="filteredAnalyses.length === 0" class="text-center py-12">
          <p>No analyses found in your library.</p>
          <button
            @click="createNew"
            class="bg-gray-100 px-4 py-2 rounded hover:bg-gray-200 transition"
          >
            Create your first analysis
          </button>
        </div>

        <div v-else class="bg-white rounded-lg overflow-hidden shadow-md">
          <table class="w-full border-collapse">
            <thead class="bg-gray-100">
              <tr>
                <th class="text-left p-4 font-semibold text-tertiary">Title</th>
                <th class="text-left p-4 font-semibold text-tertiary">Created</th>
                <th class="text-left p-4 font-semibold text-tertiary">Status</th>
                <th class="text-left p-4 font-semibold text-tertiary">Results</th>
                <th class="text-left p-4 font-semibold text-tertiary">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="analysis in filteredAnalyses"
                :key="analysis.id"
                @click="openAnalysis(analysis.id)"
                class="cursor-pointer transition hover:bg-gray-100"
              >
                <td class="p-4 border-t border-gray-100 flex flex-col">
                  <strong>{{ analysis.title }}</strong>
                  <span v-if="analysis.description" class="text-sm text-gray-500 mt-1">
                    {{ analysis.description }}
                  </span>
                </td>
                <td class="p-4 border-t border-gray-100">{{ formatDate(analysis.createdAt) }}</td>
                <td class="p-4 border-t border-gray-100">
                  <span
                    :class="[
                      'px-3 py-1 rounded-full text-sm font-semibold',
                      analysis.status === 'completed' && 'bg-green-100 text-green-800',
                      analysis.status === 'processing' && 'bg-yellow-100 text-yellow-800',
                      analysis.status === 'failed' && 'bg-red-100 text-red-800',
                    ]"
                  >
                    {{ analysis.status }}
                  </span>
                </td>
                <td class="p-4 border-t border-gray-100">
                  <div v-if="analysis.sentiment" class="flex h-5 rounded overflow-hidden">
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
                </td>
                <td class="p-4 border-t border-gray-100 flex gap-2 items-center" @click.stop>
                  <button
                    @click="openAnalysis(analysis.id)"
                    class="bg-transparent border-none cursor-pointer text-xl p-1"
                    title="View"
                  >
                    <Eye class="w-5 h-5 inline" />
                  </button>
                  <button
                    @click="deleteAnalysis(analysis.id)"
                    class="bg-transparent border-none cursor-pointer text-xl p-1 hover:opacity-70"
                    title="Delete"
                  >
                    <Trash2 class="w-5 h-5 inline" />
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { Eye, Trash2 } from 'lucide-vue-next';
import Swal from 'sweetalert2';
import Sidebar from '@/components/Sidebar.vue';
import ThemeSwitcher from '@/components/ThemeSwitcher.vue';

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
}

const router = useRouter();
const analyses = ref<Analysis[]>([]);
const loading = ref(true);
const error = ref<string | null>(null);
const searchQuery = ref('');
const filterStatus = ref('');

const filteredAnalyses = computed(() => {
  return analyses.value.filter((analysis) => {
    const matchesSearch = analysis.title.toLowerCase().includes(searchQuery.value.toLowerCase());
    const matchesStatus = !filterStatus.value || analysis.status === filterStatus.value;
    return matchesSearch && matchesStatus;
  });
});

const fetchAnalyses = async () => {
  try {
    loading.value = true;
    error.value = null;
    const response = await fetch('/api/analyses');
    if (!response.ok) throw new Error('Failed to fetch analyses');
    analyses.value = await response.json();
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'An error occurred';
  } finally {
    loading.value = false;
  }
};

const openAnalysis = (id: string) => {
  router.push({ name: 'analysis', params: { id } });
};

const createNew = () => {
  router.push({ name: 'search' });
};

const deleteAnalysis = async (id: string) => {
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
    // Mock handling
    if (result.isConfirmed) {
      try {
        const response = await fetch(`/api/analyses/${id}`, { method: 'DELETE' });
        if (!response.ok) throw new Error('Failed to delete');
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
