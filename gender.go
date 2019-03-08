package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/kako-jun/gender/gender-core"
)

func parseArgs() (keyword string, commandArgs []string, err error) {
	flag.Parse()
	if flag.NArg() < 1 {
		err = errors.New("invalid argument")
		return
	}

	args := flag.Args()
	keyword = args[0]
	return
}

// entry point
func main() {
	keyword, args, err := parseArgs()
	if err != nil {
		fmt.Println("error:", err)
		fmt.Println("usage:")
		fmt.Println("  gender [a search term]")
		return
	}

	gender.Exec(keyword, args)
}
