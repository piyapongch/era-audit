package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/antchfx/xmlquery"
)

func main() {
	file, err := os.Create("result.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	f, err1 := os.Open("../sample_data/wolves_foxml.xml")
	doc, err1 := xmlquery.Parse(f)
	if err1 != nil {
		panic(err)
	}

	var cb []string
	channel := xmlquery.FindOne(doc, "//foxml:digitalObject")
	n := channel.SelectAttr("PID")
	fmt.Printf("PID: %s\n", n)
	cb = append(cb, n)

	channel = xmlquery.FindOne(doc, "/foxml:digitalObject/foxml:datastream[2]/foxml:datastreamVersion[last()]/foxml:xmlContent/dc/dc:subject[last()]")
	n = channel.InnerText()
	fmt.Printf("subject: %s\n", n)
	cb = append(cb, n)
	n = ""
	cb = append(cb, n)

	n = ""
	cb = append(cb, n)

	err = writer.Write(cb)
	checkError("Cannot write to file", err)
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
