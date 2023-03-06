package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"path/filepath"

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
	Text string
	Type LineType
}

var lineTypeRegex = map[LineType]regexp.Regexp{titleLine: *regexp.MustCompile(`^#\s.*?$`),
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
	if err := validateLines(lines); err != nil {
		log.Fatal(err)
	}

	noteFile := createNoteFile(lines)

	noteYaml, err2 := yaml.Marshal(&noteFile)
	if err2 != nil {
		log.Fatal(err2)
	}

	output := file.Name()
	output = output[0 : len(output)-len(filepath.Ext(output))]
	output = fmt.Sprintf("%s.yaml", output)
	if err := ioutil.WriteFile(output, noteYaml, 0666); err != nil {
		log.Fatal(err)
	}

	file.Close()
}

func createNoteFile(lines []Line) *NoteFile {
	title := ""
	subtitle := ""
	page := ""
	notes := []Note{}

	for _, l := range lines {
		switch l.Type {
		case titleLine:
			title = l.Text
		case subtitleLine:
			subtitle = l.Text
		case pageLine:
			page = l.Text
		case noteLine:
			note := Note{Text: l.Text, Subtitle: subtitle, Page: page, Tag: []string{"dummy"}}
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
			log.Fatal(err)
		}

		lines = append(lines, l)
	}

	return lines
}

func validateLines(lines []Line) error {
	currType := emptyLine
	prevType := emptyLine

	for _, l := range lines {
		currType = l.Type

		if currType == titleLine {
			if prevType != emptyLine {
				return errors.New("title starting with # has to be in the first line")
			}
		} else if currType == noteLine {
			if prevType != pageLine {
				return errors.New(fmt.Sprintf("note has to follow page number, note: %s", l.Text))
			}
		}

		if currType != emptyLine {
			prevType = currType
		}
	}

	return nil
}

func getLine(text string) (Line, error) {
	var lineTypesOrder = [...]LineType{titleLine, subtitleLine, pageLine, emptyLine, noteLine}
	for _, t := range lineTypesOrder {
		rgx := lineTypeRegex[t]

		if rgx.MatchString(text) == true {
			text = strings.TrimLeft(text, "#")
			text = strings.TrimSpace(text)

			return Line{Text: text, Type: t}, nil
		}
	}

	return Line{Text: "", Type: emptyLine}, errors.New("unknown line type")
}
