package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/palats/glop/parser"
	"github.com/palats/glop/repl"
	"github.com/palats/glop/runtime"
)

var (
	historyFilename = flag.String(
		"history_filename", "~/.glophistory",
		"File to use for prompt history. '~/' is expanded to current user home "+
			"directory. Set to empty string to disable history loading/saving.")
	inputCode = flag.String(
		"e", "",
		"Execute the provided code. If equal to '-', execute code from stdin.")
)

func main() {
	// Check arguments
	flag.Parse()

	if len(flag.Args()) > 1 {
		fmt.Fprintf(os.Stderr, "Only one file argument allowed.\n")
		os.Exit(1)
	}

	if len(flag.Args()) > 0 && *inputCode != "" {
		fmt.Fprintf(os.Stderr, "A filename cannot be specified together with -e.\n")
		os.Exit(1)
	}

	// Run REPL if nothing else is specified
	if *inputCode == "" && len(flag.Args()) == 0 {
		if err := repl.REPL(*historyFilename); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

	// By default run the code from -e flag.
	var input = *inputCode

	if len(flag.Args()) > 0 {
		// Get code from a file if specified
		data, err := ioutil.ReadFile(flag.Args()[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to read %q: %v\n", flag.Args()[0], err)
			os.Exit(1)
		}
		input = string(data)
	} else if input == "-" {
		// Get code from stdin if asked.
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to read stdin: %v", err)
			os.Exit(1)
		}
		input = string(data)
	}

	// Parse & exec.
	root, err := parser.Parse(parser.NewSource(input))
	if err != nil {
		fmt.Fprintf(os.Stderr, "parsing failed: %v", err)
		os.Exit(1)
	}

	ctx := runtime.NewContext(nil)
	if _, err = ctx.TryEval(root); err != nil {
		fmt.Fprintf(os.Stderr, "execution failed: %v", err)
		os.Exit(1)
	}
}
