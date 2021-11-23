package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	// Вывод перевернутой фразы: "Hello, OTUS!"
	var revercedString = stringutil.Reverse("Hello, OTUS!")
	fmt.Println(revercedString)
}
