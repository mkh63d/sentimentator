<template>
  <div v-if="searchStore.visible" class="fixed inset-0 z-50 flex items-center justify-center">
    <!-- Backdrop -->
    <div class="absolute inset-0 bg-black/30" @click="closeWindow()"></div>

    <!-- Modal content -->
    <div
      class="relative w-[680px] max-w-[90%] m-4 p-6 border border-primary/20 rounded-lg bg-background shadow-lg"
    >
      <!-- Close button -->
      <button
        class="cursor-pointer absolute right-3 top-3 p-1.5 text-tertiary hover:text-primary rounded-full hover:bg-primary/10 transition-colors"
        @click="closeWindow()"
        :aria-label="t('common.close')"
      >
        <X size="20" />
      </button>

      <h2 class="mb-4 text-xl font-bold text-tertiary">{{ t('modal.searchTitle') }}</h2>

      <!-- Single link input -->
      <div class="flex gap-2 mb-4">
        <input
          class="flex-1 p-2.5 border border-primary/20 rounded-lg bg-secondary/50 text-tertiary placeholder:text-tertiary/50 focus:outline-none focus:border-primary/50"
          v-model="current"
          @keyup.enter="addLink"
          :class="{ 'border-red-600': !isValidUrl(current) && current }"
          :placeholder="t('modal.enterUrl')"
          :aria-label="t('common.add')"
        />
        <button
          class="px-4 py-2 transition-colors border-0 rounded-lg cursor-pointer bg-primary text-background disabled:opacity-50 disabled:cursor-not-allowed hover:bg-primary/90"
          @click="addLink"
          :disabled="!current"
        >
          {{ t('common.add') }}
        </button>
      </div>
      <p class="mb-3 text-sm text-red-600" v-if="!isValidUrl(current) && current">
        {{ t('modal.validUrl') }}
      </p>

      <!-- Validation / hint -->
      <p class="mb-3 text-sm text-red-600" v-if="invalidCount > 0">
        {{ invalidCount }} {{ t('modal.invalidLinks') }}
      </p>

      <!-- Links list -->
      <ul class="list-none p-0 m-0 mb-4 max-h-[300px] overflow-y-auto" v-if="links.length">
        <li
          class="flex items-center gap-2 py-2 border-b border-primary/10"
          v-for="(link, i) in links"
          :key="link"
        >
          <a
            class="no-underline break-all text-textBase hover:text-primary/80"
            :href="link"
            target="_blank"
            rel="noopener noreferrer"
            >{{ link }}</a
          >
          <button
            class="cursor-pointer ml-auto text-tertiary hover:text-primary p-1.5 rounded-full hover:bg-primary/10 transition-colors"
            @click="removeLink(i)"
            :aria-label="t('common.delete')"
          >
            <X size="16" />
          </button>
        </li>
      </ul>

      <!-- Actions -->
      <div class="flex gap-3">
        <button
          class="px-5 py-2 transition-colors rounded-lg cursor-pointer bg-primary text-background disabled:opacity-50 disabled:cursor-not-allowed hover:bg-primary/90"
          @click="submit"
          :disabled="!links.length"
        >
          {{ t('common.analyse') }}
        </button>
        <button
          class="px-5 py-2 transition-colors rounded-lg cursor-pointer bg-tertiary/10 text-tertiary disabled:opacity-50 disabled:cursor-not-allowed hover:bg-tertiary/20"
          @click="clearAll"
          :disabled="!links.length"
        >
          {{ t('common.clear') }}
        </button>
      </div>

      <LoadingOverlay v-model:visible="isLoading" v-model:message="loadingMessage" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useSearchStore } from '@/stores/searchStore';
import { X } from 'lucide-vue-next';
import LoadingOverlay from '@/components/LoadingOverlay.vue';
import { createAnalysis } from '@/api/analyses';
import { fireTimedAlert } from '@/lib/sweetAlerts';
import { useI18n } from 'vue-i18n';

const isLoading = ref(false);
const loadingMessage = ref('');
const searchStore = useSearchStore();
const { t } = useI18n();

const emit = defineEmits(['analyze']);

const current = ref('');
const bulk = ref('');
const links = ref([]);
const whitelist = ref([
  'youtube.com',
  'youtu.be',
  'tiktok.com',
  'instagram.com',
  'googlevideo.com',
  'ytimg.com',
  'tiktokv.com',
]);

const closeWindow = () => {
  clearAll();
  searchStore.hide();
};

const isValidUrl = (value) => {
  try {
    const u = new URL(value.trim());
    return (
      ['http:', 'https:'].includes(u.protocol) &&
      whitelist.value.some((domain) => u.hostname.includes(domain))
    );
  } catch {
    return false;
  }
};

const addLink = () => {
  const val = current.value.trim();
  if (links.value.length >= 1) {
    fireTimedAlert(t('modal.linkLimitReached'), t('modal.premiumFeature'), 2000, false);
    return;
  }

  if (!val) return;
  if (isValidUrl(val) && !links.value.includes(val)) {
    links.value.push(val);
    current.value = '';
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
  loadingMessage.value = t('modal.creatingAnalyses');
  isLoading.value = true;
  try {
    await new Promise((resolve) => setTimeout(resolve, 1000)); // Simulate delay
    for (const link of links.value) {
      await createAnalysis(link);
    }
  } catch (error) {
    console.error('Error during analysis creation:', error);
  }
  isLoading.value = false;
  loadingMessage.value = '';
  searchStore.hide();
  clearAll();
  emit('analyze');
};
</script>
