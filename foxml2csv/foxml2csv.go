package foxml2csv

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/antchfx/xpath"
)

func Run(input string, writer *csv.Writer) {

	// open foxml
	f, err1 := os.Open(input)
	doc, err1 := xmlquery.Parse(f)
	if err1 != nil {
		fmt.Printf("Error: %s, %s\n", input, err1)
		return
	}

	// declare variable
	var val string
	var csv []string

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
	csv = append(csv, val)

	// # DCQ datastreams
	expr, err2 := xpath.Compile("count(/foxml:digitalObject/foxml:datastream[@ID='DCQ']/foxml:datastreamVersion)")
	if val = ""; err2 != nil {
		checkError("Could not find DCQ datastream!", err2)
	} else {
		dcq := expr.Evaluate(xmlquery.CreateXPathNavigator(doc)).(float64)
		val = fmt.Sprintf("%.0f", dcq)
	}
	csv = append(csv, val)

	// # DC datastreams
	expr, err2 = xpath.Compile("count(/foxml:digitalObject/foxml:datastream[@ID='DC']/foxml:datastreamVersion)")
	if val = ""; err2 != nil {
		checkError("Could not find DC datastream!", err2)
	} else {
		dc := expr.Evaluate(xmlquery.CreateXPathNavigator(doc)).(float64)
		val = fmt.Sprintf("%.0f", dc)
	}
	csv = append(csv, val)

	// # DS datastreams
	expr, err2 = xpath.Compile("count(/foxml:digitalObject/foxml:datastream[starts-with(@ID,'DS')]/foxml:datastreamVersion)")
	if val = ""; err2 != nil {
		checkError("Could not find DS datastream!", err2)
	} else {
		ds := expr.Evaluate(xmlquery.CreateXPathNavigator(doc)).(float64)
		val = fmt.Sprintf("%.0f", ds)
	}
	csv = append(csv, val)

	// # LICENSE datastreams
	expr, err2 = xpath.Compile("count(/foxml:digitalObject/foxml:datastream[@ID='LICENSE']/foxml:datastreamVersion)")
	if val = ""; err2 != nil {
		checkError("Could not find DS datastream!", err2)
	} else {
		ds := expr.Evaluate(xmlquery.CreateXPathNavigator(doc)).(float64)
		val = fmt.Sprintf("%.0f", ds)
	}
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
	csv = append(csv, val)

	// DCQ version changes
	val = ""
	for i, n := range xmlquery.Find(doc, "/foxml:digitalObject/foxml:datastream/foxml:datastreamVersion/foxml:xmlContent/audit:auditTrail/audit:record/audit:componentID[starts-with(text(),'DCQ')]") {
		if n != nil {
			if txt := n.Parent.SelectElement("audit:justification").InnerText(); txt != "" {
				if i > 0 {
					val = val + "|"
				}
				val += txt
			}
		}
	}
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
	csv = append(csv, val)

	// DS MimeType
	val = ""
	for i, n := range xmlquery.Find(doc, "/foxml:digitalObject/foxml:datastream/foxml:datastreamVersion[starts-with(@ID,'DS')]/@MIMETYPE") {
		if n != nil {
			if txt := n.InnerText(); txt != "" {
				if i > 0 {
					val = val + "|"
				}
				val += txt
			}
		}
	}
	csv = append(csv, val)

	// DS version changes
	val = ""
	for i, n := range xmlquery.Find(doc, "/foxml:digitalObject/foxml:datastream/foxml:datastreamVersion/foxml:xmlContent/audit:auditTrail/audit:record/audit:componentID[starts-with(text(),'DS')]") {
		if n != nil {
			if txt := n.Parent.SelectElement("audit:justification").InnerText(); txt != "" {
				if i > 0 {
					val = val + "|"
				}
				val += txt
			}
		}
	}
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
	csv = append(csv, val)

	// LICENSE version changes
	val = ""
	for i, n := range xmlquery.Find(doc, "/foxml:digitalObject/foxml:datastream/foxml:datastreamVersion/foxml:xmlContent/audit:auditTrail/audit:record/audit:componentID[starts-with(text(),'LICENSE')]") {
		if n != nil {
			if txt := n.Parent.SelectElement("audit:justification").InnerText(); txt != "" {
				if i > 0 {
					val = val + "|"
				}
				val += txt
			}
		}
	}
	csv = append(csv, val)

	// write to csv file
	err1 = writer.Write(csv)
	checkError("Could not write to file", err1)
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
