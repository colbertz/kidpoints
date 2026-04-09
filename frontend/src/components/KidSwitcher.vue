<script setup lang="ts">
import { useStore } from '../stores/app';

const state = useStore();

function selectKid(kid: typeof state.currentKid) {
  if (kid) {
    state.setCurrentKid(kid);
  }
}
</script>

<template>
  <div class="flex gap-2 p-4">
    <button
      v-for="kid in state.kids"
      :key="kid.id"
      @click="selectKid(kid)"
      class="px-4 py-2 rounded-lg border transition-all flex items-center"
      :class="state.currentKid?.id === kid.id
        ? 'border-purple-500 bg-purple-50 text-purple-700'
        : 'border-gray-200 bg-white text-gray-600 hover:border-purple-300'"
    >
      <img
        v-if="kid.avatar.startsWith('/uploads')"
        :src="kid.avatar"
        class="w-8 h-8 rounded-full mr-2 object-cover"
      />
      <span v-else class="text-2xl mr-2">{{ kid.avatar }}</span>
      <span class="font-medium">{{ kid.name }}</span>
    </button>
  </div>
</template>
