# Rule Engine By Go

通过go-yacc实现的规则引擎，解析算式并计算得到结果。支持输入自定义参数，算数运算，逻辑运算，deciaml和一些内置函数（后续会根据需要扩展）。


## 安装

```sh
go get github.com/uyouii/rule_engine
```

## 要求

rule_engine library requires Go version `>=1.7`

## 示例

示例项目： [rule_engine_example](https://github.com/Uyouii/rule_engine_example)

```go
package main

import (
	"fmt"

	"github.com/uyouii/rule_engine"
)

func main() {

	params := []*rule_engine.Param{
		rule_engine.GetParam("i", 100),
		rule_engine.GetParam("f", 3.5),
		rule_engine.GetParam("s", "hello world"),
		rule_engine.GetParam("b", false),
		rule_engine.GetParam("d", 3.3),
		// get decimal from string
		rule_engine.GetParamWithType("d2", rule_engine.ValueTypeDecimal, "3.3"),
	}

	// use decimal: true
	praser, err := rule_engine.GetNewPraser(params, true)
	if err != nil {
		panic(err)
	}

	exampleList := []string{
		// integrate
		`4 * (2 + 3) - 5 * 3`,
		// float
		`3.0 * (2.5 - 3)`,
		// logic
		`1 < 2 and 2 < 3 and 4.0 != 4.0001`,
		`1 > 2 && 1 < 2`,
		// func exapmle
		`max(min(10.0, 20, 30), len("vstr"))`,
		`min(len("test"), abs(-4.5), min(5,6))`,
		`upper("abc") == "ABC"`,
		`startWith("hello world", "hel")`,
		`int(97) + 3 == max(100, -1)`,
		`regexMatch("^0[xX][a-fA-F0-9]+", "0xasf4")`,
		`string({{d}}) == "3.3"`,
		// use param
		`min(len({{s}}), {{i}}, {{f}}, {{d}})`,
		`int({{i}} / {{f}})`,
		`{{d}} * 10 - 3 - int({{i}} / {{f}})`,
		`{{d2}} - {{d}}`,
		// if else
		`{{d}} * 10 if len({{s}}) > 10 else {{f}} / 10`,
	}

	for _, example := range exampleList {
		res, err := praser.Parse(example)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%v  --> %v\n", example, res.Value)
	}
}

```

**输出**

```sh
4 * (2 + 3) - 5 * 3  --> 5
3.0 * (2.5 - 3)  --> -1.5
1 < 2 and 2 < 3 and 4.0 != 4.0001  --> true
1 > 2 && 1 < 2  --> false
max(min(10.0, 20, 30), len("vstr"))  --> 10
min(len("test"), abs(-4.5), min(5,6))  --> 4
upper("abc") == "ABC"  --> true
startWith("hello world", "hel")  --> true
int(97) + 3 == max(100, -1)  --> true
regexMatch("^0[xX][a-fA-F0-9]+", "0xasf4")  --> true
string({{d}}) == "3.3"  --> true
min(len({{s}}), {{i}}, {{f}}, {{d}})  --> 3.3
int({{i}} / {{f}})  --> 28
{{d}} * 10 - 3 - int({{i}} / {{f}})  --> 2
{{d2}} - {{d}}  --> 0
{{d}} * 10 if len({{s}}) > 10 else {{f}} / 10  --> 33
```

## go pkg 文档

https://pkg.go.dev/github.com/uyouii/rule_engine

## 接口

```go
// 1. use GetNewPraser to get a New Praser
// the params are the variable will be used in the calculation
// is set useDecimal, all the float in the param and calculation will be changed to decimal.
func GetNewPraser(params []*Param, useDecimal bool) (*Praser, error)

// 2. use Parse to get the result
func (p *Praser) Parse(str string) (*TokenNode, error)

// for example
praser, _ := rule_engine.GetNewPraser(params, true)
res, _ := praser.Parse(`1 + 1`)
fmt.Printf(res.Value)

2
```

#### 设置变量

Param用作向Praser传递变量：

```go
type Param struct {
	Name  string      // value name
	Type  ValueType   // value type
	Value interface{} // value
}
```

如果只设置 `Name`and `Value`,  `Praser` 会尝试从传入的`Value`中解析出`Type`.

```go
// for example
p1 := GetParamWithType("x", rule_engine.ValueTypeDecimal, "3.3")
variable {{x}} will be prase to decimal 3.3
```

如果同时设置 `Name` , `Value` 和 `Type`, `Praser` 会根据传入的`Type`重新解析 `Value` 。

```go
// for example
p2 := GetParam("y", "3.3")
variable {{y}} will be prase to string "3.3"
```

传入变量的使用可以参考`支持传入变量`一节。

### 获取结果

最终 `Parse`接口会返回 `TokenNode` 作为结果。

```go
type TokenNode struct {
	ValueType ValueType   // result type, can see ValueType
	Value     interface{} // result value
}
```

```go
// if want get interface{} res, can use
func(t *TokenNode) GetValue() interface{}

// if want get detail type value, can use
func (t *TokenNode) GetInt() int64
func (t *TokenNode) GetBool() bool
func (t *TokenNode) GetFloat() float64
func (t *TokenNode) GetDecimal() decimal.Decimal
func (t *TokenNode) GetString() string
```

## 功能实现

### 支持类型

| Type    | ValueType in Code |
| ------- | ----------------- |
| bool    | ValueTypeBool     |
| string  | ValueTypeString   |
| int     | ValueTypeInteger  |
| float   | ValueTypeFloat    |
| decimal | ValueTypeDecimal  |

> 注意：decimal类型相关的实现依赖  https://github.com/shopspring/decimal

#### 类型归约

在运算中算数类型的值会隐式归约为更加精确的类型。

如果 `int` 和 `float`一起计算, `int`会先转换为 `float`, 并且最终结果是 `float`.

如果 `int` 或者 `float` 和`decimal`一起计算, 都会被转换为`decimal`, 最终结果也是 `decimal`.

```go
int >> float >> decimal
```

例如:

```go
int + float = float
float * decimal = decimal
decimal - int = decimal
max(int, float, deciamal) = decimal
```

### 支持运算符

| Operator        | Name              | Support Types       |
| --------------- | ----------------- | ------------------- |
| `()`            | Parentheses       | ALL                 |
| `{{var_name}}`  | External Variable | ALL                 |
| `-`             | Negative          | int, float, decimal |
| `!` `not`       | Not               | bool                |
| `+`             | Addition          | int, float, decimal |
| `-`             | Subtraction       | int, float, decimal |
| `*`             | Multiplication    | int, float, decimal |
| `/`             | Division          | int, float, decimal |
| `%`             | Mod               | int                 |
| `>`             | Larger            | int, float, decimal |
| `>=`            | Larger or Equal   | int, float, decimal |
| `<`             | Less              | int, float, decimal |
| `<=`            | Less or Equal     | int, float, decimal |
| `==`            | Equal             | ALL                 |
| `!=`            | NotEqual          | ALL                 |
| `and` `&&`      | And               | bool                |
| `or` `\|\|`     | Or                | bool                |
| `x if c else y` | Ternary operator  | ALL, `c` must bool  |

#### 运算符优先级

优先级从上到下依次递减。

```go
() {{var_name}}
! not -(Negative)
* / %
+ -
> >= < <=
== !=
if else (ternary operator)
and && or ||
```

### 支持传入变量

`rule_engine` 支持在运算中传入外部变量, 通过使用 `{{}}`(双打括号) 包裹变量名表示。

传入的变量名和类型需要通过`Param`在创建`Praser`的时候传递。

使用示例:

```go
d: type decimal, value 3.3
i: type int, value 100

{{d}} * 100 + 3 - int({{i}} * 10 / 3)  --> 0
```

传入变量的类型可以是 `int`, `float`, `decimal`, `bool`, `string`

### 函数

#### 支持的内置函数列表

| Function Name | Descrption                          |
| ------------- | ----------------------------------- |
| len()         | length of the string                |
| min()         | min of the args                     |
| max()         | max of the args                     |
| abs()         | Abs                                 |
| upper()       | Upper of the string                 |
| lower()       | Lower of the string                 |
| startWith()   | check string start with some prefix |
| endWith()     | check string end with some suffix   |
| regexMatch()  | Regex Match                         |
| int()         | change arg to int type              |
| float()       | change arg to float type            |
| decimal()     | change arg to decimal type          |
| string()      | change arg to string type           |

#### len()

```go
// length of the string
// param {string} input string
// return {bool}
bool len(str)

e.g.
len("test")
4
```

#### min()

```go
// min of the args
// param {int/bool/decimal} 
// return {int/bool/decimal}, result type will be the most precise type.
any min(x, y, z, ...)

e.g.
min(1, len("test"), 6.8)
6.8
```

#### max()

```go
// max of the args
// param {int/bool/decimal} 
// return {int/bool/decimal}, result type will be the most precise type.
any max(x, y, z, ...)

e.g.
max(1, len("test"), 6.8)
6.8
```

#### abs()

```go
// Abs of the arg
// param {int/bool/decimal} x
// return {int/bool/decimal}, result type accornding to the input type
any abs(x)

e.g.
abs(-1.1)
1.1
```

#### upper()

```go
// the upper string of the input
// param {string} s
// return {string}
string upper(s string)

e.g.
upper("Hello World")
"HELLO WORLD"
```

#### lower()

```go
// the lower string of the input
// param {string} s
// return {string}
string lower(s string)

e.g.
lower("Hello World")
"hello world"
```

#### startWith()

```go
// check s start with prefix
// param {string} s
// param {string} prefix
// return {bool}
bool startWith(s string, prefix string)

e.g.
startWith("Hello World", "Hello")
true
```

#### endWith()

```go
// check s start with suffix
// param {string} s
// param {string} suffix
// return {bool}
bool endWith(s string, suffix string)

e.g.
startWith("Hello World", "World")
true
```

#### regexMatch()

```go
// whether the string s contains any match of the regular expression pattern.
// param {string} pattern, the regex pattern
// param {string} the check string
// return {bool}
bool regexMatch(pattern string, s string)

e.g.
regexMatch("^test$", "test")
true

regexMatch("(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]","https://www.baidu.com")
true
```

#### int()

```go
// change value to integer type
// param {int/float/decimal/string} v
// return {int}
int int(v any)

e.g.
int("100")
100

int(33.33)
33
```

#### float()

```go
// change value to float type
// param {int/float/decimal/string} v
// return {float}
float float(v any)

e.g.
float("33.3")
33.3

float(100)
100
```

#### decimal()

```go
// change value to decimal type
// param {int/float/decimal/string} v
// return {decimal}
decimal decimal(v any)

e.g.
decimal("33.3")
33.3

decimal(100)
100
```

#### string()

```go
// change value to string type
// param {int/float/decimal/string} v
// return {string}
string string(v any)

e.g.
string(33.3)
"33.3"

string(100)
"100"
```

### BNF 范式

`rule_engine`解析语法的BNF范式，描述了如何解析输入的字符串并且归约得到结果。

```go
top :
    TRANSLATION_UNIT

TRANSLATION_UNIT :
    LOGIC_EXPR END

LOGIC_EXPR :
    LOGIC_OR_EXPR

LOGIC_OR_EXPR :
    LOGIC_AND_EXPR
    | LOGIC_OR_EXPR OR LOGIC_AND_EXPR

LOGIC_AND_EXPR :
    THIRD_OPER_EXPR
    | LOGIC_AND_EXPR AND THIRD_OPER_EXPR

THIRD_OPER_EXPR :
    EQUAL_EXPR
    | EQUAL_EXPR IF THIRD_OPER_EXPR ELSE THIRD_OPER_EXPR

EQUAL_EXPR :
    RELATION_EXPR
    | EQUAL_EXPR EQ RELATION_EXPR
    | EQUAL_EXPR NE RELATION_EXPR

RELATION_EXPR :
    ADD_EXPR
    | ADD_EXPR '<' RELATION_EXPR
    | ADD_EXPR '>' RELATION_EXPR
    | ADD_EXPR LE RELATION_EXPR
    | ADD_EXPR GE RELATION_EXPR

ADD_EXPR :
    MUL_EXPR
    | ADD_EXPR '+' MUL_EXPR
    | ADD_EXPR '-' MUL_EXPR

MUL_EXPR :
    UNARY_EXPR
    | MUL_EXPR '*' UNARY_EXPR
    | MUL_EXPR '/' UNARY_EXPR
    | MUL_EXPR '%' UNARY_EXPR

UNARY_EXPR :
    POST_EXPR
    | '-' PRIMARY_EXPR
    | NOT PRIMARY_EXPR

POST_EXPR :
    PRIMARY_EXPR
    | IDENTIFIER '(' ARGUMENT_EXPRSSION_LIST ')'
    | IDENTIFIER '(' ')'

ARGUMENT_EXPRSSION_LIST :
    LOGIC_EXPR
    | ARGUMENT_EXPRSSION_LIST ',' LOGIC_EXPR

PRIMARY_EXPR :
    INTEGER
    | FLOAT
    | BOOL
    | STRING
    | ERROR
    | '(' LOGIC_EXPR ')'
    | VALUE_EXPR

VALUE_EXPR :
    IDLEFT VAR_NAME IDRIGHT

VAR_NAME :
    IDENTIFIER
    | VAR_NAME '.' IDENTIFIER
    | VAR_NAME '.' INTEGER
```

## FAQ

### Why need rule engine?



## TODO