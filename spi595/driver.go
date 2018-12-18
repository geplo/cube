package spi595

import (
	"math/rand"
	"sync"
	"time"

	"github.com/geplo/cube"
	"github.com/geplo/cube/scenes"
	"github.com/pkg/errors"

	"gobot.io/x/gobot"
)

// Events / Commands.
const (
	Step  = "step"
	Scene = "scene"
)

// Driver .
type Driver struct {
	sync.Mutex
	name       string
	connection *Adaptor

	gobot.Eventer
	gobot.Commander

	cube.Cube
	scene scenes.Scene
	next  time.Time
}

// NewDriver .
func NewDriver(a *Adaptor, cube cube.Cube, scene scenes.Scene) *Driver {
	d := &Driver{
		name:       gobot.DefaultName("ledcube"),
		connection: a,
		Eventer:    gobot.NewEventer(),
		Commander:  gobot.NewCommander(),

		Cube:  cube,
		scene: scene,
	}
	d.AddEvent(Step)
	d.AddCommand(Step, func(map[string]interface{}) interface{} { return d.Step() })
	d.AddCommand(Scene, func(params map[string]interface{}) interface{} { d.Scene(params[Scene].(scenes.Scene)); return nil })
	return d
}

// Step .
func (d *Driver) Step() time.Duration {
	d.Lock()
	defer d.Unlock()
	next := d.scene.Step(d.Cube)
	d.next = time.Now().Add(next)
	return next
}

// Scene .
func (d *Driver) Scene(scene scenes.Scene) {
	d.Lock()
	defer d.Unlock()
	d.scene = scene
	d.next = time.Now()
}

// Connection returns the Connection of the device.
// Implements gobot.Driver / gobot.Device interface.
func (d *Driver) Connection() gobot.Connection { return d.connection }

// Start initializes the driver.
// Implements gobot.Driver / gobot.Device interface.
func (d *Driver) Start() (err error) {
	// Seed the random generator.
	rand.Seed(time.Now().UnixNano())
	println("driver start")
	return nil
}

// Halt stops the driver.
// Implements gobot.Driver / gobot.Device interface.
func (d *Driver) Halt() error {
	println("driver stop")
	d.Cube.Clear()
	return nil
}

// Name returns the name of the device.
// Implements gobot.Driver / gobot.Device interface.
func (d *Driver) Name() string { return d.name }

// SetName sets the name of the device.
// Implements gobot.Driver / gobot.Device interface.
func (d *Driver) SetName(n string) { d.name = n }

// Draw displays the cube.
func (d *Driver) Draw() error {
	start := time.Now()

	if d.next.Before(start) {
		d.Publish(d.Event(Step), start)
		d.next = start.Add(d.scene.Step(d.Cube))
	}

	if err := d.connection.renderCube(d.Cube); err != nil {
		return errors.Wrap(err, "renderCube")
	}
	return nil
}
