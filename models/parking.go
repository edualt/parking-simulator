package models

import (
	"image/color"
	"sync"

	"github.com/oakmound/oak/v4/render"
)

var (
	rows    = 4
	columns = 5
)

type Parking struct {
	slots         []*ParkingSlot
	queueCars     []Car
	mutex         sync.Mutex
	availableCond *sync.Cond
}

func NewParking() *Parking {
	slots := make([]*ParkingSlot, 4*5)
	queueCars := make([]Car, 0)
	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			x := 410 - i*90
			y := 210 + j*45
			slots[i*columns+j] = NewParkingSlot(float64(x), float64(y), float64(x+30), float64(y+30), i+1)
			line := render.NewLine(float64(x), float64(y), float64(x+30), float64(y+30), color.White)
			render.Draw(line, 2)
		}
	}

	p := &Parking{
		slots:     slots,
		queueCars: queueCars,
	}
	p.availableCond = sync.NewCond(&p.mutex)
	return p
}

func (p *Parking) GetSlots() []*ParkingSlot {
	return p.slots
}

func (p *Parking) safeExecute(action func()) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	action()
}

func (p *Parking) GetAvailableParkingSlot() *ParkingSlot {
	var spot *ParkingSlot
	p.safeExecute(func() {
		spot = p.findAvailableSpot()
		if spot != nil {
			spot.SetIsAvailable(false)
		}
	})
	return spot
}

func (p *Parking) findAvailableSpot() *ParkingSlot {
	for _, spot := range p.slots {
		if spot.GetIsAvailable() {
			return spot
		}
	}
	p.availableCond.Wait()
	return p.findAvailableSpot()
}

func (p *Parking) MakeParkingSlotAvailable(spot *ParkingSlot) {
	p.safeExecute(func() {
		spot.SetIsAvailable(true)
		p.availableCond.Signal()
	})
}

func (p *Parking) GetCarsInQueue() (cars []Car) {
	p.safeExecute(func() {
		cars = p.queueCars
	})
	return cars
}
