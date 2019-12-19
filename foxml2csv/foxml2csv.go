package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/antchfx/xpath"
)

func main() {

	// create output file
	file, err := os.Create("results.csv")
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

	// declare vars
	var csv []string
	var val string

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
	csv = append(csv, "DS IDs")
	csv = append(csv, "DS MD5s")
	csv = append(csv, "LICENSE IDs")

	// write header to csv file
	err = writer.Write(csv)
	checkError("Could not write to file", err)

	// clear array
	csv = nil

	// PID
	n := xmlquery.FindOne(doc, "/foxml:digitalObject/@PID")
	if val = ""; n != nil {
		val = n.InnerText()
	}
	fmt.Printf("PID: %s\n", val)
	csv = append(csv, val)

	// UUID
	fmt.Printf("UUID: %s\n", val)
	csv = append(csv, val)

	// handle
	val = ""
	for _, n := range xmlquery.Find(doc, "foxml:digitalObject/foxml:datastream[@ID='DCQ']/foxml:datastreamVersion[last()]/foxml:xmlContent/dc/dc:identifier") {
		if n != nil {
			if val = n.InnerText(); strings.HasPrefix(val, "http://hdl.handle.net") {
				break
			}
		}
	}
	if val == "" {
		for _, n := range xmlquery.Find(doc, "foxml:digitalObject/foxml:datastream[@ID='DCQ']/foxml:datastreamVersion[last()]/foxml:xmlContent/dc/dcterms:identifier") {
			if n != nil {
				if val = n.InnerText(); strings.HasPrefix(val, "http://hdl.handle.net") {
					break
				}
			}
		}
	}
	fmt.Printf("handle: %s\n", val)
	csv = append(csv, val)

	// titles
	val = ""
	for i, n := range xmlquery.Find(doc, "/foxml:digitalObject/foxml:datastream[@ID='DCQ']/foxml:datastreamVersion[last()]/foxml:xmlContent/dc/dc:title") {
		if n != nil {
			if txt := n.InnerText(); txt != "" {
				if i > 0 {
					val = val + "|"
				}
				val += txt
			}
		}
	}
	fmt.Printf("titles: %s\n", val)
	csv = append(csv, val)

	// deposit date
	n = xmlquery.FindOne(doc, "/foxml:digitalObject/foxml:datastream[@ID='DCQ']/foxml:datastreamVersion[last()]/foxml:xmlContent/dc/dcterms:dateaccepted")
	if val = ""; n == nil {
		n = xmlquery.FindOne(doc, "/foxml:digitalObject/foxml:objectProperties/foxml:property[@NAME='info:fedora/fedora-system:def/model#createdDate']/@VALUE")
		if val = ""; n != nil {
			val = n.InnerText()
		}
	} else {
		val = n.InnerText()
	}
	fmt.Printf("deposit date: %s\n", val)
	csv = append(csv, val)

	// # DCQ datastreams
	expr, err2 := xpath.Compile("count(/foxml:digitalObject/foxml:datastream[@ID='DCQ']/foxml:datastreamVersion)")
	if val = ""; err2 != nil {
		checkError("Could not find DCQ datastream!", err2)
	} else {
		dcq := expr.Evaluate(xmlquery.CreateXPathNavigator(doc)).(float64)
		val = fmt.Sprintf("%.0f", dcq)
	}
	fmt.Printf("# DCQ datasreams: %s\n", val)
	csv = append(csv, val)

	// # DC datastreams
	expr, err2 = xpath.Compile("count(/foxml:digitalObject/foxml:datastream[@ID='DC']/foxml:datastreamVersion)")
	if val = ""; err2 != nil {
		checkError("Could not find DC datastream!", err2)
	} else {
		dc := expr.Evaluate(xmlquery.CreateXPathNavigator(doc)).(float64)
		val = fmt.Sprintf("%.0f", dc)
	}
	fmt.Printf("# DC datasreams: %s\n", val)
	csv = append(csv, val)

	// # DS datastreams
	expr, err2 = xpath.Compile("count(/foxml:digitalObject/foxml:datastream[starts-with(@ID,'DS')]/foxml:datastreamVersion)")
	if val = ""; err2 != nil {
		checkError("Could not find DS datastream!", err2)
	} else {
		ds := expr.Evaluate(xmlquery.CreateXPathNavigator(doc)).(float64)
		val = fmt.Sprintf("%.0f", ds)
	}
	fmt.Printf("# DS datasreams: %s\n", val)
	csv = append(csv, val)

	// # LICENSE datastreams
	expr, err2 = xpath.Compile("count(/foxml:digitalObject/foxml:datastream[@ID='LICENSE']/foxml:datastreamVersion)")
	if val = ""; err2 != nil {
		checkError("Could not find DS datastream!", err2)
	} else {
		ds := expr.Evaluate(xmlquery.CreateXPathNavigator(doc)).(float64)
		val = fmt.Sprintf("%.0f", ds)
	}
	fmt.Printf("# LICENSE datastreams: %s\n", val)
	csv = append(csv, val)

	// DCQ IDs
	val = ""
	for i, n := range xmlquery.Find(doc, "/foxml:digitalObject/foxml:datastream[@ID='DCQ']/foxml:datastreamVersion/@ID") {
		if n != nil {
			if txt := n.InnerText(); txt != "" {
				if i > 0 {
					val = val + "|"
				}
				val += txt
			}
		}
	}
	fmt.Printf("DCQ IDs: %s\n", val)
	csv = append(csv, val)

	// DCQ MD5s
	val = ""
	for i, n := range xmlquery.Find(doc, "/foxml:digitalObject/foxml:datastream[@ID='DCQ']/foxml:datastreamVersion") {
		if n != nil {
			e := n.SelectElement("foxml:contentDigest")
			a := ""
			if e != nil {
				a = e.SelectAttr("DIGEST")
			}
			if i > 0 {
				val = val + "|"
			}
			val += a
		}
	}
	fmt.Printf("DCQ MD5s: %s\n", val)
	csv = append(csv, val)

	// DS IDs
	val = ""
	for i, n := range xmlquery.Find(doc, "/foxml:digitalObject/foxml:datastream/foxml:datastreamVersion[starts-with(@ID,'DS')]/@ID") {
		if n != nil {
			if txt := n.InnerText(); txt != "" {
				if i > 0 {
					val = val + "|"
				}
				val += txt
			}
		}
	}
	fmt.Printf("DS IDs: %s\n", val)
	csv = append(csv, val)

	// DS MD5s
	val = ""
	for i, n := range xmlquery.Find(doc, "/foxml:digitalObject/foxml:datastream/foxml:datastreamVersion[starts-with(@ID,'DS')]") {
		if n != nil {
			e := n.SelectElement("foxml:contentDigest")
			a := ""
			if e != nil {
				a = e.SelectAttr("DIGEST")
			}
			if i > 0 {
				val = val + "|"
			}
			val += a
		}
	}
	fmt.Printf("DS MD5s: %s\n", val)
	csv = append(csv, val)

	// LICENSE IDs
	val = ""
	for i, n := range xmlquery.Find(doc, "/foxml:digitalObject/foxml:datastream/foxml:datastreamVersion[starts-with(@ID,'LICENSE')]/@ID") {
		if n != nil {
			if i > 0 {
				val = val + "|"
			}
			val += n.InnerText()
		}
	}
	fmt.Printf("LICENSE IDs: %s\n", val)
	csv = append(csv, val)

	// write to csv file
	err = writer.Write(csv)
	checkError("Could not write to file", err)
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
