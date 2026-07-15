# Handoff Document - FileBrowser Quantum v1.4.0.6-fde

Generated: 2026-07-14
Next session focus: Design and implement v1.4.0.6-fde (mobile fixes + compress control features + Bug B re-investigation)

---

## Context

FileBrowser Quantum is a custom fork of FileBrowser Quantum (gtsteffaniak/filebrowser)
maintained by MASTER at GitHub: github.com/fde-lander/filebrowser.git

The project adds custom features on top of the original FileBrowser:
- Image viewer with dual-buffer transition engine (crossfade/fade/instant)
- Tap navigation (left/right 40% zones, 200ms double-tap guard)
- Mobile nav button persistence + opacity slider
- Image compression (WebP encoder, 3 tiers, ZSTD backup, polling progress)
- Extract-to-new-folder feature

Current deployed version: v1.4.0.5-fde-hotfix (tag: v1.4.0.5-fde-hotfix-latest)
This version fixed 3 bugs (A, C, F) but Bug B fix FAILED completely and was reverted.
MASTER tested on mobile and found 3 NEW bugs (G, H, I) not visible on PC.
MASTER also requested 3 new compress control features (pause/continue/cancel).

The v1.4.0.6-next-version-plan.md has been updated to 772 lines with 9 chapters
covering all issues, requirements, and the DOUBT investigation protocol.

---

## Key Decisions

- **Transition engine is LOCKED at v1.4.0.4-hotfix:** Do NOT modify transition engine
  again. v1.4.0.5 attempt caused visual quality regression for logged-in users.
  Reverted in commit 20558745. Only revisit with fundamentally different approach.

- **Bug B fix FAILED and was REVERTED:** Previous AI wrote mounted() + recoverQueueStatus()
  in CompressImages.vue WITHOUT reading Prompts.vue. MASTER deployed and saw ZERO
  user-visible result. Reverted in commit 079cbeda. v1.4.0.6 MUST re-investigate
  from scratch using DOUBT protocol (Ch9).

- **DOUBT Principle (MASTER Directive):** All v1.4.0.6 work MUST follow:
  1. Dispatch subagents to research each issue (read source code, find line numbers)
  2. Main agent verifies subagent findings by reading actual code independently
  3. Design fix only with verified facts (cite specific code lines as evidence)
  4. MASTER deploys and tests before any fix is declared done
  5. NEVER assume Vue lifecycle, NEVER use silent catch{}, NEVER guess

- **Mobile fixes must NOT break PC:** Bug G/H/I are mobile-specific. All fixes must
  preserve existing PC behavior. Swipe-to-navigate on mobile is GOOD - do not break it.

- **Compress control features are NEW:** Pause/Continue/Cancel require backend
  queueMgr changes + new API endpoints + frontend UI buttons. Design carefully.

- **Legacy CSS: annotate, do NOT delete:** MASTER directive (2026-07-14). Old CSS
  classes that might affect functionality should be annotated with LEGACY comment,
  not removed. Only delete after confirmed safe.

- **Docker distribution: tar files only:** Custom Docker images are NEVER pushed to
  Docker Hub. Build locally -> docker save -> scp tar to target server -> docker load.
  Tar files MUST NOT be deleted (deployment artifacts).

---

## Remaining Tasks (v1.4.0.6-fde)

### P0 - Critical (3 items)
- [ ] **Bug B: Compression progress lost after dialog close (REOPENED)**
  Must read Prompts.vue FIRST. Verify mounted() fires on reopen. Use DOUBT protocol.
  See Ch2.2 + Ch8.1. 7 hypotheses documented, 5 alternative approaches listed.

- [ ] **Bug G: Mobile tap navigation completely fails (NEW)**
  Tap zones + nav buttons do nothing on mobile. Likely click vs touch event issue.
  Must NOT break PC navigation or mobile swipe. See Ch7.1.

- [ ] **Ch4.1: Queue List API**
  GET /compress-images/queue - returns ALL queue items with status.
  MASTER directive: "要有一个功能可以取得所有加入到队列的压缩列表"

### P1 - High Priority (7 items)
- [ ] **Bug H: Mobile nav button has no transition (NEW)** - Black screen flash.
  See Ch7.2. Must ensure transition function is called on mobile nav path.
- [ ] **Bug I: Mobile compress preview layout broken (NEW)** - Images overlap
  tips text on mobile portrait. See Ch7.3. Need mobile CSS media queries.
- [ ] **Bug E: No success notification** - Should work after Bug C fix. Verify.
- [ ] **Bug D: Queue count confusion** - Depends on Bug B refinement.
- [ ] **Ch8.2: Pause compression** - POST /compress-images/pause. Current file
  finishes, then worker pauses. Queue preserved.
- [ ] **Ch8.3: Continue compression** - POST /compress-images/resume. Worker
  resumes from next file in queue.
- [ ] **Ch8.4: Cancel compression** - POST /compress-images/cancel. Current file
  finishes, then entire queue discarded.

### P3 - Low Priority
- [ ] CSS duplicate cleanup (annotate only per MASTER directive)
- [ ] Transition fix for anonymous users (needs fundamentally new approach)

