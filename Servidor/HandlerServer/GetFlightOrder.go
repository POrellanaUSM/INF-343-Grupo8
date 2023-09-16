package Handler

import (
	"bytes"
	"io"
	"net/http"
	"net/url"


	"github.com/gin-gonic/gin"
)

func GetFlightOrder(c *gin.Context) {
	// Obtener token
	accessToken := GetToken()

	orderID := c.Param("orderID")
	// fmt.Println("orderID:", orderID) // Agrega esta línea para depurar
	encodedOrderID := url.QueryEscape(orderID)
	apiURL := "https://test.api.amadeus.com/v1/booking/flight-orders/" + encodedOrderID

	// Crear solicitud GET
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Hubo un error al realizar la búsqueda"})
		return
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	// Solicitud HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Hubo un error al realizar la búsqueda"})
		return
	}
	defer resp.Body.Close()

	// buffer para almacenar respuesta
	var responseBody bytes.Buffer

	// Copiar respuesta  al buffer
	_, err = io.Copy(&responseBody, resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Hubo un error al realizar la búsqueda"})
		return
	}


	c.Header("Content-Type", "application/json")

	if resp.StatusCode == http.StatusOK {
		// Código 200
		c.Data(http.StatusOK, "application/json", responseBody.Bytes())
	} else if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		// Códigos 4XX
		c.Data(resp.StatusCode, "application/json", []byte(`{"message": "Hubo un error al realizar la búsqueda"}`))
	} else if resp.StatusCode >= 500 {
		// Códigos 5XX
		c.Data(resp.StatusCode, "application/json", []byte(`{"message": "Hubo un error al realizar la búsqueda"}`))
	}
}
