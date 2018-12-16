package planeshift

import (
	"math/rand"

	"github.com/geplo/cube"
)

var (
	timer          int
	looped         bool
	loading        bool
	planeDirection cube.Direction
	planePosition  int
)

// Scene describes the plane shift scene.
func Scene(c cube.Cube) {
	if loading {
		c.Clear()
		axis := cube.Axis(rand.Int() % 3)
		planePosition = (rand.Int()) % 2 * 7
		c.SetPlane(axis, planePosition)
		switch axis {
		case cube.AxisX:
			if planePosition == 0 {
				planeDirection = cube.PosX
			} else {
				planeDirection = cube.NegX
			}
		case cube.AxisY:
			if planePosition == 0 {
				planeDirection = cube.PosY
			} else {
				planeDirection = cube.NegY
			}
		case cube.AxisZ:
			if planePosition == 0 {
				planeDirection = cube.PosZ
			} else {
				planeDirection = cube.NegZ
			}
		}
		timer = 0
		looped = false
		loading = false
	}

	timer++
	if timer > 220 {
		timer = 0
		c.Shift(planeDirection)
		if planeDirection%2 == 0 {
			planePosition++
			if planePosition == 7 {
				if looped {
					loading = true
				} else {
					planeDirection++
					looped = true
				}
			}
		} else {
			planePosition--
			if planePosition == 0 {
				if looped {
					loading = true
				} else {
					planeDirection--
					looped = true
				}
			}
		}
	}
}

/*
package planeshift

import (
	"math/rand"
	"time"

	"github.com/geplo/cube"
)

type Scene struct {
	timer          int
	looped         bool
	loading        bool
	planeDirection cube.Direction
	planePosition  int
}

func New() Scene {
	return Scene{loading: true}
}

// Scene describes the plane shift scene.
func (s Scene) Step(c cube.Cube) time.Duration {
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
		s.timer = 0
		s.looped = false
		s.loading = false
	}

	s.timer++
	if s.timer > 2200 {
		s.timer = 0
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
	}
	return 0
}
*/
