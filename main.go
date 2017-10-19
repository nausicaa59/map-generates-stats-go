package main

import (
  "fmt"
  "time"
  "os"
  "sort"
  "bufio"
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

type StatGenerale struct {
    NbPseudo            uint
    NbSujet             uint
    NbReponse           uint
    NbRepSujet          float64
    ReponseByYear       []LabelSerie
    ReponseByLastMouth  []LabelSerie
    Sujets              []Sujet
    Auteurs             []Auteur
    AnalyseTextuel      []WordOccurence
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
        case "-statsgeneral":
            GenerateStatGeneral(bdd)            
        default:
            fmt.Println("Aucun mode valide choisi, fin du programme")
            os.Exit(0)
    }

    end := time.Now()
    fmt.Println(argsWithProg, len(argsWithProg), end.Sub(start))
}


func GenerateStatGeneral(bdd Impl) {
    var stats StatGenerale
    stats.NbPseudo  = uint(bdd.CountAuteur())
    stats.NbSujet   = uint(bdd.CountSujets())
    stats.NbReponse = uint(bdd.CountNbReponseSujets())
    stats.NbRepSujet = float64(float64(stats.NbReponse) / float64(stats.NbSujet))
    stats.ReponseByYear = bdd.CountNbReponseByYear()  
    stats.ReponseByLastMouth = bdd.CountNbReponseByLastMouth()
    stats.Sujets = bdd.TopSujets()
    stats.Auteurs = bdd.TopAuteurs() 
    stats.AnalyseTextuel = GenerateAllWords(bdd) 
    WriteOutPutStatGeneral(stats)
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
        if FileExiste("similaires/", pseudos[i], 3) {
            continue
        }

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

func GenerateAllWords(bdd Impl) []WordOccurence{
    var words map[string]int
    words = make(map[string]int)
    excludeWord := getExcludeWordFile("input/excludeWord.csv")
    compteur := 0

    file, err := os.Open("input/sujets.csv")
    if err != nil {
        panic("erreur")
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        var t []string
        t = append(t, scanner.Text())
        GenrateMapFromUrl(words, t)
        compteur += 1
        
        if compteur % 1000000 == 0 {
            fmt.Println(compteur)
        }
    }

    if err := scanner.Err(); err != nil {
        panic("erreur")
    }


    final := ConvertMapWordToSlice(words)
    final = cleanArrayWords(final, excludeWord)
    sort.Slice(final, func(i, j int) bool {
        return final[i].Nb > final[j].Nb
    })

    return final[0:100]
}
