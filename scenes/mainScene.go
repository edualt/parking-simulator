package scenes

import (
	"image/color"
	"log"
	"math/rand"
	"simulador/models"
	"sync"
	"time"

	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
)

var (
	rows    = 4
	columns = 5

	parking   = models.NewParking()
	gateMutex sync.Mutex
)

func MainScene(ctx *scene.Context) {
	prepareParking(ctx)

	event.GlobalBind(ctx, event.Enter, func(enterPayload event.EnterPayload) event.Response {
		for i := 0; i < 100; i++ {
			go run(ctx)
			sleepRandomDuration(1000, 2000)
		}
		return 0
	})
}

func prepareParking(ctx *scene.Context) {
	backgroundRender, err := render.LoadSprite("assets/images/OIP.jpg")
	if err != nil {
		log.Fatal(err)
	}

	entities.New(
		ctx,
		// make full screen
		entities.WithRenderable(backgroundRender),
		entities.WithDrawLayers([]int{-1}),
	)

	area := floatgeom.NewRect2(100, 130, 500, 455)
	entities.New(ctx, entities.WithRect(area), entities.WithColor(color.RGBA{86, 101, 115, 255}), entities.WithDrawLayers([]int{0}))

	leftLine := render.NewLine(445, 120, 445, 130, color.White)
	render.Draw(leftLine, 0)
	rightLine := render.NewLine(500, 120, 500, 130, color.White)
	render.Draw(rightLine, 0)

	// slot lines have to be in layer 1
	for _, spot := range parking.GetSlots() {
		area := spot.GetArea()
		areaX1 := area.Min.X()
		areaY1 := area.Min.Y()
		areaX2 := area.Max.X()
		areaY2 := area.Max.Y()

		topLine := render.NewLine(areaX1, areaY1, areaX2, areaY1, color.White)
		render.Draw(topLine, 0)

		leftLine := render.NewLine(areaX1, areaY1, areaX1, areaY2, color.White)
		render.Draw(leftLine, 0)

		bottomLine := render.NewLine(areaX1, areaY2, areaX2, areaY2, color.White)
		render.Draw(bottomLine, 0)
	}
}

func run(ctx *scene.Context) {
	c := createCar(ctx)
	parkCar(c)
	exitCar(c)
	removeCar(c)
}

func createCar(ctx *scene.Context) *models.Car {
	c := models.NewCar(ctx)
	models.AddCar(c)
	c.AddToQueue()
	return c
}

func parkCar(c *models.Car) {
	spotAvailable := parking.GetAvailableParkingSlot()
	withMutex(&gateMutex, c.EntryParking)
	c.Park(spotAvailable)
	sleepRandomDuration(40000, 50000)
	c.LeaveSpot()
	parking.MakeParkingSlotAvailable(spotAvailable)
	c.Leave(spotAvailable)
}

func exitCar(c *models.Car) {
	withMutex(&gateMutex, c.ExitParking)
	c.GoAway()
}

func removeCar(c *models.Car) {
	c.Remove()
	models.RemoveCar(c)
}

func withMutex(m *sync.Mutex, action func()) {
	m.Lock()
	defer m.Unlock()
	action()
}

func sleepRandomDuration(min, max int) {
	randomGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomDuration := time.Duration(randomGen.Intn(max-min+1) + min)
	time.Sleep(time.Millisecond * randomDuration)
}
