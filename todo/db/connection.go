package db

import (
	"todo/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB	

func Connect(){
	dns:="host=localhost user=postgres password=root dbname=test port=5555 sslmode=disable TimeZone=Europe/Minsk"
	connection,err:=gorm.Open(postgres.Open(dns),&gorm.Config{})
	if err!=nil{
		panic(err)
	}
	DB = connection
	connection.AutoMigrate(&models.User{})
	connection.AutoMigrate(&models.Tasks{})
}