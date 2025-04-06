package repository

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
	"servers_status/model"
	"strconv"
)

var (
	DB  *gorm.DB
	err error
)

func Init() {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s",
		os.Getenv("DB_HOST"), port, os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	DB, err = gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "bot_servers.", // schema name
			SingularTable: false,
		}})
	if err != nil {
		panic(err)
	}
	fmt.Println("Database connection is succesful")

	err = DB.AutoMigrate(
		&model.Channels{},
		&model.Servers{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Migration successful")

}
