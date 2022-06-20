package warn

import (
	"flag"
	"fmt"
	"warn/trace"
	"io/ioutil"
	"log"
	"os"
)

type Logger interface {
	Println(...interface{})
	Printf(string, ...interface{})
	Output(int, string) error
}

var (
	warn      = log.New(os.Stderr, "E", log.LstdFlags|log.Lshortfile)
	flVerbose = flag.Bool("verbose", false, "enable debug output")
	flTrace   = flag.Bool("trace", false, "enable trace output")
)

// init is a function that sets prefix(D), output(stdout) and file & line to std logger.
func init() {
	log.SetPrefix("D")
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// Init is a function that discards std logger output if verbose flag was not set.
func Init() {
	if !*flVerbose {
		log.SetOutput(ioutil.Discard)
	}

	trace.Init(*flTrace)
}

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) { warn.Output(2, fmt.Sprintf(format, v...)) }

// Print calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) { warn.Output(2, fmt.Sprint(v...)) }

// Println calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) { warn.Output(2, fmt.Sprintln(v...)) }

// Fatal calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func Fatal(v ...interface{}) { warn.Output(2, fmt.Sprint(v...)); os.Exit(1) }

// IsVerbose returns true if `verbose` flag was set.
func IsVerbose() bool { return *flVerbose }

// IsTrace returns true if `trace` flag was set.
func IsTrace() bool { return *flTrace }

// Std returns standard Warn logger.
func Std() Logger { return warn }
