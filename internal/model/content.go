package model

import (
	"encoding/json"
	"time"
)

// баннер представляет собой  JSON-документ неопределенной структуры
type Content map[string]interface{}

type ContentTime struct {
	Body Content
	Time time.Time
}

func (c ContentTime) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c ContentTime) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &c)
}
