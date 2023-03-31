package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	str := "Hello, OTUS!"
	reversedString := stringutil.Reverse(str)
	fmt.Println(reversedString)
}
