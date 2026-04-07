<script setup lang="ts">
import { ref } from 'vue';
import { useStore } from '../stores/app';
import type { Behavior } from '../types';
import * as api from '../api';

const state = useStore();
const showForm = ref(false);
const editingBehavior = ref<Behavior | null>(null);
const form = ref({
  name: '',
  type: 'add' as 'add' | 'sub',
  points: 1,
});

function openAddForm() {
  editingBehavior.value = null;
  form.value = { name: '', type: 'add', points: 1 };
  showForm.value = true;
}

function openEditForm(behavior: Behavior) {
  editingBehavior.value = behavior;
  form.value = {
    name: behavior.name,
    type: behavior.type,
    points: behavior.points,
  };
  showForm.value = true;
}

async function handleSubmit() {
  try {
    if (editingBehavior.value) {
      await api.updateBehavior(editingBehavior.value.id, form.value);
    } else {
      await api.createBehavior(form.value);
    }
    await state.fetchBehaviors();
    showForm.value = false;
  } catch (e) {
    alert((e as Error).message);
  }
}

async function handleDelete(id: number) {
  if (!confirm('确定要删除这个行为吗？')) return;
  try {
    await api.deleteBehavior(id);
    await state.fetchBehaviors();
  } catch (e) {
    alert((e as Error).message);
  }
}
</script>

<template>
  <div class="p-4">
    <div class="flex justify-between items-center mb-4">
      <h2 class="text-lg font-semibold">行为管理</h2>
      <button
        @click="openAddForm"
        class="px-3 py-1 bg-purple-500 text-white rounded-lg hover:bg-purple-600"
      >
        + 添加行为
      </button>
    </div>

    <!-- Form Modal -->
    <div v-if="showForm" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 w-96">
        <h3 class="text-lg font-semibold mb-4">
          {{ editingBehavior ? '编辑行为' : '添加行为' }}
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
            <label class="block text-sm text-gray-600 mb-1">类型</label>
            <select v-model="form.type" class="w-full px-3 py-2 border rounded-lg">
              <option value="add">加分</option>
              <option value="sub">减分</option>
            </select>
          </div>
          <div>
            <label class="block text-sm text-gray-600 mb-1">积分</label>
            <input
              v-model.number="form.points"
              type="number"
              required
              min="1"
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

    <!-- Behavior List -->
    <div class="space-y-2">
      <div
        v-for="behavior in state.behaviors"
        :key="behavior.id"
        class="flex items-center justify-between p-3 bg-white rounded-lg border"
        :class="behavior.type === 'add' ? 'border-green-200' : 'border-red-200'"
      >
        <div>
          <div class="font-medium">{{ behavior.name }}</div>
          <div class="text-sm" :class="behavior.type === 'add' ? 'text-green-600' : 'text-red-600'">
            {{ behavior.type === 'add' ? '+' : '-' }}{{ behavior.points }} 分
          </div>
        </div>
        <div class="flex gap-2">
          <button
            @click="openEditForm(behavior)"
            class="px-3 py-1 text-sm text-purple-600 hover:bg-purple-50 rounded"
          >
            编辑
          </button>
          <button
            @click="handleDelete(behavior.id)"
            class="px-3 py-1 text-sm text-red-600 hover:bg-red-50 rounded"
          >
            删除
          </button>
        </div>
      </div>
      <div v-if="state.behaviors.length === 0" class="text-center text-gray-400 py-8">
        暂无行为配置
      </div>
    </div>
  </div>
</template>
