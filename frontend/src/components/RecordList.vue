<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useStore } from '../stores/app';
import type { Record } from '../types';
import * as api from '../api';

const state = useStore();
const records = ref<Record[]>([]);
const loading = ref(false);
const filterKidId = ref<number | null>(null);
const showReverseModal = ref(false);
const reverseRecordId = ref<number | null>(null);
const reverseReason = ref('误操作');

async function fetchRecords() {
  loading.value = true;
  try {
    const result = await api.getRecords(filterKidId.value || undefined);
    records.value = result || [];
  } catch (e) {
    console.error('Failed to fetch records:', e);
    records.value = [];
  } finally {
    loading.value = false;
  }
}

function openReverseModal(id: number) {
  reverseRecordId.value = id;
  reverseReason.value = '误操作';
  showReverseModal.value = true;
}

async function handleReverse() {
  if (!reverseRecordId.value || !reverseReason.value.trim()) return;
  try {
    await api.reverseRecord(reverseRecordId.value, reverseReason.value.trim());
    showReverseModal.value = false;
    await fetchRecords();
    await state.fetchKids();
  } catch (e) {
    alert((e as Error).message);
  }
}

function formatDate(dateStr: string) {
  const date = new Date(dateStr);
  return date.toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  });
}

function getKidName(kidId: number) {
  const kid = state.kids.find(k => k.id === kidId);
  return kid ? `${kid.name}` : `小朋友 ${kidId}`;
}

function getBehaviorEmoji(record: Record) {
  if (record.reversed) return '🚫';
  if (record.points > 0) return '✨';
  return '📝';
}

onMounted(() => {
  fetchRecords();
});
</script>

<template>
  <div class="space-y-4">
    <!-- Header -->
    <div class="card p-4">
      <div class="flex justify-between items-center">
        <h2 class="text-lg font-bold text-deep-blue flex items-center gap-2">
          <span>📋</span>
          <span>行为记录</span>
        </h2>
        <div class="flex gap-2 items-center">
          <select
            v-model="filterKidId"
            @change="fetchRecords"
            class="px-3 py-2 border-2 border-sky/30 rounded-xl text-sm bg-sky/10 focus:border-ocean transition-colors"
          >
            <option :value="null">全部小朋友</option>
            <option v-for="kid in state.kids" :key="kid.id" :value="kid.id">
              {{ kid.avatar }} {{ kid.name }}
            </option>
          </select>
          <button
            @click="fetchRecords"
            class="btn btn-ghost p-2"
            title="刷新"
          >
            🔄
          </button>
        </div>
      </div>
    </div>

    <!-- Loading state -->
    <div v-if="loading" class="card p-8 text-center">
      <div class="text-4xl mb-3 animate-float">⏳</div>
      <p class="text-gray-400">加载中...</p>
    </div>

    <!-- Empty state -->
    <div v-else-if="records.length === 0" class="card p-8 text-center">
      <div class="text-5xl mb-4">📭</div>
      <p class="text-gray-400 font-medium">暂无记录</p>
      <p class="text-gray-400 text-sm mt-1">快去添加积分变化吧</p>
    </div>

    <!-- Records list -->
    <div v-else class="space-y-3">
      <TransitionGroup name="list">
        <div
          v-for="(record, index) in records"
          :key="record.id"
          class="card p-4 transition-all duration-300"
          :class="[
            record.reversed ? 'opacity-50' : '',
            record.points > 0 ? 'border-l-4 border-l-grass' : 'border-l-4 border-l-coral'
          ]"
          :style="{ animationDelay: `${index * 0.05}s` }"
        >
          <div class="flex justify-between items-start">
            <div class="flex-1">
              <!-- Behavior name with emoji -->
              <div class="flex items-center gap-2 mb-1">
                <span class="text-lg">{{ getBehaviorEmoji(record) }}</span>
                <span
                  class="font-semibold"
                  :class="record.reversed ? 'line-through text-gray-400' : 'text-deep-blue'"
                >
                  {{ record.behavior_name }}
                </span>
                <!-- Points badge -->
                <span
                  class="badge text-sm font-bold px-2 py-0.5"
                  :class="record.points > 0 ? 'badge-green' : 'badge-red'"
                >
                  {{ record.points > 0 ? '+' : '' }}{{ record.points }}
                </span>
              </div>

              <!-- Kid name & time -->
              <div class="text-sm text-gray-500 ml-7">
                {{ getKidName(record.kid_id) }} · {{ formatDate(record.created_at) }}
              </div>

              <!-- Reversed info -->
              <div v-if="record.reversed" class="mt-2 ml-7 text-sm text-coral flex items-center gap-1">
                <span>🚫</span>
                <span>已撤销: {{ record.reversed_reason }}</span>
                <span class="text-gray-400">· {{ record.reversed_at }}</span>
              </div>
            </div>

            <!-- Reverse button (only for non-reversed) -->
            <button
              v-if="!record.reversed"
              @click="openReverseModal(record.id)"
              class="text-gray-400 hover:text-coral transition-colors p-2 rounded-lg hover:bg-coral/10"
              title="撤销"
            >
              <span class="text-lg">↩️</span>
            </button>
          </div>
        </div>
      </TransitionGroup>
    </div>

    <!-- Reverse Modal -->
    <Teleport to="body">
      <div v-if="showReverseModal" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50 backdrop-blur-sm">
        <div class="card-elevated p-6 w-80 mx-4 animate-bounce-in">
          <h3 class="text-lg font-bold text-deep-blue mb-4 flex items-center gap-2">
            <span>🔄</span>
            <span>撤销记录</span>
          </h3>
          <div class="mb-4">
            <label class="block text-sm font-semibold text-gray-600 mb-2">撤销原因</label>
            <input
              v-model="reverseReason"
              type="text"
              class="w-full px-4 py-3 rounded-xl border-2 border-sky/30 focus:border-ocean transition-colors"
              autofocus
            />
          </div>
          <div class="flex gap-3 justify-end">
            <button
              @click="showReverseModal = false"
              class="btn btn-ghost px-4 py-2"
            >
              取消
            </button>
            <button
              @click="handleReverse"
              class="btn btn-sub px-4 py-2"
            >
              确认撤销
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<style scoped>
/* List transition animations */
.list-enter-active {
  animation: slideInRight 0.35s ease-out forwards;
}

.list-leave-active {
  animation: slideInRight 0.25s ease-in reverse forwards;
}

.list-move {
  transition: transform 0.35s ease;
}
</style>
