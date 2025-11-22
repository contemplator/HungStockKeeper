# HungStockKeeper

## 專案簡介

這是一個用來追蹤個人台灣股票投資的網站，目標是方便記錄與追蹤自己的投資狀況。

### 初期主要功能

1. **持股紀錄**
   - 記錄已購買的股票（股票代號/名稱）
   - 買進日期
   - 買進價格
   - 目前預估損益
   - 購買券商

2. **追蹤清單**
   - 追蹤的股票目前股價
   - 追蹤時的股價
   - 追蹤原因
   - 是否如預期漲跌

### 未來規劃功能
- 計算便宜價
- 盈在表
- 三大財報表摘要

## 技術選型
- 前端：Angular
- 後端：Go
- 資料庫：PostgreSQL
- 定期排程：Go 伺服器負責定期抓取最新股票價格並更新資料庫
- 單一 Git Repo 管理前後端與資料庫 schema

## 專案架構建議

```
HungStockKeeper/
├── frontend/      # Angular 前端專案
├── backend/       # Go 後端 API 專案
└── README.md      # 專案說明文件
```

- `frontend/`：放置 Angular 專案原始碼，負責 UI 與與後端 API 溝通
- `backend/`：放置 Go 專案原始碼，負責 API、資料庫存取、排程任務
- `README.md`：專案說明文件

## 開始開發建議步驟
1. 初始化 Git repository，建立 frontend 與 backend 資料夾
2. 使用 Angular CLI 建立前端專案
3. 使用 Go module 初始化後端專案，設計 API 與資料庫 schema
4. 設定 PostgreSQL 資料庫，設計資料表
5. 實作基本 API 與前端畫面
6. 實作排程功能，定期更新股票價格

---

## 啟動方式說明

### 前端（Angular）

1. **Node.js 版本需求**：請使用 Node.js v22 以上
2. 安裝依賴：
   ```sh
   cd frontend
   npm install
   ```
3. 啟動開發伺服器：
   ```sh
   npm start
   # 或
   ng serve
   ```

### 後端（Go）

1. 安裝依賴：
   ```sh
   cd backend
   go mod tidy
   ```
2. 啟動伺服器：
   ```sh
   go run main.go
   ```

---

## 資料庫 Migration 管理（golang-migrate）

本專案使用 [`golang-migrate`](https://github.com/golang-migrate/migrate) 來管理 PostgreSQL schema 版本。

### migration 檔案命名規則

- 檔案放在 `db/migrations/` 目錄下
- 命名格式為：
  ```
  [順序]_[描述].up.sql
  [順序]_[描述].down.sql
  ```
  其中 `[順序]` 建議使用 `YYYYMMDDhhmm` 的時間戳記，例如：
  ```
  202511221518_init_schema.up.sql
  202511221518_init_schema.down.sql
  ```
  - `.up.sql`：升級（建立/修改資料表）
  - `.down.sql`：降級（還原/刪除資料表）

### 常用 migrate 指令

**升級（執行所有 up migration）**
```sh
migrate -path db/migrations -database "postgresql://[user]:[password]@[host]:[port]/[database]?sslmode=disable" up
```

**降級（回滾一個 migration）**
```sh
migrate -path db/migrations -database "postgresql://[user]:[password]@[host]:[port]/[database]?sslmode=disable" down 1
```

**指定步數升級/降級**
```sh
# 升級 2 步
migrate -path db/migrations -database "postgresql://[user]:[password]@[host]:[port]/[database]?sslmode=disable" up 2
# 降級 2 步
migrate -path db/migrations -database "postgresql://[user]:[password]@[host]:[port]/[database]?sslmode=disable" down 2
```

**更多用法請參考官方文件：** https://github.com/golang-migrate/migrate

---

歡迎一起貢獻與討論！