package main

import (
  "fmt"
  "encoding/json"
  "io/ioutil"
  "os"
)


type Conf struct {
    DbHost          string
    DbUsername      string
    DbPassword      string
    DbName          string
}


func IsValideConf(c Conf) bool {
    if c.DbHost == "" {
        return false
    }

    if c.DbUsername == "" {
        return false
    }

    if c.DbPassword == "" {
        return false
    }

    if c.DbName == "" {
        return false
    }

    return true
}


func LoadConf() Conf {
    file, e := ioutil.ReadFile("conf.env")
    if e != nil {
        fmt.Printf("File error: %v\n", e)
        os.Exit(1)
    }

    var conf Conf
    json.Unmarshal(file, &conf) 

    if !IsValideConf(conf) {
        fmt.Printf("Configuration invalide")
        os.Exit(1)        
    }

    return conf   
}
