package evaluator

import (
	"fmt"
	"testing"

	"sudocoding.xyz/interpreter_in_go/src/object"
)

func Test_Quote(t *testing.T) {
	for _, test := range []struct {
		input    string
		expected string
	}{
		{`quote(5)`, `5`},
		{`quote(5 + 8)`, `(5 + 8)`},
		{`quote(foobar)`, `foobar`},
		{`quote(foobar + barfoo)`, `(foobar + barfoo)`},
	} {
		t.Run(fmt.Sprintf("Test Quote for %s", test.input), func(t *testing.T) {
			evaluated := testEval(test.input)

			quote, ok := evaluated.(*object.Quote)
			eq(t, true, ok, fmt.Sprintf("Failed to typecase obj of type %s to *object.Quote", quote.Type()))
			notEq(t, nil, &quote.Node, "Expected non nil node")
			eq(t, test.expected, quote.Node.String(), "Expected node string didn't match")
		})
	}
}
