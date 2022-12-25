%{

package rule_engine

%}

// fields inside this union end up as the fields in a structure known
// as ${PREFIX}SymType, of which a reference is passed to the lexer.
%union{
	node *TokenNode
}

%type <node> VALUE_EXPR VAR_NAME
%type <node> PRIMARY_EXPR UNARY_EXPR POST_EXPR
%type <node> RELATION_EXPR EQUAL_EXPR
%type <node> LOGIC_OR_EXPR LOGIC_AND_EXPR LOGIC_EXPR
%type <node> ADD_EXPR MUL_EXPR
%type <node> TRANSLATION_UNIT
%type <node> ARGUMENT_EXPRSSION_LIST
%type <node> THIRD_OPER_EXPR


// same for terminals
%token <node> INTEGER FLOAT STRING
%token <node> IDLEFT IDRIGHT IDENTIFIER
%token <node> BOOL TRUE FALSE
%token <node> AND OR NOT
%token <node> LE GE EQ NE
%token <node> ERROR END
%token <node> IF ELSE

%left AND OR
%left '>' '<' LE GE EQ NE
%left '+'  '-'
%left '*'  '/'  '%'

%%

top :
	TRANSLATION_UNIT {
		lex := ruleEnginelex.(*RuleEngineLex)
		lex.resNode = $1
		return 0
	}

TRANSLATION_UNIT :
	LOGIC_EXPR END {
		$$ = $1
	}

LOGIC_EXPR :
	LOGIC_OR_EXPR {
		$$ = $1
	}

LOGIC_OR_EXPR :
	LOGIC_AND_EXPR {
		$$ = $1
	}
	| LOGIC_OR_EXPR OR LOGIC_AND_EXPR {
		lex := ruleEnginelex.(*RuleEngineLex)
		node, err := lex.oper.tokenNodeOr($1, $3)
		if err != nil {
			return lex.setErr(err)
		}
		$$ = node
	}

LOGIC_AND_EXPR :
	THIRD_OPER_EXPR {
		$$ = $1
	}
	| LOGIC_AND_EXPR AND THIRD_OPER_EXPR {
		lex := ruleEnginelex.(*RuleEngineLex)
		node, err := lex.oper.tokenNodeAnd($1, $3)
		if err != nil {
			return lex.setErr(err)
		}
		$$ = node
	}

THIRD_OPER_EXPR :
	EQUAL_EXPR {
		$$ = $1
	}
	| EQUAL_EXPR IF THIRD_OPER_EXPR ELSE THIRD_OPER_EXPR {
		lex := ruleEnginelex.(*RuleEngineLex)
		node, err := lex.oper.tokenNodeThirdOper($1, $3, $5)
		if err != nil {
			return lex.setErr(err)
		}
		$$ = node
	}

EQUAL_EXPR :
	RELATION_EXPR {
		$$ = $1
	}
	| EQUAL_EXPR EQ RELATION_EXPR {
		lex := ruleEnginelex.(*RuleEngineLex)
		node, err := lex.oper.tokenNodeEqual($1, $3)
		if err != nil {
			return lex.setErr(err)
		}
		$$ = node
	}
	| EQUAL_EXPR NE RELATION_EXPR {
		lex := ruleEnginelex.(*RuleEngineLex)
		node, err := lex.oper.tokenNodeNotEqual($1, $3)
		if err != nil {
			return lex.setErr(err)
		}
		$$ = node
	}

RELATION_EXPR :
	ADD_EXPR {
		$$  =  $1
	}
	| ADD_EXPR '<' RELATION_EXPR {
		lex := ruleEnginelex.(*RuleEngineLex)
		node, err := lex.oper.tokenNodeLess($1, $3)
		if err != nil {
			return lex.setErr(err)
		}
		$$ = node
	}
	| ADD_EXPR '>' RELATION_EXPR {
		lex := ruleEnginelex.(*RuleEngineLex)
		node, err := lex.oper.tokenNodeGreater($1, $3)
		if err != nil {
			return lex.setErr(err)
		}
		$$ = node
	}
	| ADD_EXPR LE RELATION_EXPR {
		lex := ruleEnginelex.(*RuleEngineLex)
		node, err := lex.oper.tokenNodeLessEqual($1, $3)
		// __yyfmt__.Println($1, $3)
		if err != nil {
			return lex.setErr(err)
		}
		$$ = node
	}
	| ADD_EXPR GE RELATION_EXPR {
		lex := ruleEnginelex.(*RuleEngineLex)
		node, err := lex.oper.tokenNodeGreaterEqual($1, $3)
		if err != nil {
			return lex.setErr(err)
		}
		$$ = node
	}


