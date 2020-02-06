package main

import (
	"fmt"

	"github.com/piyapongch/era-audit/greet"
)

func main() {
	greet.Hello()

	fmt.Println(greet.Shark)

	oct := greet.Octopus{
		Name:  "Jesse",
		Color: "orange",
	}

	fmt.Println(oct.String())

	oct.Reset()

	fmt.Println(oct.​String())
}
