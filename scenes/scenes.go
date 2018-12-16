package scenes

import (
	"time"

	"github.com/geplo/cube"
)

// Scene represent a simple step by step scenario.
type Scene interface {
	Step(cube.Cube) time.Duration
}
