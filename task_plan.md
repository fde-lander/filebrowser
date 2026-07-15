# Task Plan: FileBrowser Quantum - v1.4.0.6-fde

## Goal

v1.4.0.6-fde: Mobile bug fixes (G/H/I) + Compress control features (pause/continue/cancel/skip) + Bug B re-investigation + Queue List API + cumulative progress + pause auto-timeout.
All modules are independent commits with individual rollback capability.

## Current Phase

Phase 5: Implementation (IN PROGRESS - inline execution)

## Phases

### Phase 1: Deep Code Research (COMPLETED)
- Wave 1: 3 subagents dispatched + returned
- Agent A: Prompts.vue -> CompressImages lifecycle (Bug B root cause)
- Agent B: ExtendedImage.vue -> touch/click event bindings (Bug G/H root cause)
- Agent C: compress.go + httpRouter.go + compress.js -> queueMgr + worker + routes (Ch8)
- All findings verified by main agent via code reading
- Status: completed

### Phase 2: Brainstorming + Design Discussion (COMPLETED)
- 6 design topics discussed with MASTER, all confirmed:
  1. Bug B: Method C (two-layer: global status bar + dialog detail)
  2. Bug G: Conditional preventDefault (preserve swipe + transition)
  3. Queue progress: Cumulative totalFiles/totalProcessed + batch info
  4. Cancel scope: Cancel entire queue (not current batch)
  5. Skip current batch: Added to v1.4.0.6 scope
  6. Pause auto-timeout: Toggle + numeric input, default 30min, cross-session persisted
- Skip/Cancel require confirmation dialog (showPrompt system)
- Only 1 batch: hide skip button
- Status: completed

### Phase 3: Design Spec Writing (IN PROGRESS)
- 10 chapters, chapter-by-chapter PATCH
- Status: in_progress

### Phase 5: Implementation (IN PROGRESS - inline execution)
- Execution mode: INLINE (主 Agent 串行执行)
- Skills loaded: executing-plans, TDD, PWF
- Plan: ~/.hermes/docs/superpowers/plans/2026-07-15-v1.4.0.6-fde-plan.md
- Each Task = 1 git commit, individually rollbackable
- All patches < 2.5KB
- Status: in_progress

## Confirmed Design Decisions

1. Bug B: Method C two-layer (global status bar + dialog detail)
2. Bug G: Remove .prevent from @touchmove, conditional preventDefault in touchMove/touchEnd
3. Bug H: Ensure nav button calls transition path (nextPrevious.vue investigation needed)
4. Bug I: @media (max-width: 768px) responsive CSS for preview overlay
5. Queue List API: Leverage existing CompressJobStatus.Queue field
6. Cumulative progress: totalFiles/totalProcessed/batchCount/currentBatchIndex
7. Cancel = entire queue; Skip = current batch only (with batch > 1 check)
8. Pause/Resume: sync.Cond on queueMgr
9. Skip/Cancel: secondary confirmation via showPrompt system
10. Pause auto-timeout: Toggle + numeric (5-120min, default 30), persisted 4-touchpoint
11. Only 1 batch: hide skip button

## Commit Groups (9 independent commits)

1. Global status bar component (Bug B solution)
2. compress.go: backend flags + sync.Cond + cumulative stats + timeout goroutine
3. httpRouter.go: new route registration (pause/resume/cancel/skip/queue)
4. compress.js: frontend API functions
5. CompressImages.vue: control buttons + detail progress + confirm dialog
6. CompressImages.vue: preview layout mobile CSS (Bug I)
7. ExtendedImage.vue: Bug G fix (touchmove.prevent conditional)
8. Bug H fix (nextPrevious.vue, pending investigation)
9. Settings system + i18n (compressPauseTimeout)

## Key Constraints

- Transition engine LOCKED at v1.4.0.4-hotfix (do NOT modify)
- DOUBT principle: subagent research -> verify -> design -> approve
- Legacy CSS: annotate only, do NOT delete
- Docker: tar only, never push to Hub
- No local testing (Docker build -> save -> scp -> deploy)
- i18n: 3 files must sync (en/zh-cn/zh-tw)
- HARD-GATE: no sudo without MASTER approval
