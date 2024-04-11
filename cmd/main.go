package main

import (
	"context"
	"net/http"
	"os"

	"github.com/AntonCkya/banner-api/internal/db"
	"github.com/AntonCkya/banner-api/internal/route"
)

func main() {
	os.Setenv("DB", "postgres://postgres:XD_120403_1000$@localhost:5432/banner?sslmode=disable")
	ctx := context.Background()
	conn := db.New()
	http.ListenAndServe("localhost:8000", route.New(ctx, conn))
}
