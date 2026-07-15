# Handoff Document - FileBrowser Quantum v1.4.0.7 Hotfix

Generated: 2026-07-16
Next session focus: Investigate and fix v1.4.0.6 failed features (Bug G/H/B + compress control) with proper codebase research before any implementation

---

## 你是谁，你要做什么

你是接手 FileBrowser Quantum v1.4.0.7 Hotfix 工作的 AI。上一个 AI（超级猪兔兔）完成了 v1.4.0.6 的实施，但部署测试后大部分功能失败，已全部回滚。当前代码库 = v1.4.0.5-fde-hotfix + Bug I @media CSS（仅 3 行）。

v1.4.0.6 失败原因 + 回滚状态 + 下次调研方向全部记录在 v1.4.0.7 hotfix plan 文档中。你的任务：先读文档完整掌握，再用 codebase-memory + subagent 充分调研每个问题，最后设计修复方案给 MASTER 批准后实施。

---

## 第一步：必读文档（按顺序读完）

**1. v1.4.0.7 Hotfix Plan（最重要 - 你的工作指引）**
/home/hermes/workspace/filebrowser-fde/docs/v1.4.0.7-hotfix-plan.md
8 章，涵盖：v1.4.0.6 完整失败分析、回滚状态、每个 issue 的调研方向、调研协议、验证 checklist。

**2. v1.4.0.6 需求文档（了解原始需求和 MASTER 原话）**
/home/hermes/workspace/filebrowser-fde/docs/v1.4.0.6-next-version-plan.md
772 行，9 章。记录了所有 Bug G/H/I + 压缩控制功能的需求和 MASTER 测试反馈。

**3. PWF 状态文件**
/home/hermes/workspace/filebrowser-fde/task_plan.md - 当前 Phase 5 状态
/home/hermes/workspace/filebrowser-fde/findings.md - 3 subagent 调研结果 + 验证
/home/hermes/workspace/filebrowser-fde/progress.md - 完整 session log（含 v1.4.0.6 失败记录）

**4. FileBrowser 项目 Skill（项目约定 + 部署流程）**
skill_view("filebrowser-quantum-docker-setup") - 包含部署流程、代码架构、post-subagent bug 模式等

**5. 之前的实施计划（了解 v1.4.0.6 做过什么、为什么失败）**
/home/hermes/.hermes/docs/superpowers/plans/2026-07-15-v1.4.0.6-fde-plan.md
2537 行。10 个 Task。注意：这份 plan 有严重缺陷 - CompressStatusBar 冇触发代码、checkBackendStatus 放喺 items guard 后面。参考佢嘅代码结构但唔好重复佢嘅错误。

**6. 之前的设计文档（理解设计决策的来源）**
/home/hermes/.hermes/docs/superpowers/specs/2026-07-15-v1.4.0.6-fde-design.md
1230 行，10 章 spec。注意：设计本身有部分正确但实施时出咗错。

---

## 项目概况

- **项目**：FileBrowser Quantum 自定义 Fork
- **Fork**：github.com/fde-lander/filebrowser.git
- **路径**：/home/hermes/workspace/filebrowser-fde
- **分支**：v1.4.0.6-fde（已 push 到 origin/main）
- **当前代码状态**：v1.4.0.5-fde-hotfix + Bug I @media CSS（仅 3 行差异）
- **稳定 Release**：v1.4.0.5-fde-stable tag（已发布，含 Docker tar）
- **Go binary**：/usr/local/go/bin/go（不在默认 PATH，用前 export PATH=/usr/local/go/bin:$PATH）
- **gh CLI**：已安装 v2.96.0，已登录 fde-lander 账号

## 当前代码差异（v1.4.0.5-fde-hotfix vs HEAD）

唯一改动：frontend/src/components/prompts/CompressImages.vue 的 @media block 加了 3 行 CSS：
- compress-preview-overlay: flex-direction column + overflow-y auto
- compress-preview-item: max-height 30vh
- compress-preview-section: flex-direction column
作用：移动端竖屏压缩预览布局修复（Bug I）。同翻页功能零重叠。

