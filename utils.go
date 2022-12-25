package rule_engine

import (
	"fmt"
	"math"
	"reflect"

	"github.com/shopspring/decimal"
)

func getInt(t *TokenNode) int64 {
	return t.Value.(int64)
}

func getBool(t *TokenNode) bool {
	return t.Value.(bool)
}

func getFloat(t *TokenNode) float64 {
	switch t.ValueType {
	case ValueTypeInteger:
		return float64(getInt(t))
	case ValueTypeFloat:
		return t.Value.(float64)
	}
	panic("invalid operation")
}

func getDecimal(t *TokenNode) decimal.Decimal {
	switch t.ValueType {
	case ValueTypeInteger:
		return decimal.NewFromInt(t.Value.(int64))
	case ValueTypeFloat:
		return decimal.NewFromFloat(t.Value.(float64))
	case ValueTypeDecimal:
		return t.Value.(decimal.Decimal)
	}
	panic("invalid operation")
}

func getString(t *TokenNode) string {
	return t.Value.(string)
}

func isFloatEqual(x, y float64) bool {
	return math.Abs(x-y) < 0.00000000001
}

func checkOperType(t *TokenNode, oper_type OPER_TYPE, oper_name string) error {
	validTypeList, ok := OPER_VALID_TYPE[oper_type]
	if !ok {
		return GetError(ErrRuleEngineUnknownOperator, fmt.Sprintf("unkonwn operator, oper: %v", oper_name))
	}

	for _, valueType := range validTypeList {
		if valueType == t.ValueType {
			return nil
		}
	}
	valueType_str, ok := VALUE_TYPE_NAME_DICT[t.ValueType]
	if !ok {
		valueType_str = ""
	}
	return GetError(ErrRuleEngineNotSupportedOperator, fmt.Sprintf("%v not support operation: %v", valueType_str, oper_name))
}

func batchCheckOperType(token_list []*TokenNode, oper_type OPER_TYPE, oper_name string) error {
	for _, t := range token_list {
		if err := checkOperType(t, oper_type, oper_name); err != nil {
			return err
		}
	}
	return nil
}

func checkFiledType(t *TokenNode, valueType_list []ValueType) error {
	for _, filed_type := range valueType_list {
		if filed_type == t.ValueType {
			return nil
		}
	}
	return &EngineErr{ErrCode: ErrRuleEngineInvalidOperation}
}

func batchCheckFieldType(token_list []*TokenNode, valueType_list []ValueType) error {
	for _, t := range token_list {
		if err := checkFiledType(t, valueType_list); err != nil {
			return err
		}
	}
	return nil
}

func intMin(x int64, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func intMax(x int64, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func intAbs(x int64) int64 {
	if x >= 0 {
		return x
	}
	return -x
}

func getArgNumberError(needArg int, giveArg int) error {
	return GetError(ErrRuleEngineFuncArgument, fmt.Sprintf("func can only handle %v arg, but give %v", needArg, giveArg))
}

func parseParam(useDecimal bool, param *Param) error {
	rt := reflect.ValueOf(param.Value)
	if !rt.IsValid() {
		return GetError(ErrRuleEngineInvalidParam, "not valid")
	}
	if rt.Kind() == reflect.Ptr {
		return GetError(ErrRuleEngineInvalidParam,
			fmt.Sprintf("not support point args, params: %v", param.Value))
	}
	if rt.Kind() == reflect.Interface {
		rt = reflect.ValueOf(rt.Interface())
	}

	notMatchErr := GetError(ErrRuleEngineParamValueTypeNotMatch,
		fmt.Sprintf("value: %v, type: %v", param.Value, VALUE_TYPE_NAME_DICT[param.Type]))

	switch rt.Kind() {
	case reflect.Invalid:
		return GetError(ErrRuleEngineInvalidParam,
			fmt.Sprintf("invalid params: %v", param.Value))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch param.Type {
		case ValueTypeNone, ValueTypeInteger:
			param.Type = ValueTypeInteger
			param.Value = rt.Int()
		case ValueTypeDecimal:
			param.Value = decimal.NewFromInt(rt.Int())
		case ValueTypeFloat:
			param.Value = float64(rt.Int())
		default:
			return notMatchErr
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch param.Type {
		case ValueTypeNone, ValueTypeInteger:
			param.Type = ValueTypeInteger
			param.Value = int64(rt.Uint())
		case ValueTypeFloat:
			param.Value = float64(rt.Uint())
		case ValueTypeDecimal:
			param.Value = decimal.NewFromInt(int64(rt.Uint()))
		default:
			return notMatchErr
		}
	case reflect.Bool:
		if param.Type != ValueTypeNone && param.Type != ValueTypeBool {
			return notMatchErr
		}
		param.Type = ValueTypeBool
		param.Value = rt.Bool()
	case reflect.String:
		switch param.Type {
		case ValueTypeNone, ValueTypeString:
			param.Type = ValueTypeString
			param.Value = rt.String()
		case ValueTypeDecimal:
			decimalValue, err := decimal.NewFromString(rt.String())
			if err != nil {
				return GetError(ErrRuleEngineDecimalError, fmt.Sprintf("msg: %v", err))
			}
			param.Value = decimalValue
		default:
			return notMatchErr
		}
	case reflect.Float32, reflect.Float64:
		switch param.Type {
		case ValueTypeNone, ValueTypeFloat:
			if useDecimal {
				param.Type = ValueTypeDecimal
				param.Value = decimal.NewFromFloat(rt.Float())
			} else {
				param.Type = ValueTypeFloat
				param.Value = rt.Float()
			}
		case ValueTypeDecimal:
			param.Value = decimal.NewFromFloat(rt.Float())
		default:
			return notMatchErr
		}
	case reflect.Struct:
		value := rt.Interface()
		decimalValue, ok := value.(decimal.Decimal)
		if !ok {
			return GetError(ErrRuleEngineNotSupportedVarType, fmt.Sprintf("value: %v", param.Value))
		}
		param.Type = ValueTypeDecimal
		param.Value = decimalValue
	default:
		return GetError(ErrRuleEngineNotSupportedVarType, fmt.Sprintf("value: %v", param.Value))
	}
	return nil
}
