package runtime

import (
	"github.com/olekukonko/tablewriter"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

const counterVariable = "@counter"

func ExecuteMLOGFile(fileName string, objects map[string]interface{}) error {
	file, err := ioutil.ReadFile(fileName)

	if err != nil {
		return err
	}

	return ExecuteMLOGBytes(file, objects)
}

func ExecuteMLOGBytes(input []byte, objects map[string]interface{}) error {
	return ExecuteMLOG(string(input), objects)
}

func ExecuteMLOG(input string, objects map[string]interface{}) error {
	operations, err := Parse(input)

	if err != nil {
		return err
	}

	context, counter := ConstructContext(objects)

	return ExecuteContext(operations, context, counter)
}

func ConstructContext(objects map[string]interface{}) (*ExecutionContext, *Variable) {
	log.Info().Msg("parsing")

	counter := &Variable{
		Value: int64(0),
	}

	context := &ExecutionContext{
		Variables: map[string]*Variable{
			counterVariable: counter,
		},
		PrintBuffer: strings.Builder{},
		Objects:     objects,
		Metrics:     make(map[int64]*Metrics),
	}

	return context, counter
}

func ExecuteContext(operations []Operation, context *ExecutionContext, counter *Variable) error {
	metrics := viper.GetBool("metrics")

	log.Info().Msg("starting execution")

	start := time.Now()
	executedCount := int64(0)
	operationCount := int64(len(operations))
	for counter.Value.(int64) < operationCount {
		i := counter.Value.(int64)

		operation := operations[i]
		counter.Value = i + 1
		operation.Executor(context)
		executedCount++

		if metrics {
			if _, ok := context.Metrics[i]; !ok {
				context.Metrics[i] = &Metrics{
					Executions: 0,
				}
			}

			context.Metrics[i].Executions++
		}
	}

	log.Info().Dur("took", time.Since(start)).Int64("executed", executedCount).Msg("completed execution")

	if metrics {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{".", "X", "Log", "#"})
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetColWidth(1000)

		for i, operation := range operations {
			line := make([]string, 0)
			line = append(line, strconv.Itoa(operation.Line.SourceLine))

			if metric, ok := context.Metrics[int64(i)]; ok {
				line = append(line, strconv.FormatUint(metric.Executions, 10))
			} else {
				line = append(line, "0")
			}

			line = append(line, strings.Join(operation.Line.Instruction, " "))
			line = append(line, operation.Line.Comment)

			table.Append(line)
		}

		table.Render()
	}

	for _, object := range context.Objects {
		if post, ok := object.(PostExecute); ok {
			post.PostExecute()
		}
	}

	return nil
}
