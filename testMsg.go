package main

import (
    "fmt"
    "github.com/NaySoftware/go-fcm"
)

const (
     serverKey = "AIzaSyAM5yN0SNAswN6l6t6DEKv9fLRSeUaliVY"
)

func main() {
    data := map[string]string{
        "UrgenceLvl": "RED",
        "adresse": "91 chemin yolo",
    }

  ids := []string{
      "e6THtaBcNVE:APA91bESyZPEZ19jjMIpSBkry1eKAJCnYeRPsw6Dm_mMUQovH3APX4V-gSxJHHnuFK1OWhcM3dOpNw2h__sRy3HYaY5fqQ--vKwzG43WngO-XGEqO1b_X8aFM7HAioLljQH4M505RR1U",
  }



    c := fcm.NewFcmClient(serverKey)
    c.NewFcmRegIdsMsg(ids, data)
    //c.AppendDevices(xds)

    status, err := c.Send()


    if err == nil {
    status.PrintResults()
    } else {
        fmt.Println(err)
    }

}