package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	baseURL = "http://api.openweathermap.org/data/2.5/"
)

// Example /weather response
// {
// 	"coord":{
// 		"lon":-3.98,
// 		"lat":55.87
// 	},
// 	"weather":[{
// 		"id":310,
// 		"main":"Drizzle",
// 		"description":"light intensity drizzle rain",
// 		"icon":"09d"
// 	},{
// 		"id":500,
// 		"main":"Rain",
// 		"description":"light rain",
// 		"icon":"10d"
// 	}],
// 	"base":"stations",
// 	"main":{
// 		"temp":282.48,
// 		"pressure":1008,
// 		"humidity":87,
// 		"temp_min":282.15,
// 		"temp_max":283.15
// 	},
// 	"visibility":3200,
// 	"wind":{
// 		"speed":7.7,
// 		"deg":240
// 	},
// 	"clouds":{
// 		"all":90
// 	},
// 	"dt":1524473400,
// 	"sys":{
// 		"type":1,
// 		"id":5121,
// 		"message":0.0031,
// 		"country":"GB",
// 		"sunrise":1524459028,
// 		"sunset":1524512354
// 	},
// 	"id":2657613,
// 	"name":"Airdrie",
// 	"cod":200
// }

// APIClient defines an interface for an OpenWeatherMap API client
type APIClient interface {
	Weather(string) *weatherDict
}

type apiClient struct {
	httpClient *http.Client
	appID      string
}

type weatherDict struct {
	Description string
	Icon        string
	Timestamp   time.Time
}

// CoordData models OpenWeatherMap weather endpoint coord data
type CoordData struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

// WeatherData models the OpenWeatherMap weather endpoint weather data
type WeatherData struct {
	ID   int    `json:"id"`
	Main string `json:"main"`
	Desc string `json:"description"`
	Icon string `json:"icon"`
}

// Data models the OpenWeatherMap weather endpoint main data
type Data struct {
	Temp     float64 `json:"temp"`
	Pressure int     `json:"pressure"`
	Humidity int     `json:"humidity"`
	TempMin  float64 `json:"temp_min"`
	TempMax  float64 `json:"temp_max"`
}

// WindData models the OpenWeatherMap wind data
type WindData struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
}

// CloudData models the OpenWeatherMap cloud data
type CloudData struct {
	Cloudiness int `json:"all"`
}

// RainData models the OpenWeatherMap rain data
type RainData struct {
	Volume int `json:"3h"`
}

// SnowData models the OpenWeatherMap snow data
type SnowData struct {
	Volume int `json:"3h"`
}

// WeatherResponse models the OpenWeatherMap weather endpoint full weather response data
type WeatherResponse struct {
	Coords      CoordData     `json:"coord"`
	WeatherList []WeatherData `json:"weather"`
	Data        Data          `json:"main"`
	Visibility  int           `json:"visibility"`
	Wind        WindData      `json:"wind,omitempty"`
	Clouds      CloudData     `json:"clouds,omitempty"`
	Rain        RainData      `json:"rain,omitempty"`
	Snow        SnowData      `json:"snow,omitempty"`
	Timestamp   int64         `json:"dt"`
	ID          string        `json:"id"`
	LocName     string        `json:"name"`
}

// NewClient returns an API client instance
func NewClient(httpClient http.Client, appID string) APIClient {
	return &apiClient{&httpClient, appID}
}

func (client *apiClient) apiCall(url string) []byte {
	appIDURL := url + "&APPID=" + client.appID

	resp, err := client.httpClient.Get(appIDURL)

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return body
}

func (client *apiClient) Weather(location string) *weatherDict {
	url := baseURL + "weather?q=" + location
	response := client.apiCall(url)

	var jsonResp WeatherResponse
	json.Unmarshal(response, &jsonResp)

	return &weatherDict{
		jsonResp.WeatherList[0].Main,
		jsonResp.WeatherList[0].Icon,
		time.Unix(jsonResp.Timestamp, 0).UTC().In(time.Local),
	}
}
