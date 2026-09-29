# 部署文檔 — FileBrowser Quantum（FDE 自訂增強版）

本文件是 FDE 自訂版的完整部署指引。請按順序執行，不要跳步。

---

## ⚠️ 部署模式（先讀這一節）

本項目採用 **「本地建置 + tar.zst 傳輸 + 遠端載入」** 模式，流程如下：

1. 在建置機器上執行 `docker build` 產生映像
2. 用 `docker save` + `zstd` 壓縮為 `.tar.zst` 檔
3. 將 `.tar.zst` 傳輸到目標伺服器
4. 在目標伺服器上執行 `docker load` 載入映像
5. 用 `docker compose` 啟動服務

**重要：目標伺服器不需要建置環境，不需要 Git、Go、Node.js，只需要 Docker Engine。**

---

## 版本資訊

• 映像版本：**v1.4.0.5-fde-hotfix**
• Fork 倉庫：https://github.com/fde-lander/filebrowser.git
• 基準分支：main（tag v1.4.0.5-fde-hotfix）
• 已作廢版本：v1.4.0.6、v1.4.0.7、v1.5.6.1-fde（請勿使用，詳見 README）

---

## 第 0 步：備份資料目錄（強制，任何更新／升級前先做）

⚠️ **在執行任何「載入新映像、更新、升級」之前，第一步必須先備份資料目錄。**

• `cd /opt/filebrowser-fde`
• `tar -czf "backup-data-$(date +%Y%m%d-%H%M%S).tar.gz" ./data`

備份完成後，用 `ls -lh backup-data-*.tar.gz` 確認檔案非空（大小不應是 0）。

**`./data` 目錄內含：資料庫、設定檔、使用者資料、已置頂項目、縮圖快取。** 它是唯一的復原點——遺漏備份，回滾時將無法還原資料。

- 全新部署（尚無 `./data` 目錄）可跳過本步。但**已有資料的升級，本步不可省略**。

---

## 前置需求

### 目標伺服器需求

• Docker Engine 20.10 以上
• Docker Compose v2 以上（指令形式為 `docker compose`，不是舊版 `docker-compose`）
• 磁碟空間：映像約 87 MB，加上資料目錄至少預留 1 GB
• 記憶體：執行階段很輕，512 MB 以上即可穩定運作

### 需要開放的連接埠

• 容器內部固定監聽 **80**（HTTP）
• 對外連接埠由你在 `docker-compose.yaml` 中自行對應，本文件範例用 8080

### 驗證前置工具

在目標伺服器上執行以下指令，確認全部有回應：

• `docker --version`
• `docker compose version`

---

## 步驟 1：傳輸映像檔到目標伺服器

將建置好的 `.tar.zst` 檔傳輸到目標伺服器。常用方式：

**方式 A：SCP 直接傳輸**

• `scp filebrowser-fde-v1.4.0.5-fde-hotfix.tar.zst user@目標伺服器:/tmp/`

**方式 B：先傳到中轉位置再下載**

若建置機器與目標伺服器不在同網段，可先上傳到雲端儲存或中轉伺服器，再從目標伺服器下載。

### 驗證檔案完整性（推薦）

傳輸完成後，在目標伺服器上驗證 SHA256：

• `sha256sum /tmp/filebrowser-fde-v1.4.0.5-fde-hotfix.tar.zst`

將輸出與建置機器提供的 SHA256 比對，一致才繼續。

---

## 步驟 2：載入映像

在目標伺服器上執行：

• `docker load -i /tmp/filebrowser-fde-v1.4.0.5-fde-hotfix.tar.zst`

預期輸出類似：`Loaded image: filebrowser-fde:v1.4.0.5-fde-hotfix`

### 驗證映像已載入

• `docker images | grep filebrowser-fde`

預期看到 `filebrowser-fde:v1.4.0.5-fde-hotfix`，大小約 87 MB。

### 驗證映像內建版本

• `docker run --rm --entrypoint ./filebrowser filebrowser-fde:v1.4.0.5-fde-hotfix version`

預期輸出含 `v1.4.0.5-fde-hotfix`。

---

## 步驟 3：建立資料目錄與設定檔

