# FileBrowser Quantum V1.4.0 - Extract to New Folder Patch

**Date:** 2026-07-11  
**Author:** fde-lander  
**Base Version:** v1.4.0-stable  
**Branch:** custom/v1.4.0-extract-to-folder  
**Fork:** https://github.com/fde-lander/filebrowser.git  
**Status:** Initial test passed ✅

---

## 功能概述

新增「解壓到新資料夾」右键菜单功能。用户右键压缩包时，可选择此功能自动创建以压缩包名命名的文件夹并解压到其中，省去手动新建文件夹的步骤。

## 改动原因

原版解压流程需要 5 步操作（右键 -> 选解压 -> 点路径 -> 新建文件夹 -> 选择 -> 解压），过于繁琐。常见用例只需要「解压到同名文件夹」，因此新增一键完成的功能。

## 改动范围

1 个新文件 + 8 个修改文件，共 454 行新增代码。

### 后端（Go）

**文件：backend/http/archive.go**

- `unarchiveRequest` struct 新增 `CreateSubfolder bool` 字段
- 新增 `archiveFolderName()` 函数：从压缩包路径提取文件夹名（去除 .zip/.tar.gz/.tgz 扩展名）
- 新增 `resolveUniqueSubfolderName()` 函数：冲突时自动加后缀 ` (1)`、` (2)`... 最多 100 次
- `unarchiveHandler` 新增逻辑：当 `createSubfolder == true` 时，自动创建子文件夹再解压
- Swagger 注释更新：新增 createSubfolder 参数说明

### 前端（Vue/JS）

**新增文件：frontend/src/components/prompts/ExtractToFolder.vue**（340 行）

- Dup 自 Unarchive.vue，关键差异：
  - `mounted()` 预填 destPath 为 `parentPath + folderName + "/"`
  - `submit()` 发送 `createSubfolder: true` 和 `destination: parentPath`
  - 使用 `extractToFolder*` i18n key
  - 成功通知使用 `extractToFolderSuccess` 文案

**修改文件：frontend/src/components/ContextMenu.vue**

- 新增 2 个 template action（主菜单 + overflow 菜单）
- 新增 2 个 computed（showExtractToFolder + showExtractToFolderInOverflow）
- 新增 3 个 method（showExtractToFolderPrompt + FromPreview + openExtractToFolderPrompt）
- 更新 hasOverflowItems 加入新 action

**修改文件：frontend/src/components/prompts/Prompts.vue**

- 新增 import ExtractToFolder
- 新增 components 注册
- 新增 getDisplayTitle switch case

**修改文件：frontend/src/api/archive.js + resources.js**

- `unarchive()` 函数新增 `createSubfolder` 参数和 body 字段

### i18n（3 语言 × 3 key = 9 条）

| Key | en | zh-cn | zh-tw |
|-----|-----|-------|-------|
| extractToFolder | Extract to new folder | 解压到新文件夹 | 解壓到新資料夾 |
| extractToFolderDestination | Extract to: | 解压到： | 解壓到： |
| extractToFolderSuccess | Archive extracted to new folder successfully! | 压缩包已解压到新文件夹！ | 封存檔已解壓到新資料夾！ |

复用现有 key：`profileSettings.deleteAfterArchive`、`buttons.goToItem`、`general.cancel` 等。

## 技术设计要点

- **后端单次 API 调用**：createSubfolder 字段让后端一次完成建文件夹+解压，避免前端两次请求的竞态风险
- **冲突处理**：os.Stat 检测 + 自动加后缀，类似 Windows 资源管理器行为
- **deleteAfterArchive 复用**：使用现有用户偏好持久化机制，跨 session 记住用户选择
- **不影响原有功能**：createSubfolder 默认 false（Go zero-value），原有「解壓縮封存檔」功能完全不变
- **安全机制复用**：normalizeArchiveEntryName、safeExtractPath、symlinkTargetStaysUnderDest 全部不变

## Git 提交记录

| # | Commit | 说明 |
|---|--------|------|
| 1 | 4c2c4392 | feat(backend): add createSubfolder field, helpers, handler logic, Swagger docs |
| 2 | e6c39e54 | feat(frontend): add createSubfolder parameter to unarchive API functions |
| 3 | b6d8a924 | feat(i18n): add extractToFolder keys for en, zh-cn, zh-tw |
| 4 | 0da1830b | feat(frontend): add ExtractToFolder.vue prompt component |
| 5 | 8c2f8044 | feat(frontend): register ExtractToFolder in ContextMenu and Prompts |

## Docker 构建

