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
var repartiteur map[string]chan([]string)


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

    go broadcastInit(msgEmergency ,address, id)
}

func handleAndroidClient(w http.ResponseWriter, r *http.Request) {
    // recup l'id de l'emergency dans le message pour chopé le chan
    idEmergency := r.FormValue("idEmergency")
    token := r.FormValue("token")
    response := r.FormValue("response")
    fmt.Printf(response)
    fmt.Printf(idEmergency)
    fmt.Printf(token)

    rep := []{token, response}
    repartiteur[idEmergency] <- rep
}

func broadcastInit(msg map[string]string, address string, id string) { // id est aussi dans map...
    // init une liste de token potentiel pour l'intervention
    tokensPerimeter := spot(address, 2)
    memo[id] = tokensPerimeter
    repartiteur[id] = make(chan string)
    go listenResponse(id)
    go sendAndroids(tokens, msg)

}

func listenResponse(id string, numberNecessary int) {
    c := repartiteur[id]
    inCharge := []string
    for {
        rep := <- c
        switch rep[1] {
        case "ok" :
            inCharge = append(inCharge, rep[0])
            t := []{rep[0]}
            r := "go go go"
            sendAndroids(t, r)
        case "not" :
            continue
    }
    if len(inCharge) == numberNecessary {
        memo[id] = inCharge
        return
    }
}

func sendAndroids(tokens []string, msg string) {
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
    fmt.Printf("%s\n", v) // oh le format degueu
}

func main() {

    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/android", handleAndroidClient).Methods("POST")
    router.HandleFunc("/java", handleJavaClient).Methods("POST")

    catchGPS()

    fmt.Println("listening...")
    err := http.ListenAndServe(":"+os.Getenv("PORT"), router)
    if err != nil {
        panic(err)
    }
}