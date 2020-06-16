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
)

func main() {
    username, passwd := credentials()
//    fmt.Printf("Username: %s, Password: %s\n", username, passwd)
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

    passwd := getPassword("Enter Password: ")

    return strings.TrimSpace(username), strings.TrimSpace(passwd)
}

func getPassword(prompt string) string {
    fmt.Print(prompt)

    // Common settings and variables for both stty calls.
    attrs := syscall.ProcAttr{
        Dir:   "",
        Env:   []string{},
        Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
        Sys:   nil}
    var ws syscall.WaitStatus

    // Disable echoing.
    pid, err := syscall.ForkExec(
        "/bin/stty",
        []string{"stty", "-echo"},
        &attrs)
    if err != nil {
        panic(err)
    }

    // Wait for the stty process to complete.
    _, err = syscall.Wait4(pid, &ws, 0, nil)
    if err != nil {
        panic(err)
    }

    // Echo is disabled, now grab the data.
    reader := bufio.NewReader(os.Stdin)
    text, err := reader.ReadString('\n')
    if err != nil {
        panic(err)
    }

    // Re-enable echo.
    pid, err = syscall.ForkExec(
        "/bin/stty",
        []string{"stty", "echo"},
        &attrs)
    if err != nil {
        panic(err)
    }

    // Wait for the stty process to complete.
    _, err = syscall.Wait4(pid, &ws, 0, nil)
    if err != nil {
        panic(err)
    }

    return strings.TrimSpace(text)
}
