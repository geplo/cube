package cube

type Element byte

// Cube is the in-memory representation of the c.
type Cube struct {
	state [][]Element // Actual state.
	XLen  int
	YLen  int
	ZLen  int
}

// New instantiate a simple cube of the given size.
// Max size is 64.
func New(size int) Cube {
	if size > 64 {
		panic("max size for cube is 64")
	}
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
// xLen max is 64.
func NewCustom(xLen, yLen, zLen int) Cube {
	if xLen > 64 {
		panic("max size for zLen is 64")
	}
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

func (c Cube) SetVoxel(x, y, z int) {
	c.state[c.YLen-1-y][c.ZLen-1-z] |= 0x01 << uint(x)
}

func (c Cube) GetVoxel(x, y, z int) bool {
	return (c.state[c.YLen-1-y][c.ZLen-1-z] & (0x01 << uint(x))) == (0x01 << uint(x))
}

func (c Cube) Clear() {
	for y, line := range c.state {
		for z := range line {
			c.state[y][z] = 0x00
		}
	}
}

// Direction .
type Direction int

const (
	PosX Direction = iota
	NegX
	PosY
	NegY
	PosZ
	NegZ
)

func (c Cube) Shift(dir Direction) {
	switch dir {
	case PosX:
		for y := 0; y < 8; y++ {
			for z := 0; z < 8; z++ {
				c.state[y][z] = c.state[y][z] << 1
			}
		}
	case NegX:
		for y := 0; y < 8; y++ {
			for z := 0; z < 8; z++ {
				c.state[y][z] = c.state[y][z] >> 1
			}
		}
	case PosY:
		for y := 1; y < 8; y++ {
			for z := 0; z < 8; z++ {
				c.state[y-1][z] = c.state[y][z]
			}
		}
		for i := 0; i < 8; i++ {
			c.state[7][i] = 0
		}
	case NegY:
		for y := 7; y > 0; y-- {
			for z := 0; z < 8; z++ {
				c.state[y][z] = c.state[y-1][z]
			}
		}
		for i := 0; i < 8; i++ {
			c.state[0][i] = 0
		}
	case PosZ:
		for y := 0; y < 8; y++ {
			for z := 1; z < 8; z++ {
				c.state[y][z-1] = c.state[y][z]
			}
		}
		for i := 0; i < 8; i++ {
			c.state[i][7] = 0
		}
	case NegZ:
		for y := 0; y < 8; y++ {
			for z := 7; z > 0; z-- {
				c.state[y][z] = c.state[y][z-1]
			}
		}
		for i := 0; i < 8; i++ {
			c.state[i][0] = 0
		}
	default:
		panic("invalid direction")
	}
}

type Axis int

const (
	AxisX Axis = iota
	AxisY
	AxisZ
)

func (c Cube) SetPlane(axis Axis, i int) {
	for j := 0; j < 8; j++ {
		for k := 0; k < 8; k++ {
			switch axis {
			case AxisX:
				c.SetVoxel(i, j, k)
			case AxisY:
				c.SetVoxel(j, i, k)
			case AxisZ:
				c.SetVoxel(j, k, i)
			default:
				panic("invalid axis")
			}
		}
	}
}
