package Handler

import (
	"fmt"

	"bufio"
	"encoding/json"
	"io/ioutil"
	"github.com/olekukonko/tablewriter"
	"net/http"
	"os"

)

// Define las estructuras para mapear el JSON

type SegmentR struct {
	Number    string `json:"number"`
	Departure struct {
		At string `json:"at"`
	} `json:"departure"`
	Arrival struct {
		At string `json:"at"`
	} `json:"arrival"`
	Aircraft struct {
		Code string `json:"code"`
	} `json:"aircraft"`
}

type Price struct {
	Total string `json:"total"`
}

type FlightOfferR struct {
	FlightOffers []struct {
		Itineraries []struct {
			SegmentsR []SegmentR `json:"segments"`
		} `json:"itineraries"`
		Price Price `json:"price"`
	} `json:"flightOffers"`
}

func ObtenerHandler() {
	server := os.Getenv("SERVER")
    port := os.Getenv("PORT")
	var ID string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Ingrese el ID de la reserva: ")
	scanner.Scan()
	ID = scanner.Text()
	
	// Realiza una solicitud HTTP GET a la API para obtener el JSON
	apiURL := "http://"+ server +":"+ port +"/api/booking/" + ID

	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return
	}
	defer resp.Body.Close()


	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}



	var reserva ReservaJSON

	err = json.Unmarshal(body, &reserva)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}
	
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Número", "Hora de Salida", "Hora de Llegada", "Avión", "Precio Total"})

		fmt.Println("Resultados: ")
		for _, segment := range reserva.Data.FlightOffers[0].Itineraries[0].Segments {
			table.Append([]string{
				segment.CarrierCode+segment.Number,
				segment.Departure.At[11:16],
				segment.Arrival.At[11:16],
				segment.Aircraft.Code,
				reserva.Data.FlightOffers[0].Price.Total,
			})
		}
		table.Render()
		fmt.Println("Pasajeros: ")

		table_P := tablewriter.NewWriter(os.Stdout)
		table_P.SetHeader([]string{"Nombre", "Apellido"})
		for _, traveler := range reserva.Data.Travelers  {
			table_P.Append([]string{
				traveler.Name.FirstName,
				traveler.Name.LastName,
			})
		}

		table_P.Render()
	

	

}
