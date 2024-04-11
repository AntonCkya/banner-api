package db

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/AntonCkya/banner-api/internal/model"
	_ "github.com/lib/pq"
)

const (
	getBannerQuery = `
	SELECT content, is_active
	FROM Banners
	JOIN Banner_tags
	ON Banner_tags.banner_id = Banners.banner_id
	WHERE feature_id = $1 AND Banner_tags.tag_id = $2; 
	`
	getAllBannerQuery = `
	SELECT Banners.banner_id, Banners.feature_id, Banners.content, Banners.is_active, Banners.created_at, Banners.updated_at
	FROM Banners
	JOIN Banner_tags
	ON Banner_tags.banner_id = Banners.banner_id
	WHERE feature_id = $1 AND Banner_tags.tag_id = $2
	LIMIT $3
	OFFSET $4;
	`
	getAllBannerQueryNoFeature = `
	SELECT Banners.banner_id, Banners.feature_id, Banners.content, Banners.is_active, Banners.created_at, Banners.updated_at
	FROM Banners
	JOIN Banner_tags
	ON Banner_tags.banner_id = Banners.banner_id
	WHERE Banner_tags.tag_id = $1
	LIMIT $2
	OFFSET $3;
	`
	getAllBannerQueryNoTag = `
	SELECT Banners.banner_id, Banners.feature_id, Banners.content, Banners.is_active, Banners.created_at, Banners.updated_at
	FROM Banners
	WHERE feature_id = $1
	LIMIT $2
	OFFSET $3;
	`
	getAllBannerQueryNoAll = `
	SELECT Banners.banner_id, Banners.feature_id, Banners.content, Banners.is_active, Banners.created_at, Banners.updated_at
	FROM Banners
	LIMIT $1
	OFFSET $2; 
	`
	getAllTagsQuery = `
	SELECT tag_id
	FROM Banner_tags
	WHERE banner_id = $1;
	`
	postBannerQuery = `
	INSERT INTO Banners
	(feature_id, content, is_active)
	VALUES
	($1, $2, $3)
	RETURNING banner_id;
	`
	postFeatureQuery = `
	INSERT INTO Banner_tags
	(banner_id, tag_id)
	VALUES
	($1, $2);
	`
	deleteBannerQuery = `
	DELETE FROM Banners
	WHERE banner_id = $1;
	`
	patchBannerQuery = `
	UPDATE Banners
	SET
	feature_id = $1,
	content = $2,
	is_active = $3
	WHERE banner_id = $4;
	`
	deleteBannerTags = `
	DELETE FROM Banner_tags
	WHERE banner_id = $1;
	`
)

