package p

import (
	"errors"

	"warn"
	"warn/trace"
)

func testCrimsonLoggerCalls() {
	x := 10
	err := errors.New("test err")
	warn.Printf("some error happend domain_id=%d err=%s", x, err) // want "log2zap"
	trace.Print("foo bar", x, "zzz") // want "log2zap"
}
