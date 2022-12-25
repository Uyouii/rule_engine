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
	ValueTypeDecimal
	valueTypeArgs
)

var valueTypeNameDict = map[ValueType]string{
	ValueTypeBool:    "bool",
	ValueTypeFloat:   "float",
	ValueTypeString:  "string",
	ValueTypeInteger: "integer",
	valueTypeArgs:    "args",
	ValueTypeDecimal: "decimal",
}

var valueTokenToValueType = map[int]ValueType{
	INTEGER:    ValueTypeInteger,
	FLOAT:      ValueTypeFloat,
	STRING:     ValueTypeString,
	TRUE:       ValueTypeBool,
	FALSE:      ValueTypeBool,
	BOOL:       ValueTypeBool,
	IDENTIFIER: ValueTypeString,
}

type operType int

const (
	operTypeMath operType = 1 + iota
	operTypeMod
	operTypeMinus
	operTypeRelation
	operTypeEqual
	operTypeLogic
	operTypeString
	operTypeArgument
	operTypeRegex
	operTypeChangeTo
)

var operValidType = map[operType][]ValueType{
	operTypeMath:     {ValueTypeInteger, ValueTypeFloat, ValueTypeDecimal},
	operTypeMod:      {ValueTypeInteger},
	operTypeMinus:    {ValueTypeInteger, ValueTypeFloat, ValueTypeDecimal},
	operTypeRelation: {ValueTypeInteger, ValueTypeFloat, ValueTypeDecimal},
	operTypeEqual:    {ValueTypeInteger, ValueTypeFloat, ValueTypeBool, ValueTypeString, ValueTypeDecimal},
	operTypeLogic:    {ValueTypeBool},
	operTypeString:   {ValueTypeString},
	operTypeArgument: {valueTypeArgs},
	operTypeRegex:    {ValueTypeString},
	operTypeChangeTo: {ValueTypeInteger, ValueTypeFloat, ValueTypeDecimal, ValueTypeString},
}
