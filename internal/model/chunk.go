package model

import "encoding/json"

type Chunk struct {
	Id    uint   `json:"id"`
	Size  uint   `json:"size"`
	Data  string `json:"data"`
	Image string `json:"image,omitempty"`
}

func (m *Chunk) Unmarshal(b []byte) error {
	return json.Unmarshal(b, m)
}
func (m *Chunk) Marshal() ([]byte, error) {

	return json.Marshal(m)
}
