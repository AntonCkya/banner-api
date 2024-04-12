package route

import (
	"context"
	"net/http"

	"github.com/AntonCkya/banner-api/internal/cache"
	"github.com/AntonCkya/banner-api/internal/db"
)

func New(ctx context.Context, conn db.IConnect, cacheConn cache.ICacheConnect) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/user_banner", UserBannerHandler(ctx, conn, cacheConn))
	mux.Handle("/banner", BannerHandler(ctx, conn))
	mux.Handle("/banner/{id}", BannerHandlerID(ctx, conn))

	return mux
}
