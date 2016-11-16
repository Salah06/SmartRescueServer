package main

import (
    "log"
    "fmt"
    "net/http"
    "io/ioutil"
    "os"

    "github.com/gorilla/mux"
)



func handleAndroidClient(w http.ResponseWriter, r *http.Request) {
    b,err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s",b)
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