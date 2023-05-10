package server

import (
	"gin-exercise/pkg/product"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type prodReq struct {
	Price       float32 `json:"price" binding:"required"`
	Description string  `json:"description" binding:"required"`
}

func postHandler(context *gin.Context, repository *product.Repository, id *string) {
	req := prodReq{}
	err := context.BindJSON(req)
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

func getHandler(context *gin.Context, repo *product.Repository) {
	id := context.Param("id")

	prod, found := repo.Find(id)
	if found {
		context.JSON(http.StatusOK, prod)
		return
	}
	context.AbortWithStatusJSON(http.StatusNotFound, map[string]string{"error": "not found"})
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
		context.JSON(http.StatusAccepted, map[string]string{"op": "deleted"})
		return
	}
	context.AbortWithStatusJSON(http.StatusNotFound, map[string]string{"error": "not found"})
}

func SetupRoutes() *gin.Engine {
	prodRepo := product.NewProductRepository()

	router := gin.Default()
	router.POST("/product/:id", func(context *gin.Context) {
		idParam := context.Param("id")
		postHandler(context, prodRepo, &idParam)
	})

	router.POST("/product", func(context *gin.Context) {
		postHandler(context, prodRepo, nil)
	})

	router.GET("/product/:id", func(context *gin.Context) {
		getHandler(context, prodRepo)
	})

	router.GET("/product", func(context *gin.Context) {
		getManyHandler(context, prodRepo)
	})

	router.DELETE("/product/:id", func(context *gin.Context) {
		deleteHandler(context, prodRepo)
	})
	return router
}
