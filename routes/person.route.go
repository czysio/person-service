package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/czysio/person-service/controllers"
)

type PersonRoutes struct {
	personController controllers.PersonController
}

func NewRoutePerson(personController controllers.PersonController) PersonRoutes {
	return PersonRoutes{personController}
}

func (pc *PersonRoutes) PersonRoute(rg *gin.RouterGroup) {

	router := rg.Group("people")
	router.POST("/", pc.personController.CreatePerson)
	router.GET("/", pc.personController.GetAllPeople)
	router.PATCH("/:personId", pc.personController.UpdatePerson)
	router.GET("/:personId", pc.personController.GetPersonById)
	router.DELETE("/:personId", pc.personController.DeletePersonById)
}