package rule_engine

import (
	"fmt"
	"regexp"
	"strconv"
	"unicode"

	"github.com/shopspring/decimal"
)

type RuleEngineLex struct {
	str     string
	pos     int
	err     *engineErr
	resNode *TokenNode
	oper    *TokenOperator
}

func NewRuleEngineLex(str string, oper *TokenOperator) *RuleEngineLex {
	return &RuleEngineLex{
		str:  str,
		oper: oper,
	}
}

func (lex *RuleEngineLex) setErr(err error) int {
	if err == nil {
		return Success
	}
	lex.err = err.(*engineErr)
	return int(lex.err.errCode)
}

func (lex *RuleEngineLex) matchRule(str string) (int, string) {
	token, resStr := 0, ""
	for _, tokenRule := range TOKEN_RULE_LIST {
		r, _ := regexp.Compile("^" + tokenRule.reStr)
		matchStr := r.FindString(str)
		if len(matchStr) != 0 {
			token, resStr = tokenRule.token, matchStr
			break
		}
	}

	// there is some key words match the identifier
	if token == IDENTIFIER {
		for _, tokenRule := range KEY_WORD_LIST {
			r, _ := regexp.Compile("^" + tokenRule.reStr)
			matchStr := r.FindString(resStr)
			if len(matchStr) != 0 && matchStr == resStr {
				token, resStr = tokenRule.token, matchStr
				break
			}
		}
	}

	return token, resStr
}

func (lex *RuleEngineLex) Lex(lval *ruleEngineSymType) int {
	for ; lex.pos < len(lex.str); lex.pos++ {
		c := rune(lex.str[lex.pos])
		if !unicode.IsSpace(c) {
			break
		}
	}

	if lex.pos >= len(lex.str) {
		return END
	}

	token, matchStr := lex.matchRule(lex.str[lex.pos:])

	var err error

	if len(matchStr) > 0 {
		lex.pos += len(matchStr)
		lval.node = &TokenNode{
			ValueType: valueTokenToValueType[token],
		}

		switch token {
		case STRING:
			lval.node.Value = matchStr[1 : len(matchStr)-1]
		case INTEGER:
			if lval.node.Value, err = strconv.ParseInt(matchStr, 0, 64); err != nil {
				return ERROR
			}
		case FLOAT:
			if lex.oper.decimalMode {
				if lval.node.Value, err = decimal.NewFromString(matchStr); err != nil {
					return ERROR
				}
				lval.node.ValueType = ValueTypeDecimal
			} else {
				if lval.node.Value, err = strconv.ParseFloat(matchStr, 64); err != nil {
					return ERROR
				}
			}
		case TRUE:
			lval.node.Value, token = true, BOOL
		case FALSE:
			lval.node.Value, token = false, BOOL
		case IDENTIFIER:
			lval.node.Value = matchStr
		}

		return token
	}

	c := rune(lex.str[lex.pos])
	if _, ok := VALID_CHAR_SET[c]; ok {
		lex.pos += 1
		return int(c) // 直接使用这个char
	}

	return ERROR // some thing wrong
}

func (lex *RuleEngineLex) Error(s string) {
	prefix := make([]byte, lex.pos)
	for i := 0; i < lex.pos; i++ {
		prefix[i] = ' '
	}
	lex.err = GetError(ErrRuleEngineSyntaxError, fmt.Sprintf("%v, pos: %v\n%v\n%v\n", s, lex.pos, lex.str, string(prefix)+"^"))
}

func (lex *RuleEngineLex) getErrCode() int {
	if lex.err != nil {
		return int(lex.err.errCode)
	}
	return 0
}
