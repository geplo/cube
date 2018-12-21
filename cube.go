package cube

// Element represent the X axis. Z and Y are from the 2d array state.
type Element byte

// Cube is the in-memory representation of the c.
type Cube struct {
	state [][]Element // Actual state.
	XLen  int
	YLen  int
	ZLen  int
}

// New instantiate a simple cube of the given size.
func New(size int) Cube {
	state := make([][]Element, size)
	for i := 0; i < size; i++ {
		state[i] = make([]Element, size)
	}
	return Cube{
		state: state,
		XLen:  size,
		YLen:  size,
		ZLen:  size,
	}
}

// NewCustom instantiate a rectangular prism.
func NewCustom(xLen, yLen, zLen int) Cube {
	state := make([][]Element, yLen)
	for i := 0; i < yLen; i++ {
		state[i] = make([]Element, zLen)
	}
	return Cube{
		state: state,
		XLen:  xLen,
		YLen:  yLen,
		ZLen:  zLen,
	}
}

// SetVoxel turns on the given point in the cube.
func (c Cube) SetVoxel(x, y, z int) {
	c.state[c.YLen-1-y][c.ZLen-1-z] |= 0x01 << uint(x)
}

// GetVoxel returns the value of the requested point in the cube.
func (c Cube) GetVoxel(x, y, z int) bool {
	return (c.state[c.YLen-1-y][c.ZLen-1-z] & (0x01 << uint(x))) == (0x01 << uint(x))
}

// Clear turns off the whole cube,
func (c Cube) Clear() {
	for y, line := range c.state {
		for z := range line {
			c.state[y][z] = 0x00
		}
	}
}

// Axis enum type.
type Axis int

// Axis enum values.
const (
	AxisX Axis = iota
	AxisY
	AxisZ
)

// SetPlane turns on the Nth plane following then given axis.
func (c Cube) SetPlane(axis Axis, n int) {
	for j := 0; j < 8; j++ {
		for k := 0; k < 8; k++ {
			switch axis {
			case AxisX:
				c.SetVoxel(n, j, k)
			case AxisY:
				c.SetVoxel(j, n, k)
			case AxisZ:
				c.SetVoxel(j, k, n)
			default:
				panic("invalid axis")
			}
		}
	}
}

// Direction enum type.
type Direction bool

// Direction enum values.
const (
	Pos Direction = true
	Neg Direction = false
)

// AxisVector represent an axis with a direction.
type AxisVector struct {
	Axis
	Direction
}

// AxisVector enum values.
var (
	PosX = AxisVector{AxisX, Pos}
	NegX = AxisVector{AxisX, Neg}
	PosY = AxisVector{AxisY, Pos}
	NegY = AxisVector{AxisY, Neg}
	PosZ = AxisVector{AxisZ, Pos}
	NegZ = AxisVector{AxisZ, Neg}
)

// Shift translates the cube following the given direction.
func (c Cube) Shift(dir AxisVector) {
	switch dir {
	case PosX:
		for y := 0; y < c.YLen; y++ {
			for z := 0; z < c.ZLen; z++ {
				c.state[y][z] <<= 1
			}
		}
	case NegX:
		for y := 0; y < c.YLen; y++ {
			for z := 0; z < c.ZLen; z++ {
				c.state[y][z] >>= 1
			}
		}
	case PosY:
		for y := 1; y < c.YLen; y++ {
			for z := 0; z < c.ZLen; z++ {
				c.state[y-1][z] = c.state[y][z]
			}
		}
		for i := 0; i < c.ZLen; i++ {
			c.state[c.YLen-1][i] = 0
		}
	case NegY:
		for y := c.YLen - 1; y > 0; y-- {
			for z := 0; z < c.ZLen; z++ {
				c.state[y][z] = c.state[y-1][z]
			}
		}
		for i := 0; i < c.ZLen; i++ {
			c.state[0][i] = 0
		}
	case PosZ:
		for y := 0; y < c.YLen; y++ {
			for z := 1; z < c.ZLen; z++ {
				c.state[y][z-1] = c.state[y][z]
			}
		}
		for i := 0; i < c.YLen; i++ {
			c.state[i][c.ZLen-1] = 0
		}
	case NegZ:
		for y := 0; y < c.YLen; y++ {
			for z := c.ZLen - 1; z > 0; z-- {
				c.state[y][z] = c.state[y][z-1]
			}
		}
		for i := 0; i < c.YLen; i++ {
			c.state[c.YLen][0] = 0
		}
	default:
		panic("invalid direction")
	}
}
