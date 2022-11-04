package utils

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type XmlRow struct {
	XMLName     xml.Name `xml:"ROW"`
	Code        string   `xml:"res_codigo"`
	Description string   `xml:"res_descripcion_es"`
}

type XmlData struct {
	XMLName xml.Name `xml:"DATA"`
	Rows    []XmlRow `xml:"ROW"`
}

func ReadData(reader io.Reader) ([]XmlRow, error) {
	var xmlData XmlData
	if err := xml.NewDecoder(reader).Decode(&xmlData); err != nil {
		return nil, err
	}

	return xmlData.Rows, nil
}

func ReadXml(code string) (description string) {
	var resp string

	strapsFilePath, err := filepath.Abs("resources/ResCodes.xml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Open the straps.xml file
	file, err := os.Open(strapsFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	// Read the DATA file
	xmlStraps, err := ReadData(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Display The first ROW
	for i := 0; i < len(xmlStraps); i++ {
		if xmlStraps[i].Code == code {
			resp = xmlStraps[i].Description
		}
	}
	return resp
}
