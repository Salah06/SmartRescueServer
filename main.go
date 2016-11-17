package main

import (
    "fmt"
    "net/http"
    "os"
    "time"
    "log"

<<<<<<< HEAD
    "gopkg.in/zabawaba99/firego.v1"
=======
>>>>>>> 6b95c3267d638b8457374239d742696ce3d9ff1d
    "github.com/NaySoftware/go-fcm"
    "github.com/gorilla/mux"
)

const (
     serverKey = "AIzaSyAM5yN0SNAswN6l6t6DEKv9fLRSeUaliVY"
)
var tokenC = make(chan map[string]string)
var tokens []string
var idEmergency = 0
var firebase = firego.New("https://smartrescue-6e8ce.firebaseio.com/", nil)

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
    //service := r.FormValue("service")

    msgEmergency := map[string]string{
        "msg": address,
        "emergencyLevel": lvl,
    }
    go broadcast(msgEmergency)
}

func broadcast(msg map[string]string) {
    c := fcm.NewFcmClient(serverKey)
    c.NewFcmRegIdsMsg(tokens, msg)
    status, err := c.Send()

    if err == nil {
    status.PrintResults()
    } else {
        fmt.Println(err)
    }
}

func catchGPS() {
    for {
        time.Sleep(8 * time.Second)
        var v map[string]interface{}
        if err := firebase.Value(&v); err != nil {
            log.Fatal(err)
        }
        fmt.Printf("%s\n", v)
    }
}

func main() {

    tokens = append(tokens, "e6THtaBcNVE:APA91bESyZPEZ19jjMIpSBkry1eKAJCnYeRPsw6Dm_mMUQovH3APX4V-gSxJHHnuFK1OWhcM3dOpNw2h__sRy3HYaY5fqQ--vKwzG43WngO-XGEqO1b_X8aFM7HAioLljQH4M505RR1U")
    //fmt.Println(tokens)

    go catchGPS()
    
    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/android", handleAndroidClient).Methods("POST")
    router.HandleFunc("/java", handleJavaClient).Methods("POST")

    fmt.Println("listening...")
    err := http.ListenAndServe(":"+os.Getenv("PORT"), router)
    if err != nil {
        panic(err)
    }
}