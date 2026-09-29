# Progress Log (v1.4.0.6-fde)

## Session: 2026-07-15

### Session Start
- **Branch:** v1.4.0.6-fde (created from v1.4.0.5-fde @ 2c6b4490)
- **Base tag:** v1.4.0.5-fde-hotfix-latest
- **PWF + TODO created**
- **Brainstorming skill loaded**
- **Codebase-memory indexed:** 3927 nodes, 13640 edges
- **Handoff document read:** docs/handoff-v1.4.0.6-fde.md (303 lines)
- **Plan document read:** docs/v1.4.0.6-next-version-plan.md (772 lines, ALL)
- **LanceDB memory recalled:** 5 relevant entries (Bug B failure, v1.4.0.3/4/5 results)
- **filebrowser-quantum-docker-setup skill loaded**

### MASTER Directives
1. v1.4.0.6-fde is the target version
2. Use codebase-memory system for project analysis
3. Subagent research for complex code investigation
4. Give subagents sufficient skills + deterministic instructions
5. Consider truncation prevention for subagent results
6. Main agent verifies via codebase + code reading
7. NO spec writing without MASTER approval
8. NO implementation without MASTER approval
9. Discuss thoroughly before acting
10. Manage TODO properly

### Phase 1: DOUBT Code Research (COMPLETED)
- Wave 1: 3 subagents dispatched + returned (43s total)
- All findings verified by main agent via independent code reading
- Key verified findings:
  - Prompts.vue: NO keep-alive, mounted() fires on every reopen
  - ExtendedImage.vue: @touchmove.prevent kills click events on mobile
  - compress.go: queueMgr struct confirmed, NO pause/cancel mechanism exists
  - showPrompt system has built-in confirm/callback mechanism
  - compressBackup persistence: 4-touchpoint (state.js + mutations.js + users.go + CompressImages.vue)

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

### Phase 3: Design Spec Writing (COMPLETED)
- 10 chapters written, 1230 lines, 44KB
- Spec location: ~/.hermes/docs/superpowers/specs/2026-07-15-v1.4.0.6-fde-design.md
- Ch4 Bug H updated: root cause = double navigation + transitioning overlay (verified by 2 subagents)
- Self-review: 0 placeholders, 0 contradictions, 0 ambiguities
- Bug H fix: 2 independent fixes in nextPrevious.vue (double nav guard + skip transitioning for images)
- nextPrevious.vue has getters.previewType() but NO previewType computed -> use getters.previewType() directly
- Status: completed, awaiting MASTER review

### Phase 4: Implementation Plan Writing (COMPLETED)
- writing-plans skill loaded
- Plan location: ~/.hermes/docs/superpowers/plans/2026-07-15-v1.4.0.6-fde-plan.md
- 10 tasks + post-build checklist, 2537 lines, 84KB
- ALL tasks have precise old_string/new_string, complete code, grep verification, commit messages
- Self-review: spec coverage complete, 0 placeholders, type consistency verified
- Status: completed, awaiting MASTER review

### Phase 5: Handoff Document (COMPLETED)
- Handoff skill loaded
- Handoff doc: docs/handoff-v1.4.0.6-fde-implementation.md (8KB)
- Covers: project context, 10 task overview, key decisions, CRITICAL warnings,
  post-build checklist, artifacts, environment, recommended skills
- Status: completed, ready for implementation AI

### Phase 5: Implementation (IN PROGRESS)
- Session start: 2026-07-16
- Execution mode: INLINE (主 Agent 串行)
- Skills loaded: PWF, executing-plans, TDD
- Plan: 10 Tasks, each = 1 git commit
- Status: 10/10 tasks completed + Post-Build verification

#### All Tasks COMPLETED
- Task 1: Bug G (87b16212) - ExtendedImage.vue touchmove.prevent
- Task 2: Bug H (50dda77c) - nextPrevious.vue double nav + transitioning
- Task 3: Bug I (b129eda3) - CompressImages.vue @media CSS
- Task 4: Backend compress control (e87a7eb4) - compress.go + users.go
- Task 5: Routes (8421e7df) - httpRouter.go 5 new routes
- Task 6: Frontend API (eac03203) - compress.js 5 new functions
- Task 7: ConfirmAction.vue (65ce705f) - generic confirm dialog
- Task 8: CompressStatusBar.vue (1e558660) - global status bar
- Task 9: CompressImages.vue (ba2f4801) - control buttons + mounted check
- Task 10: Settings + i18n (6b14db97) - 7 files, 21 keys x 3 languages

