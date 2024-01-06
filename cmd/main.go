package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type Line struct {
	Id      string
	TypeLOB string
}

var line1 Line = Line{Id: "firstId", TypeLOB: "LOBTypeA"}
var line2 Line = Line{Id: "secondId", TypeLOB: "LOBTypeB"}

type Lines []*Line // not a fan of this

var sampleLines Lines = Lines{&line1, &line2}

var dynamicLines Lines

type Policy struct {
	number string //GUID  (UGH!)
	lines  Lines
}

type PolicyOption func(*Policy)

func main() {
	aPolicy := NewPolicy("Policy1234",
		WithLines(sampleLines)) //Hard code value works however I trying to get values from csv.I can read them in loop but can't append
	//	GetLines(Lines))

	fmt.Println("Print Lines HardCoded") //Print Lines HardCoded
	for _, line := range aPolicy.lines {
		fmt.Println("Line ID -" + line.TypeLOB)
	}

}

func NewPolicy(number string, opts ...PolicyOption) *Policy {
	p := &Policy{
		number: number,
		lines:  make([]*Line, 0),
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func WithLines(lines Lines) PolicyOption {

	return func(p *Policy) {

		Lineslist := readLineObjectsCSVFile("Objects/Line.csv")
		for _, line := range Lineslist {
			//p.lines = append(p.lines, lines)
			fmt.Println("Print Dynamic")
			fmt.Println("Line ID -" + line.TypeLOB)
			//	p.lines = append(p.lines, line) //Can't make this work

		}
		for _, line := range lines {
			p.lines = append(p.lines, line)
		}
	}
}

func (PolicyOption) GetLines(lines Lines) PolicyOption {
	return func(p *Policy) {
		for _, line := range lines {
			p.lines = append(p.lines, line)
		}
		/*
			lines := readLineObjectsCSVFile("Objects/Line.csv")
			for _, lines := range lines {
				//p.lines = append(p.lines, lines)
				fmt.Println("Line ID -" + lines.TypeLOB)

			}*/ //The comment code works and brings desired output
		//	p.lines = append(p.lines, line) //Unable to make this work
	}
}

func readLineObjectsCSVFile(filePath string) (lines []Line) {
	isFirstRow := true
	headerMap := make(map[string]int)

	// Load a csv file.
	f, _ := os.Open(filePath)

	// Create a new reader.
	r := csv.NewReader(f)
	for {
		// Read row
		record, err := r.Read()
		defer f.Close()

		// Stop at EOF.
		if err == io.EOF {
			break
		}

		checkError("Some other error occurred", err)

		// Handle first row case
		if isFirstRow {
			isFirstRow = false

			// Add mapping: Column/property name --> record index
			for i, v := range record {
				headerMap[v] = i
			}

			// Skip next code
			continue
		}

		// Create new coverage and add to persons array
		lines = append(lines, Line{
			TypeLOB: record[headerMap["TypeLOB"]],
			Id:      record[headerMap["Id"]],
		})
	}
	return

}

func checkError(message string, err error) {
	// Error Logging
	if err != nil {
		//log.Fatal(message, err)
		//log := zerolog.New(os.Stdout).With().Logger()
		//log.Debug().Str(" checkError %s ", message)
	}
}
