<div align="center">

  [![Apache-2.0 License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)

  <h3>FileBrowser Quantum（FDE 自訂增強版）</h3>
  基於 FileBrowser Quantum v1.4.0-stable<br/>
  最新穩定版：v1.4.0.5-fde-hotfix<br/><br/>
</div>

## FDE FORK 說明

本倉庫是 [FileBrowser Quantum](https://github.com/gtsteffaniak/filebrowser) 的自訂增強版本，基於 v1.4.0-stable tag 建立，維護者為 **FDE**。

**Fork 倉庫**：https://github.com/fde-lander/filebrowser.git

**上游項目**：https://github.com/gtsteffaniak/filebrowser （原版 Quantum）

**官方文檔**：https://filebrowserquantum.com/en/docs/getting-started/docker

本 Fork 不向上游提交 PR，所有改動僅供自用。改動以最小侵入方式實現，保持與上游代碼結構相容。

**基準版本**：v1.4.0.5-fde-hotfix — 本 Fork 目前唯一採用並建議使用的版本。

---

## ⚠️ 重要提醒：部署模式已變更

本項目自即日起採用 **「目標伺服器拉取倉庫 → Docker Build → Docker Compose」** 模式。

• 在目標伺服器上 `git clone` 本倉庫
• 在目標伺服器上執行 `docker build` 產生映像
• 用 `docker compose` 直接使用該映像啟動服務

**本倉庫不再提供預先編譯的 tar 映像檔。** 倉庫內所有 `*.tar` 已列入 `.gitignore`，不會被追蹤、不會被推送。`git clone` 之後目錄內沒有 tar 檔是預期行為。

**完整詳細的部署指引（含前置需求、建置指令、設定檔、疑難排解、低記憶體伺服器注意事項、回滾流程）請務必閱讀：**

📄 **[DEPLOY.md](DEPLOY.md)**

部署前請先完整讀過該文件再動手，尤其是「資源需求與低記憶體伺服器注意事項」章節。

---

## 歷史版本說明（日期倒序）

| 日期 | 版本 | 內容摘要 | 狀態 |
| --- | --- | --- | --- |
| 2026-07-16 | v1.4.0.5-fde-hotfix | 移動端壓縮預覽佈局修正（Bug I）；首個 GitHub Release 發佈 | 採用（最新穩定版） |
| 2026-07-16 | v1.4.0.6 | 行動端 Bug G/H/B 修復 ＋ 壓縮控制功能（暫停／繼續／取消／跳過）等 | 已作廢（部署後失效，全數回滾） |
| 2026-07-16 | v1.4.0.7 | 原定為 v1.4.0.6 的後續修復計劃 | 已作廢（未實施，直接跳過） |
| 2026-07-14 | v1.4.0.5 | 壓縮系統改為佇列式輪詢、備份開關跨 session 持久化、資料夾展開、PNG 低檔改 WebP Q75、預覽原圖清晰度修正、進度計數修正、資料夾壓縮修正、行動端豎屏相容 | 採用（已併入最新版） |
| 2026-07-13 | v1.4.0.4-hotfix | 備份路徑三級退避、預覽目錄展開、輪詢欄位對齊、預覽 UI 重做、decode-first 過渡架構 | 採用（已併入） |
| 2026-07-12 | v1.4.0.3 | 過渡架構重寫、壓縮 API 前後端對齊、資料夾右鍵偵測、Admin 權限閘門 | 採用（已併入） |
| 2026-07-11 | v1.4.0.2 | 圖片檢視器增強（雙緩衝／快取池／漸變過渡／點擊翻頁）＋ 圖片壓縮（三檔位／自動備份） | 採用（已併入） |
| 2026-07-11 | v1.4.0.1 | 解壓到新資料夾 | 採用（已併入） |

---

## 快速部署

**本項目的正式部署方式為「目標伺服器拉取倉庫 → Docker Build → Docker Compose」。**

完整步驟請閱讀 **[DEPLOY.md](DEPLOY.md)**。以下僅為速覽。

步驟 1：在目標伺服器取得源碼

• `git clone https://github.com/fde-lander/filebrowser.git /opt/filebrowser-fde`
• `cd /opt/filebrowser-fde`

步驟 2：在目標伺服器建置映像

• `docker build --build-arg="VERSION=v1.4.0.5-fde-hotfix" --build-arg="REVISION=$(git rev-parse --short HEAD)" -t filebrowser-fde:v1.4.0.5-fde-hotfix -f _docker/Dockerfile .`

步驟 3：建立 `data/config.yaml`（必須，否則無法啟動）

• 內容見 DEPLOY.md 步驟 3

步驟 4：建立 `docker-compose.yaml`，`image` 設為 `filebrowser-fde:v1.4.0.5-fde-hotfix`

步驟 5：啟動

• `docker compose up -d`

**注意：本倉庫不包含預先編譯的 tar 映像檔。** 所有 `*.tar` 已列入 `.gitignore`。若你偏好「先建置再傳映像」的做法，請參考 DEPLOY.md 的「資源需求與低記憶體伺服器注意事項」章節（做法二）。

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

詳細配置請參考 [官方文檔](https://filebrowserquantum.com/en/docs/getting-started/config)。

---

## 舊版 Release 附件說明（歷史記錄）

以下說明過去以 tar 檔分發時期的情況，供追溯用途。

先前已發佈的 Release，其附件命名沿用了舊編號，實際情況如下：

- **Release 標籤**：v1.4.0.5-fde-stable
- **附件檔名**：filebrowser-fde-v1.4.0.6.tar
- **附件內部映像標籤**：filebrowser-fde:v1.4.0.6

以上檔名與映像標籤均為舊編號殘留，與實際內容不符。經映像分層內容驗證確認：

- 附件內容 = v1.4.0.5-fde-hotfix ＋ Bug I 移動端預覽佈局修正
- 不含任何 v1.4.0.6 失效代碼（全域狀態列、壓縮控制路由、暫停逾時設定全部不存在）
- 附件 SHA256 與本地建置產物完全一致

**新部署模式已不再需要 tar 檔。** 若你要使用此歷史附件，載入後建議重新打上正確標籤：

• `docker load -i filebrowser-fde-v1.4.0.6.tar`
• `docker tag filebrowser-fde:v1.4.0.6 filebrowser-fde:v1.4.0.5-fde-hotfix`
• 之後在 `docker-compose.yaml` 使用 `filebrowser-fde:v1.4.0.5-fde-hotfix`

---

## 已作廢版本（請勿使用）

以下版本已作廢，不會再維護，亦不會出現在任何後續升級路徑中。

**v1.4.0.6 — 已作廢**

- 作廢日期：2026-07-16
- 作廢原因：部署實測後大部分功能失效。全域狀態列從未觸發、壓縮控制按鈕全部不可見、點擊翻頁過渡品質下降、導覽按鈕停止顯示。
- 處置方式：已於同日全數回滾，代碼庫回復至 v1.4.0.5-fde-hotfix。不修補、不發佈，後續版本升級將直接跳過此版本。

**v1.4.0.7 — 已作廢**

- 作廢日期：2026-07-16
- 作廢原因：原定為 v1.4.0.6 的後續修復版本。因 v1.4.0.6 整條路線放棄，未進入實施階段即終止。
- 處置方式：不實施、不發佈，後續版本升級將直接跳過此版本。

未來如需升級，將由 v1.4.0.5-fde-hotfix 之後另起新編號，不再沿用 v1.4.0.6／v1.4.0.7。

---

## 本 Fork 自訂功能清單

以下是本 Fork 相比原版 Quantum v1.4.0-stable 的全部自訂改動。

### v1.4.0.5-fde-hotfix - 移動端預覽佈局修正（最新穩定版）

**Bug I - 移動端壓縮預覽佈局**：修正手機豎屏下壓縮預覽 overlay 與項目區塊的排列問題（3 行 CSS）。

- 新增 @media (max-width: 768px) 響應式規則
- 預覽 overlay 改為垂直排列 + 可滾動
- 預覽項目限高，避免撐破視窗
- 與翻頁功能零重疊

### v1.4.0.5 - 壓縮系統改為佇列式架構

**後端**：

- 廢棄 SSE 進度串流（progressManager 已移除），改為佇列式背景壓縮
- 新增壓縮佇列管理器 compressQueueManager（enqueue / dequeue / updateProgress / finishCurrent / getStatus）
- 新增狀態端點：GET /compress-images/status（取代舊 GET /compress-images/progress）
- 新增背景 worker：compressWorker 逐項處理佇列，避免長連線佔用資源
- 備份路徑解析與三級退避（同級目錄 → 上一層 → source 根目錄）
- 目錄遞迴展開為圖片清單

**前端**：

- 以 pollStatus() 輪詢（3 秒間隔）取代 subscribeProgress() SSE 訂閱
- 備份開關跨 session 持久化（compressBackup 加入自動持久化允許清單），預設 ON
- 資料夾展開為內部圖片清單
- 預覽原圖清晰度修正 + 佈局溢出修正
- 進度計數修正（不再卡在 N-1/N）
- 資料夾壓縮修正（selectedFileList 改用展開後清單）
- PNG 低檔改用 WebP Q75
- 行動端豎屏相容
- 預覽模式 checkbox 移至檔案清單下方
- 舊 CSS 標註為 LEGACY（保留不刪除）

**過渡引擎**：維持 v1.4.0.4-hotfix 行為，不作改動。

### v1.4.0.4 - Hotfix（備份路徑 + 預覽 UI + 過渡 decode-first）

**Issue 1 - 備份路徑解析 + 3 級退避**：
- 後端 backupPath 透過 resolveCompressPath 解析到真實檔案系統路徑
- 3 級退避：同級目錄 → 上一層目錄 → source 根目錄
- 全部失敗則中止壓縮（backup-first 原則）
- finishEvent 新增 BackupFallback 欄位通知前端

**Issue 2 - compressPreviewHandler 目錄展開**：
- 右鍵資料夾預覽不再 500 錯誤
- os.Stat + filepath.Walk 取第一張圖片做預覽

**Issue 3 - 輪詢欄位對齊**：
- 事件名稱與欄位映射以後端為準
- progress 欄位映射：processed -> current, current -> currentFile
- complete 欄位映射：success/skipped/failed + backupPath + backupFallback
- error 欄位映射：data.error 相容 data.message

**Issue 4 - 預覽 UI 全新設計**：
- checkbox「開啟預覽模式」toggle
- 預覽模式 ON：點擊檔案打開 overlay（左原圖 + 右壓縮預覽 + loading + fade-in）
- 全螢幕 overlay：黑色背景 + 返回按鈕 + 點擊外部返回
- 原圖用 getPreviewURL + blob URL（壓縮圖）+ beforeUnmount cleanup
- 預覽模式 OFF：正常勾選／取消選擇

**Issue 5 - 過渡 decode-first 架構**：
- 新增 waitForDecode() helper：decode().then() + 3s 超時 + 舊瀏覽器 fallback
- navigateToImage 重寫：set src -> waitForDecode -> swapBuffers（統一路徑）
- swapBuffers 三種模式全部改為「decode 完成後同時漸變」-> 永不到黑
- fade_to_black 改名為 fade（500ms 柔和漸變，不再有黑屏）
- 保留 fade_to_black 向後相容
- CSS transform 衝突修正：translate -> translate3d（恢復 GPU 合成）

### v1.4.0.3 - Bug Fix（過渡架構重寫 + API 對齊 + 權限門控）

- 圖片過渡架構重寫：CSS 管理過渡動畫 + generation token 防 stale callback
- 壓縮 API 前後端對齊：URL / 欄位名 / 參數格式全部以後端為準
- 資料夾右鍵偵測：isDir || type === 'directory' 雙重偵測
- Admin 權限門控：解壓到新資料夾 + 壓縮圖片均需 Admin 權限

### v1.4.0.2 - 圖片檢視器增強 + 圖片壓縮

**功能 A - 圖片檢視器增強**：
- 雙緩衝架構：imgA/imgB 兩個 img 元素交叉顯示，組件永不被銷毀
- 3 張快取池：Set 結構快取上一張 + 當前 + 下一張，來回翻頁零等待
- 漸變過渡引擎：可擴展註冊表模式，支援 crossfade / fade (soft) / instant
- 點擊翻頁：左右各 40% 區域，200ms 雙擊偵測延遲，拖拽 >10px 取消翻頁
- 移動端按鈕持久化 + 獨立透明度滑桿
- 設定持久化：所有設定透過 NonAdminEditable struct 跨 session 記憶

**功能 B - 圖片壓縮**：
- 右鍵批量壓縮：支援資料夾／多選檔案
- 3 檔壓縮級別：低檔（WebP Q75）、中檔（Q65）、高檔（Q55）
- PNG 特殊路徑：低檔改用 WebP Q75
- ZSTD 自動備份：tar.zst 格式，壓縮前先備份（backup-first 設計）
- 3 級退避：同級目錄 -> 上一層 -> source 根目錄
- 保底邏輯：壓縮後如 ≥ 原圖大小則跳過
- Admin 權限門控：僅 Admin 使用者可見可操作
- i18n：en / zh-cn / zh-tw 三語言

### v1.4.0.1 - 解壓到新資料夾

**功能**：右鍵壓縮檔時新增「解壓到新資料夾」選項，自動建立以壓縮檔名命名的資料夾並解壓到其中。

- 後端：archive.go 新增 CreateSubfolder 欄位 + archiveFolderName() + resolveUniqueSubfolderName()
- 前端：新增 ExtractToFolder.vue 組件 + ContextMenu 註冊 + i18n（en/zh-cn/zh-tw）
- 衝突處理：同名資料夾自動加後綴 (1)、(2)... 最多 100 次
- deleteAfterArchive 複用：使用現有使用者偏好持久化機制
- 安全：複用 normalizeArchiveEntryName + safeExtractPath + symlink 檢查

---

## 技術棧

- 後端：Go 1.25（http handlers + SQLite index + WebP/ZSTD）
- 前端：Vue 3 / Vite（雙緩衝圖片檢視器 + 壓縮彈窗 + 輪詢）
- Docker：3 階段建置（Go backend -> Node frontend -> Alpine final）
- 部署：目標伺服器 git clone → docker build → docker compose（不推 Docker Hub，不分發 tar）

## 環境資訊

- 目標伺服器：1 vCPU core, 2GB RAM, 2GB ZRAM + 2GB SWAP
- Docker volume 映射：/root/qbb/downloads:/folder（容器內）
- 反向代理：Nginx Proxy Manager（Docker，跨伺服器）

## 上游項目致謝

本 Fork 基於以下優秀開源項目：
- [FileBrowser Quantum](https://github.com/gtsteffaniak/filebrowser) by gtsteffaniak
- [原版 FileBrowser](https://github.com/filebrowser/filebrowser)

## License

Apache-2.0（與上游一致）
