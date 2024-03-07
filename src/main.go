package main

import (
	"flag"
	"fmt"
	"os"

	"sudocoding.xyz/interpreter_in_go/src/execute"
	"sudocoding.xyz/interpreter_in_go/src/repl"
)

var replIt = flag.Bool("repl", true, "run repl mode")
var exeFile = flag.String("exe", "", "execute file")

func main() {
	flag.Parse()

	if *exeFile != "" {
		execute.Execute(*exeFile)
		return
	}

	fmt.Println("Welcome to Monkie Lang!!")
	if *replIt {
		fmt.Println("Starting Repl. Type `exit` to quit ")
		repl.Start(os.Stdin, os.Stdout)
		return
	}
}
