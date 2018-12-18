package spi595

import (
	"github.com/geplo/cube"
	"github.com/pkg/errors"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/spi"
)

// Adaptor .
type Adaptor struct {
	name string
	Config

	// Parent connector.
	connector spi.Connector

	// SPI handler.
	connection spi.Connection
}

// NewAdaptor .
func NewAdaptor(parent spi.Connector, options ...func(Config)) *Adaptor {
	a := &Adaptor{
		name:      gobot.DefaultName("ledCube"),
		Config:    NewConfig(),
		connector: parent,
	}
	for _, option := range options {
		option(a)
	}
	return a
}

// Name implements the gobot.Adaptor / gobot.Connection interface.
func (a *Adaptor) Name() string { return a.name }

// SetName implements the gobot.Adaptor / gobot.Connection interface.
func (a *Adaptor) SetName(name string) { a.name = name }

// Connect implements the gobot.Adaptor / gobot.Connection interface.
func (a *Adaptor) Connect() error {
	println("adaptor start")
	var (
		bus              = a.GetBusOrDefault(0)     // Bus 0 (/dev/spidev0).
		chip             = a.GetChipOrDefault(0)    // CEO0 (/dev/spidev0.0).
		mode             = a.GetModeOrDefault(0)    // SPI_MODE_0.
		bits             = a.GetBitsOrDefault(8)    // 8 bits per word.
		maxSpeed         = a.GetSpeedOrDefault(8e6) // 8MHz.
		getSpiConnection = spi.GetSpiConnection
	)

	// If we have a parent connector, use it.
	if a.connector != nil {
		bus = a.GetBusOrDefault(a.connector.GetSpiDefaultBus())
		chip = a.GetChipOrDefault(a.connector.GetSpiDefaultChip())
		mode = a.GetModeOrDefault(a.connector.GetSpiDefaultMode())
		bits = a.GetBitsOrDefault(a.connector.GetSpiDefaultBits())
		maxSpeed = a.GetSpeedOrDefault(a.connector.GetSpiDefaultMaxSpeed())
		getSpiConnection = a.connector.GetSpiConnection
	}

	// Initialize SPI.
	hdlr, err := getSpiConnection(bus, chip, mode, bits, maxSpeed)
	if err != nil {
		return errors.Wrap(err, "GetSpiConnection")
	}
	a.connection = hdlr

	return nil
}

// Finalize implements the gobot.Adaptor / gobot.Connection interface.
func (a *Adaptor) Finalize() error {
	println("adaptor stop")
	return a.connection.Close()
}

func (a *Adaptor) mapCube(src cube.Cube) cube.Cube {
	dst := cube.NewCustom(src.XLen, src.YLen, src.ZLen)

	// For each point of the cube, map x/y/z to match the defined hardware wiring.
	for x := 0; x < dst.XLen; x++ {
		for y := 0; y < dst.YLen; y++ {
			for z := 0; z < dst.ZLen; z++ {
				if src.GetVoxel(x, y, z) {
					dst.SetVoxel(a.XMap(z, x), a.YMap(x, y), a.ZMap(x, z))
				}
			}
		}
	}

	return dst
}

func (a *Adaptor) renderCube(c cube.Cube) error {
	c = a.mapCube(c)

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
		if err := a.connection.Tx(tx, nil); err != nil {
			return errors.Wrap(err, "spi.Tx")
		}
	}

	return nil
}
