package Handler

import (
	"fmt"

	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

// //////////////////////////////////////////////////////////////////////////////////////////////////////////
// Estructura principal que representa el JSON completo.

// lo que recibo de la segunda parte
type RespuestaPreciosJSON struct {
	Data struct {
		FlightOffers []FlightOffer `json:"flightOffers"`
	}
	Diccionario DiccionarioData `json:"dictionaries"`
}

type RespuestaJSON struct {
	Meta        MetaInfo        `json:"meta"`
	Datos       []FlightOffer   `json:"data"`
	Diccionario DiccionarioData `json:"dictionaries"`
}

type RevisarPrecioJSON struct {
	Data struct {
		Type         string        `json:"type"`
		FlightOffers []FlightOffer `json:"flightOffers"`
	} `json:"data"`
}

type ReservaJSON struct {
	Data struct {
		Type               string             `json:"type"`
		FlightOffers       []FlightOffer      `json:"flightOffers"`
		Travelers          []Traveler         `json:"travelers"`
		Remarks            interface{}        `json:"remarks"`
		TicketingAgreement TicketingAgreement `json:"ticketingAgreement"`
		Contacts           interface{}        `json:"contacts"`
	} `json:"data"`
}

type RespuestaReservaJSON struct {
	Data struct {
		Type string `json:"type"`
		Id   string `json:"id"`
	} `json:"data"`
}

type Traveler struct {
	ID          string `json:"id"`
	DateOfBirth string `json:"dateOfBirth"`
	Name        struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	} `json:"name"`
	Gender  string `json:"gender"`
	Contact struct {
		EmailAddress string  `json:"emailAddress"`
		Phones       []Phone `json:"phones"`
	} `json:"contact"`
	Documents interface{} `json:"documents"`
}

type Phone struct {
	DeviceType         string `json:"deviceType"`
	CountryCallingCode string `json:"countryCallingCode"`
	Number             string `json:"number"`
}

type TicketingAgreement struct {
	Option string `json:"option"`
	Delay  string `json:"delay"`
}

// Información de metadatos.
type MetaInfo struct {
	Count int `json:"count"`
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
}

// Detalles de la oferta de vuelo.
type FlightOffer struct {
	Type                     string            `json:"type"`
	ID                       string            `json:"id"`
	Source                   string            `json:"source"`
	InstantTicketingRequired bool              `json:"instantTicketingRequired"`
	NonHomogeneous           bool              `json:"nonHomogeneous"`
	OneWay                   bool              `json:"oneWay"`
	LastTicketingDate        string            `json:"lastTicketingDate"`
	LastTicketingDateTime    string            `json:"lastTicketingDateTime"`
	NumberOfBookableSeats    int               `json:"numberOfBookableSeats"`
	Itineraries              []Itinerary       `json:"itineraries"`
	Price                    PriceInfo         `json:"price"`
	PricingOptions           PricingOptions    `json:"pricingOptions"`
	ValidatingAirlineCodes   []string          `json:"validatingAirlineCodes"`
	TravelerPricings         []TravelerPricing `json:"travelerPricings"`
}

// Detalles de la ruta de vuelo.
type Itinerary struct {
	Duration string    `json:"duration"`
	Segments []Segment `json:"segments"`
}

// Detalles de un segmento de vuelo.
type Segment struct {
	Departure       AirportInfo `json:"departure"`
	Arrival         AirportInfo `json:"arrival"`
	CarrierCode     string      `json:"carrierCode"`
	Number          string      `json:"number"`
	Aircraft        Aircraft    `json:"aircraft"`
	Duration        string      `json:"duration"`
	ID              string      `json:"id"`
	NumberOfStops   int         `json:"numberOfStops"`
	BlacklistedInEU bool        `json:"blacklistedInEU"`
}

// Información del aeropuerto.
type AirportInfo struct {
	IATACode string `json:"iataCode"`
	Terminal string `json:"terminal"`
	At       string `json:"at"`
}

// Información de la aeronave.
type Aircraft struct {
	Code string `json:"code"`
}

// Información de operación.

// Información de precios.
type PriceInfo struct {
	Currency   string `json:"currency"`
	Total      string `json:"total"`
	Base       string `json:"base"`
	Fees       []Fee  `json:"fees"`
	GrandTotal string `json:"grandTotal"`
}

// Detalles de las tarifas.
type Fee struct {
	Amount string `json:"amount"`
	Type   string `json:"type"`
}

// Opciones de precios.
type PricingOptions struct {
	FareType                []string `json:"fareType"`
	IncludedCheckedBagsOnly bool     `json:"includedCheckedBagsOnly"`
}

