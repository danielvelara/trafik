package trafik

import (
	"math/rand"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Car struct {
	game        *Game
	direction   int
	destination int
	turn        int
	speed       float64
	distance    float64
	xPos        float64
	yPos        float64
	run         bool
	turned      bool
	light       bool
	pass        bool
	img         ebiten.Image
	semaphore   *Semaphore
	orientation []string
	vuelta      []string
}

var wg sync.WaitGroup

func CarInit(g *Game, speed float64, dir int, destination int, s *Semaphore, i int) *Car {

	rand.Seed(time.Now().UnixNano())
	min := 3
	max := 5

	car := Car{
		game:      g,
		speed:     float64(rand.Intn(max-min) + min),
		direction: dir,
		distance:  0,
		xPos:      50,
		yPos:      325,
		run:       true,
		semaphore: s,
		turned:    false,
		turn:      destination,
		pass:      false,
	}

	car.orientation = []string{"West", "South", "East", "North"}
	car.vuelta = []string{"", "Derecha", "Frente", "Izquierda"}
	car.destination = ((car.semaphore.position) + car.turn) % 4
	car.light = car.semaphore.state

	switch dir := car.direction; dir { // Load Image
	case 0: // West -> East
		img, _, _ := ebitenutil.NewImageFromFile("assets/carWE.png", ebiten.FilterDefault)
		car.img = *img
		car.xPos = -240 // giro 250
		car.yPos = 480
	case 1: // South -> North
		img, _, _ := ebitenutil.NewImageFromFile("assets/carSN.png", ebiten.FilterDefault)
		car.img = *img
		car.xPos = 480
		car.yPos = 1110 // giro 250
	case 2: // East -> West
		img, _, _ := ebitenutil.NewImageFromFile("assets/carEW.png", ebiten.FilterDefault)
		car.img = *img
		car.xPos = 1110 // giro 325
		car.yPos = 400
	case 3: // North -> South
		img, _, _ := ebitenutil.NewImageFromFile("assets/carNS.png", ebiten.FilterDefault)
		car.img = *img
		car.xPos = 400
		car.yPos = -230 // giro 325
	}
	go car.checkSemaphore()
	go car.matchSpeed()

	return &car
}

func (c *Car) checkSemaphore() {
	for true {
		time.Sleep(time.Duration(50) * time.Millisecond)
		c.light = c.semaphore.state
		if !c.light {
			if !c.run {
				c.carStart()
			}
		} else {
			if c.run {
				c.carStop()
			}
		}
	}
}

func (c *Car) Update(dTime int) error {

	// Turn Around
	if c.run {
		// If car has crossed semaphore
		if (c.distance >= 530 && c.distance < 1050) && !c.pass {
			c.pass = true
			c.dequeueW()
		} else if c.distance >= 1180 { // Car is done traveling
			c.dequeue() // Delete Items
		}

		/*
			if está en rango de vuelta && destino es hacia donde se da la vuelta
		*/
		if (c.distance >= 650 && c.distance <= 690) && c.turn == 1 { // vuelta derecha
			c.direction = (c.destination + 2) % 4
			// fmt.Println(c.distance)
			imgdir := []string{"assets/carWE.png", "assets/carSN.png", "assets/carEW.png", "assets/carNS.png"}
			img, _, _ := ebitenutil.NewImageFromFile(imgdir[c.direction], ebiten.FilterDefault)
			c.img = *img
		} else if (c.distance >= 700 && c.distance <= 750) && c.turn == 3 { // vuelta izquierda
			c.direction = (c.destination + 2) % 4
			// fmt.Println(c.distance)
			imgdir := []string{"assets/carWE.png", "assets/carSN.png", "assets/carEW.png", "assets/carNS.png"}
			img, _, _ := ebitenutil.NewImageFromFile(imgdir[c.direction], ebiten.FilterDefault)
			c.img = *img
		}

		// Change car direction
		switch d := c.direction; d {
		case 0: // East
			c.xPos += c.speed
			c.distance += c.speed
		case 1: // North
			c.yPos -= c.speed
			c.distance += c.speed
		case 2: // West
			c.xPos -= c.speed
			c.distance += c.speed
		case 3: // South
			c.yPos += c.speed
			c.distance += c.speed
		}
	}
	return nil
}

func (c *Car) Draw(screen *ebiten.Image) error {

	cDo := &ebiten.DrawImageOptions{}
	cDo.GeoM.Translate(c.xPos, c.yPos)

	screen.DrawImage(&c.img, cDo)

	return nil
}

func (c *Car) matchSpeed() {
	for true {
		time.Sleep(time.Millisecond * 10)
		i := int(c.exitQueuePos())
		if i > 0 { // si hay algun auto
			if c.pass { // pasando el semaforo revisa la velocida dde otros
				cr := c.game.carQueues[c.destination][i-1] // guarda referencia al carro frente a el
				if (c.distance - cr.distance) <= 250 {
					if c.speed > cr.speed { // iguala su velocidad
						c.speed = cr.speed // si es mayor
					}
				}
			} else { // si no ha pasado se fija de la lista de espera
				i := int(c.queuePos())
				if i > 0 {
					cr := c.semaphore.carsAtLight[i-1] // guarda referencia al carro frente a el
					if (c.distance - cr.distance) <= 250 {
						if c.speed > cr.speed { // iguala su velocidad
							c.speed = cr.speed // si es mayor
						}
					}
				}
			}
		}
	}
}

func (c *Car) queuePos() int {
	for i, cr := range c.semaphore.carsAtLight {
		if c == cr {
			return i
		}
	}
	return -1
}
func (c *Car) exitQueuePos() int {
	for i, cr := range c.game.carQueues[c.destination] {
		if c == cr {
			return i
		}
	}

	return -1
}
func (c *Car) atPos() bool {
	pos := float64(c.queuePos())
	// comparar distancia recorrida, contra distancia de semaforo segun pos
	// Ej. Posición 2 tiene que estar a = distancia de semaforo - distancia de un carro * pos
	//waitPos := c.dis < (290 - 90*pos)
	dist := (500 - 90*pos)
	if c.distance < dist {
		return false
	} else if c.distance == dist {
		return true
	} else {
		return true
	}
}

func (c *Car) carStop() {
	if !c.atPos() || c.pass {
		c.run = true
	} else {
		c.run = false
	}
}

func (c *Car) carStart() {
	c.run = true
}

func (c *Car) dequeueW() {
	time.Sleep(50 * time.Millisecond)

	if len(c.semaphore.carsAtLight) > 0 {

		c.queue()
		i := c.queuePos()
		// elimina el primer elemento
		c.semaphore.carsAtLight = append(c.semaphore.carsAtLight[:i], c.semaphore.carsAtLight[i+1:]...)
	}
}

func (c *Car) queue() {
	dest := c.destination
	c.game.carQueues[dest] = append(c.game.carQueues[dest], c)
	c.game.hud.q[c.destination]++
}

func (c *Car) dequeue() {
	if len(c.game.carQueues[c.destination]) > 0 {
		i := c.queuePos()
		if i > 0 { // remote ith element
			c.game.carQueues[c.destination] = append(c.game.carQueues[c.destination][:i], c.game.carQueues[c.destination][i+1:]...)
			c.semaphore.cars = append(c.semaphore.cars[i:], c.semaphore.cars[i+1:]...)
		} else { // remove first element
			c.game.carQueues[c.destination] = c.game.carQueues[c.destination][1:]
			c.semaphore.cars = c.semaphore.cars[1:]

		}
		c.game.hud.q[c.destination]--
		c.game.hud.currentCars--

	}
}
