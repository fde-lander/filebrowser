# Findings: v1.4.0.6-fde Code Research

## Wave 1: COMPLETE (3 subagents + main agent verification)

### Agent A: Prompts.vue -> CompressImages lifecycle (Bug B)

**Key Finding: mounted() DOES fire on every dialog reopen**

- L2: v-for + :key=prompt.id -> each open creates new ID -> new component instance
- L71-74: <component :is="prompt.name"> dynamic rendering
- NO <keep-alive> anywhere in 991 lines (confirmed by main agent)
- Close = splice from state.prompts -> Vue destroys component
- Reopen = new prompt entry + new ID -> Vue creates new instance -> mounted() fires

**Bug B real root cause: CompressImages.vue mounted() (L309-318) has NO status check**
- L311: if no items -> return (no backend check)
- L315-317: only reads backup pref + expandItems()
- No pollStatus() call, no recoverQueueStatus()
- Previous AI's fix was reverted (commit 079cbeda), current code has zero recovery logic

**Conclusion: mounted() approach is viable. Add pollStatus() check in mounted().**

### Agent B: ExtendedImage.vue -> touch/click events (Bug G/H)

**Bug G root cause: @touchmove.prevent + touchEnd.preventDefault() double suppression**

- L2: @touchmove.prevent="touchMove" -> .prevent modifier = auto preventDefault() on ALL touchmove
- L1049-1054: touchEnd() calls event.preventDefault() when scale===1 + changedTouches>0
- Mobile tap sequence: touchstart -> touchmove (1px micro-move) -> touchend -> click
- touchmove.prevent kills synthetic click -> handleImageClick (L555) never runs
- Swipe works because it uses touch events directly (touchStart/touchMove/touchEnd -> finishEdgeGesture)

**handleImageClick (L555-586):**
- Guards: imageTapNavEnabled, scale > 1, wasDragging
- 200ms double-tap guard via tapNavTimeout
- Position: clientX - rect.left, < 0.4 = prev, > 0.6 = next
- Emits 'navigate-previous' or 'navigate-next'

**Transition path:**
- src watcher (L1061) -> navigateToImage (L455) -> waitForDecode -> swapBuffers (L481)
- swapBuffers: sets opacity=0 then display=block before rAF -> potential black flash if not decoded

**Bug H: Nav buttons are in parent component (nextPrevious.vue), not ExtendedImage.vue**
- Need to investigate nextPrevious.vue event bindings
- Nav button click -> emit navigate -> parent changes src -> transition path fires

**No device detection anywhere. Zero navigator.maxTouchPoints checks.**

### Agent C: compress.go + httpRouter.go + compress.js -> queueMgr (Ch8)

**queueMgr struct (L90-96):**
- mu sync.RWMutex
- queue []QueueItem
- current *QueueItem
- status CompressJobStatus
- workerActive bool

**QueueItem (L62-74):** ID, Files, Level, Quality, Source, Backup, BackupPath, BackupName, SourceRoot, Status, AddedAt

**CompressJobStatus (L76-88):** Status, CurrentFile, Processed, Total, Skipped, Failed, SavedBytes, BackupPath, BackupFallback, QueueLength, Queue []QueueItem

**compressWorker (L653-737):**
- L654: setWorkerActive(true), L655: defer setWorkerActive(false)
- L657: for {} infinite loop
- L658: dequeue() -> nil = return
- L697-733: for i, filePath := range item.Files (inner file loop)
- L735: finishCurrent()
- NO pause/cancel/stop mechanism exists

**setWorkerActive (L161-174):** if deactivating + queue empty, only sets "idle" if status is "running" (Bug C fix preserves "completed")

**Routes (httpRouter.go L147-149):**
- POST /compress-images/preview (withUser)
- POST /compress-images (withUser)
- GET /compress-images/status (withUser)
- Pattern: api.HandleFunc("METHOD /path", withUser(handler))

**compress.js (98 lines):**
- previewCompress() L15-44
- startCompress() L59-84
- pollStatus() L91-98 -> GET compress-images/status, returns response.json()

## Persistence System (verified by main agent)

**compressBackup 4-touchpoint pattern:**
1. Backend: users.go L174 -> CompressBackup bool
2. Frontend state: state.js L67 -> compressBackup: false
3. Auto-persist: mutations.js L606 -> "compressBackup" in allowlist
4. Read: CompressImages.vue L315 -> state.user?.compressBackup ?? true

**NOTE: compressBackup is NOT in Profile.vue usersApi.update field list (L450-469)**
- It relies solely on auto-persist allowlist mechanism
- New setting compressPauseTimeout should follow same pattern

**showPrompt system (mutations.js L327-368):**
- Already has confirm/callback mechanism
- showPrompt({ name, confirm: (result) => {}, props: {} })
- closeTopPrompt (L310-326) closes top
- download.js L70-78 has working example
- NO need to build new ConfirmDialog component - just register new prompt component
