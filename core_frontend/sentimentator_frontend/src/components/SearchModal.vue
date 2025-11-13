<template>
  <div v-if="modelValue" class="fixed inset-0 z-50 flex items-center justify-center">
    <!-- Backdrop -->
    <div class="absolute inset-0 bg-background/30" @click="$emit('update:modelValue', false)"></div>

    <!-- Modal content -->
    <div class="relative w-[680px] max-w-[90%] m-4 p-6 border border-primary/20 rounded-lg bg-background shadow-lg">
      <!-- Close button -->
      <button
        class="absolute right-3 top-3 p-1.5 text-tertiary hover:text-primary rounded-full hover:bg-primary/10 transition-colors"
        @click="$emit('update:modelValue', false)"
        aria-label="Close"
      >
        <X size="20" />
      </button>

      <h2 class="mb-4 text-xl font-bold text-tertiary">Add Links for Analysis</h2>

      <!-- Single link input -->
      <div class="flex gap-2 mb-4">
        <input
          class="flex-1 p-2.5 border border-primary/20 rounded-lg bg-secondary/50 text-tertiary placeholder:text-tertiary/50 focus:outline-none focus:border-primary/50"
          v-model="current"
          @keyup.enter="addLink"
          :class="{ 'border-red-500': current && !isValidUrl(current) }"
          placeholder="Enter a URL (press Enter or Add)"
          aria-label="Add link"
        />
        <button
          class="py-2 px-4 border-0 bg-primary text-background rounded-lg cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed hover:bg-primary/90 transition-colors"
          @click="addLink"
          :disabled="!current"
        >
          Add
        </button>
      </div>

      <!-- Validation / hint -->
      <p class="text-red-600 mb-3 text-sm" v-if="invalidCount > 0">
        {{ invalidCount }} invalid link(s) ignored.
      </p>

      <!-- Links list -->
      <ul class="list-none p-0 m-0 mb-4 max-h-[300px] overflow-y-auto" v-if="links.length">
        <li
          class="flex items-center gap-2 py-2 border-b border-primary/10"
          v-for="(link, i) in links"
          :key="link"
        >
          <a
            class="text-textBase break-all no-underline hover:text-primary/80"
            :href="link"
            target="_blank"
            rel="noopener noreferrer"
            >{{ link }}</a
          >
          <button
            class="ml-auto text-tertiary hover:text-primary p-1.5 rounded-full hover:bg-primary/10 transition-colors"
            @click="removeLink(i)"
            aria-label="Remove link"
          >
            <X size="16" />
          </button>
        </li>
      </ul>

      <!-- Actions -->
      <div class="flex gap-3">
        <button
          class="py-2 px-5 bg-primary text-background rounded-lg cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed hover:bg-primary/90 transition-colors"
          @click="submit"
          :disabled="!links.length"
        >
          Analyze
        </button>
        <button
          class="py-2 px-5 bg-tertiary/10 text-tertiary rounded-lg cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed hover:bg-tertiary/20 transition-colors"
          @click="clearAll"
          :disabled="!links.length"
        >
          Clear
        </button>
      </div>

      <LoadingOverlay />
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useLoadingStore } from '@/stores/loadingStore';
import { X } from 'lucide-vue-next';
import LoadingOverlay from '@/components/LoadingOverlay.vue';

const loading = useLoadingStore();

const props = defineProps({
  modelValue: {
    type: Boolean,
    required: true,
  },
});

const emit = defineEmits(['analyze', 'add-link', 'clear-all', 'update:modelValue']);

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
  emit('update:modelValue', false); // Close popup after submission

  // Mock Loading
  loading.show('Analyzing links...');
  await new Promise((resolve) => setTimeout(resolve, 1000));
  loading.hide();
};
</script>
