package Handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FlightOffersPricingHandler(c *gin.Context) {
	// Obtener Token
	accessToken := GetToken()

	// Leer la solicitud
	var requestBody bytes.Buffer
	_, err := requestBody.ReadFrom(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Hubo un error al realizar la búsqueda"})
		return
	}

	// URL
	apiURL := "https://test.api.amadeus.com/v1/shopping/flight-offers/pricing"

	// Crear solicitud POST
	req, err := http.NewRequest("POST", apiURL, &requestBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Hubo un error al realizar la búsqueda"})
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	// Solicitud HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Hubo un error al realizar la búsqueda"})
		return
	}
	defer resp.Body.Close()

	// Leer la respuesta JSON 
	var response json.RawMessage
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Hubo un error al realizar la búsqueda"})
		return
	}

	
	if resp.StatusCode == http.StatusOK {
		// Código 200
		c.JSON(http.StatusOK, response)
	} else if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		// Códigos 4XX
		c.JSON(resp.StatusCode, gin.H{"message": "Hubo un error al realizar la búsqueda"})
	} else if resp.StatusCode >= 500 {
		// Códigos 5XX
		c.JSON(resp.StatusCode, gin.H{"message": "Hubo un error al realizar la búsqueda"})
	}
}
