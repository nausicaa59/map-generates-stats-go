package main

import (
	"sort"
    "github.com/xrash/smetrics"
    "regexp"
    "strings"
	"os"
	"io"
	"bufio"
	"encoding/csv"
	"fmt"
)

type DistString struct {
	S string
	Dist float64
}

type WordOccurence struct {
	Word string
	Nb int
}


func getExcludeWordFile(path string) []string {
	csvFile, _ := os.Open(path)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	words := []string{}

	for {
		line, error := reader.Read()
		
		if error == io.EOF {
			break
		}

		if error != nil {
			fmt.Println(error)
		}

		if len(line) > 0 {
			words = append(words, line[0])
		}
	}

	return words
}


func CalcDistPseudo(c string, pseudos []string) []DistString {
	var dists []DistString

	for _, v := range pseudos {
		x := DistString{}
		x.S = v
		x.Dist = smetrics.Jaro(c, v)
		dists = append(dists, x)
	}

	sort.Slice(dists, func(i, j int) bool {
	    return dists[i].Dist > dists[j].Dist
	})

	return dists[0:60]
}


func CleanUrls(urls []string) []string {
	var finals []string
	r, _ := regexp.Compile("http://www.jeuxvideo.com/forums/[0-9]*-[0-9]*-[0-9]*-[0-9]*-[0-9]*-[0-9]*-[0-9]*-(.*).htm")
	
	for _, v := range urls {
		found := r.FindStringSubmatch(v)
		if found != nil {
			finals = append(finals, found[1])
		}
	}

	return finals
}

func ConvertMapWordToSlice(m map[string]int) []WordOccurence {
	var conversion []WordOccurence

	for k,v := range m {
		var t WordOccurence
		t.Word = k
		t.Nb = v
		conversion = append(conversion, t)
	}

	return conversion
}

func AnalyseUrls(urls []string) []WordOccurence {
	var words map[string]int
	words = make(map[string]int)
	urlsCleans := CleanUrls(urls)

	for _,v := range urlsCleans {
		for _, w := range strings.Split(v, "-") {
			if _, ok:= words[w]; ok {
				words[w] += 1
			} else {
				words[w] = 1
			}
		}
	}

	final := ConvertMapWordToSlice(words)
	sort.Slice(final, func(i, j int) bool {
	    return final[i].Nb > final[j].Nb
	})

	return final
}

func cleanArrayWords(src []WordOccurence, exclude []string) []WordOccurence {
	var clean []WordOccurence

	for _,v := range src {
		isValide := true
		for _,e := range exclude {
			if v.Word == e {
				isValide = false
			}
		}

		if(isValide) {
			clean = append(clean, v)
		}
	}

	return clean
}


func analyseTextuelSujets(sujets []Sujet, exclude []string) []WordOccurence {
	var urls []string

	for _,v := range sujets {
		urls = append(urls, v.Url)
	}

	final := AnalyseUrls(urls)
	return cleanArrayWords(final, exclude)
}