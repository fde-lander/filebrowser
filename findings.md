# Findings & Decisions (Round 2)

## User Requirements

### Feature A: Image Viewer Enhancement

**Problem 1: Black screen flash on navigation**
- When viewing next image, current image disappears instantly → black screen → slowly loads next
- Need: preload next image + smooth transition (no abrupt disappearance)

**Problem 2: Mobile navigation buttons disappear too fast**
- On Android browser, left/right nav buttons vanish quickly
- Need: setting (default OFF) to keep buttons persistent + opacity control (default opaque)

**Problem 3: Tap navigation**
- Need: setting (default ON) - tap right side = next image, tap left-center = previous
- Must NOT interfere with double-tap zoom or drag-to-pan

**Problem 4: Smooth transitions with preload**
- Need: setting (default ON) - preload next image in background
- When user navigates, if next image is preloaded → smooth transition (fade, slide, etc)
- If not yet loaded → wait indicator → transition once ready
- Need: multiple transition effect options for user to choose

### Feature B: Image Compression

**Right-click "Compress Images" action**
- Works on folder (auto-selects all images) or single file
- Dialog with:
  - Image file list (pre-selected, user can deselect)
  - 3 compression levels (lossless or near-lossless for web viewing)
  - Preview: show selected image at chosen compression level vs original
  - ZSTD backup toggle (default ON): creates ZSTD backup of all images before compressing
  - Confirmation dialog before execution
- Progress indication during compression
- Completion notification

## Code Research Results (Inline Analysis)

### Image Viewer Architecture

**Core Components:**
- Preview.vue (712 lines) — main preview container, decides previewType (image/video/pdf/etc)
- ExtendedImage.vue (823 lines) — THE actual image rendering component
  - Uses <img> tag with JS-based loading (loadFullImage sets src via $refs)
  - Has thumbnail placeholder system (cachedThumbnailUrl shows while full image loads)
  - Has LoadingSpinner overlay while loading
  - Supports zoom/pan/double-tap/edge-gesture swipe navigation
  - fullImageLoaded boolean controls display: 'none' -> 'block'
- nextPrevious.vue (1187 lines) — navigation buttons + edge detection zones
  - showNav controls button visibility, auto-hides after 3 seconds (line 490-496)
  - navigationTimeout = setTimeout(... , 3000) — this is the "buttons disappear" problem
  - Left/right edge zones (5em wide, middle 50% height) for hover/touch detection
  - Already has link rel="prefetch" for prev/next URLs (line 91-92) but NOT true image preloading
- PopupPreview.vue (183 lines) — hover thumbnail popup (not for full image viewing)

**Image Loading Pipeline:**
1. Preview.vue: raw computed builds URL from resourcesApi (line 170-214)
2. ExtendedImage.vue: mounted() -> loadFullImage() -> sets img.src via $refs (line 227-239)
3. On navigation: state.navigation.isTransitioning = true -> Preview.vue shows LoadingSpinner overlay (line 4-6) -> ExtendedImage v-if="!isTransitioning" REMOVES the component entirely -> old image DISAPPEARS -> new route loads -> new ExtendedImage mounts -> black screen -> loads from scratch
4. THIS IS THE BLACK SCREEN FLASH ROOT CAUSE: v-if="isTransitioning" destroys ExtendedImage on every navigation

**Settings System (User Preferences):**
- Backend: users.go NonAdminEditable struct (line 130-169) — all user preferences are bool fields here
- Frontend: constants.js settings array (line 7-16) — defines settings page tabs
- Settings pages: frontend/src/views/settings/ — Profile.vue is where user toggles live
- ToggleSwitch.vue — reusable toggle component
- Existing pattern: deleteAfterArchive (bool), darkMode (bool), singleClick (bool), etc.
- To add new setting: add field to NonAdminEditable struct + frontend toggle in Profile settings + i18n keys

### Backend Image Processing Capabilities

**go.mod key dependencies:**
- golang.org/x/image v0.41.0 — standard Go image processing library (supports JPEG, PNG, GIF, BMP, TIFF decoding/encoding)
- github.com/kovidgoyal/imaging v1.8.21 — HIGH-LEVEL image processing (resize, crop, adjust quality, blur, sharpen)
  - This is EXCELLENT for image compression: imaging.Encode(img, w, imaging.JPEG, &jpeg.Options{Quality: 75})
- github.com/spf13/afero v1.15.0 — filesystem abstraction (used for all file operations)

**NO ZSTD library in go.mod** — need to add github.com/klauspost/compress/zstd

**Archive system (reference pattern for batch operations):**
- archive.go: unarchiveHandler accepts JSON body, processes files, returns JSON response
- Pattern: parse request -> validate -> execute -> return result
- For image compression: similar pattern, new endpoint /api/compress-images

**Preview system already has image processing:**
- backend/preview/image.go — generates thumbnails with imaging library
- backend/preview/image_enum.go — image format enum
- Can reference this code for compression implementation

