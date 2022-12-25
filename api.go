package rule_engine

import (
	"reflect"

	"github.com/shopspring/decimal"
)

type Praser struct {
	operator *TokenOperator
}

func GetNewPraser(params []*Param, useDecimal bool) (*Praser, error) {
	oper := &TokenOperator{
		decimalMode: useDecimal,
		paramMap:    make(map[string]*Param),
	}

	for _, param := range params {
		if param == nil {
			continue
		}
		err := parseParam(useDecimal, param)
		if err != nil {
			return nil, err
		}
		oper.paramMap[param.Key] = param
	}
	return &Praser{operator: oper}, nil
}

// TODO: add catch panic
// if want use decimal to handle float, can set use Decimal
func (p *Praser) Parse(str string) (*TokenNode, error) {
	lex := NewRuleEngineLex(str, p.operator)

	if res := ruleEngineParse(lex); res == Success {
		return lex.resNode, nil
	}
	return nil, lex.err
}

type Param struct {
	Key   string
	Type  ValueType
	Value interface{}
}

func GetParam(key string, valueType ValueType, value interface{}) *Param {
	return &Param{Key: key, Type: valueType, Value: value}
}

type TokenNode struct {
	ValueType ValueType
	Value     interface{}
}

func GetTokenNode(valueType ValueType, value interface{}) *TokenNode {
	return &TokenNode{ValueType: valueType, Value: value}
}

func (t *TokenNode) GetValue() interface{} {
	return t.Value
}

func (t *TokenNode) CheckValue(v interface{}) bool {
	rt := reflect.ValueOf(v)
	if !rt.IsValid() {
		return false
	}
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	if rt.Kind() == reflect.Interface {
		rt = reflect.ValueOf(rt.Interface())
	}
	switch t.ValueType {
	case ValueTypeString:
		if rt.Kind() != reflect.String {
			return false
		}
		return rt.String() == t.Value.(string)
	case ValueTypeInteger:
		switch rt.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return t.Value.(int64) == rt.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return t.Value.(int64) == int64(rt.Uint())
		case reflect.Struct:
			decimalValue, ok := v.(decimal.Decimal)
			if !ok {
				return false
			}
			return getDecimal(t).Equal(decimalValue)
		default:
			return false
		}
	case ValueTypeFloat:
		switch rt.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return t.Value.(float64) == float64(rt.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return t.Value.(float64) == float64(rt.Uint())
		case reflect.Float32, reflect.Float64:
			return isFloatEqual(t.Value.(float64), rt.Float())
		case reflect.Struct:
			decimalValue, ok := v.(decimal.Decimal)
			if !ok {
				return false
			}
			return getDecimal(t).Equal(decimalValue)
		default:
			return false
		}
	case ValueTypeBool:
		if checkValue, ok := v.(bool); !ok {
			return false
		} else {
			return t.Value.(bool) == checkValue
		}
	case ValueTypeDecimal:
		switch rt.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return decimal.NewFromInt(rt.Int()).Equal(getDecimal(t))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return decimal.NewFromInt(int64(rt.Uint())).Equal(getDecimal(t))
		case reflect.Float32, reflect.Float64:
			return decimal.NewFromFloat(rt.Float()).Equal(getDecimal(t))
		case reflect.Struct:
			decimalValue, ok := v.(decimal.Decimal)
			if !ok {
				return false
			}
			return decimalValue.Equal(getDecimal(t))
		default:
			return false
		}
	}
	return false
}
