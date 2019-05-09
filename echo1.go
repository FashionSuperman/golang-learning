package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	var s, temp string
	for i := 0; i < len(os.Args); i++ {
		s += temp + os.Args[i]
		temp = " "
	}
	fmt.Println(s)

	s = ""
	temp = ""
	//or
	for _, arg := range os.Args[:] {
		s += temp + arg
		temp = " "
	}
	fmt.Println(s)

	//or
	fmt.Println(strings.Join(os.Args, " "))
}
