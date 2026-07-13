# Task Plan: FileBrowser Quantum - Custom Patches (Round 2)

## Goal

Two major feature upgrades:
1. **Image Viewer Enhancement** - preload + transition + tap navigation + persistent mobile buttons
2. **Image Compression** - right-click batch compress with preview + ZSTD backup

## Current Phase

Phase 3 (v1.4.0.4 Hotfix - Deep Code Research + Design)

## Project Context

- **Base Version:** v1.4.0-stable (locked, no upstream changes)
- **Fork URL:** https://github.com/fde-lander/filebrowser.git
- **Project Path:** /home/hermes/workspace/filebrowser-fde
- **Branch:** v1.4.0.2-image-viewer-compression (7 commits: 6 feature + 1 hotfix)
- **Previous Round:** Extract-to-folder feature complete (5 commits, 9 files, 454 insertions)
- **Docker image:** filebrowser-fde:v1.4.0.2 (84MB tar at project root)
- **Go binary:** /usr/local/go/bin/go (NOT in default PATH)
- **Subagent delegation:** config fixed (provider: custom:iturn), working

## Phases

### Phase 1: Design & Brainstorming (COMPLETED)
- [x] Subagent code research: image viewer frontend (3 agents parallel)
- [x] Subagent code research: image compression backend possibilities
- [x] Clarify requirements with user
- [x] Propose approaches for each feature
- [x] Write design spec document (789 lines, 8 chapters)
- [x] User reviews spec
- [x] Implementation plan written (1939 lines, 12 tasks)
- **Status:** completed

### Phase 1.5: Implementation Wave 1-4 (COMPLETED)
- [x] Wave 1: 3 subagents parallel (Task 1+2 / 6+7 / 9+11) - ALL PASS
- [x] Wave 2: 1 subagent (Task 3+4+5 backend compress.go) - PASS
- [x] Wave 3: 1 subagent (Task 8+10 frontend integration) - PASS
- [x] Wave 4: Main agent (Task 12 Dockerfile + verification) - PASS
- [x] Go build + vet + mod verify: ALL PASS
- [x] Docker build + save: PASS (84MB)
- **Status:** completed

### Phase 1.6: Hotfix 1 (COMPLETED)
- [x] Deploy + test by MASTER
- [x] Bug 1 fix: i18n key paths (settings.* -> profileSettings.*)
- [x] Bug 2 fix: imagePreload/imageTapNav default true in Profile.vue
- [x] Bug 3 partial fix: src watcher calls doTransition() + Preview.vue no isTransitioning
- [x] Bug 4 fix: CompressImages.vue items prop + mounted() null guard
- **Status:** completed (commit ad5f2774)

### Phase 2: Bug Fix & Permission Enhancement (COMPLETED)
- [x] Subagent parallel analysis: Bug A/B/C/D + permission system
- [x] Brainstorming: design fix approaches + permission gate design
- [x] Design spec + implementation plan
- [x] Implementation (main agent inline, 5 commits)
- [x] Build + deploy + test by MASTER
- **Status:** completed (v1.4.0.3 deployed)

### Phase 3: v1.4.0.4 Hotfix (CURRENT)
- [x] Deep code research: 3 subagents parallel
- [x] Brainstorming: design fix approaches for 5 remaining issues
- [x] Present design to user, get approval
- [x] Write design spec (8 chapters, 911 lines)
- [x] Implementation (main agent inline, 18 tasks, 5 phases)
- [x] Build + deploy + test (Docker image built + saved)
- **Status:** COMPLETED (v1.4.0.4 hotfix ready for deployment)

### Phase 3: Build & Deploy
- [ ] Go build + Docker build
- [ ] Deploy + test
- **Status:** pending

## Key Decisions

- Locked to v1.4.0-stable base (user confirmed no upstream changes needed)
- Main branch is default (merged from custom branch, clean state)
- AVIF excluded: encoding >120s/image, 2.5GB RAM (not feasible on 4GB server)
- All compression uses WebP except low-tier PNG uses OxiPNG CLI
- GIF files always skipped (animation too complex)
- Safety net: if compressed >= original size, keep original
- ZSTD backup: tar.zst format (tar wrapping zstd, NOT just .zst)
- Settings persist via NonAdminEditable struct
- Go bool zero-value = false: imagePreload/imageTapNav need frontend ?? true
- i18n keys for image viewer settings live under profileSettings.* namespace
- Preview.vue must NOT set isTransitioning=true for image navigation
- ExtendedImage.vue src watcher must call doTransition() not loadFullImage()
- CompressImages.vue requires items prop passed via showPrompt({ props: { items } })
- Subagent delegation provider: custom:iturn (format: custom:<name>)
- Main agent does critical code fixes, subagent only for research/deterministic tasks
- NEW: Permission gate required for extract-to-folder + compress-images features

## Notes

- User requirements: see findings.md for detailed feature requirements
- Use subagents (max 3 parallel) for code research
- Design phase only - do NOT implement until user explicitly says so
- NEW requirement: permission control for extract + compress features
- Handoff document: /tmp/handoff-r6jEEE.md
- LanceDB memory: ID d96afe8c (bug status), ID 49271d6b (project facts), ID f13604b0 (v1.4.0.3 results)
- Branch: v1.4.0.4-hotfix (created from v1.4.0.2-image-viewer-compression @ 3d973404)
