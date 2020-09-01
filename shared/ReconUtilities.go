package shared

import (
	"encoding/json"
)

func ToJsonString(resp interface{}) (string, error) {
	bytes, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
