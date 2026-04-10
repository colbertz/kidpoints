<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';
import { useStore } from '../stores/app';

const state = useStore();
const isAnimating = ref(false);
const displayPoints = ref(0);

// Animate points when they change
// Sync displayPoints when kids data is loaded
watch(() => state.kids, (kids) => {
  if (kids.length > 0 && state.currentKid) {
    displayPoints.value = state.currentKid.points;
  }
}, { immediate: true });

// Animate points when they change
watch(() => state.currentKid?.points, (newPoints, oldPoints) => {
  if (newPoints === undefined || oldPoints === undefined) return;

  if (newPoints !== oldPoints) {
    animatePoints(oldPoints, newPoints);
  }
});

onMounted(() => {
  if (state.currentKid) {
    displayPoints.value = state.currentKid.points;
  }
});

function animatePoints(from: number, to: number) {
  isAnimating.value = true;
  const duration = 600;
  const startTime = performance.now();

  function update(currentTime: number) {
    const elapsed = currentTime - startTime;
    const progress = Math.min(elapsed / duration, 1);

    // Easing: ease-out with bounce
    const eased = 1 - Math.pow(1 - progress, 3);
    displayPoints.value = Math.round(from + (to - from) * eased);

    if (progress < 1) {
      requestAnimationFrame(update);
    } else {
      displayPoints.value = to;
      setTimeout(() => {
        isAnimating.value = false;
      }, 200);
    }
  }

  requestAnimationFrame(update);
}

// Floating particles for celebration
const particles = ref<{ id: number; x: number; y: number; emoji: string }[]>([]);
let particleId = 0;

function createCelebration() {
  const emojis = ['⭐', '✨', '🌟', '💫'];
  for (let i = 0; i < 6; i++) {
    const id = ++particleId;
    particles.value.push({
      id,
      x: Math.random() * 100 - 50,
      y: 0,
      emoji: emojis[Math.floor(Math.random() * emojis.length)],
    });
    setTimeout(() => {
      particles.value = particles.value.filter(p => p.id !== id);
    }, 1000);
  }
}

watch(() => state.currentKid?.points, (newVal, oldVal) => {
  if (newVal !== undefined && oldVal !== undefined && newVal > oldVal) {
    createCelebration();
  }
});
</script>

<template>
  <div class="card-elevated p-6 relative overflow-hidden">
    <!-- Decorative clouds -->
    <div class="absolute top-2 right-4 text-4xl opacity-20 animate-float">☁️</div>
    <div class="absolute top-8 left-2 text-2xl opacity-10 animate-float" style="animation-delay: 0.5s;">☁️</div>

    <!-- Decorative stars -->
    <div class="absolute bottom-3 left-4 text-xl opacity-20">✨</div>
    <div class="absolute top-4 right-16 text-lg opacity-15">⭐</div>

    <div class="relative">
      <div v-if="state.currentKid" class="flex flex-col items-center">
        <!-- Avatar with glow -->
        <div class="relative mb-4">
          <div
            v-if="state.currentKid.avatar.startsWith('/uploads')"
            class="w-24 h-24 rounded-full object-cover border-4 border-white shadow-lg"
            :class="isAnimating ? 'scale-110' : ''"
            style="transition: transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);"
          >
            <img :src="state.currentKid.avatar" class="w-full h-full rounded-full" />
          </div>
          <span
            v-else
            class="text-6xl block drop-shadow-lg transition-transform duration-300"
            :class="isAnimating ? 'scale-110' : ''"
            style="transition: transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);"
          >
            {{ state.currentKid.avatar }}
          </span>
          <!-- Glow ring when animating -->
          <div
            v-if="isAnimating"
            class="absolute inset-0 rounded-full border-4 border-sunny animate-pulse-soft"
          ></div>
        </div>

        <!-- Name -->
        <h2 class="text-lg font-bold text-deep-blue mb-2">
          {{ state.currentKid.name }}
        </h2>

        <!-- Points display -->
        <div class="relative py-6 px-12">
          <!-- Celebration particles -->
          <div
            v-for="particle in particles"
            :key="particle.id"
            class="absolute text-2xl pointer-events-none animate-float"
            :style="{
              left: `calc(50% + ${particle.x}px)`,
              top: `${particle.y}px`,
            }"
          >
            {{ particle.emoji }}
          </div>

          <!-- Points number with gradient -->
          <div
            class="text-7xl font-black transition-transform duration-300"
            :class="isAnimating ? 'scale-110' : ''"
            style="
              background: linear-gradient(135deg, var(--deep-blue) 0%, var(--ocean) 50%, var(--sky) 100%);
              -webkit-background-clip: text;
              -webkit-text-fill-color: transparent;
              background-clip: text;
            "
          >
            {{ displayPoints }}
          </div>

          <!-- Label -->
          <div class="text-center mt-2">
            <span class="text-sm font-semibold text-gray-400 uppercase tracking-wider">积分</span>
          </div>
        </div>

        <!-- Decorative badge -->
        <div class="mt-2 px-4 py-1 bg-gradient-to-r from-sky/20 to-ocean/20 rounded-full">
          <span class="text-sm font-semibold text-deep-blue">☁️ 继续加油哦 ☁️</span>
        </div>
      </div>

      <!-- Empty state -->
      <div v-else class="text-center py-12">
        <div class="text-5xl mb-4">👶</div>
        <p class="text-gray-400 font-medium">还没有小朋友</p>
        <p class="text-gray-400 text-sm mt-1">去"小朋友"页面添加吧</p>
      </div>
    </div>
  </div>
</template>
