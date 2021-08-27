package cli

import (
	"encoding/json"
	"github.com/pkg/errors"
	"os"
)

func ConstructObjectsFromConfig() (map[string]interface{}, error) {
	objects := make(map[string]interface{})

	file, err := os.ReadFile("config.json")

	if err != nil {
		return nil, errors.Wrap(err, "could not open config.json file")
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, errors.Wrap(err, "failed unmarshaling json")
	}

	for _, object := range config.Objects {
		switch object.Type {
		case ObjectMessage:
			objects[object.Name], err = NewMessage(object)
			break
		case ObjectDisplay:
			objects[object.Name], err = NewDisplay(object)
			break
		case ObjectMemory:
			objects[object.Name], err = NewMemory(object)
			break
		default:
			return nil, errors.New("unknown object type: " + string(object.Type))
		}

		if err != nil {
			return nil, err
		}
	}

	return objects, nil
}
