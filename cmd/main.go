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
	conn := db.New()
	cacheConn := cache.New()
	http.ListenAndServe("localhost:8000", route.New(ctx, conn, cacheConn))
}
