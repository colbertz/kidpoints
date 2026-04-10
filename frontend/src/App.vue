<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useStore } from './stores/app';
import KidSwitcher from './components/KidSwitcher.vue';
import PointsDisplay from './components/PointsDisplay.vue';
import BehaviorButtons from './components/BehaviorButtons.vue';
import LuckyWheel from './components/LuckyWheel.vue';
import BehaviorManager from './components/BehaviorManager.vue';
import PrizeManager from './components/PrizeManager.vue';
import RecordList from './components/RecordList.vue';
import KidManager from './components/KidManager.vue';
import * as api from './api';

const state = useStore();
const activeTab = ref<'main' | 'kids' | 'behaviors' | 'prizes' | 'records'>('main');

const tabs = [
  { key: 'main', label: '积分', icon: '☁️' },
  { key: 'kids', label: '小朋友', icon: '🎈' },
  { key: 'behaviors', label: '行为', icon: '✨' },
  { key: 'prizes', label: '奖项', icon: '🎁' },
  { key: 'records', label: '记录', icon: '📋' },
] as const;

onMounted(async () => {
  await state.fetchKids();
  await state.fetchBehaviors();
  try {
    state.prizes = await api.getPrizes();
  } catch (e) {
    console.error('Failed to fetch prizes:', e);
  }
});
</script>

<template>
  <div class="min-h-screen relative overflow-hidden">
    <!-- Animated sky background layers -->
    <div class="fixed inset-0 pointer-events-none select-none" aria-hidden="true">
      <!-- Gradient orbs -->
      <div class="absolute top-0 left-1/4 w-96 h-96 rounded-full bg-gradient-to-br from-sky/30 to-transparent blur-3xl animate-pulse-soft"></div>
      <div class="absolute top-1/3 right-0 w-80 h-80 rounded-full bg-gradient-to-bl from-sunny/20 to-transparent blur-3xl animate-pulse-soft" style="animation-delay: 1s;"></div>
      <div class="absolute bottom-1/4 left-0 w-72 h-72 rounded-full bg-gradient-to-tr from-mint/25 to-transparent blur-3xl animate-pulse-soft" style="animation-delay: 2s;"></div>

      <!-- Floating clouds - large background clouds -->
      <div class="absolute top-16 -left-20 text-8xl opacity-20 drift">☁️</div>
      <div class="absolute top-32 right-10 text-6xl opacity-15 drift-slow">☁️</div>
      <div class="absolute top-64 left-20 text-5xl opacity-10 drift" style="animation-delay: -3s;">☁️</div>
      <div class="absolute bottom-40 right-1/4 text-7xl opacity-15 drift-slow" style="animation-delay: -6s;">☁️</div>

      <!-- Floating stars/sparkles -->
      <div class="absolute top-20 left-1/3 text-xl shimmer-soft" style="animation-delay: 0s;">✨</div>
      <div class="absolute top-40 right-1/4 text-lg shimmer-soft" style="animation-delay: 1.5s;">⭐</div>
      <div class="absolute top-56 left-1/4 text-base shimmer-soft" style="animation-delay: 0.8s;">🌟</div>
      <div class="absolute bottom-1/3 right-1/3 text-xl shimmer-soft" style="animation-delay: 2.2s;">✨</div>
      <div class="absolute top-48 left-1/2 text-lg shimmer-soft" style="animation-delay: 1.2s;">💫</div>
      <div class="absolute bottom-1/2 left-1/3 text-base shimmer-soft" style="animation-delay: 0.5s;">⭐</div>

      <!-- Small decorative elements -->
      <div class="absolute top-1/2 right-8 text-2xl opacity-10 drift" style="animation-delay: -2s;">🌸</div>
      <div class="absolute bottom-1/3 left-8 text-xl opacity-12 drift-slow">🦋</div>
      <div class="absolute top-1/4 left-12 text-lg opacity-10 drift" style="animation-delay: -5s;">🌺</div>
    </div>

    <!-- Header -->
    <header class="glass sticky top-0 z-40 border-b border-white/30">
      <div class="text-center py-3 px-4 relative">
        <!-- Decorative elements on header -->
        <div class="absolute left-4 top-1/2 -translate-y-1/2 text-lg opacity-40">🌈</div>
        <div class="absolute right-4 top-1/2 -translate-y-1/2 text-lg opacity-40">🌈</div>
        <h1 class="text-xl font-black text-deep-blue tracking-wide drop-shadow-sm">
          ⭐ 积分小助手 ⭐
        </h1>
      </div>
    </header>

    <!-- Main Content -->
    <main class="max-w-2xl mx-auto px-4 py-4 relative z-10">
      <Transition name="page" mode="out-in">
        <!-- Main Tab: Points Display -->
        <template v-if="activeTab === 'main'">
          <div class="space-y-4">
            <KidSwitcher />
            <PointsDisplay />
            <BehaviorButtons />
            <LuckyWheel />
          </div>
        </template>

        <!-- Kids Tab -->
        <template v-else-if="activeTab === 'kids'">
          <KidManager />
        </template>

        <!-- Behaviors Tab -->
        <template v-else-if="activeTab === 'behaviors'">
          <BehaviorManager />
        </template>

        <!-- Prizes Tab -->
        <template v-else-if="activeTab === 'prizes'">
          <PrizeManager />
        </template>

        <!-- Records Tab -->
        <template v-else-if="activeTab === 'records'">
          <RecordList />
        </template>
      </Transition>
    </main>

    <!-- Bottom Navigation -->
    <nav class="glass fixed bottom-0 left-0 right-0 border-t border-white/40 z-50 safe-area-bottom">
      <div class="max-w-2xl mx-auto flex justify-around items-center py-2">
        <button
          v-for="tab in tabs"
          :key="tab.key"
          @click="activeTab = tab.key"
          class="flex flex-col items-center justify-center py-2 px-4 rounded-2xl transition-all duration-200 min-w-[60px] relative"
          :class="activeTab === tab.key
            ? 'text-deep-blue scale-105'
            : 'text-gray-500 hover:text-sky active:scale-95'"
        >
          <span class="text-2xl transition-transform duration-200" :class="activeTab === tab.key ? 'animate-bounce-in' : ''">
            {{ tab.icon }}
          </span>
          <span class="text-xs font-semibold mt-1">{{ tab.label }}</span>
          <!-- Active indicator -->
          <div
            v-if="activeTab === tab.key"
            class="absolute -bottom-1 w-8 h-1 bg-gradient-to-r from-sky to-ocean rounded-full shadow-lg shadow-ocean/30"
          ></div>
        </button>
      </div>
    </nav>
  </div>
</template>

<style scoped>
.safe-area-bottom {
  padding-bottom: env(safe-area-inset-bottom, 0);
}

/* Ensure content doesn't get hidden behind bottom nav on notched devices */
@media (max-width: 767px) {
  main {
    padding-bottom: calc(90px + env(safe-area-inset-bottom, 0));
  }
}
</style>
