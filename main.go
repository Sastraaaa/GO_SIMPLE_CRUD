package main

import (
	"database/sql"
	"go-crud/handler"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func main() {
	dbHost:="localhost"
	dbUser:="root"
	dbPass:=""
	dbPort:="3306"
	dbName:="go_crud"
	dsn:=dbUser + `:` + dbPass + `@tcp(` + dbHost + `:` + dbPort +`)/` + dbName

	db,err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	handler := handler.InitArticle(db)
	echoServer := echo.New()

	echoServer.GET("/articles",handler.FetchArticles)
	echoServer.POST("/articles", handler.Insert)
	echoServer.GET("/articles/:id", handler.Get)
	echoServer.DELETE("/articles/:id", handler.Delete)
	echoServer.PUT("/articles/:id", handler.Update)

	echoServer.Start(":8080")

}