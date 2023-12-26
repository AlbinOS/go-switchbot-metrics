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
	"os/signal"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/nasa9084/go-switchbot"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var switchbotClient *switchbot.Client

type DeviceValue struct {
	DeviceId    string  `json:"deviceId"`
	HubDeviceId string  `json:"hubDeviceId"`
	DeviceName  string  `json:"deviceName"`
	DeviceType  string  `json:"deviceType"`
	Battery     int     `json:"battery"`
	Humidity    int     `json:"humidity"`
	Temperature float64 `json:"temperature"`
}

type Metrics struct {
	DevicesValue []DeviceValue `json:"devicesValue"`
}

func handler(c *fiber.Ctx) error {
	pdev, _, err := switchbotClient.Device().List(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("fetching device list failed!")
		return fiber.NewError(fiber.StatusServiceUnavailable, fmt.Errorf("error getting device list: %v", err).Error())
	}

	var result Metrics
	for _, d := range pdev {
		if d.Type != switchbot.HubMini {
			deviceStatus, err := switchbotClient.Device().Status(context.Background(), d.ID)

			if err != nil {
				log.Error().Str("deviceId", d.ID).Str("hubId", d.Hub).Str("deviceName", d.Name).Str("deviceType", string(d.Type)).Err(err).Msg("fetching device status failed!")
				return fiber.NewError(fiber.StatusServiceUnavailable, fmt.Errorf("fetching device status failed for %s (%s): %v", d.Name, d.Type, err).Error())
			} else {
				result.DevicesValue = append(result.DevicesValue,
					DeviceValue{
						DeviceId:    d.ID,
						HubDeviceId: d.Hub,
						DeviceName:  d.Name,
						DeviceType:  string(d.Type),
						Battery:     deviceStatus.Battery,
						Humidity:    deviceStatus.Humidity,
						Temperature: deviceStatus.Temperature,
					},
				)
			}
		}
	}

	return c.JSON(result)
}

func Init() {
	// Init logger
	logLevel := zerolog.DebugLevel
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).Level(logLevel)

	// Initialize the SwitchBot client
	switchbotClient = switchbot.New(viper.GetString("switchbot_openapi_token"), viper.GetString("switchbot_secret_key"))

	// Init Fiber
	app := fiber.New()

	// Listen for SIGINT (Ctrl+C) and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	var serverShutdown sync.WaitGroup

	go func() {
		<-c
		log.Info().Msg("Stopping go-switchbot-metrics serve...")
		serverShutdown.Add(1)
		defer serverShutdown.Done()
		_ = app.ShutdownWithTimeout(60 * time.Second)
	}()

	// Middlewares
	app.Use(recover.New())
	app.Use(fiberLogger.New())

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Pong!")
	})
	app.Get("/metrics", handler)
	app.Get("/monitoring", monitor.New(monitor.Config{Title: "go-switchbot-metrics Monitoring Page"}))

	// Start server
	if err := app.Listen(fmt.Sprintf("%s:%s", viper.GetString("bind_ip"), viper.GetString("bind_port"))); err != nil {
		log.Panic().Err(err).Msg("go-switchbot-metrics serve unable to start! ")
	}
	serverShutdown.Wait()
	log.Info().Msg("go-switchbot-metrics serve stopped.")
}
