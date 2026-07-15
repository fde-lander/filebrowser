# Handoff Document - FileBrowser Quantum v1.4.0.6-fde Implementation

Generated: 2026-07-15
Next session focus: Execute the v1.4.0.6-fde implementation plan (10 tasks, deterministic)

---

## 你是谁，你要做什么

你是接手 FileBrowser Quantum v1.4.0.6-fde 实施工作的 AI。上一个 AI 已经完成了：
- 完整的代码调研（DOUBT 原则：3 个 subagent 调研 + 主 Agent 独立验证）
- 6 轮设计讨论（与 MASTER 确认所有方案）
- 设计文档 Spec（10 章，1230 行）
- 实施计划 Plan（10 个 Task，2537 行，84KB）

**你的任务是：按照实施计划文档逐个 Task 执行，每个 Task = 一个 git commit。**

---

## 第一步：必读文档（按顺序读完）

**1. 实施计划（最重要 - 你的工作指引）**
/home/hermes/.hermes/docs/superpowers/plans/2026-07-15-v1.4.0.6-fde-plan.md
2537 行，84KB。包含 10 个 Task，每个 Task 有精确的 old_string/new_string、完整代码、grep 验证、commit message。
**你必须逐字逐行读懂每个 Task 的每个 Step，严格按照 old_string/new_string 执行。**

**2. 设计文档（理解设计决策）**
/home/hermes/.hermes/docs/superpowers/specs/2026-07-15-v1.4.0.6-fde-design.md
1230 行，44KB。包含 10 章设计 spec，记录了所有根因分析、方案选择、MASTER 原话。

**3. PWF 状态文件（了解项目状态）**
- /home/hermes/workspace/filebrowser-fde/task_plan.md - 当前 Phase 4 completed
- /home/hermes/workspace/filebrowser-fde/findings.md - 调研结果（3 subagent + 验证）
- /home/hermes/workspace/filebrowser-fde/progress.md - 完整 session log

**4. FileBrowser 项目 Skill（了解项目约定）**
skill_view("filebrowser-quantum-docker-setup") - 包含部署流程、代码架构、post-subagent bug 模式等

---

## 项目概况

- **项目**：FileBrowser Quantum 自定义 Fork
- **Fork**：github.com/fde-lander/filebrowser.git
- **路径**：/home/hermes/workspace/filebrowser-fde
- **分支**：v1.4.0.6-fde（已创建，从 v1.4.0.5-fde @ 2c6b4490）
- **当前版本**：v1.4.0.5-fde-hotfix（tag: v1.4.0.5-fde-hotfix-latest）
- **目标版本**：v1.4.0.6-fde
- **Go binary**：/usr/local/go/bin/go（不在默认 PATH，用前 export PATH=/usr/local/go/bin:$PATH）

---

## 10 个 Task 概览

| # | Task | Priority | Files | Steps |
|---|------|----------|-------|-------|
| 1 | Bug G: 移动端点击翻页失效 | TIER 1 | ExtendedImage.vue | 10 |
| 2 | Bug H: 移动端翻页按钮双重导航 | TIER 1 | nextPrevious.vue | 12 |
| 3 | Bug I: 移动端压缩预览布局 | P1 | CompressImages.vue CSS | 4 |
| 4 | 后端 compress 控制系统 | P0 | compress.go | 15 |
| 5 | 新路由注册 | P0 | httpRouter.go | 4 |
| 6 | 前端 API 函数 | P0 | compress.js | 3 |
| 7 | 确认弹窗组件 | P1 | ConfirmAction.vue + Prompts.vue | 5 |
| 8 | 全局状态条组件 | P0 | CompressStatusBar.vue + Prompts.vue | 5 |
| 9 | 压缩对话框控制按钮 | P0 | CompressImages.vue | 8 |
| 10 | 设置系统 + i18n | P1 | 7 files | 14 |

**执行顺序**：Task 1-2 先做（Tier 1），然后 3-10 按顺序。每个 Task 是独立 commit，可单独回滚。

---

## 关键决策（MASTER 已确认）

1. **Bug B（进度状态丢失）**：方案 C 两层设计 - 全局状态条 + CompressImages mounted 检查
2. **Bug G（移动端点击失效）**：移除 @touchmove.prevent，改条件性 preventDefault。保护 swipe + transition 不变
3. **Bug H（移动端黑屏闪烁）**：根因是 double navigation + transitioning overlay。修复：touchState.triggered 检查 + 图片类型跳过 transitioning
4. **Bug I（预览布局破碎）**：@media (max-width: 768px) 响应式 CSS
5. **队列进度**：累计统计 totalFiles/totalProcessed/batchCount/currentBatchIndex
6. **取消范围**：取消整个队列（唔系当前批次）
7. **跳过当前批次**：新增功能，只有 queueLength > 0 先显示
8. **跳过/取消**：需要二次确认弹窗（ConfirmAction.vue）
9. **暂停自动超时**：可设置（0=禁用，5-120 分钟，默认 30），跨 session 持久化
10. **独立 commit**：每个 Task 一个 commit，方便单独回滚

