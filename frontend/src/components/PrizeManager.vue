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
  <div class="space-y-4">
    <!-- Header -->
    <div class="card p-4">
      <div class="flex justify-between items-center">
        <h2 class="text-lg font-bold text-deep-blue flex items-center gap-2">
          <span>🎁</span>
          <span>奖项管理</span>
        </h2>
        <button
          @click="openAddForm"
          class="btn btn-primary text-sm px-4 py-2"
        >
          + 添加奖项
        </button>
      </div>
    </div>

    <!-- Form Modal -->
    <Teleport to="body">
      <div v-if="showForm" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50 backdrop-blur-sm">
        <div class="card-elevated p-6 w-96 mx-4 animate-bounce-in">
          <h3 class="text-lg font-bold text-deep-blue mb-4 flex items-center gap-2">
            <span>{{ editingPrize ? '✏️' : '➕' }}</span>
            <span>{{ editingPrize ? '编辑奖项' : '添加奖项' }}</span>
          </h3>
          <form @submit.prevent="handleSubmit" class="space-y-4">
            <div>
              <label class="block text-sm font-semibold text-gray-600 mb-2">名称</label>
              <input
                v-model="form.name"
                type="text"
                required
                class="w-full px-4 py-3 rounded-xl border-2 border-sky/30 focus:border-ocean transition-colors"
                placeholder="例如：一根棒棒糖"
              />
            </div>
            <div>
              <label class="block text-sm font-semibold text-gray-600 mb-2">图标</label>
              <EmojiPicker v-model="form.icon" />
            </div>
            <div>
              <label class="block text-sm font-semibold text-gray-600 mb-2">概率 (%)</label>
              <input
                v-model.number="form.probability"
                type="number"
                required
                min="1"
                max="100"
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

    <!-- Prize List -->
    <div class="space-y-3">
      <div
        v-for="(prize, index) in state.prizes"
        :key="prize.id"
        class="card p-4 animate-slide-up"
        :style="{ animationDelay: `${index * 0.05}s` }"
      >
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-3 flex-1">
            <!-- Prize icon -->
            <div class="w-14 h-14 rounded-2xl bg-gradient-to-br from-sunny/30 to-ocean/20 flex items-center justify-center text-3xl">
              {{ prize.icon || '🎁' }}
            </div>
            <div class="flex-1">
              <div class="font-bold text-deep-blue">{{ prize.name }}</div>
              <div class="flex items-center gap-2 mt-1">
                <div class="flex-1 h-2 bg-sky/20 rounded-full overflow-hidden">
                  <div
                    class="h-full bg-gradient-to-r from-ocean to-sky rounded-full transition-all duration-500"
                    :style="{ width: `${prize.probability}%` }"
                  ></div>
                </div>
                <span class="text-sm font-semibold text-ocean">{{ prize.probability }}%</span>
              </div>
            </div>
          </div>
          <div class="flex gap-2 ml-4">
            <button
              @click="openEditForm(prize)"
              class="btn btn-ghost p-2 text-sm"
            >
              ✏️
            </button>
            <button
              @click="handleDelete(prize.id)"
              class="btn btn-ghost p-2 text-sm hover:text-coral"
            >
              🗑️
            </button>
          </div>
        </div>
      </div>

      <!-- Empty state -->
      <div v-if="state.prizes.length === 0" class="card p-8 text-center">
        <div class="text-5xl mb-4">🎁</div>
        <p class="text-gray-400 font-medium">还没有奖项配置</p>
        <p class="text-gray-400 text-sm mt-1">点击上方按钮添加第一个奖项</p>
      </div>
    </div>
  </div>
</template>
