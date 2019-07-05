package main

import (
	"log"

	"github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)

type Movie struct {
    ID int `gorm:"AUTO_INCREMENT"`
    Title string `gorm:"type:varchar(100);unique_index"`
}

func checkError(err error) {
    if err != nil {
        panic(err)
    }
}

func main() {
	db, err := gorm.Open("mysql", "root:@/test?charset=utf8")
	defer db.Close()
    checkError(err)

	var movie Movie
	var movies []Movie
	db.First(&movie, 30170448)
	log.Println(movie)
	log.Println(movie.ID, movie.Title)

	db.Order("id").Limit(3).Find(&movies)
	log.Println(movies)
	log.Println(movies[0].ID)

	db.Order("id desc").Limit(3).Offset(1).Find(&movies)
	log.Println(movies)
	db.Select("title").Find(&movies, 30170448)
	log.Println(movies)
	db.Select("title").First(&movies, "title = ?", "四个春天")
    log.Println(movie)

	var count int64
	db.Where("id = ?", 30170448).Or("title = ?", "四个春天").Find(&movies).Count(&count)
	log.Println(count)
}
