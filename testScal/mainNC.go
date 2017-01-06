package main

import (
    "fmt"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
)

var chans []chan int
var tokenDelete = make(chan string)
var memo = make(map[string][]string)

func handleAndroidClientEmergency(w http.ResponseWriter, r *http.Request) {
    idEmergency := r.FormValue("id")

    tokenDelete <- idEmergency
}

func fillEmergency(id string) {
    tokensPerimeter := fillTokens(100)
    memo[id] = make([]string, 100)
    memo[id] = tokensPerimeter
}

func fillTokens(n int) []string {
    var tokensTmp []string
    for i := 0; i < n; i++ {
        tokensTmp = append(tokensTmp, strconv.Itoa(i))
    }
    return tokensTmp
}

func checkEmpty() {
    acc := 0
    for i := 0; i < 1000; i++ {
        if len(memo[strconv.Itoa(i)]) == 1 {
            acc += 1
        }
    }
    if acc == 1000 {
        // si on arrive là tout les urgences on 1 token
        fmt.Println("TOUT LES MESSAGES SONT ARRIVES !")
    }
}

func deleteToken() {
    for {
        id := <- tokenDelete
        //fmt.Printf("---%s---", id)
        memo[id] = memo[id][:len(memo[id])-1]
        //fmt.Print("-1")
        checkEmpty()  // juste pour check mais grosse perte de temps? en faite non -_-
    }
}

func main() {

    for i := 0; i < 1000; i++ {
        fillEmergency(strconv.Itoa(i))
    }

    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/androidEmergency", handleAndroidClientEmergency).Methods("POST")

    //fmt.Println(memo["998"][99])
    go checkEmpty()  // fatal error: concurrent map read and map write
    go deleteToken()

    fmt.Println("listening...")
    err := http.ListenAndServe(":1234", router)
    if err != nil {
        panic(err)
    }
}