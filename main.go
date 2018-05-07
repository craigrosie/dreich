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
	AppID string `json:"app_id"`
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
	dreich := cli.NewApp()
	dreich.Name = "dreich"
	dreich.Version = "0.0.1"
	dreich.Description = "A weather CLI tool"

	dreich.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "emoji, e",
			Usage: "show weather as emoji",
		},
	}

	raw, err := ioutil.ReadFile(tryExpandPath(defaultConfigFileLocation))
	if err != nil {
		log.Fatal("Reading config:", err)
	}

	var cfg Config
	json.Unmarshal(raw, &cfg)

	weatherClient := NewClient(http.Client{}, cfg.AppID)

	dreich.Action = func(c *cli.Context) error {
		weather := weatherClient.Weather("London,uk")
		if c.Bool("emoji") {
			emoji.Println(emojiMap[weather.Icon])
		} else {
			log.Println(weather.Description)
		}
		return nil
	}

	err = dreich.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
