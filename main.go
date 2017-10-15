package main

import (
  "fmt"
  "time"
)

type PseudoProfils struct {
    Infos                Auteur
    AnalyseTextuel      []WordOccurence
    SujetByYear         []LabelSerie
    SujetByLastMouth    []LabelSerie
    Sujets              []Sujet
    Similaires          []Similaire
}

type Similaire struct {
    Pseudo      string
    ID          uint
    Nb_messages uint
    Pourc       float64
    Img_lien    string
}

func main() {
    conf := LoadConf()
    start := time.Now() 
    bdd := Impl{}
    bdd.InitDB(conf)
    bdd.InitSchema()
    /*pseudos := bdd.GetAllPseudo()
    GenerateProfils(pseudos, bdd)*/
    end := time.Now()
    fmt.Println(end.Sub(start))
}


func GenerateProfils(pseudos []string, bdd Impl) {
    nb := len(pseudos)
    excludeWord := getExcludeWordFile("input/excludeWord.csv")

    for i := 0; i < nb; i++ {
        start := time.Now()        
        p := PseudoProfils{}        
        p.Infos             = bdd.GetAuteurByPseudo(pseudos[i])
        p.Sujets            = bdd.GetSujetByAuteur(int(p.Infos.ID))
        p.SujetByYear       = StatSujetsByYear(p.Sujets)
        p.SujetByLastMouth  = StatSujetsByLastMouth(p.Sujets)
        p.AnalyseTextuel    = analyseTextuelSujets(p.Sujets, excludeWord)
        p.Similaires        = getSimilaires(pseudos[i], pseudos, bdd)
        WriteOutPut(p)
        end := time.Now()
        fmt.Println(end.Sub(start), p.Infos.Pseudo)

        if(i > 20) {
            break
        }
    }
}

func getSimilaires(target string, pseudos []string, bdd Impl) []Similaire {
    var s []Similaire
    pFounds := CalcDistPseudo(target, pseudos)

    for _,v := range pFounds {
        tempo := Similaire{}
        infos := bdd.GetAuteurByPseudo(v.S)
        tempo.Pseudo        = infos.Pseudo
        tempo.ID            = infos.ID
        tempo.Nb_messages   = infos.Nb_messages
        tempo.Img_lien      = infos.Img_lien
        tempo.Pourc         = v.Dist * 100
        s = append(s, tempo)
    }

    return s
}
