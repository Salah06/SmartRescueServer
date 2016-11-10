//
// gosimpleserv.go  -- A simple server
//
//           Author: Erick Gallesio [eg@unice.fr]
//    Creation date: 17-Oct-2016 14:42 (eg)
// Last file update: 17-Oct-2016 15:50 (eg)
//
// +++
package main

import (
    "net"
    "bufio"
    "log"
    "fmt"
    "gopkg.in/zabawaba99/firego.v1"
)


var f = firego.New("https://testrescue-d8b04.firebaseio.com/", nil)

func pushValue() {
    v := "bar"
    pushedFirego, err := f.Push(v)
    if err != nil {
        log.Fatal(err)
    }

    var bar string
    if err := pushedFirego.Value(&bar); err != nil {
        log.Fatal(err)
    }

    // prints "https://my-firebase-app.firebaseIO.com/-JgvLHXszP4xS0AUN-nI: bar"
    fmt.Printf("%s: %s\n", pushedFirego, bar)
}

func getValue() {
    var v map[string]interface{}
    if err := f.Value(&v); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s\n", v)
}

func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()

    reader := bufio.NewReader(conn)
    for {
        line,err := reader.ReadString('\n')
        if err != nil {
            log.Println(err)
            return
        }
        fmt.Fprintf(conn, "You said => %s", line);
    }
}

func main() {
    //f.Auth("null")
    //pushValue()
    //getValue()

    listener, err := net.Listen("tcp", "localhost:1234")
    if err != nil {
        log.Fatal(err)
    }

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Println(err)
            continue
        }
        go handleConnection(conn)
    }
}
