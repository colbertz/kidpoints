<script setup lang="ts">
import { useStore } from '../stores/app';

const state = useStore();
</script>

<template>
  <div class="text-center py-8">
    <div v-if="state.currentKid" class="inline-flex flex-col items-center">
      <span class="text-6xl mb-2">{{ state.currentKid.avatar }}</span>
      <span class="text-xl text-gray-600 mb-4">{{ state.currentKid.name }}</span>
      <div class="text-7xl font-bold text-purple-600">
        {{ state.currentKid.points }}
      </div>
      <span class="text-gray-400 mt-2">积分</span>
    </div>
    <div v-else class="text-gray-400 py-8">
      暂无小朋友
    </div>
  </div>
</template>
