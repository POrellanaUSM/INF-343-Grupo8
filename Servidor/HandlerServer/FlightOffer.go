package Handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetFlightOffers(c *gin.Context) {
	// Obtener token
	accessToken := GetToken()

	// Parámetros URL
	originLocationCode := c.DefaultQuery("originLocationCode", "")
	destinationLocationCode := c.DefaultQuery("destinationLocationCode", "")
	departureDate := c.DefaultQuery("departureDate", "")
	adults := c.DefaultQuery("adults", "")

	// URL
	apiURL := "https://test.api.amadeus.com/v2/shopping/flight-offers" +
		"?originLocationCode=" + originLocationCode +
		"&destinationLocationCode=" + destinationLocationCode +
		"&departureDate=" + departureDate +
		"&adults=" + adults +
		"&includedAirlineCodes=H2,LA,JA" +
		"&nonStop=true" +
		"&currencyCode=CLP" +
		"&travelClass=ECONOMY"

	//  solicitud HTTP 
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Hubo un error al realizar la búsqueda"})
		return
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	// Realizar la solicitud HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Hubo un error al realizar la búsqueda"})
		return
	}
	defer resp.Body.Close()

	
	c.Header("Content-Type", "application/json")

	if resp.StatusCode == http.StatusOK {
		// Código 200
		// Leer la respuesta HTTP sin ioutil.ReadAll
		c.DataFromReader(http.StatusOK, resp.ContentLength, "application/json", resp.Body, nil)
	} else if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		// Códigos 4XX
		c.Data(resp.StatusCode, "application/json", []byte(`{"message": "Hubo un error al realizar la búsqueda"}`))
	} else if resp.StatusCode >= 500 {
		// Códigos 5XX
		c.Data(resp.StatusCode, "application/json", []byte(`{"message": "Hubo un error al realizar la búsqueda"}`))
	}
}