---

## v1.4.0.6 失败总结（详见 v1.4.0.7-hotfix-plan.md Chapter 7）

### 5 个根本失败原因

1. **CompressStatusBar 冇触发代码** - 组件建了注册了但冇任何地方调用 showPrompt({ name: "compressStatusBar" })。Plan gap。
2. **checkBackendStatus 放喺 items guard 后面** - mounted() 先 return 咗，checkBackendStatus 永远行唔到。同 v1.4.0.5 同样错误。
3. **setNavigationTransitioning skip 有未知副作用** - Bug H Fix 2 跳过图片的 setNavigationTransitioning，导致导航按钮完全唔显示 + tap 翻页过渡质量下降。冇调查所有 isTransitioning 消费者。
4. **冇做触发链路追踪** - 所有验证只做 grep + go build，冇追踪「乜嘢触发乜嘢显示」。
5. **tap vs swipe 过渡差异未调查** - 唔知点解 swipe 丝滑但 tap 唔丝滑。

### 7 个待解决问题（v1.4.0.7 范围）

| # | Issue | Severity | 状态 |
|---|-------|----------|------|
| 1 | Tap 翻页过渡质量下降 | P1 HIGH | 需调研 tap vs swipe 代码路径差异 |
| 2 | 导航按钮完全唔显示 | P0 CRITICAL | 需调研 isTransitioning 所有消费者 |
| 3 | 压缩状态栏从未触发 | P0 CRITICAL | 需设计触发点 + 实现触发代码 |
| 4 | 控制按钮唔显示 | P0 CRITICAL | checkBackendStatus 需移到 items guard 之前 |
| 5 | 对话框异常消失 | P0 CRITICAL | checkBackendStatus 设 compressing=true 导致替换 UI |
| 6 | 有队列就无法再压缩 | P0 CRITICAL | compressing=true 隐藏了压缩按钮 |
| 7 | 当前文件名无显示 | P1 HIGH | 需验证 API response currentFile 字段 |

---

## 关键约束（CRITICAL）

1. **Transition engine 锁定**：navigateToImage / swapBuffers / finishTransition 绝对唔好改。v1.4.0.5 改过一次搞烂咗，回滚咗。v1.4.0.6 又改过一次又搞烂咗。两次都失败。第三次必须用完全唔同嘅方法。
2. **PC 点击翻页 = TIER 1**：所有修复零影响 PC
3. **唔好用 silent catch{}**：必须加 console.error()
4. **唔好删 legacy CSS**：标注 LEGACY 注释就得
5. **i18n 3 文件同步**：en.json, zh-cn.json, zh-tw.json
6. **每次 patch < 2.5KB**
7. **无本地测试**：只能 Docker build -> save -> scp -> MASTER 部署测试
8. **HARD-GATE**：没有 MASTER 批准绝对不能 sudo 或提权
9. **Docker tar 唔好删**：tar 文件系部署工件
10. **MASTER 用 Matrix 端**：输出禁止用代码块、Markdown 标题、表格。用 emoji + 粗体 + 分隔线
11. **subagent 调研结果必须主 Agent 验证**：读实际代码确认每个 line number
12. **CLI 命令必须验证**：唔好凭记忆脑补参数，用前先 --help 确认

---

## 环境信息

- **服务器**：2 核 + 4GB RAM（Ryzen 9 7950X 分配）
- **目标部署服务器**：1 核 2GB RAM + 2GB ZRAM + 2GB SWAP（远程）
- **Docker volume**：/root/qbb/downloads:/folder
- **Docker build 流程**：docker build -f _docker/Dockerfile -t filebrowser-fde:v1.4.0.6 . -> docker save -> scp -> docker load -> docker-compose up -d
- **Dockerfile 位置**：_docker/Dockerfile（唔喺 root 目录！）
- **Codebase 已索引**：codebase-memory-mcp 有 3927 节点 / 13640 边
- **gh CLI**：v2.96.0 已安装，已登录 fde-lander（GITHUB_TOKEN via env）