FileBrowser Quantum **必須有 config.yaml**，沒有會直接啟動失敗。這是最常見的失敗原因。

### 建立目錄

• `mkdir -p /opt/filebrowser-fde/data`

### 建立 data/config.yaml

新增檔案 `data/config.yaml`，最小可用內容如下：

server:
  sources:
    - path: "/srv"
      config:
        defaultEnabled: true

auth:
  adminUsername: admin
  adminPassword: admin

**結構要點（極重要）：**

• `sources` **必須縮排在 `server:` 底下**。若寫在檔案最上層，會報 `Field validation for 'Sources' failed` 並啟動失敗
• `path: "/srv"` 是**容器內**路徑，對應 compose 的 volume 掛載目標
• 縮排必須用空格，**不可用 Tab**。混用會導致 YAML 解析失敗

### 可選：加入快取目錄

若要啟用縮圖快取（建議），完整內容改為：

server:
  cacheDir: /home/filebrowser/data/tmp
  sources:
    - path: "/srv"
      config:
        defaultEnabled: true

auth:
  adminUsername: admin
  adminPassword: admin

---

## 步驟 4：建立 docker-compose.yaml

在 `/opt/filebrowser-fde/` 新增 `docker-compose.yaml`：

services:
  filebrowser:
    image: filebrowser-fde:v1.4.0.5-fde-hotfix
    container_name: filebrowser-fde
    environment:
      - TZ=Asia/Taipei
    volumes:
      - /你的/實際/檔案/目錄:/srv
      - ./data:/home/filebrowser/data
    ports:
      - 8080:80
    restart: unless-stopped

**各項說明：**

• `image` 必須與步驟 2 載入的標籤完全一致，否則 compose 會嘗試去 Docker Hub 拉取而失敗
• **不要寫 `build:` 區塊**，映像已由步驟 2 載入，compose 直接使用
• 第一個 volume 的冒號左邊是你的實際檔案目錄（主機路徑），右邊固定為 `/srv`（對應 config.yaml 的 path）
• 第二個 volume 掛載資料目錄，內含 config.yaml、database.db、快取
• `8080:80` 左邊可改成你要的對外連接埠，右邊固定 80
• `container_name` 可自訂，但不可與現有容器重複

---

## 步驟 5：啟動與驗證

### 啟動

• `cd /opt/filebrowser-fde`
• `docker compose up -d`

### 確認容器狀態

• `docker compose ps`

預期狀態為 `Up`（不是 `Restarting` 或 `Exited`）。

### 查看日誌

• `docker compose logs --tail=50`

若啟動正常，會看到服務監聽訊息，不會有 FATAL 錯誤。

### 健康檢查

• `docker inspect --format='{{.State.Health.Status}}' filebrowser-fde`

預期輸出 `healthy`。首次啟動需等待約 30 秒讓健康檢查通過。

### 瀏覽器驗證

• 開啟 `http://你的伺服器IP:8080`
• 用 `admin` / `admin` 登入
• **登入後第一件事：立即修改 admin 密碼**

### 功能驗證清單

登入後依序確認：

1. 左側能看到你掛載的檔案目錄
2. 點開一張圖片，圖片檢視器正常顯示，左右翻頁無黑屏
3. 右鍵圖片或資料夾，選單出現「壓縮圖片」（需 Admin 權限）
4. 右鍵壓縮檔，選單出現「解壓到新資料夾」（需 Admin 權限）
5. 手機豎屏開啟壓縮視窗，預覽區塊排列正常（Bug I 修正項目）

---

## 步驟 6：日後更新版本

**⚠️ 更新前必須先完成「第 0 步：備份資料目錄」。** 備份做完才做下列更新步驟。

接著更新版本：

• 將新版本的 `.tar.zst` 傳輸到目標伺服器
• `docker load -i /tmp/filebrowser-fde-<新版本>.tar.zst`
• 編輯 `docker-compose.yaml`，把 `image` 改成新版本標籤
• `docker compose down && docker compose up -d`

**建議保留舊映像以便回滾**，不要載入新版後就刪除舊版。

---

## 回滾

**⚠️ 回滾原則：先還原資料，再退版本。**

