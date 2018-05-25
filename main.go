package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	emoji "gopkg.in/kyokomi/emoji.v1"
	"gopkg.in/urfave/cli.v1"
)

const (
	defaultConfigFileLocation = "~/.dreich.conf.json"
)

// Config holds config values for the app
type Config struct {
	AppID           string `json:"app_id"`
	DefaultLocation string `json:"default_location,omitempty"`
}

var emojiMap = map[string]string{
	"01d": ":sunny:",
	"01n": ":crescent_moon:",
	"02d": ":partly_sunny:",
	"02n": ":partly_sunny:",
	"03d": ":cloud:",
	"03n": ":cloud:",
	"04d": ":cloud:",
	"04n": ":cloud:",
	"09d": ":umbrella:",
	"09n": ":umbrella:",
	"10d": ":umbrella:",
	"10n": ":umbrella:",
	"11d": ":zap:",
	"11n": ":zap:",
	"13d": ":snowflake:",
	"13n": ":snowflake:",
	"50d": ":foggy:",
	"50n": ":foggy:",
}

// tryExpandPath attempts to expand a given path and returns the expanded path
// if successful. Otherwise, if expansion failed, the original path is returned.
// From https://github.com/Rican7/define/blob/master/internal/config/config.go
func tryExpandPath(path string) string {
	if expanded, err := homedir.Expand(path); nil == err {
		path = expanded
	}

	return path
}

func init() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer, "%s\n", c.App.Version)
	}
}

func main() {

	err := os.MkdirAll(tryExpandPath("~/.dreich/cache/"), os.ModeDir)
	if err != nil {
		log.Println("Could not create cache directory:", err)
	}

	dreich := cli.NewApp()
	dreich.Name = "dreich"
	dreich.Version = "0.0.1"
	dreich.Usage = "A weather CLI tool"

	dreich.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "appid, a",
			Usage:  "OpenWeatherMap `APPID`",
			EnvVar: "OPEN_WEATHER_MAP_APPID",
		},
		cli.BoolFlag{
			Name:  "emoji, e",
			Usage: "show weather as emoji",
		},
		cli.StringFlag{
			Name:  "location, l",
			Usage: "define `location` to get weather for",
		},
		cli.BoolFlag{
			Name:  "tomorrow, t",
			Usage: "get the forecast for tomorrow",
		},
	}

	raw, err := ioutil.ReadFile(tryExpandPath(defaultConfigFileLocation))
	if err != nil {
		log.Fatal("Reading config:", err)
	}

	var cfg Config
	json.Unmarshal(raw, &cfg)

	dreich.Action = func(c *cli.Context) error {
		location := "London,uk"
		appID := ""

		if cfg.AppID != "" {
			appID = cfg.AppID
		}

		if c.String("appid") != "" {
			appID = c.String("appid")
		}

		weatherClient := NewClient(http.Client{}, appID)

		if cfg.DefaultLocation != "" {
			location = cfg.DefaultLocation
		}

		if c.String("location") != "" {
			location = c.String("location")
		}

		if c.Bool("tomorrow") {
			forecast := weatherClient.Tomorrow(location)

			for _, f := range *forecast {
				if c.Bool("emoji") {
					emoji.Println(f.Timestamp.Format("15:04") + "\t" + emojiMap[f.Icon])
				} else {
					fmt.Println(f.Timestamp.Format("15:04") + "\t" + f.Description)
				}
			}
		} else {
			weather := weatherClient.Weather(location)
			if c.Bool("emoji") {
				emoji.Println(emojiMap[weather.Icon])
			} else {
				log.Println(weather.Description)
			}
		}
		return nil
	}

	err = dreich.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
