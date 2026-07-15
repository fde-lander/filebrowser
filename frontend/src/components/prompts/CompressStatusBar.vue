<template>
  <div class="compress-status-bar" v-if="status">
    <div class="status-info" @click="openDetail">
      <span class="status-icon">{{ statusIcon }}</span>
      <span class="status-text">{{ statusText }}</span>
      <span class="status-file" v-if="status.currentFile">{{ status.currentFile }}</span>
    </div>
    <div class="status-buttons">
      <button v-if="status.status === 'running'" class="status-btn" @click="handlePause" :title="$t('compress.pause')">
        <i class="material-symbols">pause</i>
      </button>
      <button v-if="status.status === 'paused'" class="status-btn" @click="handleResume" :title="$t('compress.resume')">
        <i class="material-symbols">play_arrow</i>
      </button>
      <button v-if="status.queueLength > 0 && (status.status === 'running' || status.status === 'paused')" class="status-btn" @click="handleSkip" :title="$t('compress.skipBatch')">
        <i class="material-symbols">skip_next</i>
      </button>
      <button v-if="status.status === 'running' || status.status === 'paused'" class="status-btn" @click="handleCancel" :title="$t('compress.cancelAll')">
        <i class="material-symbols">close</i>
      </button>
    </div>
  </div>
</template>

<script>
import { mutations } from "@/store";
import { notify } from "@/notify";
import { pollStatus, pauseCompress, resumeCompress, cancelCompress, skipBatch } from "@/api/compress";

export default {
  name: "CompressStatusBar",
  data() {
    return {
      status: null,
      pollTimer: null,
      dismissTimer: null,
    };
  },
  computed: {
    statusIcon() {
      if (this.status?.status === "running") return "🔄";
      if (this.status?.status === "paused") return "⏸";
      if (this.status?.status === "completed") return "✅";
      if (this.status?.status === "cancelled") return "✕";
      return "";
    },
    statusText() {
      if (!this.status) return "";
      const totalProc = this.status.totalProcessed || 0;
      const totalFiles = this.status.totalFiles || 0;
      const percent = totalFiles > 0 ? Math.round((totalProc / totalFiles) * 100) : 0;
      const batchIdx = this.status.currentBatchIndex || 0;
      const batchCnt = this.status.batchCount || 0;
      let text = `${totalProc}/${totalFiles} (${percent}%)`;
      if (batchCnt > 1) {
        text += ` | ${this.$t('compress.batchInfo', { current: batchIdx, total: batchCnt })}`;
      }
      if (this.status.status === "paused" && this.status.pausedAt) {
        const elapsed = Math.floor((Date.now() - new Date(this.status.pausedAt).getTime()) / 60000);
        text += ` | ${this.$t('compress.pauseTimeoutCountdown', { min: elapsed })}`;
      }
      return text;
    },
  },
  mounted() {
    this.startPolling();
  },
  beforeUnmount() {
    this.stopPolling();
  },
  methods: {
    startPolling() {
      this.fetchStatus();
      this.pollTimer = setInterval(() => this.fetchStatus(), 3000);
    },
    stopPolling() {
      if (this.pollTimer) {
        clearInterval(this.pollTimer);
        this.pollTimer = null;
      }
      if (this.dismissTimer) {
        clearTimeout(this.dismissTimer);
        this.dismissTimer = null;
      }
    },
    async fetchStatus() {
      try {
        const s = await pollStatus();
        this.status = s;
        if (s.status === "idle" || s.status === "completed" || s.status === "cancelled") {
          if (!this.dismissTimer) {
            this.dismissTimer = setTimeout(() => {
              mutations.closeTopPrompt();
            }, 5000);
          }
        }
      } catch (err) {
        console.error("CompressStatusBar poll failed:", err);
      }
    },
    openDetail() {
      mutations.showPrompt({ name: "compressImages", props: { items: [] } });
    },
    async handlePause() {
      try {
        await pauseCompress();
      } catch (err) {
        console.error("Pause failed:", err);
        notify.showError(this.$t("prompts.compressControlError"));
      }
    },
    async handleResume() {
      try {
        await resumeCompress();
      } catch (err) {
        console.error("Resume failed:", err);
        notify.showError(this.$t("prompts.compressControlError"));
      }
    },
    handleSkip() {
      const remaining = (this.status?.total || 0) - (this.status?.processed || 0);
      mutations.showPrompt({
        name: "confirmAction",
        props: {
          title: this.$t("prompts.skipBatchTitle"),
          message: this.$t("prompts.skipBatchMessage", { count: remaining }),
          danger: true,
          confirmText: this.$t("prompts.skipBatchConfirm"),
        },
        confirm: async () => {
          mutations.closeTopPrompt();
          try {
            await skipBatch();
          } catch (err) {
            console.error("Skip batch failed:", err);
            notify.showError(this.$t("prompts.compressControlError"));
          }
        },
      });
    },
    handleCancel() {
      const remaining = (this.status?.totalFiles || 0) - (this.status?.totalProcessed || 0);
      mutations.showPrompt({
        name: "confirmAction",
        props: {
          title: this.$t("prompts.cancelAllTitle"),
          message: this.$t("prompts.cancelAllMessage", { count: remaining }),
          danger: true,
          confirmText: this.$t("prompts.cancelAllConfirm"),
        },
        confirm: async () => {
          mutations.closeTopPrompt();
          try {
            await cancelCompress();
          } catch (err) {
            console.error("Cancel failed:", err);
            notify.showError(this.$t("prompts.compressControlError"));
          }
        },
      });
    },
  },
};
</script>

<style scoped>
.compress-status-bar {
  position: fixed;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  max-width: 500px;
  width: 100%;
  background: rgba(30, 30, 30, 0.95);
  color: white;
  padding: 10px 16px;
  border-radius: 8px 8px 0 0;
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.3);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  z-index: 1000;
}
.status-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  flex: 1;
  min-width: 0;
}
.status-icon {
  font-size: 1.1em;
}
.status-text {
  font-size: 0.9em;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.status-file {
  font-size: 0.8em;
  opacity: 0.7;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.status-buttons {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
}
.status-btn {
  background: rgba(255, 255, 255, 0.1);
  border: none;
  color: white;
  width: 36px;
  height: 36px;
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}
.status-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}
.status-btn .material-symbols {
  font-size: 20px;
}
@media (max-width: 768px) {
  .compress-status-bar {
    flex-direction: column;
    border-radius: 0;
    padding: 8px 12px;
  }
  .status-info {
    width: 100%;
  }
  .status-buttons {
    width: 100%;
    justify-content: space-around;
  }
  .status-btn {
    width: 44px;
    height: 44px;
  }
  .status-btn .material-symbols {
    font-size: 24px;
  }
}
</style>
