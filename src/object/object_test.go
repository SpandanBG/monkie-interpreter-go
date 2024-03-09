package object

import (
	"strings"
	"testing"
)

func eq[T comparable](t *testing.T, expected T, actual T, msg ...string) {
	if expected != actual {
		t.Fatalf("%s\nexpected: %+v\nactual: %+v\n", strings.Join(msg, " "), expected, actual)
	}
}

func notEq[T comparable](t *testing.T, expected T, actual T, msg ...string) {
	if expected == actual {
		t.Fatalf("%s\nexpected not: %+v\nactual: %+v\n", strings.Join(msg, " "), expected, actual)
	}
}

func Test_ObjectHashing(t *testing.T) {
	helloWorld_1 := &String{Value: "Hello World"}
	helloWorld_2 := &String{Value: "Hello World"}
	diff_1 := &String{Value: "Hola"}
	diff_2 := &String{Value: "Hola"}

	eq(t, helloWorld_1.Hash(), helloWorld_1.Hash())
	eq(t, diff_1.Hash(), diff_2.Hash())
	notEq(t, helloWorld_1.Hash(), diff_1.Hash())
	notEq(t, helloWorld_2.Hash(), diff_2.Hash())
}

func Test_HashFetching(t *testing.T) {
	name1 := &String{Value: "name"}
	monkie := &String{Value: "monkie"}

	pairs := map[HashKey]Object{}

	pairs[name1.Hash()] = monkie
	name2 := &String{Value: "name"}

	eq(t, monkie, pairs[name2.Hash()].(*String))
	eq(t, monkie.Inspect(), pairs[name2.Hash()].Inspect())
}
