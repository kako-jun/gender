package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/kako-jun/gender/gender-core"
)

// app version
var Version string = "1.0.0"

func parseArgs() (string, bool, bool, bool, bool, bool, bool, bool, bool, bool, bool, bool, bool, error) {
	var (
		versionFlag bool
		exactFlag   bool
		closestFlag bool
		simpleFlag  bool
		jsonFlag    bool
		arFlag      bool
		frFlag      bool
		deFlag      bool
		hiFlag      bool
		itFlag      bool
		ptFlag      bool
		ruFlag      bool
		esFlag      bool
	)

	flag.BoolVar(&versionFlag, "version", false, "print version number")
	flag.BoolVar(&exactFlag, "exact", false, "exact match")
	flag.BoolVar(&closestFlag, "closest", false, "print only the closest candidate")
	flag.BoolVar(&simpleFlag, "simple", false, "hide language label")
	flag.BoolVar(&jsonFlag, "json", false, "print in JSON format")
	flag.BoolVar(&arFlag, "ar", false, "search in Arabic")
	flag.BoolVar(&frFlag, "fr", false, "search in French")
	flag.BoolVar(&deFlag, "de", false, "search in German")
	flag.BoolVar(&hiFlag, "hi", false, "search in Hindi")
	flag.BoolVar(&itFlag, "it", false, "search in Italian")
	flag.BoolVar(&ptFlag, "pt", false, "search in Portuguese")
	flag.BoolVar(&ruFlag, "ru", false, "search in Russian")
	flag.BoolVar(&esFlag, "es", false, "search in Spanish")

	flag.Parse()
	args := flag.Args()

	if versionFlag {
		fmt.Println(Version)
		os.Exit(0)
	}

	var err error
	if flag.NArg() < 1 {
		err = errors.New("invalid argument")
		// return
	}

	keyword := ""
	if flag.NArg() >= 1 {
		keyword = args[0]
	}

	return keyword, exactFlag, closestFlag, simpleFlag, jsonFlag, arFlag, frFlag, deFlag, hiFlag, itFlag, ptFlag, ruFlag, esFlag, err
}

// entry point
func main() {
	keyword, exactFlag, closestFlag, simpleFlag, jsonFlag, arFlag, frFlag, deFlag, hiFlag, itFlag, ptFlag, ruFlag, esFlag, err := parseArgs()
	if err != nil {
		fmt.Println("error:", err)
		fmt.Println("usage:")
		fmt.Println("  gender ([options]) [keyword]")
		return
	}

	gender.Exec(keyword, exactFlag, closestFlag, simpleFlag, jsonFlag, arFlag, frFlag, deFlag, hiFlag, itFlag, ptFlag, ruFlag, esFlag)
}
