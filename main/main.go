package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/piyapongch/era-audit/foxml2csv"
)

func main() {

	// parse command line agruments
	var input string
	var output string
	flag.StringVar(&input, "i", ".", "a foxml input file directory")
	flag.StringVar(&output, "o", "results.csv", "a result csv file")
	flag.Parse()
	fmt.Println("i:", input)
	fmt.Println("o:", output)

	// create output file
	file, err := os.Create(output)
	if err != nil {
		log.Fatal("Could not create output file!", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// declare vars
	var csv []string
	var i int

	// create csv header
	csv = append(csv, "PID")
	csv = append(csv, "UUID")
	csv = append(csv, "handle")
	csv = append(csv, "title")
	csv = append(csv, "deposit date")
	csv = append(csv, "# DCQ datastreams")
	csv = append(csv, "# DC datastreams")
	csv = append(csv, "# DS datastreams")
	csv = append(csv, "# LICENSE datastreams")
	csv = append(csv, "DCQ IDs")
	csv = append(csv, "DCQ MD5s")
	csv = append(csv, "DCQ version changes")
	csv = append(csv, "DS IDs")
	csv = append(csv, "DS MD5s")
	csv = append(csv, "DS MIME Type")
	csv = append(csv, "DS version changes")
	csv = append(csv, "LICENSE IDs")
	csv = append(csv, "LICENSE version changes")

	// write header to csv file
	err = writer.Write(csv)
	if err != nil {
		log.Fatal("Could not write csv header!", err)
	}

	files, err := ioutil.ReadDir(input)
	if err != nil {
		log.Fatal("Could not read directory!", err)
	}

	for _, file := range files {
		i++
		fmt.Printf("# %d ", i)
		foxml2csv.Run(input+"/"+file.Name(), writer)
	}
}
