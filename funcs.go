package rule_engine

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

func (o *TokenOperator) tokenNodeArg(arg *TokenNode) (*TokenNode, error) {
	res := &TokenNode{
		ValueType: ValueTypeArgs,
	}
	tokenList := []*TokenNode{arg}

	res.Value = tokenList
	return res, nil
}

func (o *TokenOperator) tokenNodeArgList(argList *TokenNode, arg *TokenNode) (*TokenNode, error) {
	if err := checkOperType(argList, OPER_TYPE_ARGUMENT, "args"); err != nil {
		return nil, err
	}

	res := *argList

	res.Value = append(res.Value.([]*TokenNode), arg)
	return &res, nil
}

func (o *TokenOperator) tokenNodeFunc(funcNode *TokenNode, argNode *TokenNode) (*TokenNode, error) {
	argList := []*TokenNode{}
	if argNode != nil {
		if err := checkOperType(argNode, OPER_TYPE_ARGUMENT, "args"); err != nil {
			return nil, err
		}
		argList = argNode.Value.([]*TokenNode)
	}
	return o.tokenHandleFunc(funcNode, argList)
}

func (o *TokenOperator) tokenHandleFunc(funcNode *TokenNode, argList []*TokenNode) (*TokenNode, error) {
	funcName := funcNode.Value.(string)

	switch funcName {
	case "len":
		return o.funcLen(argList)
	case "min":
		return o.funcMin(argList)
	case "max":
		return o.funcMax(argList)
	case "abs":
		return o.funcAbs(argList)
	case "regexMatch":
		return o.funcRegexMatch(argList)
	case "upper":
		return o.funcUpper(argList)
	case "lower":
		return o.funcLower(argList)
	case "startWith":
		return o.funcStartWith(argList)
	case "endWith":
		return o.funcEndWith(argList)
	case "int":
		return o.funcInt(argList)
	case "float":
		return o.funcFloat(argList)
	case "decimal":
		return o.funcDecimal(argList)
	case "string":
		return o.funcString(argList)
	default:
		return nil, GetError(ErrRuleEngineUnkonwnFunc, fmt.Sprintf("unknown func name: %v", funcName))
	}
}

func (o *TokenOperator) funcString(argList []*TokenNode) (*TokenNode, error) {
	if len(argList) != 1 {
		return nil, getArgNumberError(1, len(argList))
	}

	arg := argList[0]
	if err := checkOperType(arg, OPER_TYPE_CHANGE_TO, "string"); err != nil {
		return nil, err
	}

	switch arg.ValueType {
	case ValueTypeInteger:
		return GetTokenNode(ValueTypeString, fmt.Sprintf("%v", getInt(arg))), nil
	case ValueTypeFloat:
		return GetTokenNode(ValueTypeString, fmt.Sprintf("%v", getFloat(arg))), nil
	case ValueTypeDecimal:
		return GetTokenNode(ValueTypeString, fmt.Sprintf("%v", getDecimal(arg))), nil
	case ValueTypeString:
		return GetTokenNode(ValueTypeString, getString(arg)), nil
	}
	return nil, nil
}

func (o *TokenOperator) funcDecimal(argList []*TokenNode) (*TokenNode, error) {
	if len(argList) != 1 {
		return nil, getArgNumberError(1, len(argList))
	}

	arg := argList[0]
	if err := checkOperType(arg, OPER_TYPE_CHANGE_TO, "decimal"); err != nil {
		return nil, err
	}

	switch arg.ValueType {
	case ValueTypeInteger:
		return GetTokenNode(ValueTypeDecimal, decimal.NewFromInt(getInt(arg))), nil
	case ValueTypeFloat:
		return GetTokenNode(ValueTypeDecimal, decimal.NewFromFloat(getFloat(arg))), nil
	case ValueTypeDecimal:
		return GetTokenNode(ValueTypeDecimal, getDecimal(arg)), nil
	case ValueTypeString:
		value, err := decimal.NewFromString(getString(arg))
		if err != nil {
			return nil, GetError(ErrRuleEngineFuncArgument, fmt.Sprintf("invalid string arg in func decimal, arg: %v", arg.Value))
		}
		return GetTokenNode(ValueTypeDecimal, value), nil
	}
	return nil, nil
}

