<template>
  <div class="card-content">
    <!-- Progress View -->
    <div v-if="compressing" class="compress-progress-view">
      <LoadingSpinner size="medium" />
      <div class="compress-progress-bar-container">
        <div class="compress-progress-bar" :style="{ width: progressPercent + '%' }"></div>
      </div>
      <p class="compress-progress-text">
        {{ $t("prompts.compressProgress", { current: progressData.current || 0, total: progressData.total || 0 }) }}
      </p>
      <p v-if="progressData.currentFile" class="compress-progress-file">
        {{ progressData.currentFile }}
      </p>
    </div>

    <!-- Main Dialog -->
    <div v-else>
      <!-- Confirmation Dialog -->
      <div v-if="showConfirm" class="compress-confirm-overlay" @click.self="showConfirm = false">
        <div class="compress-confirm-dialog">
          <p>{{ $t("prompts.compressConfirm", { count: selectedFileCount }) }}</p>
          <div class="compress-confirm-actions">
            <button type="button" class="button button--flat button--grey" @click="showConfirm = false">
              {{ $t("general.cancel") }}
            </button>
            <button type="button" class="button button--flat" @click="doCompress">
              {{ $t("general.ok") }}
            </button>
          </div>
        </div>
      </div>

      <!-- File List Area -->
      <div class="compress-file-list">
        <p class="prompts-label">{{ $t("prompts.compressImages") }} ({{ selectedFileCount }})</p>
        <div class="compress-groups">
          <div
            v-for="(group, folder) in groupedFiles"
            :key="folder"
            class="compress-group"
          >
            <div class="compress-group-header clickable" @click="toggleGroup(folder)">
              <i class="material-symbols">{{ collapsedGroups[folder] ? "expand_more" : "chevron_right" }}</i>
              <span class="compress-group-folder">{{ folder }}</span>
              <span class="compress-group-count">({{ group.length }})</span>
            </div>
            <div v-show="!collapsedGroups[folder]" class="compress-group-items">
              <label
                v-for="file in group"
                :key="file.path"
                class="compress-file-item"
                :class="{ 'compress-file-item--selected': selectedFiles[file.path] }"
              >
                <input
                  type="checkbox"
                  v-model="selectedFiles[file.path]"
                />
                <img
                  v-if="file.thumbnailUrl"
                  :src="file.thumbnailUrl"
                  class="compress-file-thumb"
                  loading="lazy"
                  alt=""
                />
                <i v-else class="material-symbols compress-file-icon">image</i>
                <span class="compress-file-name">{{ file.name }}</span>
                <span class="compress-file-size">{{ formatSize(file.size) }}</span>
              </label>
            </div>
          </div>
        </div>
      </div>

      <!-- Tier Selection -->
      <div class="compress-tier-section">
        <p class="prompts-label">{{ $t("prompts.compressQuality") }}</p>
        <div class="compress-tier-buttons">
          <button
            type="button"
            class="button compress-tier-btn"
            :class="{ 'compress-tier-btn--active': selectedTier === 'low' }"
            @click="selectTier('low')"
          >
            {{ $t("prompts.compressTierLow") }}
          </button>
          <button
            type="button"
            class="button compress-tier-btn"
            :class="{ 'compress-tier-btn--active': selectedTier === 'medium' }"
            @click="selectTier('medium')"
          >
            {{ $t("prompts.compressTierMedium") }}
          </button>
          <button
            type="button"
            class="button compress-tier-btn"
            :class="{ 'compress-tier-btn--active': selectedTier === 'high' }"
            @click="selectTier('high')"
          >
            {{ $t("prompts.compressTierHigh") }}
          </button>
        </div>
      </div>

      <!-- Advanced Settings -->
      <div class="compress-advanced-section">
        <button
          type="button"
          class="button button--flat compress-advanced-toggle"
          @click="showAdvanced = !showAdvanced"
        >
          <i class="material-symbols">{{ showAdvanced ? "expand_less" : "expand_more" }}</i>
          {{ $t("prompts.compressAdvanced") }}
        </button>
        <div v-if="showAdvanced" class="compress-advanced-content">
          <label class="compress-quality-slider-label">
            {{ $t("prompts.compressQualityLevel") }}: {{ quality }}
          </label>
          <input
            type="range"
            class="compress-quality-slider"
            :min="qualityRange.min"
            :max="qualityRange.max"
            v-model.number="quality"
            @change="updatePreview"
          />
          <span class="compress-quality-range">{{ qualityRange.min }} - {{ qualityRange.max }}</span>
        </div>
      </div>

      <!-- Preview Comparison -->
      <div v-if="previewUrls.original || previewUrls.compressed" class="compress-preview-section">
        <p class="prompts-label">{{ $t("prompts.compressPreview") }}</p>
        <div class="compress-preview-container">
          <div class="compress-preview-item">
            <span class="compress-preview-label">{{ $t("prompts.compressOriginal") }}</span>
            <img
              v-if="previewUrls.original"
              :src="previewUrls.original"
              class="compress-preview-img clickable"
              @click="zoomImage(previewUrls.original)"
              alt="Original"
            />
            <span v-if="previewData.originalSize" class="compress-preview-size">
              {{ formatSize(previewData.originalSize) }}
            </span>
          </div>
          <div class="compress-preview-item">
            <span class="compress-preview-label">{{ $t("prompts.compressCompressed") }}</span>
            <img
              v-if="previewUrls.compressed"
              :src="previewUrls.compressed"
              class="compress-preview-img clickable"
              @click="zoomImage(previewUrls.compressed)"
              alt="Compressed"
            />
            <span v-if="previewData.compressedSize" class="compress-preview-size">
              {{ formatSize(previewData.compressedSize) }}
              <span v-if="previewData.savingsPercent" class="compress-savings">
                (−{{ previewData.savingsPercent }}%)
              </span>
            </span>
          </div>
        </div>
      </div>

      <!-- Options -->
      <div class="compress-options settings-items">
        <ToggleSwitch
          class="item"
          v-model="backupEnabled"
          :name="$t('prompts.compressBackup')"
          :description="$t('prompts.compressBackupDescription')"
        />
        <p v-if="backupEnabled" class="compress-backup-name">
          {{ $t("prompts.compressBackupName") }}: <code>{{ backupFileName }}</code>
        </p>
      </div>
    </div>
  </div>

  <!-- Zoom Overlay -->
  <div v-if="zoomUrl" class="compress-zoom-overlay" @click="zoomUrl = null">
    <img :src="zoomUrl" class="compress-zoom-img" alt="Zoomed" />
  </div>

  <!-- Action Buttons -->
  <div class="card-actions">
    <template v-if="!compressing">
      <button
        type="button"
        class="button button--flat button--grey"
        @click="closeTopPrompt"
        :aria-label="$t('general.cancel')"
      >
        {{ $t("general.cancel") }}
      </button>
      <button
        type="button"
        class="button button--flat"
        :disabled="selectedFileCount === 0"
        @click="confirmCompress"
      >
        {{ $t("prompts.compressImages") }}
      </button>
    </template>
  </div>
