package example

import subexample "github.com/vieolo/shiraz/example/sub_example"

func FnOne() string {
	return subexample.FnFive()
}

func FnTwo() string {
	return "two"
}
