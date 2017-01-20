package main

import (
    "fmt"
    "net/http"
    "os"
    "log"
    "strconv"

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
var idEmergency int


func handleJavaClient(w http.ResponseWriter, r *http.Request) {
    lvl := r.FormValue("emergencyLevel")
    address := r.FormValue("address")
    //service := r.FormValue("service")

    msgEmergency := map[string]string{
        "id" : strconv.Itoa(idEmergency),
        "msg": address,
        "emergencyLevel": lvl,
    }
    idEmergency += 1

    msgFinal := map[string]interface{}{
        "request" : "caca",
        "data": msgEmergency,
    }

    go broadcastInit(msgFinal ,address, strconv.Itoa(idEmergency))
}

func handleAndroidClient(w http.ResponseWriter, r *http.Request) {
    // recup l'id de l'emergency dans le message pour chopé le chan
    idEmergency := r.FormValue("idEmergency")
    token := r.FormValue("token")
    response := r.FormValue("response")
    fmt.Printf(response)
    fmt.Printf(idEmergency)
    fmt.Printf(token)

    rep := []string{token, response}
    repartiteur[idEmergency] <- rep
}

func broadcastInit(msg map[string]interface{}, address string, id string) { // id est aussi dans map...
    // init une liste de token potentiel pour l'intervention
    tokensPerimeter := spot(address, 2)
    memo[id] = tokensPerimeter
    repartiteur[id] = make(chan []string)
    go listenResponse(id, 10)
    go sendAndroids(tokens, msg)

}

func listenResponse(id string, numberNecessary int) {
    c := repartiteur[id]
    inCharge := []string{}
    for {
        rep := <- c
        switch rep[1] {
        case "ok" :
            inCharge = append(inCharge, rep[0])
            t := []string{rep[0]}
            r := map[string]interface{}{ 
                "msg" : "go go go",
            }
            rf := map[string]interface{}{ 
                "request" : "pipi",
                "data" : r,
            }
            sendAndroids(t, rf)
        case "ko" :
            continue
        }
    }
    if len(inCharge) == numberNecessary {
        memo[id] = inCharge
        return
    }
}

func sendAndroids(tokens []string, msg map[string]interface{}) {
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
    t := v["ambulance-1"]
    t = t.(map[string]interface{})
    fmt.Printf("%s\n", t)
    //fmt.Printf("%s\n", t["token"])
}

func main() {

    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/android", handleAndroidClient).Methods("POST")
    router.HandleFunc("/java", handleJavaClient).Methods("POST")

    catchGPS() // a virer

    fmt.Println("listening...")
    err := http.ListenAndServe(":"+os.Getenv("PORT"), router)
    if err != nil {
        panic(err)
    }
}