package models

type Car struct {
	Id           int
	IsWaiting    bool
	IsParked     bool
	IsEntering   bool
	IsGettingOut bool
}

func NewCar(id int, isWaiting, isParked, isEntering, isGettingOut bool) *Car {
	return &Car{
		Id:           id,
		IsWaiting:    isWaiting,
		IsParked:     isParked,
		IsEntering:   isEntering,
		IsGettingOut: isGettingOut,
	}
}

func (c *Car) GetID() int {
	return c.Id
}

func (c *Car) GetStatus() string {
	status := "Unknown"
	if c.IsWaiting {
		status = "Waiting"
	} else if c.IsParked {
		status = "Parked"
	} else if c.IsEntering {
		status = "Entering"
	} else if c.IsGettingOut {
		status = "Getting Out"
	}
	return status
}