### Backend Image Compression - Subagent Research Results

**Pure Go libraries (no CGO, Docker-friendly):**
- WebP: github.com/deepteams/webp (pure Go, 98% browser support, 25-34% smaller than JPEG)
- AVIF: github.com/gen2brain/avif (WASM-based, no CGO, 94% browser support, 50% smaller than JPEG)
- JPEG XL: github.com/gen2brain/jpegxl (WASM, but only 35% browser support - SKIP)
- ZSTD: github.com/klauspost/compress/zstd (pure Go, tar.zst archive pattern available)
- Palette: github.com/esimov/colorquant (pure Go, configurable palette + dithering)
- JPEG stdlib: image/jpeg (built-in, quality 1-100, no MozJPEG trellis optimization)

**CLI binary dependencies (better compression, need Docker install):**
- MozJPEG: cjpeg CLI (10-20% better than stdlib JPEG at same quality, trellis quantization)
- OxiPNG: oxipng CLI (Rust, lossless PNG optimization, 10-30% reduction)
- pngquant: pngquant CLI (PNG palette reduction, 40-70% smaller for 24/32-bit PNG)

**Skip for web use:**
- QOI: no browser support, worse than PNG - not suitable for FileBrowser web viewing
- JPEG XL: only 35% browser support - can't display in browser

**Recommended technology stack:**
- Primary conversion format: WebP (pure Go, best browser support, great compression)
- Secondary format: AVIF (CGO-free, even better compression, 94% browser support)
- PNG lossless optimization: OxiPNG CLI + esimov/colorquant palette reduction
- JPEG optimization: MozJPEG CLI (cjpeg) or stdlib image/jpeg as fallback
- ZSTD backup: klauspost/compress/zstd (tar.zst archive)

**CRITICAL: Compression ratios depend on source format!**

Key findings from subagent research:

- WebP q80 vs JPEG at same quality: only 25-34% smaller (NOT 60-70%)
- WebP q80 vs PNG 24-bit: 75-85% smaller (because PNG is lossless)
- WebP quality scale ≠ JPEG quality scale (WebP q80 ≈ JPEG q77 visual quality)
- AVIF q50 vs JPEG: ~60-65% smaller (NOT 75-85%), but with visible artifacts
- AVIF q55-60 vs JPEG: ~55-65% smaller, still good quality (recommended minimum for web)
- AVIF encoding speed: 1-4s per image at speed 6, ~2.5GB RAM per encode
- Go stdlib JPEG q92 re-encoding already-compressed JPEG: 0% to NEGATIVE (file gets LARGER!)
- Go stdlib JPEG q92 only works for uncompressed sources (PNG->JPEG): 75-85% reduction
- OxiPNG -o2: 30-50% reduction on unoptimized PNG (lossless)
- OxiPNG -o4: 5-10% extra over -o2 (6x slower)
- pngquant --quality=65-80: 70-80% reduction on 24-bit PNG (lossy)
- pngquant + oxipng pipeline: 65-80% total on 24-bit PNG
- pngquant first, then oxipng (correct order)

**Revised tier parameters based on actual data:**
- Low tier: JPEG->WebP q85 (25-30% smaller, near-lossless), PNG->OxiPNG -o2 (30-50% lossless)
- Medium tier: JPEG->WebP q75 (30-34% smaller), PNG->pngquant q70-85 + oxipng (60-75%)
- High tier: all->AVIF q55 (55-60% smaller vs JPEG, 87-89% vs PNG)

**REAL TEST DATA (actual measurements, not estimates):**

WebP compression (real images via subagent Go test):
- PNG -> WebP Q70: 98% reduction (1.5MB -> 35KB)
- PNG -> WebP Q75: 98% reduction
- PNG -> WebP Q80: 97% reduction
- JPEG Q70 -> WebP Q70: 36-41% reduction
- JPEG Q75 -> WebP Q75: 33-38% reduction
- JPEG Q80 -> WebP Q80: 16-24% reduction
- JPEG Q85 -> WebP Q85: -1% to 8% (minimal or LARGER)
- JPEG Q90 -> WebP Q90: NEGATIVE (file gets larger!)
- Lossless WebP: always larger than JPEG source
- KEY: WebP Q70-Q75 is the sweet spot for JPEG sources

PNG optimization (real CLI tests):
- OxiPNG -o2: 3-50% (lossless, depends on image)
- OxiPNG -o4: 4-72% (lossless, better but slower)
- pngquant q65-80: 62-74% (lossy, good for photos)
- pngquant q70-85: 57-69% (lossy, balanced)
- pngquant FAILS on smooth gradients (file gets LARGER)
- pngquant + oxipng pipeline: 61-72% (best for photos)
- OxiPNG alone best for flat UI graphics (60-72% lossless)

AVIF: encoding too slow for batch on 4GB server (>120s timeout per image)

---
*Update this file after every 2 view/browser/search operations*
