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
import * as api from './api';

const state = useStore();
const activeTab = ref<'main' | 'behaviors' | 'prizes' | 'records'>('main');

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
  <div class="min-h-screen bg-gray-50">
    <header class="bg-white shadow-sm">
      <h1 class="text-2xl font-bold text-center py-4 text-purple-600">
        积分管理
      </h1>
      <!-- Tab Navigation -->
      <nav class="flex border-t">
        <button
          v-for="tab in ['main', 'behaviors', 'prizes', 'records'] as const"
          :key="tab"
          @click="activeTab = tab"
          class="flex-1 py-3 text-center text-sm transition-colors"
          :class="activeTab === tab ? 'text-purple-600 border-b-2 border-purple-600' : 'text-gray-500 hover:text-gray-700'"
        >
          {{ tab === 'main' ? '积分' : tab === 'behaviors' ? '行为' : tab === 'prizes' ? '奖项' : '记录' }}
        </button>
      </nav>
    </header>

    <main class="max-w-2xl mx-auto">
      <!-- Main Tab: Points Display -->
      <template v-if="activeTab === 'main'">
        <KidSwitcher />
        <PointsDisplay />
        <BehaviorButtons />
        <LuckyWheel />
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
    </main>
  </div>
</template>