</template>

<script>
import { mutations } from "@/store";
import { notify } from "@/notify";
import { previewCompress, startCompress, subscribeProgress } from "@/api/compress.js";
import LoadingSpinner from "@/components/LoadingSpinner.vue";
import ToggleSwitch from "@/components/settings/ToggleSwitch.vue";

// Tier quality ranges
const TIER_RANGES = {
  low: { min: 60, max: 100, default: 85 },
  medium: { min: 30, max: 80, default: 60 },
  high: { min: 10, max: 50, default: 30 },
};

export default {
  name: "compressImages",
  components: { LoadingSpinner, ToggleSwitch },
  props: {
    items: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      selectedTier: "medium",
      quality: TIER_RANGES.medium.default,
      showAdvanced: false,
      previewFile: null,
      previewUrls: { original: null, compressed: null },
      previewData: { originalSize: 0, compressedSize: 0, savingsPercent: 0 },
      backupEnabled: false,
      compressing: false,
      showConfirm: false,
      zoomUrl: null,
      progressData: { current: 0, total: 0, currentFile: null },
      selectedFiles: {},
      collapsedGroups: {},
      progressSubscription: null,
    };
  },
  computed: {
    selectedFileCount() {
      return Object.values(this.selectedFiles).filter(Boolean).length;
    },
    selectedFileList() {
      return this.items.filter((item) => this.selectedFiles[item.path]);
    },
    groupedFiles() {
      const groups = {};
      for (const file of this.items) {
        const folder = this.getFolder(file.path);
        if (!groups[folder]) groups[folder] = [];
        groups[folder].push(file);
      }
      return groups;
    },
    qualityRange() {
      const range = TIER_RANGES[this.selectedTier] || TIER_RANGES.medium;
      return range;
    },
    backupFileName() {
      const firstFile = this.items[0];
      if (!firstFile) return "";
      const parts = firstFile.path.split("/");
      const fileName = parts[parts.length - 1];
      // Naming rule: original.zst (e.g., photo.jpg -> photo.jpg.zst)
      return `${fileName}.zst`;
    },
    progressPercent() {
      if (!this.progressData.total) return 0;
      return Math.round((this.progressData.current / this.progressData.total) * 100);
    },
  },
  mounted() {
    // Defensive guard: items may not be passed
    if (!this.items || !Array.isArray(this.items) || this.items.length === 0) {
      return;
    }
    // Select all files by default
    for (const item of this.items) {
      this.selectedFiles[item.path] = true;
    }
    // Auto-preview the first file
    this.previewFile = this.items[0];
    this.updatePreview();
  },
  beforeUnmount() {
    if (this.progressSubscription) {
      this.progressSubscription.close();
    }
  },
  methods: {
    closeTopPrompt() {
      mutations.closeTopPrompt();
    },
    getFolder(path) {
      const parts = path.split("/");
      parts.pop();
      return parts.join("/") || "/";
    },
    toggleGroup(folder) {
      this.collapsedGroups[folder] = !this.collapsedGroups[folder];
    },
    selectTier(tier) {
      this.selectedTier = tier;
      this.quality = TIER_RANGES[tier].default;
      this.updatePreview();
    },
    async updatePreview() {
      if (!this.previewFile) return;
      const file = this.previewFile;
      try {
        const result = await previewCompress({
          source: file.source,
          path: file.path,
          tier: this.selectedTier,
          quality: this.quality,
        });
        this.previewUrls.original = result.previewUrl || null;
        this.previewUrls.compressed = result.compressedPreviewUrl || null;
        this.previewData.originalSize = result.originalSize || file.size || 0;
        this.previewData.compressedSize = result.compressedSize || 0;
        if (this.previewData.originalSize > 0 && this.previewData.compressedSize > 0) {
          const savings = Math.round(
            (1 - this.previewData.compressedSize / this.previewData.originalSize) * 100
          );
          this.previewData.savingsPercent = savings > 0 ? savings : 0;
        }
      } catch (err) {
        // Silent fail on preview - not critical
        console.error("Preview error:", err);
      }
    },
    confirmCompress() {
      this.showConfirm = true;
    },
    async doCompress() {
      this.showConfirm = false;
      this.compressing = true;
      const files = this.selectedFileList;
      this.progressData.total = files.length;
      this.progressData.current = 0;
      const source = files[0]?.source;
      const filePaths = files.map((f) => ({ path: f.path }));
      try {
        const result = await startCompress({
          source,
          files: filePaths,
          tier: this.selectedTier,
          quality: this.quality,
          backup: this.backupEnabled,
        });
        // Subscribe to progress
        this.progressSubscription = subscribeProgress(result.jobId, {
          onProgress: (data) => {
            this.progressData = data;
          },
          onComplete: (data) => {
            this.compressing = false;
            this.progressData.current = this.progressData.total;
            mutations.setReload(true);
            mutations.closeTopPrompt();
            notify.showSuccess(
              this.$t("prompts.compressComplete", {
                count: data.totalFiles || files.length,
                saved: this.formatSize(data.totalSaved || 0),
              })
            );
          },
          onError: (data) => {
            this.compressing = false;
            notify.showError(data.message || "Compression failed");
          },
        });
      } catch (err) {
        this.compressing = false;
        console.error(err);
      }
    },
    formatSize(bytes) {
      if (!bytes || bytes === 0) return "0 B";
      const units = ["B", "KB", "MB", "GB"];
      let size = bytes;
      let unitIndex = 0;
      while (size >= 1024 && unitIndex < units.length - 1) {
        size /= 1024;
        unitIndex++;
      }
      return `${size.toFixed(unitIndex === 0 ? 0 : 1)} ${units[unitIndex]}`;
    },
    zoomImage(url) {
      this.zoomUrl = url;
    },
  },
};
</script>

