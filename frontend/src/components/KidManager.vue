<script setup lang="ts">
import { ref } from 'vue';
import { useStore } from '../stores/app';
import type { Kid } from '../types';
import * as api from '../api';

const state = useStore();
const showForm = ref(false);
const editingKid = ref<Kid | null>(null);
const form = ref({
  name: '',
  avatar: '',
  points: 10,
});
const avatarPreview = ref('');
const uploading = ref(false);
const fileInput = ref<HTMLInputElement | null>(null);

function openAddForm() {
  editingKid.value = null;
  form.value = { name: '', avatar: '', points: 10 };
  avatarPreview.value = '';
  showForm.value = true;
}

function openEditForm(kid: Kid) {
  editingKid.value = kid;
  form.value = {
    name: kid.name,
    avatar: kid.avatar,
    points: kid.points,
  };
  avatarPreview.value = kid.avatar.startsWith('/uploads') ? kid.avatar : '';
  showForm.value = true;
}

function triggerFileInput() {
  fileInput.value?.click();
}

async function handleAvatarChange(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file) return;

  uploading.value = true;
  try {
    const avatarUrl = await api.uploadAvatar(file);
    form.value.avatar = avatarUrl;
    avatarPreview.value = avatarUrl;  // Use directly, vite proxy handles /uploads
  } catch (e) {
    alert((e as Error).message);
  } finally {
    uploading.value = false;
  }
}

async function handleSubmit() {
  try {
    if (editingKid.value) {
      await api.updateKid(editingKid.value.id, form.value);
    } else {
      await api.createKid(form.value);
    }
    await state.fetchKids();
    showForm.value = false;
  } catch (e) {
    alert((e as Error).message);
  }
}

async function handleDelete(id: number) {
  if (!confirm('确定要删除这个小朋友吗？')) return;
  try {
    await api.deleteKid(id);
    await state.fetchKids();
  } catch (e) {
    alert((e as Error).message);
  }
}
</script>

<template>
  <div class="p-4">
    <div class="flex justify-between items-center mb-4">
      <h2 class="text-lg font-semibold">小朋友管理</h2>
      <button
        @click="openAddForm"
        class="px-3 py-1 bg-purple-500 text-white rounded-lg hover:bg-purple-600"
      >
        + 添加小朋友
      </button>
    </div>

    <!-- Form Modal -->
    <div v-if="showForm" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 w-96">
        <h3 class="text-lg font-semibold mb-4">
          {{ editingKid ? '编辑小朋友' : '添加小朋友' }}
        </h3>
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <div>
            <label class="block text-sm text-gray-600 mb-1">名字</label>
            <input
              v-model="form.name"
              type="text"
              required
              class="w-full px-3 py-2 border rounded-lg"
            />
          </div>
          <div>
            <label class="block text-sm text-gray-600 mb-1">头像</label>
            <div class="flex items-center gap-3">
              <div
                class="w-16 h-16 rounded-full border-2 border-dashed border-gray-300 flex items-center justify-center overflow-hidden bg-gray-50"
              >
                <img
                  v-if="avatarPreview"
                  :src="avatarPreview"
                  class="w-full h-full object-cover"
                />
                <span v-else class="text-3xl text-gray-400">{{ form.avatar || '👤' }}</span>
              </div>
              <input
                ref="fileInput"
                type="file"
                accept="image/*"
                @change="handleAvatarChange"
                class="hidden"
              />
              <button
                type="button"
                @click="triggerFileInput"
                :disabled="uploading"
                class="px-3 py-2 text-sm bg-gray-100 hover:bg-gray-200 rounded-lg disabled:opacity-50"
              >
                {{ uploading ? '上传中...' : '选择图片' }}
              </button>
            </div>
          </div>
          <div>
            <label class="block text-sm text-gray-600 mb-1">积分</label>
            <input
              v-model.number="form.points"
              type="number"
              required
              min="0"
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

    <!-- Kid List -->
    <div class="space-y-2">
      <div
        v-for="kid in state.kids"
        :key="kid.id"
        class="flex items-center justify-between p-3 bg-white rounded-lg border border-gray-200"
      >
        <div class="flex items-center gap-3">
          <div class="w-12 h-12 rounded-full overflow-hidden bg-gray-100 flex items-center justify-center">
            <img
              v-if="kid.avatar.startsWith('/uploads')"
              :src="kid.avatar"
              class="w-full h-full object-cover"
            />
            <span v-else class="text-2xl">{{ kid.avatar }}</span>
          </div>
          <div>
            <div class="font-medium">{{ kid.name }}</div>
            <div class="text-sm text-purple-600">{{ kid.points }} 分</div>
          </div>
        </div>
        <div class="flex gap-2">
          <button
            @click="openEditForm(kid)"
            class="px-3 py-1 text-sm text-purple-600 hover:bg-purple-50 rounded"
          >
            编辑
          </button>
          <button
            @click="handleDelete(kid.id)"
            class="px-3 py-1 text-sm text-red-600 hover:bg-red-50 rounded"
          >
            删除
          </button>
        </div>
      </div>
      <div v-if="state.kids.length === 0" class="text-center text-gray-400 py-8">
        暂无小朋友
      </div>
    </div>
  </div>
</template>
