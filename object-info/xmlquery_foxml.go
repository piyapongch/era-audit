package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/antchfx/xmlquery"
)

func main() {

	// create output file
	file, err := os.Create("result.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// open foxml
	f, err1 := os.Open("../sample_data/foxml_wolves.xml")
	doc, err1 := xmlquery.Parse(f)
	if err1 != nil {
		panic(err1)
	}

	// csv string array
	var csv []string
	var val string

	// PID
	n := xmlquery.FindOne(doc, "/foxml:digitalObject")
	if n != nil {
		val = n.SelectAttr("PID")
	}
	fmt.Printf("PID: %s\n", val)
	csv = append(csv, val)

	// handle
	for i, n := range xmlquery.Find(doc, "/foxml:digitalObject/foxml:datastream[@ID='DCQ']/foxml:datastreamVersion[last()]/foxml:xmlContent/dc/dcterms:identifier") {
		if val = ""; n != nil {
			if hd := n.InnerText(); strings.HasPrefix(hd, "http://hdl.handle.net") {
				val = hd
			}
		}
		fmt.Printf("handle: #%d %s\n", i, val)
		csv = append(csv, val)
	}

	err = writer.Write(csv)
	checkError("Cannot write to file", err)
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
