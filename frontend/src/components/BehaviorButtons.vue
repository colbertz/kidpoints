<script setup lang="ts">
import { computed } from 'vue';
import { useStore } from '../stores/app';

const state = useStore();

const addBehaviors = computed(() => state.behaviors.filter(b => b.type === 'add'));
const subBehaviors = computed(() => state.behaviors.filter(b => b.type === 'sub'));

async function handleAdd(behaviorId: number) {
  await state.addPoints(behaviorId);
}

async function handleSub(behaviorId: number) {
  await state.subPoints(behaviorId);
}
</script>

<template>
  <div class="p-4 space-y-6">
    <div v-if="addBehaviors.length > 0">
      <h3 class="text-sm font-medium text-green-600 mb-3">加分行为</h3>
      <div class="flex flex-wrap gap-2">
        <button
          v-for="behavior in addBehaviors"
          :key="behavior.id"
          @click="handleAdd(behavior.id)"
          class="px-4 py-2 rounded-lg bg-green-500 text-white hover:bg-green-600 transition-colors"
        >
          {{ behavior.name }}+{{ behavior.points }}
        </button>
      </div>
    </div>

    <div v-if="subBehaviors.length > 0">
      <h3 class="text-sm font-medium text-red-600 mb-3">减分行为</h3>
      <div class="flex flex-wrap gap-2">
        <button
          v-for="behavior in subBehaviors"
          :key="behavior.id"
          @click="handleSub(behavior.id)"
          class="px-4 py-2 rounded-lg bg-red-500 text-white hover:bg-red-600 transition-colors"
        >
          {{ behavior.name }}-{{ behavior.points }}
        </button>
      </div>
    </div>

    <div v-if="state.behaviors.length === 0" class="text-gray-400 text-center py-4">
      暂无行为配置
    </div>
  </div>
</template>
