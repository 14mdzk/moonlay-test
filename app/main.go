package main

import (
	"fmt"
	"log"
	_domain "moonlay-test/domain"
	_listHandler "moonlay-test/list/delivery"
	_listRepository "moonlay-test/list/repository"
	_listUsecase "moonlay-test/list/usecase"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Run service on DEBUG mode")
	}
}

func main() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("database.host"), viper.GetString("database.port"), viper.GetString("database.user"), viper.GetString("database.pass"), viper.GetString("database.dbname"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Println(err)
	}

	var listTable _domain.List
	db.AutoMigrate(&listTable)

	sqlDB, err := db.DB()
	sqlDB.Ping()
	if err != nil {
		log.Println(err)
	}

	defer func() {
		err := sqlDB.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	e := echo.New()
	e.Use(middleware.CORS())

	listRepo := _listRepository.NewListSQLRepository(db)
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	listUsecase := _listUsecase.NewListUsecase(listRepo, timeoutContext)
	_listHandler.NewListHttpHandler(e, listUsecase)

	log.Fatal(e.Start(viper.GetString("server.address")))
}
