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
                @click="handleFileClick(file, $event)"
              >
                <input
                  v-if="!previewMode"
                  type="checkbox"
                  v-model="selectedFiles[file.path]"
                />
                <i class="material-symbols compress-file-icon">image</i>
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
          />
          <span class="compress-quality-range">{{ qualityRange.min }} - {{ qualityRange.max }}</span>
        </div>
      </div>

      <!-- Preview mode toggle -->
      <div class="compress-preview-toggle">
        <label>
          <input type="checkbox" v-model="previewMode" />
          {{ $t("prompts.compressPreviewToggle") }}
        </label>
      </div>

      <!-- Options -->
      <div class="compress-options settings-items">
        <ToggleSwitch
          class="item"
          v-model="backupEnabled"
          @change="onBackupToggle"
          :name="$t('prompts.compressBackup')"
          :description="$t('prompts.compressBackupDescription')"
        />
        <p v-if="backupEnabled" class="compress-backup-name">
          {{ $t("prompts.compressBackupName") }}: <code>{{ backupFileName }}</code>
        </p>
      </div>
    </div>
  </div>

  <!-- Preview Overlay -->
  <div v-if="previewOverlay" class="compress-preview-overlay">
    <div class="compress-preview-header">
      <button class="compress-preview-back-btn" @click="closePreview">
        <i class="material-symbols">arrow_back</i>
        {{ $t("prompts.compressPreviewBack") }}
      </button>
    </div>
    <div class="compress-preview-body">
      <div class="compress-preview-item" @click="openFullscreen(previewUrls.original)">
        <span class="compress-preview-label">{{ $t("prompts.compressOriginal") }}</span>
        <img v-if="previewUrls.original" :src="previewUrls.original"
             class="compress-preview-img" alt="Original" />
        <span v-if="previewData.originalSize" class="compress-preview-size">
          {{ formatSize(previewData.originalSize) }}
        </span>
      </div>
      <div class="compress-preview-item" @click="openFullscreen(previewUrls.compressed)">
        <span class="compress-preview-label">{{ $t("prompts.compressCompressed") }}</span>
        <LoadingSpinner v-if="previewLoading" size="small" />
        <img v-if="previewUrls.compressed && !previewLoading"
             :src="previewUrls.compressed"
             class="compress-preview-img compress-preview-fade-in"
             alt="Compressed" />
        <span v-if="previewData.compressedSize && !previewLoading" class="compress-preview-size">
          {{ formatSize(previewData.compressedSize) }}
          <span v-if="previewData.savingsPercent" class="compress-savings">
            (−{{ previewData.savingsPercent }}%)
          </span>
        </span>
      </div>
    </div>
  </div>

  <!-- Fullscreen Overlay -->
  <div v-if="fullscreenUrl" class="compress-fullscreen-overlay" @click.self="closeFullscreen">
    <img :src="fullscreenUrl" class="compress-fullscreen-img" />
    <button class="compress-fullscreen-back" @click="closeFullscreen">
      <i class="material-symbols">arrow_back</i>
    </button>
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
import { state, mutations } from "@/store";
import { notify } from "@/notify";
import { previewCompress, startCompress, pollStatus } from "@/api/compress.js";
import { getPreviewURL, fetchFiles } from "@/api/resources.js";
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
      previewMode: false,
      previewLoading: false,
      previewOverlay: false,
      fullscreenUrl: null,
      backupEnabled: true,
      expandingItems: false,
      expandedItems: [],
      compressing: false,
      showConfirm: false,
      zoomUrl: null,
      progressData: { current: 0, total: 0, currentFile: null },
      selectedFiles: {},
      collapsedGroups: {},
      pollTimer: null,
    };
  },
  computed: {
    selectedFileCount() {
      return Object.values(this.selectedFiles).filter(Boolean).length;
    },
    selectedFileList() {
      const items = this.expandedItems.length > 0 ? this.expandedItems : this.items;
      return items.filter((item) => this.selectedFiles[item.path]);
    },
    groupedFiles() {
      const groups = {};
      const items = this.expandedItems.length > 0 ? this.expandedItems : this.items;
      for (const file of items) {
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
      const parts = firstFile.path.split("/").filter(Boolean);
      const firstName = parts[parts.length - 1] || "backup";
      const count = this.items.length;
      if (count === 1) {
        return `${firstName}.tar.zst`;
      }
      const others = count - 1;
      return `${firstName}_and_${others}_other${others > 1 ? 's' : ''}.tar.zst`;
    },
    backupPath() {
      const firstFile = this.items[0];
      if (!firstFile) return "/";
      const parts = firstFile.path.split("/");
      parts.pop();
      return parts.join("/") || "/";
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
    // Read backup preference from store (default ON)
    this.backupEnabled = state.user?.compressBackup ?? true;
    // Expand folders to individual image files
    this.expandItems();
  },
  beforeUnmount() {
    this.cleanupBlobUrls();
    if (this.pollTimer) {
      clearInterval(this.pollTimer);
    }
  },
  methods: {
    closeTopPrompt() {
      mutations.closeTopPrompt();
    },
    onBackupToggle() {
      mutations.updateCurrentUser({ compressBackup: this.backupEnabled });
    },
    isCompressableImage(path) {
      const ext = path.split('.').pop().toLowerCase();
      return ['jpg', 'jpeg', 'png', 'bmp', 'webp'].includes(ext);
    },
    async expandItems() {
      this.expandingItems = true;
      const expanded = [];
      for (const item of this.items) {
        if (item.isDir || item.type === 'directory') {
          try {
            const result = await fetchFiles(item.source, item.path);
            for (const f of (result.items || [])) {
              if (!f.isDir && f.type !== 'directory' && this.isCompressableImage(f.path)) {
                expanded.push(f);
              }
            }
          } catch (e) { /* skip folder on error */ }
        } else {
          if (this.isCompressableImage(item.path)) {
            expanded.push(item);
          }
        }
      }
      this.expandedItems = expanded;
      for (const item of expanded) {
        this.selectedFiles[item.path] = true;
      }
      this.expandingItems = false;
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
    },
    handleFileClick(file, event) {
      if (this.previewMode) {
        event.preventDefault();
        this.openPreview(file);
      }
    },
    async openPreview(file) {
      this.previewOverlay = true;
      this.previewLoading = true;
      this.previewFile = file;

      // Revoke previous blob URLs
      this.cleanupBlobUrls();

      // Get original image URL (instant)
      this.previewUrls.original = getPreviewURL(file.source, file.path, file.modified) + '&size=original';

      // Reset compressed side
      this.previewUrls.compressed = null;
      this.previewData.originalSize = file.size || 0;
      this.previewData.compressedSize = 0;
      this.previewData.savingsPercent = 0;

      try {
        const result = await previewCompress({
          source: file.source,
          path: file.path,
          level: this.selectedTier,
          quality: this.quality,
        });

        this.previewUrls.compressed = URL.createObjectURL(result.blob);
        this.previewData.originalSize = result.originalSize || file.size || 0;
        this.previewData.compressedSize = result.compressedSize || 0;

        if (this.previewData.originalSize > 0 && this.previewData.compressedSize > 0) {
          const savings = Math.round(
            (1 - this.previewData.compressedSize / this.previewData.originalSize) * 100
          );
          this.previewData.savingsPercent = savings > 0 ? savings : 0;
        }
      } catch (err) {
        console.error("Preview error:", err);
        this.previewUrls.compressed = null;
      } finally {
        this.previewLoading = false;
      }
    },
    closePreview() {
      this.previewOverlay = false;
      this.cleanupBlobUrls();
      this.previewUrls.original = null;
      this.previewUrls.compressed = null;
    },
    openFullscreen(url) {
      if (url) this.fullscreenUrl = url;
    },
    closeFullscreen() {
      this.fullscreenUrl = null;
    },
    cleanupBlobUrls() {
      if (this.previewUrls.compressed && this.previewUrls.compressed.startsWith('blob:')) {
        URL.revokeObjectURL(this.previewUrls.compressed);
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
      const filePaths = files.map((f) => f.path);
      try {
        const result = await startCompress({
          source,
          files: filePaths,
          level: this.selectedTier,
          quality: this.quality,
          backup: this.backupEnabled,
          backupPath: this.backupPath,
          backupName: this.backupFileName,
        });

        // Start polling for status
        this.pollTimer = setInterval(async () => {
          try {
            const status = await pollStatus();
            if (status.status === "running") {
              this.progressData = {
                current: status.processed || 0,
                total: status.total || 0,
                currentFile: status.currentFile || null,
              };
            } else if (status.status === "completed") {
              clearInterval(this.pollTimer);
              this.pollTimer = null;
              this.compressing = false;
              this.progressData.current = this.progressData.total;
              mutations.setReload(true);
              mutations.closeTopPrompt();

              const savedStr = this.formatSize(status.savedBytes || 0);
              const successCount = status.processed - (status.skipped || 0) - (status.failed || 0);
              const failedCount = status.failed || 0;

              if (failedCount > 0 && successCount <= 0) {
                notify.showError(
                  this.$t("prompts.compressFailed", { count: failedCount })
                );
              } else {
                let msgParts = [
                  this.$t("prompts.compressComplete", {
                    count: successCount,
                    saved: savedStr,
                  })
                ];
                if (status.backupPath) {
                  msgParts.push(this.$t("prompts.compressBackupPath", {
                    path: status.backupPath
                  }));
                  if (status.backupFallback) {
                    msgParts.push(this.$t("prompts.compressBackupFallback"));
                  }
                }
                notify.showSuccess(msgParts.join(" "));
              }
            } else if (status.status === "failed") {
              clearInterval(this.pollTimer);
              this.pollTimer = null;
              this.compressing = false;
              notify.showError(this.$t("prompts.compressFailed", { count: status.failed }));
            }
          } catch (err) {
            console.error("Status poll error:", err);
            // Keep polling - network error doesn't mean backend failed
          }
        }, 3000);
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

/* Preview Mode Toggle */
.compress-preview-toggle {
  padding: 0.5em 0;
}

/* Preview Overlay */
.compress-preview-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: var(--background);
  z-index: 100;
  display: flex;
  flex-direction: column;
}
.compress-preview-header {
  padding: 1em;
  border-bottom: 1px solid var(--divider);
}
.compress-preview-back-btn {
  display: flex;
  align-items: center;
  gap: 0.5em;
  background: none;
  border: none;
  cursor: pointer;
  color: var(--textPrimary);
  font-size: 1em;
}
.compress-preview-body {
  display: flex;
  flex: 1;
  overflow: hidden;
}
.compress-preview-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 1em;
  border-right: 1px solid var(--divider);
  cursor: pointer;
  min-height: 0;
}
.compress-preview-item:last-child {
  border-right: none;
}
.compress-preview-img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}
.compress-preview-fade-in {
  animation: compressFadeIn 300ms ease;
}
@keyframes compressFadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}
.compress-preview-size {
  flex: 0 0 auto;
  margin-top: 0.5em;
  padding: 0.5em;
  font-size: 0.9em;
  color: var(--textSecondary);
}

/* Fullscreen Overlay */
.compress-fullscreen-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: #000;
  z-index: 200;
  display: flex;
  align-items: center;
  justify-content: center;
}
.compress-fullscreen-img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}
.compress-fullscreen-back {
  position: fixed;
  bottom: 1em;
  right: 1em;
  z-index: 201;
  background: rgba(255,255,255,0.2);
  border: none;
  color: white;
  cursor: pointer;
  padding: 0.5em 1em;
  border-radius: 4px;
}

@media (max-width: 768px) {
  .compress-preview-body { flex-direction: column; }
  .compress-preview-img { max-height: 35vh; }
  .compress-tier-buttons { flex-wrap: wrap; }
  .compress-tier-btn { flex: 1 1 auto; min-width: 80px; }
  .compress-fullscreen-back { min-width: 44px; min-height: 44px; font-size: 1.5em; }
  .compress-file-name { white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
  .compress-queue-bar { font-size: 0.85em; padding: 0.5em; }
}
</style>
