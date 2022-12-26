package rule_engine

import (
	"fmt"
)

type TokenOperator struct {
	decimalMode bool
	varMap      map[string]*TokenNode
}

func (o *TokenOperator) tokenNodeAdd(x, y *TokenNode) (*TokenNode, error) {
	if err := batchCheckOperType([]*TokenNode{x, y}, operTypeMath, "+"); err != nil {
		return nil, err
	}

	if x.ValueType == ValueTypeInteger && y.ValueType == ValueTypeInteger {
		return GetTokenNode(ValueTypeInteger, x.GetInt()+y.GetInt()), nil
	}

	if o.decimalMode || x.ValueType == ValueTypeDecimal || y.ValueType == ValueTypeDecimal {
		res := x.GetDecimal().Add(y.GetDecimal())
		return GetTokenNode(ValueTypeDecimal, res), nil
	}

	return GetTokenNode(ValueTypeFloat, x.GetFloat()+y.GetFloat()), nil
}

func (o *TokenOperator) tokenNodeSub(x, y *TokenNode) (*TokenNode, error) {
	if err := batchCheckOperType([]*TokenNode{x, y}, operTypeMath, "-"); err != nil {
		return nil, err
	}

	if x.ValueType == ValueTypeInteger && y.ValueType == ValueTypeInteger {
		return GetTokenNode(ValueTypeInteger, x.GetInt()-y.GetInt()), nil
	}

	if o.decimalMode || x.ValueType == ValueTypeDecimal || y.ValueType == ValueTypeDecimal {
		res := x.GetDecimal().Sub(y.GetDecimal())
		return GetTokenNode(ValueTypeDecimal, res), nil
	}

	return GetTokenNode(ValueTypeFloat, x.GetFloat()-y.GetFloat()), nil
}

func (o *TokenOperator) tokenNodeMul(x, y *TokenNode) (*TokenNode, error) {
	if err := batchCheckOperType([]*TokenNode{x, y}, operTypeMath, "*"); err != nil {
		return nil, err
	}

	if x.ValueType == ValueTypeInteger && y.ValueType == ValueTypeInteger {
		return GetTokenNode(ValueTypeInteger, x.GetInt()*y.GetInt()), nil
	}

	if o.decimalMode || x.ValueType == ValueTypeDecimal || y.ValueType == ValueTypeDecimal {
		res := x.GetDecimal().Mul(y.GetDecimal())
		return GetTokenNode(ValueTypeDecimal, res), nil
	}

	return GetTokenNode(ValueTypeFloat, x.GetFloat()*y.GetFloat()), nil
}

func (o *TokenOperator) tokenNodeDiv(x, y *TokenNode) (*TokenNode, error) {
	if err := batchCheckOperType([]*TokenNode{x, y}, operTypeMath, "/"); err != nil {
		return nil, err
	}

	if y.GetDecimal().IsZero() {
		return nil, GetError(ErrRuleEngineDivideByZero, "divide by zero")
	}

	if x.ValueType == ValueTypeInteger && y.ValueType == ValueTypeInteger {
		return GetTokenNode(ValueTypeInteger, x.GetInt()/y.GetInt()), nil
	}

	if o.decimalMode || x.ValueType == ValueTypeDecimal || y.ValueType == ValueTypeDecimal {
		res := x.GetDecimal().Div(y.GetDecimal())
		return GetTokenNode(ValueTypeDecimal, res), nil
	}

	return GetTokenNode(ValueTypeFloat, x.GetFloat()/y.GetFloat()), nil
}

func (o *TokenOperator) tokenNodeMod(x, y *TokenNode) (*TokenNode, error) {
	if err := batchCheckOperType([]*TokenNode{x, y}, operTypeMod, "%"); err != nil {
		return nil, err
	}

	if y.GetInt() == 0 {
		return nil, GetError(ErrRuleEngineDivideByZero, "divide by zero")
	}

	return GetTokenNode(ValueTypeInteger, x.GetInt()%y.GetInt()), nil
}

