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
          <X size="16" />
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
import { X } from 'lucide-vue-next';
import LoadingOverlay from '@/components/LoadingOverlay.vue';

const loading = useLoadingStore();

const emit = defineEmits(['analyze', 'add-link', 'clear-all']);

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
    emit('add-link', val);
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
  emit('clear-all', null);
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

  emit('analyze', links.value.slice());

  // Mock Loading
  loading.show('Analyzing links...');
  await new Promise((resolve) => setTimeout(resolve, 1000));
  loading.hide();
};
</script>
