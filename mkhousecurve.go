package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/getlantern/golog"
)

var (
	log = golog.LoggerFor("mkhousecurve")
)

var (
	in      = flag.String("in", "", "name of input file")
	out     = flag.String("out", "housecurve.txt", "name of output file")
	refDB   = flag.Float64("refdb", 84, "reference level in dB")
	comment = flag.String("comment", "", "comment to include at top of house curve")
)

func main() {
	flag.Parse()
	if *in == "" {
		log.Fatalf("Please specify an input file")
	}
	commentString := *comment
	if commentString == "" {
		commentString = "house curve generated with mkhousecurve from " + *in
	}

	inFile, err := os.Open(*in)
	if err != nil {
		log.Fatalf("Unable to open input file: %v", err)
	}
	defer inFile.Close()

	outFile, err := os.Create(*out)
	if err != nil {
		log.Fatalf("Unable to open output file: %v", err)
	}
	defer outFile.Close()

	csvIn := csv.NewReader(inFile)
	csvIn.Comma = '\t'
	csvIn.FieldsPerRecord = 3
	csvOut := csv.NewWriter(outFile)
	csvOut.Comma = '\t'

	err = csvOut.Write([]string{"# " + commentString})
	if err != nil {
		log.Fatalf("Error on writing comment: %v", err)
	}

readLoop:
	for {
		row, err := csvIn.Read()
		if err != nil {
			switch t := err.(type) {
			case *csv.ParseError:
				if t.Err == csv.ErrFieldCount {
					// ignore
					continue readLoop
				}
			default:
				if err != io.EOF {
					log.Fatalf("Unexpected error reading from input: %v", err)
				}
				csvOut.Flush()
				log.Debugf("Wrote house curve to %v", *out)
				return
			}
		}

		if strings.Contains(row[0], "*") {
			// ignore header
			continue
		}

		hz := row[0]
		db, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			log.Fatalf("Unable to parse dB level from %v", row)
		}
		err = csvOut.Write([]string{hz, fmt.Sprint(db - *refDB)})
		if err != nil {
			log.Fatalf("Unable to write row to output: %v", err)
		}
	}
}
