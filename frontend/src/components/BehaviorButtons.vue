<script setup lang="ts">
import { ref, computed } from 'vue';
import { useStore } from '../stores/app';

const state = useStore();

const addBehaviors = computed(() => state.behaviors.filter(b => b.type === 'add'));
const subBehaviors = computed(() => state.behaviors.filter(b => b.type === 'sub'));

// Floating text animations
const floatingTexts = ref<{ id: number; text: string; x: number; y: number; type: 'add' | 'sub' }[]>([]);
let floatId = 0;

async function handleAdd(behaviorId: number, event: MouseEvent) {
  const behavior = state.behaviors.find(b => b.id === behaviorId);
  if (!behavior) return;

  // Show floating +N text
  showFloatingText(`+${behavior.points}`, event, 'add');

  await state.addPoints(behaviorId);
}

async function handleSub(behaviorId: number, event: MouseEvent) {
  const behavior = state.behaviors.find(b => b.id === behaviorId);
  if (!behavior) return;

  // Show floating -N text
  showFloatingText(`-${behavior.points}`, event, 'sub');

  await state.subPoints(behaviorId);
}

function showFloatingText(text: string, event: MouseEvent, type: 'add' | 'sub') {
  const rect = (event.target as HTMLElement).getBoundingClientRect();
  const id = ++floatId;

  floatingTexts.value.push({
    id,
    text,
    x: rect.left + rect.width / 2,
    y: rect.top,
    type,
  });

  setTimeout(() => {
    floatingTexts.value = floatingTexts.value.filter(t => t.id !== id);
  }, 800);
}
</script>

<template>
  <div class="card p-5">
    <!-- Section: Add behaviors -->
    <div v-if="addBehaviors.length > 0" class="mb-6">
      <h3 class="text-sm font-bold text-gray-500 mb-3 flex items-center gap-2">
        <span class="w-2 h-2 rounded-full bg-gradient-to-r from-grass to-mint"></span>
        加分行为
      </h3>
      <div class="flex flex-wrap gap-3">
        <button
          v-for="(behavior, index) in addBehaviors"
          :key="behavior.id"
          @click="handleAdd(behavior.id, $event)"
          class="btn btn-add text-base font-bold px-5 py-3 rounded-2xl shadow-md hover:shadow-lg"
          :style="{ animationDelay: `${index * 0.05}s` }"
        >
          <span>{{ behavior.name }}</span>
          <span class="ml-1 text-white/90">+{{ behavior.points }}</span>
        </button>
      </div>
    </div>

    <!-- Section: Subtract behaviors -->
    <div v-if="subBehaviors.length > 0">
      <h3 class="text-sm font-bold text-gray-500 mb-3 flex items-center gap-2">
        <span class="w-2 h-2 rounded-full bg-gradient-to-r from-coral to-red-400"></span>
        减分行为
      </h3>
      <div class="flex flex-wrap gap-3">
        <button
          v-for="(behavior, index) in subBehaviors"
          :key="behavior.id"
          @click="handleSub(behavior.id, $event)"
          class="btn btn-sub text-base font-bold px-5 py-3 rounded-2xl shadow-md hover:shadow-lg"
          :style="{ animationDelay: `${index * 0.05}s` }"
        >
          <span>{{ behavior.name }}</span>
          <span class="ml-1 text-white/90">-{{ behavior.points }}</span>
        </button>
      </div>
    </div>

    <!-- Empty state -->
    <div v-if="state.behaviors.length === 0" class="text-center py-8">
      <div class="text-4xl mb-3">✨</div>
      <p class="text-gray-400 font-medium">还没有行为配置</p>
      <p class="text-gray-400 text-sm mt-1">去"行为"页面添加加分/减分行为</p>
    </div>

    <!-- Floating text animations -->
    <Teleport to="body">
      <TransitionGroup name="float">
        <div
          v-for="ft in floatingTexts"
          :key="ft.id"
          class="floating-text"
          :class="ft.type === 'add' ? 'text-grass' : 'text-coral'"
          :style="{
            left: `${ft.x}px`,
            top: `${ft.y}px`,
            transform: 'translate(-50%, -50%)',
          }"
        >
          {{ ft.text }}
        </div>
      </TransitionGroup>
    </Teleport>
  </div>
</template>

<style scoped>
.float-enter-active {
  animation: floatUp 0.8s ease-out forwards;
}

.float-leave-active {
  display: none;
}
</style>
