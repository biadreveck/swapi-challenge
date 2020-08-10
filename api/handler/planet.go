package handler

import (
	"b2w/swapi-challenge/domain"
	"net/http"

	"b2w/swapi-challenge/api/presenter"
	"b2w/swapi-challenge/domain/entity/planet"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreatePlanetRoutes(router *gin.Engine, manager planet.Manager) {
	planet := router.Group("/v1/planets")
	{
		planet.POST("", createPlanet(manager))
		planet.GET("", getPlanets(manager))
		planet.GET("/:id", getPlanet(manager))
		planet.DELETE("/:id", deletePlanet(manager))
	}
}

func createPlanet(manager planet.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var addPlanet presenter.AddPlanetCommand
		err := c.BindJSON(&addPlanet)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unexpected JSON format", "params": addPlanet})
			return
		}

		p := addPlanet.ToModel()
		err = manager.Insert(&p)
		if err != nil {
			if err == domain.ErrConflict {
				c.JSON(http.StatusConflict, gin.H{"error": "A planet with specified params already exists", "params": addPlanet})
			} else if err == domain.ErrBadParamInput {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planet input params", "params": addPlanet})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while saving planet on database"})
			}

			return
		}

		c.JSON(http.StatusCreated, gin.H{"data": presenter.NewPlanetResult(p)})
	}
}

func getPlanet(manager planet.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unexpected ID format", "params": idParam})
			return
		}

		p, err := manager.GetById(id)
		if err != nil {
			if err == domain.ErrNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Planet not found", "params": idParam})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while getting planet from database"})
			}

			return
		}

		c.JSON(http.StatusOK, gin.H{"data": presenter.NewPlanetResult(p)})
	}
}

func deletePlanet(manager planet.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unexpected ID format", "params": idParam})
			return
		}

		err = manager.Delete(id)
		if err != nil {
			if err == domain.ErrNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Planet not found", "params": idParam})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while removing planet from database"})
			}

			return
		}

		c.Status(http.StatusNoContent)
	}
}

func getPlanets(manager planet.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")

		var result interface{}
		var err error

		if name == "" {
			var pList []planet.Planet
			pList, err = manager.FindAll()
			result = presenter.NewPlanetResultSlice(pList)
		} else {
			var p planet.Planet
			p, err = manager.GetByName(name)
			result = presenter.NewPlanetResult(p)

			if err != nil && err == domain.ErrNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Planet not found", "params": name})
				return
			}
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching planets from database"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": result})
	}
}
