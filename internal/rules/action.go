package rules

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

type Action interface {
	Apply() error
	Pprint()
}

func getActionByName(name string) (Action, error) {
	switch name {
	case "migrate":
		return Migrate{}, nil
	default:
		errMessage := fmt.Sprintf("Unknown action type '%s'", name)
		return nil, errors.New(errMessage)
	}
}

func actionDecodeHook(sourceType, destType reflect.Type, raw interface{}) (interface{}, error) {
	// TODO: find a better way to compare these types
	if fmt.Sprintf("%s", destType) == "rules.Action" {
		var err error
		var result Action

		rawMap := raw.(map[interface{}]interface{})
		result, err = getActionByName(rawMap["type"].(string))
		if err != nil {
			return nil, err
		}

		mapstructure.Decode(raw, &result)
		return result, nil
	}
	return raw, nil
}
