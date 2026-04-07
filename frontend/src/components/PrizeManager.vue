<script setup lang="ts">
import { ref } from 'vue';
import { useStore } from '../stores/app';
import type { Prize } from '../types';
import * as api from '../api';
import EmojiPicker from './ui/EmojiPicker.vue';

const state = useStore();
const showForm = ref(false);
const editingPrize = ref<Prize | null>(null);
const form = ref({
  name: '',
  icon: '',
  probability: 10,
});

function openAddForm() {
  editingPrize.value = null;
  form.value = { name: '', icon: '', probability: 10 };
  showForm.value = true;
}

function openEditForm(prize: Prize) {
  editingPrize.value = prize;
  form.value = {
    name: prize.name,
    icon: prize.icon,
    probability: prize.probability,
  };
  showForm.value = true;
}

async function handleSubmit() {
  try {
    if (editingPrize.value) {
      await api.updatePrize(editingPrize.value.id, form.value);
    } else {
      await api.createPrize(form.value);
    }
    const prizes = await api.getPrizes();
    state.prizes = prizes;
    showForm.value = false;
  } catch (e) {
    alert((e as Error).message);
  }
}

async function handleDelete(id: number) {
  if (!confirm('确定要删除这个奖项吗？')) return;
  try {
    await api.deletePrize(id);
    const prizes = await api.getPrizes();
    state.prizes = prizes;
  } catch (e) {
    alert((e as Error).message);
  }
}
</script>

<template>
  <div class="p-4">
    <div class="flex justify-between items-center mb-4">
      <h2 class="text-lg font-semibold">奖项管理</h2>
      <button
        @click="openAddForm"
        class="px-3 py-1 bg-purple-500 text-white rounded-lg hover:bg-purple-600"
      >
        + 添加奖项
      </button>
    </div>

    <!-- Form Modal -->
    <div v-if="showForm" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 w-96">
        <h3 class="text-lg font-semibold mb-4">
          {{ editingPrize ? '编辑奖项' : '添加奖项' }}
        </h3>
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <div>
            <label class="block text-sm text-gray-600 mb-1">名称</label>
            <input
              v-model="form.name"
              type="text"
              required
              class="w-full px-3 py-2 border rounded-lg"
            />
          </div>
          <div>
            <label class="block text-sm text-gray-600 mb-1">图标</label>
            <EmojiPicker v-model="form.icon" />
          </div>
          <div>
            <label class="block text-sm text-gray-600 mb-1">概率 (%)</label>
            <input
              v-model.number="form.probability"
              type="number"
              required
              min="1"
              max="100"
              class="w-full px-3 py-2 border rounded-lg"
            />
          </div>
          <div class="flex gap-2 justify-end">
            <button
              type="button"
              @click="showForm = false"
              class="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded-lg"
            >
              取消
            </button>
            <button
              type="submit"
              class="px-4 py-2 bg-purple-500 text-white rounded-lg hover:bg-purple-600"
            >
              保存
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Prize List -->
    <div class="space-y-2">
      <div
        v-for="prize in state.prizes"
        :key="prize.id"
        class="flex items-center justify-between p-3 bg-white rounded-lg border border-gray-200"
      >
        <div class="flex items-center gap-3">
          <span class="text-2xl">{{ prize.icon }}</span>
          <div>
            <div class="font-medium">{{ prize.name }}</div>
            <div class="text-sm text-gray-500">
              概率: {{ prize.probability }}%
            </div>
          </div>
        </div>
        <div class="flex gap-2">
          <button
            @click="openEditForm(prize)"
            class="px-3 py-1 text-sm text-purple-600 hover:bg-purple-50 rounded"
          >
            编辑
          </button>
          <button
            @click="handleDelete(prize.id)"
            class="px-3 py-1 text-sm text-red-600 hover:bg-red-50 rounded"
          >
            删除
          </button>
        </div>
      </div>
      <div v-if="state.prizes.length === 0" class="text-center text-gray-400 py-8">
        暂无奖项配置
      </div>
    </div>
  </div>
</template>
