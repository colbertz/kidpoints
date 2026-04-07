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
  return kid ? `${kid.avatar} ${kid.name}` : `小朋友 ${kidId}`;
}

onMounted(() => {
  fetchRecords();
});
</script>

<template>
  <div class="p-4">
    <div class="flex justify-between items-center mb-4">
      <h2 class="text-lg font-semibold">行为记录</h2>
      <div class="flex gap-2 items-center">
        <select
          v-model="filterKidId"
          @change="fetchRecords"
          class="px-3 py-1 border rounded-lg text-sm"
        >
          <option :value="null">全部小朋友</option>
          <option v-for="kid in state.kids" :key="kid.id" :value="kid.id">
            {{ kid.avatar }} {{ kid.name }}
          </option>
        </select>
        <button
          @click="fetchRecords"
          class="px-3 py-1 text-purple-600 hover:bg-purple-50 rounded-lg"
        >
          刷新
        </button>
      </div>
    </div>

    <div v-if="loading" class="text-center py-8 text-gray-400">
      加载中...
    </div>

    <div v-else-if="records.length === 0" class="text-center text-gray-400 py-8">
      暂无记录
    </div>

    <div v-else class="space-y-2">
      <div
        v-for="record in records"
        :key="record.id"
        class="p-3 bg-white rounded-lg border"
        :class="record.reversed ? 'border-gray-200 opacity-50' : record.points > 0 ? 'border-green-200' : 'border-red-200'"
      >
        <div class="flex justify-between items-start">
          <div>
            <div class="flex items-center gap-2">
              <span class="font-medium">{{ record.behavior_name }}</span>
              <span
                class="text-sm px-2 py-0.5 rounded"
                :class="record.points > 0 ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'"
              >
                {{ record.points > 0 ? '+' : '' }}{{ record.points }}
              </span>
            </div>
            <div class="text-sm text-gray-500 mt-1">
              {{ getKidName(record.kid_id) }} · {{ formatDate(record.created_at) }}
            </div>
            <div v-if="record.reversed" class="text-sm text-red-500 mt-1">
              已撤销: {{ record.reversed_reason }} · {{ record.reversed_at }}
            </div>
          </div>
          <button
            v-if="!record.reversed"
            @click="openReverseModal(record.id)"
            class="px-2 py-1 text-xs text-gray-500 hover:bg-gray-100 rounded"
          >
            撤销
          </button>
        </div>
      </div>
    </div>

    <!-- Reverse Modal -->
    <div v-if="showReverseModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 w-80">
        <h3 class="text-lg font-semibold mb-4">撤销记录</h3>
        <div class="mb-4">
          <label class="block text-sm text-gray-600 mb-1">撤销原因</label>
          <input
            v-model="reverseReason"
            type="text"
            class="w-full px-3 py-2 border rounded-lg"
            autofocus
          />
        </div>
        <div class="flex gap-2 justify-end">
          <button
            @click="showReverseModal = false"
            class="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded-lg"
          >
            取消
          </button>
          <button
            @click="handleReverse"
            class="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600"
          >
            确认撤销
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
