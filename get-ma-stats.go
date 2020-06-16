package main

import (
    "archive/zip"
    "bytes"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "bufio"
    "os"
    "strings"
    "syscall"
    "golang.org/x/crypto/ssh/terminal"

)

func main() {
    username, passwd := credentials()
    fmt.Printf("Username: %s, Password: %s\n", username, passwd)
    client := &http.Client{}
    req, err := http.NewRequest("GET", "https://cloud.redhat.com/api/xavier/administration/report/csv", nil)
    req.SetBasicAuth(username, passwd)
    resp, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
    if err != nil {
        log.Fatal(err)
    }

    // Read all the files from zip archive
    for _, zipFile := range zipReader.File {
        fmt.Println("Reading file:", zipFile.Name)
        unzippedFileBytes, err := readZipFile(zipFile)
        if err != nil {
            log.Println(err)
            continue
        }

        _ = unzippedFileBytes // this is unzipped file bytes
    }
}

func readZipFile(zf *zip.File) ([]byte, error) {
    f, err := zf.Open()
    if err != nil {
        return nil, err
    }
    defer f.Close()
    return ioutil.ReadAll(f)
}


func credentials() (string, string) {
    reader := bufio.NewReader(os.Stdin)

    fmt.Print("Enter Username: ")
    username, _ := reader.ReadString('\n')

    fmt.Print("Enter Password: ")
    bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
    if err == nil {
        fmt.Println("\nPassword typed: " + string(bytePassword))
    }
    password := string(bytePassword)

    return strings.TrimSpace(username), strings.TrimSpace(password)
}
