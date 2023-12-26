/*
Copyright © 2023 Albin Gilles <gilles.albin@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package serve

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/nasa9084/go-switchbot"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	openToken   = ""
	secretKey   = ""
	influxToken = ""
)

var c *switchbot.Client
var client influxdb2.Client
var writeAPI api.WriteAPIBlocking

func fetch() {
	// get physical devices and show
	pdev, _, _ := c.Device().List(context.Background())

	for _, d := range pdev {
		if d.Type != switchbot.HubMini {
			deviceStatus, err := c.Device().Status(context.Background(), d.ID)

			if err != nil {
				fmt.Printf("%s\t%s\t%s\n", d.Type, d.Name, err)
			} else {
				timestamp := time.Now()
				// Create a point using the full parameters constructor
				pt := influxdb2.NewPoint("temperature",
					map[string]string{"device_name": d.Name, "device_type": string(d.Type)},
					map[string]interface{}{"battery": deviceStatus.Battery, "humidity": deviceStatus.Humidity, "temperature": deviceStatus.Temperature},
					timestamp)
				// Write the point immediately
				writeAPI.WritePoint(context.Background(), pt)
			}
		}
	}

	// Ensure that background processes finish
	client.Close()
}

func Init() {
	// Init logger
	logLevel := zerolog.DebugLevel
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).Level(logLevel)

	// Initialize the SwitchBot client
	c = switchbot.New(openToken, secretKey)

	// Initialize the InfluxDB client
	client = influxdb2.NewClient("http://localhost:8086", influxToken)
	writeAPI = client.WriteAPIBlocking("GoPex", "Home")

	// Init scheduler
	scheduler := gocron.NewScheduler(time.UTC)

	// Schedule fetchCurrentBIS every 10 seconds
	fetchJob, err := scheduler.Every(10).Seconds().Do(fetch)
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("Job error: %v!", fetchJob))
	}

	// Start scheduler
	log.Info().Msg("Starting go-switchbot-influx fetcher...")
	scheduler.StartBlocking()
	log.Info().Msg("go-switchbot-influx stopped!")
}
