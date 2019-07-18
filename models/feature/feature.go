package feature

import (
	"encoding/json"
	"io/ioutil"

	shared "github.com/agile-work/srv-mdl-shared"
	"github.com/agile-work/srv-mdl-shared/models/translation"
)

// Features define a list of features from a module
type Features struct {
	list map[string]feature `validate:"required,dive,required"`
}

// MarshalJSON custom marshal function to return only the map of available features
func (f Features) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.list)
}

type feature struct {
	Mode        string                             `json:"mode" validate:"required"`
	Name        translation.Translation            `json:"name" validate:"required"`
	Description translation.Translation            `json:"description" validate:"required"`
	Permissions map[string]translation.Translation `json:"permissions" validate:"required"`
}

// Load validates the json file and marshall to a struct
func Load(path string) (*Features, error) {
	f := &Features{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &f.list); err != nil {
		return nil, err
	}

	if err := shared.Validate.Struct(f); err != nil {
		return nil, err
	}

	return f, nil
}
