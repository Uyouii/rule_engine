package rule_engine

import (
	"fmt"
)

type TokenOperator struct {
	decimalMode bool
	paramMap    map[string]*Param
}

func (o *TokenOperator) tokenNodeAdd(x, y *TokenNode) (*TokenNode, error) {
	if err := batchCheckOperType([]*TokenNode{x, y}, operTypeMath, "+"); err != nil {
		return nil, err
	}

	if x.ValueType == ValueTypeInteger && y.ValueType == ValueTypeInteger {
		return GetTokenNode(ValueTypeInteger, getInt(x)+getInt(y)), nil
	}

	if o.decimalMode || x.ValueType == ValueTypeDecimal || y.ValueType == ValueTypeDecimal {
		res := getDecimal(x).Add(getDecimal(y))
		return GetTokenNode(ValueTypeDecimal, res), nil
	}

	return GetTokenNode(ValueTypeFloat, getFloat(x)+getFloat(y)), nil
}

func (o *TokenOperator) tokenNodeSub(x, y *TokenNode) (*TokenNode, error) {
	if err := batchCheckOperType([]*TokenNode{x, y}, operTypeMath, "-"); err != nil {
		return nil, err
	}

	if x.ValueType == ValueTypeInteger && y.ValueType == ValueTypeInteger {
		return GetTokenNode(ValueTypeInteger, getInt(x)-getInt(y)), nil
	}

	if o.decimalMode || x.ValueType == ValueTypeDecimal || y.ValueType == ValueTypeDecimal {
		res := getDecimal(x).Sub(getDecimal(y))
		return GetTokenNode(ValueTypeDecimal, res), nil
	}

	return GetTokenNode(ValueTypeFloat, getFloat(x)-getFloat(y)), nil
}

func (o *TokenOperator) tokenNodeMul(x, y *TokenNode) (*TokenNode, error) {
	if err := batchCheckOperType([]*TokenNode{x, y}, operTypeMath, "*"); err != nil {
		return nil, err
	}

	if x.ValueType == ValueTypeInteger && y.ValueType == ValueTypeInteger {
		return GetTokenNode(ValueTypeInteger, getInt(x)*getInt(y)), nil
	}

	if o.decimalMode || x.ValueType == ValueTypeDecimal || y.ValueType == ValueTypeDecimal {
		res := getDecimal(x).Mul(getDecimal(y))
		return GetTokenNode(ValueTypeDecimal, res), nil
	}

	return GetTokenNode(ValueTypeFloat, getFloat(x)*getFloat(y)), nil
}

func (o *TokenOperator) tokenNodeDiv(x, y *TokenNode) (*TokenNode, error) {
	if err := batchCheckOperType([]*TokenNode{x, y}, operTypeMath, "/"); err != nil {
		return nil, err
	}

	if getDecimal(y).IsZero() {
		return nil, GetError(ErrRuleEngineDivideByZero, "divide by zero")
	}

	if x.ValueType == ValueTypeInteger && y.ValueType == ValueTypeInteger {
		return GetTokenNode(ValueTypeInteger, getInt(x)/getInt(y)), nil
	}

	if o.decimalMode || x.ValueType == ValueTypeDecimal || y.ValueType == ValueTypeDecimal {
		res := getDecimal(x).Div(getDecimal(y))
		return GetTokenNode(ValueTypeDecimal, res), nil
	}

	return GetTokenNode(ValueTypeFloat, getFloat(x)/getFloat(y)), nil
}

func (o *TokenOperator) tokenNodeMod(x, y *TokenNode) (*TokenNode, error) {
	if err := batchCheckOperType([]*TokenNode{x, y}, operTypeMod, "%"); err != nil {
		return nil, err
	}

	if getInt(y) == 0 {
		return nil, GetError(ErrRuleEngineDivideByZero, "divide by zero")
	}

	return GetTokenNode(ValueTypeInteger, getInt(x)%getInt(y)), nil
}

