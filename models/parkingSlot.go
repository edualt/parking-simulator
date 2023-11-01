package models

import (
	"github.com/oakmound/oak/v4/alg/floatgeom"
)

type ParkingSlot struct {
	area                 *floatgeom.Rect2
	isAvailable          bool
	directionsForParking []*struct {
		Direction string
		Location  float64
	}
	directionsForLeaving []*struct {
		Direction string
		Location  float64
	}
}

func NewParkingSlot(x, y, x2, y2 float64, column int) *ParkingSlot {
	directionsForParking := getDirectionForParking(x, y, column)
	directionsForLeaving := getDirectionsForLeaving()
	area := floatgeom.NewRect2(x, y, x2, y2)

	return &ParkingSlot{
		area:                 &area,
		isAvailable:          true,
		directionsForParking: directionsForParking,
		directionsForLeaving: directionsForLeaving,
	}
}

func getDirectionForParking(x, y float64, column int) []*struct {
	Direction string
	Location  float64
} {
	var directions []*struct {
		Direction string
		Location  float64
	}

	leftLocation := map[int]float64{
		1: 445,
		2: 355,
		3: 265,
		4: 175,
	}

	if location, ok := leftLocation[column]; ok {
		directions = append(directions, &struct {
			Direction string
			Location  float64
		}{
			"left", location,
		})
	}

	directions = append(directions, &struct {
		Direction string
		Location  float64
	}{
		"down", y + 5,
	})
	directions = append(directions, &struct {
		Direction string
		Location  float64
	}{
		"left", x + 5,
	})

	return directions
}

func getDirectionsForLeaving() []*struct {
	Direction string
	Location  float64
} {
	var directions []*struct {
		Direction string
		Location  float64
	}

	directions = append(directions, &struct {
		Direction string
		Location  float64
	}{
		"down", 425,
	})
	directions = append(directions, &struct {
		Direction string
		Location  float64
	}{
		"right", 475,
	})
	directions = append(directions, &struct {
		Direction string
		Location  float64
	}{
		"up", 135,
	})

	return directions
}

func (p *ParkingSlot) GetArea() *floatgeom.Rect2 {
	return p.area
}

func (p *ParkingSlot) GetDirectionsForParking() []*struct {
	Direction string
	Location  float64
} {
	return p.directionsForParking
}

func (p *ParkingSlot) GetDirectionsForLeaving() []*struct {
	Direction string
	Location  float64
} {
	return p.directionsForLeaving
}

func (p *ParkingSlot) GetIsAvailable() bool {
	return p.isAvailable
}

func (p *ParkingSlot) SetIsAvailable(isAvailable bool) {
	p.isAvailable = isAvailable
}
