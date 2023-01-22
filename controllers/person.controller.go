package controllers

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/czysio/person-service/db/sqlc"
	"github.com/czysio/person-service/schemas"
)

type People interface {
	CreatePerson(*gin.Context)
	UpdatePerson(*gin.Context)
	GetPersonById(*gin.Context)
	GetAllPeople(*gin.Context)
	DeletePersonById(*gin.Context)
}

type PersonController struct {
	db  *db.Queries
	ctx context.Context
}

func NewPersonController(db *db.Queries, ctx context.Context) *PersonController {
	return &PersonController{db, ctx}
}

func (pc *PersonController) CreatePerson(ctx *gin.Context) {
	var payload *schemas.CreatePerson

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	now := time.Now()
	args := &db.CreatePersonParams{
		FirstName:     payload.FirstName,
		Surname:  payload.Surname,
		Email:   payload.Email,
		Nickname:     payload.Nickname,
		CreatedAt: now,
		UpdatedAt: now,
	}

	person, err := pc.db.CreatePerson(ctx, *args)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "person": person})
}

func (pc *PersonController) UpdatePerson(ctx *gin.Context) {
	var payload *schemas.UpdatePerson
	personId := ctx.Param("personId")

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	now := time.Now()
	args := &db.UpdatePersonParams{
		ID:        uuid.MustParse(personId),
		FirstName: sql.NullString{String: payload.FirstName, Valid: payload.FirstName != ""},
		Surname:   sql.NullString{String: payload.Surname, Valid: payload.Surname != ""},
		Email:     sql.NullString{String: payload.Email, Valid: payload.Email != ""},
		Nickname:  sql.NullString{String: payload.Nickname, Valid: payload.Nickname != ""},
		UpdatedAt: sql.NullTime{Time: now, Valid: true},
	}

	person, err := pc.db.UpdatePerson(ctx, *args)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No person with that ID exists"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "person": person})
}

func (pc *PersonController) GetPersonById(ctx *gin.Context) {
	personId := ctx.Param("personId")

	person, err := pc.db.GetPersonById(ctx, uuid.MustParse(personId))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No person with that ID exists"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "person": person})
}

func (pc *PersonController) GetAllPeople(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	args := &db.ListPeopleParams{
		Limit:  int32(intLimit),
		Offset: int32(offset),
	}

	people, err := pc.db.ListPeople(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if people == nil {
		people = []db.Person{}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(people), "data": people})
}

func (pc *PersonController) DeletePersonById(ctx *gin.Context) {
	personId := ctx.Param("personId")

	_, err := pc.db.GetPersonById(ctx, uuid.MustParse(personId))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No person with that ID exists"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	err = pc.db.DeletePerson(ctx, uuid.MustParse(personId))
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"status": "success"})
}