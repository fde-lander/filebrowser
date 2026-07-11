import { notify } from "@/notify";
import { getApiPath } from "@/utils/url.js";
import { fetchURL } from "./utils";

/**
 * Preview image compression for a single file without applying it.
 * Returns estimated compressed size and a preview thumbnail URL.
 * @param {Object} opts
 * @param {string} opts.source - Source name where the file lives
 * @param {string} opts.path - File path
 * @param {string} opts.level - Compression level: "low" | "medium" | "high"
 * @param {number} [opts.quality] - Quality value (1-100), overrides level default
 * @returns {Promise<Object>} { originalSize, compressedSize, previewUrl }
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
    return response.json();
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
 * Subscribe to compression job progress via Server-Sent Events (SSE).
 * Returns an object with the EventSource and a cleanup function.
 * @param {string} taskId - The task ID returned by startCompress
 * @param {Object} handlers
 * @param {function(Object)} [handlers.onProgress] - Called with { current, total, currentFile, compressedSize, originalSize }
 * @param {function(Object)} [handlers.onComplete] - Called with { totalFiles, totalSaved }
 * @param {function(Object)} [handlers.onError] - Called with { message }
 * @returns {{ eventSource: EventSource, close: function() }}
 */
export function subscribeProgress(taskId, handlers = {}) {
  if (!taskId) {
    throw new Error("taskId is required");
  }
  const apiPath = getApiPath(`compress-images/progress`, { taskId });
  const eventSource = new EventSource(apiPath);

  if (handlers.onProgress) {
    eventSource.addEventListener("progress", (event) => {
      try {
        handlers.onProgress(JSON.parse(event.data));
      } catch (e) {
        console.error("Error parsing progress SSE:", e);
      }
    });
  }

  if (handlers.onComplete) {
    eventSource.addEventListener("complete", (event) => {
      try {
        handlers.onComplete(JSON.parse(event.data));
      } catch (e) {
        console.error("Error parsing complete SSE:", e);
      }
      eventSource.close();
    });
  }

  if (handlers.onError) {
    eventSource.addEventListener("error", (event) => {
      // Check if this is a data error event (from the server) or a connection error
      if (event.data) {
        try {
          handlers.onError(JSON.parse(event.data));
        } catch (e) {
          handlers.onError({ message: "Compression job failed" });
        }
      }
      eventSource.close();
    });
  }

  // Handle raw EventSource connection errors (network failures)
  eventSource.onerror = () => {
    if (handlers.onError) {
      handlers.onError({ message: "Connection lost to compression progress stream" });
    }
    eventSource.close();
  };

  return {
    eventSource,
    close() {
      eventSource.close();
    },
  };
}