#### Post-Build Verification
- go build: PASS ✅
- go vet: PASS ✅
- go mod verify: PASS ✅
- Grep sweep: ALL PASS ✅
- JSON validity: EN/CN/TW ALL PASS ✅
- Commit count: 10 task commits ✅
- Docker build: SUCCESS ✅ (84MB tar)
- Docker save: SUCCESS ✅
- Build cache pruned: 3.691GB reclaimed ✅
- Tar: filebrowser-fde-v1.4.0.6.tar (84MB)
- Status: READY FOR DEPLOYMENT

---
## Session: 2026-09-29 — README Rewrite for v1.4.0.5-fde-hotfix

### MASTER Directives
1. Only v1.4.0.5-fde-hotfix works — MASTER uses only this version
2. LATEST RELEASE already published; problem = README not updated yet
3. v1.4.0.6 and everything after = VOIDED; future upgrades skip them entirely
4. README must be fully Traditional Chinese
5. Signature: FDE only, no other names
6. Voided versions use recommended wording (not deleted silently)
7. Release handling: OPTION A — do not touch published assets, document the situation
8. README front section needs dedicated FDE FORK block + dated version history (descending date order)

### Environment Correction
- Previous PWF goal (v1.4.0.6) was stale. Goal updated to reflect v1.4.0.5-fde-hotfix baseline.
- Project registry status was `completed` while work was actually unfinished — noted, not changed.

### Release Verification (read-only, no modifications)
- Release tag: v1.4.0.5-fde-stable, published 2026-07-15, points to commit 2c6b4490
- Asset filename `filebrowser-fde-v1.4.0.6.tar` and inner image tag `filebrowser-fde:v1.4.0.6` are stale naming from the voided version
- Decisive layer-content scan proved asset = v1.4.0.5-fde-hotfix + Bug I CSS:
  - Bug I marker `flex-direction:column!important` present 2x (older builds: 0)
  - Voided-version markers absent: no `compress-images/pause`, no `compressStatusBar`, no `compressPauseTimeout`
  - SHA256 matches local build: b00e914105dc2c6bea95c7dc56be370d236f45cfad7cfaa218a97b2039908399
  - Image built-in Version/CommitSHA are EMPTY (build did not pass -ldflags values)
- Correction to earlier assumption: tag `v1.4.0.5-fde-stable` itself does NOT contain Bug I CSS (Bug I landed in HEAD only, added to the rebuilt image).

### Deliverable
- README.md rewritten (6976 -> 11520 bytes), full Traditional Chinese
- Signature: FDE only
- New structure: FDE FORK intro -> dated version history (descending) -> deploy steps -> release asset naming note -> voided versions -> feature list -> tech stack -> environment -> credits
- Backup: README.md.bak-20260929-pre-tc-rewrite.bak (byte count verified identical to original)

### Verification
- Traditional-Chinese check via built-in conversion table (1493 chars): no simplified characters
- Stale reference sweep: 0 hits for `filebrowser-fde-v1.4.0.4.tar`, `filebrowser-fde:v1.4.0.4`, `正在開發中`
- Signature sweep: no AI / mascot names present
- v1.4.0.5 SSE removal confirmed against source: progressManager functions deleted, subscribeProgress replaced by pollStatus

### Open Item
- ~~README not yet committed to git~~ RESOLVED in the next session below — committed as a2302c66 and pushed to origin/main.

---

## Session: 2026-09-29 (cont.) — DEPLOY.md, Release Publication, MAIN Lock

