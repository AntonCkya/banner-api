package model

import "time"

type BannerNoID struct {
	Tag_ids    []int   `json:"tag_ids"`
	Feature_id int     `json:"feature_id"`
	Content    Content `json:"content"`
	Is_active  bool    `json:"is_active"`
}

type Banner struct {
	Banner_id  int       `json:"banner_id"`
	Tag_ids    []int     `json:"tag_ids"`
	Feature_id int       `json:"feature_id"`
	Content    Content   `json:"content"`
	Is_active  bool      `json:"is_active"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}
