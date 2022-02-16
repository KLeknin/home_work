package main

import (
	"flag"
	"fmt"
)

var (
	from, to      string
	limit, offset int64
)

func check(e error) {
	if e != nil {
		fmt.Printf("Error: %s\n", e.Error())
	}
}

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	// Place your code here.
	println(from, to, offset, limit)

	err := Copy(from, to, offset, limit)
	check(err)
}
