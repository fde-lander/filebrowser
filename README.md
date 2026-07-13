<div align="center">

  [![Apache-2.0 License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)

  <h3>FileBrowser Quantum (Custom Fork)</h3>
  自定义增强版 - 基于 FileBrowser Quantum v1.4.0-stable<br/><br/>
</div>

## 关于本 Fork

本仓库是 [FileBrowser Quantum](https://github.com/gtsteffaniak/filebrowser) 的自定义增强版本，基于 v1.4.0-stable tag 创建。

**Fork 仓库**：https://github.com/fde-lander/filebrowser.git

**上游项目**：https://github.com/gtsteffaniak/filebrowser （原版 Quantum）

**官方文档**：https://filebrowserquantum.com/en/docs/getting-started/docker

本 Fork 不向上游提交 PR，所有改动仅供自用。改动以最小化侵入方式实现，保持与上游代码结构兼容。

## 快速部署

Docker 镜像通过 tar 文件分发，不推送到 Docker Hub。

**部署步骤：**

步骤 1：加载 Docker 镜像
- docker load -i filebrowser-fde-v1.4.0.4.tar

步骤 2：docker-compose.yaml 中 image 设置为 filebrowser-fde:v1.4.0.4

步骤 3：docker compose down && docker compose up -d

**当前版本**：v1.4.0.4

**最小配置（config.yaml）**：
```yaml
server:
  sources:
    - path: "/srv"
      config:
        defaultEnabled: true

auth:
  adminUsername: admin
  adminPassword: admin
```

详细配置请参考 [官方文档](https://filebrowserquantum.com/en/docs/getting-started/config)。

---

## 本 Fork 自定义功能清单

以下是本 Fork 相比原版 Quantum v1.4.0-stable 的全部自定义改动：

### v1.4.0.1 - 解壓到新資料夾

**功能**：右键压缩包时新增「解壓到新資料夾」选项，自动创建以压缩包名命名的文件夹并解压到其中。

- 后端：archive.go 新增 CreateSubfolder 字段 + archiveFolderName() + resolveUniqueSubfolderName()
- 前端：新增 ExtractToFolder.vue 组件 + ContextMenu 注册 + i18n（en/zh-cn/zh-tw）
- 冲突处理：同名文件夹自动加后缀 (1)、(2)... 最多 100 次
- deleteAfterArchive 复用：使用现有用户偏好持久化机制
- 安全：复用 normalizeArchiveEntryName + safeExtractPath + symlink 检查

### v1.4.0.2 - 图片查看器增强 + 图片压缩

**功能 A - 图片查看器增强**：
- 双缓冲架构：imgA/imgB 两个 img 元素交叉显示，组件永不被销毁
- 3 张缓存池：Set 结构缓存上一张+当前+下一张，来回翻页零等待
- 渐变过渡引擎：可扩展注册表模式，支持 crossfade / fade (soft) / instant
- 点击翻页：左右各 40% 区域，200ms 双击检测延迟，拖拽 >10px 取消翻页
- 移动端按钮持久化 + 独立透明度滑块
- 设置持久化：所有设置通过 NonAdminEditable struct 跨 session 记忆

**功能 B - 图片压缩**：
- 右键批量压缩：支持文件夹 / 多选文件
- 3 档压缩级别：低档（WebP Q75）、中档（Q65）、高档（Q55）
- PNG 特殊路径：低档 OxiPNG 无损 + 中高档 WebP
- ZSTD 自动备份：tar.zst 格式，压缩前先备份（backup-first 设计）
- 3 级退避：同级目录 -> 上一层 -> source 根目录
- 保底逻辑：压缩后如 ≥ 原图大小则跳过
- Admin 权限门控：仅 Admin 用户可见可操作
- i18n：en / zh-cn / zh-tw 三语言

### v1.4.0.3 - Bug Fix（过渡架构重写 + API 对齐 + 权限门控）

- 图片过渡架构重写：CSS 管理过渡动画 + generation token 防 stale callback
- 压缩 API 前后端对齐：URL / 字段名 / 参数格式全部以后端为准
- 文件夹右键检测：isDir || type === 'directory' 双重检测
- Admin 权限门控：解压到新文件夹 + 压缩图片均需 Admin 权限

### v1.4.0.4 - Hotfix（备份路径 + 预览 UI + 过渡 decode-first）

**Issue 1 - 备份路径解析 + 3 级退避**：
- 后端 backupPath 通过 resolveCompressPath 解析到真实文件系统路径
- 3 级退避：同级目录 -> 上一层目录 -> source 根目录
- 全部失败则中止压缩（backup-first 原则）
- finishEvent 新增 BackupFallback 字段通知前端

**Issue 2 - compressPreviewHandler 目录展开**：
- 右键文件夹预览不再 500 错误
- os.Stat + filepath.Walk 取第一张图片做预览

**Issue 3 - SSE 字段对齐**：
- SSE event name: "complete" -> "finish"（匹配后端）
- progress 字段映射：processed -> current, current -> currentFile
- complete 字段映射：success/skipped/failed + backupPath + backupFallback
- error 字段映射：data.error 兼容 data.message

**Issue 4 - 预览 UI 全新设计**：
- checkbox「开启预览模式」toggle
- 预览模式 ON：点击文件打开 overlay（左原图 + 右压缩预览 + loading + fade-in）
- 全屏 overlay：黑色背景 + 返回按钮 + 点击外部返回
- 原图用 getPreviewURL + blob URL（压缩图）+ beforeUnmount cleanup
- 预览模式 OFF：正常勾选/取消选择

**Issue 5 - 过渡 decode-first 架构**：
- 新增 waitForDecode() helper：decode().then() + 3s 超时 + 旧浏览器 fallback
- navigateToImage 重写：set src -> waitForDecode -> swapBuffers（统一路径）
- swapBuffers 三种模式全部改为「decode 完成后同时渐变」-> 永不到黑
- fade_to_black 改名为 fade（500ms 柔和渐变，不再有黑屏）
- 保留 fade_to_black 向后兼容
- CSS transform 冲突修复：translate -> translate3d（恢复 GPU 合成）

### 正在开发中（v1.4.0.5 计划）

详见 `docs/v1.4.0.5-next-version-plan.md`

- 后端后台压缩系统（CRITICAL）：compressJobManager 单例 + 废弃 SSE 改为 3 秒轮询
- 备份打勾持久化 + 默认 ON
- 文件夹文件列表展开（前端展开目录为内部图片列表）
- 预览原图模糊修复（&size=original）
- 进度准确性（处理后更新 + 准确计数）

---

## 技术栈

- 后端：Go 1.25（http handlers + SQLite index + WebP/OxiPNG/ZSTD）
- 前端：Vue 3 / Vite（双缓冲图片查看器 + 压缩弹窗 + SSE/轮询）
- Docker：3 阶段构建（Go backend -> Node frontend -> Alpine final）
- 部署：Docker tar save/load（不推 Docker Hub）

## 环境信息

- 目标服务器：1 vCPU core, 2GB RAM, 2GB ZRAM + 2GB SWAP
- Docker volume 映射：/root/qbb/downloads:/folder（容器内）
- 反向代理：Nginx Proxy Manager（Docker，跨服务器）

## 版本发布

| 版本 | Docker Image | Tar 文件 | 说明 |
| --- | --- | --- | --- |
| v1.4.0.1 | filebrowser-fde:custom | filebrowser-fde-custom.tar | 解压到新文件夹 |
| v1.4.0.2 | filebrowser-fde:v1.4.0.2 | filebrowser-fde-v1.4.0.2.tar | 图片查看器 + 压缩 |
| v1.4.0.3 | - | - | Bug fix（已含在后续版本） |
| v1.4.0.4 | filebrowser-fde:v1.4.0.4 | filebrowser-fde-v1.4.0.4.tar | Hotfix + decode-first |

## 上游项目致谢

本 Fork 基于以下优秀开源项目：
- [FileBrowser Quantum](https://github.com/gtsteffaniak/filebrowser) by gtsteffaniak
- [原版 FileBrowser](https://github.com/filebrowser/filebrowser)

## License

Apache-2.0（与上游一致）
