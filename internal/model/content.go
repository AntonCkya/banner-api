package model

import "encoding/json"

// баннер представляет собой  JSON-документ неопределенной структуры
type Content map[string]interface{}

func (c Content) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c Content) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &c)
}
