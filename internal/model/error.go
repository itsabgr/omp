package model

import (
	"encoding/json"
	"github.com/itsabgr/omp/internal/utils"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func MarshalError(message string, code int) []byte {
	return utils.Must((&Error{Code: code, Message: message}).Marshal())
}

func (m *Error) Unmarshal(b []byte) error {
	return json.Unmarshal(b, m)
}
func (m *Error) Marshal() ([]byte, error) {
	return json.Marshal(m)
}
