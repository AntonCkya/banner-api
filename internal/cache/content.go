package cache

import (
	"context"
	"time"

	"github.com/AntonCkya/banner-api/internal/model"
)

func (r *CacheConnect) SetContent(ctx context.Context, key string, content model.Content) error {
	contentTime := model.ContentTime{
		Body: content,
		Time: time.Now(),
	}
	return r.Conn.Set(ctx, key, contentTime, 0).Err()
}

func (r *CacheConnect) GetContent(ctx context.Context, key string) (string, error) {
	return r.Conn.Get(ctx, key).Result()
}

func (r *CacheConnect) DeleteContent(ctx context.Context, key string) error {
	return r.Conn.Del(ctx, key).Err()
}
