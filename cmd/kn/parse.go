package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
)

type LineType int64

const (
	Title LineType = iota
	Subtitle
	Page
	Empty
	Note
)

type Line struct {
	text  string
	type_ LineType
}

var lineTypeRegex = map[LineType]regexp.Regexp{
	Title:    *regexp.MustCompile(`^#\s.*?$`),
	Subtitle: *regexp.MustCompile(`^##\s.*?$`),
	Page:     *regexp.MustCompile(`^\d+?$`),
	Empty:    *regexp.MustCompile(`^\s*?$`),
	Note:     *regexp.MustCompile(`^.*?$`),
}

func ParseToYaml(file *os.File) {
	lines := parseLines(file)
	fmt.Print(lines)
	// TODO: Create yaml.
}

func parseLines(file *os.File) *[]Line {
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	var text string
	lines := []Line{}

	for fileScanner.Scan() {
		text = fileScanner.Text()

		l, err := getLine(text)
		if err != nil {
			log.Fatalln("cannot parse given file")
		}

		lines = append(lines, l)
	}

	return &lines
}

func getLine(text string) (Line, error) {
	fmt.Println(text)
	for type_, rgx := range lineTypeRegex {
		if rgx.MatchString(text) == true {
			return Line{text: text, type_: type_}, nil
		}
	}

	return Line{text: "", type_: Empty}, errors.New("unknown line type")
}
