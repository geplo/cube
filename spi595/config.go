package spi595

import "gobot.io/x/gobot/drivers/spi"

type config struct {
	spi.Config

	// Hardware mapping.
	xMap [][]int
	yMap [][]int
	zMap [][]int
}

// Config .
type Config interface {
	spi.Config

	WithXMap([][]int)
	WithYMap([][]int)
	WithZMap([][]int)

	XMap(z, x int) int
	YMap(x, y int) int
	ZMap(x, z int) int
}

// NewConfig .
func NewConfig() Config {
	return &config{
		Config: spi.NewConfig(),
	}
}

// WithXMap .
func (c *config) WithXMap(xMap [][]int) { c.xMap = xMap }

// WithYMap .
func (c *config) WithYMap(yMap [][]int) { c.yMap = yMap }

// WithZMap .
func (c *config) WithZMap(zMap [][]int) { c.zMap = zMap }

// WithXMap .
func WithXMap(xMap [][]int) func(c Config) { return func(c Config) { c.WithXMap(xMap) } }

// WithYMap .
func WithYMap(yMap [][]int) func(c Config) { return func(c Config) { c.WithYMap(yMap) } }

// WithZMap .
func WithZMap(zMap [][]int) func(c Config) { return func(c Config) { c.WithZMap(zMap) } }

// XMap .
func (c config) XMap(z, x int) int {
	if c.xMap == nil {
		return x
	}
	return c.xMap[z][x]
}

// YMap .
func (c config) YMap(x, y int) int {
	if c.yMap == nil {
		return y
	}
	return c.xMap[x][y]
}

// ZMap .
func (c config) ZMap(x, z int) int {
	if c.zMap == nil {
		return z
	}
	return c.zMap[x][z]
}

// WithBus sets which bus to use as a optional param.
func WithBus(bus int) func(Config) { return func(c Config) { c.WithBus(bus) } }

// WithChip sets which chip to use as a optional param.
func WithChip(chip int) func(Config) { return func(c Config) { c.WithChip(chip) } }

// WithMode sets which mode to use as a optional param.
func WithMode(mode int) func(Config) { return func(c Config) { c.WithMode(mode) } }

// WithBits sets how many bits to use as a optional param.
func WithBits(bits int) func(Config) { return func(c Config) { c.WithBits(bits) } }

// WithSpeed sets what speed to use as a optional param.
func WithSpeed(speed int64) func(Config) { return func(c Config) { c.WithSpeed(speed) } }
