package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	Tomorrow(string) *[]weatherDict
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

// Example forecast response
// {
// 	"city":{
// 		"id":1851632,
// 		"name":"Shuzenji",
// 		"coord":{
// 			"lon":138.933334,
// 			"lat":34.966671
// 		},
// 		"country":"JP",
// 		"cod":"200",
// 		"message":0.0045,
// 		"cnt":38,
// 		"list":[{
// 			"dt":1406106000,
// 			"main":{
// 				"temp":298.77,
// 				"temp_min":298.77,
// 				"temp_max":298.774,
// 				"pressure":1005.93,
// 				"sea_level":1018.18,
// 				"grnd_level":1005.93,
// 				"humidity":87,
// 				"temp_kf":0.26
// 			},
// 			"weather":[{
// 				"id":804,
// 				"main":"Clouds",
// 				"description":"overcast clouds",
// 				"icon":"04d"
// 			}],
// 			"clouds":{
// 				"all":88
// 			},
// 			"wind":{
// 				"speed":5.71,
// 				"deg":229.501
// 			},
// 			"sys":{
// 				"pod":"d"
// 			},
// 			"dt_txt":"2014-07-23 09:00:00"
// 		}]
// 	}
// }

// ForecastElement models the OpenWeatherMap forecast endpoint forecast data
type ForecastElement struct {
	Timestamp     int64         `json:"dt"`
	Data          Data          `json:"main"`
	WeatherList   []WeatherData `json:"weather"`
	Wind          WindData      `json:"wind,omitempty"`
	Clouds        CloudData     `json:"clouds,omitempty"`
	Rain          RainData      `json:"rain,omitempty"`
	Snow          SnowData      `json:"snow,omitempty"`
	TimestampText string        `json:"dt_txt"`
}

// CityData models the OpenWeatherMap forecast endpoint city data
type CityData struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	Coords  CoordData `json:"coord"`
	Country string    `json:"country"`
}

// ForecastResponse models the OpenWeatherMap forecast endpoint response data
type ForecastResponse struct {
	City      CityData          `json:"city"`
	Forecasts []ForecastElement `json:"list"`
}

// NewClient returns an API client instance
func NewClient(httpClient http.Client, appID string) APIClient {
	return &apiClient{&httpClient, appID}
}

func (client *apiClient) apiCall(url string) []byte {
	urlHash := sha1.New()
	urlHash.Write([]byte(url))

	cachePath := tryExpandPath(fmt.Sprintf("~/.dreich/cache/%x.json", urlHash.Sum(nil)))
	if stat, err := os.Stat(cachePath); !os.IsNotExist(err) {
		now := time.Now()
		diff := now.Sub(stat.ModTime())
		if diff.Seconds() < 300 {
			raw, err := ioutil.ReadFile(cachePath)
			if err != nil {
				log.Println("Could not read cache file:", err)
			} else {
				log.Println("Reading from cache!")
				return raw
			}
		}
	}

	appIDURL := url + "&APPID=" + client.appID

	resp, err := client.httpClient.Get(appIDURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if err := ioutil.WriteFile(cachePath, body, 0644); err != nil {
		log.Println("Could not write cache file:", err)
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

func (client *apiClient) Forecast(location string) *[]weatherDict {
	url := baseURL + "forecast?q=" + location
	response := client.apiCall(url)

	var jsonResp ForecastResponse
	json.Unmarshal(response, &jsonResp)

	weatherDicts := []weatherDict{}
	for _, w := range jsonResp.Forecasts {
		weatherDicts = append(weatherDicts, weatherDict{
			w.WeatherList[0].Main,
			w.WeatherList[0].Icon,
			time.Unix(w.Timestamp, 0).UTC().In(time.Local),
		})
	}

	return &weatherDicts
}

func (client *apiClient) Tomorrow(location string) *[]weatherDict {
	forecastData := client.Forecast(location)

	tom := time.Now().AddDate(0, 0, 1)
	dayAfter := tom.AddDate(0, 0, 1)

	start := time.Date(tom.Year(), tom.Month(), tom.Day(), 0, 0, 0, 0, tom.Location())
	end := time.Date(dayAfter.Year(), dayAfter.Month(), dayAfter.Day(), 0, 0, 0, 0, dayAfter.Location())

	tomorrowData := []weatherDict{}
	for _, w := range *forecastData {
		if w.Timestamp.After(start) && w.Timestamp.Before(end) {
			tomorrowData = append(tomorrowData, w)
		}
	}

	return &tomorrowData
}
