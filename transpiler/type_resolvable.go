package transpiler

import (
	"context"
	"strconv"
)

type Value struct {
	Value string
}

func (m *Value) GetValue() string {
	return m.Value
}

func (m *Value) PreProcess(context.Context, *Global, *Function) error {
	return nil
}

func (m *Value) PostProcess(context.Context, *Global, *Function) error {
	return nil
}

func (m *Value) String() string {
	return m.Value
}

type NormalVariable struct {
	Name           string
	CalculatedName string
}

func (m *NormalVariable) PreProcess(_ context.Context, global *Global, function *Function) error {
	if m.CalculatedName == "" {
		if m.Name == "_" {
			m.CalculatedName = "@_"
		} else {
			if _, ok := global.Constants[m.Name]; ok {
				m.CalculatedName = m.Name
			} else {
				m.CalculatedName = "_" + function.Name + "_" + m.Name
			}
		}
	}
	return nil
}

func (m *NormalVariable) PostProcess(context.Context, *Global, *Function) error {
	return nil
}

func (m *NormalVariable) GetValue() string {
	if m.CalculatedName == "" {
		panic("PreProcess not called on NormalVariable (" + m.Name + ")")
	}
	return m.CalculatedName
}

type DynamicVariable struct {
	Name string
}

func (m *DynamicVariable) PreProcess(_ context.Context, _ *Global, function *Function) error {
	if m.Name == "" {
		suffix := function.VariableCounter
		function.VariableCounter += 1
		m.Name = "_" + function.Name + "_" + strconv.Itoa(suffix)
	}
	return nil
}

func (m *DynamicVariable) PostProcess(context.Context, *Global, *Function) error {
	return nil
}

func (m *DynamicVariable) GetValue() string {
	if m.Name == "" {
		panic("PreProcess not called on DynamicVariable")
	}
	return m.Name
}