- **Image tag:** filebrowser-fde:custom
- **Tar 文件:** filebrowser-fde-custom.tar (82MB)
- **构建命令:** `docker build --build-arg="VERSION=custom-v1.4.0-extract-to-folder" --build-arg="REVISION=<sha>" -t filebrowser-fde:custom -f _docker/Dockerfile .`

## 部署方式

1. `scp filebrowser-fde-custom.tar user@server:/tmp/`
2. `docker load -i /tmp/filebrowser-fde-custom.tar`
3. docker-compose.yaml 中 image 改为 `filebrowser-fde:custom`
4. `docker compose down && docker compose up -d`

## 测试结果

- ✅ Go build PASS
- ✅ Docker build PASS
- ✅ 初步功能测试：解压到新文件夹 + 自动删除原压缩包
- ✅ 代码审查：无涉及图片预览/加载逻辑

## 回滚方式

docker-compose.yaml 中 image 改回 `gtstef/filebrowser:1.4.0-stable`，restart 即可。

## 相关文档

- Spec: /home/hermes/.hermes/docs/superpowers/specs/2026-07-11-extract-to-folder-design.md
- Plan: /home/hermes/.hermes/docs/superpowers/plans/2026-07-11-extract-to-folder-plan.md
- PWF: /home/hermes/workspace/filebrowser-fde/ (task_plan.md, findings.md, progress.md)

---

# V1.4.0.2 - Image Viewer Enhancement + Image Compression

**Date:** 2026-07-11
**Author:** fde-lander
**Base Version:** v1.4.0-stable (locked, no upstream changes)
**Branch:** v1.4.0.2-image-viewer-compression (merged to main)
**Fork:** https://github.com/fde-lander/filebrowser.git
**Status:** Built, pending remote server testing

---

## 功能概述

两大功能升级：

**功能 A - 图片查看器增强：** 消除翻页黑屏闪烁，新增双缓冲预加载、3 张缓存池、渐变过渡效果、点击翻页、移动端按钮持久化 + 透明度调节。

**功能 B - 图片压缩：** 右键批量压缩图片，支持 3 档压缩级别、实时预览、ZSTD 自动备份、SSE 进度推送。

## 改动原因

**图片查看器：** 原版翻页时 ExtendedImage 组件被销毁重建，导致黑屏闪烁。移动端导航按钮 3 秒后消失，翻页不便。

**图片压缩：** 部分图片文件过大（10+MB），需要批量压缩节省空间。需要保底备份机制确保安全。

## 改动范围

4 个新建文件 + 12 个修改文件，6 个 commit。

### 新建文件

| 文件 | 说明 |
|------|------|
| backend/http/compress.go (676 行) | 压缩 API：encoder 选择 + 预览 handler + 执行 handler + ZSTD 备份 + SSE 进度 |
| frontend/src/components/prompts/CompressImages.vue (708 行) | 压缩弹窗：文件列表 + 档位 + 预览 + 进度 |
| frontend/src/api/compress.js | 前端 API 层：previewCompress + startCompress + subscribeProgress |
| frontend/src/components/settings/RangeSlider.vue | 可复用滑块组件 |

### 修改文件

| 文件 | 改动 |
|------|------|
| backend/database/users/users.go | +8 字段到 NonAdminEditable |
| backend/http/httpRouter.go | +3 个 API 路由 |
| backend/go.mod + go.sum | +3 个 Go 依赖 |
| frontend/src/views/files/Preview.vue | 移除 v-if 销毁 ExtendedImage |
| frontend/src/components/files/ExtendedImage.vue | 双缓冲 + 缓存池 + 过渡引擎 + 点击翻页 |
| frontend/src/components/files/nextPrevious.vue | 按钮持久化 + 透明度 |
| frontend/src/components/ContextMenu.vue | +压缩图片右键选项 |
| frontend/src/components/prompts/Prompts.vue | 注册 CompressImages |
| frontend/src/views/settings/Profile.vue | +图片查看器设置区域 |
| frontend/src/i18n/{en,zh-cn,zh-tw}.json | +~30 条 i18n × 3 语言 |
| _docker/Dockerfile | +apk add oxipng pngquant |

### 新增 Go 依赖

- github.com/deepteams/webp v1.2.7（WebP 编码，pure Go）
- github.com/klauspost/compress/zstd（ZSTD 备份，pure Go）
- github.com/esimov/colorquant（PNG 调色板，pure Go）

## 技术设计要点

### 图片查看器