資料目錄若已被新版本寫入而損壞（或格式不兼容），單獨退回舊映像是不夠的，必須一併還原備份資料。

### 情境 A：新版本有問題，回到上一個映像（資料無損）

• 編輯 `docker-compose.yaml`，把 `image` 改回舊版本標籤
• `docker compose down && docker compose up -d`

### 情境 A-2：資料也被改壞，回滾映像＋還原資料

• 停服務：`docker compose down`
• 還原資料：`tar -xzf backup-data-<時間戳>.tar.gz -C /opt/filebrowser-fde`（覆蓋 `./data`）
• 把 `image` 改回舊版本標籤
• `docker compose up -d`

### 情境 B：完全放棄，回到官方映像

• 編輯 `docker-compose.yaml`，把 `image` 改為 `gtstef/filebrowser:1.4.0-stable`
• `docker compose down && docker compose up -d`

**資料目錄 `./data` 不要刪除**，內含你的資料庫與設定。回滾時優先從 `backup-data-*.tar.gz` 還原。

---

## 疑難排解

**症狀：`docker compose up` 時嘗試從 Docker Hub 拉取映像而失敗**

• 原因：compose 中的 `image` 名稱與實際載入的標籤不一致
• 解法：`docker images | grep filebrowser` 確認實際標籤，改成一致

**症狀：日誌出現 `FATAL: config file does not exist`**

• 原因：沒有建立 `data/config.yaml`
• 解法：回到步驟 3 建立該檔案

**症狀：`Field validation for 'Sources' failed`**

• 原因：`sources` 寫在 YAML 最上層
• 解法：`sources` 必須縮排在 `server:` 底下

**症狀：容器不斷重啟，狀態為 `Restarting`**

• 原因：config.yaml 格式錯誤（最常見是混用 Tab 與空格）
• 解法：檢查縮排，全部改用空格

**症狀：首次登入彈出「A new database was created」**

• 原因：全新安裝的正常現象
• 解法：直接關閉提示，修改密碼後即為正常狀態

**症狀：介面顯示的目錄是空的**

• 原因：volume 掛載來源路徑錯誤或不存在
• 解法：確認 compose 第一個 volume 的冒號左邊是實際存在的主機目錄；確認 config.yaml 的 `path` 是 `/srv`

**症狀：`container_name` 衝突**

• 原因：舊容器未清理
• 解法：`docker compose down` 後再啟動；必要時 `docker rm 容器名`

**症狀：權限錯誤，無法寫入資料庫**

• 原因：掛載目錄權限不符（容器內以 uid 1000 執行）
• 解法：確認 `./data` 目錄對 uid 1000 可寫

---

## 附錄：關鍵路徑與連接埠對照

• 映像內應用程式路徑：`/home/filebrowser/filebrowser`
• 映像內設定檔路徑：`/home/filebrowser/data/config.yaml`
• 映像內資料庫路徑：`/home/filebrowser/data/database.db`
• 映像內快取路徑（若已設定）：`/home/filebrowser/data/tmp`
• 容器內部監聽連接埠：`80`
• 映像內執行使用者：`filebrowser`（uid 1000）
• 健康檢查端點：`http://localhost:80/health`

---

## 附錄：建置機器資訊（僅供參考，目標伺服器不需要）

以下為建置時的環境資訊，記錄用途，目標伺服器不需要任何這些工具：

• 建置機器：Linux 6.12.73+deb13-cloud-amd64, 2 核 CPU, 3.8G RAM + 7.8G swap
• Docker Engine 29.3.0 + Buildx v0.31.1
• zstd v1.5.7
• 建置指令：docker build --build-arg="VERSION=v1.4.0.5-fde-hotfix" --build-arg="REVISION=$(git rev-parse --short HEAD)" -t filebrowser-fde:v1.4.0.5-fde-hotfix -f _docker/Dockerfile .
• 打包指令：docker save filebrowser-fde:v1.4.0.5-fde-hotfix | zstd -5 -o filebrowser-fde-v1.4.0.5-fde-hotfix.tar.zst

---

*本文件對應版本：v1.4.0.5-fde-hotfix*
