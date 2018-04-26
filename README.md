# Dreich

A command-line weather app, written in Go.

*Dreich*: _(adjective_) (especially of weather) dreary; bleak _(Scottish)_

## Installation

```bash
go get github.com/craigrosie/dreich
```

## Configuration

`dreich` looks for a configuration file at `~/.dreich.conf.json`, and expects the following format:

```json
{
    "app_id": "<OpenWeatherMap API Key>"
}
```

You can follow the process at [https://openweathermap.org/appid](https://openweathermap.org/price) to obtain an API key.

`dreich` has been built using the Free tier of OpenWeatherMap. To view their other (paid) tiers, visit [https://openweathermap.org/price](https://openweathermap.org/price).
