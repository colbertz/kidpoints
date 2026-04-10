<script setup lang="ts">
import { ref } from 'vue';
import { useStore } from '../stores/app';

const state = useStore();
const switchingKidId = ref<number | null>(null);

function selectKid(kid: typeof state.currentKid) {
  if (!kid || kid.id === state.currentKid?.id) return;

  switchingKidId.value = kid.id;
  setTimeout(() => {
    state.setCurrentKid(kid);
    switchingKidId.value = null;
  }, 150);
}
</script>

<template>
  <div class="card p-4">
    <div class="flex gap-3 justify-center items-center flex-wrap">
      <button
        v-for="kid in state.kids"
        :key="kid.id"
        @click="selectKid(kid)"
        class="relative flex flex-col items-center p-3 rounded-3xl transition-all duration-300 min-w-[90px]"
        :class="state.currentKid?.id === kid.id
          ? 'bg-gradient-to-b from-sky/20 to-ocean/10 scale-105'
          : 'hover:bg-gray-50 opacity-70 hover:opacity-100'"
      >
        <!-- Avatar container with decorative frame -->
        <div
          class="relative w-16 h-16 mb-2 transition-all duration-300"
          :class="[
            switchingKidId === kid.id ? 'scale-90 rotate-6' : '',
            state.currentKid?.id === kid.id ? 'scale-110' : ''
          ]"
        >
          <!-- Glow effect for selected -->
          <div
            v-if="state.currentKid?.id === kid.id"
            class="absolute inset-0 rounded-full bg-gradient-to-r from-sky to-ocean blur-md opacity-50 -z-10"
          ></div>

          <!-- Avatar border ring -->
          <div
            class="absolute inset-0 rounded-full border-4 transition-colors duration-200"
            :class="state.currentKid?.id === kid.id ? 'border-ocean' : 'border-gray-200'"
          ></div>

          <!-- Avatar image or emoji -->
          <div class="w-full h-full rounded-full overflow-hidden bg-white flex items-center justify-center">
            <img
              v-if="kid.avatar.startsWith('/uploads')"
              :src="kid.avatar"
              class="w-full h-full object-cover"
            />
            <span v-else class="text-4xl">{{ kid.avatar }}</span>
          </div>

          <!-- Selected indicator star -->
          <div
            v-if="state.currentKid?.id === kid.id"
            class="absolute -top-1 -right-1 w-6 h-6 bg-sunny rounded-full flex items-center justify-center shadow-md animate-bounce-in"
          >
            <span class="text-sm">⭐</span>
          </div>
        </div>

        <!-- Name -->
        <span
          class="text-sm font-bold transition-colors duration-200"
          :class="state.currentKid?.id === kid.id ? 'text-deep-blue' : 'text-gray-600'"
        >
          {{ kid.name }}
        </span>

        <!-- Points badge -->
        <div
          class="mt-1 px-2 py-0.5 rounded-full text-xs font-semibold"
          :class="state.currentKid?.id === kid.id
            ? 'bg-ocean/20 text-deep-blue'
            : 'bg-gray-100 text-gray-500'"
        >
          {{ kid.points }}分
        </div>
      </button>
    </div>

    <!-- Empty state -->
    <div v-if="state.kids.length === 0" class="text-center py-6">
      <div class="text-4xl mb-2">🎈</div>
      <p class="text-gray-400 text-sm">还没有小朋友，快去添加吧</p>
    </div>
  </div>
</template>
