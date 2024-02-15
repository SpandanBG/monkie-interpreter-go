package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"sudocoding.xyz/interpreter_in_go/src/lexer"
	"sudocoding.xyz/interpreter_in_go/src/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()

		if line == "exit" {
			fmt.Println("Byeee")
			return
		}

		l := lexer.New_V2(strings.NewReader(line))

		for tok := l.NextToken_V2(); tok.Type != token.EOF; tok = l.NextToken_V2() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
