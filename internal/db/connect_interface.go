package db

import (
	"database/sql"
	"os"

	"github.com/AntonCkya/banner-api/internal/model"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type IConnect interface {
	GetBanner(feature_id int, tag_id int, isAdmin bool) (model.Content, error)
	GetAllBanners(feature_id int, tag_id int, limit int, offset int) ([]model.Banner, error)
	PostBanner(banner model.BannerNoID) (int, error)
	PatchBanner(id int, new_banner model.BannerNoID) error
	DeleteBanner(id int) error
}

func New() IConnect {
	godotenv.Load("local.env")
	connURL := os.Getenv("DB")
	db, _ := sql.Open("postgres", connURL)
	c := Connect{
		Conn: db,
	}
	return &c
}
