package server

import (
	"errors"
	"gin-exercise/pkg/db"
	"gin-exercise/pkg/product"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type prodReq struct {
	Price       float32 `json:"price" binding:"required"`
	Description string  `json:"description" binding:"required"`
}

func postHandler(context *gin.Context, repository *product.Repository, id *string) {
	req := prodReq{}
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

func getHandler(context *gin.Context, repo *gorm.DB) {
	id := context.Param("id")

	var prdEnt db.ProductEntity
	found := repo.First(&prdEnt, "id = ?", id)

	if errors.Is(found.Error, gorm.ErrRecordNotFound) {
		context.AbortWithStatusJSON(http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}

	if err := found.Error; err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	toProduct, _ := prdEnt.ToProduct()
	context.JSON(http.StatusOK, toProduct)
}

func getManyHandler(context *gin.Context, repo *product.Repository) {
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

func deleteHandler(context *gin.Context, repo *product.Repository) {
	id := context.Param("id")

	found := repo.Delete(id)
	if found {
		context.Status(http.StatusNoContent)
		return
	}
	context.AbortWithStatusJSON(http.StatusNotFound, map[string]string{"error": "not found"})
}

func SetupRoutes() *gin.Engine {
	prodRepo := product.NewProductRepository()
	database := db.ProductsDatabase()

	router := gin.Default()
	router.POST("/product/:id", func(context *gin.Context) {
		idParam := context.Param("id")
		postHandler(context, prodRepo, &idParam)
	})

	router.POST("/product", func(context *gin.Context) {
		postHandler(context, prodRepo, nil)
	})

	router.GET("/product/:id", func(context *gin.Context) {
		getHandler(context, database)
	})

	router.GET("/product", func(context *gin.Context) {
		getManyHandler(context, prodRepo)
	})

	router.DELETE("/product/:id", func(context *gin.Context) {
		deleteHandler(context, prodRepo)
	})
	return router
}
