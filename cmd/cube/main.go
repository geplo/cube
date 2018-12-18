package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/geplo/cube"
	"github.com/geplo/cube/scenes"
	"github.com/geplo/cube/scenes/planeshift"
	"github.com/geplo/cube/scenes/rain"
	"github.com/geplo/cube/spi595"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/raspi"
)

var (
	// X wiring on Z axis.
	xMap = [][]int{{0, 1, 2, 3, 4, 5, 6, 7}, {7, 6, 5, 4, 3, 2, 1, 0}, {0, 1, 2, 3, 4, 5, 6, 7}, {7, 6, 5, 4, 3, 2, 1, 0}, {0, 1, 2, 3, 4, 5, 6, 7}, {7, 6, 5, 4, 3, 2, 1, 0}, {0, 1, 2, 3, 4, 5, 6, 7}, {7, 6, 5, 4, 3, 2, 1, 0}}
	// Y wiring on X axis.
	yMap = [][]int{{0, 1, 2, 3, 4, 5, 6, 7}, {0, 1, 2, 3, 4, 5, 6, 7}, {0, 1, 2, 3, 4, 5, 6, 7}, {0, 1, 2, 3, 4, 5, 6, 7}, {0, 1, 2, 3, 4, 5, 6, 7}, {0, 1, 2, 3, 4, 5, 6, 7}, {0, 1, 2, 3, 4, 5, 6, 7}, {0, 1, 2, 3, 4, 5, 6, 7}}
	// Z wiring on X axis.
	zMap = [][]int{{1, 0, 3, 2, 5, 4, 7, 6}, {1, 0, 3, 2, 5, 4, 7, 6}, {1, 0, 3, 2, 5, 4, 7, 6}, {1, 0, 3, 2, 5, 4, 7, 6}, {1, 0, 3, 2, 5, 4, 7, 6}, {1, 0, 3, 2, 5, 4, 7, 6}, {1, 0, 3, 2, 5, 4, 7, 6}, {1, 0, 3, 2, 5, 4, 7, 6}}
)

func main() {
	platform := raspi.NewAdaptor()
	adaptor := spi595.NewAdaptor(platform,
		spi595.WithXMap(xMap),
		spi595.WithYMap(yMap),
		spi595.WithZMap(zMap),
		spi595.WithBus(0),     // Bus 0 (/dev/spidev0).
		spi595.WithChip(0),    // CEO0 (/dev/spidev0.0).
		spi595.WithMode(0),    // SPI_MODE_0.
		spi595.WithBits(8),    // 8 bits per word.
		spi595.WithSpeed(8e6), // 8MHz.
	)
	driver := spi595.NewDriver(adaptor,
		cube.New(8),
		planeshift.New(),
	)

	sceneCtors := []func() scenes.Scene{
		planeshift.New,
		rain.New,
	}

	robot := gobot.NewRobot("spicube",
		[]gobot.Connection{adaptor},
		[]gobot.Device{driver},
	)
	robot.AddEvent("stop")

	robot.Work = func() {
		ticker1 := gobot.Every(5*time.Microsecond, func() {
			driver.Draw()
		})
		i := 0
		ticker2 := gobot.Every(10*time.Second, func() {
			i %= len(sceneCtors)
			driver.Scene(sceneCtors[i]())
			i++
		})
		if err := robot.On(robot.Event("stop"), func(interface{}) { ticker1.Stop(); ticker2.Stop() }); err != nil {
			log.Fatalf("On stop hook setup error: %s\n", err)
		}
	}

	// Start with autorun = false to properly handle stop.
	if err := robot.Start(false); err != nil {
		log.Fatalf("robot.Start error: %s\n", err)
	}

	// Handle ctrl-c.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	signal.Stop(c)
	close(c)

	// Cleanup.
	robot.Publish(robot.Event("stop"), time.Now())
	if err := robot.Stop(); err != nil {
		log.Fatalf("robot.Stop error: %s\n", err)
	}
	println("end")
	<-make(chan struct{})
}
