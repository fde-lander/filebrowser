package http

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/deepteams/webp"
	"github.com/esimov/colorquant"
	"github.com/gtsteffaniak/filebrowser/backend/indexing"
	"github.com/gtsteffaniak/go-logger/logger"
	"github.com/klauspost/compress/zstd"
)

// --- Request/Response Types ---

type compressPreviewRequest struct {
	Source  string `json:"source"`
	Path    string `json:"path"`
	Level   string `json:"level"`   // low, medium, high
	Quality int    `json:"quality"` // webp quality
}

type compressRequest struct {
	Source     string   `json:"source"`
	Files      []string `json:"files"`
	Level      string   `json:"level"`
	Quality    int      `json:"quality"`
	Backup     bool     `json:"backup"`
	BackupPath string   `json:"backupPath"`
	BackupName string   `json:"backupName"`
}

type compressResponse struct {
	TaskID  string `json:"taskId"`
	Message string `json:"message"`
}

type progressEvent struct {
	Current   string `json:"current"`
	Processed int    `json:"processed"`
	Total     int    `json:"total"`
	Percent   int    `json:"percent"`
	Skipped   int    `json:"skipped"`
	Failed    int    `json:"failed"`
}

type finishEvent struct {
	Success        int    `json:"success"`
	Skipped        int    `json:"skipped"`
	Failed         int    `json:"failed"`
	BackupPath     string `json:"backupPath"`
	BackupFallback bool   `json:"backupFallback"`
}

type errorEvent struct {
	File  string `json:"file"`
	Error string `json:"error"`
}

// --- Compression Result ---

type compressionResult struct {
	Data    []byte
	Format  string // "webp" or "png"
	Skipped bool
	Err     error
}

// --- Progress Manager (for SSE) ---

type progressManager struct {
	mu     sync.Mutex
	chans  map[string]chan progressEvent
	dones  map[string]chan finishEvent
	errors map[string]chan errorEvent
}

var progressMgr = &progressManager{
	chans:  make(map[string]chan progressEvent),
	dones:  make(map[string]chan finishEvent),
	errors: make(map[string]chan errorEvent),
}

func (pm *progressManager) createTask(taskID string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.chans[taskID] = make(chan progressEvent, 100)
	pm.dones[taskID] = make(chan finishEvent, 1)
	pm.errors[taskID] = make(chan errorEvent, 100)
}

func (pm *progressManager) closeTask(taskID string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	if ch, ok := pm.chans[taskID]; ok {
		close(ch)
		delete(pm.chans, taskID)
	}
	if ch, ok := pm.dones[taskID]; ok {
		close(ch)
		delete(pm.dones, taskID)
	}
	if ch, ok := pm.errors[taskID]; ok {
		close(ch)
		delete(pm.errors, taskID)
	}
}

func (pm *progressManager) sendProgress(taskID string, evt progressEvent) {
	pm.mu.Lock()
	ch, ok := pm.chans[taskID]
	pm.mu.Unlock()
	if ok {
		select {
		case ch <- evt:
		default: // drop if buffer full
		}
	}
}

func (pm *progressManager) sendFinish(taskID string, evt finishEvent) {
	pm.mu.Lock()
	ch, ok := pm.dones[taskID]
	pm.mu.Unlock()
	if ok {
		select {
		case ch <- evt:
		default:
		}
	}
}

func (pm *progressManager) getProgressChan(taskID string) chan progressEvent {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	return pm.chans[taskID]
}

func (pm *progressManager) getDoneChan(taskID string) chan finishEvent {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	return pm.dones[taskID]
}

// --- Encoder Selection ---

func detectFormat(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg":
		return "jpeg"
	case ".png":
		return "png"
	case ".gif":
		return "gif"
	case ".bmp":
		return "bmp"
	case ".tif", ".tiff":
		return "tiff"
	case ".webp":
		return "webp"
	default:
		return "unknown"
	}
}

func isImageFile(path string) bool {
	f := detectFormat(path)
	return f != "unknown" && f != "gif" && f != "tiff"
}

func shouldResize(level string, width int) (int, bool) {
	switch level {
	case "medium":
		if width > 4000 {
			return 3000, true
		}
	case "high":
		if width > 3000 {
			return 2000, true
		}
	}
	return 0, false
}

func compressImage(srcPath string, level string, quality int) compressionResult {
	format := detectFormat(srcPath)

	// GIF: skip
	if format == "gif" {
		return compressionResult{Skipped: true}
	}

	// Read original file
	origData, err := os.ReadFile(srcPath)
	if err != nil {
		return compressionResult{Err: fmt.Errorf("read file: %w", err)}
	}

	// PNG + low tier: OxiPNG path
	if format == "png" && level == "low" {
		return compressWithOxiPNG(srcPath, origData)
	}

	// All other cases: WebP
	return compressWithWebP(origData, format, level, quality)
}