ADD_EXPR :
	MUL_EXPR {
		$$  =  $1
	}
	| MUL_EXPR '+' ADD_EXPR {
		lex := ruleEnginelex.(*RuleEngineLex)
		node, err := lex.oper.tokenNodeAdd($1, $3)
		if err != nil {
			return lex.setErr(err)
		}
		$$ = node
	}
	| MUL_EXPR '-' ADD_EXPR {
		lex := ruleEnginelex.(*RuleEngineLex)
		node, err := lex.oper.tokenNodeSub($1, $3)
		if err != nil {
			return lex.setErr(err)
		}
		$$ = node
	}

MUL_EXPR :
	UNARY_EXPR {
		$$ = $1
	}
	| UNARY_EXPR '*' MUL_EXPR {
		lex := ruleEnginelex.(*RuleEngineLex)
		node, err := lex.oper.tokenNodeMul($1, $3)
		if err != nil {
			return lex.setErr(err)
		}
		$$ = node
	}
	| UNARY_EXPR '/' MUL_EXPR {
		lex := ruleEnginelex.(*RuleEngineLex)
		node, err := lex.oper.tokenNodeDiv($1, $3)
		if err != nil {
			return lex.setErr(err)
		}
		$$ = node
	}
	| UNARY_EXPR '%' MUL_EXPR {
		lex := ruleEnginelex.(*RuleEngineLex)
		node, err := lex.oper.tokenNodeMod($1, $3)
		if err != nil {
			return lex.setErr(err)
		}
		$$ = node
	}

UNARY_EXPR :
	POST_EXPR {
		$$ = $1
	}
	| '-' PRIMARY_EXPR {
		lex := ruleEnginelex.(*RuleEngineLex)
		node, err := lex.oper.tokenNodeMinus($2)
		if err != nil {
			return lex.setErr(err)
		}
		$$ = node
	}
	| NOT PRIMARY_EXPR {
		lex := ruleEnginelex.(*RuleEngineLex)
		node, err := lex.oper.tokenNodeNot($2)
		if err != nil {
			return lex.setErr(err)
		}
		$$ = node
	}

POST_EXPR :
	PRIMARY_EXPR {
		$$ = $1
	}
	| IDENTIFIER '(' ARGUMENT_EXPRSSION_LIST ')' {
		lex := ruleEnginelex.(*RuleEngineLex)
		node, err := lex.oper.tokenNodeFunc($1, $3)
		if err != nil {
			return lex.setErr(err)
		}
		// __yyfmt__.Println(*node)
		$$ = node
	}
	| IDENTIFIER '(' ')' {
		lex := ruleEnginelex.(*RuleEngineLex)
		node, err := lex.oper.tokenNodeFunc($1, nil);
		if err != nil {
			return lex.setErr(err)
		}
		$$ = node
	}

ARGUMENT_EXPRSSION_LIST :
	LOGIC_EXPR {
		lex := ruleEnginelex.(*RuleEngineLex)
		node, err := lex.oper.tokenNodeArg($1)
		if err != nil {
			return lex.setErr(err)
		}
		$$ = node
	}
	| ARGUMENT_EXPRSSION_LIST ',' LOGIC_EXPR {
		lex := ruleEnginelex.(*RuleEngineLex)
		node, err := lex.oper.tokenNodeArgList($1, $3);
		if err != nil {
			return lex.setErr(err)
		}
		$$ = node
	}


PRIMARY_EXPR :
	INTEGER {
		// __yyfmt__.Println($1)
		$$ = $1
	}
	| FLOAT {
		$$ = $1
	}
	| BOOL {
		$$ = $1
	}
	| STRING {
		$$ = $1
	}
	| ERROR {
		ruleEnginelex.Error("syntax error")
		return ruleEnginelex.(*RuleEngineLex).getErrCode()
	}
	| '(' LOGIC_EXPR ')' {
		$$  =  $2
	}
	| VALUE_EXPR {
		$$ = $1
	}

VALUE_EXPR :
	IDLEFT VAR_NAME IDRIGHT {
		// __yyfmt__.Println($2)
		lex := ruleEnginelex.(*RuleEngineLex)
		node, err := lex.oper.tokenNodeVar($2)
		if err != nil {
			return lex.setErr(err)
		}
		$$ = node
	}

VAR_NAME :
	IDENTIFIER {
		$$ = $1
	}
	| VAR_NAME '.' IDENTIFIER {
		lex := ruleEnginelex.(*RuleEngineLex)
		node, err := lex.oper.tokenNodeVarName($1, $3);
		if err != nil {
			return lex.setErr(err)
		}
		$$ = node
	}
	| VAR_NAME '.' INTEGER {
		// __yyfmt__.Println($1, $3)
		lex := ruleEnginelex.(*RuleEngineLex)
		node, err := lex.oper.tokenNodeVarName($1, $3)
		if err != nil {
			return lex.setErr(err)
		}
		$$ = node
	}


%%      /*  start  of  programs  */