# 部署文檔 — FileBrowser Quantum（FDE 自訂增強版）

本文件是 FDE 自訂版的完整部署指引。請按順序執行，不要跳步。

---

## ⚠️ 部署模式（先讀這一節）

本項目採用 **「目標伺服器直接建置」** 模式，流程如下：

1. 在目標伺服器上 `git clone` 或 `git pull` 本倉庫
2. 在目標伺服器上執行 `docker build` 產生映像
3. 用 `docker compose` 直接使用該映像啟動服務

**重要：本倉庫不提供、也不包含任何預先編譯的 tar 映像檔。**

- 倉庫內所有 `*.tar` 已列入 `.gitignore`，不會被追蹤、不會被推送
- 因此 `git clone` 之後目錄內不會有任何 tar 檔，這是預期行為，不是缺失
- 不需要另外尋找或下載 tar 檔，直接在伺服器上建置即可

---

## 版本資訊

• 基準版本：**v1.4.0.5-fde-hotfix**（本 Fork 唯一採用並建議使用的版本）
• Fork 倉庫：https://github.com/fde-lander/filebrowser.git
• 預設分支：main
• 已作廢版本：v1.4.0.6、v1.4.0.7（請勿使用，詳見 README）

---

## 前置需求

### 伺服器需求

• Docker Engine 20.10 以上
• Docker Compose v2 以上（指令形式為 `docker compose`，不是舊版 `docker-compose`）
• Git
• 磁碟空間：建置過程需要約 6 至 10 GB 可用空間（含建置快取），完成後映像約 85 MB
• 記憶體：建議 2 GB 以上（低記憶體伺服器請先讀「資源需求」章節）

### 需要開放的連接埠

• 容器內部固定監聽 **80**（HTTP）
• 對外連接埠由你在 `docker-compose.yaml` 中自行對應，本文件範例用 8080

### 驗證前置工具

在伺服器上執行以下指令，確認全部有回應：

• `docker --version`
• `docker compose version`
• `git --version`

---

## 步驟 1：取得源碼

選擇一個工作目錄（範例用 `/opt/filebrowser-fde`）：

建立工作目錄並進入：

• `mkdir -p /opt/filebrowser-fde`
• `cd /opt/filebrowser-fde`

複製本倉庫：

• `git clone https://github.com/fde-lander/filebrowser.git .`

確認目前位於正確版本（v1.4.0.5-fde-hotfix 內容）：

• `git log -1 --oneline`
• `git branch --show-current`

預期：分支為 `main`，最新提交為 v1.4.0.5-fde-hotfix 系列內容。

---

## 步驟 2：建置映像

### 確認建置檔存在

在建置前先確認 Dockerfile 與必要輸入檔存在：

• `ls -l _docker/Dockerfile`
• `ls -l backend/go.mod backend/go.sum backend/config.yaml`
• `ls -l frontend/package.json frontend/package-lock.json`

以上檔案全部必須存在，缺任何一個都會導致建置失敗。

### 完整版建置指令（建議使用）

完整版包含 FFmpeg、exiftool、oxipng、pngquant，支援影片與文件預覽，以及全部圖片壓縮功能。

在倉庫根目錄執行：

• `docker build --build-arg="VERSION=v1.4.0.5-fde-hotfix" --build-arg="REVISION=$(git rev-parse --short HEAD)" -t filebrowser-fde:v1.4.0.5-fde-hotfix -f _docker/Dockerfile .`

指令說明：

• `--build-arg="VERSION=..."` 寫入映像內建版本字串，顯示在介面與 `filebrowser version` 輸出
• `--build-arg="REVISION=..."` 寫入提交雜湊，方便日後核對來源
• `-t filebrowser-fde:v1.4.0.5-fde-hotfix` 指定映像名稱與標籤，之後 compose 會用到同一個名字
• `-f _docker/Dockerfile` 指定建置檔路徑（本項目的 Dockerfile 放在 `_docker/` 子目錄，不是根目錄）
• 最後的 `.` 是建置上下文，必須是倉庫根目錄

**注意：最後的 `.` 必須是倉庫根目錄。** 因為 Dockerfile 內用 `COPY ./backend` 與 `COPY ./frontend/`，上下文指錯會直接複製失敗。

### 建置過程會做什麼

Dockerfile 為三階段建置：

