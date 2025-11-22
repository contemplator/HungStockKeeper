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

歡迎一起貢獻與討論！