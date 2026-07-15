# Handoff Document

Generated: 2026-07-14
Next session focus: FileBrowser v1.4.0.6-fde — properly fix Bug B (status lost after dialog close) with thorough investigation, AND add queue list API per MASTER directive

---

## ⭐ CRITICAL: First thing to read

**v1.4.0.6 Next Version Plan (READ FIRST — contains EVERYTHING):**
/home/hermes/workspace/filebrowser-fde/docs/v1.4.0.6-next-version-plan.md
- 300+ lines, 13KB
- Ch2.2 Bug B has the failed fix code + 7 hypotheses + 5 alternative approaches + step-by-step investigation protocol
- Ch5 fix priority table (3 FIXED, 1 REVERTED, 3 PENDING)
- ALL MASTER quotes recorded verbatim

**Recent handoff (v1.4.0.5-fde implementation):**
/tmp/handoff-FQzls0.md
- Background on previous version
- Useful for understanding what came before

---

## Context

FileBrowser Quantum v1.4.0.5-fde implemented 8 commits. Real-world testing
by MASTER revealed 4 bugs:

- **Bug A** (folder compress does nothing): ✅ FIXED (commit 40aab284)
- **Bug B** (status lost after dialog close): ❌ FAILED — fix written based
  on assumptions, MASTER deployed and saw NO status bar appearing. Reverted
  in commit 079cbeda.
- **Bug C** (progress stuck at N-1/N): ✅ FIXED (commit d1b66eea)
- **Bug F** (UI changes invisible): ✅ FIXED (commit 56d5b6b2)
  - Checkbox moved to below file list
  - Old CSS annotated as LEGACY
  - Dead .compress-file-thumb removed

Transition engine was reverted to v1.4.0.4 (commit 20558745) because
v1.4.0.5 fix caused visual quality flicker for logged-in users.

**Current state:** v1.4.0.5-fde-hotfix on tag `v1.4.0.5-fde-hotfix-latest`
(commit 2c6b4490). Force-pushed to main. Docker tar built and saved.

---

## Key Decisions

- **Bug A fix kept:** selectedFileList now uses expandedItems when available.
  Folder compression works for single + multiple folders.
- **Bug B fix reverted:** The mounted() + recoverQueueStatus() approach was
  a complete failure. Next AI must investigate WHY before writing any fix.
  Likely causes (DO NOT assume): Prompts.vue keep-alive wrapping, silent
  catch{} block hiding errors, mounted not firing on reopen, race condition
  with idle status.
- **Bug C fix kept:** setWorkerActive(false) only resets to "idle" if status
  is still "running". "completed" status preserved for frontend poll.
