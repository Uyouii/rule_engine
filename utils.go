package rule_engine

import (
	"fmt"
	"math"
	"reflect"

	"github.com/shopspring/decimal"
)

func isFloatEqual(x, y float64) bool {
	return math.Abs(x-y) < 0.00000000001
}

func checkOperType(t *TokenNode, oper_type operType, oper_name string) error {
	validTypeList, ok := operValidType[oper_type]
	if !ok {
		return GetError(ErrRuleEngineUnknownOperator, fmt.Sprintf("unkonwn operator, oper: %v", oper_name))
	}

	for _, valueType := range validTypeList {
		if valueType == t.ValueType {
			return nil
		}
	}
	valueType_str, ok := valueTypeNameDict[t.ValueType]
	if !ok {
		valueType_str = ""
	}
	return GetError(ErrRuleEngineNotSupportedOperator, fmt.Sprintf("%v not support operation: %v", valueType_str, oper_name))
}

func batchCheckOperType(token_list []*TokenNode, oper_type operType, oper_name string) error {
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

func parseParam(useDecimal bool, param *Param) (*TokenNode, error) {
	rt := reflect.ValueOf(param.Value)
	if !rt.IsValid() {
		return nil, GetError(ErrRuleEngineInvalidParam, "not valid")
	}
	if rt.Kind() == reflect.Ptr {
		return nil, GetError(ErrRuleEngineInvalidParam,
			fmt.Sprintf("not support point args, params: %v", param.Value))
	}
	if rt.Kind() == reflect.Interface {
		rt = reflect.ValueOf(rt.Interface())
	}

	notMatchErr := GetError(ErrRuleEngineParamValueTypeNotMatch,
		fmt.Sprintf("value: %v, type: %v", param.Value, valueTypeNameDict[param.Type]))

	var resType ValueType
	var resValue interface{}

	switch rt.Kind() {
	case reflect.Invalid:
		return nil, GetError(ErrRuleEngineInvalidParam,
			fmt.Sprintf("invalid params: %v", param.Value))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch param.Type {
		case ValueTypeNone, ValueTypeInteger:
			resType, resValue = ValueTypeInteger, rt.Int()
		case ValueTypeDecimal:
			resType, resValue = ValueTypeDecimal, decimal.NewFromInt(rt.Int())
		case ValueTypeFloat:
			resType, resValue = ValueTypeFloat, float64(rt.Int())
		default:
			return nil, notMatchErr
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch param.Type {
		case ValueTypeNone, ValueTypeInteger:
			resType, resValue = ValueTypeInteger, int64(rt.Uint())
		case ValueTypeFloat:
			resType, resValue = ValueTypeFloat, float64(rt.Uint())
		case ValueTypeDecimal:
			resType, resValue = ValueTypeDecimal, decimal.NewFromInt(int64(rt.Uint()))
		default:
			return nil, notMatchErr
		}
	case reflect.Bool:
		if param.Type != ValueTypeNone && param.Type != ValueTypeBool {
			return nil, notMatchErr
		}
		resType, resValue = ValueTypeBool, rt.Bool()
	case reflect.String:
		switch param.Type {
		case ValueTypeNone, ValueTypeString:
			resType, resValue = ValueTypeString, rt.String()
		case ValueTypeDecimal:
			decimalValue, err := decimal.NewFromString(rt.String())
			if err != nil {
				return nil, GetError(ErrRuleEngineDecimalError, fmt.Sprintf("msg: %v", err))
			}
			resType, resValue = ValueTypeDecimal, decimalValue
		default:
			return nil, notMatchErr
		}
	case reflect.Float32, reflect.Float64:
		switch param.Type {
		case ValueTypeNone, ValueTypeFloat:
			if useDecimal {
				resType, resValue = ValueTypeDecimal, decimal.NewFromFloat(rt.Float())
			} else {
				resType, resValue = ValueTypeFloat, rt.Float()
			}
		case ValueTypeDecimal:
			resType, resValue = ValueTypeDecimal, decimal.NewFromFloat(rt.Float())
		default:
			return nil, notMatchErr
		}
	case reflect.Struct:
		value := rt.Interface()
		decimalValue, ok := value.(decimal.Decimal)
		if !ok {
			return nil, GetError(ErrRuleEngineNotSupportedVarType, fmt.Sprintf("value: %v", param.Value))
		}
		resType, resValue = ValueTypeDecimal, decimalValue
	default:
		return nil, GetError(ErrRuleEngineNotSupportedVarType, fmt.Sprintf("value: %v", param.Value))
	}
	return GetTokenNode(resType, resValue), nil
}
