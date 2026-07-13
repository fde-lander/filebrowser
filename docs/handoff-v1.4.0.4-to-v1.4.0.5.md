# Handoff Document

Generated: 2026-07-13
Next session focus: Implement v1.4.0.5 hotfix - rebuild compress system (background goroutine + 3-second polling, ABANDON SSE), folder expansion, backup persistence, preview fix, progress accuracy

---

## ⭐ CRITICAL: 第一个 AI 必须读嘅文档

**首要必读**：`/home/hermes/workspace/filebrowser-fde/docs/v1.4.0.5-next-version-plan.md`

呢份文档详细记录咗 v1.4.0.4 测试发现嘅所有问题 + 完整嘅修复方案架构。**开工前必读**，里面有 7 个章节、问题根因、修复代码示例、架构图、12 项测试清单。

**⚠️ 最重要：v1.4.0.4 SSE 系统永久性损坏** - 一次断开后所有新压缩尝试 1-2 秒就失败。**必须彻底废弃 SSE，改用 3 秒轮询**（详见 v1.4.0.5 文档 Ch2 + Ch6）。

---

## Context

**项目**：FileBrowser Quantum 自定义 Fork（基于 v1.4.0-stable）

**当前版本**：v1.4.0.4-hotfix（已经部署测试，已发现关键 bug）

**下一版本**：v1.4.0.5（需要根据测试反馈重建压缩系统）

**目标服务器**：1 vCPU core, 2GB RAM, 2GB ZRAM + 2GB SWAP（单核服务器，必须考虑长时间压缩）

**Docker volume 映射**：`/root/qbb/downloads:/folder`（容器内路径）

**测试反馈总结**：
- ✅ 过渡 fade (soft) 模式睇漫画效果好
- ⚠️ 压缩功能有重大缺陷（见下方 "Remaining Tasks"）
- ✅ 预览模式压缩预览图 + 百分比显示 OK
- ⚠️ 预览原图非常模糊（getPreviewURL 默认返回缩略图）
- ⚠️ 文件夹右键压缩只显示 1 条记录（应展开为内部所有图片）
- ⚠️ backupEnabled 默认不打勾，应默认 ON 或持久化

---

## Key Decisions

### 决策 1：彻底废弃 SSE，改用 3 秒轮询

**Reasoning**：v1.4.0.4 测试发现 SSE 断开后整个压缩系统永久瘫痪（所有新尝试 1-2 秒就失败）。SSE handler goroutine 可能永久阻塞 + progressMgr channel 死锁 + 资源泄漏。MASTER 明确要求：**整个压缩系统必须重建**，前端用 3 秒轮询取代 SSE。

### 决策 2：compressJobManager 单例

**Reasoning**：压缩任务必须在后端独立于 HTTP request 生命周期。JobManager 持有 taskID、total、processed、success、skipped、failed、backupPath、backupFallback、status 等状态。前端关闭弹窗 = clearInterval（停止轮询），后端继续运行。重新打开弹窗时调 GET /compress-images/status 查询当前状态。

### 决策 3：单实例强制

**Reasoning**：单核服务器同时只能压缩 1 组，否则 ZRAM/SWAP thrashing。Backend compressHandler 调 compressJobManager.isRunning() 检查；如果已有 job 运行返回 409 Conflict。前端根据 status API 禁用 Compress 按钮。

### 决策 4：进度发送改为处理后（After Processing）

**Reasoning**：原代码 progress 发送在前、处理在后，导致 processed 计数不准（用 loop index 而非 success+skipped+failed 实际计数）。改为先处理文件再发 progress + 用准确 processed 计数。

### 决策 5：预览原图加 size=original

**Reasoning**：CompressImages.vue openPreview() 调用 getPreviewURL() 但冇加 size 参数，导致返回低质量缩略图。Preview.vue 嘅 raw() computed 有 `&size=original` 嘅正确做法。CompressImages.vue 必须加上。

### 决策 6：备份 toggle 默认 ON + 持久化

**Reasoning**：通过 NonAdminEditable struct 加 CompressBackup 字段（参考 deleteAfterArchive 模式）。前端 mounted() 读取 `state.user.compressionSettings?.compressBackup ?? true`，@change 触发 mutations.updateUserField 持久化。

### 决策 7：前端展开文件夹

**Reasoning**：当前 items prop 直接来自 state.selected，文件夹对象不会被展开。后端 compressHandler 已有目录展开逻辑但前端列表只显示文件夹本身。需要前端在 mounted() 调 fetchListing API 展开所有文件夹为内部图片文件。

