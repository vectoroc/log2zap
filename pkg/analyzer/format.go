package analyzer

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"strconv"
	"strings"
)

func expr2name(expr ast.Expr) (string, bool) {
	switch v := expr.(type) {
	case *ast.Ident:
		// use var key
		return v.String(), true

	case *ast.CallExpr:
		cs, err := callSelector(v)
		if err != nil {
			return "", false
		}

		return cs.String(), true
	default:
		return "", false
	}
}

func expr2string(pass *analysis.Pass, expr ast.Expr) (string, error) {
	buf := &bytes.Buffer{}
	err := format.Node(buf, pass.Fset, expr)
	return buf.String(), err
}

func formatZapVar(pass *analysis.Pass, arg ast.Expr, v formatVar) (string, error) {
	typ := pass.TypesInfo.TypeOf(arg)

	if typ.String() == "error" {
		val, err := expr2string(pass, arg)
		return `zap.Error(` + val + `)`, err
	}

	if ce, ok := arg.(*ast.CallExpr); ok && strings.Join(selector2string(ce.Fun), ".") == "err.Error" {
		return `zap.Error(err)`, nil
	}

	type2method := map[types.Type]string{
		types.Typ[types.Int]:     "Int",
		types.Typ[types.Int32]:   "Int32",
		types.Typ[types.Uint32]:  "Uint32",
		types.Typ[types.Float32]: "Float32",
		types.Typ[types.Float64]: "Float64",
		types.Typ[types.String]:  "String",
	}

	name := v.key
	method, ok := type2method[typ]

	if !ok {
		method = "Any"
	}

	if name == "" {
		name, _ = expr2name(arg)
	}

	val, err := expr2string(pass, arg)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`zap.%s(%s, %s)`, method, strconv.Quote(name), val), nil
}
