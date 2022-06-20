package analyzer

import (
	"fmt"
	"go/ast"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "log2zap",
	Doc:      "Find std-like logger calls and make suggests to replace it with uber.zap",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspector.Preorder(nodeFilter, func(node ast.Node) {
		ce := node.(*ast.CallExpr)

		switch ce.Fun.(type) {
		case *ast.ArrayType, *ast.ParenExpr, *ast.IndexExpr:
			// possible not a logger call, ignore
			return
		}

		lc, err := callSelector(ce)
		if err != nil {
			raw, _ := expr2string(pass, ce)
			pass.ReportRangef(ce, "call: %v: %#v, %s", err, ce.Fun, raw)
			return
		}

		// ignore non-logger calls
		if !lc.isLoggerCall() {
			return
		}

		if !lc.isPrintf() {
			pass.Report(analysis.Diagnostic{
				Pos:     ce.Pos(),
				Message: "log2zap",
				SuggestedFixes: []analysis.SuggestedFix{{
					TextEdits: []analysis.TextEdit{{
						Pos:     ce.Fun.Pos(),
						End:     ce.Fun.End(),
						NewText: []byte(`zap.S().` + lc.zapLevel()),
					}},
				}},
			})
			return
		}

		lit, ok := ce.Args[0].(*ast.BasicLit)
		if !ok {
			raw, _ := expr2string(pass, ce.Args[0])
			pass.ReportRangef(ce, "expected 1st arg to be format string literal, got %s", raw)
			return
		}

		format := strings.Trim(lit.Value, `"`)
		vars := parseFormat(format)
		msg := cleanUpFormatString(format, vars)

		var zapArgs []string
		for i, arg := range ce.Args[1:] {
			v, err := formatZapVar(pass, arg, vars[i])
			if err != nil {
				pass.ReportRangef(arg, "failed to format arg")
				return
			}
			zapArgs = append(zapArgs, v)
		}

		pass.Report(analysis.Diagnostic{
			Pos:     ce.Pos(),
			Message: "log2zap",
			SuggestedFixes: []analysis.SuggestedFix{{
				Message: lc.zapLevel() + " call",
				TextEdits: []analysis.TextEdit{{
					Pos:     ce.Pos(),
					End:     ce.End(),
					NewText: []byte(fmt.Sprintf(`zap.L().%s(%s, %s)`, lc.zapLevel(), strconv.Quote(msg), strings.Join(zapArgs, ", "))),
				}},
			}},
			Related: nil,
		})
	})

	return nil, nil
}