### MASTER Directives
1. RELEASE asset: use CHOICE A — filebrowser-fde-v1.4.0.5-fde-hotfix.tar
2. MAIN is LOCKED to tag V1.4.0.5-FDE-HOTFIX
3. Old release (if not v1.4.0.5-fde-hotfix) must be voided; publish ours as LATEST
4. If no compiled TAR exists, skip the release entirely
5. Future deployment model changed: target server pulls repo, builds image, runs docker compose
6. A complete, detailed DEPLOY document is MANDATORY for this new model

### Decisive Tar Verification (full extraction, per-file scan)
- `filebrowser-fde-v1.4.0.5-fde-hotfix.tar`: RepoTag v1.4.0.5-fde-hotfix, built-in Version v1.4.0.5-fde-hotfix,
  JS bundle contains `recoverQueueStatus` x3 (the Bug B attempt), CSS has NO Bug I rules. Built 07-14 08:53.
- `filebrowser-fde-v1.4.0.5.tar`: RepoTag v1.4.0.5, built-in Version v1.4.0.5, no Bug I.
- `filebrowser-fde-v1.4.0.6.tar`: RepoTag v1.4.0.6, built-in Version EMPTY, CSS contains
  `flex-direction:column!important` x2 (Bug I). Built 07-16 06:19 after rollback.
- NOTE: dist JS/CSS are stored gzipped inside the layer; plain grep misses them.
  Must gunzip inner assets to see markers. Recorded as a pitfall.

### Deliverables
- DEPLOY.md (449 lines, Traditional Chinese): server-side build model, prerequisites,
  build command with VERSION/REVISION args, config.yaml, docker-compose, verification list,
  update flow, rollback (3 scenarios), resource requirements + 3 low-memory workarounds
  (add swap / build elsewhere + transfer / slim variant), 11 troubleshooting entries,
  image variant comparison, path and port reference.
- README.md: added prominent "deployment model changed" notice pointing to DEPLOY.md;
  quick-deploy section rewritten for the new flow; old tar instructions moved to a
  historical-record section; tech stack deployment line updated.
- Commit a2302c66 pushed to origin/main (fast-forward, no force).

### Release Operations (executed)
- Tag `v1.4.0.5-fde-hotfix` created at main HEAD a2302c66 and pushed.
- Release `v1.4.0.5-fde-hotfix` created (id 399096490), asset uploaded:
  filebrowser-fde-v1.4.0.5-fde-hotfix.tar (87,063,040 bytes), digest
  sha256:abf54ddfa5f6844d36692425ecec15e2d550440c1c9d1aab750d89c7308a4980
- Old release `v1.4.0.5-fde-stable` voided by converting to DRAFT and prepending a
  deprecation notice (reversible; asset preserved). No asset was deleted.
- GitHub now reports LATEST = v1.4.0.5-fde-hotfix.
- Post-release verification: asset downloadable (HTTP 200), Content-Length matches,
  digest matches local file, tar structure valid.

### Pitfall Recorded
- GitHub asset upload must target `https://uploads.github.com/...`, NOT `api.github.com`.
  Using the API host returns HTTP 404.
- Large asset upload (83 MB) exceeds the 600s foreground limit in this environment;
  run it as a background process.

### Status
- All 6 PWF phases complete and archived.
- MAIN locked to v1.4.0.5-fde-hotfix tag (commit a2302c66).

---

## Session: 2026-09-29 — Upstream Upgrade Scout (v1.5.6-stable)

### Role & Scope (MASTER directive)
- Role: 上游升级侦察员 (upstream upgrade scout).
- Mission: find upstream latest stable tag as target = v1.5.6-stable.
- Baseline: v1.4.0.5-fde-hotfix tag (MAIN locked).
- DO NOT implement. Deliverable = readiness plan only (value / conflict / how).
- End state: PWF records + research doc + handoff, then PAUSE project for next AI group.

### Recon executed (READ-ONLY)
- Loaded skills: planning-with-files, cross-version-codebase-comparison,
  git-fork-maintenance, filebrowser-quantum-docker-setup.
