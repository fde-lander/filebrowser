# Task Plan: FileBrowser Quantum - Custom Patches (Round 2)

## Goal

Two major feature upgrades:
1. **Image Viewer Enhancement** - preload + transition + tap navigation + persistent mobile buttons
2. **Image Compression** - right-click batch compress with preview + ZSTD backup

## Current Phase

Phase 1 (Design & Brainstorming - in progress)

## Project Context

- **Base Version:** v1.4.0-stable (locked, no upstream changes)
- **Fork URL:** https://github.com/fde-lander/filebrowser.git
- **Project Path:** /home/hermes/workspace/filebrowser-fde
- **Branch:** main (default, contains Round 1 extract-to-folder patch)
- **Previous Round:** Extract-to-folder feature complete (5 commits, 9 files, 454 insertions)

## Phases

### Phase 1: Design & Brainstorming (CURRENT)
- [ ] Subagent code research: image viewer frontend (3 agents parallel)
- [ ] Subagent code research: image compression backend possibilities
- [ ] Clarify requirements with user
- [ ] Propose approaches for each feature
- [ ] Write design spec document
- [ ] User reviews spec
- **Status:** in_progress

### Phase 2: Writing-Plans
- [ ] Invoke writing-plans skill
- [ ] Break spec into implementation tasks
- **Status:** pending

### Phase 3: Implementation
- [ ] Backend changes
- [ ] Frontend changes
- **Status:** pending

### Phase 4: Build & Deploy
- [ ] Go build + Docker build
- [ ] Deploy + test
- **Status:** pending

### Phase 5: Delivery
- [ ] Documentation + LanceDB
- **Status:** pending

## Key Decisions

| Decision | Rationale |
|----------|-----------|
| Locked to v1.4.0-stable base | User confirmed no upstream changes needed |
| Main branch is default | Merged from custom branch, clean state |

## Notes

- User requirements: see findings.md for detailed feature requirements
- Use subagents (max 3 parallel) for code research
- Design phase only - do NOT implement until user explicitly says so
