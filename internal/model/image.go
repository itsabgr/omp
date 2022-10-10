package model

import "encoding/json"

type Image struct {
	Sha256    string `json:"sha256"`
	Size      uint   `json:"size"`
	ChunkSize uint   `json:"chunk_size"`
}

func (m *Image) Unmarshal(b []byte) error {
	return json.Unmarshal(b, m)
}
func (m *Image) Marshal() ([]byte, error) {
	return json.Marshal(m)
}
