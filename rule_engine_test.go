package rule_engine

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
)

type CheckUnit struct {
	str     string
	res     interface{}
	errCode int
}

type RuleEngineTest struct {
	t      testing.TB
	praser Praser
}

func GetNewRuleEngineTest(t testing.TB, params []*Param, useDecimal bool) (*RuleEngineTest, error) {
	praser, err := GetNewPraser(params, useDecimal)
	if err != nil {
		return nil, err
	}
	return &RuleEngineTest{
		t:      t,
		praser: *praser,
	}, nil
}

const DEBUG = false

func (rt *RuleEngineTest) check(checkCase *CheckUnit) {
	res, err := rt.praser.Parse(checkCase.str)
	if err != nil && err.(*engineErr).errCode != checkCase.errCode {
		rt.t.Fatalf("check errcode failed, input: %v, res_err: %v, expect_err_code: %v",
			checkCase.str, err, checkCase.errCode)
		return
	}

	if err != nil {
		if DEBUG {
			fmt.Printf("pass: %v\n err: %v\n", checkCase.str, err)
		}
		return
	}
	if !res.CheckValue(checkCase.res) {
		rt.t.Fatalf("check res value failed, input: %v, res_value: %v, expect_value: %v",
			checkCase.str, res.GetValue(), checkCase.res)
		return
	}
	if DEBUG {
		fmt.Printf("pass: %v\n res: %v\n", checkCase.str, res.GetValue())
	}
}

func (rt *RuleEngineTest) batchCheck(checkList *[]CheckUnit) {
	for _, checkCase := range *checkList {
		rt.check(&checkCase)
	}
}

func TestRuleEngineString(t *testing.T) {
	checkList := []CheckUnit{
		{`"test"`, "test", 0},
		{`'test'`, "test", 0},
		{`"'test'"`, `'test'`, 0},
		{`"'te\'st'"`, `'te\'st'`, 0},
		{`"te\"st"`, `te\"st`, 0},
	}
	rt, err := GetNewRuleEngineTest(t, nil, false)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	rt.batchCheck(&checkList)
}

func TestFuleEngineInteger(t *testing.T) {
	checkList := []CheckUnit{
		{"1 + 1", int64(2), 0},
		{"10 - 5", int64(5), 0},
		{"3 * 3* 3", int64(27), 0},
		{"1 + 8 / 4", int64(3), 0},
		{"10 % 3", int64(1), 0},
		{"9 / 4", int64(2), 0},

		{"3 + 2 * 2", int64(7), 0},
		{"3 - 2 % 3", int64(1), 0},

		{"2 * (3 + 2)", int64(10), 0},

		{"2 + -2", int64(0), 0},
		{"2 - -2", int64(4), 0},
		{"2 % -3", int64(2), 0},
		{"-2 % 3", int64(-2), 0},

		{"0x100", int64(16 * 16), 0},
		{"0x100 / 8", int64(32), 0},
		{"0x100 % 3", int64(1), 0},
		{"2 * 2 + 2 - 2 * 2 * 2", int64(-2), 0},
		{"010 * 2", int64(16), 0},
		{"0x10 * 4", int64(64), 0},
		{"010 * 010", int64(64), 0},
		{"010 * 010 - 0x10", int64(48), 0},
		{"010 * 010 + 0x10 * 4", int64(128), 0},
		{"010 * 010 - 0x10 * 4", int64(0), 0},

		{"1 / 0", 0, int(ErrRuleEngineDivideByZero)},
		{"1 % 0", 0, int(ErrRuleEngineDivideByZero)},

		{"1 + 2 - ", 0, int(ErrRuleEngineSyntaxError)},

		{`10 - 3 + 28`, 35, 0},
		{`10 + 3 - 28`, -15, 0},
	}

	rt, err := GetNewRuleEngineTest(t, nil, false)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	rt.batchCheck(&checkList)
}

func TestRuleEngineFloat(t *testing.T) {
	checkList := []CheckUnit{
		{"1.0 + 3.3", 4.3, 0},
		{"1.0 - 3.3", -2.3, 0},
		{"3.0 * 3.3", 9.9, 0},
		{"3.3 / 3.0", 1.1, 0},

		{"0.8 + 5", 5.8, 0},
		{"1 + 0.5", 1.5, 0},
		{"1 - 0.5", 0.5, 0},
		{"0.5 - 1", -0.5, 0},
		{"0.5 * 3", 1.5, 0},
		{"8 * 0.8", 6.4, 0},
		{"2.5 / 3", 2.5 / 3, 0},
		{"3 / 0.3", 10, 0},

		{"3 - 2 / 5.0", 2.6, 0},
		{"1 + 2.5 * 3", 8.5, 0},
		{"3.0 * (2.5 - 3)", -1.5, 0},
		{"3.0 * 2.5 - 3", 4.5, 0},
		{"2.1892 * 100", 218.92, 0},
		{"27413748 / 1000.0", 27413.748, 0},

		{"-3 / 0.3", -10, 0},
		{"-5.0 - -0.3 * 10", -2, 0},

		{"8.8 / 0", 0, int(ErrRuleEngineDivideByZero)},
		{"1.1 % 2", 0, int(ErrRuleEngineNotSupportedOperator)},
		{"1 % 2.1", 0, int(ErrRuleEngineNotSupportedOperator)},
	}

	rt, err := GetNewRuleEngineTest(t, nil, false)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	rt.batchCheck(&checkList)
}

