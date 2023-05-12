package server

import (
	"gin-exercise/pkg/product/broker"
	"gin-exercise/pkg/product/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ProdReq struct {
	Price       float32 `json:"price" binding:"required"`
	Description string  `json:"description" binding:"required"`
}

func postHandler(context *gin.Context, producer *broker.ProductEventProducer, id *string) {
	req := ProdReq{}
	err := context.BindJSON(&req)
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusBadRequest,
			map[string]string{"error": "malformed json"})
		return
	}

	producer.SendCreationRequest(id, req.Price, req.Description)
	context.Status(http.StatusAccepted)
}

func getHandler(context *gin.Context, repo *db.Repository) {
	id := context.Param("id")

	prod, found := repo.Find(id)
	if found == false {
		context.AbortWithStatusJSON(http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	context.JSON(http.StatusOK, prod)
}

func getManyHandler(context *gin.Context, repo *db.Repository) {
	limitQueryParam := context.Query("limit")
	if limitQueryParam == "" {
		limitQueryParam = "5"
	}
	limit, err := strconv.Atoi(limitQueryParam)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	res, getErr := repo.GetMany(limit)
	if getErr != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, res)
}

func deleteHandler(context *gin.Context, producer *broker.ProductEventProducer) {
	id := context.Param("id")

	producer.SendDeletionRequest(id)
	context.Status(http.StatusAccepted)
}

func SetupRoutes() *gin.Engine {
	database := db.NewProductsDatabase()
	producer, err := broker.NewEventProducer()
	if err != nil {
		panic("unable to create the event producer")
	}

	router := gin.Default()
	router.POST("/product/:id", func(context *gin.Context) {
		idParam := context.Param("id")
		postHandler(context, producer, &idParam)
	})

	router.POST("/product", func(context *gin.Context) {
		postHandler(context, producer, nil)
	})

	router.GET("/product/:id", func(context *gin.Context) {
		getHandler(context, database)
	})

	router.GET("/product", func(context *gin.Context) {
		getManyHandler(context, database)
	})

	router.DELETE("/product/:id", func(context *gin.Context) {
		deleteHandler(context, producer)
	})
	return router
}
