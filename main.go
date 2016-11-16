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
    fmt.Printf(token)
}

func handleJavaClient(w http.ResponseWriter, r *http.Request) {
    lvl := r.FormValue("emergencyLevel")
    address := r.FormValue("address")
    service := r.FormValue("service")
    fmt.Println("%s, %s, %s", service, lvl, address)
}

func main() {

    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/android", handleAndroidClient).Methods("POST")
    router.HandleFunc("/java", handleJavaClient).Methods("POST")

    fmt.Println("listening...")
    err := http.ListenAndServe(":"+os.Getenv("PORT"), router)
    if err != nil {
        panic(err)
    }
}