func compressWithWebP(origData []byte, format string, level string, quality int) compressionResult {
	// Decode image
	reader := bytes.NewReader(origData)
	img, _, err := image.Decode(reader)
	if err != nil {
		// Try PNG decoder explicitly
		reader2 := bytes.NewReader(origData)
		img, err = png.Decode(reader2)
		if err != nil {
			return compressionResult{Err: fmt.Errorf("decode image: %w", err)}
		}
	}

	// Optional resize
	if newW, doResize := shouldResize(level, img.Bounds().Dx()); doResize {
		img = resizeImage(img, newW)
	}

	// Encode to WebP
	var buf bytes.Buffer
	err = webp.Encode(&buf, img, &webp.EncoderOptions{
		Quality: float32(quality),
		Method:  4, // balanced speed/compression
	})
	if err != nil {
		return compressionResult{Err: fmt.Errorf("webp encode: %w", err)}
	}

	// Safety net: if compressed >= original, skip
	if buf.Len() >= len(origData) {
		return compressionResult{Skipped: true}
	}

	return compressionResult{
		Data:   buf.Bytes(),
		Format: "webp",
	}
}

func compressWithOxiPNG(srcPath string, origData []byte) compressionResult {
	// Write temp file for oxipng
	tmpFile, err := os.CreateTemp("", "oxipng-*.png")
	if err != nil {
		return compressionResult{Err: fmt.Errorf("create temp: %w", err)}
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	_, err = tmpFile.Write(origData)
	tmpFile.Close()
	if err != nil {
		return compressionResult{Err: fmt.Errorf("write temp: %w", err)}
	}

	// Optional: palette reduction using colorquant
	paletteReduced := tryPaletteReduction(tmpPath, origData)

	// Run oxipng -o2 --strip safe
	cmd := exec.Command("oxipng", "-o", "2", "--strip", "safe", tmpPath)
	err = cmd.Run()
	if err != nil {
		// If oxipng not available, return original as skipped
		if paletteReduced {
			result, _ := os.ReadFile(tmpPath)
			if len(result) > 0 && len(result) < len(origData) {
				return compressionResult{Data: result, Format: "png"}
			}
		}
		return compressionResult{Skipped: true}
	}

	// Read result
	result, err := os.ReadFile(tmpPath)
	if err != nil {
		return compressionResult{Err: fmt.Errorf("read result: %w", err)}
	}

	// Safety net
	if len(result) >= len(origData) {
		return compressionResult{Skipped: true}
	}

	return compressionResult{
		Data:   result,
		Format: "png",
	}
}

func tryPaletteReduction(path string, origData []byte) bool {
	reader := bytes.NewReader(origData)
	img, err := png.Decode(reader)
	if err != nil {
		return false
	}

	// Only reduce if >256 unique colors
	bounds := img.Bounds()
	if hasFewColors(img, bounds) {
		return false
	}

	// Create paletted image with 256 colors
	dst := image.NewPaletted(bounds, nil)
	colorquant.NoDither.Quantize(img, dst, 256, false, true)

	// Write back to file
	f, err := os.Create(path)
	if err != nil {
		return false
	}
	defer f.Close()
	png.Encode(f, dst)
	return true
}

func hasFewColors(img image.Image, bounds image.Rectangle) bool {
	colorSet := make(map[[4]int]bool)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			key := [4]int{int(r >> 8), int(g >> 8), int(b >> 8), int(a >> 8)}
			colorSet[key] = true
			if len(colorSet) > 256 {
				return false
			}
		}
	}
	return true
}

func resizeImage(src image.Image, newWidth int) image.Image {
	bounds := src.Bounds()
	oldWidth := bounds.Dx()
	if oldWidth <= newWidth {
		return src
	}
	ratio := float64(newWidth) / float64(oldWidth)
	newHeight := int(float64(bounds.Dy()) * ratio)
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	// Simple nearest-neighbor resize
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX := int(float64(x) / ratio)
			srcY := int(float64(y) / ratio)
			dst.Set(x, y, src.At(srcX, srcY))
		}
	}
	return dst
}

// --- Preview Handler ---