- **双缓冲架构：** imgA/imgB 两个 img 元素，翻页时交叉显示，组件永不被销毁
- **3 张缓存池：** Map 结构缓存上一张+当前+下一张，来回翻页零等待
- **渐变过渡引擎：** 可扩展注册表模式，支持 crossfade/fade_to_black/instant，未来可加新效果
- **点击翻页：** 左右各 40% 区域，200ms 双击检测延迟，拖拽 >10px 取消翻页
- **按钮持久化：** 移动端可设按钮不消失，独立透明度滑块，hover 时临时提升

### 图片压缩

- **统一编码器：** 所有档位用 WebP（deepteams/webp），低档 PNG 特殊路径用 OxiPNG + palette reduction
- **3 档预设：** 低档（30-50%）、中档（55-75%）、高档（75%+），基于实测数据
- **保底逻辑：** 压缩后如 ≥ 原图大小则跳过，保留原图
- **ZSTD 备份：** tar.zst 格式，放在操作发起路径
- **SSE 进度：** 实时推送当前文件名、已处理/总数、百分比
- **设置记忆：** 档位、quality、备份开关持久化到用户设置

## Git 提交记录

| # | Commit | 说明 |
|---|--------|------|
| 1 | 600afcff | feat(backend): add image viewer + compression settings fields + Go dependencies |
| 2 | 030770aa | feat(frontend): add dual-buffer + 3-image cache pool + transition engine + tap navigation to ExtendedImage |
| 3 | e40604df | feat(frontend): add CompressImages dialog + compress API layer + i18n keys for en/zh-cn/zh-tw |
| 4 | 906493bc | feat(backend): add compress.go - encoder + preview + execute + ZSTD backup + SSE progress + routes |
| 5 | e38e1646 | feat(frontend): Preview.vue keep ExtendedImage mounted + nextPrevious persistence/opacity + RangeSlider + ContextMenu/Prompts/Profile integration |
| 6 | b9433cf2 | feat(docker): add oxipng + pngquant to Dockerfile for PNG optimization |

## Docker 构建

- **Image tag:** filebrowser-fde:v1.4.0.2
- **Tar 文件:** filebrowser-fde-v1.4.0.2.tar (84MB)
- **构建命令:** docker build --build-arg="VERSION=v1.4.0.2" --build-arg="REVISION=<sha>" -t filebrowser-fde:v1.4.0.2 -f _docker/Dockerfile .

## 部署方式

1. scp filebrowser-fde-v1.4.0.2.tar user@server:/tmp/
2. docker load -i /tmp/filebrowser-fde-v1.4.0.2.tar
3. docker-compose.yaml 中 image 改为 filebrowser-fde:v1.4.0.2
4. docker compose down && docker compose up -d
5. 按 docs/deploy-v1.4.0.2.md 中 23 项测试清单验证

## 验证结果

- ✅ go build PASS
- ✅ go vet PASS
- ✅ go mod verify PASS
- ✅ Docker build PASS
- ⏳ 远程服务器功能测试：待主人部署后验证

## 压缩档位实测数据

| 档位 | 源格式 | 编码器 | 默认 Q | 实测压缩率 |
|------|--------|--------|--------|-----------|
| 低档 | JPEG | WebP | 75 | 33-38% |
| 低档 | PNG | OxiPNG -o2 + palette | N/A | 40-72% (无损) |
| 中档 | JPEG | WebP | 65 | 45-50% |
| 中档 | PNG | WebP | 90 | ~94% (视觉无损) |
| 高档 | JPEG | WebP | 55 | 55-65% |
| 高档 | PNG | WebP | 75 | ~97% |

## 相关文档

- Spec: ~/.hermes/docs/superpowers/specs/2026-07-11-image-viewer-compression-design.md
- Plan: ~/.hermes/docs/superpowers/plans/2026-07-11-image-viewer-compression-plan.md
- Deploy: docs/deploy-v1.4.0.2.md
- PWF: task_plan.md, findings.md, progress.md

---

# v1.4.0.3 - Bug Fix & Permission Enhancement

**Date:** 2026-07-12
**Author:** fde-lander
**Base Version:** v1.4.0-stable
**Branch:** v1.4.0.2-image-viewer-compression
**Docker Image:** filebrowser-fde:v1.4.0.3 (84MB)
**Status:** Ready for deployment

---

## 修复内容

### Bug A (CRITICAL): 图片快速翻页时过渡不稳定

**问题：** 快速翻页几次后图片消失，显示缩略图 + 转圈，然后放大。三种过渡模式都受影响。

