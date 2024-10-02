package main

import (
	"fmt"
	"log"
	"github.com/hadocdev/hadmeds/internal/xlsx"
)

func main() {
    err := xlsx.DownloadXlsxFiles(true)
	if err != nil {
	    log.Fatal(err)
	}
	filename := "imported.xlsx"
	fieldsCount := 11
	hasHeaders := true
	data, err := xlsx.ExtractDataFromXlsx(filename, fieldsCount, hasHeaders, xlsx.ProcessRowImported)
    if err != nil { log.Fatal(err) }
	fmt.Println(len(xlsx.Generics))
    fmt.Println(data[3])
}
