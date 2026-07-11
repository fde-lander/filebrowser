import { notify } from "@/notify";
import { getApiPath } from "@/utils/url.js";
import { fetchURL } from "./utils";

/**
 * Preview image compression for a single file without applying it.
 * Returns estimated compressed size and a preview thumbnail URL.
 * @param {Object} opts
 * @param {string} opts.source - Source name where the file lives
 * @param {string} opts.path - File path
 * @param {string} opts.tier - Compression tier: "low" | "medium" | "high"
 * @param {number} [opts.quality] - Quality value (1-100), overrides tier default
 * @returns {Promise<Object>} { originalSize, compressedSize, previewUrl }
 */
export async function previewCompress(opts) {
  const { source, path, tier, quality } = opts;
  if (!source || !path || !tier) {
    throw new Error("source, path, and tier are required");
  }
  const body = {
    source,
    path,
    tier,
    ...(quality !== undefined && quality !== null && { quality }),
  };
  try {
    const apiPath = getApiPath("resources/compress/preview");
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
 * @param {Array<{path: string}>} opts.files - Array of file paths to compress
 * @param {string} opts.tier - Compression tier: "low" | "medium" | "high"
 * @param {number} [opts.quality] - Quality value (1-100), overrides tier default
 * @param {boolean} [opts.backup] - If true, create .zst backup of original files
 * @returns {Promise<Object>} { jobId, totalFiles }
 */
export async function startCompress(opts) {
  const { source, files, tier, quality, backup } = opts;
  if (!source || !files?.length || !tier) {
    throw new Error("source, files, and tier are required");
  }
  const body = {
    source,
    files,
    tier,
    ...(quality !== undefined && quality !== null && { quality }),
    ...(backup && { backup: true }),
  };
  try {
    const apiPath = getApiPath("resources/compress");
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
 * @param {string} jobId - The job ID returned by startCompress
 * @param {Object} handlers
 * @param {function(Object)} [handlers.onProgress] - Called with { current, total, currentFile, compressedSize, originalSize }
 * @param {function(Object)} [handlers.onComplete] - Called with { totalFiles, totalSaved }
 * @param {function(Object)} [handlers.onError] - Called with { message }
 * @returns {{ eventSource: EventSource, close: function() }}
 */
export function subscribeProgress(jobId, handlers = {}) {
  if (!jobId) {
    throw new Error("jobId is required");
  }
  const apiPath = getApiPath(`resources/compress/progress`, { jobId });
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
