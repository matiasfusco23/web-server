package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

type ProductCreation struct {
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

var products []Product

func populate_products(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return errors.New("Could not open file")
	}
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return errors.New("Could not read file")
	}
	defer file.Close()
	unmarshall_err := json.Unmarshal([]byte(fileContent), &products)
	if unmarshall_err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return unmarshall_err
	}
	return nil
}

type Persona struct {
	Nombre   string `json:"Nombre"`
	Apellido string `json:"Apellido"`
}

func Ping(c *gin.Context) {
	c.String(200, "pong")
}
func GetAllProducts(c *gin.Context) {
	c.JSON(200, products)

}
func GetProductById(c *gin.Context) {
	for _, product := range products {
		param_id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(404, "Product not found")
			return

		}
		if product.Id == param_id {
			c.JSON(200, product)
			return
		}

	}
	c.String(404, "Product not found")

}

func GetProductsWithPriceGreaterThan(c *gin.Context) {
	var products_with_desired_price []Product
	for _, product := range products {
		min_price, err := strconv.ParseFloat(c.Query("min_price"), 64)

		if err != nil {
			c.String(404, "Please provide a valid price")
			return
		}

		if product.Price > min_price {
			products_with_desired_price = append(products_with_desired_price, product)
		}

	}
	c.JSON(200, products_with_desired_price)
}

func CreateProduct(c *gin.Context) {
	var request ProductCreation
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	lastIndex := len(products) - 1
	lastProductId := products[lastIndex].Id
	//TODO:
	/*
		Ningún dato puede estar vacío, exceptuando is_published (vacío indica un valor false).
		El campo code_value debe ser único para cada producto.
		La fecha de vencimiento debe tener el formato: XX/XX/XXXX, además debemos verificar que día, mes y año sean valores válidos.
	*/

	newProduct := Product{lastProductId + 1, request.Name, request.Quantity, request.CodeValue, request.IsPublished, request.Expiration, request.Price}
	products = append(products, newProduct)
	c.String(200, "Product added!")
}
func main() {
	router := gin.Default()
	populate_products("products.json")

	router.GET("/ping", Ping)

	router.GET("products", GetAllProducts)

	router.GET("products/:id", GetProductById)

	router.GET("products/withPriceGreaterThan", GetProductsWithPriceGreaterThan)

	router.POST("products/", CreateProduct)

	router.Run(":8081")
}
