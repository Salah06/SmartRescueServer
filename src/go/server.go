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
)

// +slide
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
// ---
// Local Variables:
// compile-command: "go run *.go"
// End:

// pour tester : nc localhost 1234

// on met go pour lancer une fonction en arriere plan (goroutine)