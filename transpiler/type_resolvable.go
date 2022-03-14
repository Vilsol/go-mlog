package transpiler

import (
	"context"
	"strconv"
	"strings"
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

type InlineVariable struct {
	Value Resolvable
}

func (m *InlineVariable) GetValue() string {
	return strings.Trim(m.Value.GetValue(), "\"")
}

func (m *InlineVariable) PreProcess(ctx context.Context, g *Global, function *Function) error {
	return m.Value.PreProcess(ctx, g, function)
}

func (m *InlineVariable) PostProcess(ctx context.Context, g *Global, function *Function) error {
	return m.Value.PostProcess(ctx, g, function)
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
				baseName := "_" + function.Name + "_" + m.Name
				if _, ok := function.ScopeVariableCounter[baseName]; ok {
					m.CalculatedName = baseName + "_" + strconv.Itoa(function.ScopeVariableCounter[baseName])
				} else {
					m.CalculatedName = baseName
				}
				function.ScopeVariableCounter[baseName] = function.ScopeVariableCounter[baseName] + 1
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
