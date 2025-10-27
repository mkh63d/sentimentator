<template>
  <section class="search-panel">
    <h2>Links</h2>

    <!-- Single link input -->
    <div class="row">
      <input
        v-model="current"
        @keyup.enter="addLink"
        :class="{ invalid: current && !isValidUrl(current) }"
        placeholder="Enter a URL (press Enter or Add)"
        aria-label="Add link"
      />
      <button @click="addLink" :disabled="!current">Add</button>
    </div>

    <!-- Validation / hint -->
    <p class="hint" v-if="invalidCount > 0">{{ invalidCount }} invalid link(s) ignored.</p>

    <!-- Links list -->
    <ul class="links" v-if="links.length">
      <li v-for="(link, i) in links" :key="link">
        <a :href="link" target="_blank" rel="noopener noreferrer">{{ link }}</a>
        <button class="remove" @click="removeLink(i)" aria-label="Remove link">✕</button>
      </li>
    </ul>

    <!-- Actions -->
    <div class="actions">
      <button @click="submit" :disabled="!links.length">Analyze</button>
      <button @click="clearAll" :disabled="!links.length">Clear</button>
    </div>

    <LoadingOverlay />
  </section>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useLoadingStore } from '@/stores/useLoadingStore';
import LoadingOverlay from '@/components/LoadingOverlay.vue';

const loading = useLoadingStore();

const emit = defineEmits(['analyze']);

const current = ref('');
const bulk = ref('');
const links = ref([]);

const isValidUrl = (value) => {
  try {
    const u = new URL(value.trim());
    return ['http:', 'https:'].includes(u.protocol);
  } catch {
    return false;
  }
};

const addLink = () => {
  const val = current.value.trim();
  if (!val) return;
  if (isValidUrl(val) && !links.value.includes(val)) {
    links.value.push(val);
  }
  current.value = '';
};

const parseBulk = () => {
  const lines = bulk.value
    .split(/\r?\n/)
    .map((s) => s.trim())
    .filter(Boolean);

  let added = 0;
  for (const l of lines) {
    if (isValidUrl(l) && !links.value.includes(l)) {
      links.value.push(l);
      added++;
    }
  }
  // remove processed lines from textarea
  if (added) {
    // keep only lines that were invalid or duplicates
    bulk.value = lines.filter((l) => !isValidUrl(l) || links.value.indexOf(l) === -1).join('\n');
  }
};

const removeLink = (index) => {
  links.value.splice(index, 1);
};

const clearAll = () => {
  links.value = [];
  bulk.value = '';
  current.value = '';
};

const invalidCount = computed(() => {
  const all = [
    ...bulk.value
      .split(/\r?\n/)
      .map((s) => s.trim())
      .filter(Boolean),
  ];
  return all.filter((l) => !isValidUrl(l)).length;
});

const submit = async () => {
  if (!links.value.length) return;

  // Mock Loading
  loading.show('Analyzing links...');
  await new Promise((resolve) => setTimeout(resolve, 1000));
  loading.hide();

  emit('analyze', links.value.slice());
};
</script>

<style scoped>
.search-panel {
  max-width: 680px;
  margin: 1rem auto;
  padding: 1rem;
  border: 1px solid #e6e6e6;
  border-radius: 6px;
  background: #fafafa;
  font-family:
    system-ui,
    -apple-system,
    'Segoe UI',
    Roboto,
    'Helvetica Neue',
    Arial;
}
h2 {
  margin: 0 0 0.75rem 0;
  font-size: 1.05rem;
}

.row {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 0.75rem;
}
input {
  flex: 1;
  padding: 0.5rem;
  border: 1px solid #ccc;
  border-radius: 4px;
}
textarea {
  width: 100%;
  min-height: 80px;
  padding: 0.5rem;
  border: 1px solid #ccc;
  border-radius: 4px;
  resize: vertical;
}
button {
  padding: 0.45rem 0.7rem;
  border: none;
  background: #2d8cf0;
  color: white;
  border-radius: 4px;
  cursor: pointer;
}
button:disabled {
  opacity: 0.5;
  cursor: default;
}
.invalid {
  border-color: #e55;
}

.links {
  list-style: none;
  padding: 0;
  margin: 0 0 0.75rem 0;
}
.links li {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.35rem 0;
  border-bottom: 1px dashed #eee;
}
.links a {
  color: #0366d6;
  word-break: break-all;
  text-decoration: none;
}
.remove {
  margin-left: auto;
  background: transparent;
  color: #888;
  font-size: 0.9rem;
  padding: 0 0.4rem;
}
.actions {
  display: flex;
  gap: 0.5rem;
}
.hint {
  color: #a00;
  margin: 0 0 0.5rem 0;
}
</style>