// @Summary Preview image compression
// @Description Compress a single image and return the result as a blob for preview
// @Tags compress
// @Accept json
// @Produce octet-stream
// @Param body body compressPreviewRequest true "Preview request"
// @Router /api/compress-images/preview [post]
func compressPreviewHandler(w http.ResponseWriter, r *http.Request, d *requestContext) (int, error) {
	if !d.user.Permissions.Admin {
		return http.StatusForbidden, fmt.Errorf("admin permission required for image compression")
	}
	var req compressPreviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return http.StatusBadRequest, fmt.Errorf("invalid request body: %v", err)
	}

	if req.Path == "" {
		return http.StatusBadRequest, fmt.Errorf("path is required")
	}

	// Resolve real path via source/user scope
	realPath, err := resolveCompressPath(req.Source, req.Path, d)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Read original
	origData, err := os.ReadFile(realPath)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to read file: %v", err)
	}

	// Compress
	result := compressImage(realPath, req.Level, req.Quality)
	if result.Err != nil {
		return http.StatusInternalServerError, result.Err
	}

	if result.Skipped {
		// Return original as-is
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("X-Original-Size", fmt.Sprintf("%d", len(origData)))
		w.Header().Set("X-Compressed-Size", fmt.Sprintf("%d", len(origData)))
		w.Header().Set("X-Skipped", "true")
		w.Write(origData)
		return 0, nil
	}

	// Set headers
	contentType := "image/webp"
	if result.Format == "png" {
		contentType = "image/png"
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("X-Original-Size", fmt.Sprintf("%d", len(origData)))
	w.Header().Set("X-Compressed-Size", fmt.Sprintf("%d", len(result.Data)))
	w.Write(result.Data)
	return 0, nil
}

// --- ZSTD Backup ---

func createBackup(filePaths []string, outputPath string) error {
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create backup file: %w", err)
	}
	defer outFile.Close()

	enc, err := zstd.NewWriter(outFile, zstd.WithEncoderLevel(zstd.SpeedDefault))
	if err != nil {
		return fmt.Errorf("create zstd encoder: %w", err)
	}
	defer enc.Close()

	tw := tar.NewWriter(enc)
	defer tw.Close()

	for _, filePath := range filePaths {
		err = addFileToTar(tw, filePath)
		if err != nil {
			return fmt.Errorf("add file %s to tar: %w", filePath, err)
		}
	}
	return nil
}

func addFileToTar(tw *tar.Writer, filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return filepath.Walk(filePath, func(walkPath string, walkInfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if walkInfo.IsDir() {
				return nil
			}
			return addSingleFileToTar(tw, walkPath)
		})
	}
	return addSingleFileToTar(tw, filePath)
}

func addSingleFileToTar(tw *tar.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	header.Name = filePath

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	_, err = io.Copy(tw, file)
	return err
}

// --- Execute Handler ---

// @Summary Compress images
// @Description Start batch image compression with optional ZSTD backup
// @Tags compress
// @Accept json
// @Produce json
// @Param body body compressRequest true "Compress request"
// @Router /api/compress-images [post]
func compressHandler(w http.ResponseWriter, r *http.Request, d *requestContext) (int, error) {
	if !d.user.Permissions.Admin {
		return http.StatusForbidden, fmt.Errorf("admin permission required for image compression")
	}
	var req compressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return http.StatusBadRequest, fmt.Errorf("invalid request body: %v", err)
	}

	if len(req.Files) == 0 {
		return http.StatusBadRequest, fmt.Errorf("no files specified")
	}

	// Resolve all file paths to real paths
	realPaths := make([]string, 0, len(req.Files))
	for _, userPath := range req.Files {
		realPath, err := resolveCompressPath(req.Source, userPath, d)
		if err != nil {
			logger.Errorf("compress: failed to resolve path %s: %v", userPath, err)
			continue
		}
		info, err := os.Stat(realPath)
		if err != nil {
			logger.Errorf("compress: stat failed for %s: %v", realPath, err)
			continue
		}
		if info.IsDir() {
			filepath.Walk(realPath, func(walkPath string, walkInfo os.FileInfo, err error) error {
				if err != nil || walkInfo.IsDir() {
					return nil
				}
				if isImageFile(walkPath) {
					realPaths = append(realPaths, walkPath)
				}
				return nil
			})
		} else {
			realPaths = append(realPaths, realPath)
		}
	}

	if len(realPaths) == 0 {
		return http.StatusBadRequest, fmt.Errorf("no valid files to compress")
	}

	// Resolve backup path to real filesystem path before goroutine
	var realBackupPath string
	var sourceRootPath string
	if req.Backup && req.BackupPath != "" {
		var errBackup error
		realBackupPath, errBackup = resolveCompressPath(req.Source, req.BackupPath, d)
		if errBackup != nil {
			return http.StatusBadRequest, fmt.Errorf("failed to resolve backup path: %v", errBackup)
		}
		sourceRootPath, _ = resolveCompressPath(req.Source, "/", d)
	}

	taskID := fmt.Sprintf("compress-%d", time.Now().UnixNano())
	progressMgr.createTask(taskID)

	// Run in background
	go func() {
		defer progressMgr.closeTask(taskID)

		backupPathDisplay := ""
		backupFallback := false
		// Step 1: ZSTD backup with 3-level fallback
		if req.Backup {
			// Level 1: same-level directory (original backupPath)
			realBackupFull := filepath.Join(realBackupPath, req.BackupName)
			err := createBackup(realPaths, realBackupFull)

			if err != nil {
				// Level 2: try parent directory
				parentDir := filepath.Dir(realBackupPath)
				realBackupFull = filepath.Join(parentDir, req.BackupName)
				err = createBackup(realPaths, realBackupFull)
				backupFallback = true
			}

			if err != nil && sourceRootPath != "" {
				// Level 3: try source root
				realBackupFull = filepath.Join(sourceRootPath, req.BackupName)
				err = createBackup(realPaths, realBackupFull)
				backupFallback = true
			}

			if err != nil {
				// All 3 levels failed - abort compression
				logger.Errorf("compress: backup failed at all levels: %v", err)
				progressMgr.sendFinish(taskID, finishEvent{
					Failed:     len(realPaths),
					BackupPath: "",
				})
				return
			}

			backupPathDisplay = filepath.Join(req.BackupPath, req.BackupName)
		}

		// Step 2: Serial compression
		success, skipped, failed := 0, 0, 0
		total := len(realPaths)

		for i, filePath := range realPaths {
			progressMgr.sendProgress(taskID, progressEvent{
				Current:   filepath.Base(filePath),
				Processed: i,
				Total:     total,
				Percent:   int(float64(i) / float64(total) * 100),
				Skipped:   skipped,
				Failed:    failed,
			})

			result := compressImage(filePath, req.Level, req.Quality)
			if result.Err != nil {
				failed++
				continue
			}
			if result.Skipped {
				skipped++
				continue
			}

			// Write compressed file
			newPath := getCompressedPath(filePath, result.Format)
			err := os.WriteFile(newPath, result.Data, 0644)
			if err != nil {
				failed++
				continue
			}

			// Remove original if path changed
			if newPath != filePath {
				os.Remove(filePath)
			}
			success++
		}

		// Send finish event
		progressMgr.sendFinish(taskID, finishEvent{
			Success:        success,
			Skipped:        skipped,
			Failed:         failed,
			BackupPath:     backupPathDisplay,
			BackupFallback: backupFallback,
		})
	}()

	// Return task ID immediately
	return renderJSON(w, r, compressResponse{
		TaskID:  taskID,
		Message: "Compression started",
	})
}

