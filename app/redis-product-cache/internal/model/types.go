// internal/model/types.go
package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONString for []string
type JSONString []string

func (js *JSONString) Scan(value interface{}) error {
	if value == nil {
		*js = make(JSONString, 0)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan JSONString")
	}
	return json.Unmarshal(bytes, js)
}

func (js JSONString) Value() (driver.Value, error) {
	if js == nil {
		return nil, nil
	}
	return json.Marshal(js)
}

// JSONMap for map[string]string
type JSONMap map[string]string

func (jm *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*jm = make(JSONMap)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan JSONMap")
	}
	return json.Unmarshal(bytes, jm)
}

func (jm JSONMap) Value() (driver.Value, error) {
	if jm == nil {
		return nil, nil
	}
	return json.Marshal(jm)
}
