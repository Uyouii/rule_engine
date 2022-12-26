package rule_engine

import "fmt"

const (
	Success = iota
	ErrRuleEngineFuncArgument
	ErrRuleEngineUnkonwnFunc
	ErrRuleEngineRegexMatch
	ErrRuleEngineNotSupportedVarType
	ErrRuleEngineSyntaxError
	ErrRuleEngineDivideByZero
	ErrRuleEngineUnknownVarName
	ErrRuleEngineInvalidVarType
	ErrRuleEngineUnknownOperator
	ErrRuleEngineNotSupportedOperator
	ErrRuleEngineInvalidOperation
	ErrRuleEngineInvalidParam
	ErrRuleEngineParamValueTypeNotMatch
	ErrRuleEngineDecimalError
)

var ERROR_MSG_MAP = map[int]string{
	Success:                             "ok",
	ErrRuleEngineFuncArgument:           "func args error",
	ErrRuleEngineUnkonwnFunc:            "unkonwn func",
	ErrRuleEngineRegexMatch:             "regex error",
	ErrRuleEngineNotSupportedVarType:    "not supported type",
	ErrRuleEngineSyntaxError:            "syntax error",
	ErrRuleEngineDivideByZero:           "divide by zero",
	ErrRuleEngineUnknownVarName:         "unknown variable name",
	ErrRuleEngineInvalidVarType:         "invalid variable type",
	ErrRuleEngineUnknownOperator:        "unknown operator",
	ErrRuleEngineNotSupportedOperator:   "not supported operator",
	ErrRuleEngineInvalidOperation:       "invalid operation",
	ErrRuleEngineInvalidParam:           "invalid parameter",
	ErrRuleEngineParamValueTypeNotMatch: "parameter value type not match",
	ErrRuleEngineDecimalError:           "error handle decimal",
}

type engineErr struct {
	errCode int
	errMsg  string
}

func (t *engineErr) Error() string {
	return fmt.Sprintf("[err]: code %v, %v, [err_msg]: %v", t.errCode, ERROR_MSG_MAP[t.errCode], t.errMsg)
}

func GetError(code int, msg string) *engineErr {
	return &engineErr{errCode: code, errMsg: msg}
}
