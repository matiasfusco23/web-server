package main

import "github.com/gin-gonic/gin"

// import "server/cmd/handlers/product"
import "server/cmd/server/handlers"
func main() {
	server := gin.Default()
	// creo un grupo para productos
	productsGroup := server.Group("/products")
	// llamo a la creacion de un router.
	productRouter := handlers.NewProductRouter(productsGroup)

	// ejecuto el metodo que crea las rutas para que esten registradas
	productRouter.ProductRoutes()

	// corro mi server
	server.Run(":8081")
}