func (o *TokenOperator) funcFloat(argList []*TokenNode) (*TokenNode, error) {
	if len(argList) != 1 {
		return nil, getArgNumberError(1, len(argList))
	}

	arg := argList[0]
	if err := checkOperType(arg, OPER_TYPE_CHANGE_TO, "float"); err != nil {
		return nil, err
	}

	switch arg.ValueType {
	case ValueTypeInteger:
		return GetTokenNode(ValueTypeFloat, float64(getInt(arg))), nil
	case ValueTypeFloat:
		return GetTokenNode(ValueTypeFloat, getFloat(arg)), nil
	case ValueTypeDecimal:
		return GetTokenNode(ValueTypeFloat, getDecimal(arg).InexactFloat64()), nil
	case ValueTypeString:
		value, err := strconv.ParseFloat(getString(arg), 64)
		if err != nil {
			return nil, GetError(ErrRuleEngineFuncArgument, fmt.Sprintf("invalid string arg in func float, arg: %v", arg.Value))
		}
		return GetTokenNode(ValueTypeFloat, value), nil
	}
	return nil, nil
}

func (o *TokenOperator) funcInt(argList []*TokenNode) (*TokenNode, error) {
	if len(argList) != 1 {
		return nil, getArgNumberError(1, len(argList))
	}

	arg := argList[0]
	if err := checkOperType(arg, OPER_TYPE_CHANGE_TO, "int"); err != nil {
		return nil, err
	}

	switch arg.ValueType {
	case ValueTypeInteger:
		return GetTokenNode(ValueTypeInteger, getInt(arg)), nil
	case ValueTypeFloat:
		return GetTokenNode(ValueTypeInteger, int64(getFloat(arg))), nil
	case ValueTypeDecimal:
		return GetTokenNode(ValueTypeInteger, getDecimal(arg).IntPart()), nil
	case ValueTypeString:
		value, err := strconv.ParseInt(getString(arg), 0, 64)
		if err != nil {
			return nil, GetError(ErrRuleEngineFuncArgument, fmt.Sprintf("invalid string arg in func int, arg: %v", arg.Value))
		}
		return GetTokenNode(ValueTypeInteger, value), nil
	}
	return nil, nil
}

func (o *TokenOperator) funcLower(argList []*TokenNode) (*TokenNode, error) {
	if len(argList) != 1 {
		return nil, getArgNumberError(1, len(argList))
	}

	arg := argList[0]
	if arg.ValueType != ValueTypeString {
		return nil, GetError(ErrRuleEngineFuncArgument, "lower func can onle handle string")
	}

	return GetTokenNode(ValueTypeString, strings.ToLower(getString(arg))), nil
}

func (o *TokenOperator) funcUpper(argList []*TokenNode) (*TokenNode, error) {
	if len(argList) != 1 {
		return nil, getArgNumberError(1, len(argList))
	}

	arg := argList[0]
	if arg.ValueType != ValueTypeString {
		return nil, GetError(ErrRuleEngineFuncArgument, "upper func can onle handle string")
	}

	return GetTokenNode(ValueTypeString, strings.ToUpper(getString(arg))), nil
}

func (o *TokenOperator) funcRegexMatch(argList []*TokenNode) (*TokenNode, error) {
	if len(argList) != 2 {
		return nil, getArgNumberError(2, len(argList))
	}

	if err := batchCheckOperType(argList, OPER_TYPE_REGEX, "regexMatch"); err != nil {
		return nil, err
	}

	pattern := getString(argList[0])
	s := getString(argList[1])

	matched, err := regexp.MatchString(pattern, s)
	if err != nil {
		return nil, GetError(ErrRuleEngineRegexMatch, fmt.Sprintf("regexMatch failed, %v", err))
	}
	return GetTokenNode(ValueTypeBool, matched), nil
}

