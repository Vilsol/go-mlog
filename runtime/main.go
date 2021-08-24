package runtime

import (
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"strings"
	"time"
)

const counterVariable = "@counter"

func ExecuteMLOGFile(fileName string, objects map[string]interface{}) (string, error) {
	file, err := ioutil.ReadFile(fileName)

	if err != nil {
		return "", err
	}

	return ExecuteMLOGBytes(file, objects)
}

func ExecuteMLOGBytes(input []byte, objects map[string]interface{}) (string, error) {
	return ExecuteMLOG(string(input), objects)
}

func ExecuteMLOG(input string, objects map[string]interface{}) (string, error) {
	log.Info().Msg("parsing")

	operations, err := Parse(input)

	if err != nil {
		return "", err
	}

	counter := &Variable{
		Value: int64(0),
	}

	context := &ExecutionContext{
		Variables: map[string]*Variable{
			counterVariable: counter,
		},
		PrintBuffer: strings.Builder{},
		Objects:     objects,
	}

	start := time.Now()

	log.Info().Msg("starting execution")

	executedCount := int64(0)
	operationCount := int64(len(operations))
	for counter.Value.(int64) < operationCount {
		operation := operations[counter.Value.(int64)]
		counter.Value = counter.Value.(int64) + 1
		operation(context)
		executedCount++
	}

	log.Info().Dur("took", time.Now().Sub(start)).Int64("executed", executedCount).Msg("completed execution")

	for _, object := range objects {
		if post, ok := object.(PostExecute); ok {
			post.PostExecute()
		}
	}

	return "", nil
}
