<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { useStore } from '../stores/app';
import type { DrawResult } from '../types';
import * as api from '../api';

const state = useStore();
const spinning = ref(false);
const rotation = ref(0);
const result = ref<DrawResult | null>(null);
const showResult = ref(false);
const wheelCanvas = ref<HTMLCanvasElement | null>(null);

// Color palette for sectors - 9 distinct colors
const sectorColors = [
  '#FF6B6B', // 红
  '#4ECDC4', // 青绿
  '#45B7D1', // 蓝
  '#96CEB4', // 薄荷绿
  '#FFEAA7', // 淡黄
  '#DDA0DD', // 梅红
  '#98D8C8', // 绿松石
  '#F7DC6F', // 金黄
  '#BB8FCE', // 紫
];

// Calculate total probability for normalization
const totalProbability = computed(() => {
  return state.prizes.reduce((sum, p) => sum + p.probability, 0);
});

// Calculate sector angles for each prize
const sectors = computed(() => {
  let currentAngle = 0;
  return state.prizes.map(prize => {
    const angle = (prize.probability / totalProbability.value) * 360;
    const sector = {
      prize,
      startAngle: currentAngle,
      endAngle: currentAngle + angle,
    };
    currentAngle += angle;
    return sector;
  });
});

// Draw the wheel on canvas
function drawWheel() {
  const canvas = wheelCanvas.value;
  if (!canvas) return;

  const ctx = canvas.getContext('2d');
  if (!ctx) return;

  const centerX = canvas.width / 2;
  const centerY = canvas.height / 2;
  const radius = Math.min(centerX, centerY) - 5;

  // Clear canvas
  ctx.clearRect(0, 0, canvas.width, canvas.height);

  // Draw each sector
  sectors.value.forEach((sector, index) => {
    const startAngle = ((sector.startAngle - 90) * Math.PI) / 180; // -90 to start from top
    const endAngle = ((sector.endAngle - 90) * Math.PI) / 180;

    // Draw sector
    ctx.beginPath();
    ctx.moveTo(centerX, centerY);
    ctx.arc(centerX, centerY, radius, startAngle, endAngle);
    ctx.closePath();
    ctx.fillStyle = sectorColors[index % sectorColors.length];
    ctx.fill();
    ctx.strokeStyle = 'rgb(232 121 249)';
    ctx.lineWidth = 2;
    ctx.stroke();

    const midAngle = (startAngle + endAngle) / 2;
    const iconRadius = radius * 0.65;
    const iconX = centerX + Math.cos(midAngle) * iconRadius;
    const iconY = centerY + Math.sin(midAngle) * iconRadius;

    ctx.font = '26px Arial';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    ctx.fillText(sector.prize.icon, iconX, iconY);
  });

  // Draw center circle
  ctx.beginPath();
  ctx.arc(centerX, centerY, 35, 0, Math.PI * 2);
  ctx.fillStyle = '#7c3aed';
  ctx.fill();
  ctx.strokeStyle = '#1A202C';
  ctx.lineWidth = 3;
  ctx.stroke();
}

async function spin() {
  if (!state.currentKid) {
    alert('请先选择一个小朋友');
    return;
  }
  if (state.currentKid.points < 2) {
    alert('积分不足，需要至少 2 分才能抽奖');
    return;
  }
  if (state.prizes.length === 0) {
    alert('暂无奖项配置');
    return;
  }
  if (spinning.value) return;

  spinning.value = true;
  result.value = null;
  showResult.value = false;

  try {
    const drawResult = await api.draw(state.currentKid.id);
    result.value = drawResult;

    // Find the winning sector
    const prizeIndex = state.prizes.findIndex(p => p.id === drawResult.id);
    const sector = sectors.value[prizeIndex];
    const targetAngle = sector.startAngle + (sector.endAngle - sector.startAngle) / 2;

    // Calculate rotation - pointer is at top (0 degrees visually)
    const extraRotations = Math.floor(Math.random() * 3 + 4) * 360;
    const finalRotation = rotation.value + extraRotations + (360 - targetAngle);

    rotation.value = finalRotation;

    setTimeout(() => {
      showResult.value = true;
      spinning.value = false;
    }, 4000);
  } catch (e) {
    alert((e as Error).message);
    spinning.value = false;
  }
}

function closeResult() {
  showResult.value = false;
  result.value = null;
}

onMounted(() => {
  if (state.prizes.length === 0) {
    api.getPrizes().then(prizes => {
      (state as any).prizes = prizes;
      drawWheel();
    }).catch(console.error);
  } else {
    drawWheel();
  }
});

// Redraw wheel when prizes change
watch(() => state.prizes, () => {
  drawWheel();
}, { deep: true });
</script>

<template>
  <div class="p-4">
    <h2 class="text-lg font-semibold mb-4">抽奖转盘</h2>

    <div class="flex gap-6 items-center justify-center">
      <!-- Wheel Container -->
      <div class="relative">
        <!-- Pointer -->
        <div class="absolute top-0 left-1/2 -translate-x-1/2 -translate-y-1 z-20">
          <div
            class="w-0 h-0 border-l-8 border-r-8 border-t-12 border-l-transparent border-r-transparent border-t-purple-600"
            style="border-top-width: 20px;"></div>
        </div>

        <!-- Wheel -->
        <div class="rounded-full overflow-hidden transition-transform origin-center"
          :style="{ transform: `rotate(${rotation}deg)`, transitionDuration: spinning ? '4s' : '0ms', transitionTimingFunction: 'ease-out' }">
          <canvas ref="wheelCanvas" width="480" height="480" />
        </div>

        <!-- Center Button -->
        <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 z-10">
          <button @click="spin" :disabled="spinning"
            class="w-16 h-16 rounded-full bg-purple-600 text-white font-bold shadow-lg flex items-center justify-center text-sm"
            :class="spinning ? 'opacity-50 cursor-not-allowed' : 'hover:bg-purple-700'">
            <span v-if="!spinning">抽奖</span>
            <span v-else>转...</span>
          </button>
        </div>
      </div>

      <!-- Prize List -->
      <div class="bg-white rounded-lg border border-gray-200 p-4 w-48">
        <h3 class="text-sm font-medium text-gray-600 mb-3">奖项列表</h3>
        <div class="space-y-2">
          <div v-for="prize in state.prizes" :key="prize.id" class="flex items-center gap-2 text-sm">
            <span class="text-xl">{{ prize.icon }}</span>
            <span class="flex-1 truncate">{{ prize.name }}</span>
            <span class="text-gray-500">{{ prize.probability }}%</span>
          </div>
        </div>
        <div class="mt-3 pt-3 border-t text-sm text-gray-500">
          每次抽奖消耗 2 积分
        </div>
      </div>
    </div>

    <p class="text-center text-gray-500 text-sm mt-4">
      {{ state.currentKid ? `${state.currentKid.avatar} ${state.currentKid.name} 可以抽奖` : '请先选择小朋友' }}
    </p>

    <!-- Result Modal -->
    <div v-if="showResult && result" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div class="bg-white rounded-xl p-8 text-center max-w-sm">
        <div class="text-6xl mb-4">{{ result.icon }}</div>
        <h3 class="text-xl font-bold text-purple-600 mb-2">恭喜获得</h3>
        <p class="text-2xl font-bold mb-6">{{ result.name }}</p>
        <button @click="closeResult" class="px-6 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700">
          确定
        </button>
      </div>
    </div>
  </div>
</template>
