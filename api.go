package rule_engine

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Praser struct {
	operator *TokenOperator
}

// if need use decimal to handle float, set useDecimal: true
// if don't require decimal, can set useDeicimal: false, will use float in calculate.
// if not spefic the type in params, Praser will try to prase the type from value.
// if spefic the type in params, Praser will use this type and analyze the value.
func GetNewPraser(params []*Param, useDecimal bool) (*Praser, error) {
	oper := &TokenOperator{
		decimalMode: useDecimal,
		varMap:      make(map[string]*TokenNode),
	}

	for _, param := range params {
		if param == nil {
			continue
		}
		node, err := parseParam(useDecimal, param)
		if err != nil {
			return nil, err
		}
		oper.varMap[param.Name] = node
	}
	return &Praser{operator: oper}, nil
}

func (p *Praser) Parse(str string) (*TokenNode, error) {
	lex := NewRuleEngineLex(str, p.operator)

	if res := ruleEngineParse(lex); res == Success {
		return lex.resNode, nil
	}
	return nil, lex.err
}

func (p *Praser) CheckValue(node *TokenNode, v interface{}) bool {
	param := &Param{Value: v}
	vnode, err := parseParam(p.operator.decimalMode, param)
	if err != nil {
		return false
	}
	return node.Compare(vnode)
}

type Param struct {
	Name  string      // value name
	Type  ValueType   // value type
	Value interface{} // value
}

func GetParam(key string, value interface{}) *Param {
	return &Param{Name: key, Value: value}
}

func GetParamWithType(key string, valueType ValueType, value interface{}) *Param {
	return &Param{Name: key, Type: valueType, Value: value}
}

type TokenNode struct {
	ValueType ValueType   // result type, can see ValueType
	Value     interface{} // result value
}

func GetTokenNode(valueType ValueType, value interface{}) *TokenNode {
	return &TokenNode{ValueType: valueType, Value: value}
}

func (t *TokenNode) GetValue() interface{} {
	return t.Value
}

func (t *TokenNode) GetInt() int64 {
	switch t.ValueType {
	case ValueTypeInteger:
		return t.Value.(int64)
	case ValueTypeFloat:
		return int64(t.Value.(float64))
	case ValueTypeDecimal:
		return t.GetDecimal().IntPart()
	}
	panic(fmt.Sprintf("invalid type change, from %v to int, value: %v",
		valueTypeNameDict[t.ValueType], t.Value))
}

func (t *TokenNode) GetBool() bool {
	switch t.ValueType {
	case ValueTypeBool:
		return t.Value.(bool)
	}
	panic(fmt.Sprintf("invalid type change, from %v to bool, value: %v",
		valueTypeNameDict[t.ValueType], t.Value))
}

func (t *TokenNode) GetFloat() float64 {
	switch t.ValueType {
	case ValueTypeInteger:
		return float64(t.GetInt())
	case ValueTypeFloat:
		return t.Value.(float64)
	case ValueTypeDecimal:
		return t.GetDecimal().InexactFloat64()
	}
	panic(fmt.Sprintf("invalid type change, from %v to bool, value: %v",
		valueTypeNameDict[t.ValueType], t.Value))
}

func (t *TokenNode) GetDecimal() decimal.Decimal {
	switch t.ValueType {
	case ValueTypeInteger:
		return decimal.NewFromInt(t.Value.(int64))
	case ValueTypeFloat:
		return decimal.NewFromFloat(t.Value.(float64))
	case ValueTypeDecimal:
		return t.Value.(decimal.Decimal)
	}
	panic(fmt.Sprintf("invalid type change, from %v to bool, value: %v",
		valueTypeNameDict[t.ValueType], t.Value))
}

func (t *TokenNode) GetString() string {
	switch t.ValueType {
	case ValueTypeString:
		return t.Value.(string)
	default:
		return fmt.Sprintf("%v", t.Value)
	}
}

func (x *TokenNode) Compare(y *TokenNode) bool {
	if x.ValueType == ValueTypeNone || y.ValueType == ValueTypeNone {
		return false
	}

	if x.ValueType == ValueTypeBool || y.ValueType == ValueTypeBool {
		if x.ValueType != ValueTypeBool || y.ValueType != ValueTypeBool {
			return false
		}
		return x.GetBool() == y.GetBool()
	}

	if x.ValueType == ValueTypeString || y.ValueType == ValueTypeString {
		if x.ValueType != ValueTypeString || y.ValueType != ValueTypeString {
			return false
		}
		return x.GetString() == y.GetString()
	}

	if x.ValueType == ValueTypeDecimal || y.ValueType == ValueTypeDecimal {
		return x.GetDecimal().Equal(y.GetDecimal())
	}

	if x.ValueType == ValueTypeFloat || y.ValueType == ValueTypeFloat {
		return isFloatEqual(x.GetFloat(), y.GetFloat())
	}

	if x.ValueType == ValueTypeInteger && y.ValueType == ValueTypeInteger {
		return x.GetInt() == y.GetInt()
	}
	return false
}