<style scoped>
.card-content {
  position: relative;
}

.prompts-label {
  margin-top: 1em;
  margin-bottom: 0.25em;
  font-weight: 500;
}

/* File List */
.compress-file-list {
  max-height: 250px;
  overflow-y: auto;
  border: 1px solid var(--borderColor);
  border-radius: 4px;
  padding: 0.5em;
}

.compress-groups {
  display: flex;
  flex-direction: column;
  gap: 0.25em;
}

.compress-group-header {
  display: flex;
  align-items: center;
  gap: 0.25em;
  padding: 0.25em;
  font-weight: 500;
  font-size: 0.9em;
  color: var(--textSecondary);
}

.compress-group-header .material-symbols {
  font-size: 18px;
}

.compress-group-count {
  opacity: 0.7;
  font-size: 0.85em;
}

.compress-group-items {
  padding-left: 1.5em;
}

.compress-file-item {
  display: flex;
  align-items: center;
  gap: 0.5em;
  padding: 0.25em 0.5em;
  cursor: pointer;
  border-radius: 4px;
  transition: background 0.15s;
}

.compress-file-item:hover {
  background: var(--surfacePrimaryAlt);
}

.compress-file-item--selected {
  background: var(--surfacePrimaryAlt);
}

.compress-file-thumb {
  width: 32px;
  height: 32px;
  object-fit: cover;
  border-radius: 4px;
  flex-shrink: 0;
}

