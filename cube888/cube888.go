package cube888

import (
	"math/rand"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/spi"
)

type Driver struct {
	name       string
	connector  spi.Connector
	connection spi.Connection
	spi.Config

	cube        [8][8]uint8
	randomTimer uint16
	timer       uint64
	loading     bool

	planePosition  uint
	planeDirection shiftType
	looped         bool
}

func NewDriver(a spi.Connector, options ...func(spi.Config)) *Driver {
	d := &Driver{
		name:      gobot.DefaultName("spicube"),
		connector: a,
		Config:    spi.NewConfig(),
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
func (d *Driver) Start() error {
	var (
		bus   = d.GetBusOrDefault(d.connector.GetSpiDefaultBus())
		chip  = d.GetChipOrDefault(d.connector.GetSpiDefaultChip())
		mode  = d.GetModeOrDefault(d.connector.GetSpiDefaultMode())
		bits  = d.GetBitsOrDefault(d.connector.GetSpiDefaultBits())
		speed = d.GetSpeedOrDefault(d.connector.GetSpiDefaultMaxSpeed())
	)

	connection, err := d.connector.GetSpiConnection(bus, chip, mode, bits, speed)
	if err != nil {
		return err
	}
	d.connection = connection

	return nil
}

// Halt stops the driver.
// Implements gobot.Driver / gobot.Device interface.
func (d *Driver) Halt() (err error) {
	d.connection.Close()
	return
}

// Name returns the name of the device.
// Implements gobot.Driver / gobot.Device interface.
func (d *Driver) Name() string { return d.name }

// SetName sets the name of the device.
// Implements gobot.Driver / gobot.Device interface.
func (d *Driver) SetName(n string) { d.name = n }

// Transfer uses SPI to send tx and receive rx.
// tx and rx must be of the same size.
func (d *Driver) Transfer(tx, rx []byte) error {

	return nil
}

func (d *Driver) Step() {
	d.randomTimer++
	d.planeBoing()

	// for _, line := range d.cube {
	// 	for _, elem := range line {
	// 		fmt.Printf("0b%08b ", elem)
	// 	}
	// 	fmt.Printf("\n")
	// }
	// fmt.Printf("\n")
}

func (d *Driver) rain() {
	if d.loading {
		d.clearCube()
		d.loading = false
	}
	d.timer++
	if d.timer > 260 {
		d.timer = 0
		d.shift(d.planeDirection)
		numDrops := rand.Int() % 5
		for i := 0; i < numDrops; i++ {
			d.setVoxel(uint(rand.Int())%8, 7, uint(rand.Int())%8)
		}
	}
}

func (d *Driver) clearCube() {
	for i, elem := range d.cube {
		for j := range elem {
			d.cube[i][j] = 0
		}
	}
}
