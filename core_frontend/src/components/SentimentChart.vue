<template>
  <div class="relative w-full max-w-full mx-0 my-auto">
    <Doughnut :data="chartData" :options="chartOptions" />
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue';
import { Doughnut } from 'vue-chartjs';
import {
  Chart as ChartJS,
  ArcElement,
  Tooltip,
  Legend,
  type ChartData,
  type ChartOptions,
} from 'chart.js';

ChartJS.register(ArcElement, Tooltip, Legend);

interface Props {
  positive?: number;
  neutral?: number;
  negative?: number;
}

const props = withDefaults(defineProps<Props>(), {
  positive: 65,
  neutral: 25,
  negative: 10,
});

const chartData = computed<ChartData<'doughnut'>>(() => ({
  labels: ['Positive', 'Neutral', 'Negative'],
  datasets: [
    {
      data: [props.positive, props.neutral, props.negative],
      backgroundColor: ['#10b981', '#6b7280', '#ef4444'],
      borderColor: ['#059669', '#4b5563', '#dc2626'],
      borderWidth: 2,
    },
  ],
}));

const chartOptions = computed<ChartOptions<'doughnut'>>(() => ({
  responsive: true,
  maintainAspectRatio: true,
  plugins: {
    legend: {
      position: 'bottom',
      labels: {
        padding: 15,
        font: {
          size: 12,
        },
      },
    },
    tooltip: {
      callbacks: {
        label: (context) => {
          const label = context.label || '';
          const value = context.parsed;
          return `${label}: ${value}%`;
        },
      },
    },
  },
}));
</script>