### 决策 8：3 个过渡模式改为同时渐变

**Reasoning**：所有 3 种模式（instant、crossfade、fade）都因 decode 异步导致黑闪。改为先 waitForDecode → toRef 已 decode 完 → 同时渐变两图（fade 模式 500ms 同时渐变，唔再中间黑屏）。fade_to_black 改名为 fade（向后兼容两者都支持）。

---

## Remaining Tasks

### Priority 1 - CRITICAL：v1.4.0.5 后台压缩系统（详见 v1.4.0.5 文档 Ch2）

- [ ] 创建 compressJobManager 单例（compress.go package-level struct with mutex）
- [ ] compressHandler 调 isRunning() 检查并发拒绝（返回 409）
- [ ] 新增 GET /compress-images/status endpoint（httpRouter.go + compress.go）
- [ ] 移除 SSE 全部代码（删除 compressProgressHandler 或改为 status polling）
- [ ] 删除 subscribeProgress frontend（api/compress.js）
- [ ] 前端实现 pollStatus() 函数：3 秒 setInterval + clearInterval

### Priority 2 - HIGH：备份 toggle 持久化（Ch3）

- [ ] backend/database/users/users.go NonAdminEditable 加 CompressBackup bool
- [ ] CompressImages.vue mounted() 读取 store + ?? true fallback
- [ ] ToggleSwitch @change 调 mutations.updateUserField

### Priority 3 - HIGH：文件夹展开（Ch4）

- [ ] CompressImages.vue expandItems() 方法：fetchListing 展开文件夹
- [ ] isImagePath() helper（jpg/jpeg/png/bmp/webp）
- [ ] 模板改用 expandedItems 替代 items
- [ ] 添加 expandingItems loading state
- [ ] 新 i18n key: prompts.compressScanning

### Priority 4 - HIGH：预览原图模糊（Ch5）

- [ ] CompressImages.vue openPreview() 加 `&size=original` 到 getPreviewURL
- [ ] 或 fallback 用 getDownloadURL(file.source, file.path, true)
- [ ] 全屏 overlay 同步使用全分辨率 URL

### Priority 5 - HIGH：进度准确性 + 移除 SSE（Ch6）

- [ ] compress.go progress 发送改为处理后（用 success+skipped+failed 准确计数）
- [ ] 删除 progressMgr 全部 SSE channel 代码
- [ ] 改为 compressJobManager.updateProgress() 直接更新状态

### 测试清单（v1.4.0.5 完成后）

- [ ] 压缩 250 张图，关闭弹窗，重开 -> 显示运行中进度
- [ ] 压缩完成后即使弹窗关闭也完成
- [ ] 同时只 1 个 job，2nd attempt 显示 "already running"
- [ ] Backup toggle 跨 session 持久化（默认 ON）
- [ ] 文件夹选择展开为所有图片
- [ ] 多文件夹：所有图片列出，按文件夹分组
- [ ] 预览原图全分辨率（不模糊）
- [ ] 全屏显示全分辨率
- [ ] 进度计数准确（匹配实际文件数）
- [ ] 新用户默认设置正常工作

---

## Blockers

无。但需注意：

1. **无本地 npm/node 环境** - 前端构建只喺 Docker 多阶段构建内进行。所有改动必须理论上验证。
2. **无本地测试** - MASTER 部署到远程服务器后才测试。每次部署成本高。
3. **单核 2GB RAM 限制** - 250+ 图片压缩时间长（可能几分钟），架构必须支持长时间后台运行。
4. **SSE 旧代码可能彻底损坏** - 调查 progressMgr 实现时可能发现死锁，需要彻底删除 SSE 系统而非修补。

---

## Artifacts

### ⭐ 必读文档

- **v1.4.0.5 详细计划**：`/home/hermes/workspace/filebrowser-fde/docs/v1.4.0.5-next-version-plan.md` — **第一个 AI 必读！** 7 个章节 + 完整代码示例 + 架构图 + 12 项测试清单

### 设计文档

- Spec：`/home/hermes/.hermes/docs/superpowers/specs/2026-07-13-v1.4.0.4-hotfix-design.md` — v1.4.0.4 设计（已实施）
- Plan：`/home/hermes/.hermes/docs/superpowers/plans/2026-07-13-v1.4.0.4-hotfix-plan.md` — v1.4.0.4 实施计划（已实施）

### 项目文档（项目目录内）

