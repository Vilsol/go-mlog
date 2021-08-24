package runtime

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"strconv"
)

func (c ExecutionContext) Resolve(val string) interface{} {
	if variable, ok := c.Variables[val]; ok {
		return variable.Value
	}

	// Check if basic type
	if val[0] == '"' && val[len(val)-1] == '"' {
		return val[1 : len(val)-1]
	}

	if n, err := strconv.ParseInt(val, 10, 64); err == nil {
		return n
	}

	if n, err := strconv.ParseFloat(val, 64); err == nil {
		return n
	}

	// Else null
	return "null"
}

func (c ExecutionContext) ResolveStr(val string) string {
	result := c.Resolve(val)

	log.Trace().
		Interface("val", val).
		Interface("result", result).
		Msgf("resolving string (got %T)", result)

	switch cast := result.(type) {
	case string:
		return cast
	case bool:
		if cast {
			return "true"
		}

		return "false"
	case int:
		return strconv.FormatInt(int64(cast), 10)
	case int8:
		return strconv.FormatInt(int64(cast), 10)
	case int16:
		return strconv.FormatInt(int64(cast), 10)
	case int32:
		return strconv.FormatInt(int64(cast), 10)
	case int64:
		return strconv.FormatInt(cast, 10)
	case uint:
		return strconv.FormatUint(uint64(cast), 10)
	case uint8:
		return strconv.FormatUint(uint64(cast), 10)
	case uint16:
		return strconv.FormatUint(uint64(cast), 10)
	case uint32:
		return strconv.FormatUint(uint64(cast), 10)
	case uint64:
		return strconv.FormatUint(cast, 10)
	case float32:
		return strconv.FormatFloat(float64(cast), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(cast, 'f', -1, 64)
	}

	panic(fmt.Sprintf("unkonwn value type: %t", result))
}

func (c ExecutionContext) ResolveFloat(val string) float64 {
	result := c.Resolve(val)

	log.Trace().
		Interface("val", val).
		Interface("result", result).
		Msgf("resolving float64 (got %T)", result)

	switch cast := result.(type) {
	case string:
		if v, err := strconv.ParseFloat(val, 64); err == nil {
			return v
		}

		return 0
	case int:
		return float64(cast)
	case int8:
		return float64(cast)
	case int16:
		return float64(cast)
	case int32:
		return float64(cast)
	case int64:
		return float64(cast)
	case uint:
		return float64(cast)
	case uint8:
		return float64(cast)
	case uint16:
		return float64(cast)
	case uint32:
		return float64(cast)
	case uint64:
		return float64(cast)
	case float32:
		return float64(cast)
	case float64:
		return cast
	}

	panic(fmt.Sprintf("unkonwn value type: %t", result))
}

func (c ExecutionContext) ResolveInt(val string) int64 {
	result := c.Resolve(val)

	log.Trace().
		Interface("val", val).
		Interface("result", result).
		Msgf("resolving int64 (got %T)", result)

	switch cast := result.(type) {
	case string:
		if v, err := strconv.ParseFloat(val, 64); err == nil {
			return int64(v)
		}

		return 0
	case bool:
		if cast {
			return 1
		}

		return 0
	case int:
		return int64(cast)
	case int8:
		return int64(cast)
	case int16:
		return int64(cast)
	case int32:
		return int64(cast)
	case int64:
		return cast
	case uint:
		return int64(cast)
	case uint8:
		return int64(cast)
	case uint16:
		return int64(cast)
	case uint32:
		return int64(cast)
	case uint64:
		return int64(cast)
	case float32:
		return int64(cast)
	case float64:
		return int64(cast)
	}

	panic(fmt.Sprintf("unkonwn value type: %t", result))
}

func (c ExecutionContext) IsNumber(val string) bool {
	result := c.Resolve(val)

	switch result.(type) {
	case int:
		return true
	case int8:
		return true
	case int16:
		return true
	case int32:
		return true
	case int64:
		return true
	case uint:
		return true
	case uint8:
		return true
	case uint16:
		return true
	case uint32:
		return true
	case uint64:
		return true
	case float32:
		return true
	case float64:
		return true
	}

	return false
}

func (c ExecutionContext) Set(variable string, value interface{}) {
	if _, ok := c.Variables[variable]; !ok {
		c.Variables[variable] = &Variable{}
	}

	// Special condition for counter variable
	if variable == counterVariable {
		if _, ok := value.(int64); !ok {
			value = c.ResolveInt(fmt.Sprintf("%v", value))
		}
	}

	c.Variables[variable].Value = value

	log.Trace().Interface("var", variable).Interface("val", value).Msgf("set (type %T)", value)
}

func (c ExecutionContext) Object(name string) (interface{}, error) {
	if obj, ok := c.Objects[name]; ok {
		return obj, nil
	}

	return nil, errors.New(fmt.Sprintf("no object with name \"%s\" is connected", name))
}

func (c ExecutionContext) Message(name string) (Message, error) {
	obj, err := c.Object(name)

	if err != nil {
		return nil, err
	}

	if message, ok := obj.(Message); ok {
		return message, nil
	}

	return nil, errors.New(fmt.Sprintf("object with name \"%s\" is not a message (is %t)", name, obj))
}

func (c ExecutionContext) Display(name string) (Display, error) {
	obj, err := c.Object(name)

	if err != nil {
		return nil, err
	}

	if display, ok := obj.(Display); ok {
		return display, nil
	}

	return nil, errors.New(fmt.Sprintf("object with name \"%s\" is not a display (is %t)", name, obj))
}
