<template>
  <div class="confirm-action">
    <div class="confirm-header">
      <i class="material-symbols" v-if="danger">warning</i>
      <span class="confirm-title">{{ title }}</span>
    </div>
    <div class="confirm-message">{{ message }}</div>
    <div class="confirm-buttons">
      <button class="button button--cancel" @click="cancel">
        {{ cancelText || $t('prompts.confirmCancel') }}
      </button>
      <button class="button" :class="{ 'button--danger': danger }" @click="confirm">
        {{ confirmText || $t('prompts.confirmConfirm') }}
      </button>
    </div>
  </div>
</template>

<script>
import { mutations } from "@/store";

export default {
  name: "ConfirmAction",
  props: {
    title: { type: String, default: "" },
    message: { type: String, default: "" },
    confirmText: { type: String, default: "" },
    cancelText: { type: String, default: "" },
    danger: { type: Boolean, default: false },
  },
  methods: {
    confirm() {
      mutations.closeTopPrompt();
    },
    cancel() {
      mutations.closeTopPrompt();
    },
  },
};
</script>

<style scoped>
.confirm-action {
  padding: 20px;
  text-align: center;
}
.confirm-header {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-bottom: 12px;
}
.confirm-header .material-symbols {
  color: #f44336;
  font-size: 24px;
}
.confirm-title {
  font-weight: 600;
  font-size: 1.1em;
}
.confirm-message {
  margin-bottom: 20px;
  color: var(--textSecondary);
  line-height: 1.5;
}
.confirm-buttons {
  display: flex;
  gap: 12px;
  justify-content: center;
}
.confirm-buttons .button {
  min-width: 100px;
  min-height: 44px;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
}
.button--cancel {
  background: transparent;
  border: 1px solid var(--borderColor, #ccc);
  color: var(--textPrimary);
}
.button--danger {
  background: #f44336;
  color: white;
  border: none;
}
@media (max-width: 768px) {
  .confirm-buttons .button {
    flex: 1;
    min-width: 0;
  }
}
</style>
