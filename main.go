package main

import (
    "fmt"
    "net/http"
    "os"
    "strconv"

    "github.com/NaySoftware/go-fcm"
    "github.com/gorilla/mux"
)

const (
     serverKey = "AIzaSyAM5yN0SNAswN6l6t6DEKv9fLRSeUaliVY"
)
var tokenC = make(chan map[string]string)
var tokens []string
var idEmergency = 0

func handleAndroidClient(w http.ResponseWriter, r *http.Request) {
    vehiculeId := r.FormValue("vehiculeId")
    token := r.FormValue("token")
    fmt.Printf(vehiculeId)
    fmt.Printf(token)

    tokens = append(tokens, token)
}

func handleJavaClient(w http.ResponseWriter, r *http.Request) {
    lvl := r.FormValue("emergencyLevel")
    address := r.FormValue("address")
    service := r.FormValue("service")
    fmt.Printf("%s, %s, %s", service, lvl, address)

    msgEmergency := strconv.Itoa(idEmergency) + ";" + lvl + ";" + address
    go broadcast(msgEmergency)
}

func broadcast(msg string) {
    c := fcm.NewFcmClient(serverKey)
    c.NewFcmRegIdsMsg(tokens, msg)
    status, err := c.Send()

    if err == nil {
    status.PrintResults()
    } else {
        fmt.Println(err)
    }
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