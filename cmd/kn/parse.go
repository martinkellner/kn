package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

type LineType int64

const (
	titleLine LineType = iota
	subtitleLine
	pageLine
	emptyLine
	noteLine
)

type Line struct {
	text  string
	type_ LineType
}

var lineTypeRegex = map[LineType]regexp.Regexp{
	titleLine:    *regexp.MustCompile(`^#\s.*?$`),
	subtitleLine: *regexp.MustCompile(`^##\s.*?$`),
	pageLine:     *regexp.MustCompile(`^\d+?$`),
	emptyLine:    *regexp.MustCompile(`^\s*?$`),
	noteLine:     *regexp.MustCompile(`^.*?$`),
}

type Note struct {
	Text     string
	Subtitle string
	Page     string
	Tag      []string
}

type NoteFile struct {
	Author string
	Title  string
	Notes  []Note
}

func ParseToYaml(file *os.File) {
	lines := parseLines(file)
	noteFile := createNoteFile(lines)

	noteYaml, err := yaml.Marshal(&noteFile)
	if err != nil {
		log.Fatal(err)
	}

	err2 := ioutil.WriteFile(fmt.Sprintf("%s.yaml", file.Name()), noteYaml, 0666)
	if err2 != nil {
		log.Fatal(err2)
	}

	file.Close()
}

func createNoteFile(lines []Line) *NoteFile {
	title := ""
	subtitle := ""
	page := ""
	notes := []Note{}

	for _, l := range lines {
		if l.type_ == titleLine {
			title = l.text
		} else if l.type_ == subtitleLine {
			subtitle = l.text
		} else if l.type_ == pageLine {
			page = l.text
		} else if l.type_ == noteLine {
			note := Note{Text: l.text, Subtitle: subtitle, Page: page, Tag: []string{"dummy"}}
			notes = append(notes, note)
		}
	}

	return &NoteFile{Author: "", Title: title, Notes: notes}
}

func parseLines(file *os.File) []Line {
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

	return lines
}

func getLine(text string) (Line, error) {
	for type_, rgx := range lineTypeRegex {
		if rgx.MatchString(text) == true {
			return Line{text: text, type_: type_}, nil
		}
	}

	return Line{text: "", type_: emptyLine}, errors.New("unknown line type")
}
