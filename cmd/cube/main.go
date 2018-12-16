package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/geplo/cube"
	"github.com/geplo/cube/scenes/planeshift"
	"github.com/pkg/errors"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/spi"
	"gobot.io/x/gobot/platforms/raspi"
)

var (
	adaptor    gobot.Connection
	spiHandler spi.Connection
	c          cube.Cube
	mappedCube cube.Cube
	// scene      interface{ Step(cube.Cube) time.Duration }
	scene func(cube.Cube)
)

// setup is called before the main loop.
func setup() error {
	// Instantiate the cube.
	c = cube.New(8)
	// Pre-alloc the hardware mapped cube.
	mappedCube = cube.New(8)

	// Use a raspberry pi.
	a := raspi.NewAdaptor()

	// Initialize SPI.
	hdlr, err := a.GetSpiConnection(
		0,   // Bus 0.
		0,   // CEO0.
		0,   // SPI_MODE_0.
		8,   // 8 bits per words.
		8e6, // 8MHz.
	)
	if err != nil {
		return errors.Wrap(err, "GetSpiConnection")
	}
	adaptor = a
	spiHandler = hdlr

	// Seed the random generator.
	rand.Seed(time.Now().UnixNano())

	// Set the scene to use.
	scene = planeshift.Scene
	return nil
}

var cc = 0

var next time.Time

func loop() error {
	start := time.Now()

	// if next.Before(start) {
	// Step the scene.
	scene(c)

	// Render it.
	if err := renderCube(spiHandler, c, mappedCube); err != nil {
		return errors.Wrap(err, "renderCube")
	}
	// }

	if cc%200 == 0 {
		fmt.Printf("%s\n", time.Since(start))
	}
	cc++

	//time.Sleep(time.Until(next))
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

func mapCube(src, dst cube.Cube) cube.Cube {
	// Make sure the dst is empty.
	dst.Clear()

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

func renderCube(spiHandler spi.Connection, c, mappedCube cube.Cube) error {
	c = mapCube(c, mappedCube)

	// We send one plane at a time to SPI. 1 cathode, ZLen anones.
	// TODO: Handle XLen > 8.
	tx := make([]byte, 1+c.ZLen)
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
		[]gobot.Connection{adaptor},
		work,
	)

	if err := robot.Start(); err != nil {
		log.Fatalf("robot.Start error: %s\n", err)
	}
}
