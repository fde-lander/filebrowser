# Task Plan: FileBrowser Quantum - v1.4.0.6-fde

## Goal

Maintain FileBrowser Quantum FDE fork. v1.4.0.5-fde-stable is the only usable release (README rewrite in progress). v1.4.0.6 and v1.4.0.7 are VOIDED - future upgrades will skip them entirely.

## Current Phase

Phase 5: README Rewrite for Stable Release

## Phases

<!-- CURRENT PHASES MARKER - Insert new phases BELOW (Do not modify) -->

### Phase 5: README Rewrite for Stable Release
- [x] 5.1 Verify release state and image content by layer scan
- [x] 5.2 Draft Traditional Chinese README with FDE signature only
- [x] 5.3 Add FDE Fork intro and dated version history section
- [x] 5.4 Update deploy steps to latest stable release
- [x] 5.5 Document release asset naming situation
- [x] 5.6 Document voided versions with recommended wording
- [x] 5.7 Self-review for stale version references
- **Status:** in_progress

### Phase 6: Deploy Doc and Latest Release
- [ ] 6.1 Write complete Traditional Chinese DEPLOY.md for server-side build flow
- [ ] 6.2 Add prominent deployment model notice section to README
- [ ] 6.3 Verify build inputs and prerequisites are documented
- [ ] 6.4 Commit README and DEPLOY docs to git
- [ ] 6.5 Plan git tag and GitHub release operations for MASTER approval
- **Status:** pending

## Archived phase

<!-- ARCHIVED PHASES MARKER - Move completed phases BELOW (Do not modify) -->

### Archived phase 1: Deep Code Research
- [x] 1.1 Wave 1: 3 subagents dispatched + returned
- [x] 1.2 Agent A: Prompts.vue -> CompressImages lifecycle (Bug B root cause)
- [x] 1.3 Agent B: ExtendedImage.vue -> touch/click event bindings (Bug G/H root cause)
- [x] 1.4 Agent C: compress.go + httpRouter.go + compress.js -> queueMgr + worker + routes (Ch8)
- [x] 1.5 All findings verified by main agent via code reading
- **Status:** complete

### Archived phase 2: Brainstorming + Design Discussion
- [x] 2.1 6 design topics discussed with MASTER, all confirmed:
- [x] 2.2 1. Bug B: Method C (two-layer: global status bar + dialog detail)
- [x] 2.3 2. Bug G: Conditional preventDefault (preserve swipe + transition)
- [x] 2.4 3. Queue progress: Cumulative totalFiles/totalProcessed + batch info
- [x] 2.5 4. Cancel scope: Cancel entire queue (not current batch)
- [x] 2.6 5. Skip current batch: Added to v1.4.0.6 scope
- [x] 2.7 6. Pause auto-timeout: Toggle + numeric input, default 30min, cross-session persisted
- [x] 2.8 Skip/Cancel require confirmation dialog (showPrompt system)
- [x] 2.9 Only 1 batch: hide skip button
- **Status:** complete

### Archived phase 3: Design Spec Writing
- [x] 3.1 10 chapters, chapter-by-chapter PATCH
- **Status:** complete

### Archived phase 4: Implementation
- [x] 4.1 Execution mode: INLINE (主 Agent 串行执行)
- [x] 4.2 Skills loaded: executing-plans, TDD, PWF
- [x] 4.3 Plan: ~/.hermes/docs/superpowers/plans/2026-07-15-v1.4.0.6-fde-plan.md
- [x] 4.4 Each Task = 1 git commit, individually rollbackable
- [x] 4.5 All patches < 2.5KB
- **Status:** complete

## Key Questions

1. Phase 4 gap: original v1.4.0.6 plan numbered phases 1/2/3/5 (Phase 4 was Round 1 Build & Deploy, superseded by Phase 5 Implementation which includes deployment logic). Not a data loss - verify against git history d7008dbf
## Decisions Made

| Decision | Rationale |
|----------|-----------|
| Bug B: Method C two-layer (global status bar + dialog detail) | Design confirmed with MASTER during brainstorming |
| Bug G: Remove .prevent from @touchmove, conditional preventDefault in touchMove/touchEnd | Preserve swipe + transition |
| Bug H: Ensure nav button calls transition path (nextPrevious.vue investigation needed) | Design confirmed with MASTER |
| Bug I: @media (max-width: 768px) responsive CSS for preview overlay | Design confirmed with MASTER |
| Queue List API: Leverage existing CompressJobStatus.Queue field | Design confirmed with MASTER |
| Cumulative progress: totalFiles/totalProcessed/batchCount/currentBatchIndex | Design confirmed with MASTER |
| Cancel = entire queue; Skip = current batch only (with batch > 1 check) | Design confirmed with MASTER |
| Pause/Resume: sync.Cond on queueMgr | Design confirmed with MASTER |
| Skip/Cancel: secondary confirmation via showPrompt system | Design confirmed with MASTER |
| Pause auto-timeout: Toggle + numeric (5-120min, default 30), persisted 4-touchpoint | Design confirmed with MASTER |
| Only 1 batch: hide skip button | Design confirmed with MASTER |

## Errors Encountered

| Error | Attempt | Resolution |
|-------|---------|------------|

## Notes
- Commit Group 1: Global status bar component (Bug B solution)
- Commit Group 2: compress.go backend flags + sync.Cond + cumulative stats + timeout goroutine
- Commit Group 3: httpRouter.go new route registration (pause/resume/cancel/skip/queue)
- Commit Group 4: compress.js frontend API functions
- Commit Group 5: CompressImages.vue control buttons + detail progress + confirm dialog
- Commit Group 6: CompressImages.vue preview layout mobile CSS (Bug I)
- Commit Group 7: ExtendedImage.vue Bug G fix (touchmove.prevent conditional)
- Commit Group 8: Bug H fix (nextPrevious.vue, pending investigation)
- Commit Group 9: Settings system + i18n (compressPauseTimeout)
- Key Constraint: Transition engine LOCKED at v1.4.0.4-hotfix (do NOT modify)
- Key Constraint: DOUBT principle - subagent research -> verify -> design -> approve
- Key Constraint: Legacy CSS - annotate only, do NOT delete
- Key Constraint: Docker - tar only, never push to Hub
- Key Constraint: No local testing (Docker build -> save -> scp -> deploy)
- Key Constraint: i18n 3 files must sync (en/zh-cn/zh-tw)
- Key Constraint: HARD-GATE - no sudo without MASTER approval
