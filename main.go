package main

import (
    "log"
    "fmt"
    "net/http"
    "io/ioutil"
    "os"
    "bytes"

    "github.com/gorilla/mux"
)



func handleAndroidClient(w http.ResponseWriter, r *http.Request) {
    b,err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Fatal(err)
    }
    n := bytes.Index(b, []byte{0})
    fmt.Println(n)
    fmt.Println("%s", b)
    fmt.Println("caca")
}

func main() {

    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/android", handleAndroidClient)

    fmt.Println("listening...")
    err := http.ListenAndServe(":"+os.Getenv("PORT"), router)
    if err != nil {
        panic(err)
    }
}