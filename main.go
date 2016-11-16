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
    vehiculeId := r.FormValue("vehiculeId")
    fmt.Printf("%s", vehiculeId)
}

func main() {

    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/android", handleAndroidClient).Methods("POST")

    fmt.Println("listening...")
    err := http.ListenAndServe(":"+os.Getenv("PORT"), router)
    if err != nil {
        panic(err)
    }
}