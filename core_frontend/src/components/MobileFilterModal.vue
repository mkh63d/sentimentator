<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition-opacity duration-300 ease-out"
      leave-active-class="transition-opacity duration-200 ease-in"
      enter-from-class="opacity-0"
      leave-to-class="opacity-0"
    >
      <div
        v-if="modelValue"
        class="fixed inset-0 z-50 flex items-end justify-center"
        @click.self="close"
      >
        <div class="fixed inset-0 bg-black/50" @click="close" aria-hidden="true"></div>
      </div>
    </Transition>

    <Transition
      enter-active-class="transition-transform duration-300 ease-out"
      leave-active-class="transition-transform duration-200 ease-in"
      enter-from-class="translate-y-full"
      leave-to-class="translate-y-full"
    >
      <div
        v-if="modelValue"
        class="fixed bottom-0 left-0 right-0 z-60 w-full bg-background rounded-t-3xl shadow-xl p-6 max-h-[85vh] overflow-y-auto"
        @click.stop
      >
        <div class="flex items-center justify-between mb-6">
          <h2 class="text-2xl font-bold text-tertiary">{{ t('modal.filterTitle') }}</h2>
          <button
            @click="close"
            class="p-2 rounded-full hover:bg-secondary transition"
            :aria-label="t('common.close')"
          >
            <X class="w-6 h-6 text-tertiary" />
          </button>
        </div>

        <div class="space-y-4">
          <div>
            <label for="mobile-search" class="block text-sm font-medium text-tertiary mb-2">
              {{ t('common.search') }}
            </label>
            <input
              id="mobile-search"
              v-model="localSearchQuery"
              type="text"
              :placeholder="t('analysis.searchAnalyses')"
              class="w-full px-4 py-3 text-base border rounded-lg border-tertiary bg-background text-tertiary focus:outline-none focus:ring-2 focus:ring-primary"
            />
          </div>

          <div>
            <label for="mobile-status" class="block text-sm font-medium text-tertiary mb-2">
              {{ t('modal.status') }}
            </label>
            <select
              id="mobile-status"
              v-model="localFilterStatus"
              class="w-full px-4 py-3 text-base border rounded-lg cursor-pointer border-tertiary bg-background text-tertiary focus:outline-none focus:ring-2 focus:ring-primary"
            >
              <option value="">{{ t('analysis.allStatuses') }}</option>
              <option value="completed">{{ t('analysis.completed') }}</option>
              <option value="processing">{{ t('analysis.processing') }}</option>
              <option value="failed">{{ t('analysis.failed') }}</option>
            </select>
          </div>
        </div>

        <div class="flex gap-3 mt-6">
          <button
            @click="clearFilters"
            class="flex-1 px-4 py-3 text-base font-semibold transition rounded-lg bg-secondary text-tertiary hover:opacity-80"
          >
            {{ t('common.clear') }}
          </button>
          <button
            @click="applyFilters"
            class="flex-1 px-4 py-3 text-base font-semibold transition rounded-lg bg-primary text-background hover:opacity-80"
          >
            {{ t('common.apply') }}
          </button>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { X } from 'lucide-vue-next';
import { useI18n } from 'vue-i18n';

interface Props {
  modelValue: boolean;
  searchQuery: string;
  filterStatus: string;
}

const { t } = useI18n();

interface Emits {
  (e: 'update:modelValue', value: boolean): void;
  (e: 'update:searchQuery', value: string): void;
  (e: 'update:filterStatus', value: string): void;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

const localSearchQuery = ref(props.searchQuery);
const localFilterStatus = ref(props.filterStatus);

watch(
  () => props.searchQuery,
  (newVal) => {
    localSearchQuery.value = newVal;
  },
);

watch(
  () => props.filterStatus,
  (newVal) => {
    localFilterStatus.value = newVal;
  },
);

const close = () => {
  emit('update:modelValue', false);
};

const applyFilters = () => {
  emit('update:searchQuery', localSearchQuery.value);
  emit('update:filterStatus', localFilterStatus.value);
  close();
};

const clearFilters = () => {
  localSearchQuery.value = '';
  localFilterStatus.value = '';
  emit('update:searchQuery', '');
  emit('update:filterStatus', '');
  close();
};
</script>
