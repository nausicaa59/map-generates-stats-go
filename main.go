package main

import (
  "fmt"
  "time"
  "os"
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
    Banni       uint
}





func main() {
    conf := LoadConf()
    start := time.Now() 
    bdd := Impl{}
    bdd.InitDB(conf)
    bdd.InitSchema()

    argsWithProg := os.Args[1:]
    if len(argsWithProg) == 0 {
        fmt.Println("Aucun mode choisi, fin du programme")
        os.Exit(0)
    }

    switch argsWithProg[0] {
        case "-profils":
            pseudos := bdd.GetAllPseudo()
            GenerateProfils(pseudos, bdd)
        case "-similaires":
            pseudos := bdd.GetAllPseudo()
            GenerateSimilaires(pseudos, bdd)
        default:
            fmt.Println("Aucun mode valide choisi, fin du programme")
            os.Exit(0)
    }

    end := time.Now()
    fmt.Println(argsWithProg, len(argsWithProg), end.Sub(start))
}


func GenerateProfils(pseudos []string, bdd Impl) {
    nb := len(pseudos)
    excludeWord := getExcludeWordFile("input/excludeWord.csv")

    for i := 0; i < nb; i++ {
        start := time.Now() 
        pourc := (float64(i)/float64(nb))*100

        p := PseudoProfils{}        
        p.Infos             = bdd.GetAuteurByPseudo(pseudos[i])
        p.Sujets            = bdd.GetSujetByAuteur(int(p.Infos.ID))
        p.SujetByYear       = StatSujetsByYear(p.Sujets)
        p.SujetByLastMouth  = StatSujetsByLastMouth(p.Sujets)
        p.AnalyseTextuel    = analyseTextuelSujets(p.Sujets, excludeWord)
        WriteOutPutProfils(p)

        end := time.Now()
        fmt.Println(pourc, end.Sub(start), p.Infos.Pseudo)
    }
}


func GenerateSimilaires(pseudos []string, bdd Impl) {
    nb := len(pseudos)
    for i := 0; i < nb; i++ {
        start := time.Now() 
        pourc := (float64(i)/float64(nb))*100
        
        similaires := getSimilaires(pseudos[i], pseudos, bdd)
        WriteOutPutSimilaire(pseudos[i], similaires)
        
        end := time.Now()
        fmt.Println(pourc, end.Sub(start), pseudos[i])
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
        tempo.Banni         = infos.Banni
        tempo.Pourc         = v.Dist * 100
        s = append(s, tempo)
    }

    return s
}