1. 第一階段：拉取 FFmpeg 映像，取得 ffmpeg 與 ffprobe 執行檔
2. 第二階段（Go 後端）：Alpine + Go 1.25，安裝 gcc 與 musl-dev，以 CGO 啟用方式編譯後端
3. 第三階段（前端）：Node 20 以上，安裝 npm 依賴後執行 `vite build`
4. 最終階段：Alpine，合併前兩階段產物，安裝執行期套件

### 建置時間預估

• 4 核心以上、8 GB 記憶體：約 8 至 15 分鐘
• 1 核心、2 GB 記憶體：可能長達 40 至 90 分鐘，且有可能因記憶體不足而失敗（詳見「資源需求」章節）

### 建置成功判定

建置結束時應看到成功訊息。之後驗證映像：

• `docker images | grep filebrowser-fde`

預期看到 `filebrowser-fde:v1.4.0.5-fde-hotfix`，大小約 85 MB。

進一步驗證映像內建版本：

• `docker run --rm --entrypoint ./filebrowser filebrowser-fde:v1.4.0.5-fde-hotfix version`

預期輸出含 `v1.4.0.5-fde-hotfix`。

---

## 步驟 3：建立資料目錄與設定檔

FileBrowser Quantum **必須有 config.yaml**，沒有會直接啟動失敗。這是最常見的失敗原因。

### 建立目錄

• `mkdir -p /opt/filebrowser-fde/data`

### 建立 data/config.yaml

新增檔案 `data/config.yaml`，最小可用內容如下：

```
server:
  sources:
    - path: "/srv"
      config:
        defaultEnabled: true

auth:
  adminUsername: admin
  adminPassword: admin
```

**結構要點（極重要）：**

• `sources` **必須縮排在 `server:` 底下**。若寫在檔案最上層，會報 `Field validation for 'Sources' failed` 並啟動失敗
• `path: "/srv"` 是**容器內**路徑，對應 compose 的 volume 掛載目標
• 縮排必須用空格，**不可用 Tab**。混用會導致 YAML 解析失敗

### 可選：加入快取目錄

若要啟用縮圖快取（建議），完整內容改為：

```
server:
  cacheDir: /home/filebrowser/data/tmp
  sources:
    - path: "/srv"
      config:
        defaultEnabled: true

auth:
  adminUsername: admin
  adminPassword: admin
```

---

## 步驟 4：建立 docker-compose.yaml

在 `/opt/filebrowser-fde/` 新增 `docker-compose.yaml`：

```
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
```

**各項說明：**

• `image` 必須與步驟 2 的 `-t` 標籤完全一致，否則 compose 會嘗試去 Docker Hub 拉取而失敗
• **不要寫 `build:` 區塊**，本文件模式是先用步驟 2 建好映像，再由 compose 直接使用
• 第一個 volume 的冒號左邊是你的實際檔案目錄（主機路徑），右邊固定為 `/srv`（對應 config.yaml 的 path）
• 第二個 volume 掛載資料目錄，內含 config.yaml、database.db、快取
• `8080:80` 左邊可改成你要的對外連接埠，右邊固定 80
• `container_name` 可自訂，但不可與現有容器重複

---

## 步驟 5：啟動與驗證

### 啟動

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

當有新版本時：

• `cd /opt/filebrowser-fde`
• `git pull`
• 重新執行步驟 2 的建置指令，只把 `-t` 標籤改成新版本號
• 編輯 `docker-compose.yaml`，把 `image` 改成新標籤
• `docker compose down && docker compose up -d`

**建議保留舊映像以便回滾**，不要建置完就刪除。

---

## 回滾

### 情境 A：新版本有問題，回到上一個映像

• 編輯 `docker-compose.yaml`，把 `image` 改回舊版本標籤
• `docker compose down && docker compose up -d`

### 情境 B：原始碼也要退回

• `cd /opt/filebrowser-fde`
• `git log --oneline` 找到目標提交
• `git checkout <目標提交雜湊>`
• 重新建置並啟動

### 情境 C：完全放棄，回到官方映像

• 編輯 `docker-compose.yaml`，把 `image` 改為 `gtstef/filebrowser:1.4.0-stable`
• `docker compose down && docker compose up -d`

資料目錄 `./data` 不要刪除，內含你的資料庫與設定。

---

## 資源需求與低記憶體伺服器注意事項

### 建置階段是資源瓶頸，不是執行階段

• **執行階段**很輕：單核心、數百 MB 記憶體即可穩定運作
• **建置階段**很重：Go 以 CGO 編譯加上 npm 打包，尖峰記憶體需求可達 1.5 至 2 GB 以上