.compress-file-icon {
  font-size: 24px;
  color: var(--textSecondary);
  flex-shrink: 0;
}

.compress-file-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 0.9em;
}

.compress-file-size {
  font-size: 0.85em;
  color: var(--textSecondary);
  flex-shrink: 0;
}

/* Tier Selection */
.compress-tier-section {
  margin-top: 1em;
}

.compress-tier-buttons {
  display: flex;
  gap: 0.5em;
}

.compress-tier-btn {
  flex: 1;
  border: 1px solid var(--borderColor);
  background: transparent;
}

.compress-tier-btn--active {
  background: var(--primaryColor);
  color: var(--textPrimaryAlt);
  border-color: var(--primaryColor);
}

/* Advanced */
.compress-advanced-section {
  margin-top: 1em;
}

.compress-advanced-toggle {
  display: flex;
  align-items: center;
  gap: 0.25em;
  padding: 0.25em 0;
}

.compress-advanced-toggle .material-symbols {
  font-size: 20px;
}

.compress-advanced-content {
  padding: 0.5em 0;
}

.compress-quality-slider-label {
  font-size: 0.9em;
  margin-bottom: 0.5em;
  display: block;
}

.compress-quality-slider {
  width: 100%;
  cursor: pointer;
}

.compress-quality-range {
  font-size: 0.8em;
  color: var(--textSecondary);
}

/* Preview */
.compress-preview-section {
  margin-top: 1em;
}

.compress-preview-container {
  display: flex;
  gap: 1em;
}

.compress-preview-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.25em;
}

.compress-preview-label {
  font-size: 0.85em;
  font-weight: 500;
  color: var(--textSecondary);
}

.compress-preview-img {
  max-width: 100%;
  max-height: 150px;
  object-fit: contain;
  border: 1px solid var(--borderColor);
  border-radius: 4px;
}

.compress-preview-size {
  font-size: 0.85em;
  color: var(--textSecondary);
}

.compress-savings {
  color: var(--primaryColor);
  font-weight: 500;
}

/* Options */
.compress-options {
  margin-top: 1em;
}

.compress-backup-name {
  font-size: 0.85em;
  color: var(--textSecondary);
  margin-top: 0.5em;
}

.compress-backup-name code {
  background: var(--surfacePrimaryAlt);
  padding: 0.1em 0.3em;
  border-radius: 3px;
  font-size: 0.9em;
}

/* Confirmation Dialog */
.compress-confirm-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10;
}

.compress-confirm-dialog {
  background: var(--surfaceSecondary);
  border-radius: 8px;
  padding: 1.5em;
  max-width: 80%;
  text-align: center;
}

.compress-confirm-actions {
  display: flex;
  gap: 0.5em;
  justify-content: center;
  margin-top: 1em;
}

/* Progress View */
.compress-progress-view {
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1em;
  padding: 2em 0;
  min-height: 200px;
}

.compress-progress-bar-container {
  width: 100%;
  height: 8px;
  background: var(--surfacePrimaryAlt);
  border-radius: 4px;
  overflow: hidden;
}

.compress-progress-bar {
  height: 100%;
  background: var(--primaryColor);
  transition: width 0.3s ease;
}

.compress-progress-text {
  font-weight: 500;
  margin: 0;
}

.compress-progress-file {
  font-size: 0.85em;
  color: var(--textSecondary);
  margin: 0;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* Zoom Overlay */
.compress-zoom-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.85);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.compress-zoom-img {
  max-width: 90vw;
  max-height: 90vh;
  object-fit: contain;
}
</style>