- Started PWF project filebrowser-fde (was paused).
- Added upstream remote (git remote add upstream git@gtsteffaniak/filebrowser).
- Fetched upstream tags: v1.5.1..v1.5.6-stable, v2.0.x-beta, v2.0.x-preview.
- Confirmed target v1.5.6-stable = 5d9b4df2 (2026-09-04); latest commit is GHSA-55mw-cwg7-m8f5 security fix.
- Measured divergence: merge-base 21dfc616 ("Beta/v1.4.3"), 39 upstream commits, 84 fork commits, non-FF.
- Simulated 3-way merge (git merge-tree): 23 files / 30 hunks real conflicts.
- Mapped GHSA fixes: 9 security patches in the upstream drift.
- Confirmed fork core assets (7 files) untouched by upstream.

### Deliverables produced
- docs/v1.5.6-upgrade-scout-report.md (10.9KB, Traditional Chinese): full
  readiness assessment — value, conflict surface, top 3 semantic conflicts,
  9-step implementation plan, gates, open questions.
- findings.md: appended "Upstream Upgrade Scout — v1.5.6-stable" section.
- progress.md: this session log.

### Conclusion (deterministic wording per MASTER preference)
- Upgrade is FEASIBLE (measured, not "maybe").
- Conflicts QUANTIFIED: 23 files / 30 hunks (measured, not "maybe conflict").
- PRIMARY value = 9 GHSA security fixes (fork baseline behind 8+ CVEs).
- fork core assets safe: 7 files zero conflict.
- Biggest risk: pinned-items feature collision (upstream natively merged it).

### Handoff & pause
- Project is paused; next AI group reads docs/v1.5.6-upgrade-scout-report.md
  + task_plan phase 7 (scout) before implementing the upgrade.

---

## Session: 2026-09-29 (cont.) — v1.5.6.1-fde Upgrade IMPLEMENTATION

### Role & Scope (MASTER directive)
- Role: 实施组 (implementation group), NOT spec-writer.
- Mission: upgrade v1.4.0.5-fde-hotfix -> new branch v1.5.6.1-fde via Option C
  (three-way merge), preserving fork patches + core assets + 9 GHSA fixes.
- DO NOT write spec. DOUBT-FIRST full probe, then implement.
- BUILD policy: NO local docker build. Local = edit + commit + push only.
  MASTER updates repo on remote server, docker builds there, uses new image.

### Phase 8 results (conflict deep-dive + pinned ruling) — COMPLETED
- CORRECTED scout report errors:
  1. Real conflict = 20 files (report said 23). 3 "extra" are non-conflicting
     both-changed-different-region merges needing post-merge verify (httpRouter.go,
     resource.go, mutations.js).
  2. mutations.js: NEITHER side has togglePinnedItem (report was wrong). Real
     conflict = upstream setObjectProperty deep-copy refactor.
  3. pinned "collision" = FALSE. upstream natively CONTAINS fork's pinned
     superset (share/public-path support). Ruling: upstream everywhere.
- MISSED 8th core asset found: backend/database/users/users.go adds 8 persisted
  settings fields (ImagePreload/ImageTransition/ImageTapNav/PersistentNavButtons/
  NavButtonOpacity/CompressLevel/CompressQuality/CompressBackup) — MUST BE OURS.
- Final matrix: 20 conflict files -> 17 upstream / 1 ours (users.go) + 7 core
  assets / merge-both (Prompts.vue, i18n x3, httpRouter.go routes).
- Deliverables: findings.md Phase 8/8b sections + docs/phase8-pinned-ruling.md.

### Environment (this session)
- go 1.26.2 available (compiler verification).
- node v24.14.1 + pnpm 10.33.2 via fnm (frontend build verification).
- swag NOT installed (swagger codegen — handle in conflict resolution).
- Local go build / pnpm build are COMPILER VERIFICATION, not docker build
  (BUILD policy = no docker image build locally).

### Phase 10/11 results (branch + merge + conflict resolution) — COMPLETED
- Branch v1.5.6.1-fde created from a2302c66 (baseline tag deref confirmed equal).
- Pre-merge commit 64ad7181: upgrade prep docs (scout report + phase8 ruling).
- `git merge v1.5.6-stable` executed -> 22 conflict files (matches matrix).
- Resolution (all verified, no leftover markers anywhere):
  - UPSTREAM (--theirs): adapters/files/files.go, common/utils/file.go, go.mod,
    go.sum, http/pinnedItems.go, swagger x3, api/utils.test.js, ContextMenu.vue,
    files/FileList.vue, utils/object.js, utils/sort.js, views EpubViewer +
    OnlyOfficeEditor + Files.vue.
  - OURS (--ours): database/users/users.go (8 settings fields), README.md.
  - MERGE-BOTH (manual): Prompts.vue (fork ExtractToFolder+CompressImages import
    restored after --theirs wiped them), i18n en/zh-cn/zh-tw (compressImages +
    sendToApp both kept).
