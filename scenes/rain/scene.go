package rain

import (
	"math/rand"
	"time"

	"github.com/geplo/cube"
	"github.com/geplo/cube/scenes"
)

// Scene holds the state.
type Scene struct {
	loading bool
}

// New instantiate a new scene.
func New() scenes.Scene {
	return &Scene{
		loading: true,
	}
}

// Step implements the scenes.Scene interface.
func (s *Scene) Step(c cube.Cube) time.Duration {
	if s.loading {
		c.Clear()
		s.loading = false
	}
	c.Shift(cube.NegY)
	numDrops := rand.Int() % 5
	for i := 0; i < numDrops; i++ {
		c.SetVoxel(rand.Int()%c.XLen, c.YLen-1, rand.Int()%c.ZLen)
	}

	return 60 * time.Millisecond
}