---

## Blockers

- **No local testing possible:** Target server is 1 vCPU 2GB RAM (remote).
  This machine can only Docker build -> save -> tar -> scp to target -> load.
  MASTER deploys and tests manually. Every fix requires a full deploy cycle.

- **Prompts.vue NOT YET INVESTIGATED:** Previous AI did NOT read Prompts.vue
  before writing Bug B fix. This is the #1 reason the fix failed. Next AI
  MUST read this file FIRST. Check: keep-alive? v-if? dynamic component?
  Does mounted() fire on dialog reopen? Use console.log to verify.

- **No automated test framework:** Project has zero tests. All verification
  via: go build + go vet (backend), grep markers (frontend), MASTER manual
  deployment test (final gate).

- **Single-core target server:** 1 vCPU + 2GB RAM + 2GB ZRAM + 2GB SWAP.
  Backend goroutine + polling must be lightweight. Avoid CPU-heavy patterns.
  SSE was already replaced with polling due to permanent failure on single-core.

- **i18n changes must sync 3 files:** en.json, zh-cn.json, zh-tw.json.
  Every new UI string must be added to all 3 simultaneously.

---

## Artifacts

### Primary Documents (READ FIRST, in order)

- **Plan doc (CRITICAL):** /home/hermes/workspace/filebrowser-fde/docs/v1.4.0.6-next-version-plan.md
  772 lines, 9 chapters. Contains ALL bugs, requirements, MASTER quotes,
  DOUBT protocol, subagent investigation plan, file reading checklist.
  Ch2.2 = Bug B failure analysis (7 hypotheses + 5 alternative approaches).
  Ch7 = Mobile bugs G/H/I (NEW, from MASTER mobile testing).
  Ch8 = Compress control features (pause/continue/cancel).
  Ch9 = DOUBT principle + subagent investigation protocol.

- **Previous handoff:** /home/hermes/workspace/filebrowser-fde/docs/handoff-v1.4.0.5-fde-hotfix-to-v1.4.0.6.md
  194 lines. Historical context from v1.4.0.5 implementation session.

- **Design spec (v1.4.0.5):** /home/hermes/.hermes/docs/superpowers/specs/2026-07-14-v1.4.0.5-fde-design.md
  730 lines. Root cause analysis + impact assessment + rollback plan for v1.4.0.5.

- **Implementation plan (v1.4.0.5, COMPLETED):** /home/hermes/.hermes/docs/superpowers/plans/2026-07-14-v1.4.0.5-fde-plan.md
  1139 lines. Reference for what was implemented. Do NOT redo these tasks.

### Source Files (must read before implementing each fix)

- frontend/src/components/prompts/Prompts.vue (991 lines)
  How CompressImages is mounted. keep-alive? v-if? lifecycle on close/reopen.
  **Previous AI did NOT read this - that is why Bug B fix failed.**

- frontend/src/components/prompts/CompressImages.vue (922 lines)
  Compress dialog: file list, preview overlay, mounted(), polling, progress UI.
  Bug B fix was here (reverted). Bug I preview layout is here.

- frontend/src/components/files/ExtendedImage.vue (1160 lines)
  Image viewer: dual-buffer imgA/imgB, transition engine, tap zones, nav buttons,
  touch/click event handlers. Bug G + Bug H are here.

- frontend/src/views/settings/Profile.vue (525 lines)
  User settings: transition mode, nav button persistence, opacity, backup toggle.

- backend/http/compress.go (778 lines)
  queueMgr struct, compressWorker goroutine, status handler, all compress logic.
  Pause/Continue/Cancel features need changes here.

- backend/http/httpRouter.go (416 lines)
  Route registration. New API endpoints (pause/resume/cancel/queue) go here.

- frontend/src/api/compress.js (98 lines)
  pollStatus() + all compress API call functions. New pause/resume/cancel calls.

- i18n files (all 911 lines each):
  frontend/src/i18n/en.json
  frontend/src/i18n/zh-cn.json
  frontend/src/i18n/zh-tw.json

### Git State

- **Project path:** /home/hermes/workspace/filebrowser-fde
- **Branch:** v1.4.0.5-fde
- **Current tag:** v1.4.0.5-fde-hotfix-latest (commit 2c6b4490)
- **Remote:** github.com/fde-lander/filebrowser.git
- **Default branch:** main (force-pushed to 2c6b4490)

Git commit history (most recent first):
  2c6b4490 docs: v1.4.0.6 Ch2.2 + Ch5 - document Bug B fix failure + revert + lessons
  079cbeda Revert "fix: status lost after dialog close - recover queue on mount"
  358ebe0a docs: v1.4.0.6 Ch2.2 + Ch5 - record MASTER feedback on Bug B
  56d5b6b2 fix: UI changes now visible - move checkbox + annotate legacy CSS   [Bug F]
  cbb83649 fix: status lost after dialog close - recover queue on mount        [Bug B - REVERTED]
  d1b66eea fix: progress stuck at N-1/N - preserve completed status           [Bug C]
  40aab284 fix: folder compress does nothing - selectedFileList uses expanded  [Bug A]
  adade51d docs: v1.4.0.6-fde Ch4-5 - MASTER directives
  20558745 revert: transition engine back to v1.4.0.4-hotfix behavior
  753bef4b docs: v1.4.0.6-fde next version plan
  ... (older v1.4.0.5-fde commits below)

