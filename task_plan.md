# Task Plan: FileBrowser Quantum - v1.5.6.1-fde

## Goal

Maintain FileBrowser Quantum FDE fork. v1.5.6.1-fde branch ready (merged upstream v1.5.6-stable, all verification PASS). DEPLOYMENT MODEL: LOCAL docker build then docker save then tar.zst level 5 compression then transfer to remote server then docker load. Remote server does NOT build. Full rollback = previous git tag + backup data dir.

## Current Phase

Phase 15: Void v1.5.6.1-fde and Restore v1.4.0.5-fde-hotfix

## Phases

<!-- CURRENT PHASES MARKER - Insert new phases BELOW (Do not modify) -->

### Phase 15: Void v1.5.6.1-fde and Restore v1.4.0.5-fde-hotfix
- [ ] 15.1 Switch to main branch
- [ ] 15.2 Delete local v1.5.6.1-fde branch
- [ ] 15.3 Delete remote v1.5.6.1-fde branch
- [ ] 15.4 Remove docker image filebrowser-fde:v1.5.6.1-fde
- [ ] 15.5 Remove tar.zst deliverable file
- [ ] 15.6 Update README: void v1.5.6.1-fde, restore v1.4.0.5-fde-hotfix as LATEST
- [ ] 15.7 Update DEPLOY.md back to v1.4.0.5-fde-hotfix baseline
- [ ] 15.8 Commit and push README + DEPLOY.md changes
- [ ] 15.9 Record lesson to LanceDB: code-level verification != functional verification
- [ ] 15.10 Update progress.md with void session log
- [ ] 15.11 Pause project
- **Status:** in_progress

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

### Archived phase 5: README Rewrite for Stable Release
- [x] 5.1 Verify release state and image content by layer scan
- [x] 5.2 Draft Traditional Chinese README with FDE signature only
- [x] 5.3 Add FDE Fork intro and dated version history section
- [x] 5.4 Update deploy steps to latest stable release
- [x] 5.5 Document release asset naming situation
- [x] 5.6 Document voided versions with recommended wording
- [x] 5.7 Self-review for stale version references
- **Status:** complete

### Archived phase 6: Deploy Doc and Latest Release
- [x] 6.1 Write complete Traditional Chinese DEPLOY.md for server-side build flow
- [x] 6.2 Add prominent deployment model notice section to README
- [x] 6.3 Verify build inputs and prerequisites are documented
- [x] 6.4 Commit README and DEPLOY docs to git
- [x] 6.5 Plan git tag and GitHub release operations for MASTER approval
- **Status:** complete

### Archived phase 7: Upstream Upgrade Scout (v1.5.6-stable)
- [x] 7.1 Add upstream remote + fetch latest tags (confirm v1.5.6-stable target)
- [x] 7.2 Measure divergence: merge-base, commit ranges, FF feasibility
- [x] 7.3 Quantify conflict surface via merge-tree (files + hunks)
- [x] 7.4 Map GHSA security fixes + assess upgrade value
- [x] 7.5 Confirm fork core assets untouched by upstream
- [x] 7.6 Write research doc + PWF records (findings/progress)
- [x] 7.7 Handoff + pause project for next AI group
- **Status:** complete

### Archived phase 8: Conflict deep-dive and pinned ruling
- [x] 8.1 Diff pinnedItems.go handler fork vs upstream
- [x] 8.2 Analyze files.go symlink refactor vs ShowPinnedItems branch
- [x] 8.3 Analyze mutations.js togglePinnedItem vs setObjectProperty deep copy
- [x] 8.4 Analyze ContextMenu.vue 3-hunk conflict
- [x] 8.5 Map i18n key collision across en/zh-cn/zh-tw
- [x] 8.6 Assess go.mod/go.sum dependency conflict
- [x] 8.7 Assess bolt DB compat for pinned data
- [x] 8.8 Present ruling summary to MASTER for approval
- **Status:** complete