func (o *TokenOperator) funcLen(argList []*TokenNode) (*TokenNode, error) {
	if len(argList) != 1 {
		return nil, getArgNumberError(1, len(argList))
	}

	arg := argList[0]
	if arg.ValueType != ValueTypeString {
		return nil, GetError(ErrRuleEngineFuncArgument, "len func can onle handle string")
	}

	return GetTokenNode(ValueTypeInteger, int64(len(getString(arg)))), nil
}

func (o *TokenOperator) funcMin(argList []*TokenNode) (*TokenNode, error) {
	if len(argList) < 2 {
		return nil, GetError(ErrRuleEngineFuncArgument,
			fmt.Sprintf("min func take at least 2 arg, but give %v", len(argList)))
	}

	if err := batchCheckOperType(argList, OPER_TYPE_MATH, "min"); err != nil {
		return nil, err
	}

	res := GetTokenNode(argList[0].ValueType, argList[0].Value)

	for _, arg := range argList {
		if res.ValueType == ValueTypeInteger && arg.ValueType == ValueTypeInteger {
			res.Value = intMin(getInt(res), getInt(arg))
		} else if o.decimalMode || res.ValueType == ValueTypeDecimal || arg.ValueType == ValueTypeDecimal {
			res.Value = decimal.Min(getDecimal(res), getDecimal(arg))
			res.ValueType = ValueTypeDecimal
		} else {
			res.Value = math.Min(getFloat(res), getFloat(arg))
			res.ValueType = ValueTypeFloat
		}
	}

	return res, nil
}

func (o *TokenOperator) funcMax(argList []*TokenNode) (*TokenNode, error) {
	if len(argList) < 2 {
		return nil, GetError(ErrRuleEngineFuncArgument,
			fmt.Sprintf("max func take at least 2 arg, but give %v", len(argList)))
	}

	if err := batchCheckOperType(argList, OPER_TYPE_MATH, "max"); err != nil {
		return nil, err
	}

	res := GetTokenNode(argList[0].ValueType, argList[0].Value)

	for _, arg := range argList {
		if res.ValueType == ValueTypeInteger && arg.ValueType == ValueTypeInteger {
			res.Value = intMax(getInt(res), getInt(arg))
		} else if o.decimalMode || res.ValueType == ValueTypeDecimal || arg.ValueType == ValueTypeDecimal {
			res.Value = decimal.Max(getDecimal(res), getDecimal(arg))
			res.ValueType = ValueTypeDecimal
		} else {
			res.Value = math.Max(getFloat(res), getFloat(arg))
			res.ValueType = ValueTypeFloat
		}
	}

	return res, nil
}

func (o *TokenOperator) funcAbs(argList []*TokenNode) (*TokenNode, error) {
	if len(argList) != 1 {
		return nil, getArgNumberError(1, len(argList))
	}

	arg := argList[0]
	if err := checkOperType(arg, OPER_TYPE_MATH, "abs"); err != nil {
		return nil, err
	}

	switch arg.ValueType {
	case ValueTypeInteger:
		return GetTokenNode(ValueTypeInteger, intAbs(getInt(arg))), nil
	case ValueTypeFloat:
		return GetTokenNode(ValueTypeFloat, math.Abs(getFloat(arg))), nil
	case ValueTypeDecimal:
		return GetTokenNode(ValueTypeDecimal, getDecimal(arg).Abs()), nil
	}
	return nil, nil
}

func (o *TokenOperator) funcStartWith(argList []*TokenNode) (*TokenNode, error) {
	if len(argList) != 2 {
		return nil, getArgNumberError(2, len(argList))
	}

	if err := batchCheckOperType(argList, OPER_TYPE_STRING, "startWith"); err != nil {
		return nil, err
	}

	res := strings.HasPrefix(argList[0].Value.(string), argList[1].Value.(string))

	return GetTokenNode(ValueTypeBool, res), nil
}

func (o *TokenOperator) funcEndWith(argList []*TokenNode) (*TokenNode, error) {
	if len(argList) != 2 {
		return nil, getArgNumberError(2, len(argList))
	}

	if err := batchCheckOperType(argList, OPER_TYPE_STRING, "endWith"); err != nil {
		return nil, err
	}

	res := strings.HasSuffix(argList[0].Value.(string), argList[1].Value.(string))

	return GetTokenNode(ValueTypeBool, res), nil
}
