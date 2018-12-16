package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/pkg/errors"

	"github.com/geplo/cube"
	"github.com/geplo/cube/scenes"
	"github.com/geplo/cube/scenes/planeshift"
	"github.com/geplo/cube/scenes/rain"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/spi"
)

var (
	spiHandler spi.Connection
	c          cube.Cube
	scene      scenes.Scene
	next       time.Time
)

// setup is called before the main loop.
func setup() error {
	// Instantiate the cube.
	c = cube.New(8)

	// Initialize SPI.
	hdlr, err := spi.GetSpiConnection(
		0, 0, // Bus 0, CEO0, i.e. /dev/spidev0.0.
		0,   // SPI_MODE_0.
		8,   // 8 bits per words.
		8e6, // 8MHz.
	)
	if err != nil {
		return errors.Wrap(err, "GetSpiConnection")
	}
	spiHandler = hdlr

	// Seed the random generator.
	rand.Seed(time.Now().UnixNano())

	// Set the scene to use.
	scene = planeshift.New()
	scene = rain.New()

	return nil
}

func loop() error {
	start := time.Now()

	if next.Before(start) {
		// Step the scene.
		next = start.Add(scene.Step(c))
	}

	// Render the cube.
	if err := renderCube(spiHandler, c); err != nil {
		return errors.Wrap(err, "renderCube")
	}

	// Repeat.
	return nil
}

var (
	// X wiring on Z axis.
	xMap = [][]int{{0, 1, 2, 3, 4, 5, 6, 7}, {7, 6, 5, 4, 3, 2, 1, 0}, {0, 1, 2, 3, 4, 5, 6, 7}, {7, 6, 5, 4, 3, 2, 1, 0}, {0, 1, 2, 3, 4, 5, 6, 7}, {7, 6, 5, 4, 3, 2, 1, 0}, {0, 1, 2, 3, 4, 5, 6, 7}, {7, 6, 5, 4, 3, 2, 1, 0}}
	// Y wiring on X axis.
	yMap = [][]int{{0, 1, 2, 3, 4, 5, 6, 7}, {0, 1, 2, 3, 4, 5, 6, 7}, {0, 1, 2, 3, 4, 5, 6, 7}, {0, 1, 2, 3, 4, 5, 6, 7}, {0, 1, 2, 3, 4, 5, 6, 7}, {0, 1, 2, 3, 4, 5, 6, 7}, {0, 1, 2, 3, 4, 5, 6, 7}, {0, 1, 2, 3, 4, 5, 6, 7}}
	// Z wiring on X axis.
	zMap = [][]int{{1, 0, 3, 2, 5, 4, 7, 6}, {1, 0, 3, 2, 5, 4, 7, 6}, {1, 0, 3, 2, 5, 4, 7, 6}, {1, 0, 3, 2, 5, 4, 7, 6}, {1, 0, 3, 2, 5, 4, 7, 6}, {1, 0, 3, 2, 5, 4, 7, 6}, {1, 0, 3, 2, 5, 4, 7, 6}, {1, 0, 3, 2, 5, 4, 7, 6}}
)

func mapCube(src cube.Cube) cube.Cube {
	// Make sure the dst is empty.
	dst := cube.NewCustom(src.XLen, src.YLen, src.ZLen)

	// For each point of the cube, map x/y/z to match the defined hardware wiring.
	for x := 0; x < dst.XLen; x++ {
		for y := 0; y < dst.YLen; y++ {
			for z := 0; z < dst.ZLen; z++ {
				if src.GetVoxel(x, y, z) {
					dst.SetVoxel(xMap[z][x], yMap[x][y], zMap[x][z])
				}
			}
		}
	}
	return dst
}

func renderCube(spiHandler spi.Connection, c cube.Cube) error {
	c = mapCube(c)

	// We send one plane at a time to SPI. 1 cathode, ZLen anones.
	// TODO: Handle XLen > 8.
	tx := make([]byte, 1+c.ZLen)

	// The first "word" is the cathode layer, which is the Y axis.
	for y := 0; y < c.YLen; y++ {
		tx[0] = 0x01 << uint(y)

		for z := 0; z < c.ZLen; z++ {
			tx[z+1] = 0
			for x := 0; x < c.XLen; x++ {
				if c.GetVoxel(x, y, z) {
					tx[z+1] |= 0x01 << uint(x)
				}
			}
		}
		if err := spiHandler.Tx(tx, nil); err != nil {
			return errors.Wrap(err, "spi.Tx")
		}
	}

	return nil
}

func main() {
	if err := setup(); err != nil {
		log.Fatalf("setup error: %s\n\n%+v\n", errors.Cause(err), err)
	}

	work := func() {
		gobot.Every(5*time.Microsecond, func() {
			if err := loop(); err != nil {
				log.Fatalf("loop error: %s\n\n%+v\n", errors.Cause(err), err)
			}
		})
	}

	robot := gobot.NewRobot("spicube",
		work,
	)

	if err := robot.Start(); err != nil {
		log.Fatalf("robot.Start error: %s\n", err)
	}
}
