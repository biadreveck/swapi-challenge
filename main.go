package main

import (
	"b2w/swapi-challenge/api/handler"
	"b2w/swapi-challenge/api/middleware"
	"b2w/swapi-challenge/config"
	"b2w/swapi-challenge/domain/entity/planet"
	"b2w/swapi-challenge/infra/database"
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	// Lendo as configurações
	config.ReadConfig()

	// Conectando com o banco de dados
	dbConfig := config.Data.Database
	dbClient, err := database.NewClient(dbConfig)
	if err != nil {
		log.Fatalln(err)
	}

	connectDatabaseClient(dbClient, dbConfig)
	defer disconnectDatabaseClient(dbClient, dbConfig)
	db := dbClient.Database(dbConfig.DBName)

	// Criando os repositórios e gerenciadores
	planetDbRepo := planet.NewMongoRepository(db)
	planetSWApiRepo := planet.NewSWApiRepository()
	planetManager := planet.NewManager(planetDbRepo, planetSWApiRepo)

	// Criando as rotas da API
	router := gin.Default()
	router.Use(middleware.Cors())

	handler.CreatePlanetRoutes(router, planetManager)

	router.Run(config.Data.Server.Address)
}

func connectDatabaseClient(dbClient database.ClientHelper, dbConfig config.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), dbConfig.ConnectionTimeout)
	defer cancel()
	err := dbClient.Connect(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	// Checar se realmente está conectado no banco
	err = dbClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalln(err)
	}
}

func disconnectDatabaseClient(dbClient database.ClientHelper, dbConfig config.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), dbConfig.ConnectionTimeout)
	defer cancel()
	err := dbClient.Disconnect(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}
