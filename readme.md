# Rule Engine By Go
用go-yacc实现的规则引擎，

### exapmle

```python
{{fieldA}} >= 10 and
{{filedA}} <= ({{fieldB}} + 5) * 2.5 and
{{fieldA}} % 1000 == 0

len({{fieldS}}) <= 40 and
regexMatch("^\s*http(s)?://.*shopee\.com", {{fieldS}})

{{fildC}} >= 100 if upper({{ENV}}) == "LIVE" else true

upper({{ENV}}) != "LIVE" or {{fildC}} >= 100
```



支持运算符：

```python
类型:      bool string int float
外部传值:   {{var_name}}
括号:      ( )
单目运算符: - ! not
算数运算符: + - * / %
比较运算符: > >= < <= == !=
逻辑运算符: and && or ||
函数操作:  len min max abs regexMatch lower upper
三目运算符: x if c else y
```



### BNF of ruleengine

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
    | MUL_EXPR '+' ADD_EXPR
    | MUL_EXPR '-' ADD_EXPR

MUL_EXPR :
    UNARY_EXPR
    | UNARY_EXPR '*' MUL_EXPR
    | UNARY_EXPR '/' MUL_EXPR
    | UNARY_EXPR '%' MUL_EXPR

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









