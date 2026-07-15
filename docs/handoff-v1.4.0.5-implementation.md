# Handoff Document: FileBrowser Quantum v1.4.0.5-fde Implementation

Generated: 2026-07-14
Next session focus: Execute the v1.4.0.5-fde implementation plan (9 fix modules, each as independent commit)

---

## ⭐ CRITICAL: First thing to read

**Implementation Plan (STEP-BY-STEP with code):**
/home/hermes/.hermes/docs/superpowers/plans/2026-07-14-v1.4.0.5-fde-plan.md
- 1139 lines, 36KB, 10 Tasks with exact code, commands, line numbers
- Contains CRITICAL CONTEXT section for project setup, Go path, Docker commands
- EVERY step has: file path, line number, before/after code, verify command, commit message

**Design Spec (full rationale):**
/home/hermes/.hermes/docs/superpowers/specs/2026-07-14-v1.4.0.5-fde-design.md
- 730 lines, 28KB, 10 chapters with root cause analysis + impact assessment

**Research Findings (verified ground truth):**
/home/hermes/workspace/filebrowser-fde/findings.md
- All Wave 1+2+3+4 research results with exact line numbers

---

## Context

FileBrowser Quantum custom fork v1.4.0.5-fde - the 5th and final stability upgrade.
9 subagents across 4 waves performed exhaustive code research. All findings verified
by main agent reading actual code. 9 topics discussed with MASTER, all confirmed.
Design spec and implementation plan written and self-reviewed.

**Current state:** All planning complete. AWAITING MASTER's explicit "implement" instruction
before writing any code. MASTER said: "写Spec文档之前停下等我指示" and
"全过程不能实施除非我主动跟你说实施".

---

## Key Decisions

- **Queue system (not single instance):** MASTER confirmed direction B - compress queue
  where user can always add files. Bottom thin bar shows queue status.
- **Transition default 'fade':** Changed from 'crossfade' to 'fade' (500ms, softer).
  MASTER requested this as default for all users including anonymous/share.
- **PNG low tier = WebP Q75:** Remove OxiPNG branch entirely. Was hanging forever.
- **Share link fix approach:** Fix transition engine fault tolerance, NOT change share URL
  construction. 6 fixes: waitForDecode error handling, onImageError guard,
  finishTransition verification, clear old buffer src, default fade, imgB @load.
- **Each module = independent commit:** Can rollback individually without affecting others.
- **Transition fix rollback:** `git checkout v1.4.0.4-hotfix -- <files>` restores exact
  v1.4.0.4 behavior. MASTER requires this specific rollback capability.
- **Thumbnail dead code removal = separate commit:** MASTER specifically requested this
  be independently rollbackable.
- **No WebP encoding changes:** Quality is priority. MASTER explicitly forbade
  encoding_speed/library/quality changes.
- **Build cache cleanup:** `docker builder prune -f` after each Docker build.
- **MASTER's read-only user workaround:** MASTER found that creating a read-only user
  with preset image viewer settings (preload, fade transition) gives anonymous share
  users a near-identical experience. This is independent of code fixes and serves as
  a fallback if transition fix is rolled back.

---

## Remaining Tasks

Execute the implementation plan (Task 0 through Task 9):

- [ ] Task 0: Pre-implementation setup (tag + backup files + baseline build)
- [ ] Task 1: Transition Fix - ExtendedImage.vue + Profile.vue (6 steps)
- [ ] Task 2: SSE -> Queue Polling - compress.go + httpRouter.go + compress.js + CompressImages.vue (10 steps)
- [ ] Task 3: PNG Low Tier -> WebP - compress.go (2 steps)
- [ ] Task 4: Preview Original + Layout - CompressImages.vue (3 steps)
- [ ] Task 5: Backup Persistence - state.js + mutations.js + Profile.vue + CompressImages.vue (5 steps)
- [ ] Task 6: Folder Expansion + Checkbox - CompressImages.vue + i18n (5 steps)
- [ ] Task 7a: SavedBytes - compress.go (3 steps)
- [ ] Task 7b: Thumbnail Dead Code - CompressImages.vue (2 steps)
- [ ] Task 8: Mobile CSS - CompressImages.vue (2 steps)
- [ ] Task 9: Final Build + Verify + Docker (9 steps)

---

## Blockers

- **NO implementation until MASTER says "implement"** - this is a HARD-GATE
- No local testing possible - deployment to remote server via Docker image only
- Target server: 1 vCPU, 2GB RAM, 2GB ZRAM + 2GB SWAP (single-core constraints)
- Go binary at /usr/local/go/bin/go (NOT in default PATH)
- CompressImages.vue is touched by 7 commits - MUST apply in order
- Never delete Docker tar files (deployment artifacts)

---

## Artifacts

- Design Spec: /home/hermes/.hermes/docs/superpowers/specs/2026-07-14-v1.4.0.5-fde-design.md - 10 chapters, 730 lines, root cause + fix plan + rollback
- Implementation Plan: /home/hermes/.hermes/docs/superpowers/plans/2026-07-14-v1.4.0.5-fde-plan.md - 1139 lines, step-by-step with code, MUST READ FIRST
- Research Findings: /home/hermes/workspace/filebrowser-fde/findings.md - All 4 waves of subagent research, verified
- PWF task_plan: /home/hermes/workspace/filebrowser-fde/task_plan.md - Phase 5 = Implementation
- PWF progress: /home/hermes/workspace/filebrowser-fde/progress.md - Session log
- v1.4.0.5 plan doc: /home/hermes/workspace/filebrowser-fde/docs/v1.4.0.5-next-version-plan.md - 1443 lines, all MASTER feedback + design decisions
- Project path: /home/hermes/workspace/filebrowser-fde
- Git branch: v1.4.0.5-fde (from v1.4.0.4-hotfix @ 73940a56)
- GitHub: https://github.com/fde-lander/filebrowser.git
- Handoff (previous): /tmp/handoff-5EV2Ba.md - v1.4.0.4 to v1.4.0.5 transition

---

## Environment Notes

- Go binary: /usr/local/go/bin/go (use full path in all commands)
- Go module root: /home/hermes/workspace/filebrowser-fde/backend/
- Docker build: `sudo docker build --build-arg="VERSION=v1.4.0.5" ...`
- Docker save: `docker save filebrowser-fde:v1.4.0.5 -o filebrowser-fde-v1.4.0.5.tar`
- Build cache cleanup: `docker builder prune -f` (after EACH build)
- Docker volume on target: /root/qbb/downloads:/folder
- i18n: 3 files MUST sync: en.json, zh-cn.json, zh-tw.json
- Use `patch` tool for code changes, never sed/awk
- MASTER uses Matrix (Element) - output must be pure text, NO code blocks, NO markdown headers

---

## Recommended Skills

- superpowers/executing-plans: For executing the implementation plan task-by-task
- superpowers/subagent-driven-development: For dispatching implementation tasks to subagents
- superpowers/systematic-debugging: If build errors occur during implementation
- hermes-agent: For any Hermes-specific configuration questions
