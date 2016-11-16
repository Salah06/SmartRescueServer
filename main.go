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
    // "net"
    // "bufio"
    "log"
    "fmt"
    "net/http"
    "io/ioutil"
    //"os"
    "bytes"

    "gopkg.in/zabawaba99/firego.v1"
    "github.com/gorilla/mux"
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

func handleJavaClient(w http.ResponseWriter, r *http.Request) {
    b,err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("%s", b)
    fmt.Println(r.Methode)

}

func handleAndroidClient(w http.ResponseWriter, r *http.Request) {
    token := make(chan string)

    b,err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Fatal(err)
    }
    n := bytes.Index(b, []byte{0})
    token <- string(b[:n])
    go tokenPrinter(token)
}

func tokenPrinter(token chan string) {
    for {
        t := <- token
        fmt.Println("token = %s", t)
    }
}

func main() {
    //f.Auth("null")
    //pushValue()
    //getValue()

    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/android", handleAndroidClient)
    router.HandleFunc("/java", handleJavaClient)

    fmt.Println("listening...")
    //err := http.ListenAndServe(":"+os.Getenv("PORT"), router)
    err := http.ListenAndServe(":1234", router)
    if err != nil {
        panic(err)
    }
}