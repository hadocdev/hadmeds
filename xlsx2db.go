package main

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type Cell struct {
	DataType string `xml:"t,attr"`
	Value    string `xml:"v"`
}

type Row struct {
	Cells []Cell `xml:"c"`
}

type SheetData struct {
	Rows []Row `xml:"row"`
}

type WorkSheet struct {
	SheetData SheetData `xml:"sheetData"`
}

type StringItem struct {
	Text string `xml:"t"`
}

type SharedStringTable struct {
	UniqueCount int          `xml:"uniqueCount,attr"`
	StringItems []StringItem `xml:"si"`
}

func extractDataFromXlsx(filename string, fieldsCount int, hasHeaders bool) ([][]string, error) {
	var data [][]string

	var zipfile *zip.ReadCloser
	var err error
	zipfile, err = zip.OpenReader(filename)
	if err != nil {
		return data, err
	}
	defer zipfile.Close()

	var sheetFile, stringsFile io.ReadCloser
	for _, f := range zipfile.File {
		if f.Name == "xl/worksheets/sheet1.xml" {
			sheetFile, err = f.Open()
			if err != nil {
				return data, err
			}
		} else if f.Name == "xl/sharedStrings.xml" {
			stringsFile, err = f.Open()
			if err != nil {
				return data, err
			}
		}
	}
	defer sheetFile.Close()
	defer stringsFile.Close()

	var xmlDecoder *xml.Decoder
	var ws WorkSheet
	var sst SharedStringTable
	xmlDecoder = xml.NewDecoder(stringsFile)
	err = xmlDecoder.Decode(&sst)
	if err != nil {
		return data, err
	}

	xmlDecoder = xml.NewDecoder(sheetFile)
	err = xmlDecoder.Decode(&ws)
	if err != nil {
		return data, err
	}
	data = make([][]string, len(ws.SheetData.Rows))
	for i, row := range ws.SheetData.Rows {
		if len(row.Cells) < fieldsCount {
			continue
		}
		if hasHeaders {
			hasHeaders = false
			continue
		}
		data[i] = make([]string, len(row.Cells))
		for j, cell := range row.Cells {
			if cell.DataType == "s" {
				index, err := strconv.Atoi(cell.Value)
				if err != nil {
					return data, err
				}
				if index < 0 || index >= sst.UniqueCount {
					return data, fmt.Errorf("Out of bounds in SharedStringTable.")
				}
				data[i][j] = sst.StringItems[index].Text
			} else {
				data[i][j] = cell.Value
			}
		}
	}
	return data, nil
}

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatal("No filename provided.")
	}
	filename := args[1]
	fieldsCount := 11
	hasHeaders := true
	data, err := extractDataFromXlsx(filename, fieldsCount, hasHeaders)
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range data {
		if len(item) > 0 {
			fmt.Println(item)
		}
	}
}
