package example

import subone "github.com/vieolo/shiraz/example/sub_one"

func FnOne() string {
	return subone.FnFive()
}

func FnTwo() string {
	return "two"
}
