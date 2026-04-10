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
  <div class="min-h-screen bg-gradient-sky">
    <!-- Header -->
    <header class="bg-white/80 backdrop-blur-sm sticky top-0 z-40 border-b border-sky/20">
      <div class="text-center py-3 px-4">
        <h1 class="text-xl font-black text-deep-blue tracking-wide">
          ⭐ 积分小助手 ⭐
        </h1>
      </div>
    </header>

    <!-- Main Content -->
    <main class="max-w-2xl mx-auto px-4 py-4">
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
    <nav class="fixed bottom-0 left-0 right-0 bg-white/95 backdrop-blur-md border-t border-sky/20 z-50 safe-area-bottom">
      <div class="max-w-2xl mx-auto flex justify-around items-center py-2">
        <button
          v-for="tab in tabs"
          :key="tab.key"
          @click="activeTab = tab.key"
          class="flex flex-col items-center justify-center py-2 px-4 rounded-2xl transition-all duration-200 min-w-[60px]"
          :class="activeTab === tab.key
            ? 'text-deep-blue scale-105'
            : 'text-gray-400 hover:text-sky active:scale-95'"
        >
          <span class="text-2xl transition-transform duration-200" :class="activeTab === tab.key ? 'animate-bounce-in' : ''">
            {{ tab.icon }}
          </span>
          <span class="text-xs font-semibold mt-1">{{ tab.label }}</span>
          <!-- Active indicator -->
          <div
            v-if="activeTab === tab.key"
            class="absolute -bottom-1 w-8 h-1 bg-gradient-to-r from-sky to-ocean rounded-full"
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
    padding-bottom: calc(80px + env(safe-area-inset-bottom, 0));
  }
}
</style>
