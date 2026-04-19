package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

const (
	cacheKeyPrefix     = "url:"
	pgUniqueViolation  = pq.ErrorCode("23505")
	maxCodeRetries     = 10
)

type URLRequest struct {
	URL string `json:"url" binding:"required"`
}

type URLData struct {
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
}

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	ctx     = context.Background()
	rdb     *redis.Client
	db      *sql.DB
	base    string // resolved once at startup
)

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func initPostgres() {
	var err error
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		envOr("POSTGRES_HOST", "localhost"),
		envOr("POSTGRES_PORT", "5432"),
		envOr("POSTGRES_USER", "user"),
		envOr("POSTGRES_PASSWORD", "pass"),
		envOr("POSTGRES_DB", "shortener"),
	)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("❌ 連線 Postgres 失敗:", err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS short_urls (
			code       VARCHAR(10) PRIMARY KEY,
			url        TEXT        NOT NULL,
			created_at TIMESTAMP   DEFAULT now()
		);`)
	if err != nil {
		log.Fatal("❌ 建表失敗:", err)
	}
	log.Println("✅ Postgres 已連線")
}

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     envOr("REDIS_ADDR", "localhost:6379"),
		Password: envOr("REDIS_PASSWORD", ""),
	})
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatal("❌ 連線 Redis 失敗:", err)
	}
	log.Println("✅ Redis 已連線")
}

func randString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func cacheURL(code, url, createdAt string, ttl time.Duration) error {
	data := URLData{URL: url, CreatedAt: createdAt}
	jsonValue, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, cacheKeyPrefix+code, jsonValue, ttl).Err()
}

func isDuplicateKey(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == pgUniqueViolation
}

func main() {
	if v := os.Getenv("BASE_URL"); v != "" {
		base = v
	} else {
		base = "http://localhost:8080"
	}

	initPostgres()
	initRedis()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// Serve built frontend
	r.Static("/assets", "./static/assets")
	r.StaticFile("/", "./static/index.html")

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/stats", func(c *gin.Context) {
		var count int
		if err := db.QueryRow("SELECT COUNT(*) FROM short_urls").Scan(&count); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve stats"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"shortened_url_count": count})
	})

	r.GET("/stats/today", func(c *gin.Context) {
		var count int
		err := db.QueryRow(
			"SELECT COUNT(*) FROM short_urls WHERE created_at >= current_date",
		).Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve stats"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"shortened_url_count_today": count})
	})

	r.POST("/shorten", func(c *gin.Context) {
		var req URLRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Generate a unique code; rely on DB primary-key constraint as the
		// authoritative uniqueness check — Redis may not hold expired codes.
		var code string
		var createdAt time.Time
		for i := 0; i < maxCodeRetries; i++ {
			code = randString(8)
			err := db.QueryRow(
				"INSERT INTO short_urls (code, url) VALUES ($1, $2) RETURNING created_at",
				code, req.URL,
			).Scan(&createdAt)
			if err == nil {
				break
			}
			// Any error other than a duplicate-key violation is fatal.
			if !isDuplicateKey(err) {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save to DB"})
				return
			}
			if i == maxCodeRetries-1 {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate unique code"})
				return
			}
		}

		// Cache write is best-effort; the URL is already durably stored in DB.
		if err := cacheURL(code, req.URL, createdAt.Format(time.RFC3339), 7*24*time.Hour); err != nil {
			log.Printf("⚠️ 快取寫入失敗 (code=%s): %v", code, err)
		}

		c.JSON(http.StatusOK, gin.H{"short_url": base + "/" + code})
	})

	r.GET("/:code", func(c *gin.Context) {
		code := c.Param("code")

		val, err := rdb.Get(ctx, cacheKeyPrefix+code).Result()
		if err == nil {
			var data URLData
			if err := json.Unmarshal([]byte(val), &data); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid cache entry"})
				return
			}
			c.Redirect(http.StatusFound, data.URL)
			return
		}

		var url string
		var createdAt time.Time
		err = db.QueryRow(
			"SELECT url, created_at FROM short_urls WHERE code=$1", code,
		).Scan(&url, &createdAt)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
			return
		}

		// Best-effort cache re-population; don't block the redirect on failure.
		if err := cacheURL(code, url, createdAt.Format(time.RFC3339), 24*time.Hour); err != nil {
			log.Printf("⚠️ 快取回填失敗 (code=%s): %v", code, err)
		}

		c.Redirect(http.StatusFound, url)
	})

	r.Run(":" + envOr("PORT", "8080"))
}
