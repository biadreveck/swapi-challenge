package api

import (
	"b2w/swapi-challenge/api/handler"
	"b2w/swapi-challenge/api/middleware"
	"b2w/swapi-challenge/domain/entity/planet"

	"github.com/gin-gonic/gin"
)

func SetupRouter(pManager planet.Manager) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.Cors())

	handler.CreatePlanetRoutes(router, pManager)

	return router
}
