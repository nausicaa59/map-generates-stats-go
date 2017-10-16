package main

import (
  "fmt"
  "encoding/json"
  "io/ioutil"
  "os"
)


func WriteOutPutProfils(p PseudoProfils) bool {
    output, err := json.Marshal(p)
    if(err != nil) {
        fmt.Println("erreur lors de la sérialisation de ", p.Infos.Pseudo, ":", err)
        return false
    }

    path := GeneratePathFile("profils/", p.Infos.Pseudo, 3)
    err = ioutil.WriteFile(path, output, 0644)
    if(err != nil) {
        fmt.Println("erreur lors de l'écriture de ", p.Infos.Pseudo, ":", err)
        return false
    }

    return true
}

func WriteOutPutSimilaire(pseudo string, s []Similaire) bool {
    output, err := json.Marshal(s)
    if(err != nil) {
        fmt.Println("erreur lors de la sérialisation de ", pseudo, ":", err)
        return false
    }

    path := GeneratePathFile("similaires/", pseudo, 3)
    err = ioutil.WriteFile(path, output, 0644)
    if(err != nil) {
        fmt.Println("erreur lors de l'écriture de ", pseudo, ":", err)
        return false
    }

    return true
}


func GeneratePathFile(subFolder string, pseudo string, nbLvl int) string {
    path := "output/" + subFolder

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