- 7 core assets verified byte-identical to fork baseline (0 diff).
- httpRouter.go auto-merged cleanly: fork 3 compress routes + upstream
  share/pinnedItems + media/subtitles routes all present.

### CRITICAL find during build (fork deps dropped by upstream go.mod)
- fork compress.go imports 4 third-party libs NOT in upstream go.mod:
  deepteams/webp v1.2.7, esimov/colorquant v1.0.0, klauspost/compress/zstd,
  go-logger/logger. Taking upstream go.mod DROPPED them -> go build failed.
- FIX: go get webp@v1.2.7 + colorquant@v1.0.0 + klauspost/compress/zstd,
  then go mod tidy. This is the "upstream deps + fork core asset" intersection
  that MUST be preserved. Documented for future upgrades.

### Phase 12 results (build verification) — COMPLETED
- go build ./... : PASS (exit 0)
- go vet ./... : PASS (exit 0)
- go mod verify : PASS (all modules verified)
- go test ./... : PASS (all packages ok, 0 FAIL)
- npm install : PASS (417 packages; 1 high-severity audit warning noted, NOT
  introduced by merge — pre-existing)
- npm run typecheck (vue-tsc) : PASS (0 errors)
- npm run test (vitest) : PASS (9 files / 50 tests)
- npm run build (vite) : PASS (1021 modules, 19.53s)
- GHSA markers verified in tree: JoinScopedIndexPath/BoundIndexPath (symlink),
  ReplaceAll backslash sanitize, shareStore rename.
- fork compress feature verified in built bundle: compressImages x7 + 
  CompressImages x1 markers in index bundle.

### Phase 13 results (deploy doc + version history) — COMPLETED
- README.md: bump latest stable -> v1.5.6.1-fde; add version history entry
  (9 GHSA + FDE features preserved + public share pinned); 基准版本 note updated.
- DEPLOY.md: version bump (v1.5.6.1-fde); added MANDATORY data-dir backup
  step before upgrade (tar -czf backup-data-<ts>.tar.gz ./data); added
  rollback scenario A-2 (restore data + image); checkout-tag instructions
  in step 1; git fetch --tags added to update flow.
- swagger: REGENERATED via `go tool swag init --output swagger/docs`
  (swag v1.16.6 project-locked). fork compress routes (3) now in docs.
  func init block removed per makefile sed rule. go build still PASS.

### Git commits (v1.5.6.1-fde branch)
- 64ad7181: docs — upgrade prep (scout report + phase8 ruling + PWF tracking)
- f498fe3d: merge — upstream v1.5.6-stable into v1.5.6.1-fde (22 conflicts)
- 1dd3b1cc: docs — version history + DEPLOY.md backup/rollback + swagger regen

### Final verification summary (ALL PASS)
- go build ./... ✅ | go vet ✅ | go mod verify ✅ | go test ✅ (all ok)
- npm typecheck (vue-tsc) ✅ | npm test (vitest 50/50) ✅ | npm build (1021 modules) ✅
- 7 core assets byte-identical to fork baseline ✅
- 9 GHSA markers present in tree ✅
- 0 leftover conflict markers ✅
- fork compress deps (webp/colorquant/zstd) re-added + build passes ✅

### STATUS: v1.5.6.1-fde branch PUSHED to origin (2026-09-29)

### BUILD OOM 诊断（npm run build:docker FAIL, exit 134, node OOMErrorHandler）

症状：远程 docker build 在 `RUN npm run build:docker`（Dockerfile:22）崩溃，
node `v8::internal::V8::FatalProcessOutOfMemory` → Aborted (core dumped) exit 134。

