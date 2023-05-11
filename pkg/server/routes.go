package server

import (
	"encoding/json"
	"fmt"
	"gin-exercise/pkg/product"
	"gin-exercise/pkg/product/db"
	"gin-exercise/pkg/util"
	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
)

type ProdReq struct {
	Price       float32 `json:"price" binding:"required"`
	Description string  `json:"description" binding:"required"`
}

func postKafkaHandler(context *gin.Context, id *string) {
	if id == nil {
		id = util.GetPtr(db.GenerateNewId())
	}
	req := ProdReq{}
	err := context.BindJSON(&req)
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusBadRequest,
			map[string]string{"error": "malformed json"})
		return
	}

	config := sarama.NewConfig()
	p, err := sarama.NewAsyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	body, _ := json.Marshal(req)

	msg := sarama.ProducerMessage{
		Topic: "products",
		Key:   sarama.StringEncoder(*id),
		Value: sarama.StringEncoder(body),
	}

	p.Input() <- &msg
}

func postHandler(context *gin.Context, repository *db.Repository, id *string) {
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

	router := gin.Default()
	router.POST("/product/:id", func(context *gin.Context) {
		idParam := context.Param("id")
		postHandler(context, database, &idParam)
	})

	router.POST("/product", func(context *gin.Context) {
		postKafkaHandler(context, nil)
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
