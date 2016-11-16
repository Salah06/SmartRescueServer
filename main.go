package main

import (
    "fmt"
    "net/http"
    "os"

    "github.com/gorilla/mux"
)



func handleAndroidClient(w http.ResponseWriter, r *http.Request) {
    vehiculeId := r.FormValue("vehiculeId")
    token := r.FormValue("token")
    fmt.Printf(vehiculeId)
    fmt.Printf("coucou")
    fmt.Printf(token)
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