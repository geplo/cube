package cube888

import (
	"math/rand"
	"time"

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
		name:      gobot.DefaultName("Cube888"),
		connector: a,
		Config:    spi.NewConfig(),
	}
	for _, option := range options {
		option(d)
	}
	return d
}

// Connection returns the Connection of the device.
func (d *Driver) Connection() gobot.Connection { return d.connection.(gobot.Connection) }

// Start initializes the driver.
func (d *Driver) Start() (err error) {
	bus := d.GetBusOrDefault(d.connector.GetSpiDefaultBus())
	chip := d.GetChipOrDefault(d.connector.GetSpiDefaultChip())
	mode := d.GetModeOrDefault(d.connector.GetSpiDefaultMode())
	bits := d.GetBitsOrDefault(d.connector.GetSpiDefaultBits())
	maxSpeed := d.GetSpeedOrDefault(d.connector.GetSpiDefaultMaxSpeed())

	d.connection, err = d.connector.GetSpiConnection(bus, chip, mode, bits, maxSpeed)
	if err != nil {
		return err
	}

	d.loading = true
	d.randomTimer = 0
	rand.Seed(time.Now().UnixNano())

	return nil
}

// Halt stops the driver.
func (d *Driver) Halt() (err error) {
	d.connection.Close()
	return
}

// Name returns the name of the device.
func (d *Driver) Name() string { return d.name }

// SetName sets the name of the device.
func (d *Driver) SetName(n string) { d.name = n }

// Draw displays the cube.
func (d *Driver) Draw() error {
	tx := [9]byte{}

	for i := (uint)(0); i < 8; i++ {
		tx[0] = 0x01 << i
		for j := 0; j < 8; j++ {
			tx[j+1] = d.cube[7-i][j+1]
			tx[j+2] = d.cube[7-i][j]
			j++
		}
		if err := d.connection.Tx(tx[:], nil); err != nil {
			return err
		}
	}

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

type shiftType int

const (
	posX shiftType = iota
	negX
	posY
	negY
	posZ
	negZ
)

func (d *Driver) shift(dir shiftType) {
	switch dir {
	case posX:
		for y := 0; y < 8; y++ {
			for z := 0; z < 8; z++ {
				d.cube[y][z] = d.cube[y][z] << 1
			}
		}
	case negX:
		for y := 0; y < 8; y++ {
			for z := 0; z < 8; z++ {
				d.cube[y][z] = d.cube[y][z] >> 1
			}
		}
	case posY:
		for y := 1; y < 8; y++ {
			for z := 0; z < 8; z++ {
				d.cube[y-1][z] = d.cube[y][z]
			}
		}
		for i := 0; i < 8; i++ {
			d.cube[7][i] = 0
		}
	case negY:
		for y := 7; y > 0; y-- {
			for z := 0; z < 8; z++ {
				d.cube[y][z] = d.cube[y-1][z]
			}
		}
		for i := 0; i < 8; i++ {
			d.cube[0][i] = 0
		}
	case posZ:
		for y := 0; y < 8; y++ {
			for z := 1; z < 8; z++ {
				d.cube[y][z-1] = d.cube[y][z]
			}
		}
		for i := 0; i < 8; i++ {
			d.cube[i][7] = 0
		}
	case negZ:
		for y := 0; y < 8; y++ {
			for z := 7; z > 0; z-- {
				d.cube[y][z] = d.cube[y][z-1]
			}
		}
		for i := 0; i < 8; i++ {
			d.cube[i][0] = 0
		}
	default:
		panic("invalid shift dir")
	}
}

func (d *Driver) setVoxel(x, y, z uint) {
	d.cube[7-y][7-z] |= 0x01 << x
}

type axisType int

const (
	xAxis axisType = iota
	yAxis
	zAxis
)

func (d *Driver) setPlane(axis axisType, i uint) {
	for j := uint(0); j < 8; j++ {
		for k := uint(0); k < 8; k++ {
			switch axis {
			case xAxis:
				d.setVoxel(i, j, k)
			case yAxis:
				d.setVoxel(j, i, k)
			case zAxis:
				d.setVoxel(j, k, i)
			default:
				panic("invalid axis")
			}
		}
	}
}

func (d *Driver) planeBoing() {
	if d.loading {
		d.clearCube()
		axis := axisType(rand.Int() % 3)
		d.planePosition = uint(rand.Int()) % 2 * 7
		d.setPlane(axis, d.planePosition)
		switch axis {
		case xAxis:
			if d.planePosition == 0 {
				d.planeDirection = posX
			} else {
				d.planeDirection = negX
			}
		case yAxis:
			if d.planePosition == 0 {
				d.planeDirection = posY
			} else {
				d.planeDirection = negY
			}
		case zAxis:
			if d.planePosition == 0 {
				d.planeDirection = posZ
			} else {
				d.planeDirection = negZ
			}
		default:
			panic("invalid random axis??")
		}
		d.timer = 0
		d.looped = false
		d.loading = false
	}

	d.timer++
	if d.timer > 220 {
		d.timer = 0
		d.shift(d.planeDirection)
		if d.planeDirection%2 == 0 {
			d.planePosition++
			if d.planePosition == 7 {
				if d.looped {
					d.loading = true
				} else {
					d.planeDirection++
					d.looped = true
				}
			}
		} else {
			d.planePosition--
			if d.planePosition == 0 {
				if d.looped {
					d.loading = true
				} else {
					d.planeDirection--
					d.looped = true
				}
			}
		}
	}
}
