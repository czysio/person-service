package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/czysio/person-service/config"
	"github.com/czysio/person-service/controllers"
	dbConn "github.com/czysio/person-service/db/sqlc"
	"github.com/czysio/person-service/routes"

	_ "github.com/lib/pq"
)

var (
	server *gin.Engine
	db     *dbConn.Queries
	ctx    context.Context

	PersonController controllers.PersonController
	PersonRoutes     routes.PersonRoutes
)

func init() {
	ctx = context.TODO()
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	conn, err := sql.Open(config.PostgreDriver, config.PostgresSource)
	if err != nil {
		log.Fatalf("could not connect to postgres database: %v", err)
	}

	db = dbConn.New(conn)

	fmt.Println("PostgreSQL connected successfully...")

	PersonController = *controllers.NewPersonController(db, ctx)
	PersonRoutes = routes.NewRoutePerson(PersonController)

	server = gin.Default()
}

func main() {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{config.Origin}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")

	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Welcome!"})
	})

	PersonRoutes.PersonRoute(router)
	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": fmt.Sprintf("Route %s not found", ctx.Request.URL)})
	})
	log.Fatal(server.Run(":" + config.ServerPort))
}
