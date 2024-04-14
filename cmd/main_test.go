package main

import (
	"context"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/AntonCkya/banner-api/internal/cache"
	"github.com/AntonCkya/banner-api/internal/db"
	"github.com/AntonCkya/banner-api/internal/model"
	"github.com/AntonCkya/banner-api/internal/route"
)

type TestCase struct {
	Number     int
	URL        string
	Token      string
	StatusCode int
	Response   string
}

func TestUserBanner(t *testing.T) {
	cases := []TestCase{
		{
			//Не существующий токен
			Number:     1,
			URL:        "http://localhost:8000/user_banner?tag_id=0&feature_id=0&use_last_revision=true",
			Token:      "Bad as helld",
			StatusCode: 401,
			Response:   "{\"error\": \"Пользователь не авторизован\"}",
		},
		{
			//Не существующий баннер
			Number:     2,
			URL:        "http://localhost:8000/user_banner?tag_id=0&feature_id=0&use_last_revision=true",
			Token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjAsImlzQWRtaW4iOmZhbHNlLCJuYW1lIjoiSm9obiBTbWl0aCJ9.l5PzO3w-1HuD3zoq85oujufTSzFFTaV9zmSU3zQpxNo",
			StatusCode: 404,
			Response:   "{\"error\": \"Баннер не найден\"}",
		},
		{
			//Все норм
			Number:     3,
			URL:        "http://localhost:8000/user_banner?tag_id=120403&feature_id=777&use_last_revision=true",
			Token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjAsImlzQWRtaW4iOnRydWUsIm5hbWUiOiJKYW5lIERvZSJ9.-A6bF348vryjwST2vccaW2sgGO6bh7AzmmABdiGKhz0",
			StatusCode: 200,
			Response:   `{"name":"Anton","status":"sleep"}`,
		},
		{
			//Баннер выключен, а токен не админский
			//По логике это Not Found
			Number:     4,
			URL:        "http://localhost:8000/user_banner?tag_id=120403&feature_id=777&use_last_revision=true",
			Token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjAsImlzQWRtaW4iOmZhbHNlLCJuYW1lIjoiSm9obiBTbWl0aCJ9.l5PzO3w-1HuD3zoq85oujufTSzFFTaV9zmSU3zQpxNo",
			StatusCode: 404,
			Response:   "{\"error\": \"Баннер не найден\"}",
		},
	}
	ctx := context.Background()
	conn := db.New(true)
	conn.PostBanner(model.BannerNoID{
		Tag_ids:    []int{120403},
		Feature_id: 777,
		Is_active:  false,
		Content: model.Content{
			"name":   "Anton",
			"status": "sleep",
		},
	})
	cacheConn := cache.New()
	for _, value := range cases {
		r := httptest.NewRequest("GET", value.URL, nil)
		w := httptest.NewRecorder()
		r.Header.Add("token", value.Token)
		handlerfunc := route.UserBannerHandler(ctx, conn, cacheConn)
		handlerfunc(w, r)
		if w.Code != value.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d", value.Number, w.Code, value.StatusCode)
		}
		res := w.Result()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("[%d] failed to read body: %v", value.Number, err)
		}
		bodyStr := string(body)
		if bodyStr != value.Response {
			t.Errorf("[%d] wrong Response: got %+v, expected %+v",
				value.Number, bodyStr, value.Response)
		}
	}
}
