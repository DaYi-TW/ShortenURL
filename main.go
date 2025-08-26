package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

/*
網址暫存到記憶體
*/

var urlStore = make(map[string]string)

/*
定義輸入結構
*/

type UrlRequest struct {
	Url string `json:"url" binding:"required"`
}

type UrlData struct {
	Url       string `json:"url"`
	CreatedAt string `json:"created_at"`
}

// ----------------- 初始化 -----------------
func initPostgres() {
	var err error
	connStr := "host=localhost port=5432  user=user password=pass dbname=shortener sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("❌ 連線 Postgres 失敗:", err)
	}

	// 建立資料表（如果還沒建立）
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS short_urls (
		code VARCHAR(10) PRIMARY KEY,
		url TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT now()
	);`)
	if err != nil {
		log.Fatal("❌ 建表失敗:", err)
	}
	log.Println("✅ Postgres 已連線")
}

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // docker-compose service name
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("❌ 連線 Redis 失敗:", err)
	}
	log.Println("✅ Redis 已連線")
}

/*
亂碼生成
*/
// Base62 字元集
var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	ctx     = context.Background()
	rdb     = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	db *sql.DB
)

// 生成隨機短碼
func RandString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	initPostgres()
	initRedis()

	r := gin.Default()

	// 短網址生成
	r.POST("/shorten", func(c *gin.Context) {
		var req UrlRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 產生唯一短碼
		var code string
		for {
			code = RandString(8)
			if _, err := rdb.Get(ctx, code).Result(); err == redis.Nil {
				break
			} else if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "redis error"})
				return
			}
		}

		// 寫入 Postgres
		_, err := db.Exec("INSERT INTO short_urls (code, url) VALUES ($1, $2)", code, req.Url)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save to DB"})
			return
		}

		// 建立 JSON 資料 (含創建時間)
		data := UrlData{
			Url:       req.Url,
			CreatedAt: time.Now().Format(time.RFC3339),
		}
		jsonValue, _ := json.Marshal(data)

		// 存入 Redis
		if err := rdb.Set(ctx, code, jsonValue, 7*24*time.Hour).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save"})
			return
		}

		shortUrl := "http://localhost:8080/" + code
		c.JSON(http.StatusOK, gin.H{"short_url": shortUrl})
	})

	//短網址重定向
	r.GET("/:code", func(c *gin.Context) {
		code := c.Param("code")

		// 先查 Redis
		val, err := rdb.Get(ctx, code).Result()
		if err == nil {
			// 命中快取
			var data UrlData
			json.Unmarshal([]byte(val), &data)
			c.Redirect(http.StatusFound, data.Url)
			return
		}

		// Redis 沒有 → 查 Postgres
		var url string
		err = db.QueryRow("SELECT url FROM short_urls WHERE code=$1", code).Scan(&url)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
			return
		}

		// 存回 Redis（避免下次又查 DB）
		data := UrlData{
			Url:       url,
			CreatedAt: time.Now().Format(time.RFC3339),
		}
		jsonValue, _ := json.Marshal(data)
		rdb.Set(ctx, code, jsonValue, 24*time.Hour)

		c.Redirect(http.StatusFound, url)
	})

	//健康檢查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	//取得所有短網址的數量
	r.GET("/stats", func(c *gin.Context) {
		count := 0
		iter := rdb.Scan(ctx, 0, "", 0).Iterator()
		for iter.Next(ctx) {
			count++
		}
		if err := iter.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve stats"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"shortened_url_count": count})
	})

	//今日創建的短網址數量
	r.GET("/stats/today", func(c *gin.Context) {
		count := 0
		today := time.Now().Format("2006-01-02")
		iter := rdb.Scan(ctx, 0, "", 0).Iterator()
		for iter.Next(ctx) {
			var data UrlData
			jsonData, err := rdb.Get(ctx, iter.Val()).Result()
			if err == nil {
				json.Unmarshal([]byte(jsonData), &data)
				if data.CreatedAt == today {
					count++
				}
			}
		}
		if err := iter.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve stats"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"shortened_url_count_today": count})
	})

	r.Run(":8080")
}
