package server

import (
	"gin-exercise/pkg/product"
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

	producer.SendEvent(id, req.Price, req.Description)
	context.Status(http.StatusAccepted)
}

func _postHandler(context *gin.Context, repository *db.Repository, id *string) {
	req := ProdReq{}
	err := context.BindJSON(&req)
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusBadRequest,
			map[string]string{"error": "malformed json"})
		return
	}

	prod, repErr := repository.Store(id, req.Price, req.Description)
	if repErr != nil {
		errorPayload := map[string]string{"error": repErr.Error()}
		httpStatus := http.StatusInternalServerError
		if repErr.Error() == product.IdTakenError {
			httpStatus = http.StatusBadRequest
		}
		context.AbortWithStatusJSON(httpStatus, errorPayload)
		return
	}

	context.JSON(http.StatusOK, prod)
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

func deleteHandler(context *gin.Context, repo *db.Repository) {
	id := context.Param("id")

	found := repo.Delete(id)
	if found {
		context.Status(http.StatusNoContent)
		return
	}
	context.AbortWithStatusJSON(http.StatusNotFound, map[string]string{"error": "not found"})
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
		deleteHandler(context, database)
	})
	return router
}
