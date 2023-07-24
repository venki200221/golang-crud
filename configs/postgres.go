package configs

import (
	"fmt"
	"log"
	"go.io/gorm"
	"gorm.io/driver/postgres"
)

var DB1 *gorm.DB
var err error

type book struct{
	gorm.Model
	Title string `json:"title"`
	Author string  `json:"author"`
}

func postgresConnect(){
	host:="localhost"
	port:="5432"
	dbName:="postgres"
	dbUser:="postgres"
	password:="postgres"
	dsn:=fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable ",
	host,
	port,
	dbUser,
	dbName,
	password,)
	DB1,err=gorm.Open(postgres.Open(dsn),&gorm.Config{})
	DB1.AutoMigrate(&book{})
	if err!=nil{
		log.Fatal("error connecting to database",err)
	}
	fmt.Println("connected to database")
}
