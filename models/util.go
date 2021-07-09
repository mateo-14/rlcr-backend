package models

import (
	"encoding/json"
)

type Jsoner interface {
	ToMap() (map[string]interface{}, error)
}

type ModelBase struct{}

func FromMap(data map[string]interface{}, v Jsoner) (*[]byte, error) {
	datab, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &datab, json.Unmarshal(datab, &v)
}

func toMap(m Jsoner) (map[string]interface{}, error) {
	var modelMap map[string]interface{}
	if modelMapb, err := json.Marshal(&m); err != nil {
		return modelMap, err
	} else {
		err := json.Unmarshal(modelMapb, &modelMap)
		return modelMap, err
	}
}