func TestRuleEngineDecimal(t *testing.T) {
	checkList := []CheckUnit{
		{"1.0 + 3.3", 4.3, 0},
		{"1.0 - 3.3", -2.3, 0},
		{"3.0 * 3.3", 9.9, 0},
		{"3.3 / 3.0", 1.1, 0},

		{"0.8 + 5", 5.8, 0},
		{"1 + 0.5", 1.5, 0},
		{"1 - 0.5", 0.5, 0},
		{"0.5 - 1", -0.5, 0},
		{"0.5 * 3", 1.5, 0},
		{"8 * 0.8", 6.4, 0},
		{"2.5 / 3", decimal.NewFromFloat(2.5).Div(decimal.NewFromFloat(3)), 0},
		{"3 / 0.3", 10, 0},

		{"3 - 2 / 5.0", 2.6, 0},
		{"1 + 2.5 * 3", 8.5, 0},
		{"3.0 * (2.5 - 3)", -1.5, 0},
		{"3.0 * 2.5 - 3", 4.5, 0},
		{"2.1892 * 100", 218.92, 0},
		{"27413748 / 1000.0", 27413.748, 0},

		{"-3 / 0.3", -10, 0},
		{"-5.0 - -0.3 * 10", -2, 0},

		{"8.8 / 0", 0, int(ErrRuleEngineDivideByZero)},
		{"1.1 % 2", 0, int(ErrRuleEngineNotSupportedOperator)},
		{"1 % 2.1", 0, int(ErrRuleEngineNotSupportedOperator)},
	}

	// check decimal
	rt, err := GetNewRuleEngineTest(t, nil, true)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	rt.batchCheck(&checkList)
}

func TestRuleEngineRelation(t *testing.T) {
	checkList := []CheckUnit{
		{"100 * 3 == 300", true, 0},
		{"5.6 - 1.6 == 4", true, 0},
		{"5.6 - 1.6 > 4", false, 0},
		{"5.6  + -1.6 > 4", false, 0},
		{"5.6  - 10 > -4.4", false, 0},
		{"5.6  - 10 < -4.4", false, 0},
		{"5.6  - 10 == -4.4", true, 0},
		{"5.6  - 10 >= -4.4", true, 0},
		{"5.6  - 10 <= -4.4", true, 0},
		{"0x10 == 16", true, 0},
		{"0x10 == 020", true, 0},
		{"true > false", 0, int(ErrRuleEngineNotSupportedOperator)},
		{"1 + 1 == 2 == true", true, 0},
		{" 3.0 / 3  == 6.6 / 6.6 ", true, 0},
		{" 0x10 > 0x30 ", false, 0},
		{" 1 != 2 ", true, 0},
	}

	rt, err := GetNewRuleEngineTest(t, nil, false)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	rt.batchCheck(&checkList)

	rt, err = GetNewRuleEngineTest(t, nil, true)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	rt.batchCheck(&checkList)
}

func TestRuleEngineLogic(t *testing.T) {
	checkList := []CheckUnit{
		{"1 > 2 or 1 < 2", true, 0},
		{"1 > 2 || 1 < 2", true, 0},
		{"1 > 2 and 1 < 2", false, 0},
		{"1 > 2 && 1 < 2", false, 0},
		{"1 and true", 0, int(ErrRuleEngineNotSupportedOperator)},
		{"1 && true", 0, int(ErrRuleEngineNotSupportedOperator)},
		{"true AND false", false, 0},
		{"true OR false", true, 0},
		{"1 < 2 and 2 < 3 and 4.0 != 4.0001", true, 0},
		{"1 < 2 or 2.0 < 3 or 4.00001 > 4", true, 0},
		{"1 < 2 or 3.0 > 2.5 and 2.5 != 2.4", true, 0},
		{"not(1 > 2) and not (2.5 < 2.0)", true, 0},
		{"not true == false", true, 0},
		{"not true == true", false, 0},
		{"not (not true) == false", false, 0},
	}

	rt, err := GetNewRuleEngineTest(t, nil, false)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	rt.batchCheck(&checkList)

	rt, err = GetNewRuleEngineTest(t, nil, true)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	rt.batchCheck(&checkList)
}