---

## CRITICAL 警告（必须遵守）

1. **Transition engine 锁定**：navigateToImage / swapBuffers / finishTransition 完全唔好改。Bug G/H 只改事件绑定，唔改 transition 逻辑
2. **PC 点击翻页 = TIER 1**：所有修复必须零影响 PC 点击翻页体验。plan 入面每个 Task 有 PC Impact Analysis
3. **唔好用 silent catch{}**：所有 catch 要加 console.error()
4. **唔好删 legacy CSS**：标注 LEGACY 注释就得，唔好删
5. **i18n 3 文件同步**：en.json, zh-cn.json, zh-tw.json 必须同步
6. **每次 patch < 2.5KB**：避免 LLM context overload
7. **无本地测试**：只能 Docker build -> save -> scp -> MASTER 部署测试。所以代码必须理论验证通过先 build
8. **HARD-GATE**：没有 MASTER 批准绝对不能 sudo 或提权
9. **Docker tar 唔好删**：tar 文件系部署工件
10. **MASSER 用 Matrix 端**：输出禁止用代码块、Markdown 标题、表格。用 emoji + 粗体 + 分隔线

---

## 实施后验证

Plan 文档末尾有完整的 Post-Build Verification Checklist：
- go build + go vet + go mod verify
- 14 个 grep 验证（每个有预期数量）
- 3 个 JSON 验证
- commit count 检查（预期 10 个）
- Docker build + save + prune

---

## Blockers

无 blockers。所有调研已完成，所有设计已确认，所有 old_string 已用实际代码验证。

---

## 环境信息

- **服务器**：2 核 + 4GB RAM（Ryzen 9 7950X 分配）
- **目标部署服务器**：1 核 2GB RAM + 2GB ZRAM + 2GB SWAP（远程）
- **Docker volume**：/root/qbb/downloads:/folder
- **Docker build 流程**：docker build -> docker save -> scp -> docker load -> docker-compose up -d
- **Codebase 已索引**：codebase-memory-mcp 有 3927 节点 / 13640 边

---

## Artifacts

- **Design Spec**: ~/.hermes/docs/superpowers/specs/2026-07-15-v1.4.0.6-fde-design.md - 10 章，1230 行，44KB
- **Implementation Plan**: ~/.hermes/docs/superpowers/plans/2026-07-15-v1.4.0.6-fde-plan.md - 10 Task + checklist，2537 行，84KB
- **PWF task_plan**: /home/hermes/workspace/filebrowser-fde/task_plan.md - Phase 4 completed
- **PWF findings**: /home/hermes/workspace/filebrowser-fde/findings.md - 3 subagent 调研 + 验证结果
- **PWF progress**: /home/hermes/workspace/filebrowser-fde/progress.md - 完整 session log
- **Previous handoff**: /home/hermes/workspace/filebrowser-fde/docs/handoff-v1.4.0.6-fde.md - 上一任调研 AI 的交接文档
- **Implementation handoff**: /home/hermes/workspace/filebrowser-fde/docs/handoff-v1.4.0.6-fde-implementation.md - 实施交接文档（本文件）
- **v1.4.0.6 需求文档**: /home/hermes/workspace/filebrowser-fde/docs/v1.4.0.6-next-version-plan.md - 772 行，9 章，MASTER 原话
- **Git state**: branch v1.4.0.6-fde, base v1.4.0.5-fde-hotfix-latest (2c6b4490)

---

## Recommended Skills

- **superpowers/executing-plans**: 按计划逐 Task 执行
- **superpowers/subagent-driven-development**: 如果需要并行执行（但本 plan 建议串行，因为 Task 之间有依赖）
- **filebrowser-quantum-docker-setup**: 项目特定约定、部署流程、bug 模式
- **codebase-memory-cli**: 用 codebase 查询代码关系（已索引）
- **superpowers/systematic-debugging**: 如果 go build 失败或 grep 验证唔通过

---

## LanceDB Memory References

- ID 1a6f69b5: Bug B failure lesson（搜索 "FileBrowser Bug B mounted"）
- ID f13604b0: v1.4.0.3 test results
- ID 499be97c: v1.4.0.4 MASTER feedback
- ID a75e91ad: SSE permanent failure bug

---

## 最后提醒

1. **先读完 Plan 文档再开始**。不要跳读。
2. **每个 Step 都有 old_string/new_string**。直接复制到 patch 工具用。不要自己写代码。
3. **如果 old_string 匹配唔到**：read_file 读当前代码确认实际内容，调整 old_string。
4. **每个 Task 完成后**：运行 grep 验证 + go build（后端 Task）+ commit。
5. **所有 Task 完成后**：运行 Post-Build Verification Checklist。
6. **最后**：Docker build -> save -> prune。等 MASTER 部署测试。
7. **绝对不要**：改 transition engine、删 legacy CSS、用 silent catch、跳过 grep 验证。
