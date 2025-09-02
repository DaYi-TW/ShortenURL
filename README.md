# URL Shortener Service

一個使用 **Golang + Gin + PostgreSQL + Redis** 建置的短網址服務。
主要功能包括： - 短網址生成
- 短網址重定向
- 健康檢查
- 短網址統計

------------------------------------------------------------------------

## 🚀 功能介紹

1.  **短網址生成 (`POST /shorten`)**
    -   接收一個長網址，生成隨機 8 碼短碼。
    -   短碼會存入 PostgreSQL 與 Redis（Redis 有效期 7 天）。
2.  **短網址重定向 (`GET /:code`)**
    -   使用短碼查詢原始網址。
    -   先查 Redis，若無則查 PostgreSQL 並回填快取。
    -   302 重定向至原始網址。
3.  **健康檢查 (`GET /health`)**
    -   回傳服務狀態 `{"status": "ok"}`。
4.  **統計資訊**
    -   `GET /stats`：目前 Redis 中短網址數量。
    -   `GET /stats/today`：今日新建短網址數量。

------------------------------------------------------------------------

## 🛠️ 技術架構

-   **語言/框架**：Golang + Gin
-   **資料庫**：PostgreSQL
-   **快取**：Redis（加速短碼查詢）
-   **JSON 格式**：`UrlData` 結構（含 `url` 與 `created_at`）

------------------------------------------------------------------------

## 📦 安裝與啟動

### 1. 環境需求

-   Go 1.20+
-   PostgreSQL 14+
-   Redis 6+

### 2. 設定資料庫

啟動 PostgreSQL，並建立資料庫：

``` sql
CREATE DATABASE shortener;
```

### 3. 修改連線設定

在 `initPostgres()` 中修改連線字串：

``` go
connStr := "host=localhost port=5432 user=user password=pass dbname=shortener sslmode=disable"
```

### 4. 啟動 Redis

``` bash
redis-server
```

### 5. 啟動服務

``` bash
go run main.go
```

------------------------------------------------------------------------

## 📡 API 使用方式

### 1. 建立短網址

**Request**

``` bash
POST /shorten
Content-Type: application/json

{
  "url": "https://example.com/long-url"
}
```

**Response**

``` json
{
  "short_url": "http://localhost:8080/abc123xy"
}
```

------------------------------------------------------------------------

### 2. 使用短碼跳轉

**Request**

``` bash
GET /abc123xy
```

**Response** - 302 Redirect → `https://example.com/long-url`

------------------------------------------------------------------------

### 3. 健康檢查

``` bash
GET /health
```

**Response**

``` json
{"status":"ok"}
```

------------------------------------------------------------------------

### 4. 總數統計

``` bash
GET /stats
```

**Response**

``` json
{"shortened_url_count": 12}
```

------------------------------------------------------------------------

### 5. 今日新增數量

``` bash
GET /stats/today
```

**Response**

``` json
{"shortened_url_count_today": 3}
```

------------------------------------------------------------------------

## ⚠️ 注意事項

-   Redis 目前設置快取有效期為 **7 天**，過期會回退至 PostgreSQL 查詢。
-   `/stats/today` 目前比對方式為 `CreatedAt`，使用 **字串比對
    YYYY-MM-DD**，如果需要精確判斷，建議改為 timestamp。
-   短碼長度固定為 **8 碼 Base62**，若需縮短或調整請修改
    `RandString()`。
