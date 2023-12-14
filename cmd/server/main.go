package main

import "github.com/gin-gonic/gin"

// import "server/cmd/handlers/product"
import "server/cmd/server/handlers"
func main() {
	router := gin.Default()
	handlers.PopulateProducts("products.json")
	router.GET("/ping", handlers.Ping)

	router.GET("products", handlers.GetAllProducts)

	router.GET("products/:id", handlers.GetProductById)

	router.GET("products/withPriceGreaterThan", handlers.GetProductsWithPriceGreaterThan)

	router.POST("products/", handlers.CreateProduct)

	router.Run(":8081")
}
