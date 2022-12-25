package rule_engine

import (
	"fmt"
)

var VALID_CHAR_SET = map[rune]struct{}{
	'(': {},
	')': {},
	'+': {},
	'-': {},
	'*': {},
	'/': {},
	'%': {},
	'>': {},
	'<': {},
	',': {},
	'.': {},
}

const L = `[a-zA-Z_]`
const H = `[a-fA-F0-9]`
const E = `([Ee][+-]?[0-9]+)`
const P = `([Pp][+-]?[0-9]+)`
const FS = `(f|F|l|L)`
const IS = `((u|U)|(u|U)?(l|L|ll|LL)|(l|L|ll|LL)(u|U))`

type tokenRule struct {
	token int
	reStr string
}

var TOKEN_RULE_LIST = [...]tokenRule{
	{LE, "<="},
	{GE, ">="},
	{EQ, "=="},
	{NE, "!="},
	{AND, "&&"},
	{NOT, "!"},
	{OR, `\|\|`},
	{IDLEFT, "{{"},
	{IDRIGHT, "}}"},
	{STRING, `\"(\\.|[^\\"\n])*\"`},
	{STRING, `\'(\\.|[^\\'\n])*\'`},
	{FLOAT, fmt.Sprintf(`[0-9]+%v%v?`, E, FS)},
	{FLOAT, fmt.Sprintf(`[0-9]+\.[0-9]+%v?%v?`, E, FS)},
	{FLOAT, fmt.Sprintf(`[0-9]+\.[0-9]*%v?%v?`, E, FS)},
	{INTEGER, fmt.Sprintf(`0[xX]%v+%v?`, H, IS)},
	{INTEGER, fmt.Sprintf(`0[0-7]*%v?`, IS)},
	{INTEGER, fmt.Sprintf(`[1-9][0-9]*%v?`, IS)},
	{IDENTIFIER, fmt.Sprintf(`%v(%v|[0-9])*`, L, L)},
}

var KEY_WORD_LIST = [...]tokenRule{
	{AND, "AND"},
	{AND, "[A|a]nd"},
	{OR, "OR"},
	{OR, "[o|O]r"},
	{NOT, "[N|n]ot"},
	{NOT, "NOT"},
	{TRUE, "TRUE"},
	{TRUE, "[T|t]rue"},
	{FALSE, "FALSE"},
	{FALSE, "[F|f]alse"},
	{IF, "[I|i]f"},
	{IF, "IF"},
	{ELSE, "[E|e]lse"},
	{ELSE, "ELSE"},
}

type ValueType int

const (
	ValueTypeNone ValueType = iota
	ValueTypeInteger
	ValueTypeFloat
	ValueTypeBool
	ValueTypeString
	ValueTypeArgs
	ValueTypeDecimal
)

var VALUE_TOKEN_TO_VALUE_TYPE = map[int]ValueType{
	INTEGER:    ValueTypeInteger,
	FLOAT:      ValueTypeFloat,
	STRING:     ValueTypeString,
	TRUE:       ValueTypeBool,
	FALSE:      ValueTypeBool,
	BOOL:       ValueTypeBool,
	IDENTIFIER: ValueTypeString,
}

type OPER_TYPE int

const (
	OPER_TYPE_MATH OPER_TYPE = 1 + iota
	OPER_TYPE_MOD
	OPER_TYPE_MINUS
	OPER_TYPE_RELATION
	OPER_TYPE_EQUAL
	OPER_TYPE_LOGIC
	OPER_TYPE_STRING
	OPER_TYPE_ARGUMENT
	OPER_TYPE_REGEX
	OPER_TYPE_CHANGE_TO
)

var OPER_VALID_TYPE = map[OPER_TYPE][]ValueType{
	OPER_TYPE_MATH:      {ValueTypeInteger, ValueTypeFloat, ValueTypeDecimal},
	OPER_TYPE_MOD:       {ValueTypeInteger},
	OPER_TYPE_MINUS:     {ValueTypeInteger, ValueTypeFloat, ValueTypeDecimal},
	OPER_TYPE_RELATION:  {ValueTypeInteger, ValueTypeFloat, ValueTypeDecimal},
	OPER_TYPE_EQUAL:     {ValueTypeInteger, ValueTypeFloat, ValueTypeBool, ValueTypeString, ValueTypeDecimal},
	OPER_TYPE_LOGIC:     {ValueTypeBool},
	OPER_TYPE_STRING:    {ValueTypeString},
	OPER_TYPE_ARGUMENT:  {ValueTypeArgs},
	OPER_TYPE_REGEX:     {ValueTypeString},
	OPER_TYPE_CHANGE_TO: {ValueTypeInteger, ValueTypeFloat, ValueTypeDecimal, ValueTypeString},
}

var VALUE_TYPE_NAME_DICT = map[ValueType]string{
	ValueTypeBool:    "bool",
	ValueTypeFloat:   "float",
	ValueTypeString:  "string",
	ValueTypeInteger: "integer",
	ValueTypeArgs:    "args",
	ValueTypeDecimal: "decimal",
}
