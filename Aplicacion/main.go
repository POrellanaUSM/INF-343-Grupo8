package main

import (
	"Tarea-1-Go/Aplicacion/Handler"
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Cargar las variables de entorno desde el archivo .env
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error cargando el archivo .env: %v", err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Bienvenido a goTravel!")
	for {

		fmt.Println("1. Realizar búsqueda.")
		fmt.Println("2. Obtener reserva.")
		fmt.Println("3. Salir")

		var opcion int
		fmt.Print("Seleccione una opción: ")
		scanner.Scan()
		opcionStr := scanner.Text()
		_, err := fmt.Sscan(opcionStr, &opcion)

		if err != nil {
			fmt.Println("Error al leer la opción:", err)
			continue
		}

		switch opcion {
		case 1:
			Handler.BusquedaHandler()
		case 2:
			Handler.ObtenerHandler()
		case 3:
			fmt.Println("Gracias por usar goTravel!")
			return
		default:
			fmt.Println("Opción no válida. Por favor, seleccione una opción válida.")
		}

	}
}
