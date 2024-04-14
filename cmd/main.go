package main

import (
	"context"
	"net/http"

	"github.com/AntonCkya/banner-api/internal/cache"
	"github.com/AntonCkya/banner-api/internal/db"
	"github.com/AntonCkya/banner-api/internal/route"
)

func main() {
	ctx := context.Background()
	conn := db.New(false)
	cacheConn := cache.New()
	http.ListenAndServe("localhost:8080", route.New(ctx, conn, cacheConn))
}
