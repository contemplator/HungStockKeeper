-- 啟用 UUID 擴充功能 (可選，若未來想用 UUID 當主鍵)
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. 用戶資料表
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2. 券商資料表 (預設選單)
CREATE TABLE brokerages (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

-- 插入預設券商資料
INSERT INTO brokerages (name) VALUES 
('新光證券'), 
('永豐金大戶投'),
('元大證券'),
('凱基證券'),
('富邦證券'),
('國泰證券'),
('群益證券'),
('玉山證券'),
('兆豐證券'),
('台新證券');

-- 3. 持股紀錄資料表
CREATE TABLE holdings (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stock_code VARCHAR(20) NOT NULL, -- 股票代號
    stock_name VARCHAR(100),         -- 股票名稱
    buy_date DATE NOT NULL,          -- 買進日期
    buy_price DECIMAL(10, 2) NOT NULL, -- 買進價格
    quantity INTEGER NOT NULL,       -- 股數
    brokerage_id INTEGER REFERENCES brokerages(id) ON DELETE SET NULL, -- 購買券商
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 4. 追蹤清單資料表
CREATE TABLE watchlist (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stock_code VARCHAR(20) NOT NULL, -- 股票代號
    stock_name VARCHAR(100),         -- 股票名稱
    tracked_price DECIMAL(10, 2),    -- 追蹤時的股價
    reason TEXT,                     -- 追蹤原因
    is_trend_expected BOOLEAN,       -- 是否如預期漲跌 (True: 如預期, False: 不如預期, Null: 尚未評估)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 建立索引以加速查詢
CREATE INDEX idx_holdings_user_id ON holdings(user_id);
CREATE INDEX idx_watchlist_user_id ON watchlist(user_id);