### Archived phase 9: Upgrade design spec
- [x] 9.1 Write upgrade design spec doc (segmented)
- [x] 9.2 Self-review spec for placeholder, consistency, scope
- [x] 9.3 Confirm full conflict ruling matrix with 20 files
- [x] 9.4 Verify no missed fork assets beyond 7 plus users.go
- **Status:** complete

### Archived phase 10: Open branch and merge upstream
- [x] 10.1 Open v1.5.6.1-fde branch from fork baseline
- [x] 10.2 git merge v1.5.6-stable while keeping conflicts
- **Status:** complete

### Archived phase 11: Resolve conflicts and preserve assets
- [x] 11.1 Resolve core asset conflicts as ours
- [x] 11.2 Resolve security fix conflicts as upstream
- [x] 11.3 Resolve pinned handler per ruling from deep-dive stage
- [x] 11.4 Resolve i18n conflicts keeping both key sets
- [x] 11.5 Resolve go.mod go.sum with upstream deps
- [x] 11.6 Resolve remaining files per scout report
- **Status:** complete

### Archived phase 12: Build verification go and frontend
- [x] 12.1 Run go build plus go vet
- [x] 12.2 Run go mod verify plus tidy
- [x] 12.3 Run pnpm build frontend
- [x] 12.4 Verify 7 core assets behavior preserved
- [x] 12.5 Verify 9 GHSA patches present in tree
- **Status:** complete

### Archived phase 13: Deploy doc and version history update
- [x] 13.1 Rewrite DEPLOY.md for remote build flow with backup, git update, build, run, rollback
- [x] 13.2 Update README version history
- [x] 13.3 Commit all docs to git
- **Status:** complete

### Archived phase 14: Local Docker Build and tar.zst Delivery
- [x] 14.1 Remove debug image filebrowser-fde:v1.5.6.1-fde-debug
- [x] 14.2 Remove stale build cache (docker builder prune)
- [x] 14.3 Local docker build with correct VERSION and REVISION args
- [x] 14.4 Verify image: version output, sanity checks (ffmpeg/ffprobe/exiftool)
- [x] 14.5 Docker save to tar then compress with zstd level 5
- [x] 14.6 Verify tar.zst integrity (zstd -t)
- [x] 14.7 Rewrite DEPLOY.md for tar.zst transfer + docker load model
- [x] 14.8 Update .gitignore for tar.zst pattern
- [x] 14.9 Commit DEPLOY.md + .gitignore changes to git
- [x] 14.10 Clean build cache after successful build (free disk space)
- [x] 14.11 Report deliverable path, size, sha256, and deployment instructions to MASTER
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
| BUILD policy: LOCAL docker build then docker save then tar.zst level 5 then transfer to remote then docker load. Remote server does NOT build. | MASTER directive 2026-09-30: remote build on 2G RAM server FAILED twice (OOM even with 8G swap). Remote build is architecturally infeasible for CGO Go + npm+vite. Switched to local build + tar.zst delivery model. |
| Upgrade strategy: Option C (three-way merge). Open v1.5.6.1-fde from fork baseline, git merge v1.5.6-stable. Conflict rules: 7 core assets=ours, security fixes=upstream, pinned per diff ruling. | MASTER approved Option C 2026-09-29. Preserves fork patches + all 9 GHSA security fixes in one merge. |

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
- Key Constraint: LOCAL docker build is the build method. Build locally, save as tar.zst, transfer to remote, docker load on remote. Remote server does NOT build.
- Key Constraint: Deployment = transfer tar.zst to remote, docker load, docker compose up. NO remote git clone or docker build.
- Key Constraint: i18n 3 files must sync (en/zh-cn/zh-tw)
- Key Constraint: HARD-GATE - no sudo without MASTER approval
- Key Constraint: BUILD policy = LOCAL docker build + docker save + tar.zst level 5. Deliverable = filebrowser-fde-v1.5.6.1-fde.tar.zst. Remote server: docker load + docker compose up. Rollback = previous image tag + backup data dir.
- DEPLOY.md must document: tar.zst transfer, docker load, data dir backup, config.yaml setup, docker-compose.yaml, verification, update flow, full rollback mechanism
