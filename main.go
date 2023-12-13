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

func main() {
	router := gin.Default()
	populate_products("products.json")

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.POST("/saludo", func(c *gin.Context) {
		reqBody := new(Persona)
		err := c.BindJSON(&reqBody)
		if err != nil {
			fmt.Println("error: ", err.Error())
			c.String(400, fmt.Sprintf("Error in body: %s err", err.Error()))
			return
		}
		resultado := fmt.Sprintf("Hola %s %s", reqBody.Nombre, reqBody.Apellido)
		c.String(200, resultado)
		fmt.Println("Ejecucion exitosa")

	})

	//puedo definir endpoints agrupados
	grupo := router.Group("/perfil")
	grupo.GET("/foto")
	grupo.GET("/info")

	var empleados = map[string]string{
		"1": "pepito",
		"2": "jaimito",
	}

	router.GET("empleados/:id", func(c *gin.Context) {
		empleado, ok := empleados[c.Param("id")]
		nombre := c.Query("nombre")
		fmt.Println("Id param received: ", c.Param("id"))
		if ok && nombre == empleado {
			c.String(200, "Nombre: %s, ID: %s", empleado, c.Param("id"))
		} else {
			c.String(404, "No se encontro el empleado, por favor dar su nombre por query param")
		}
	})

	router.GET("products", func(c *gin.Context) {
		c.JSON(200, products)

	})

	router.GET("products/:id", func(c *gin.Context) {
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

	})

	router.GET("products/withPriceGreaterThan", func(c *gin.Context) {
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
	})

	router.POST("products/", func(c *gin.Context) {
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
	})

	router.Run(":8081")
}