func TestRuleEngineFunc(t *testing.T) {
	checkList := []CheckUnit{
		{`len("test")`, int64(4), 0},
		{`len("test") == 4`, true, 0},
		{`len("test") == len('test')`, true, 0},
		{`len("test") - len('test2') == -1`, true, 0},
		{`max(1,2,5.5,4,8)`, 8.0, 0},
		{`min(1,2,5.5,4,8)`, 1.0, 0},
		{`max(min(10.0, 20, 30), len("test"))`, 10, 0},
		{`min(min(10.0, 20, 30), len("test"))`, 4, 0},
		{`min(1,2,5.5,4,8) == 1`, true, 0},
		{`abs(-1)`, int64(1), 0},
		{`abs(-5.5)`, 5.5, 0},
		{`abs(3.0 - 100)`, 97, 0},
		{`max(len("test"), abs(-3), min(5,6))`, int64(5), 0},
		{`min(len("test"), abs(-4.5), min(5,6))`, 4, 0},

		{`upper("abc") == "ABC"`, true, 0},
		{`lower("CN")`, "cn", 0},
		{`len("HelloWorld") <= 5`, false, 0},
		{`len("HelloWorld")`, int64(10), 0},

		{`test()`, 0, int(ErrRuleEngineUnkonwnFunc)},
		{`min()`, 0, int(ErrRuleEngineFuncArgument)},
		{`max(1)`, 0, int(ErrRuleEngineFuncArgument)},

		{`startWith("hello world", "hell")`, true, 0},
		{`startWith("hello world", "ell")`, false, 0},

		{`endWith("hello world", "world")`, true, 0},
		{`startWith("hello world", "worl")`, false, 0},

		{`int(3.33) == 3`, true, 0},
		{`int("1878") - 8 == 187 * 10`, true, 0},
		{`int(97) + 3 == max(100, -1)`, true, 0},
		{`int(decimal(3.33)) -3 == 0`, true, 0},

		{`float(300)`, 300.0, 0},
		{`float("3.33") * 100 == int(1000 / 3.0)`, true, 0},
		{`float(decimal(-3.33)) * 100 == -(1000 / 3)`, true, 0},

		{`decimal(200) / 3`, decimal.NewFromFloat(200).Div(decimal.NewFromInt(3)), 0},
		{`decimal("3.33")`, 3.33, 0},
		{`decimal(3.33) * 10 - 0.3 == int(100.0 / 3)`, true, 0},
		{`int(decimal(3.33) * 10) == int(100.0 / 3)`, true, 0},

		{`string(decimal("3.33"))`, "3.33", 0},
		{`string(100.234) == "100.234"`, true, 0},
		{`decimal(string(3.33))`, 3.33, 0},
		{`len(string(int(100.0 / 3)))`, 2, 0},
	}

	rt, err := GetNewRuleEngineTest(t, nil, false)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	rt.batchCheck(&checkList)

	rt, err = GetNewRuleEngineTest(t, nil, true)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	rt.batchCheck(&checkList)
}

func TestRuleEngineVar(t *testing.T) {
	params := []*Param{
		GetParamWithType("value", ValueTypeInteger, 100),
		GetParam("int_value2", 33),
		GetParamWithType("float_value", ValueTypeFloat, 5.5),
		GetParam("float_value2", 3.33),
		GetParamWithType("str_value", ValueTypeString, "test"),
		GetParam("true_value", true),
		GetParam("false_value", false),
		GetParam("a.b.c", 10),
		GetParamWithType("a.1", ValueTypeInteger, int64(1)),
		GetParamWithType("d1", ValueTypeDecimal, decimal.NewFromInt(100)),
		GetParam("d2", decimal.NewFromFloat(3.33)),
		GetParamWithType("d3", ValueTypeDecimal, "3.33"),
	}

	checkList := []CheckUnit{
		{`len({{str_value}})`, int64(4), 0},
		{`len({{str_value}}) - {{value}}`, int64(-96), 0},
		{`len({{str_value}}) - {{value}} == -96 and {{true_value}} == true`, true, 0},
		{`max(len({{str_value}}), {{value}}, {{float_value}})`, 100, 0},
		{`min(len({{str_value}}), {{value}}, {{float_value}})`, 4, 0},
		{`{{true_value}} or {{false_value}}`, true, 0},
		{`{{true_value}} and {{false_value}}`, false, 0},
		{`{{a.b.c}}`, int64(10), 0},
		{`{{.b.c}}`, 0, int(ErrRuleEngineSyntaxError)},
		{`{{a.b.}}`, 0, int(ErrRuleEngineSyntaxError)},
		{`{{unknown_name}}`, 0, int(ErrRuleEngineUnknownVarName)},
		{`{{a.1}}`, int64(1), 0},
		{`{{d1}} * {{d2}}`, 333, 0},
		{`{{d1}} * {{d2}} - 3.0 == {{int_value2}} * {{value}} / 10`, true, 0},
		{`{{d1}} * {{d2}} - 3 == {{float_value2}} * {{d1}} - ({{int_value2}} - 3)/ 10.0`, true, 0},
		{`{{d3}} == {{d2}}`, true, 0},
	}

	rt, err := GetNewRuleEngineTest(t, params, false)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	rt.batchCheck(&checkList)

	rt, err = GetNewRuleEngineTest(t, params, true)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	rt.batchCheck(&checkList)
}

