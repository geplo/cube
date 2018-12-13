package main

import (
	"time"

	"github.com/geplo/cube/cube888"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/spi"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	a := raspi.NewAdaptor()

	adc := cube888.NewDriver(a,
		spi.WithSpeed(8e6),
		spi.WithMode(0),
		spi.WithBits(8),
		spi.WithChip(0),
		spi.WithBus(0),
	)

	work := func() {
		gobot.Every(5*time.Microsecond, func() {
			adc.Step()
			adc.Draw()
		})
	}

	robot := gobot.NewRobot("spicube",
		[]gobot.Connection{a},
		[]gobot.Device{adc},
		work,
	)

	robot.Start()
}
