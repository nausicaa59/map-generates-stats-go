package main

import (
  "fmt"
  "encoding/json"
  "io/ioutil"
  "os"
)


func FileExist(p string) bool {
    path := GeneratePathFile(p, 3)
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return false
    }

    return true
}


func WriteOutPut(p PseudoProfils) bool {
    output, err := json.Marshal(p)
    if(err != nil) {
        fmt.Println("erreur lors de la sérialisation de ", p.Infos.Pseudo, ":", err)
        return false
    }

    path := GeneratePathFile(p.Infos.Pseudo, 3)
    err = ioutil.WriteFile(path, output, 0644)
    if(err != nil) {
        fmt.Println("erreur lors de l'écriture de ", p.Infos.Pseudo, ":", err)
        return false
    }

    return true
}


func GeneratePathFile(pseudo string, nbLvl int) string {
    path := "output/"

    for i := 0; i < nbLvl; i++ {
        path += string(pseudo[i]) + "/"
      
        if(!folderExists(path)) {
            os.MkdirAll(path, 0777)
        }
    }

    return path + pseudo + ".json"
}


func folderExists(path string) (bool) {
    _, err := os.Stat(path)
    if err == nil { return true }
    if os.IsNotExist(err) { return false }
    return true
}
