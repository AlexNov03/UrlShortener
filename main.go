package main

import (
	"flag"
	"fmt"
)

func main() {
	inMemory := flag.Bool("in-memory", true, "defines, whether app use in-memory or postgres db")

	flag.Parse()

	if *inMemory == false {
		fmt.Println("app is using postgres db")
	} else {
		fmt.Println("app is using in-memory db")
	}
}
