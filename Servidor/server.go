package main

import (
	Handler "Tarea-1-Go/Servidor/HandlerServer"
	"fmt"
	
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)







func main() {
	// Cargar las variables de entorno desde el archivo .env
    if err := godotenv.Load(".env"); err != nil {
        log.Fatalf("Error cargando el archivo .env: %v", err)
    }

    port := os.Getenv("PORT")




	// Obtener el token al iniciar el servidor
	if err := Handler.FetchToken(); err != nil {
		fmt.Println("Error al obtener el token:", err)
		return
	}

	r := gin.Default()

	r.GET("/api/search", Handler.GetFlightOffers)
	r.POST("/api/pricing", Handler.FlightOffersPricingHandler)
	r.POST("/api/booking", Handler.CreateFlightOrderHandler)
	r.GET("/api/booking/:orderID", Handler.GetFlightOrder)

	
	fmt.Printf("API escuchando en el puerto %s\n", port)
	r.Run(":" + port)
}
