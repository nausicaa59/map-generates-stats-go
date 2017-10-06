package main

import (
  "fmt"
  "time"
  /*"encoding/json"
  "io/ioutil"*/
)

type Fragment struct {
    Start int
    End   int
}


func main() {
    bdd := Impl{}
    bdd.InitDB()
    bdd.InitSchema()
    pseudos := bdd.GetAllPseudo()
    pseudos = pseudos[40500:40510]
    worker(pseudos, bdd)
}


func worker(pseudos []string, bdd Impl) {
    nb := len(pseudos)
    for i := 0; i < nb; i++ {
        start := time.Now()
        infos := bdd.GetAuteurByPseudo(pseudos[i])
        sujets := bdd.GetSujetByAuteur(int(infos.ID))
        sujetsByYear := StatSujetsByYear(sujets)
        sujetsByLastMounth := StatSujetsByLastMouth(sujets)
        analyseTextuelSujets(sujets)
        end := time.Now()
        fmt.Println(end.Sub(start), infos.Pseudo, len(sujets), len(sujetsByYear), sujetsByLastMounth)
    }
}


/*
CalcDistPseudo(pseudos[i], pseudos)
nb := len(pseudos)
for i := 0; i < nb; i++ {
    start := time.Now()
    CalcDistPseudo(pseudos[i], pseudos)
    end := time.Now()
    fmt.Println(end.Sub(start))
}

func fragmenter(t []string, nb int) []Fragment {
    var final []Fragment
    compl := len(t) % nb
    tranche := (len(t) - compl) / nb

    for i := 1; i <= nb; i++ {
        tempo := Fragment{}
        tempo.Start = tranche * (i - 1)
        tempo.End = tranche * i
        final = append(final, tempo)
    }

    final[nb-1].End += compl
    return final
}
fmt.Println(end.Sub(start))

end := time.Now()
err = ioutil.WriteFile(a.Pseudo + ".json", j, 0644)
if(err != nil) {
  panic("Erreur lors de l'écriture")
}
*/