### Docker Images

- filebrowser-fde:v1.4.0.5-fde-hotfix (current deployed)
- filebrowser-fde:v1.4.0.5 (pre-hotfix, has Bug A/C/F unfixed)
- filebrowser-fde:v1.4.0.4 (previous stable)
- **Tar:** filebrowser-fde-v1.4.0.5-fde-hotfix.tar (84MB, current deploy artifact)

### LanceDB Memory References

- ID 1a6f69b5: Bug B failure lesson (search "FileBrowser Bug B mounted")
- ID f13604b0: v1.4.0.3 test results (search "FileBrowser Quantum v1.4.0.3 test")
- ID 499be97c: v1.4.0.4 MASTER feedback (search "FileBrowser v1.4.0.4 test feedback")
- ID 49271d6b: v1.4.0.2 complete feature list (search "FileBrowser Quantum v1.4.0.2 custom fork")

---

## Environment Notes

- **Go binary:** /usr/local/go/bin/go (NOT in default PATH - always use full path
  or export PATH=/usr/local/go/bin:$PATH before go commands)

- **Target deployment server:** 1 vCPU, 2GB RAM, 2GB ZRAM + 2GB SWAP
  Docker volume: /root/qbb/downloads:/folder
  Server is remote - cannot test locally

- **Docker build workflow:**
  1. Build: docker build -t filebrowser-fde:v1.4.0.6 .
  2. Save: docker save filebrowser-fde:v1.4.0.6 -o filebrowser-fde-v1.4.0.6.tar
  3. Cleanup: docker builder prune -f (MANDATORY after each build)
  4. SCP tar to target server
  5. Load: docker load -i filebrowser-fde-v1.4.0.6.tar
  6. Update docker-compose.yaml image tag
  7. Restart: docker-compose up -d

- **MASTER uses Matrix (Element) client:** Output MUST be:
  - NO code blocks (```) - Element renders them as garbage
  - NO Markdown headers (#, ##) - Element does not render them
  - NO tables - Element does not render them
  - USE bold, emoji, separator lines for structure
  - USE linear lists with bullet points
  - Pure text commands, not code blocks

- **Code editing rules:**
  - Use patch tool for all edits (NOT sed/awk)
  - Each patch < 2.5KB (avoid LLM context overload)
  - go build + go vet after backend changes
  - grep verification after frontend changes
  - i18n: 3 files must sync (en.json, zh-cn.json, zh-tw.json)
  - Do NOT delete legacy CSS - annotate with LEGACY comment
  - Do NOT modify transition engine (locked at v1.4.0.4-hotfix)
  - Do NOT use silent catch{} - add console.error

- **HARD-GATE:** No sudo or privilege escalation without MASTER approval

---

## Recommended Skills

- **superpowers/brainstorming:** Required before designing any fix. Explore
  alternatives, understand trade-offs. Especially for Bug B (5 approaches listed
  in Ch2.2) and mobile fix strategy.

- **superpowers/systematic-debugging:** Required for Bug B/G/H investigation.
  Read code -> form hypothesis -> verify with code read -> design fix.

- **superpowers/writing-plans:** After research phase, write detailed
  implementation plan with exact line numbers + code + verification commands.

- **superpowers/executing-plans:** When executing the implementation plan
  with review checkpoints.

- **superpowers/subagent-driven-development:** For parallel research tasks.
  Dispatch subagents to read different source files simultaneously, then
  verify their findings independently (DOUBT protocol Ch9).

- **filebrowser-quantum-docker-setup:** Project-specific skill with codebase
  architecture reference, deployment instructions, and post-subagent bug patterns.

---

## CRITICAL WARNINGS for Next AI

1. DO NOT assume Vue component lifecycle. READ Prompts.vue FIRST.
   Check for keep-alive, v-if rendering, dynamic components.
   Previous AI skipped this -> Bug B fix produced ZERO user-visible result.

2. DO NOT use silent catch{} blocks. Add console.error() at minimum.
   Silent catches hide ALL errors. MASTER will never know what failed.

3. DO NOT claim a fix works without real deployment test by MASTER.
   Code compiling != working. Build passing != user-visible effect.
   MASTER must deploy + test + confirm before declaring done.

4. DO NOT delete legacy CSS classes. Annotate them with:
   /* LEGACY from v1.4.0.4 - per MASTER directive: keep + annotate only */

5. DO NOT modify transition engine. v1.4.0.4-hotfix is the locked baseline.
   Any transition change risks visual regression for logged-in users.

6. DO NOT break mobile swipe-to-navigate. MASTER confirmed swipe transition
   is BETTER than original. Mobile tap fix must NOT affect swipe behavior.

7. DO NOT trust subagent research blindly. Main agent MUST independently
   read the actual code at claimed line numbers to verify every finding.

8. DO NOT skip reading any file that your fix touches. Read it fully first.

