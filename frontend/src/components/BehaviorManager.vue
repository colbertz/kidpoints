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
  <div class="space-y-4">
    <!-- Header -->
    <div class="card p-4">
      <div class="flex justify-between items-center">
        <h2 class="text-lg font-bold text-deep-blue flex items-center gap-2">
          <span>✨</span>
          <span>行为管理</span>
        </h2>
        <button
          @click="openAddForm"
          class="btn btn-primary text-sm px-4 py-2"
        >
          + 添加行为
        </button>
      </div>
    </div>

    <!-- Form Modal -->
    <Teleport to="body">
      <div v-if="showForm" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50 backdrop-blur-sm">
        <div class="card-elevated p-6 w-96 mx-4 animate-bounce-in">
          <h3 class="text-lg font-bold text-deep-blue mb-4 flex items-center gap-2">
            <span>{{ editingBehavior ? '✏️' : '➕' }}</span>
            <span>{{ editingBehavior ? '编辑行为' : '添加行为' }}</span>
          </h3>
          <form @submit.prevent="handleSubmit" class="space-y-4">
            <div>
              <label class="block text-sm font-semibold text-gray-600 mb-2">名称</label>
              <input
                v-model="form.name"
                type="text"
                required
                class="w-full px-4 py-3 rounded-xl border-2 border-sky/30 focus:border-ocean transition-colors"
                placeholder="例如：主动收拾玩具"
              />
            </div>
            <div>
              <label class="block text-sm font-semibold text-gray-600 mb-2">类型</label>
              <div class="grid grid-cols-2 gap-2">
                <button
                  type="button"
                  @click="form.type = 'add'"
                  class="py-3 px-4 rounded-xl font-semibold transition-all flex items-center justify-center gap-2"
                  :class="form.type === 'add'
                    ? 'bg-gradient-to-r from-grass to-mint text-white shadow-md'
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200'"
                >
                  <span>➕</span>
                  <span>加分</span>
                </button>
                <button
                  type="button"
                  @click="form.type = 'sub'"
                  class="py-3 px-4 rounded-xl font-semibold transition-all flex items-center justify-center gap-2"
                  :class="form.type === 'sub'
                    ? 'bg-gradient-to-r from-coral to-red-400 text-white shadow-md'
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200'"
                >
                  <span>➖</span>
                  <span>减分</span>
                </button>
              </div>
            </div>
            <div>
              <label class="block text-sm font-semibold text-gray-600 mb-2">积分</label>
              <input
                v-model.number="form.points"
                type="number"
                required
                min="1"
                class="w-full px-4 py-3 rounded-xl border-2 border-sky/30 focus:border-ocean transition-colors"
              />
            </div>
            <div class="flex gap-3 justify-end pt-2">
              <button
                type="button"
                @click="showForm = false"
                class="btn btn-ghost px-4"
              >
                取消
              </button>
              <button
                type="submit"
                class="btn btn-primary px-6"
              >
                保存
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- Behavior List -->
    <div class="space-y-3">
      <div
        v-for="(behavior, index) in state.behaviors"
        :key="behavior.id"
        class="card p-4 animate-slide-up border-l-4"
        :class="[
          behavior.type === 'add' ? 'border-l-grass' : 'border-l-coral',
          { 'opacity-50': !behavior }
        ]"
        :style="{ animationDelay: `${index * 0.05}s` }"
      >
        <div class="flex items-center justify-between">
          <div class="flex-1">
            <div class="flex items-center gap-2 mb-1">
              <span class="text-xl">
                {{ behavior.type === 'add' ? '✨' : '📝' }}
              </span>
              <span class="font-bold text-deep-blue">{{ behavior.name }}</span>
            </div>
            <div class="flex items-center gap-2 ml-7">
              <span
                class="badge text-sm font-bold"
                :class="behavior.type === 'add' ? 'badge-green' : 'badge-red'"
              >
                {{ behavior.type === 'add' ? '+' : '-' }}{{ behavior.points }} 分
              </span>
            </div>
          </div>
          <div class="flex gap-2">
            <button
              @click="openEditForm(behavior)"
              class="btn btn-ghost p-2 text-sm"
            >
              ✏️
            </button>
            <button
              @click="handleDelete(behavior.id)"
              class="btn btn-ghost p-2 text-sm hover:text-coral"
            >
              🗑️
            </button>
          </div>
        </div>
      </div>

      <!-- Empty state -->
      <div v-if="state.behaviors.length === 0" class="card p-8 text-center">
        <div class="text-5xl mb-4">✨</div>
        <p class="text-gray-400 font-medium">还没有行为配置</p>
        <p class="text-gray-400 text-sm mt-1">点击上方按钮添加第一个行为</p>
      </div>
    </div>
  </div>
</template>
