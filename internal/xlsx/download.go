package xlsx

import (
    "fmt"
    "io"
    "os"
    "net/http"
)

func humanReadableSize(size int64) string {
    if size < 1024 { 
        return fmt.Sprintf("%d B", size) 
    } else if size < 1024*1024 { 
        return fmt.Sprintf("%.2f kB", float64(size)/1024) 
    } else if size < 1024*1024*1024 { 
        return fmt.Sprintf("%.2f MB", float64(size)/(1024*1024)) 
    } else { 
        return fmt.Sprintf("%.2f GB", float64(size)/(1024*1024*1024)) 
    }
}

func fileExists(filename string) bool {
    _, err := os.Stat(filename)
    return !os.IsNotExist(err)
}

func downloadFile(url string, filename string, hasToCheck bool) error {
    if hasToCheck && fileExists(filename) {
        fmt.Printf("File %s already exists.\n", filename)
        return nil
    }

    file, err := os.Create(filename)
    if err != nil { return err }
    defer file.Close()
    
    res, err := http.Get(url)
    if err != nil { return err }
    if res.StatusCode < 200 && res.StatusCode > 299 { 
        return fmt.Errorf("Unsuccessful request: status %d", res.StatusCode) 
    }
    defer res.Body.Close()

    var bytesRead int64 = 0
    buffer := make([]byte, 1024)
    for {
        n, err := res.Body.Read(buffer)
        bytesRead += int64(n)
        fmt.Print("\r\033[K")
        fmt.Printf("Downloading %v: %v", filename, humanReadableSize(bytesRead))
        file.Write(buffer[:n])
        if err == io.EOF { break }
    }
    fmt.Printf("\r\033[KDownloaded %v: %v\n", filename, humanReadableSize(bytesRead))
    return nil 
}

func DownloadXlsxFiles(hasToCheck bool) error {
    local_drugs_url := "http://dgdagov.info/administrator/components/com_jcode/source/ExcelRegisterProduct.php?action=excelforDugDatabase"
    imported_drugs_url := "http://dgdagov.info/administrator/components/com_jcode/source/Excelforeign_drugs_report.php?action=excelRegisteredImportedDrugs"
    
    err := downloadFile(local_drugs_url, "local.xlsx", hasToCheck)
    if err != nil { return err }

    err = downloadFile(imported_drugs_url, "imported.xlsx", hasToCheck)
    if err != nil { return err }
    
    return nil
}
