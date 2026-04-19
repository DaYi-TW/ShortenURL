<div align="center">

# ⚡ SnapLink

### 把落落長的網址，變成一個俐落的短連結

[![Go](https://img.shields.io/badge/Go-1.24-00ADD8?style=flat-square&logo=go&logoColor=white)](https://golang.org/)
[![Vue 3](https://img.shields.io/badge/Vue-3-4FC08D?style=flat-square&logo=vue.js&logoColor=white)](https://vuejs.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14+-336791?style=flat-square&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-Cache-DC382D?style=flat-square&logo=redis&logoColor=white)](https://redis.io/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat-square&logo=docker&logoColor=white)](https://www.docker.com/)

</div>

---

## ✨ 為什麼選擇 SnapLink？

- 🚀 **毫秒級跳轉** — Redis 快取優先，直接命中不過 DB
- 🔒 **碰撞保護** — Base62 短碼以 PK 唯一性為準，衝突自動重試
- 📊 **即時統計** — 掌握總量與今日新增數字
- 🐳 **一鍵啟動** — Docker Compose 帶起整個環境，零設定煩惱
- 🎨 **精緻前端** — Vue 3 粒子動畫介面，顏值與效能兼具

---

## 🏗 技術架構

```
使用者瀏覽器
    │
    ▼
Vue 3 + Vite (前端)
    │  POST /shorten
    ▼
Go + Gin (API Server)
    ├─► Redis  ◄──── 快取命中 → 302 跳轉 (毫秒級)
    └─► PostgreSQL ◄─ 未命中 → 查詢回填 → 302 跳轉
```

| 層級 | 技術 |
|------|------|
| 後端 | Go 1.24 + Gin |
| 前端 | Vue 3 + Vite |
| 資料庫 | PostgreSQL 14+ |
| 快取 | Redis（7 天 TTL） |
| 容器化 | Docker + Docker Compose |

---

## 🚀 快速啟動

### 方式一：Docker Compose（推薦）

```bash
# 一行指令，啟動全部服務
docker-compose -f docker-compose-postgres.yaml up -d
```

服務啟動後，打開 [http://localhost:8080](http://localhost:8080) 即可使用。

---

### 方式二：本機開發

**1. 啟動依賴服務**

```bash
# PostgreSQL
docker-compose -f docker-compose-postgres.yaml up -d postgres

# Redis
redis-server
```

**2. 設定環境變數**（可複製 `.env.example`）

```bash
cp .env.example .env
```

**3. 啟動後端**

```bash
go run main.go
```

**4. 啟動前端開發伺服器**

```bash
cd frontend
npm install
npm run dev
```

---

## 🌐 環境變數

| 變數 | 預設值 | 說明 |
|------|--------|------|
| `POSTGRES_HOST` | `localhost` | PostgreSQL 主機位址 |
| `POSTGRES_PORT` | `5432` | PostgreSQL 連接埠 |
| `POSTGRES_USER` | `user` | 資料庫使用者 |
| `POSTGRES_PASSWORD` | `pass` | 資料庫密碼 |
| `POSTGRES_DB` | `shortener` | 資料庫名稱 |
| `REDIS_ADDR` | `localhost:6379` | Redis 位址 |
| `BASE_URL` | `http://localhost:8080` | 短網址前綴（對外網域） |

---

## 📡 API 文件

### `POST /shorten` — 建立短網址

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com/very/long/url/that/nobody/wants/to-type"}'
```

```json
{
  "short_url": "http://localhost:8080/aB3xYz9K"
}
```

---

### `GET /:code` — 短碼跳轉

```
GET /aB3xYz9K  →  302 Redirect → https://example.com/...
```

---

### `GET /stats` — 總數統計

```json
{ "shortened_url_count": 42 }
```

### `GET /stats/today` — 今日新增

```json
{ "shortened_url_count_today": 7 }
```

### `GET /health` — 健康檢查

```json
{ "status": "ok" }
```

---

## ☁️ 部署到 Zeabur

1. Fork 此 Repo
2. 前往 [zeabur.com](https://zeabur.com) → New Project
3. 新增 **PostgreSQL** 與 **Redis** 服務
4. 連結此 GitHub Repo（自動偵測 Dockerfile 建置）
5. 設定環境變數，並將 `BASE_URL` 設為你的 Zeabur 網域
6. 產生網域，部署完成 🎉

---

## 📁 專案結構

```
.
├── main.go                  # 後端主程式（Go + Gin）
├── Dockerfile               # 三階段 Docker build
├── docker-compose-postgres.yaml
├── frontend/                # 前端原始碼（Vue 3 + Vite）
│   └── src/
│       └── components/      # ShortenForm, ResultCard, StatsPanel...
├── static/                  # 前端 build 輸出（由 Go 直接 serve）
├── nginx/                   # Nginx 設定（選用）
└── .env.example             # 環境變數範例
```

---

<div align="center">

Made with ❤️ using Go & Vue 3

</div>