**根因：** JS 双缓冲状态机有 5 个 race condition：
1. doTransition() 无重入保护
2. toImg.complete 浏览器兼容问题
3. stale onload closure
4. imgA @load 与 doTransition onload 双重触发
5. 缓存池存 detached Image 但从未复用

**修复方案：** 架构重写 - 放弃 JS 管理过渡动画，改用 CSS transition + generation token
- 新方法 navigateToImage() 替代 doTransition()
- 新方法 swapBuffers() 用 CSS opacity 管理三种过渡模式
- generation token 防止 stale callback
- 缓存池简化为 Set<string>
- onLoad() 加 guard 防止双重触发

### Bug B (CRITICAL): 压缩图片 API 全部 404

**问题：** 预览、选档位、调 quality 都返回 404。

**根因：** 前端 API 路径、字段名、参数格式全部与后端不匹配。

**修复：** 前端 api/compress.js 全面对齐后端：
- URL 路径：resources/compress/* -> compress-images/*
- 字段名：tier -> level
- 文件格式：[{path}] -> ["path1", ...]
- 进度参数：jobId -> taskId

### Bug C (MEDIUM): 右键文件夹无「压缩图片」选项

**问题：** 右键文件夹不显示压缩菜单。

**根因：** showCompressImages 检查 item.isDir，但部分 item 对象使用 item.type === 'directory'。

**修复：** 双重检测：item.isDir || item.type === 'directory'

### Bug D (HIGH): 备份功能损坏

**问题：** backupFileName 只显示不发送，后端 os.Create("") 失败，扩展名错误，不支持目录递归。

**修复：**
- backupFileName 改为 .tar.zst 扩展名，多文件加 _and_N_others 后缀
- 新增 backupPath computed，计算同层目录
- doCompress 传递 backupPath + backupName 到后端
- 后端 addFileToTar 支持目录递归（filepath.Walk）
- 后端 compressHandler 目录展开为图片文件列表

### 新增：Admin 权限门控

**需求：** 解压到新文件夹 + 压缩图片功能必须 Admin 权限才能使用。

**实现：**
- 后端：compressHandler/compressPreviewHandler/compressProgressHandler 加 Admin 检查
- 后端：unarchiveHandler 从 Create 改为 Admin
- 前端：showCompressImages/showExtractToFolder 从 permissions.create 改为 permissions.admin

## 改动文件清单

1. backend/http/compress.go - isImageFile, addFileToTar recursion, dir expansion, Admin checks
2. backend/http/archive.go - unarchiveHandler Create -> Admin
3. frontend/src/api/compress.js - URL/field/format/param alignment
4. frontend/src/components/prompts/CompressImages.vue - backup naming, API params
5. frontend/src/components/ContextMenu.vue - folder detection, Admin permission
6. frontend/src/components/files/ExtendedImage.vue - transition architecture rewrite

## Git Commits

- b19e1960 fix: backend directory recursion + Admin permission gate for compress/extract
- f6f12a22 fix: align frontend compress API with backend routes and fields
- e90f89a2 fix: backup naming + API field alignment in CompressImages.vue
- be842d25 fix: folder detection + Admin permission gate for compress/extract
- ccd433f5 fix: rewrite image transition to CSS-managed architecture with generation token

## 相关文档

- Spec: ~/.hermes/docs/superpowers/specs/2026-07-11-bugfix-permission-design.md
- Plan: ~/.hermes/docs/superpowers/plans/2026-07-11-bugfix-permission-plan.md
- PWF: task_plan.md, findings.md, progress.md

---

# v1.4.0.4 - Hotfix: Backup Path + Preview UI + SSE + Transition

**Date:** 2026-07-13
**Author:** fde-lander
**Base Version:** v1.4.0-stable
**Branch:** v1.4.0.4-hotfix
**Docker Image:** filebrowser-fde:v1.4.0.4 (84MB)
**Status:** Ready for deployment

---

## 修复内容

### Issue 1: Backup Path Resolution + 3-Level Fallback (CRITICAL)

**问题：** 压缩功能从未成功执行。backupPath 是用户空间路径，后端直接用 filepath.Join 拼接后 os.Create，路径不存在导致失败，backup-first 设计导致整个压缩中止。

**修复：**
- goroutine 启动前用 resolveCompressPath 解析 backupPath 到真实文件系统路径
- 3 级退避：同级目录 -> 上一层目录 -> source 根目录
- 全部失败则中止压缩
- finishEvent 新增 BackupFallback 字段通知前端

### Issue 2: compressPreviewHandler 目录展开 (CRITICAL)

**问题：** 右键文件夹预览时 500 错误，os.ReadFile 尝试读取目录。

**修复：** 在 resolveCompressPath 后添加 os.Stat + IsDir 检查，目录则 filepath.Walk 取第一张图片做预览。

### Issue 3: SSE Event Name + Field Alignment (HIGH)

**问题：**
- 后端发 event: finish 但前端监听 "complete" -> onComplete 永不触发
- progress 字段 current 是文件名但前端用作计数器
- complete 事件读 data.totalFiles / data.totalSaved 但后端发 success/skipped/failed

**修复：**
- compress.js: "complete" -> "finish"
- onProgress: 映射 processed->current, current->currentFile
- onComplete: 读 success/skipped/failed + backupPath + backupFallback
- onError: 读 data.error 兼容 data.message

### Issue 4: Preview UI Redesign (HIGH)

**问题：** 后端返回二进制 blob + headers 但前端读 response.json() -> 必定失败。预览从未工作。

**修复：**
- previewCompress 改为读 response.blob() + headers
- 新增 checkbox「开启预览模式」
- 预览模式 ON: 点击文件打开 overlay（左原图 + 右压缩预览 + loading spinner + fade-in）
- 点击图片 -> 全屏 overlay（黑色背景 + 返回按钮）
- 原图 URL 用 getPreviewURL() 从 resources.js
- 压缩图用 blob URL（blob: 协议）+ beforeUnmount cleanup
- 预览模式 OFF: 正常勾选/取消选择

### Issue 5: Transition decode-first Architecture (MEDIUM)

**问题：** 三种过渡模式都在 toRef 未 decode 完就隐藏 fromRef -> 黑闪。crossfade 300ms 对大图不够，fade_to_black 有 150ms 设计性黑屏。

**修复：**
- 新增 waitForDecode() helper: decode().then() + 3s 超时 + fallback for old browsers
- navigateToImage 重写: set src -> waitForDecode -> swapBuffers（统一路径，不再分 cache hit/miss）
- swapBuffers 三种模式全部改为「decode 完成后同时渐变」-> 永不到黑
- fade_to_black 改名为 fade（500ms 柔和渐变，不再有黑屏）
- 保留 fade_to_black 向后兼容
- CSS transform 冲突修复: translate -> translate3d（恢复 GPU 合成）

## 改动文件清单

1. backend/http/compress.go - finishEvent struct + backup path resolve + 3-level fallback + compressPreviewHandler dir expansion
2. frontend/src/api/compress.js - SSE event name + previewCompress blob rewrite
3. frontend/src/components/prompts/CompressImages.vue - SSE handlers + preview UI + overlay + fullscreen + methods
4. frontend/src/components/files/ExtendedImage.vue - waitForDecode + decode-first navigateToImage + 3-mode swapBuffers + CSS fix
5. frontend/src/views/settings/Profile.vue - fade transition option
6. frontend/src/i18n/en.json - new keys
7. frontend/src/i18n/zh-cn.json - new keys
8. frontend/src/i18n/zh-tw.json - new keys

## Git Commits

- 0bc9a11d fix: add BackupFallback field to finishEvent struct
- 468faf8f fix: resolve backup path through resolveCompressPath + 3-level fallback
- 21fd2cdf fix: add directory expansion to compressPreviewHandler
- 83841304 fix: align SSE event name + rewrite previewCompress to blob
- 65f9b84a feat: preview UI redesign - checkbox + overlay + fullscreen + SSE mapping
- e454b209 fix: decode-first transition architecture - waitForDecode + simultaneous crossfade
- 04a87848 feat: add i18n keys + fade transition option

## 验证结果

- go build: PASS
- go vet: PASS
- go mod verify: PASS
- JSON validation: 3/3 PASS
- Grep sweep: ALL markers correct
- Docker build: PASS (84MB)
- Docker save: PASS (filebrowser-fde-v1.4.0.4.tar, 84MB)

## 部署方式

1. scp filebrowser-fde-v1.4.0.4.tar user@server:/tmp/
2. docker load -i /tmp/filebrowser-fde-v1.4.0.4.tar
3. docker-compose.yaml 中 image 改为 filebrowser-fde:v1.4.0.4
4. docker compose down && docker compose up -d

## 相关文档

- Spec: ~/.hermes/docs/superpowers/specs/2026-07-13-v1.4.0.4-hotfix-design.md
- Plan: ~/.hermes/docs/superpowers/plans/2026-07-13-v1.4.0.4-hotfix-plan.md
- PWF: task_plan.md, findings.md, progress.md