func TestRuleEngineRegexMatch(t *testing.T) {
	checkList := []CheckUnit{
		{`regexMatch("^test$", "test")`, true, 0},
		{`regexMatch("^test$", "test0")`, false, 0},
		{`regexMatch("^0[xX][a-fA-F0-9]+", "0xasf4")`, true, 0},
		{`regexMatch("foo.*", "seafood")`, true, 0},
		{`regexMatch("(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]",
				"https://www.baidu.com")`, true, 0},
		{`regexMatch("^\s*(DEV|TEST|UAT|STAGING|STABLE|LIVE)\s*$", " TEST ")`, true, 0},
		{`regexMatch("^\s*http(s)?://.*shopee\.com", " https://test.shopee.com")`, true, 0},
	}

	rt, err := GetNewRuleEngineTest(t, nil, false)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	rt.batchCheck(&checkList)

	rt, err = GetNewRuleEngineTest(t, nil, true)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	rt.batchCheck(&checkList)
}

func TestRuleEngineIfElse(t *testing.T) {
	params := []*Param{
		GetParamWithType("x", ValueTypeInteger, int64(100)),
		GetParamWithType("y", ValueTypeFloat, 50),
		GetParamWithType("COUNTRY", ValueTypeString, "CN"),
		GetParamWithType("field1", ValueTypeInteger, 2),
		GetParamWithType("field2", ValueTypeFloat, 7.2),
		GetParamWithType("field3", ValueTypeString, "Hel"),
	}

	checkList := []CheckUnit{
		{`1 if true else 2`, int64(1), 0},
		{`1 if false else 2`, int64(2), 0},
		{`{{x}} + 10 >= 105.3 if {{y}} == 50 else false`, true, 0},
		{`{{COUNTRY}} != "CN" or {{x}} >= 100`, true, 0},
		{`{{x}} >= 100 if {{COUNTRY}} == "CN" else true`, true, 0},
		{`{{y}} >= 100 if {{COUNTRY}} == "CN" else true`, false, 0},
		{`1 if 1 else 2`, 0, int(ErrRuleEngineInvalidOperation)},
		{"(({{field1}} > 0) and ({{field2}} > 7.8)) if len({{field3}}) >= 5 else {{field1}} + {{field2}} > 2.6", true, 0},
	}

	rt, err := GetNewRuleEngineTest(t, params, false)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	rt.batchCheck(&checkList)

	rt, err = GetNewRuleEngineTest(t, params, false)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	rt.batchCheck(&checkList)
}

func TestRuleEngineExapmle(t *testing.T) {
	checkList := []CheckUnit{
		{`3 * (5 - 2) + 1`, int64(10), 0},
	}

	rt, err := GetNewRuleEngineTest(t, nil, false)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	rt.batchCheck(&checkList)
}

func BenchmarkRule(b *testing.B) {
	checkCase := CheckUnit{
		"(({{field1}} > 0) and ({{field2}} > 7.8)) if len({{field3}}) >= 5 else {{field1}} + {{field2}} > 2.6",
		true, 0,
	}

	params := []*Param{
		GetParamWithType("x", ValueTypeInteger, int64(100)),
		GetParamWithType("y", ValueTypeInteger, 50),
		GetParamWithType("COUNTRY", ValueTypeString, "CN"),
		GetParamWithType("field1", ValueTypeInteger, int64(2)),
		GetParamWithType("field2", ValueTypeFloat, 7.2),
		GetParamWithType("field3", ValueTypeString, "Hel"),
	}
	rt, err := GetNewRuleEngineTest(b, params, false)
	if err != nil {
		b.Fatalf("%v\n", err)
	}

	for i := 0; i < b.N; i++ {
		rt.check(&checkCase)
	}
}
