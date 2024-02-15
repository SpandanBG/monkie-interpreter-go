package main

import (
	"flag"
	"fmt"
	"os"

	"sudocoding.xyz/interpreter_in_go/src/repl"
)

var replIt = flag.Bool("repl", true, "run repl mode")

func main() {
	flag.Parse()

	fmt.Println("Welcome to Monkie Lang!!")
	if *replIt {
		fmt.Println("Starting Repl. Type `exit` to quit ")
		repl.Start(os.Stdin, os.Stdout)
	}
}
