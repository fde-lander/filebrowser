# Deploy Guide: FileBrowser Quantum v1.4.0.2

**Date:** 2026-07-11
**Branch:** v1.4.0.2-image-viewer-compression (merged to main)
**Docker Image:** filebrowser-fde:v1.4.0.2 (84MB tar)

---

## 1. 部署文件

- **Docker Image tar：** /home/hermes/workspace/filebrowser-fde/filebrowser-fde-v1.4.0.2.tar
- **Image name：** filebrowser-fde:v1.4.0.2
- **Image size：** 84MB（tar），316MB（解压后磁盘）

---

## 2. 部署步骤

### 步骤 1：复制 tar 到目标服务器

在本地执行：

scp /home/hermes/workspace/filebrowser-fde/filebrowser-fde-v1.4.0.2.tar user@target-server:/tmp/

### 步骤 2：在目标服务器加载 Docker image

ssh 到目标服务器后执行：

docker load -i /tmp/filebrowser-fde-v1.4.0.2.tar

验证加载成功：

docker images filebrowser-fde:v1.4.0.2

### 步骤 3：更新 docker-compose.yaml

编辑目标服务器上的 docker-compose.yaml，将 image 改为：

image: filebrowser-fde:v1.4.0.2

### 步骤 4：重启容器

docker-compose down && docker-compose up -d

### 步骤 5：验证启动

检查容器状态：

docker ps | grep filebrowser

检查日志：

docker logs filebrowser 2>&1 | tail -20

访问 Web UI 确认正常加载。

---

## 3. 测试清单（共 23 项）

### 功能 A：图片查看器增强（10 项）

| # | 测试项 | 操作 | 预期结果 |
|---|--------|------|----------|
| T1 | 翻页无黑屏 | 打开图片 -> 点击下一页 | 无黑屏闪烁，平滑过渡 |
| T2 | 上一页缓存命中 | 看下一页后点上一页 | 即时显示，无需重新加载 |
| T3 | 来回翻页 | 快速 B->C->B->C | 全部即时，零等待 |
| T4 | 点击右侧翻页 | 点击图片右 40% 区域 | 下一张图片 |
| T5 | 点击左侧翻页 | 点击图片左 40% 区域 | 上一张图片 |
| T6 | 双击缩放不冲突 | 快速双击图片 | 缩放，不触发翻页 |
| T7 | 拖拽不翻页 | 放大后拖拽图片 | 平移，不翻页 |
| T8 | 移动端按钮持久 | Settings 开启后手机访问 | 按钮不消失 |
| T9 | 透明度调节 | Settings 拖动滑块 | 按钮实时变透明 |
| T10 | 设置持久化 | 改设置后刷新页面 | 设置保持不变 |

### 功能 B：图片压缩（13 项）

| # | 测试项 | 操作 | 预期结果 |
|---|--------|------|----------|
| T11 | 单文件右键 | 右键一张 JPEG | 出现「壓縮圖片」 |
| T12 | 文件夹右键 | 右键一个文件夹 | 出现「壓縮圖片」 |
| T13 | 非图片隐藏 | 右键非图片文件 | 无「壓縮圖片」选项 |
| T14 | 自动全选 | 打开压缩弹窗 | 所有图片已勾选 |
| T15 | 档位切换 | 切换低/中/高档 | Quality 滑块范围变化 |
| T16 | 实时预览 | 调节 quality 滑块 | 预览图片实时更新 |
| T17 | 放大预览 | 点击预览图 | 放大查看细节 |
| T18 | ZSTD 备份 | 确认压缩（备份开启） | 当前目录生成 .tar.zst |
| T19 | 进度显示 | 压缩进行中 | 显示进度条+文件名 |
| T20 | 格式转换 | 检查压缩后文件 | .jpg 变成 .webp |
| T21 | 保底逻辑 | 压缩一张已很小图片 | 原图保留不变 |
| T22 | GIF 跳过 | 压缩含 GIF 的文件夹 | GIF 保留原格式 |
| T23 | 设置记忆 | 关闭弹窗再打开 | 记住上次档位/quality/备份 |

---

## 4. 回滚方案

如果测试失败需要回滚到 v1.4.0.1：

docker tag filebrowser-fde:custom filebrowser-fde:v1.4.0.1
docker-compose down && docker-compose up -d

（假设旧 image 仍在目标服务器上）

---

## 5. 新增依赖说明

### Go 依赖（已编译到 binary，无需额外安装）
- github.com/deepteams/webp v1.2.7（WebP 编码）
- github.com/klauspost/compress/zstd（ZSTD 备份压缩）
- github.com/esimov/colorquant（PNG 调色板优化）

### Docker 系统包（已安装到 image 内）
- oxipng（PNG 无损优化，低档 PNG 用）
- pngquant（PNG 调色板减少，低档 PNG 用）

---

## 6. 新增 API Endpoints

- POST /api/compress-images/preview — 预览压缩效果（返回 blob）
- POST /api/compress-images — 执行批量压缩（返回 taskId）
- GET /api/compress-images/progress?taskId=xxx — SSE 进度推送

---

## 7. 新增用户设置

Settings -> Profile 页面新增「Image Viewer」区域：
- 圖片預載（默认 ON）
- 漸變過渡效果（默认 Crossfade）
- 點擊翻頁（默认 ON）
- 移動端按鈕常駐（默认 OFF）
- 按鈕透明度（默认 100%）

压缩弹窗内记忆：
- 上次档位（默认 Medium）
- 上次 Quality 值
- ZSTD 备份开关（默认 ON）
