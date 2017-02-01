package main

import (
    "fmt"
    "net/http"
    "os"
    "log"
    "strconv"
    "time"
    "io/ioutil"
    "strings"

    "github.com/NaySoftware/go-fcm"
    "github.com/gorilla/mux"
    "gopkg.in/zabawaba99/firego.v1" 
)

const (
     serverKey = "AIzaSyAM5yN0SNAswN6l6t6DEKv9fLRSeUaliVY"
)
var firebase = firego.New("https://smartrescue-6e8ce.firebaseio.com/", nil)
var memo  = make(map[string][]string)
var repartiteur = make(map[string]chan([]string)) 
var maxId = 0
var tokenAction []string
var reSendEmergency [][]string
var chanId = make(chan(int))

func charge() {
    acc := maxId
    for {
        chanId <- acc
        acc = acc + 1
    }
}

func getId() int {
    id := <- chanId
    return id
}

//--------------------RECOVER-------------------------------------------------------------

func saveData() {
    for {
        f, err := os.Create("/tmp/save")
        check(err)
        defer f.Close()
        save := ""

        // |id,token,token|id,token
        for id, tokens := range memo {
            save = save + "|" + id
            // tokens[0] = address si pas encore trouve tout le monde
            for _, token := range tokens {
                save = save + "," + token
            }
        }

        n3, errr := f.WriteString(save)
        check(errr)
        fmt.Printf("wrote %d bytes\n", n3)

        f.Sync()
        time.Sleep(10000 * time.Millisecond)
    }
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func recover() {
    save, err := ioutil.ReadFile("/tmp/save")
    check(err)

    split_pipe := strings.Split(string(save), "|")
    for i := 1; i < len(split_pipe); i++ {
        split_quote := strings.Split(split_pipe[i], ",")
        id := split_quote[0]
        var tab_tmp []string
        if (split_quote[2] == "LOW") || (split_quote[2] == "MEDIUM") || (split_quote[2] == "HIGH") {
            emergency_tmp := []string{split_quote[0], split_quote[1], split_quote[2]} // id,address,lvl
            reSendEmergency = append(reSendEmergency, emergency_tmp)
        } else {
            for j := 1; j < len(split_quote); j++ {
                tab_tmp = append(tab_tmp, split_quote[j])
            }
            memo[id] = tab_tmp
        }
    }
}

func reSend() {
    for i := 0; i < len(reSendEmergency); i ++ {
        handleJavaRecover(reSendEmergency[i][0], reSendEmergency[i][1], reSendEmergency[i][2])
    }
}

func handleJavaRecover(id string, address string, lvl string) {

    memo[id] = []string{address, lvl}

    msgEmergency := map[string]string{
        "idEmergency" : id,
        "address": address,
    }

    msgFinal := map[string]interface{} {
        "command" : "request",
        "data": msgEmergency,
    }

    fmt.Println("resend old emergency...")
    go broadcastInit(msgFinal ,address, id, lvl)
}

func updateData() {
    max := maxId
    for id, tokens := range memo {

        id_int, err := strconv.Atoi(id)
        check(err)
        if id_int > max {
            max = id_int
        }
        repartiteur[id] = make(chan []string)
        for _, token := range tokens {
            tokenAction = append(tokenAction, token)
        }
    }
    maxId = max + 1
}

//---------------------------------------------------------------------------------

func handleJavaClient(w http.ResponseWriter, r *http.Request) {
    lvl := r.FormValue("emergencyLevel")
    address := r.FormValue("address")
    //service := r.FormValue("service")

    id := getId()
    memo[strconv.Itoa(id)] = []string{address, lvl}

    msgEmergency := map[string]string{
        "idEmergency" : strconv.Itoa(id),
        "address": address,
    }

    msgFinal := map[string]interface{} {
        "command" : "request",
        "data": msgEmergency,
    }

    fmt.Println("Receive emergency...")
    go broadcastInit(msgFinal ,address, strconv.Itoa(id), lvl)
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

func broadcastInit(msg map[string]interface{}, address string, id string, lvlEmergency string) { // id est aussi dans map...
    // init une liste de token potentiel pour l'intervention
    tokensPerimeter := spot(address, 2)
    //fmt.Println(tokensPerimeter[0])
    //memo[id] = tokensPerimeter    // j'ajoute! je n'enleve pas !
    repartiteur[id] = make(chan []string)

    numberNecessary := 0
    switch lvlEmergency {
    case "LOW" :
        numberNecessary = 1
    case "MEDIUM" :
        numberNecessary = 2
    case "HIGH" :
        numberNecessary = 3
    }

    go listenResponse(id, numberNecessary)
    go sendAndroids(tokensPerimeter, msg)

}

func listenResponse(id string, numberNecessary int) {
    c := repartiteur[id]
    inCharge := []string{}
    for {
        rep := <- c
        switch rep[1] {
        case "OK" :
            inCharge = append(inCharge, rep[0])
            tokenAction = append(tokenAction, rep[0])
            t := []string{rep[0]}
            r := map[string]string{
                "msg" : "go go go",
            }
            rf := map[string]interface{}{ 
                "command" : "confirmEmergency",
                "data" : r,
            }
            sendAndroids(t, rf)
            fmt.Println("vehicul accept")
        case "KO" :
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
    // find vehicul in perimeter of address
    return catchGPS(10)
}

func catchGPS(n int) []string {
    var v map[string]interface{}
    if err := firebase.Value(&v); err != nil {
        log.Fatal(err)
    }

    var tokens = make([]string, n)
    acc := 0
    for k := range v {
        token_tmp := v[k].(map[string]interface{})["token"].(string)
        if stringInSlice(token_tmp, tokenAction) {
            continue
        }
        tokens[acc] = token_tmp
        acc = acc + 1
        if acc == n {
            break
        }
    }
    return tokens
}

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func main() {

    if len(os.Args) == 2 {
        if os.Args[1] == "recover" {
            fmt.Println("recover!")
            recover()
            updateData()
            reSend()
        }
    }

    go charge()

    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/android", handleAndroidClient).Methods("POST")
    router.HandleFunc("/java", handleJavaClient).Methods("POST")

    //catchGPS(1) // a virer

    fmt.Println("listening...")
    go saveData()
    //err := http.ListenAndServe(":"+os.Getenv("PORT"), router)
    err := http.ListenAndServe(":1234", router)
     if err != nil {
        panic(err)
    }
}