- **README**：`/home/hermes/workspace/filebrowser-fde/README.md` — 已重写为自定义版本文档（fork URLs + 功能清单 + 版本表）
- **CHANGELOG**：`/home/hermes/workspace/filebrowser-fde/CHANGELOG-custom.md` — 完整历史（v1.4.0.1 到 v1.4.0.4）

### PWF（项目目录内）

- `task_plan.md`：当前 Phase 3 状态
- `progress.md`：详细 session log
- `findings.md`：技术调研结果

### 修改过嘅源代码（v1.4.0.4 base）

- `backend/http/compress.go` — finishEvent struct + backup path resolve + 3-level fallback + compressPreviewHandler dir expansion
- `frontend/src/api/compress.js` — SSE event name fix + previewCompress blob rewrite
- `frontend/src/components/prompts/CompressImages.vue` — SSE handlers + preview UI redesign
- `frontend/src/components/files/ExtendedImage.vue` — waitForDecode + decode-first navigateToImage + 3-mode swapBuffers
- `frontend/src/views/settings/Profile.vue` — fade transition option
- `frontend/src/i18n/{en,zh-cn,zh-tw}.json` — new keys

### Git 状态

- Branch：`v1.4.0.4-hotfix`
- Remote main：已推送（force push，含 v1.4.0.5 计划文档 + README + tar 文件）
- Tag：`v1.4.0.4-latest` 已推送
- GitHub：https://github.com/fde-lander/filebrowser
- Default branch：main

### LanceDB 记忆 IDs

- `f13604b0`：v1.4.0.3 测试结果
- `d96afe8c`：v1.4.0.2 bug status
- `499be97c`：v1.4.0.4 测试反馈（8 个问题）
- `53200a29`：v1.4.0.5 计划概览
- `a75e91ad`：CRITICAL - SSE 永久性损坏 + 轮询方案

### Docker 镜像

- filebrowser-fde:v1.4.0.4（已部署）
- filebrowser-fde-v1.4.0.4.tar（84MB，已 git push + 部署用）
- filebrowser-fde-v1.4.0.3.tar（84MB，已 git push）
- filebrowser-fde-v1.4.0.2.tar（84MB，已 git push）
- filebrowser-fde-custom.tar（82MB，已 git push，v1.4.0.1）

---

## Environment Notes

- Go binary：`/usr/local/go/bin/go`（不在默认 PATH）
- Go module root：`/home/hermes/workspace/filebrowser-fde/backend/`
- Docker build：`docker build --build-arg="VERSION=v1.4.0.5" --build-arg="REVISION=$(git rev-parse --short HEAD)" -t filebrowser-fde:v1.4.0.5 -f _docker/Dockerfile .`
- Docker save：`docker save filebrowser-fde:v1.4.0.5 -o filebrowser-fde-v1.4.0.5.tar`
- 无 sudo docker 权限问题：用 `sudo bash -c 'docker build ...'`
- 本地磁盘清理：已清理旧 images + build cache，腾出 14GB
- v1.4.0.4 tar 文件必须保留（不能删）

### Known working patterns

- Permission persistence：参考 `deleteAfterArchive` 模式（users.go NonAdminEditable）
- Image preview URL：参考 Preview.vue raw() 加 `&size=original`
- Dual-buffer transition：ExtendedImage.vue decode-first 架构
- i18n keys：en.json + zh-cn.json + zh-tw.json 必须同步
- 用户偏好持久化：通过 mutations.updateUserField 写入

---

## Recommended Skills

下一个 AI 开工前应加载：

- **brainstorming** — 如果对 v1.4.0.5 架构有疑问，需要再次澄清设计
- **writing-plans** — 实现前需写详细实施计划（参考 v1.4.0.4 plan 模式）
- **executing-plans** — 实施计划
- **systematic-debugging** — 调查 progressMgr 实现（找 SSE 死锁根因）时使用
- **test-driven-development** — Go backend 适用（前端用 grep 结构验证）

---

## ⚠️ 重要警告

1. **SSE 系统永久性损坏** - 不要尝试修补 SSE，必须彻底删除并替换为轮询
2. **目标服务器 1 核 2GB RAM** - 长时间压缩必须支持，ZRAM/SWAP thrashing 风险
3. **不要修改主仓库 default branch** - 保持 main 为 default
4. **tar 文件必须保留** - 部署用，唔可以删
5. **使用 3 个 subagent 并行** - 代码研究时 max 3 并发
6. **所有改动必须理论验证** - 无本地测试环境，go build + grep 验证是唯一手段