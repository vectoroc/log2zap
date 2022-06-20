package analyzer

import (
	"errors"
	"go/ast"
	"strings"
)

type call struct {
	pkg, method string
}

func (c *call) String() string {
	if c.pkg == "" {
		return c.method
	}
	return c.pkg + "." + c.method
}

func selector2string(expr ast.Expr) []string {
	if ident, ok := expr.(*ast.Ident); ok {
		return []string{ident.String()}
	}

	if sel, ok := expr.(*ast.SelectorExpr); ok {
		left := selector2string(sel.X)
		return append(left, sel.Sel.Name)
	}

	return nil
}

func callSelector(ce *ast.CallExpr) (*call, error) {
	if _, ok := ce.Fun.(*ast.FuncLit); ok {
		return &call{}, nil
	}

	sel := selector2string(ce.Fun)
	if len(sel) == 0 {
		return nil, errors.New("missing selector")
	}
	return &call{
		method: sel[len(sel)-1],
		pkg:    strings.Join(sel[:len(sel)-1], "."),
	}, nil
}

func (c *call) isLoggerCall() bool {
	switch c.pkg {
	case "log", "fmt", "warn", "trace":
	// ok
	default:
		return false
	}

	switch c.method {
	case "Print", "Println", "Printf", "Fatal", "Fatalln", "Fatalf", "Panic", "Panicln", "Panicf":
		return true
	default:
		return false
	}
}

func (c *call) isPrintf() bool {
	return strings.HasSuffix(c.method, "f")
}

func (c *call) zapLevel() string {
	switch c.method {
	case "Fatal", "Fatalln", "Fatalf":
		return "Fatal"
	case "Panic", "Panicln", "Panicf":
		return "Panic"
	}

	switch c.pkg {
	case "trace":
		return "Debug"
	case "warn":
		return "Error"
	}

	return "Info"
}