根因链（已确诊）：
1. `build:docker` = `vite build`（无 rm/cp，纯 vite 打包）
2. vite.config.ts 启用了 `vite-plugin-checker` 且 `checker({ vueTsc: {...} })`，
   build 环境 (`isDevBuild=false`) 下 checker 生效。checker 的 `enableBuild`
   默认值 = `true`（实测 node_modules/vite-plugin-checker/dist/types.d.ts line 189），
   → vite build 时额外拉起 vue-tsc 类型检查进程（独立 node 进程，吃内存）。
3. VueI18nPlugin `runtimeOnly: false` → 全量预编译 30+ i18n JSON。
4. 远程服务器内存不足 → node V8 heap OOM。

本地为何通过：本地 3.8GB RAM + 7.8GB swap，内存充足。

方案：
- 方案 A（根治，推荐）：vite.config.ts 给 checker 加 `enableBuild: false`，
  部署 build 不需要类型检查（类型检查是开发期/CI 的事）。
- 方案 B（治标）：远程加临时/永久 swap（DEPLOY.md 已写完整步骤）。
- 方案 C（组合）：A + B。

待 MASTER 拍板选哪个方案。

### 方案 A 实施（enableBuild: false）— DONE + PUSHED
- MASTER 拍板方案 A（根治）。改 frontend/vite.config.ts checker 加
  `enableBuild: false`（附注释说明为何 build 要跳过类型检查）。
- 本地验证 `npm run build:docker`（= vite build，远程崩的同一步）：
  PASS 15.3s（原 19.53s，快 4s，无 checker/vue-tsc 日志 → 类型检查已跳过）。
- 验证 `npm run typecheck` 独立命令仍可用且通过（类型检查没废，
  只是不再绑进 production build）。
- commit 4057ecf5 pushed to origin/v1.5.6.1-fde（git ls-remote 确认 ref 一致）。
- 远程服务器规格：2G RAM + 4G SWAP（主人 2026-09-29 实测提供）。
  跳过后 checker 后，vite build 峰值应落在 1-1.5GB，配合 4G swap 应可过。
  Go CGO 编译（golang:alpine + mupdf/musl）是另一层峰值 ~1-1.5GB。
- 主人正在远程重跑 docker build（进行中，等结果）。
- git tag `v1.5.6.1-fde` 未打。若主人要 Release，之後再打 tag + GitHub Release。
- Full remote build flow documented in DEPLOY.md (backup + update + build + rollback).

### 远程 BUILD 结果（2026-09-29 主人实测）— 升级 FAIL，暂停
- 远程 docker build 两次都 FAIL：
  1. 默认配置（2G RAM）→ OOM
  2. 加临时 4G SWAP → 仍 OOM
- 主人关键观察：**「中前期就有错误出现，然后再 OOM」**——不是单纯
  vite build OOM 那么简单，OOM 之前已有其他错误输出。方案 A
  （enableBuild:false）已实施但未解决，说明根因另有他处，需下一 AI 组
  拿到远程完整 build log（尤其 OOM 之前的中前期报错行）再诊断。
- 本地（3.8G RAM + 7.8G swap）docker build 实测成功（sanity check 全过：
  filebrowser version / ffmpeg / ffprobe / exiftool 都 exit 0），证明
  代码 + Dockerfile 本身可变，问题集中在远程环境（内存或 Docker 层差异）。
- 交接给下一 AI 组，暂停本项目。
- Full remote build flow documented in DEPLOY.md (backup + update + build + rollback).
- `git push -u origin v1.5.6.1-fde` done; verified via git ls-remote
  (refs/heads/v1.5.6.1-fde = 1dd3b1cc). Tracking set.
- Next step (MASTER on remote server): git fetch origin -> git checkout
  v1.5.6.1-fde -> docker build -> docker compose up.
- Full remote build flow documented in DEPLOY.md (backup + update + build + rollback).
- git tag `v1.5.6.1-fde` 未打。若主人要 Release，之後再打 tag + GitHub Release。

---

## Session: 2026-09-30 — Local Build + tar.zst Delivery Model