- **Bug F fix kept:** Checkbox HTML moved. Old CSS annotated but NOT deleted
  per MASTER directive (2026-07-14 — don't delete uncertain legacy code).
- **No transition changes:** Transition engine stays at v1.4.0.4-hotfix.
- **MASTER's most important feedback:** "原因就喺根本就無。我入去压缩图片
  窗口根本从来无发现有呢个状态栏" — meaning the entire Bug B fix produced
  ZERO user-visible result. Next AI must verify in real test before
  claiming success.

---

## Remaining Tasks (v1.4.0.6-fde priorities)

- [ ] **P0 REOPEN Bug B:** Properly investigate Prompts.vue component
  lifecycle. Read it FIRST. Then test mounted() actually fires with
  console.log. Then design fix using one of 5 alternative approaches
  (recommended: persistent banner outside dialog).
- [ ] **P0 Ch4.1 Queue List API:** Add GET /compress-images/queue endpoint
  that returns ALL items in queue (not just current). MASTER directive:
  "要有一个功能可以取得所有加入到队列的压缩列表咁样找回状态"
- [ ] P1 Bug E: Verify success notification works now (should after Bug C fix)
- [ ] P1 Bug D: Queue count confusion (UX issue)
- [ ] P1 CSS duplicate cleanup (when safe per MASTER directive)

---

## Blockers

- **Bug B requires real-world test to verify:** Single-core 2GB target server.
  Cannot test locally. MASTER must deploy + test for each candidate fix.
- **Prompts.vue not yet investigated:** Previous AI did NOT read it before
  writing the failed fix. Next AI MUST read it first.
- **No automated tests:** Project has no test framework. All verification
  via go build + vet + grep + MASTER's manual deployment test.

---

## Artifacts

- **Design doc (CRITICAL — read first):** docs/v1.4.0.6-next-version-plan.md
  — Complete plan for v1.4.0.6 with Bug B failure analysis
- **Previous version plan:** docs/v1.4.0.5-next-version-plan.md
  — Historical context (1443 lines, all MASTER feedback for v1.4.0.5)
- **Implementation plan (DONE):** docs/superpowers/plans/2026-07-14-v1.4.0.5-fde-plan.md
  — Reference for what was implemented (don't redo)
- **Research findings:** docs/v1.4.0.6-next-version-plan.md Ch2-3 (combined)
  — Bug descriptions + fix attempts
- **PWF task_plan:** task_plan_hotfix.md
  — Lightweight task tracker for this hotfix session
- **PWF original task_plan:** task_plan.md — Old plan for v1.4.0.5-fde
- **LanceDB memory ID 1a6f69b5:** Bug B failure lesson (search "FileBrowser Bug B mounted")
- **Docker tar:** filebrowser-fde-v1.4.0.5-fde-hotfix.tar (84MB) — current build
- **Branch:** v1.4.0.5-fde (tag: v1.4.0.5-fde-hotfix-latest)
- **Main branch:** force-pushed to commit 2c6b4490

---

## Git Commit History (most recent first)

```
2c6b4490 docs: v1.4.0.6 Ch2.2 + Ch5 - document Bug B fix failure + revert + lessons learned
079cbeda Revert "fix: status lost after dialog close - recover queue on mount"
358ebe0a docs: v1.4.0.6 Ch2.2 + Ch5 - record MASTER feedback on Bug B
56d5b6b2 fix: UI changes now visible - move checkbox + annotate legacy CSS   [Bug F]
cbb83649 fix: status lost after dialog close - recover queue on mount        [Bug B — REVERTED]
d1b66eea fix: progress stuck at N-1/N - preserve completed status           [Bug C]
40aab284 fix: folder compress does nothing - selectedFileList uses expandedItems  [Bug A]
adade51d docs: v1.4.0.6-fde Ch4-5 - MASTER directives
20558745 revert: transition engine back to v1.4.0.4-hotfix behavior
753bef4b docs: v1.4.0.6-fde next version plan
... (older commits before v1.4.0.5-fde)
```

---

## Environment Notes

- **Project path:** /home/hermes/workspace/filebrowser-fde
- **Go binary:** /usr/local/go/bin/go (NOT in default PATH, always use full path)
- **Target server:** 1 vCPU, 2GB RAM, 2GB ZRAM + 2GB SWAP
- **Docker volume on target:** /root/qbb/downloads:/folder
- **Backend route:** GET /compress-images/status (admin permission required)
- **Frontend API:** pollStatus() in frontend/src/api/compress.js (uses fetchURL)
- **build cache cleanup:** Run `docker builder prune -f` after each Docker build
- **i18n:** 3 files must sync: en.json, zh-cn.json, zh-tw.json
- **No local testing possible** — only Docker build + remote deployment
- **MASTER uses Matrix (Element) client** — output must be pure text, NO code blocks, NO Markdown headers

---

## Recommended Skills for Next AI

- **superpowers/brainstorming:** Required before designing Bug B fix — explore
  alternatives A-E in Ch2.2, understand trade-offs
- **superpowers/systematic-debugging:** Required for Bug B investigation —
  read Prompts.vue, test mounted() lifecycle, find real root cause
- **superpowers/writing-plans:** When designing v1.4.0.6 implementation plan
- **research/real-time-news-reporter:** NOT relevant (this is code work)
- **computer-use / browser-automation:** MAY be needed to test dialog
  behavior in browser

---

## 🚨 CRITICAL WARNINGS for Next AI

1. **DO NOT assume Vue component lifecycle.** READ Prompts.vue FIRST.
   Check for `<keep-alive>`, v-if rendering, dynamic components.
2. **DO NOT use silent catch{} blocks.** Add console.error() at minimum
   so debugging is possible.
3. **DO NOT claim a fix works without real user test.** This is what
   happened with Bug B — code looked correct, build passed, deployed,
   user saw nothing.
4. **DO NOT delete legacy CSS classes.** Annotate them per MASTER's
   directive (2026-07-14). Only delete after confirmed safe.
5. **DO NOT change transition engine again.** v1.4.0.4-hotfix transition
   behavior is the baseline. Don't touch unless fundamentally new approach.

---

## Files to Read First (in order)

1. **docs/v1.4.0.6-next-version-plan.md** (Ch2.2 especially — 200 lines)
2. **frontend/src/components/prompts/Prompts.vue** (full file — not read before)
3. **frontend/src/components/prompts/CompressImages.vue** (current state)
4. **backend/http/compress.go** (compressStatusHandler + queueMgr)
5. **backend/http/httpRouter.go** (route registration)
6. **frontend/src/api/compress.js** (pollStatus function)

---

## Quick Verification Commands

```bash
cd /home/hermes/workspace/filebrowser-fde
git log --oneline -8                       # See recent commits
git tag --list | grep hotfix               # Check tags
git branch --show-current                  # Current branch (should be v1.4.0.5-fde)
cd backend && /usr/local/go/bin/go build -o /dev/null ./...  # Verify build
```