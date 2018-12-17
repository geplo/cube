package cube

import (
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

// Driver .
type Driver struct {
	// gobot interface.
	name       string
	connection spi.Connection

	cube   Cube
	mapped Cube
	scene  scenes.Scene
	next   time.Time

	// Hardware mapping.
	xMap [][]int
	yMap [][]int
	zMap [][]int
}

// WithXMap .
func (d *Driver) WithXMap(xMap [][]int) { d.xMap = xMap }

// WithYMap .
func (d *Driver) WithYMap(yMap [][]int) { d.yMap = yMap }

// WithZMap .
func (d *Driver) WithZMap(zMap [][]int) { d.zMap = zMap }

// WithXMap .
func WithXMap(xMap [][]int) func(c Config) { return func(c Config) { c.WithXMap(xMap) } }

// WithYMap .
func WithYMap(yMap [][]int) func(c Config) { return func(c Config) { c.WithYMap(yMap) } }

// WithZMap .
func WithZMap(zMap [][]int) func(c Config) { return func(c Config) { c.WithZMap(zMap) } }

// Config .
type Config interface {
	WithXMap([][]int)
	WithYMap([][]int)
	WithZMap([][]int)
}

// NewDriver .
func NewDriver(options ...func(Config)) *Driver {
	d := &Driver{
		name: gobot.DefaultName("ledcube"),
	}
	for _, option := range options {
		option(d)
	}
	return d
}

// Connection returns the Connection of the device.
// Implements gobot.Driver / gobot.Device interface.
func (d *Driver) Connection() gobot.Connection { return d.connection.(gobot.Connection) }

// Start initializes the driver.
// Implements gobot.Driver / gobot.Device interface.
func (d *Driver) Start() (err error) {
	// Instantiate the cube.
	c = New(8)

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
	d.connection = hdlr

	// Seed the random generator.
	rand.Seed(time.Now().UnixNano())

	// Set the scene to use.
	scene = planeshift.New()
	scene = rain.New()

	return nil
}

// Halt stops the driver.
// Implements gobot.Driver / gobot.Device interface.
func (d *Driver) Halt() error {
	return d.connection.Close()
}

// Name returns the name of the device.
// Implements gobot.Driver / gobot.Device interface.
func (d *Driver) Name() string { return d.name }

// SetName sets the name of the device.
// Implements gobot.Driver / gobot.Device interface.
func (d *Driver) SetName(n string) { d.name = n }

// Render displays the cube.
func (d *Driver) Render() error {
	start := time.Now()

	if d.next.Before(start) {
		// Step the scene.
		d.next = start.Add(d.scene.Step(c))
	}

	// Render the cube.
	if err := d.renderCube(); err != nil {
		return errors.Wrap(err, "renderCube")
	}

	return nil
}

func (d *Drivers) mapCube() cube.Cube {
	// Make sure the dst is empty.
	dst := NewCustom(d.cube.XLen, d.cube.YLen, d.cube.ZLen)

	// For each point of the cube, map x/y/z to match the defined hardware wiring.
	for x := 0; x < dst.XLen; x++ {
		for y := 0; y < dst.YLen; y++ {
			for z := 0; z < dst.ZLen; z++ {
				if src.GetVoxel(x, y, z) {
					dst.SetVoxel(d.xMap[z][x], d.yMap[x][y], d.zMap[x][z])
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
