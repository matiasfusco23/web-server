package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

/*
Vamos a crear un endpoint llamado /saludo. Con una pequeña estructura con nombre y apellido que al pegarle deberá responder en texto “Hola + nombre + apellido”

El endpoint deberá ser de método POST
Se deberá usar el package JSON para resolver el ejercicio
La respuesta deberá seguir esta estructura: “Hola Andrea Rivas”
La estructura deberá ser como esta:

	{
			“nombre”: “Andrea”,
			“apellido”: “Rivas”
	}
*/

type Persona struct {
	Nombre   string `json:"Nombre"`
	Apellido string `json:"Apellido"`
}

func main() {
	router := gin.Default()

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
	router.Run(":8081")
}
