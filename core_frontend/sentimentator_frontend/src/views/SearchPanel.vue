<template>
  <section
    class="max-w-[680px] mx-auto my-4 p-4 border border-gray-200 rounded-md bg-gray-50 font-sans"
  >
    <h2 class="mb-3 text-[1.05rem]">Links</h2>

    <!-- Single link input -->
    <div class="flex gap-2 mb-3">
      <input
        class="flex-1 p-2 border border-gray-300 rounded"
        v-model="current"
        @keyup.enter="addLink"
        :class="{ 'border-red-500': current && !isValidUrl(current) }"
        placeholder="Enter a URL (press Enter or Add)"
        aria-label="Add link"
      />
      <button
        class="py-1 px-3 border-0 bg-[#2d8cf0] text-white rounded cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
        @click="addLink"
        :disabled="!current"
      >
        Add
      </button>
    </div>

    <!-- Validation / hint -->
    <p class="text-red-700 mb-2" v-if="invalidCount > 0">
      {{ invalidCount }} invalid link(s) ignored.
    </p>

    <!-- Links list -->
    <ul class="list-none p-0 m-0 mb-3" v-if="links.length">
      <li
        class="flex items-center gap-2 py-1 border-b border-dashed border-gray-200"
        v-for="(link, i) in links"
        :key="link"
      >
        <a
          class="text-blue-600 break-all no-underline"
          :href="link"
          target="_blank"
          rel="noopener noreferrer"
          >{{ link }}</a
        >
        <button
          class="ml-auto bg-transparent text-gray-500 text-sm px-1"
          @click="removeLink(i)"
          aria-label="Remove link"
        >
          ✕
        </button>
      </li>
    </ul>

    <!-- Actions -->
    <div class="flex gap-2">
      <button
        class="py-1 px-3 border-0 bg-[#2d8cf0] text-white rounded cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
        @click="submit"
        :disabled="!links.length"
      >
        Analyze
      </button>
      <button
        class="py-1 px-3 border-0 bg-[#2d8cf0] text-white rounded cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
        @click="clearAll"
        :disabled="!links.length"
      >
        Clear
      </button>
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
  loading.show('Analyzing links..?');
  await new Promise((resolve) => setTimeout(resolve, 1000));
  loading.hide();

  emit('analyze', links.value.slice());
};
</script>

<style scoped>
/* .search-panel {
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
} */
/* h2 {
  margin: 0 0 0.75rem 0;
  font-size: 1.05rem;
} */
/* Tailwind: class="mb-3 text-[1.05rem]" */

/* .row {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 0.75rem;
} */
/* Tailwind: class="flex gap-2 mb-3" */

/* input {
  flex: 1;
  padding: 0.5rem;
  border: 1px solid #ccc;
  border-radius: 4px;
} */
/* Tailwind: class="flex-1 p-2 border border-gray-300 rounded" */

/* textarea {
  width: 100%;
  min-height: 80px;
  padding: 0.5rem;
  border: 1px solid #ccc;
  border-radius: 4px;
  resize: vertical;
} */
/* Tailwind: class="w-full min-h-[80px] p-2 border border-gray-300 rounded resize-y" */

/* button {
  padding: 0.45rem 0.7rem;
  border: none;
  background: #2d8cf0;
  color: white;
  border-radius: 4px;
  cursor: pointer;
} */
/* Tailwind: class="py-1 px-3 border-0 bg-blue-500 text-white rounded cursor-pointer" */
/* If you prefer exact color: use bg-[#2d8cf0] instead of bg-blue-500 */

/* button:disabled {
  opacity: 0.5;
  cursor: default;
} */
/* Tailwind: add "disabled:opacity-50 disabled:cursor-not-allowed" to the button */

/* .invalid {
  border-color: #e55;
} */
/* Tailwind: add "border-red-500" to the element (and keep border class) */

/* .links {
  list-style: none;
  padding: 0;
  margin: 0 0 0.75rem 0;
} */
/* Tailwind: class="list-none p-0 m-0 mb-3" */

/* .links li {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.35rem 0;
  border-bottom: 1px dashed #eee;
} */
/* Tailwind: class="flex items-center gap-2 py-1 border-b border-dashed border-gray-200" */

/* .links a {
  color: #0366d6;
  word-break: break-all;
  text-decoration: none;
} */
/* Tailwind: class="text-blue-600 break-all no-underline" */

/* .remove {
  margin-left: auto;
  background: transparent;
  color: #888;
  font-size: 0.9rem;
  padding: 0 0.4rem;
} */
/* Tailwind: class="ml-auto bg-transparent text-gray-500 text-sm px-1" */

/* .actions {
  display: flex;
  gap: 0.5rem;
} */
/* Tailwind: class="flex gap-2" */

/* .hint {
  color: #a00;
  margin: 0 0 0.5rem 0;
} */
/* Tailwind: class="text-red-700 mb-2" */
</style>