### MASTER Directives
1. DOUBT FIRST: investigate why remote docker build kept failing (OOM even with 2G RAM + 8G swap)
2. Cancel remote build model entirely — switch to LOCAL docker build + tar.zst delivery
3. Clean debug image, rebuild with correct naming, package as tar.zst level 5
4. Rewrite DEPLOY.md for new delivery model (docker load on remote, no build)
5. Clean builder cache to free disk space after build

### DOUBT FIRST Analysis — Why Remote Build Failed
- Previous AI never obtained remote build logs (the "early errors before OOM" MASTER observed)
- Dockerfile has hidden memory traps: full devDeps npm install (no --production),
  VueI18nPlugin runtimeOnly:false (precompiles 30+ i18n JSONs), compression plugin
  adds gzip layer during build, Go CGO_ENABLED=1 + mupdf/musl needs 1-1.5GB alone
- 2G RAM + swap thrashing: swap is 10-50x slower than RAM, Go CGO compilation
  under swap thrashing causes timeout or OOM Killer
- CONCLUSION: Remote build on 2G server is architecturally infeasible, not a tweak issue

### PWF Updates
- Goal updated: LOCAL docker build + tar.zst delivery model (remote does NOT build)
- Decision #12 updated: new BUILD policy (local build, not remote)
- Notes #13/14/17/18 updated: new deployment model (tar.zst transfer + docker load)
- Title updated: FileBrowser Quantum - v1.5.6.1-fde

### Phase 14 Execution (ALL COMPLETED)

#### 14.1 Remove debug image
- docker rmi filebrowser-fde:v1.5.6.1-fde-debug — DONE

#### 14.2 Remove stale build cache
- docker builder prune -f — 7.574GB reclaimed

#### 14.3 Local docker build
- docker build --build-arg VERSION=v1.5.6.1-fde --build-arg REVISION=4057ecf5
  -t filebrowser-fde:v1.5.6.1-fde -f _docker/Dockerfile .
- BUILD SUCCESS (exit 0), all stages passed

#### 14.4 Image verification (ALL PASS)
- Version: v1.5.6.1-fde ✅
- Commit: 4057ecf5 ✅
- ffmpeg 9.0.1 ✅
- ffprobe 9.0.1 ✅
- exiftool 13.55 ✅
- Image size: 85.1MB (content), 298MB (disk)

#### 14.5 Docker save + zstd compress
- docker save | zstd -5 -o filebrowser-fde-v1.5.6.1-fde.tar.zst
- Result: 81MB (85,087,232 bytes)
- Compression ratio: 99.67% (tar already highly compressible)

#### 14.6 Integrity verification
- zstd -t: PASS ✅
- SHA256: e5d92beeb7f260fd9f4800f95090b86fb544a15d2b1bfa9de0a80a66fdc92bce

#### 14.7 DEPLOY.md rewrite
- Full rewrite for tar.zst transfer + docker load model
- Removed all remote build instructions
- Added: SCP transfer, SHA256 verification, docker load steps

#### 14.8 .gitignore update
- Added *.tar.zst pattern

#### 14.9 Git commit + push
- Commit 46c9a146: docs + .gitignore
- Pushed to origin/v1.5.6.1-fde

#### 14.10 Build cache cleanup
- docker builder prune --all -f — total ~8GB freed
- Build Cache: 0B (fully cleaned)
- Disk: 22G available (was 14G before session)

#### 14.11 Deliverable summary
- File: /home/hermes/workspace/filebrowser-fde/filebrowser-fde-v1.5.6.1-fde.tar.zst
- Size: 81MB (85,087,232 bytes)
- SHA256: e5d92beeb7f260fd9f4800f95090b86fb544a15d2b1bfa9de0a80a66fdc92bce
- Image: filebrowser-fde:v1.5.6.1-fde (commit 4057ecf5)
- Deployment: transfer tar.zst to remote -> docker load -> docker compose up

### Docker images retained
- filebrowser-fde:v1.4.0.5-fde-hotfix (316MB) — previous stable, kept for rollback
- filebrowser-fde:v1.5.6.1-fde (298MB) — new release

### git tag v1.5.6.1-fde: NOT yet created. MASTER decides when to tag + Release.
