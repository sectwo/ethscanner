package utils

import (
	"context"
	"encoding/json"
)

func Context() context.Context {
	return context.Background()
}

func ToJson(req interface{}) (interface{}, error) {
	var res interface{}

	if bytes, err := json.Marshal(req); err != nil {
		return nil, err
	} else if err = json.Unmarshal(bytes, &res); err != nil {
		return nil, err
	} else {
		jMap := res.(map[string]interface{})
		for key, value := range jMap {
			jMap[key] = value
		}
		return jMap, nil
	}

}
