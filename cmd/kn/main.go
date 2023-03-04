package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

func main() {
	parser := argparse.NewParser("kn", "a simple utility to parse, keep and read notes")

	cmdParse := parser.NewCommand("parse", "Parse note raw file into yaml")
	cmdParseFile := cmdParse.File("f", "file", os.O_RDWR, 0600, &argparse.Options{Required: true})

	var err = parser.Parse(os.Args)
	if err != nil {
		fmt.Println(parser.Usage(err))
		os.Exit(0)
	}

	switch {
	case cmdParse.Happened():
		ParseToYaml(cmdParseFile)
	}
}
