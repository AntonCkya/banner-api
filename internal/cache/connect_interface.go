package cache

import (
	"context"
	"os"
	"strconv"

	"github.com/AntonCkya/banner-api/internal/model"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type ICacheConnect interface {
	SetContent(ctx context.Context, key string, content model.Content) error
	GetContent(ctx context.Context, key string) (string, error)
	DeleteContent(ctx context.Context, key string) error
}

func New() ICacheConnect {
	godotenv.Load("local.env")
	DBint, _ := strconv.Atoi(os.Getenv("REDISDB"))
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDISADDR"),
		Password: os.Getenv("REDISPASS"),
		DB:       DBint,
	})
	return &CacheConnect{
		Conn: client,
	}
}
