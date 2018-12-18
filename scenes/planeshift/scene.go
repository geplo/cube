package planeshift

import (
	"math/rand"
	"time"

	"github.com/geplo/cube"
	"github.com/geplo/cube/scenes"
)

type Scene struct {
	looped         bool
	loading        bool
	planeDirection cube.Direction
	planePosition  int
}

func New() scenes.Scene {
	return &Scene{loading: true}
}

// Scene describes the plane shift scene.
func (s *Scene) Step(c cube.Cube) time.Duration {
	if s.loading {
		c.Clear()
		axis := cube.Axis(rand.Int() % 3)
		s.planePosition = (rand.Int()) % 2 * 7
		c.SetPlane(axis, s.planePosition)
		switch axis {
		case cube.AxisX:
			if s.planePosition == 0 {
				s.planeDirection = cube.PosX
			} else {
				s.planeDirection = cube.NegX
			}
		case cube.AxisY:
			if s.planePosition == 0 {
				s.planeDirection = cube.PosY
			} else {
				s.planeDirection = cube.NegY
			}
		case cube.AxisZ:
			if s.planePosition == 0 {
				s.planeDirection = cube.PosZ
			} else {
				s.planeDirection = cube.NegZ
			}
		}
		s.looped = false
		s.loading = false
	}

	c.Shift(s.planeDirection)
	if s.planeDirection%2 == 0 {
		s.planePosition++
		if s.planePosition == 7 {
			if s.looped {
				s.loading = true
			} else {
				s.planeDirection++
				s.looped = true
			}
		}
	} else {
		s.planePosition--
		if s.planePosition == 0 {
			if s.looped {
				s.loading = true
			} else {
				s.planeDirection--
				s.looped = true
			}
		}
	}

	return 100 * time.Millisecond
}
