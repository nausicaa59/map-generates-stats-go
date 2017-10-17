package main

import (
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "time"
    "fmt"
    "strconv"
    "sort"
)


type Impl struct {
	DB *gorm.DB
}


type Auteur struct {
  ID uint 						`gorm:"primary_key"`
  Pseudo string  				`gorm:"column:pseudo"`
  Created_at time.Time  		`gorm:"column:created_at"`
  Updated_at time.Time 			`gorm:"column:updated_at"`
  Cheked_profil uint			`gorm:"column:cheked_profil"`
  Pays string					`gorm:"column:pays"`
  Nb_messages uint				`gorm:"column:nb_messages"`
  Img_lien string				`gorm:"column:img_lien"`
  Nb_relation uint				`gorm:"column:nb_relation"`
  Banni uint					`gorm:"column:banni"`
  Date_inscription time.Time 	`gorm:"column:date_inscription"`
  Coord_X float64 				`gorm:"column:coord_X"`
  Coord_Y float64 				`gorm:"column:coord_Y"`
}


type Sujet struct {
  ID uint 						`gorm:"primary_key"`
  Created_at time.Time  		`gorm:"column:created_at"`
  Updated_at time.Time 			`gorm:"column:updated_at"`
  Parcoured uint				`gorm:"column:parcoured"`
  Url string					`gorm:"column:url"`
  Title string					`gorm:"column:title"`
  Auteur uint					`gorm:"column:auteur"`
  Nb_reponses uint				`gorm:"column:nb_reponses"`
  Initialised_at time.Time		`gorm:"column:initialised_at"`
}


func (i *Impl) InitDB(conf Conf) {
	var err error
	c := conf.DbUsername + ":" + conf.DbPassword + "@/" + conf.DbName + "?charset=utf8&parseTime=True&loc=Local"
	i.DB, err = gorm.Open("mysql", c)
	if err != nil {
		fmt.Println(err)
		panic("Got error when connect database, the error is '%v'")
	}
}


func (i *Impl) InitSchema() {
	i.DB.AutoMigrate(&Auteur{}, &Sujet{})
}


func (i *Impl) Close() {
	fmt.Println("Fermeture")
	i.DB.Close()
}


func (i *Impl) GetAllPseudo() []string {
	var names []string
	i.DB.Model(&Auteur{}).Pluck("pseudo", &names)
	fmt.Println(len(names))

	var final []string
	for _,v := range names {
		if(v != "Pseudosupprimé") {
			final = append(final, v)
		}
	}

	return final	
}


func (i *Impl) GetAuteur(id int) Auteur {
	var a Auteur
	i.DB.Find(&a, id)
	return a
}


func (i *Impl) GetAuteurByPseudo(id string) Auteur {
	var a Auteur
	i.DB.Where("pseudo = ?", id).First(&a)
	return a
}


func (i *Impl) GetSujetByAuteur(id int) []Sujet {
	var a []Sujet
	i.DB.Where("auteur = ?", id).Find(&a)
	return a
}


func (i *Impl) CountAuteur() int {
	var count int
	i.DB.Table("auteurs").Count(&count)
	return count
}

func (i *Impl) CountSujets() int {
	var count int
	i.DB.Table("sujets").Count(&count)
	return count
}

func (i *Impl) CountNbReponseSujets() int {
	type Result struct {
	    Total int
	}

	var results []Result
	i.DB.Table("sujets").Select("sum(nb_reponses) as total").Scan(&results)

	return results[0].Total
}

func (i *Impl) CountNbReponseByYear() []LabelSerie {
	var results []LabelSerie
	i.DB.Table("sujets").Select("Year(initialised_at) as label, sum(nb_reponses) as nb").Group("label").Scan(&results)
	
	sort.Slice(results, func(i, j int) bool {
		yearI, _ := strconv.Atoi(results[i].Label)
		yearJ, _ := strconv.Atoi(results[j].Label)
	    return yearI > yearJ
	})

	return results
}

func (i *Impl) CountNbReponseByLastMouth() []LabelSerie {
	var results []LabelSerie
	i.DB.Table("sujets").Select("DATE_FORMAT(initialised_at, '%Y-%m') as label, sum(nb_reponses) as nb").Group("label").Scan(&results)
	return results[len(results)-12:]
}

func (i *Impl) TopSujets() []Sujet {
	var s []Sujet
	i.DB.Order("nb_reponses desc").Limit(30).Find(&s)
	return s	
}

func (i *Impl) TopAuteurs() []Auteur {
	var a []Auteur
	i.DB.Order("nb_messages desc").Limit(30).Find(&a)
	return a	
}

func (i *Impl) GetAllUrlSujets() []string {
	var url []string
	i.DB.Model(&Sujet{}).Pluck("url", &url)
	return url
}