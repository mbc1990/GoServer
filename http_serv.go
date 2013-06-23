package main

import (
    "fmt"
    "net"
    "log"
    "strings"
    "time"
    "os"
    "io"
)

func handle_connection(conn net.Conn) {
    fmt.Println("new connection")
    //read request
    buf := make([]byte, 2048)
    bread, err := conn.Read(buf)
    if err != nil {
        log.Printf("read error %s",err)
    }
    str := string(buf[:bread])
    fmt.Printf("%s\n\n",str)

    //parse request
    if !strings.Contains(str, "GET") {
        fmt.Printf("not a GET request")
        conn.Close()
        return
    }

    split := strings.Split(str, " ")
    path := split[1]
    if path == "/" {
        path = "index.html"
    } else {
        //get rid of '/' character
        path = path[1:]
    }

    if _, err := os.Stat(path); os.IsNotExist(err) {
        header := "HTTP/1.1 404 Not Found\r\nDate: "+time.Now().Format("2006-01-02 15:04")+"\r\nContent-Type: text/HTML\r\nContent-Length: 45\r\n\r\n<html><body>Page not found, sry</body></html>"
        fmt.Println(header)
        fmt.Println("\n\n")
        conn.Write([]byte(header))
        conn.Close()
        return
    }

    //open requested file
    fi, err := os.Open(path)
    st, _ := fi.Stat()
    fsize := st.Size()
    header := fmt.Sprintf("HTTP/1.1 200 OK\r\nDate: "+time.Now().Format("2006-01-02 15:04")+"\r\nContent-Type: text/HTML\r\nContent-Length: %d\r\n\r\n",fsize)
//    header := "HTTP/1.1 200 OK\r\nDate: "+time.Now().Format("2006-01-02 15:04")+"\r\nContent-Type: text/HTML\r\nContent-Length: "+fsize+"\r\n\r\n"
    conn.Write([]byte(header))
    fmt.Println(header)

    if err != nil {
        log.Printf("file io error %s",err)
    }
    //close the file after it's been sent 
    defer func() {
        if err := fi.Close(); err != nil {
            log.Printf("file io error %s",err)
        }
    }()

    //write requested file to the user
    wbuf := make([]byte, 1024)
    for {
        n, err := fi.Read(wbuf)
        if err != nil && err != io.EOF {
            log.Printf("error writing to conn buffer %s",err)
        }
        if n == 0 { break } //done reading/writing
        conn.Write(wbuf)
        fmt.Printf(string(wbuf))
    }
    fmt.Println("\n\n")
    conn.Close()
    return
}

func main() {
    fmt.Println("Starting server")
    ln, err := net.Listen("tcp", ":8080")
    if err != nil {
        fmt.Println("error")
    }

    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println("error accepting connection")
        }
        go handle_connection(conn)
    }
}
