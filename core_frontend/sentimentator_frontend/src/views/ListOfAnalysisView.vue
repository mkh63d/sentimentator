<template>
  <div class="px-8 py-8 max-w-[1400px] mx-auto">
    <div class="flex justify-between items-center mb-8">
      <h1 class="m-0 text-[#2c3e50]">Analysis Library</h1>
      <button @click="createNew" class="bg-blue-500 text-white border-none px-6 py-3 rounded-lg cursor-pointer font-semibold hover:bg-blue-700 transition">+ New Analysis</button>
    </div>

    <div class="flex gap-4 mb-8">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="Search analyses..."
        class="flex-1 px-3 py-3 border border-gray-200 rounded-lg text-base"
      />
      <select v-model="filterStatus" class="px-3 py-3 border border-gray-200 rounded-lg text-base">
        <option value="">All Statuses</option>
        <option value="completed">Completed</option>
        <option value="processing">Processing</option>
        <option value="failed">Failed</option>
      </select>
    </div>

    <div v-if="loading" class="text-center py-12">
      <span class="inline-block w-5 h-5 border-4 border-gray-200 border-t-blue-500 rounded-full animate-spin"></span> Loading library...
    </div>

    <div v-else-if="error" class="text-center py-12">
      <p>{{ error }}</p>
      <button @click="fetchAnalyses" class="bg-gray-200 px-4 py-2 rounded hover:bg-gray-300 transition">Retry</button>
    </div>

    <div v-else-if="filteredAnalyses.length === 0" class="text-center py-12">
      <p>No analyses found in your library.</p>
      <button @click="createNew" class="bg-gray-100 px-4 py-2 rounded hover:bg-gray-200 transition">Create your first analysis</button>
    </div>

    <div v-else class="bg-white rounded-lg overflow-hidden shadow-md">
      <table class="w-full border-collapse">
        <thead class="bg-gray-100">
          <tr>
            <th class="text-left p-4 font-semibold text-[#2c3e50]">Title</th>
            <th class="text-left p-4 font-semibold text-[#2c3e50]">Created</th>
            <th class="text-left p-4 font-semibold text-[#2c3e50]">Status</th>
            <th class="text-left p-4 font-semibold text-[#2c3e50]">Results</th>
            <th class="text-left p-4 font-semibold text-[#2c3e50]">Actions</th>
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
                  analysis.status === 'failed' && 'bg-red-100 text-red-800'
                ]"
              >
                {{ analysis.status }}
              </span>
            </td>
            <td class="p-4 border-t border-gray-100">
              <div v-if="analysis.sentiment" class="flex h-5 rounded overflow-hidden">
                <span class="h-full bg-green-600" :style="`width: ${analysis.sentiment.positive}%`"></span>
                <span class="h-full bg-red-500" :style="`width: ${analysis.sentiment.negative}%`"></span>
                <span class="h-full bg-gray-400" :style="`width: ${analysis.sentiment.neutral}%`"></span>
              </div>
              <span v-else class="text-gray-400">—</span>
            </td>
            <td class="p-4 border-t border-gray-100 flex gap-2 items-center" @click.stop>
              <button @click="openAnalysis(analysis.id)" class="bg-transparent border-none cursor-pointer text-xl p-1" title="View">
                👁️
              </button>
              <button @click="deleteAnalysis(analysis.id)" class="bg-transparent border-none cursor-pointer text-xl p-1 hover:opacity-70" title="Delete">
                🗑️
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';

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
  return analyses.value.filter(analysis => {
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
  if (!confirm('Are you sure you want to delete this analysis?')) return;

  try {
    const response = await fetch(`/api/analyses/${id}`, { method: 'DELETE' });
    if (!response.ok) throw new Error('Failed to delete');
    analyses.value = analyses.value.filter(a => a.id !== id);
  } catch (err) {
    alert('Failed to delete analysis');
  }
};

const formatDate = (dateString: string): string => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};

onMounted(() => {
  fetchAnalyses();
});
</script>
