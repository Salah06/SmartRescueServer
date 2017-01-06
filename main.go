package main

import (
    "fmt"
    "net/http"
    "os"
    "log"

    "github.com/NaySoftware/go-fcm"
    "github.com/gorilla/mux"
    "gopkg.in/zabawaba99/firego.v1" 
)

const (
     serverKey = "AIzaSyAM5yN0SNAswN6l6t6DEKv9fLRSeUaliVY"
)
var firebase = firego.New("https://smartrescue-6e8ce.firebaseio.com/", nil)
var memo map[string][]string
var tokens []string


func handleJavaClient(w http.ResponseWriter, r *http.Request) {
    id := r.FormValue("id")
    lvl := r.FormValue("emergencyLevel")
    address := r.FormValue("address")
    //service := r.FormValue("service")

    msgEmergency := map[string]string{
        "id" : id,          // !!! not yet implement
        "msg": address,
        "emergencyLevel": lvl,
    }

    go broadcastInit(msgEmergency , address, id)
}

func handleAndroidClient(w http.ResponseWriter, r *http.Request) {
    // recup l'id de l'emergency dans le message pour chopé le chan
    vehiculeId := r.FormValue("vehiculeId")
    token := r.FormValue("token")
    fmt.Printf(vehiculeId)
    fmt.Printf(token)

    tokens = append(tokens, token)
}

func broadcastInit(msg map[string]string, address string, id string) {
    // init une liste de token potentiel pour l'intervention
    tokensPerimeter := spot(address, 2)
    memo[id] = tokensPerimeter

    c := fcm.NewFcmClient(serverKey)
    c.NewFcmRegIdsMsg(tokens, msg)
    status, err := c.Send()

    if err == nil {
    status.PrintResults()
    } else {
        fmt.Println(err)
    }
}

func spot(address string, perimeter int) []string {
    catchGPS()
    // find vehicul in perimeter of address
    return tokens
}

func catchGPS() {
    var v map[string]interface{}
    if err := firebase.Value(&v); err != nil {
        log.Fatal(err)
    }
    // a print avec [token, ...] histoire de dire qu'on recup toute la liste
    fmt.Printf("%s\n", v)
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