func (c *Connect) GetBanner(feature_id int, tag_id int, isAdmin bool) (model.Content, error) {
	var content_str []byte
	var is_active bool
	row := c.Conn.QueryRow(getBannerQuery, feature_id, tag_id)
	// оно не в строку не парсится :(((
	err := row.Scan(&content_str, &is_active)
	if errors.Is(err, sql.ErrNoRows) || (is_active == false && isAdmin == false) {
		return model.Content{}, model.ErrorNotFound
	} else if err != nil {
		return model.Content{}, model.ErrorInternalServerError
	}
	var content model.Content
	json.Unmarshal(content_str, &content)
	return content, nil
}
func (c *Connect) GetAllBanners(feature_id int, tag_id int, limit int, offset int) ([]model.Banner, error) {
	if offset == -1 {
		offset = 0
	}
	var rows *sql.Rows
	var err error
	switch {
	case feature_id == -1 && tag_id == -1:
		if limit == -1 {
			rows, err = c.Conn.Query(getAllBannerQueryNoAll, nil, offset)
		} else {
			rows, err = c.Conn.Query(getAllBannerQueryNoAll, limit, offset)
		}
	case feature_id == -1:
		if limit == -1 {
			rows, err = c.Conn.Query(getAllBannerQueryNoFeature, tag_id, nil, offset)
		} else {
			rows, err = c.Conn.Query(getAllBannerQueryNoFeature, tag_id, limit, offset)
		}
	case tag_id == -1:
		if limit == -1 {
			rows, err = c.Conn.Query(getAllBannerQueryNoTag, feature_id, nil, offset)
		} else {
			rows, err = c.Conn.Query(getAllBannerQueryNoTag, feature_id, limit, offset)
		}
	default:
		if limit == -1 {
			rows, err = c.Conn.Query(getAllBannerQuery, feature_id, tag_id, nil, offset)
		} else {
			rows, err = c.Conn.Query(getAllBannerQuery, feature_id, tag_id, limit, offset)
		}
	}
	defer rows.Close()
	var banners []model.Banner
	if errors.Is(err, sql.ErrNoRows) {
		return []model.Banner{}, model.ErrorNotFound
	} else if err != nil {
		return []model.Banner{}, model.ErrorInternalServerError
	}

	for rows.Next() {
		var b model.Banner
		var content_str []byte
		err := rows.Scan(&b.Banner_id, &b.Feature_id, &content_str, &b.Is_active, &b.Created_at, &b.Updated_at)
		if err != nil {
			return []model.Banner{}, model.ErrorInternalServerError
		}
		json.Unmarshal(content_str, &b.Content)
		rowsTags, err := c.Conn.Query(getAllTagsQuery, b.Banner_id)
		if err != nil {
			return []model.Banner{}, model.ErrorInternalServerError
		}
		var tags []int
		for rowsTags.Next() {
			var t int
			err := rowsTags.Scan(&t)
			if err != nil {
				return []model.Banner{}, model.ErrorInternalServerError
			}
			tags = append(tags, t)
		}
		b.Tag_ids = tags
		banners = append(banners, b)
	}
	return banners, nil
}
func (c *Connect) PostBanner(banner model.BannerNoID) (int, error) {
	var banner_id int
	content_str, err := json.Marshal(banner.Content)
	if err != nil {
		return -1, model.ErrorBadRequest
	}
	row := c.Conn.QueryRow(postBannerQuery, banner.Feature_id, content_str, banner.Is_active)
	err = row.Scan(&banner_id)
	if err != nil {
		return -1, model.ErrorInternalServerError
	}
	for _, tag_id := range banner.Tag_ids {
		_, err = c.Conn.Exec(postFeatureQuery, banner_id, tag_id)
		if err != nil {
			return -1, model.ErrorInternalServerError
		}
	}
	return banner_id, nil
}
func (c *Connect) PatchBanner(id int, new_banner model.BannerNoID) error {
	content_str, err := json.Marshal(new_banner.Content)
	if err != nil {
		return model.ErrorBadRequest
	}
	row, err := c.Conn.Exec(patchBannerQuery, new_banner.Feature_id, content_str, new_banner.Is_active, id)
	if err != nil {
		return model.ErrorInternalServerError
	}
	count, err := row.RowsAffected()
	if err != nil {
		return model.ErrorInternalServerError
	}
	if count == 0 {
		return model.ErrorNotFound
	}
	_, err = c.Conn.Exec(deleteBannerTags, id)
	if err != nil {
		return model.ErrorInternalServerError
	}
	for _, tag_id := range new_banner.Tag_ids {
		_, err = c.Conn.Exec(postFeatureQuery, id, tag_id)
		if err != nil {
			return model.ErrorInternalServerError
		}
	}
	return nil
}
func (c *Connect) DeleteBanner(id int) error {
	row, err := c.Conn.Exec(deleteBannerQuery, id)
	if err != nil {
		return model.ErrorInternalServerError
	}
	count, err := row.RowsAffected()
	if err != nil {
		return model.ErrorInternalServerError
	}
	if count == 0 {
		return model.ErrorNotFound
	}
	return nil
}