---

## 调研工具（必须使用）

1. **codebase-memory MCP**（已索引）：
   - search_graph: 搜索函数/类/路由
   - trace_path: 追踪调用路径（calls, data_flow, cross_service）
   - get_code_snippet: 读取特定函数代码
   - query_graph: Cypher 查询

2. **Subagent delegation**（并行调研）：
   - 每个 subagent 报告 exact line numbers + actual code snippets
   - 主 Agent 必须读实际代码验证所有发现

3. **Session search**：搜索之前嘅调研结果

4. **Skills（必须加载）**：
   - skill_view("filebrowser-quantum-docker-setup") - 项目约定 + 部署
   - skill_view("planning-with-files") - PWF 状态管理
   - skill_view("superpowers/executing-plans") - 计划执行
   - skill_view("superpowers/systematic-debugging") - 系统调试
   - skill_view("codebase-memory-cli") - codebase 查询工具
   - skill_view("dogfood/cli-command-verification-rule") - CLI 命令验证

---

## Artifacts

- **v1.4.0.7 Hotfix Plan**: /home/hermes/workspace/filebrowser-fde/docs/v1.4.0.7-hotfix-plan.md - 8 章完整失败分析 + 调研方向
- **v1.4.0.6 需求文档**: /home/hermes/workspace/filebrowser-fde/docs/v1.4.0.6-next-version-plan.md - 772 行原始需求 + MASTER 原话
- **v1.4.0.6 实施 Plan**: /home/hermes/.hermes/docs/superpowers/plans/2026-07-15-v1.4.0.6-fde-plan.md - 2537 行（有缺陷，参考但唔好重复错误）
- **v1.4.0.6 设计 Spec**: /home/hermes/.hermes/docs/superpowers/specs/2026-07-15-v1.4.0.6-fde-design.md - 1230 行
- **PWF task_plan**: /home/hermes/workspace/filebrowser-fde/task_plan.md
- **PWF findings**: /home/hermes/workspace/filebrowser-fde/findings.md - 3 subagent 调研 + 验证
- **PWF progress**: /home/hermes/workspace/filebrowser-fde/progress.md - 完整 session log
- **Docker tar**: /home/hermes/workspace/filebrowser-fde/filebrowser-fde-v1.4.0.6.tar (84MB)
- **GitHub Release**: https://github.com/fde-lander/filebrowser/releases/tag/v1.4.0.5-fde-stable
- **Git state**: branch v1.4.0.6-fde, HEAD = fbe45b11, tag v1.4.0.5-fde-stable on 2c6b4490
- **Codebase index**: 3927 nodes / 13640 edges (可能需要 re-index 因为代码已回滚)

---

## 部署流程

1. Docker build: cd /home/hermes/workspace/filebrowser-fde && docker build -f _docker/Dockerfile -t filebrowser-fde:v1.4.0.6 .
2. Docker save: docker save filebrowser-fde:v1.4.0.6 -o filebrowser-fde-v1.4.0.6.tar
3. Build cache prune: docker builder prune -f
4. MASTER: scp tar 到目标服务器 -> docker load -> docker-compose up -d

---

## 最后提醒

1. **先读完所有文档再开始**。不要跳读。
2. **充分调研再动手**。v1.4.0.6 失败嘅核心原因就系调研唔够 + 乱写代码。
3. **用 codebase-memory 追踪每个代码路径**。唔好凭假设。
4. **Subagent 调研结果必须验证**。读实际代码确认。
5. **每个新组件都要追踪「乜嘢触发佢显示」**。唔好建咗组件但冇触发代码。
6. **checkBackendStatus 唔好放喺 items guard 后面**。呢个错误犯咗两次。
7. **Transition engine 唔好改**。改过两次都搞烂。
8. **isTransitioning 所有消费者都要调查清楚**。唔好跳过任何一个。
9. **设计方案要 MASTER 批准先可以实施**。
10. **绝对不要**：改 transition engine、删 legacy CSS、用 silent catch、sudo 未问主人、跳过调研。