// Detalles de precios para viajeros.
type TravelerPricing struct {
	TravelerID   string     `json:"travelerId"`
	FareOption   string     `json:"fareOption"`
	TravelerType string     `json:"travelerType"`
	Price        PriceInfo  `json:"price"`
	FareDetails  []FareInfo `json:"fareDetailsBySegment"`
}

// Detalles de tarifas.
type FareInfo struct {
	SegmentID           string `json:"segmentId"`
	Cabin               string `json:"cabin"`
	FareBasis           string `json:"fareBasis"`
	BrandedFare         string `json:"brandedFare"`
	Class               string `json:"class"`
	IncludedCheckedBags struct {
		Weight     int    `json:"weight"`
		WeightUnit string `json:"weightUnit"`
	} `json:"includedCheckedBags"`
}

// Datos de diccionario.
type DiccionarioData struct {
	Locations  map[string]LocationInfo `json:"locations"`
	Aircrafts  map[string]string       `json:"aircraft"`
	Currencies map[string]string       `json:"currencies"`
	Carriers   map[string]string       `json:"carriers"`
}

// Información de ubicación.
type LocationInfo struct {
	CityCode    string `json:"cityCode"`
	CountryCode string `json:"countryCode"`
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////

// //////////////////////////////////////////////////////////////////////////////////////////////////////////
type Pasajero struct {
	FechaNacimiento string
	Nombre          string
	Apellido        string
	Sexo            string
	Correo          string
	Telefono        string
}






func BusquedaHandler() {
	server := os.Getenv("SERVER")
    port := os.Getenv("PORT")
	scanner := bufio.NewScanner(os.Stdin)
	var respuestaJSON RespuestaJSON
	var origen string
	var destino string
	var fecha_salida string
	var num_adultos string
	var input string
	var seleccion int
	seleccion = 0
	for seleccion == 0 {
		fmt.Print("Aeropuerto de origen: ")
		scanner.Scan()
		origen = scanner.Text()
		fmt.Print("Aeropuerto de destino: ")
		scanner.Scan()
		destino = scanner.Text()
		fmt.Print("Fecha de salida: ")
		scanner.Scan()
		fecha_salida = scanner.Text()
		fmt.Print("Cantidad de adultos: ")
		scanner.Scan()
		num_adultos = scanner.Text()


		apiURL := "http://"+ server +":"+ port +"/api/search?originLocationCode=" + origen + "&destinationLocationCode=" + destino + "&departureDate=" + fecha_salida + "&adults=" + num_adultos

		// Realizar una solicitud GET
		respuesta, err := http.Get(apiURL)
		if err != nil {
			fmt.Println("Error al realizar la solicitud:", err)
			return
		}

		defer respuesta.Body.Close()
		

		cuerpo, err := ioutil.ReadAll(respuesta.Body)
		if err != nil {
			fmt.Println("Error al leer el cuerpo de la respuesta:", err)
			return
		}

		err = json.Unmarshal(cuerpo, &respuestaJSON)
		if err != nil {
			fmt.Println("Error al decodificar el JSON:", err)
			return
		}

		fmt.Println("Se obtuvieron los siguientes resultados:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Vuelo", "Número", "Hora de Salida", "Hora de Llegada", "Avión", "Precio Total"})

		for _, oferta := range respuestaJSON.Datos {
			table.Append([]string{
				oferta.ID,
				oferta.Itineraries[0].Segments[0].CarrierCode + oferta.Itineraries[0].Segments[0].Number,
				oferta.Itineraries[0].Segments[0].Departure.At[11:16],
				oferta.Itineraries[0].Segments[0].Arrival.At[11:16],
				oferta.Itineraries[0].Segments[0].Aircraft.Code,
				oferta.Price.Total,
			})
		}

		table.Render()

		fmt.Print("Seleccione un vuelo (ingrese 0 para realizar nueva búsqueda):")
		scanner.Scan()
		input = scanner.Text()
		seleccion, err = strconv.Atoi(input)
		if err != nil {
			fmt.Println("Error al obtener vuelo")
		}
	}

	selectedOffer := respuestaJSON.Datos[seleccion-1]

	
	revisarPrecioJSON := RevisarPrecioJSON{
		Data: struct {
			Type         string        `json:"type"`
			FlightOffers []FlightOffer `json:"flightOffers"`
		}{
			Type: "flight-offers-pricing",
			FlightOffers: []FlightOffer{
				selectedOffer,
			},
		},
	}


	// Codifica el nuevo JSON en una cadena
	revisarPrecioJSONStr, erro := json.Marshal(revisarPrecioJSON) //JSON que enviamos para obtener los precios
	if erro != nil {
		fmt.Println("Error al codificar el nuevo JSON:", erro)
		return
	}


	//ENVIAR JSON PREGUNTANDO PRECIO DEFINITIVO

	apiURLprecio := "http://"+ server +":"+ port +"/api/pricing"

	data := []byte(revisarPrecioJSONStr)
	//Realizar llamado a api
	RespuestaPrecio, erra := http.Post(apiURLprecio, "application/json", bytes.NewBuffer(data))
	if erra != nil {
		panic(erra)
	}
	defer RespuestaPrecio.Body.Close()
	// Leer cuerpo respuesta
	cuerpo, erre := ioutil.ReadAll(RespuestaPrecio.Body)
	if erre != nil {
		fmt.Println("Error al leer la respuesta:", erre)
		return
	}

	var respuestaPreciosJSON RespuestaPreciosJSON
	erre = json.Unmarshal(cuerpo, &respuestaPreciosJSON)
	if erre != nil {
		fmt.Println("Error al leer JSON:", erre)
		return
	}

	//BUSCAR PRECIO DEL VUELO SELECCIONADO Y OBTENER JSON DE RESPUESTA CON PRECIO

	fmt.Println("El precio total final es de: ", respuestaPreciosJSON.Data.FlightOffers[0].Price.Total)

	//generar JSON para obtener Reserva
	reservaJSON := ReservaJSON{
		Data: struct {
			Type               string             `json:"type"`
			FlightOffers       []FlightOffer      `json:"flightOffers"`
			Travelers          []Traveler         `json:"travelers"`
			Remarks            interface{}        `json:"remarks"`
			TicketingAgreement TicketingAgreement `json:"ticketingAgreement"`
			Contacts           interface{}        `json:"contacts"`
		}{
			Type: "flight-order",
			FlightOffers: []FlightOffer{
				respuestaPreciosJSON.Data.FlightOffers[0],
			},
		},
	}
	reservaJSON.Data.TicketingAgreement.Option = "DELAY_TO_CANCEL"
	reservaJSON.Data.TicketingAgreement.Delay = "6D"

	// var pasajeros []Pasajero
	var numPasajeros, err = strconv.Atoi(num_adultos)



	for i := 1; i <= numPasajeros; i++ {
		var traveler Traveler
		fmt.Printf("\nPasajero %d:\n", i)
		traveler.ID = strconv.Itoa(i)
		fmt.Print("Ingrese Fecha de nacimiento (YYYY-MM-DD): ")
		scanner.Scan()
		traveler.DateOfBirth = scanner.Text()
		fmt.Print("Ingrese Nombre: ")
		scanner.Scan()
		traveler.Name.FirstName = scanner.Text()
		fmt.Print("Ingrese Apellido: ")
		scanner.Scan()
		traveler.Name.LastName = scanner.Text()
		fmt.Print("Ingrese Sexo (MALE o FEMALE): ")
		scanner.Scan()
		traveler.Gender = scanner.Text()
		fmt.Print("Ingrese Correo: ")
		scanner.Scan()
		traveler.Contact.EmailAddress = scanner.Text()

		var numPhones int = 1

		for j := 1; j <= numPhones; j++ {
			var phone Phone
			phone.DeviceType = "MOBILE"

			phone.CountryCallingCode = "34"
			fmt.Print("Ingrese Teléfono: ")
			scanner.Scan()
			phone.Number = scanner.Text()

			traveler.Contact.Phones = append(traveler.Contact.Phones, phone)
		}
		//Agregamos viajero
		reservaJSON.Data.Travelers = append(reservaJSON.Data.Travelers, traveler)

	}

	// Generar JSON
	jsonData, err := json.MarshalIndent(reservaJSON, "", "    ")
	if err != nil {
		fmt.Println("Error al generar el JSON:", err)
		return
	}


	apiURLreserva := "http://"+ server +":"+ port +"/api/booking"

	payload := bytes.NewBuffer(jsonData)

	// Realiza la solicitud POST.
	respuestaReserva, err := http.Post(apiURLreserva, "application/json", payload)
	if err != nil {
		fmt.Println("Error al realizar la solicitud POST:", err)
		return
	}
	defer respuestaReserva.Body.Close()

	// Leer cuerpo respuesta
	cuerpoReserva, erre := ioutil.ReadAll(respuestaReserva.Body)
	if erre != nil {
		fmt.Println("Error al leer la respuesta:", erre)
		return
	}

	var respuestaReservaJSON RespuestaReservaJSON
	erre = json.Unmarshal(cuerpoReserva, &respuestaReservaJSON)
	if erre != nil {
		fmt.Println("Error al leer JSON:", erre)
		return
	}
	fmt.Println("Reserva creada con éxito: ", respuestaReservaJSON.Data.Id)
}

//////////////////////////////////////////////////////////////////////////////////////////////
