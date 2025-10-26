<template>
    <div class="analysis-library">
        <div class="library-header">
            <h1>Analysis Library</h1>
            <button @click="createNew" class="btn-primary">+ New Analysis</button>
        </div>
        
        <div class="search-filter">
            <input 
                v-model="searchQuery" 
                type="text" 
                placeholder="Search analyses..."
                class="search-input"
            />
            <select v-model="filterStatus" class="filter-select">
                <option value="">All Statuses</option>
                <option value="completed">Completed</option>
                <option value="processing">Processing</option>
                <option value="failed">Failed</option>
            </select>
        </div>
        
        <div v-if="loading" class="loading">
            <span class="spinner"></span> Loading library...
        </div>
        
        <div v-else-if="error" class="error">
            <p>{{ error }}</p>
            <button @click="fetchAnalyses" class="btn-retry">Retry</button>
        </div>
        
        <div v-else-if="filteredAnalyses.length === 0" class="empty">
            <p>No analyses found in your library.</p>
            <button @click="createNew" class="btn-secondary">Create your first analysis</button>
        </div>
        
        <div v-else class="library-table">
            <table>
                <thead>
                    <tr>
                        <th>Title</th>
                        <th>Created</th>
                        <th>Status</th>
                        <th>Results</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr 
                        v-for="analysis in filteredAnalyses" 
                        :key="analysis.id"
                        @click="openAnalysis(analysis.id)"
                        class="table-row"
                    >
                        <td class="title-cell">
                            <strong>{{ analysis.title }}</strong>
                            <span v-if="analysis.description" class="description">
                                {{ analysis.description }}
                            </span>
                        </td>
                        <td>{{ formatDate(analysis.createdAt) }}</td>
                        <td>
                            <span :class="['status-badge', analysis.status]">
                                {{ analysis.status }}
                            </span>
                        </td>
                        <td>
                            <div v-if="analysis.sentiment" class="mini-chart">
                                <span class="bar positive" :style="`width: ${analysis.sentiment.positive}%`"></span>
                                <span class="bar negative" :style="`width: ${analysis.sentiment.negative}%`"></span>
                                <span class="bar neutral" :style="`width: ${analysis.sentiment.neutral}%`"></span>
                            </div>
                            <span v-else class="no-results">—</span>
                        </td>
                        <td class="actions" @click.stop>
                            <button @click="openAnalysis(analysis.id)" class="btn-icon" title="View">
                                👁️
                            </button>
                            <button @click="deleteAnalysis(analysis.id)" class="btn-icon delete" title="Delete">
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
    router.push(`/analysis/${id}`);
};

const createNew = () => {
    router.push('/search');
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

<style scoped>
.analysis-library {
    padding: 2rem;
    max-width: 1400px;
    margin: 0 auto;
}

.library-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
}

.library-header h1 {
    margin: 0;
    color: #2c3e50;
}

.btn-primary {
    background: #3498db;
    color: white;
    border: none;
    padding: 0.75rem 1.5rem;
    border-radius: 6px;
    cursor: pointer;
    font-weight: 600;
}

.btn-primary:hover {
    background: #2980b9;
}

.search-filter {
    display: flex;
    gap: 1rem;
    margin-bottom: 2rem;
}

.search-input, .filter-select {
    padding: 0.75rem;
    border: 1px solid #ddd;
    border-radius: 6px;
    font-size: 1rem;
}

.search-input {
    flex: 1;
}

.loading {
    text-align: center;
    padding: 3rem;
}

.spinner {
    display: inline-block;
    width: 20px;
    height: 20px;
    border: 3px solid #f3f3f3;
    border-top: 3px solid #3498db;
    border-radius: 50%;
    animation: spin 1s linear infinite;
}

@keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
}

.error, .empty {
    text-align: center;
    padding: 3rem;
}

.library-table {
    background: white;
    border-radius: 8px;
    overflow: hidden;
    box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

table {
    width: 100%;
    border-collapse: collapse;
}

thead {
    background: #f8f9fa;
}

th {
    text-align: left;
    padding: 1rem;
    font-weight: 600;
    color: #2c3e50;
}

.table-row {
    cursor: pointer;
    transition: background 0.2s;
}

.table-row:hover {
    background: #f8f9fa;
}

td {
    padding: 1rem;
    border-top: 1px solid #ecf0f1;
}

.title-cell {
    display: flex;
    flex-direction: column;
}

.description {
    font-size: 0.85rem;
    color: #7f8c8d;
    margin-top: 0.25rem;
}

.status-badge {
    padding: 0.25rem 0.75rem;
    border-radius: 12px;
    font-size: 0.85rem;
    font-weight: 600;
}

.status-badge.completed {
    background: #d4edda;
    color: #155724;
}

.status-badge.processing {
    background: #fff3cd;
    color: #856404;
}

.status-badge.failed {
    background: #f8d7da;
    color: #721c24;
}

.mini-chart {
    display: flex;
    height: 20px;
    border-radius: 4px;
    overflow: hidden;
}

.bar {
    height: 100%;
}

.bar.positive { background: #27ae60; }
.bar.negative { background: #e74c3c; }
.bar.neutral { background: #95a5a6; }

.actions {
    display: flex;
    gap: 0.5rem;
}

.btn-icon {
    background: none;
    border: none;
    cursor: pointer;
    font-size: 1.2rem;
    padding: 0.25rem;
}

.btn-icon.delete:hover {
    opacity: 0.7;
}
</style>