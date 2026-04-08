<script setup lang="ts">
import { ref } from 'vue';

const props = defineProps<{
  modelValue: string;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void;
}>();

const showPicker = ref(false);

const emojis = [
  '🎮', '🦄', '🧗', '📺', '🍕', '🧘', '😀', '📝', '🍪', '🚘',
  '🧸', '✨', '💫', '🌈', '☀️', '🌙', '😴', '📚', '🎯', '🏆',
  '🎮', '🎨', '🎵', '🎬', '📝', '✏️', '🧹', '🍎', '🍕', '🍦',
  '🧸', '🎁', '🌺', '🌸', '🌻', '🌼', '🌱', '🦋', '🐱', '🐶',
  '🦁', '🐼', '🐨', '🦊', '🐸', '🐵', '🦄', '🐢', '🐠', '🦜',
  '☕', '🍵', '🧃', '🍰', '🍪', '🍩', '🍫', '🍿', '🧁', '🍭',
  '😀', '😄', '🙂', '😊', '😇', '🤗', '🤔', '😅', '😂', '🤣',
  '😘', '🥰', '😍', '🤩', '😎', '🤓', '⭐', '❤️', '💯', '🔥',
];

function selectEmoji(emoji: string) {
  emit('update:modelValue', emoji);
  showPicker.value = false;
}
</script>

<template>
  <div class="emoji-picker relative">
    <div class="flex gap-2">
      <button
        type="button"
        @click="showPicker = !showPicker"
        class="w-12 h-10 border rounded-lg text-2xl flex items-center justify-center bg-gray-50 hover:bg-gray-100"
      >
        {{ modelValue || '👍' }}
      </button>
      <input
        :value="modelValue"
        @input="emit('update:modelValue', ($event.target as HTMLInputElement).value)"
        type="text"
        placeholder="选择或输入emoji"
        class="flex-1 px-3 py-2 border rounded-lg"
      />
    </div>

    <div
      v-if="showPicker"
      class="absolute top-full left-0 mt-1 w-72 h-64 bg-white border rounded-lg shadow-lg overflow-auto z-50"
      @click.stop
    >
      <div class="grid grid-cols-8 gap-1 p-2">
        <button
          v-for="emoji in emojis"
          :key="emoji"
          type="button"
          @click="selectEmoji(emoji)"
          class="w-8 h-8 text-xl flex items-center justify-center hover:bg-gray-100 rounded"
        >
          {{ emoji }}
        </button>
      </div>
    </div>
  </div>
</template>