### 在 1 核心 2 GB 伺服器上建置的風險

可能出現的情況：

• 建置極慢（數十分鐘至一小時以上）
• 在建置後段被系統 OOM Killer 終止，建置失敗

### 建議做法（擇一）

**做法一：先加 swap 再建置（最簡單）**

• `sudo fallocate -l 4G /swapfile`
• `sudo chmod 600 /swapfile`
• `sudo mkswap /swapfile`
• `sudo swapon /swapfile`
• 驗證：`free -h`

建置完成後可保留 swap 作為保險，系統已有 ZRAM 時亦可並存。

**做法二：在較強機器上建置，再傳映像到目標伺服器**

在較強機器上建置：

• `docker build --build-arg="VERSION=v1.4.0.5-fde-hotfix" --build-arg="REVISION=$(git rev-parse --short HEAD)" -t filebrowser-fde:v1.4.0.5-fde-hotfix -f _docker/Dockerfile .`
• `docker save filebrowser-fde:v1.4.0.5-fde-hotfix -o filebrowser-fde-v1.4.0.5-fde-hotfix.tar`
• `scp filebrowser-fde-v1.4.0.5-fde-hotfix.tar user@目標伺服器:/tmp/`

在目標伺服器上載入：

• `docker load -i /tmp/filebrowser-fde-v1.4.0.5-fde-hotfix.tar`

之後步驟 3 至 5 完全相同。此做法只是把「建置」移到別處，**最終仍是 compose 直接使用該映像**。

**做法三：改用精簡版建置（節省時間與記憶體）**

精簡版不含 FFmpeg、exiftool、oxipng、pngquant。適用於不需要影片預覽的環境。

• `docker build --build-arg="VERSION=v1.4.0.5-fde-hotfix" --build-arg="REVISION=$(git rev-parse --short HEAD)" -t filebrowser-fde:v1.4.0.5-fde-hotfix -f _docker/Dockerfile.slim .`

**注意：精簡版不包含 oxipng 與 pngquant，因此圖片壓縮功能的 PNG 低檔無損路徑會不可用。** 若你需要完整的圖片壓縮功能，請使用完整版。

### 建置完成後清理

建置快取會佔用數 GB，確認新映像運作正常後可清理：

• `docker builder prune -f`

**清理前請先確認所有需要的映像都已存在**，並保留舊版本映像以便回滾。

---

## 疑難排解

**症狀：建置時 `COPY ./backend: not found` 或類似錯誤**

• 原因：建置上下文不是倉庫根目錄
• 解法：確認指令最後是 `-f _docker/Dockerfile .`，且當前目錄為倉庫根目錄（`ls backend` 應該有反應）

**症狀：`docker compose up` 時嘗試從 Docker Hub 拉取映像而失敗**

• 原因：compose 中的 `image` 名稱與實際建置的標籤不一致
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

**症狀：建置中途被終止，日誌出現 `Killed`**

• 原因：記憶體不足被 OOM Killer 終止
• 解法：參考「資源需求」章節加 swap，或改用做法二、做法三

**症狀：權限錯誤，無法寫入資料庫**

• 原因：掛載目錄權限不符（容器內以 uid 1000 執行）
• 解法：確認 `./data` 目錄對 uid 1000 可寫

---

## 附錄：映像變體對照

**完整版 — `_docker/Dockerfile`**

• 內容：Alpine + ffmpeg + ffprobe + exiftool + oxipng + pngquant
• 適用：需要影片預覽、文件預覽、完整圖片壓縮功能
• 建議：一般用途選這個

**精簡版 — `_docker/Dockerfile.slim`**

• 內容：Alpine，不含 ffmpeg / exiftool / oxipng / pngquant
• 適用：只需基本檔案瀏覽與圖片檢視
• 限制：PNG 無損壓縮路徑不可用

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

## 附錄：倉庫結構速查

• `backend/` — Go 後端原始碼（`httpRouter.go` 為路由註冊處，`http/compress.go` 為壓縮功能）
• `frontend/` — Vue 3 前端原始碼（`src/components/prompts/` 為各對話框元件）
• `_docker/` — Dockerfile 與 compose 範例
• `docs/` — 設計文件、交接文件、歷史計劃
• `README.md` — 功能清單與版本歷史
• `DEPLOY.md` — 本文件

---

*本文件對應版本：v1.4.0.5-fde-hotfix*
