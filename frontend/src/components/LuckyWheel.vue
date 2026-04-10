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

// Sky-themed color palette for sectors
const sectorColors = [
  '#81D4FA', // sky blue
  '#4FC3F7', // ocean blue
  '#FFD54F', // sunny yellow
  '#A5D6A7', // mint green
  '#B39DDB', // lavender
  '#FFCCBC', // peach
  '#80DEEA', // cyan
  '#CE93D8', // light purple
  '#F48FB1', // pink
];

// Calculate total probability for normalization
const totalProbability = computed(() => {
  return state.prizes.reduce((sum, p) => sum + p.probability, 0);
});

// Calculate sector angles for each prize
// If prizes <= 6, distribute by probability; otherwise equal distribution
const sectors = computed(() => {
  const prizeCount = state.prizes.length;
  const isEqualDistribution = prizeCount > 6;
  const sectorAngle = 360 / prizeCount;
  let currentAngle = 0;

  return state.prizes.map(prize => {
    const angle = isEqualDistribution ? sectorAngle : (prize.probability / totalProbability.value) * 360;
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

  // Draw outer glow
  const gradient = ctx.createRadialGradient(centerX, centerY, radius - 10, centerX, centerY, radius + 10);
  gradient.addColorStop(0, 'rgba(129, 212, 250, 0.3)');
  gradient.addColorStop(1, 'rgba(129, 212, 250, 0)');
  ctx.beginPath();
  ctx.arc(centerX, centerY, radius + 10, 0, Math.PI * 2);
  ctx.fillStyle = gradient;
  ctx.fill();

  // Draw each sector
  sectors.value.forEach((sector, index) => {
    const startAngle = ((sector.startAngle - 90) * Math.PI) / 180;
    const endAngle = ((sector.endAngle - 90) * Math.PI) / 180;

    // Draw sector with gradient
    ctx.beginPath();
    ctx.moveTo(centerX, centerY);
    ctx.arc(centerX, centerY, radius, startAngle, endAngle);
    ctx.closePath();

    // Create gradient for each sector
    const sectorGradient = ctx.createRadialGradient(centerX, centerY, 0, centerX, centerY, radius);
    sectorGradient.addColorStop(0, sectorColors[index % sectorColors.length]);
    sectorGradient.addColorStop(1, adjustColor(sectorColors[index % sectorColors.length], -20));

    ctx.fillStyle = sectorGradient;
    ctx.fill();

    // Sector border
    ctx.strokeStyle = 'rgba(255, 255, 255, 0.5)';
    ctx.lineWidth = 3;
    ctx.stroke();

    // Draw icon
    const midAngle = (startAngle + endAngle) / 2;
    const iconRadius = radius * 0.9;
    const iconX = centerX + Math.cos(midAngle) * iconRadius;
    const iconY = centerY + Math.sin(midAngle) * iconRadius;

    ctx.save();
    ctx.translate(iconX, iconY);
    ctx.rotate(midAngle + Math.PI / 2);
    ctx.font = '28px Arial';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    ctx.fillText(sector.prize.icon, 0, 0);
    ctx.restore();
  });

  // Draw center circle with sun design
  ctx.beginPath();
  ctx.arc(centerX, centerY, 40, 0, Math.PI * 2);
  const centerGradient = ctx.createRadialGradient(centerX - 10, centerY - 10, 0, centerX, centerY, 40);
  centerGradient.addColorStop(0, '#FFF9C4');
  centerGradient.addColorStop(0.5, '#FFD54F');
  centerGradient.addColorStop(1, '#FFB300');
  ctx.fillStyle = centerGradient;
  ctx.fill();
  ctx.strokeStyle = '#FFF';
  ctx.lineWidth = 4;
  ctx.stroke();

  // Sun rays decoration
  ctx.save();
  ctx.translate(centerX, centerY);
  for (let i = 0; i < 8; i++) {
    ctx.rotate(Math.PI / 4);
    ctx.beginPath();
    ctx.moveTo(0, 30);
    ctx.lineTo(0, 38);
    ctx.strokeStyle = 'rgba(255, 255, 255, 0.6)';
    ctx.lineWidth = 3;
    ctx.lineCap = 'round';
    ctx.stroke();
  }
  ctx.restore();
}

// Helper to adjust color brightness
function adjustColor(hex: string, amount: number): string {
  const num = parseInt(hex.replace('#', ''), 16);
  const r = Math.min(255, Math.max(0, (num >> 16) + amount));
  const g = Math.min(255, Math.max(0, ((num >> 8) & 0x00FF) + amount));
  const b = Math.min(255, Math.max(0, (num & 0x0000FF) + amount));
  return `#${((r << 16) | (g << 8) | b).toString(16).padStart(6, '0')}`;
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
  rotation.value = 0;
  result.value = null;
  showResult.value = false;

  try {
    const drawResult = await api.draw(state.currentKid.id);
    result.value = drawResult;

    // Refresh kids data to sync points display
    await state.fetchKids();

    // Find the winning sector
    const prizeIndex = state.prizes.findIndex(p => p.id === drawResult.id);
    const sector = sectors.value[prizeIndex];
    const targetAngle = sector.startAngle + (sector.endAngle - sector.startAngle) / 2;

    // Calculate rotation with more dramatic effect
    const extraRotations = Math.floor(Math.random() * 3 + 5) * 360;
    const finalRotation = rotation.value + extraRotations + (360 - targetAngle);
    rotation.value = finalRotation;

    setTimeout(() => {
      showResult.value = true;
      spinning.value = false;
    }, 4500);
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

watch(() => state.prizes, () => {
  drawWheel();
}, { deep: true });
</script>

<template>
  <div class="card p-5">
    <h2 class="text-lg font-bold text-deep-blue mb-4 text-center flex items-center justify-center gap-2">
      <span>🎁</span>
      <span>抽奖转盘</span>
      <span>🎁</span>
    </h2>

    <div class="flex flex-col items-center gap-4">
      <!-- Wheel Container -->
      <div class="relative">
        <!-- Pointer with cloud decoration -->
        <div class="absolute top-0 left-1/2 -translate-x-1/2 -translate-y-2 z-20">
          <div class="relative">
            <div class="w-0 h-0 border-l-6 border-r-6 border-t-10 border-l-transparent border-r-transparent"
              style="border-top: 16px solid var(--deep-blue);"></div>
            <div class="absolute -top-1 left-1/2 -translate-x-1/2 text-xs">☁️</div>
          </div>
        </div>

        <!-- Wheel with bounce animation -->
        <div
          class="rounded-full overflow-hidden transition-transform origin-center shadow-xl"
          :style="{
            transform: `rotate(${rotation}deg)`,
            transitionDuration: (spinning && rotation !== 0) ? '4.5s' : '0ms',
            transitionTimingFunction: 'cubic-bezier(0.34, 1.56, 0.64, 1)',
            boxShadow: '0 10px 40px rgba(25, 118, 210, 0.3)',
          }"
        >
          <canvas ref="wheelCanvas" width="420" height="420" />
        </div>

        <!-- Center Button with sun -->
        <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 z-10">
          <button
            @click="spin"
            :disabled="spinning"
            class="w-16 h-16 rounded-full font-bold shadow-lg flex items-center justify-center text-white transition-all duration-200"
            :class="spinning ? 'opacity-70 cursor-not-allowed scale-95' : 'hover:scale-105 active:scale-90'"
            style="background: linear-gradient(135deg, #FFD54F 0%, #FFB300 100%); box-shadow: 0 4px 20px rgba(255, 213, 79, 0.5);"
          >
            <span v-if="!spinning" class="text-deep-blue font-black text-sm">抽奖</span>
            <span v-else class="text-deep-blue text-lg animate-pulse">🌞</span>
          </button>
        </div>
      </div>

      <!-- Prize list -->
      <div class="w-full bg-gradient-to-b from-sky/10 to-ocean/5 rounded-2xl p-4">
        <h3 class="text-sm font-semibold text-gray-500 mb-3 text-center">🎯 奖项列表</h3>
        <div class="grid grid-cols-2 gap-2">
          <div
            v-for="prize in state.prizes"
            :key="prize.id"
            class="flex items-center gap-2 bg-sky/10 rounded-xl px-3 py-2"
          >
            <span class="text-xl">{{ prize.icon }}</span>
            <span class="flex-1 truncate text-sm font-medium text-gray-700">{{ prize.name }}</span>
            <span class="text-xs font-semibold text-ocean">{{ prize.probability }}%</span>
          </div>
        </div>
        <div class="mt-3 pt-3 border-t border-sky/20 text-center text-sm text-gray-500">
          🌟 每次抽奖消耗 2 积分
        </div>
      </div>
    </div>

    <p class="text-center text-sm text-gray-500 mt-4">
      {{ state.currentKid ? `${state.currentKid.name} 可以抽奖啦` : '请先选择小朋友' }}
    </p>

    <!-- Result Modal with cloud decorations -->
    <Teleport to="body">
      <div v-if="showResult && result" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50 backdrop-blur-sm">
        <div class="card-elevated p-8 text-center max-w-xs mx-4 relative overflow-hidden">
          <!-- Decorative clouds -->
          <div class="absolute top-2 left-4 text-3xl opacity-30 animate-float">☁️</div>
          <div class="absolute top-4 right-6 text-2xl opacity-20 animate-float" style="animation-delay: 0.3s;">☁️</div>
          <div class="absolute bottom-4 left-6 text-xl opacity-25 animate-float" style="animation-delay: 0.6s;">✨</div>

          <div class="relative">
            <!-- Celebration emoji -->
            <div class="text-6xl mb-4 animate-bounce-in">🎉</div>

            <!-- Prize icon -->
            <div class="text-7xl mb-4">{{ result.icon }}</div>

            <!-- Result text -->
            <h3 class="text-xl font-black text-deep-blue mb-2">恭喜获得</h3>
            <p class="text-2xl font-bold text-ocean mb-6">{{ result.name }}</p>

            <!-- Decorative stars -->
            <div class="flex justify-center gap-2 mb-6">
              <span class="text-2xl animate-float">⭐</span>
              <span class="text-2xl animate-float" style="animation-delay: 0.1s;">✨</span>
              <span class="text-2xl animate-float" style="animation-delay: 0.2s;">🌟</span>
            </div>

            <button
              @click="closeResult"
              class="btn btn-primary px-8 py-3 text-lg"
            >
              太棒了！🌟
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
