package Handler

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo/options"

	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"context"
	"encoding/json"
	"fmt"
)

type Co2Emission struct {
	Weight     int    `json:"weight"`
	WeightUnit string `json:"weightUnit"`
	Cabin      string `json:"cabin"`
}

type Aircraft struct {
	Code string `json:"code"`
}

type Segment struct {
	Departure     Airport       `json:"departure"`
	Arrival       Airport       `json:"arrival"`
	CarrierCode   string        `json:"carrierCode"`
	Number        string        `json:"number"`
	Aircraft      Aircraft      `json:"aircraft"`
	Duration      string        `json:"duration"`
	ID            string        `json:"id"`
	NumberOfStops int           `json:"numberOfStops"`
	Co2Emissions  []Co2Emission `json:"co2Emissions"`
}

type Itinerary struct {
	Segments []Segment `json:"segments"`
}

type Location struct {
	CityCode    string `json:"cityCode"`
	CountryCode string `json:"countryCode"`
}

type Locations struct {
	ARI Location `json:"ARI"`
	SCL Location `json:"SCL"`
}

type Airport struct {
	IataCode string `json:"iataCode"`
	At       string `json:"at"`
}

type PricingOption struct {
	FareType                []string `json:"fareType"`
	IncludedCheckedBagsOnly bool     `json:"includedCheckedBagsOnly"`
}

type Price struct {
	Currency        string `json:"currency"`
	Total           string `json:"total"`
	Base            string `json:"base"`
	Fees            []Fee  `json:"fees"`
	GrandTotal      string `json:"grandTotal"`
	BillingCurrency string `json:"billingCurrency"`
}

type Fee struct {
	Amount string `json:"amount"`
	Type   string `json:"type"`
}

type FareDetailsBySegment struct {
	SegmentID           string              `json:"segmentId"`
	Cabin               string              `json:"cabin"`
	FareBasis           string              `json:"fareBasis"`
	BrandedFare         string              `json:"brandedFare"`
	Class               string              `json:"class"`
	IncludedCheckedBags IncludedCheckedBags `json:"includedCheckedBags"`
}

type IncludedCheckedBags struct {
	Quantity int `json:"quantity"`
}

type TravelerPricing struct {
	TravelerID           string                 `json:"travelerId"`
	FareOption           string                 `json:"fareOption"`
	TravelerType         string                 `json:"travelerType"`
	Price                Price                  `json:"price"`
	FareDetailsBySegment []FareDetailsBySegment `json:"fareDetailsBySegment"`
}

type Traveler struct {
	ID          string  `json:"id"`
	DateOfBirth string  `json:"dateOfBirth"`
	Gender      string  `json:"gender"`
	Name        Name    `json:"name"`
	Contact     Contact `json:"contact"`
}

type Name struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Contact struct {
	Purpose      string  `json:"purpose"`
	Phones       []Phone `json:"phones"`
	EmailAddress string  `json:"emailAddress"`
}

type Phone struct {
	DeviceType         string `json:"deviceType"`
	CountryCallingCode string `json:"countryCallingCode"`
	Number             string `json:"number"`
}

type AssociatedRecord struct {
	Reference        string `json:"reference"`
	CreationDate     string `json:"creationDate"`
	OriginSystemCode string `json:"originSystemCode"`
	FlightOfferID    string `json:"flightOfferId"`
}

type FlightOffer struct {
	Type                   string            `json:"type"`
	ID                     string            `json:"id"`
	Source                 string            `json:"source"`
	NonHomogeneous         bool              `json:"nonHomogeneous"`
	LastTicketingDate      string            `json:"lastTicketingDate"`
	Itineraries            []Itinerary       `json:"itineraries"`
	Price                  Price             `json:"price"`
	PricingOptions         PricingOption     `json:"pricingOptions"`
	ValidatingAirlineCodes []string          `json:"validatingAirlineCodes"`
	TravelerPricings       []TravelerPricing `json:"travelerPricings"`
}

type Data struct {
	Type               string             `json:"type"`
	ID                 string             `json:"id"`
	QueuingOfficeID    string             `json:"queuingOfficeId"`
	AssociatedRecords  []AssociatedRecord `json:"associatedRecords"`
	FlightOffers       []FlightOffer      `json:"flightOffers"`
	Travelers          []Traveler         `json:"travelers"`
	TicketingAgreement TicketingAgreement `json:"ticketingAgreement"`
	AutomatedProcess   []AutomatedProcess `json:"automatedProcess"`
}

type TicketingAgreement struct {
	Option string `json:"option"`
	Delay  string `json:"delay"`
}

type AutomatedProcess struct {
	Code     string `json:"code"`
	Queue    Queue  `json:"queue"`
	OfficeID string `json:"officeId"`
}

type Queue struct {
	Number   string `json:"number"`
	Category string `json:"category"`
}

type Response struct {
	Data         Data         `json:"data"`
	Dictionaries Dictionaries `json:"dictionaries"`
}

type Dictionaries struct {
	Locations Locations `json:"locations"`
}

func insertDocument(client *mongo.Client, response Response) error {
	// Nombre de tu base de datos y colección
	databaseName := "gotravel"
	collectionName := "reservations"

	// Obtén la colección
	collection := client.Database(databaseName).Collection(collectionName)

	// Inserta el documento JSON en la colección
	_, err := collection.InsertOne(context.TODO(), response)
	return err
}
func createMongoClient() (*mongo.Client, error) {
	// Define la cadena de conexión de MongoDB
    connectionString := os.Getenv("CONNECTION_STRING")
  

	// opciones del cliente
	clientOptions := options.Client().ApplyURI(connectionString)

	// cliente de MongoDB
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Connect(context.Background())
	if err != nil {
		return nil, err
	}

	// Verifica la conexión a MongoDB
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
func CreateFlightOrderHandler(c *gin.Context) {
	// Obtener el token de acceso con TokenHandler
	accessToken := GetToken()

	// Leer la solicitud
	var requestBody bytes.Buffer
	_, err := requestBody.ReadFrom(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Hubo un error al realizar la búsqueda"})
		return
	}

	// URL
	apiURL := "https://test.api.amadeus.com/v1/booking/flight-orders"

	// solicitud POST
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

	// Leer la respuesta HTTP 
	var responseBody bytes.Buffer
	_, err = responseBody.ReadFrom(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Hubo un error al realizar la búsqueda"})
		return
	}

	// respuesta
	var responseReserva Response
	if err := json.Unmarshal(responseBody.Bytes(), &responseReserva); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Hubo un error al deserializar la respuesta"})
		fmt.Println("error: ", err)
		return
	}

	clientBD, err := createMongoClient()
	if err != nil {
		log.Fatal("Error al crear el cliente de MongoDB:", err)
	}
	defer clientBD.Disconnect(context.Background())
	// Insertar respuesta en MongoDB
	err = insertDocument(clientBD, responseReserva)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Hubo un error al guardar en la base de datos"})
		return
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		// Código 200
		c.Data(http.StatusOK, "application/json", responseBody.Bytes())
	} else if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		// Códigos 4XX
		c.JSON(resp.StatusCode, gin.H{"message": "Hubo un error al realizar la búsqueda"})
	} else if resp.StatusCode >= 500 {
		// Códigos 5XX
		c.JSON(resp.StatusCode, gin.H{"message": "Hubo un error al realizar la búsqueda"})
	}
}
