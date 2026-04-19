# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# 啟動服務
go run main.go

# 編譯
go build -o shortenurl main.go

# 啟動 PostgreSQL（使用 Docker Compose）
docker-compose -f docker-compose-postgres.yaml up -d

# 啟動 Redis
redis-server

# 下載依賴
go mod tidy
```

## Docker

```bash
# Build image
docker build -t dayi1225/shortenurl:latest .

# Push to Docker Hub
docker push dayi1225/shortenurl:latest

# 本機一鍵啟動（含 PostgreSQL + Redis）
docker-compose -f docker-compose-postgres.yaml up -d
```

`Dockerfile` 採三階段 build：
1. **frontend-builder**（node:20-alpine）— `npm run build` 輸出至 `../static`
2. **go-builder**（golang:1.24-alpine）— 靜態編譯 Go binary
3. **runtime**（distroless/static-debian12:nonroot）— 僅含 binary 與 `static/`

前端原始碼在 `frontend/`（Vue 3 + Vite），`vite.config.js` 的 `outDir` 設為 `../static`，build 後由 Go 直接 serve。

## Environment Variables

連線設定透過環境變數控制（皆有預設值）：

| 變數 | 預設值 | 說明 |
|------|--------|------|
| `POSTGRES_HOST` | `localhost` | |
| `POSTGRES_PORT` | `5432` | |
| `POSTGRES_USER` | `user` | |
| `POSTGRES_PASSWORD` | `pass` | |
| `POSTGRES_DB` | `shortener` | |
| `REDIS_ADDR` | `localhost:6379` | |
| `BASE_URL` | `http://localhost:8080` | 回傳給使用者的短網址前綴 |

PostgreSQL 預設帳密可參考 `docker-compose-postgres.yaml`（user/pass/shortener）。

## Architecture

單檔案服務（`main.go`），所有邏輯集中於此。

**技術棧：** Go + Gin + PostgreSQL + Redis

**資料流：**
- **寫入**：`POST /shorten` → 生成隨機 8 碼 Base62 短碼 → INSERT 至 PostgreSQL（以 PK 唯一性為準，衝突則重試最多 10 次）→ 寫入 Redis 快取（7 天 TTL）
- **讀取**：`GET /:code` → 先查 Redis（命中則 302 跳轉）→ 未命中則查 PostgreSQL 並回填 Redis（24 小時 TTL）→ 302 跳轉

**Redis 快取格式：** key 為 `url:<code>`，value 為 JSON `{"url":"...","created_at":"..."}`（`URLData` 結構）

**全域變數：**
- `db` — PostgreSQL 連線（`*sql.DB`）
- `rdb` — Redis 客戶端（`*redis.Client`）
- `ctx` — 共用的 `context.Background()`
- `base` — 啟動時從 `BASE_URL` 解析的短網址前綴

**API 端點：**
- `POST /shorten` — 建立短網址，回傳 `{"short_url": "..."}`
- `GET /:code` — 短碼跳轉（302）
- `GET /health` — 健康檢查
- `GET /stats` — 短網址總數（查 PostgreSQL）
- `GET /stats/today` — 今日新建數量（查 PostgreSQL，`created_at >= current_date`）

**前端靜態檔案：** `./static/index.html` 與 `./static/assets/`（由 Gin 直接 serve）。CORS 目前僅允許 `http://localhost:5173`。
