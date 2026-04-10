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
    avatarPreview.value = avatarUrl;
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
  <div class="space-y-4">
    <!-- Header -->
    <div class="card p-4">
      <div class="flex justify-between items-center">
        <h2 class="text-lg font-bold text-deep-blue flex items-center gap-2">
          <span>🎈</span>
          <span>小朋友管理</span>
        </h2>
        <button
          @click="openAddForm"
          class="btn btn-primary text-sm px-4 py-2"
        >
          + 添加小朋友
        </button>
      </div>
    </div>

    <!-- Form Modal -->
    <Teleport to="body">
      <div v-if="showForm" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50 backdrop-blur-sm">
        <div class="card-elevated p-6 w-96 mx-4 animate-bounce-in">
          <h3 class="text-lg font-bold text-deep-blue mb-4 flex items-center gap-2">
            <span>{{ editingKid ? '✏️' : '➕' }}</span>
            <span>{{ editingKid ? '编辑小朋友' : '添加小朋友' }}</span>
          </h3>
          <form @submit.prevent="handleSubmit" class="space-y-4">
            <div>
              <label class="block text-sm font-semibold text-gray-600 mb-2">名字</label>
              <input
                v-model="form.name"
                type="text"
                required
                class="w-full px-4 py-3 rounded-xl border-2 border-sky/30 focus:border-ocean transition-colors"
                placeholder="例如：小明"
              />
            </div>
            <div>
              <label class="block text-sm font-semibold text-gray-600 mb-2">头像</label>
              <div class="flex items-center gap-3">
                <!-- Avatar preview -->
                <div
                  class="w-20 h-20 rounded-full border-4 border-dashed border-sky/40 flex items-center justify-center overflow-hidden bg-gradient-to-br from-sky/20 to-ocean/10 transition-all duration-200"
                  :class="uploading ? 'opacity-50' : ''"
                >
                  <img
                    v-if="avatarPreview"
                    :src="avatarPreview"
                    class="w-full h-full object-cover"
                  />
                  <span v-else class="text-4xl">{{ form.avatar || '👤' }}</span>
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
                  class="btn btn-ghost px-4 py-2 text-sm"
                >
                  {{ uploading ? '上传中...' : '📷 选择图片' }}
                </button>
              </div>
            </div>
            <div>
              <label class="block text-sm font-semibold text-gray-600 mb-2">初始积分</label>
              <input
                v-model.number="form.points"
                type="number"
                required
                min="0"
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

    <!-- Kid List -->
    <div class="space-y-3">
      <div
        v-for="(kid, index) in state.kids"
        :key="kid.id"
        class="card p-4 animate-slide-up"
        :style="{ animationDelay: `${index * 0.05}s` }"
      >
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-3">
            <!-- Avatar -->
            <div class="relative">
              <div class="w-16 h-16 rounded-full border-4 border-white shadow-md overflow-hidden bg-gradient-to-br from-sky/30 to-ocean/20">
                <img
                  v-if="kid.avatar.startsWith('/uploads')"
                  :src="kid.avatar"
                  class="w-full h-full object-cover"
                />
                <span v-else class="text-4xl block w-full h-full flex items-center justify-center">{{ kid.avatar }}</span>
              </div>
              <!-- Star badge -->
              <div class="absolute -bottom-1 -right-1 w-7 h-7 bg-sunny rounded-full flex items-center justify-center shadow-md border-2 border-white">
                <span class="text-sm">⭐</span>
              </div>
            </div>
            <div>
              <div class="font-bold text-deep-blue text-lg">{{ kid.name }}</div>
              <div class="flex items-center gap-2 mt-1">
                <span class="badge badge-blue font-bold text-sm">
                  {{ kid.points }} 分
                </span>
              </div>
            </div>
          </div>
          <div class="flex gap-2">
            <button
              @click="openEditForm(kid)"
              class="btn btn-ghost p-2 text-lg"
            >
              ✏️
            </button>
            <button
              @click="handleDelete(kid.id)"
              class="btn btn-ghost p-2 text-lg hover:text-coral"
            >
              🗑️
            </button>
          </div>
        </div>
      </div>

      <!-- Empty state -->
      <div v-if="state.kids.length === 0" class="card p-8 text-center">
        <div class="text-5xl mb-4">🎈</div>
        <p class="text-gray-400 font-medium">还没有小朋友</p>
        <p class="text-gray-400 text-sm mt-1">点击上方按钮添加第一个小朋友</p>
      </div>
    </div>
  </div>
</template>
