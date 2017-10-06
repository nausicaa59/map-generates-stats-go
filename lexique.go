package main

import (
	"sort"
    "github.com/xrash/smetrics"
    "regexp"
    "strings"
)

type distString struct {
	s string
	dist float64
}

type WordOccurence struct {
	word string
	nb int
}


func CalcDistPseudo(c string, pseudos []string) []distString {
	var dists []distString

	for _, v := range pseudos {
		x := distString{}
		x.s = v
		x.dist = smetrics.Jaro(c, v)
		dists = append(dists, x)
	}

	sort.Slice(dists, func(i, j int) bool {
	    return dists[i].dist > dists[j].dist
	})

	return dists[0:30]
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
		t.word = k
		t.nb = v
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
	    return final[i].nb > final[j].nb
	})

	return final
}


func analyseTextuelSujets(sujets []Sujet) []WordOccurence {
	var urls []string

	for _,v := range sujets {
		urls = append(urls, v.Url)
	}

	return AnalyseUrls(urls)
}