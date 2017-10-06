package main

import (
    "strconv"
    "time"
    "sort"
)

type SerieStatItem struct {
	label string
	nb int
}

func monthFrench(x int) string {
	var m []string
	m = append(m, "Jan")
	m = append(m, "Fev")
	m = append(m, "Mar")
	m = append(m, "Avr")
	m = append(m, "Mai")
	m = append(m, "Juin")
	m = append(m, "Jui")
	m = append(m, "Aou")
	m = append(m, "Sep")
	m = append(m, "Oct")
	m = append(m, "Nov")
	m = append(m, "Dec")
	return m[x - 1]
}


func RangeYear(start int, end int) map[string]int {
	var m map[string]int
	m = make(map[string]int)

	for i := start; i <= end; i++ {
		m[strconv.Itoa(i)] = 0
	}

	return m
}


func RangeLastMonth(startY int, startM int, nbMonth int) map[string]int {
	var final map[string]int
	final = make(map[string]int)
	compteur := 0

	for y := startY; y > 0; y-- {
		for m := startM; m > 0; m-- {
			final[strconv.Itoa(y) + "-" + monthFrench(m)] = 0
			compteur += 1
			if(compteur == 12) {
				return final
			}
		}
		startM = 12
	}

	return final
}


func StatSujetsByYear(sujets []Sujet) []SerieStatItem {
	current := time.Now()
	stats := RangeYear(2004, current.Year())

	for _,v := range sujets {
		year := strconv.Itoa(v.Initialised_at.Year())
		if _, ok:= stats[year]; ok {
			stats[year] += 1
		}		
	}

	return ConvertMapToSerie(stats, true)
}


func StatSujetsByLastMouth(sujets []Sujet) []SerieStatItem {
	current := time.Now()
	stats := RangeLastMonth(current.Year(), int(current.Month()), 12)

	for _,v := range sujets {
		year := v.Initialised_at.Year()
		month := int(v.Initialised_at.Month())
		key := strconv.Itoa(year) + "-" + monthFrench(month)
		if _, ok:= stats[key]; ok {
			stats[key] += 1
		}		
	}

	return ConvertMapToSerie(stats, false)
}


func ConvertMapToSerie(m map[string]int, sortByKey bool) []SerieStatItem {
	var conversion []SerieStatItem

	for k,v := range m {
		var t SerieStatItem
		t.label = k
		t.nb = v
		conversion = append(conversion, t)
	}

	if(sortByKey){
		sort.Slice(conversion, func(i, j int) bool {
		    return conversion[i].label > conversion[j].label
		})		
	}

	return conversion
}