func getCompressedPath(origPath string, format string) string {
	if format == "png" {
		return origPath // PNG stays as PNG (OxiPNG)
	}
	// Change extension to .webp
	ext := filepath.Ext(origPath)
	if ext == "" {
		return origPath + ".webp"
	}
	return origPath[:len(origPath)-len(ext)] + ".webp"
}

// --- SSE Progress Handler ---

// @Summary Compress progress stream
// @Description Server-Sent Events stream for compression progress
// @Tags compress
// @Param taskId query string true "Task ID"
// @Router /api/compress-images/progress [get]
func compressProgressHandler(w http.ResponseWriter, r *http.Request, d *requestContext) (int, error) {
	if !d.user.Permissions.Admin {
		return http.StatusForbidden, fmt.Errorf("admin permission required for image compression")
	}
	taskID := r.URL.Query().Get("taskId")
	if taskID == "" {
		return http.StatusBadRequest, fmt.Errorf("missing taskId")
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		return http.StatusInternalServerError, fmt.Errorf("streaming not supported")
	}

	ctx := r.Context()

	for {
		select {
		case <-ctx.Done():
			return 0, nil
		case evt, ok := <-progressMgr.getProgressChan(taskID):
			if !ok {
				return 0, nil
			}
			data, _ := json.Marshal(evt)
			fmt.Fprintf(w, "event: progress\ndata: %s\n\n", data)
			flusher.Flush()
		case evt, ok := <-progressMgr.getDoneChan(taskID):
			if !ok {
				return 0, nil
			}
			data, _ := json.Marshal(evt)
			fmt.Fprintf(w, "event: finish\ndata: %s\n\n", data)
			flusher.Flush()
			return 0, nil
		}
	}
}

// --- Helpers ---

// resolveCompressPath resolves a user-facing source+path to a real filesystem path.
func resolveCompressPath(source, userPath string, d *requestContext) (string, error) {
	idx := indexing.GetIndex(source)
	if idx == nil {
		return "", fmt.Errorf("source %s not found", source)
	}

	userscope, err := d.user.GetScopeForSourceName(source)
	if err != nil {
		return "", err
	}

	fullPath := filepath.Join(userscope, userPath)

	// Access control check
	if store.Access != nil && !store.Access.Permitted(idx.Path, fullPath, d.user.Username) {
		return "", fmt.Errorf("access denied to %s", userPath)
	}

	realPath, _, err := idx.GetRealPath(fullPath)
	if err != nil {
		return "", fmt.Errorf("resolve path: %w", err)
	}

	return realPath, nil
}
