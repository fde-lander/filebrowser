import { notify } from "@/notify";
import { getApiPath } from "@/utils/url.js";
import { fetchURL } from "./utils";

/**
 * Preview image compression for a single file without applying it.
 * Returns compressed image as binary blob + size info from headers.
 * @param {Object} opts
 * @param {string} opts.source - Source name where the file lives
 * @param {string} opts.path - File path
 * @param {string} opts.level - Compression level: "low" | "medium" | "high"
 * @param {number} [opts.quality] - Quality value (1-100), overrides level default
 * @returns {Promise<Object>} { blob, originalSize, compressedSize, skipped, contentType }
 */
export async function previewCompress(opts) {
  const { source, path, level, quality } = opts;
  if (!source || !path || !level) {
    throw new Error("source, path, and level are required");
  }
  const body = {
    source,
    path,
    level,
    ...(quality !== undefined && quality !== null && { quality }),
  };
  try {
    const apiPath = getApiPath("compress-images/preview");
    const response = await fetchURL(apiPath, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });
    // Backend returns binary blob + headers, NOT JSON
    const blob = await response.blob();
    const originalSize = parseInt(response.headers.get("X-Original-Size") || "0", 10);
    const compressedSize = parseInt(response.headers.get("X-Compressed-Size") || "0", 10);
    const skipped = response.headers.get("X-Skipped") === "true";
    const contentType = response.headers.get("Content-Type") || "image/webp";
    return { blob, originalSize, compressedSize, skipped, contentType };
  } catch (err) {
    notify.showError(err.message || "Error previewing compression");
    throw err;
  }
}

/**
 * Start a batch image compression job on the server.
 * Returns a job ID that can be used to subscribe to progress updates.
 * @param {Object} opts
 * @param {string} opts.source - Source name where files live
 * @param {Array<string>} opts.files - Array of file paths to compress
 * @param {string} opts.level - Compression level: "low" | "medium" | "high"
 * @param {number} [opts.quality] - Quality value (1-100), overrides level default
 * @param {boolean} [opts.backup] - If true, create .tar.zst backup of original files
 * @param {string} [opts.backupPath] - Directory path for backup file
 * @param {string} [opts.backupName] - Filename for backup file
 * @returns {Promise<Object>} { taskId, message }
 */
export async function startCompress(opts) {
  const { source, files, level, quality, backup, backupPath, backupName } = opts;
  if (!source || !files?.length || !level) {
    throw new Error("source, files, and level are required");
  }
  const flatFiles = files.map(f => typeof f === 'string' ? f : f.path);
  const body = {
    source,
    files: flatFiles,
    level,
    ...(quality !== undefined && quality !== null && { quality }),
    ...(backup && { backup: true, backupPath, backupName }),
  };
  try {
    const apiPath = getApiPath("compress-images");
    const response = await fetchURL(apiPath, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });
    return response.json();
  } catch (err) {
    notify.showError(err.message || "Error starting compression");
    throw err;
  }
}

/**
 /**
  * Poll compression queue status via HTTP GET.
  * @returns {Promise<Object>} Queue status: { status, currentFile, processed, total, ... }
  */
 export async function pollStatus() {
   const apiPath = getApiPath("compress-images/status", {});
   const response = await fetchURL(apiPath, {});
   if (!response.ok) {
     throw new Error(`Status API returned ${response.status}`);
   }
   return response.json();
 }