func (o *TokenOperator) tokenNodeMinus(t *TokenNode) (*TokenNode, error) {
	if err := checkOperType(t, operTypeMinus, "-"); err != nil {
		return nil, err
	}
	res := &TokenNode{ValueType: t.ValueType}
	switch t.ValueType {
	case ValueTypeInteger:
		res.Value = -t.GetInt()
	case ValueTypeFloat:
		if o.decimalMode {
			res.Value = t.GetDecimal().Neg()
			res.ValueType = ValueTypeDecimal
		} else {
			res.Value = -t.GetFloat()
		}
	case ValueTypeDecimal:
		res.Value = t.GetDecimal().Neg()
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
		res.Value = x.GetInt() > y.GetInt()
	} else if o.decimalMode || x.ValueType == ValueTypeDecimal || y.ValueType == ValueTypeDecimal {
		// decimal
		res.Value = x.GetDecimal().GreaterThan(y.GetDecimal())
		return res, nil
	} else {
		// float
		res.Value = x.GetFloat() > y.GetFloat()
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
		res.Value = x.GetInt() < y.GetInt()
	} else if o.decimalMode || x.ValueType == ValueTypeDecimal || y.ValueType == ValueTypeDecimal {
		// decimal
		res.Value = x.GetDecimal().LessThan(y.GetDecimal())
		return res, nil
	} else {
		// float
		res.Value = x.GetFloat() < y.GetFloat()
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
		res.Value = x.GetInt() >= y.GetInt()
	} else if o.decimalMode || x.ValueType == ValueTypeDecimal || y.ValueType == ValueTypeDecimal {
		// decimal
		res.Value = x.GetDecimal().GreaterThanOrEqual(y.GetDecimal())
		return res, nil
	} else {
		// float
		res.Value = x.GetFloat() >= y.GetFloat()
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
		res.Value = x.GetInt() <= y.GetInt()
	} else if o.decimalMode || x.ValueType == ValueTypeDecimal || y.ValueType == ValueTypeDecimal {
		// decimal
		res.Value = x.GetDecimal().LessThanOrEqual(y.GetDecimal())
		return res, nil
	} else {
		// float
		res.Value = x.GetFloat() <= y.GetFloat()
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
		res.Value = x.GetBool() == y.GetBool()
		return res, nil
	}

	if x.ValueType == ValueTypeString || y.ValueType == ValueTypeString {
		err := batchCheckFieldType([]*TokenNode{x, y}, []ValueType{ValueTypeString})
		if err != nil {
			err.(*EngineErr).ErrMsg = "invalid equal operation for string value with other type"
			return nil, err
		}
		res.Value = x.GetString() == y.GetString()
		return res, nil
	}

	if x.ValueType == ValueTypeInteger && y.ValueType == ValueTypeInteger {
		// integer
		res.Value = x.GetInt() == y.GetInt()
	} else if o.decimalMode || x.ValueType == ValueTypeDecimal || y.ValueType == ValueTypeDecimal {
		// decimal
		res.Value = x.GetDecimal().Equal(y.GetDecimal())
		return res, nil
	} else {
		// float
		res.Value = isFloatEqual(x.GetFloat(), y.GetFloat())
	}

	return res, nil
}

func (o *TokenOperator) tokenNodeNotEqual(x, y *TokenNode) (*TokenNode, error) {
	if err := batchCheckOperType([]*TokenNode{x, y}, operTypeEqual, "!="); err != nil {
		return nil, err
	}

	res, err := o.tokenNodeEqual(x, y)
	if err == nil {
		res.Value = !res.GetBool()
	}
	return res, err
}

func (o *TokenOperator) tokenNodeAnd(x, y *TokenNode) (*TokenNode, error) {
	if err := batchCheckOperType([]*TokenNode{x, y}, operTypeLogic, "and"); err != nil {
		return nil, err
	}

	return GetTokenNode(ValueTypeBool, x.GetBool() && y.GetBool()), nil
}

func (o *TokenOperator) tokenNodeOr(x, y *TokenNode) (*TokenNode, error) {
	if err := batchCheckOperType([]*TokenNode{x, y}, operTypeLogic, "or"); err != nil {
		return nil, err
	}

	return GetTokenNode(ValueTypeBool, x.GetBool() || y.GetBool()), nil
}

func (o *TokenOperator) tokenNodeNot(t *TokenNode) (*TokenNode, error) {
	if err := checkOperType(t, operTypeLogic, "not"); err != nil {
		return nil, err
	}

	return GetTokenNode(ValueTypeBool, !t.GetBool()), nil
}

func (o *TokenOperator) tokenNodeVarName(varNameToken *TokenNode, identifier *TokenNode) (*TokenNode, error) {
	res := &TokenNode{ValueType: ValueTypeString}
	switch identifier.ValueType {
	case ValueTypeString:
		res.Value = fmt.Sprintf("%v.%v", varNameToken.Value, identifier.Value)
	case ValueTypeInteger:
		res.Value = fmt.Sprintf("%v.%v", varNameToken.Value, identifier.Value)
	default:
		return nil, GetError(ErrRuleEngineSyntaxError, fmt.Sprintf("syntax error, var: %v", identifier.Value))
	}
	return res, nil
}

func (o *TokenOperator) tokenNodeVar(t *TokenNode) (*TokenNode, error) {
	varName := t.Value.(string)

	unknownVarErr := GetError(ErrRuleEngineUnknownVarName, fmt.Sprintf("unknown var name: %v", varName))

	if o.varMap == nil {
		return nil, unknownVarErr
	}

	variable, ok := o.varMap[varName]
	if !ok {
		return nil, unknownVarErr
	}

	return GetTokenNode(variable.ValueType, variable.Value), nil
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

	condition := c.GetBool()

	if condition {
		return x, nil
	} else {
		return y, nil
	}
}