func (o *TokenOperator) tokenNodeMinus(t *TokenNode) (*TokenNode, error) {
	if err := checkOperType(t, operTypeMinus, "-"); err != nil {
		return nil, err
	}
	res := &TokenNode{ValueType: t.ValueType}
	switch t.ValueType {
	case ValueTypeInteger:
		res.Value = -getInt(t)
	case ValueTypeFloat:
		if o.decimalMode {
			res.Value = getDecimal(t).Neg()
			res.ValueType = ValueTypeDecimal
		} else {
			res.Value = -getFloat(t)
		}
	case ValueTypeDecimal:
		res.Value = getDecimal(t).Neg()
	}
	return res, nil
}

func (o *TokenOperator) tokenNodeGreater(x, y *TokenNode) (*TokenNode, error) {
	if err := batchCheckOperType([]*TokenNode{x, y}, operTypeRelation, ">"); err != nil {
		return nil, err
	}

	res := &TokenNode{
		ValueType: ValueTypeBool,
	}

	if x.ValueType == ValueTypeInteger && y.ValueType == ValueTypeInteger {
		// integer
		res.Value = getInt(x) > getInt(y)
	} else if o.decimalMode || x.ValueType == ValueTypeDecimal || y.ValueType == ValueTypeDecimal {
		// decimal
		res.Value = getDecimal(x).GreaterThan(getDecimal(y))
		return res, nil
	} else {
		// float
		res.Value = getFloat(x) > getFloat(y)
	}

	return res, nil
}

func (o *TokenOperator) tokenNodeLess(x, y *TokenNode) (*TokenNode, error) {
	if err := batchCheckOperType([]*TokenNode{x, y}, operTypeRelation, "<"); err != nil {
		return nil, err
	}

	res := &TokenNode{
		ValueType: ValueTypeBool,
	}

	if x.ValueType == ValueTypeInteger && y.ValueType == ValueTypeInteger {
		// integer
		res.Value = getInt(x) < getInt(y)
	} else if o.decimalMode || x.ValueType == ValueTypeDecimal || y.ValueType == ValueTypeDecimal {
		// decimal
		res.Value = getDecimal(x).LessThan(getDecimal(y))
		return res, nil
	} else {
		// float
		res.Value = getFloat(x) < getFloat(y)
	}

	return res, nil
}

func (o *TokenOperator) tokenNodeGreaterEqual(x, y *TokenNode) (*TokenNode, error) {
	if err := batchCheckOperType([]*TokenNode{x, y}, operTypeRelation, ">="); err != nil {
		return nil, err
	}

	res := &TokenNode{
		ValueType: ValueTypeBool,
	}

	if x.ValueType == ValueTypeInteger && y.ValueType == ValueTypeInteger {
		// integer
		res.Value = getInt(x) >= getInt(y)
	} else if o.decimalMode || x.ValueType == ValueTypeDecimal || y.ValueType == ValueTypeDecimal {
		// decimal
		res.Value = getDecimal(x).GreaterThanOrEqual(getDecimal(y))
		return res, nil
	} else {
		// float
		res.Value = getFloat(x) >= getFloat(y)
	}

	return res, nil
}

func (o *TokenOperator) tokenNodeLessEqual(x, y *TokenNode) (*TokenNode, error) {
	if err := batchCheckOperType([]*TokenNode{x, y}, operTypeRelation, "<="); err != nil {
		return nil, err
	}

	res := &TokenNode{
		ValueType: ValueTypeBool,
	}

	if x.ValueType == ValueTypeInteger && y.ValueType == ValueTypeInteger {
		// integer
		res.Value = getInt(x) <= getInt(y)
	} else if o.decimalMode || x.ValueType == ValueTypeDecimal || y.ValueType == ValueTypeDecimal {
		// decimal
		res.Value = getDecimal(x).LessThanOrEqual(getDecimal(y))
		return res, nil
	} else {
		// float
		res.Value = getFloat(x) <= getFloat(y)
	}

	return res, nil
}

