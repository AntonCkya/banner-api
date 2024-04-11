package route

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AntonCkya/banner-api/internal/db"
	"github.com/AntonCkya/banner-api/internal/model"
	"github.com/AntonCkya/banner-api/internal/token"
)

func UserBannerHandler(ctx context.Context, conn db.IConnect) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Query := r.URL.Query()

		tag_id, err := strconv.Atoi(Query.Get("tag_id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Некорректные данные"}`))
			return
		}
		feature_id, err := strconv.Atoi(Query.Get("feature_id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Некорректные данные"}`))
			return
		}
		//TODO: кэш
		//use_last_revision := Query.Get("use_last_revision")
		strToken := r.Header.Get("token")

		t := token.New(strToken)

		if !t.Exist() {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Пользователь не авторизован"}`))
			return
			// не знаю что нужно передать, чтобы вызвать 403
			// наверное нужен третий класс недо-пользователей, но это уже звучит странно
		}

		content, err := conn.GetBanner(feature_id, tag_id, t.IsAdmin())
		if errors.Is(err, model.ErrorNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "Баннер не найден"}`))
			return
		} else if errors.Is(err, model.ErrorInternalServerError) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Внутренняя ошибка сервера"}`))
			return
		}

		w.WriteHeader(http.StatusOK)
		js, err := json.Marshal(content)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Внутренняя ошибка сервера"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func BannerHandler(ctx context.Context, conn db.IConnect) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		strToken := r.Header.Get("token")
		t := token.New(strToken)
		if !t.Exist() {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Пользователь не авторизован"}`))
			return
		}
		if !t.IsAdmin() {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error": "Пользователь не имеет доступа"}`))
			return
		}
		if r.URL.Path == "/banner" {
			switch r.Method {
			case http.MethodPost:
				BannerPost(w, r, conn)
			case http.MethodGet:
				BannerGet(w, r, conn)
			}
		}
	}
}

func BannerGet(w http.ResponseWriter, r *http.Request, conn db.IConnect) {
	Query := r.URL.Query()

	tag_id, err := strconv.Atoi(Query.Get("tag_id"))
	if err != nil {
		tag_id = -1
	}
	feature_id, err := strconv.Atoi(Query.Get("feature_id"))
	if err != nil {
		feature_id = -1
	}
	limit, err := strconv.Atoi(Query.Get("limit"))
	if err != nil {
		limit = -1
	}
	offset, err := strconv.Atoi(Query.Get("offset"))
	if err != nil {
		offset = -1
	}

	banners, err := conn.GetAllBanners(feature_id, tag_id, limit, offset)

	if errors.Is(err, model.ErrorInternalServerError) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Внутренняя ошибка сервера"}`))
		return
	} else if errors.Is(err, model.ErrorNotFound) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
		return
	}

	w.WriteHeader(http.StatusOK)
	js, err := json.Marshal(banners)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Внутренняя ошибка сервера"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if len(banners) == 0 {
		w.Write([]byte("[]"))
	} else {
		w.Write(js)
	}
}

func BannerPost(w http.ResponseWriter, r *http.Request, conn db.IConnect) {
	var banner model.BannerNoID
	err := json.NewDecoder(r.Body).Decode(&banner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Некорректные данные"}`))
		return
	}

	banner_id, err := conn.PostBanner(banner)
	if errors.Is(err, model.ErrorInternalServerError) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Внутренняя ошибка сервера"}`))
		return
	} else if errors.Is(err, model.ErrorBadRequest) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Некорректные данные"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"banner_id": "%d"}`, banner_id)))
	w.Header().Set("Content-Type", "application/json")
}

func BannerHandlerID(ctx context.Context, conn db.IConnect) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Некорректные данные"}`))
			return
		}
		strToken := r.Header.Get("token")
		t := token.New(strToken)
		if !t.Exist() {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Пользователь не авторизован"}`))
			return
		}
		if !t.IsAdmin() {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error": "Пользователь не имеет доступа"}`))
			return
		}
		switch r.Method {
		case http.MethodDelete:
			BannerDelete(w, r, conn, id)
		case http.MethodPatch:
			BannerPatch(w, r, conn, id)
		}
	}
}

func BannerDelete(w http.ResponseWriter, r *http.Request, conn db.IConnect, id int) {
	err := conn.DeleteBanner(id)
	if errors.Is(err, model.ErrorInternalServerError) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Внутренняя ошибка сервера"}`))
		return
	} else if errors.Is(err, model.ErrorNotFound) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Баннер не найден"}`))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func BannerPatch(w http.ResponseWriter, r *http.Request, conn db.IConnect, id int) {
	var banner model.BannerNoID
	err := json.NewDecoder(r.Body).Decode(&banner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Некорректные данные"}`))
		return
	}
	err = conn.PatchBanner(id, banner)
	if errors.Is(err, model.ErrorInternalServerError) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Внутренняя ошибка сервера"}`))
		return
	} else if errors.Is(err, model.ErrorNotFound) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Баннер не найден"}`))
		return
	} else if errors.Is(err, model.ErrorBadRequest) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Некорректные данные"}`))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
