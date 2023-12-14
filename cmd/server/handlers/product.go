package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"server/internal/product"
	"server/internal/domain"
	"github.com/gin-gonic/gin"
)

// Una estructura que representa un router
type ProductRouter struct {
	// grupo al que pertenece este conjunto de urls
	productGroup *gin.RouterGroup
	// el service de productos
	service product.ProductService
}

// constructor del router de productos
func NewProductRouter(g *gin.RouterGroup) ProductRouter {
	// un slice que se rellena con una llamada al metodo util de carga de json
	slice := PopulateProducts("../../products.json")
	// creo un repo y le paso el slice
	repo := product.NewProductRepository(slice)
	// creo el service y le paso el repo que cree
	serv := product.NewProductService(repo)
	// creo un router y le paso el grupo, el service y el repo
	return ProductRouter{g, serv}
}


func PopulateProducts(filename string) []domain.Product {
	var products [] domain.Product
	file, err := os.Open(filename)
	/*if err != nil {
		return nil,errors.New("Could not open file")
	}*/
	fileContent, err := ioutil.ReadAll(file)
	/*if err != nil {
		fmt.Println("Error reading file:", err)
		return nil,errors.New("Could not read file")
	}*/
	defer file.Close()
	unmarshall_err := json.Unmarshal([]byte(fileContent), &products)
	/*if unmarshall_err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil,unmarshall_err
	}*/

	if err != nil && unmarshall_err != nil {
		fmt.Println(err)
		fmt.Println(unmarshall_err)
		return products
	}
	return products
}



// conjunto de rutas de URL
func (r *ProductRouter) ProductRoutes() {
	r.productGroup.GET("/ping", r.Ping())

	r.productGroup.GET("/", r.GetAllProducts())

	r.productGroup.GET("/:id", r.GetProductById())

	r.productGroup.GET("/withPriceGreaterThan", r.GetProductsWithPriceGreaterThan())

	r.productGroup.POST("/", r.CreateProduct())
}

func (r *ProductRouter) Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(200, "pong")
	}
}

func (r *ProductRouter) GetAllProducts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data := r.service.GetAllProducts()
		ctx.JSON(http.StatusOK, data)
	}
}

func (r *ProductRouter) GetProductById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			fmt.Println(err)
			ctx.String(500, "Invalid ID")
			return
		}

		data, service_err := r.service.GetById(id)

		if service_err != nil {
			ctx.String(400, service_err.Error())
		}

		ctx.JSON(http.StatusOK, data)
	}
}

func (r *ProductRouter) GetProductsWithPriceGreaterThan() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		min_price, err := strconv.ParseFloat(ctx.Query("min_price"), 64)

		if err != nil {
			fmt.Println(err)
			ctx.String(500, "Invalid min price")
			return
		}

		data, service_err := r.service.GetProductsWithPriceGreaterThan(min_price)

		if service_err != nil {
			ctx.String(400, service_err.Error())
		}

		ctx.JSON(http.StatusOK, data)
	}
}

func (r *ProductRouter) CreateProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
	var request product.ProductCreation
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

		service_err := r.service.CreateProduct(request)

		if service_err != nil {
			ctx.String(400, service_err.Error())
		}

		ctx.String(http.StatusOK, "Product added!")
	}
}