func (o *TokenOperator) tokenNodeEqual(x, y *TokenNode) (*TokenNode, error) {
	if err := batchCheckOperType([]*TokenNode{x, y}, operTypeEqual, "=="); err != nil {
		return nil, err
	}

	res := &TokenNode{ValueType: ValueTypeBool}

	if x.ValueType == ValueTypeBool || y.ValueType == ValueTypeBool {
		err := batchCheckFieldType([]*TokenNode{x, y}, []ValueType{ValueTypeBool})
		if err != nil {
			err.(*EngineErr).ErrMsg = "invalid equal operation for bool value with other type"
			return nil, err
		}
		res.Value = getBool(x) == getBool(y)
		return res, nil
	}

	if x.ValueType == ValueTypeString || y.ValueType == ValueTypeString {
		err := batchCheckFieldType([]*TokenNode{x, y}, []ValueType{ValueTypeString})
		if err != nil {
			err.(*EngineErr).ErrMsg = "invalid equal operation for string value with other type"
			return nil, err
		}
		res.Value = getString(x) == getString(y)
		return res, nil
	}

	if x.ValueType == ValueTypeInteger && y.ValueType == ValueTypeInteger {
		// integer
		res.Value = getInt(x) == getInt(y)
	} else if o.decimalMode || x.ValueType == ValueTypeDecimal || y.ValueType == ValueTypeDecimal {
		// decimal
		res.Value = getDecimal(x).Equal(getDecimal(y))
		return res, nil
	} else {
		// float
		res.Value = isFloatEqual(getFloat(x), getFloat(y))
	}

	return res, nil
}

func (o *TokenOperator) tokenNodeNotEqual(x, y *TokenNode) (*TokenNode, error) {
	if err := batchCheckOperType([]*TokenNode{x, y}, operTypeEqual, "!="); err != nil {
		return nil, err
	}

	res, err := o.tokenNodeEqual(x, y)
	if err == nil {
		res.Value = !getBool(res)
	}
	return res, err
}

func (o *TokenOperator) tokenNodeAnd(x, y *TokenNode) (*TokenNode, error) {
	if err := batchCheckOperType([]*TokenNode{x, y}, operTypeLogic, "and"); err != nil {
		return nil, err
	}

	return GetTokenNode(ValueTypeBool, getBool(x) && getBool(y)), nil
}

func (o *TokenOperator) tokenNodeOr(x, y *TokenNode) (*TokenNode, error) {
	if err := batchCheckOperType([]*TokenNode{x, y}, operTypeLogic, "or"); err != nil {
		return nil, err
	}

	return GetTokenNode(ValueTypeBool, getBool(x) || getBool(y)), nil
}

func (o *TokenOperator) tokenNodeNot(t *TokenNode) (*TokenNode, error) {
	if err := checkOperType(t, operTypeLogic, "not"); err != nil {
		return nil, err
	}

	return GetTokenNode(ValueTypeBool, !getBool(t)), nil
}

func (o *TokenOperator) tokenNodeVarName(varNameToken *TokenNode, identifier *TokenNode) (*TokenNode, error) {
	res := &TokenNode{ValueType: ValueTypeString}
	switch identifier.ValueType {
	case ValueTypeString:
		res.Value = fmt.Sprintf("%v.%v", getString(varNameToken), getString(identifier))
	case ValueTypeInteger:
		res.Value = fmt.Sprintf("%v.%v", getString(varNameToken), getInt(identifier))
	default:
		return nil, GetError(ErrRuleEngineSyntaxError, fmt.Sprintf("syntax error, var: %v", identifier.Value))
	}
	return res, nil
}

func (o *TokenOperator) tokenNodeVar(t *TokenNode) (*TokenNode, error) {
	varName := t.Value.(string)

	unknownVarErr := GetError(ErrRuleEngineUnknownVarName, fmt.Sprintf("unknown var name: %v", varName))

	if o.paramMap == nil {
		return nil, unknownVarErr
	}

	param, ok := o.paramMap[varName]
	if !ok {
		return nil, unknownVarErr
	}

	return &TokenNode{ValueType: param.Type, Value: param.Value}, nil
}

func (o *TokenOperator) tokenNodeThirdOper(x *TokenNode, c *TokenNode, y *TokenNode) (*TokenNode, error) {
	// like python third operation
	// x if c else y, c must bool value
	err := checkFiledType(c, []ValueType{ValueTypeBool})
	if err != nil {
		strValueType := valueTypeNameDict[c.ValueType]
		err.(*EngineErr).ErrMsg = fmt.Sprintf("if else condition type must bool value, but give :%v", strValueType)
		return nil, err
	}

	condition := getBool(c)

	if condition {
		return x, nil
	} else {
		return y, nil
	}
}
