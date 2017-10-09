package main

import (
    "strconv"
    "time"
    "sort"
)


type TimestampSerie struct {
	T 	int 
	Nb 	int
}


type LabelSerie struct {
	Label 	string 
	Nb 		int
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


func MonthToTimestamp(y int, m int) int {
	then := time.Date(y, time.Month(m), 1, 0, 0, 0, 0, time.UTC)
	return int(then.Unix())	
}

func YearToTimestamp(y int) int {
	return MonthToTimestamp(y, 1)
}

func TimestampToYear(t int) string {
	date := time.Unix(int64(t), 0)
	return strconv.Itoa(date.Year())
}

func TimestampToMounth(t int) string {
	date := time.Unix(int64(t), 0)
	year := date.Year()
	month := int(date.Month())
	frenchMonth := monthFrench(month)
	return strconv.Itoa(year) + "-" + frenchMonth
}


func RangeYear(start int, end int) map[int]int {
	var m map[int]int
	m = make(map[int]int)

	for i := start; i <= end; i++ {
		timestamp := YearToTimestamp(i)
		m[timestamp] = 0
	}

	return m
}


func RangeLastMonth(startY int, startM int, nbMonth int) map[int]int {
	var final map[int]int
	final = make(map[int]int)
	compteur := 0

	for y := startY; y > 0; y-- {
		for m := startM; m > 0; m-- {
			timestamp := MonthToTimestamp(y, m)
			final[timestamp] = 0			
			compteur += 1
			if(compteur == nbMonth) {
				return final
			}
		}
		startM = 12
	}

	return final
}

func ConvertTimeMapToArray(m map[int]int) []TimestampSerie {
	var timesSerie []TimestampSerie

	for k,v := range m {
		var t TimestampSerie
		t.T = k
		t.Nb = v
		timesSerie = append(timesSerie, t)
	}

	sort.Slice(timesSerie, func(i, j int) bool {
	    return timesSerie[i].T > timesSerie[j].T
	})	

	return timesSerie
}


func StatSujetsByYear(sujets []Sujet) []LabelSerie {
	current := time.Now()
	statsDict := RangeYear(2004, current.Year())
	var timesSerie []TimestampSerie
	var labelSerie []LabelSerie

	for _,v := range sujets {
		key := YearToTimestamp(v.Initialised_at.Year())
		if _, ok:= statsDict[key]; ok {
			statsDict[key] += 1
		}		
	}

	timesSerie = ConvertTimeMapToArray(statsDict)
	for _,v := range timesSerie {
		var t LabelSerie
		t.Label = TimestampToYear(v.T)
		t.Nb = v.Nb
		labelSerie = append(labelSerie, t)
	}

	return labelSerie
}


func StatSujetsByLastMouth(sujets []Sujet) []LabelSerie {
	current := time.Now()
	statsDict := RangeLastMonth(current.Year(), int(current.Month()), 12)
	var timesSerie []TimestampSerie
	var labelSerie []LabelSerie

	for _,v := range sujets {
		y := v.Initialised_at.Year()
		m := int(v.Initialised_at.Month())
		key := MonthToTimestamp(y, m)
		if _, ok:= statsDict[key]; ok {
			statsDict[key] += 1
		}		
	}

	timesSerie = ConvertTimeMapToArray(statsDict)
	for _,v := range timesSerie {
		var t LabelSerie
		t.Label = TimestampToMounth(v.T)
		t.Nb = v.Nb
		labelSerie = append(labelSerie, t)
	}

	return labelSerie
}

