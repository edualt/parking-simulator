package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/tu_paquete/models"
)

const (
	capacidadEstacionamiento = 20
	numVehiculos             = 100
)

var (
	cajones            [capacidadEstacionamiento]*models.Car
	wg                 sync.WaitGroup
	mutex              sync.Mutex
	entrada            = make(chan int, capacidadEstacionamiento)
	lugaresOcupados    int
	carrosEstacionados []*models.Car
)

func llegadaVehiculo(id int) {
	defer wg.Done()

	vehiculo := models.NewCar(id, true, false, false, false)

	fmt.Printf("Vehículo %d llega al estacionamiento.\n", vehiculo.GetID())

	// Intentar ingresar al estacionamiento
	entrada <- vehiculo.GetID()

	mutex.Lock()
	for i := 0; i < capacidadEstacionamiento; i++ {
		if cajones[i] == nil {
			cajones[i] = vehiculo
			vehiculo.IsWaiting = false
			vehiculo.IsEntering = true
			lugaresOcupados++
			carrosEstacionados = append(carrosEstacionados, vehiculo)
			mutex.Unlock()
			fmt.Printf("Vehículo %d estaciona en el cajón %d.\n", vehiculo.GetID(), i)
			vehiculo.IsEntering = false
			vehiculo.IsParked = true
			imprimirEstadoEstacionamiento()
			time.Sleep(time.Duration(rand.Intn(6)+5) * time.Second)
			mutex.Lock()
			cajones[i] = nil
			vehiculo.IsEntering = false
			vehiculo.IsParked = true
			lugaresOcupados--
			carrosEstacionados = eliminarCarro(carrosEstacionados, vehiculo)
			fmt.Printf("Vehículo %d con estado %s se va del cajón %d.\n", vehiculo.GetID(), vehiculo.GetStatus(), i)
			imprimirEstadoEstacionamiento()
			mutex.Unlock()
			<-entrada
			return
		}
	}
	mutex.Unlock()

	// El vehículo se bloquea esperando un cajón
	vehiculo.IsEntering = false
	vehiculo.IsWaiting = true
	fmt.Printf("Vehículo %d está bloqueado y haciendo fila.\n", vehiculo.GetID())

	mutex.Lock()
	fmt.Printf("Vehículo %d es el siguiente en intentar estacionar.\n", vehiculo.GetID())
	mutex.Unlock()

	// Intentar ingresar nuevamente cuando se desocupe un cajón
	mutex.Lock()
	for i := 0; i < capacidadEstacionamiento; i++ {
		if cajones[i] == nil {
			cajones[i] = vehiculo
			vehiculo.IsWaiting = false
			vehiculo.IsEntering = true
			lugaresOcupados++
			carrosEstacionados = append(carrosEstacionados, vehiculo)
			mutex.Unlock()
			fmt.Printf("Vehículo %d estaciona en el cajón %d.\n", vehiculo.GetID(), i)
			imprimirEstadoEstacionamiento()
			time.Sleep(time.Duration(rand.Intn(6)+5) * time.Second)
			mutex.Lock()
			cajones[i] = nil
			vehiculo.IsEntering = false
			vehiculo.IsParked = true
			lugaresOcupados--
			carrosEstacionados = eliminarCarro(carrosEstacionados, vehiculo)
			fmt.Printf("Vehículo %d se va del cajón %d.\n", vehiculo.GetID(), i)
			imprimirEstadoEstacionamiento()
			mutex.Unlock()
			<-entrada
			return
		}
	}
	mutex.Unlock()
}

func imprimirEstadoEstacionamiento() {
	app := app.New()

	w := app.NewWindow("Estacionamiento")

	var estadosVehiculos []string
	for _, carro := range carrosEstacionados {
		estado := ""
		if carro.IsWaiting {
			estado = "Esperando"
		} else if carro.IsEntering {
			estado = "Entrando"
		} else if carro.IsParked {
			estado = "Estacionado"
		} else if carro.IsGettingOut {
			estado = "Saliendo"
		}
		estadosVehiculos = append(estadosVehiculos, fmt.Sprintf("Vehículo %d (%s)", carro.GetID(), estado))
	}

	statusLabel := widget.NewLabel(fmt.Sprintf("Lugares ocupados: %d\nCarros estacionados:\n%s", lugaresOcupados, statesToString(estadosVehiculos)))

	w.SetContent(container.NewVBox(
		statusLabel,
	))

	w.ShowAndRun()
}

func statesToString(states []string) string {
	return fmt.Sprintf("- %s\n", states)
}

func eliminarCarro(slice []*models.Car, carro *models.Car) []*models.Car {
	for i, c := range slice {
		if c == carro {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func main() {
	rand.Seed(time.Now().UnixNano())

	wg.Add(numVehiculos)

	// Inicializar el estacionamiento
	for i := 0; i < capacidadEstacionamiento; i++ {
		cajones[i] = nil
	}

	// Iniciar vehículos
	for i := 0; i < numVehiculos; i++ {
		go llegadaVehiculo(i)
	}

	wg.Wait()

	fmt.Println("Simulación finalizada.")
}
