package trafik

import (
	"math/rand"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Semaphore struct {
	g             *Game
	cars          []*Car
	carsAtLight   []*Car
	timerRed      time.Duration
	timerGreen    time.Duration
	xPos          float64
	yPos          float64
	position      int
	dTime         int
	x             int
	carChan       chan int
	state         bool //true - verde - false rojo
	img           ebiten.Image
	remainingCars int
}

func SemInit(g *Game, pos int, wg *sync.WaitGroup) {

	semaphore := Semaphore{
		g:             g,
		timerRed:      7,
		timerGreen:    5,
		state:         true,
		position:      pos,
		x:             0,
		remainingCars: g.numCars,
	}
	//Posicionamineto de letreros vuelta en u
	switch p := semaphore.position; p {
	case 0: // Norte-Sur
		semaphore.xPos = 330
		semaphore.yPos = 550
	case 1: // Oeste-Este
		semaphore.xPos = 550
		semaphore.yPos = 550
	case 2: //
		semaphore.xPos = 550
		semaphore.yPos = 310
	case 3:
		semaphore.xPos = 330
		semaphore.yPos = 310
	}

	img, _, _ := ebitenutil.NewImageFromFile("assets/semrojo.png", ebiten.FilterDefault)

	semaphore.img = *img
	semaphore.carChan = make(chan int)
	semaphore.cars = []*Car{}
	semaphore.carsAtLight = []*Car{}

	go semaphore.queueManager()
	g.semaphores[pos] = &semaphore
	wg.Done()
}

func (s *Semaphore) queueManager() {
	for true {
		time.Sleep(time.Duration(3) * time.Second)
		if len(s.carsAtLight) < s.remainingCars {
			s.buildCar()
		}
	}

}

func (s *Semaphore) buildCar() {

	rand.Seed(time.Now().UnixNano())

	min := 1
	max := 4
	s.x++
	s.queue(CarInit(s.g, 2, s.position, rand.Intn(max-min)+min, s, s.x))

}

func (s *Semaphore) toggleLight() {
	s.state = !s.state
}

func (s *Semaphore) Update(dTime int, t time.Duration) error {

	time.Sleep(t)

	s.dTime = (s.dTime + 1) % 20
	//fmt.Println(s.state)
	if s.state == true {
		img, _, _ := ebitenutil.NewImageFromFile("assets/semrojo.png", ebiten.FilterDefault)
		s.img = *img
	} else {
		img, _, _ := ebitenutil.NewImageFromFile("assets/semver.png", ebiten.FilterDefault)
		s.img = *img
	}
	for i := 0; i < len(s.cars); i++ {
		if err := s.cars[i].Update(s.dTime); err != nil {
			s.carChan <- s.dTime
		}
	}

	return nil
}

func (s *Semaphore) Draw(screen *ebiten.Image) error {

	for i := 0; i < len(s.cars); i++ {
		if err := s.cars[i].Draw(screen); err != nil {
			return err
		}
	}
	//Semaforos
	cDo := &ebiten.DrawImageOptions{}
	cDa := &ebiten.DrawImageOptions{}
	cDo.GeoM.Translate(s.xPos, s.yPos)
	screen.DrawImage(&s.img, cDo)
	switch p := s.position; p {
	case 0:
		cDa.GeoM.Translate(s.xPos+20, s.yPos)
	case 1:
		cDa.GeoM.Translate(s.xPos, s.yPos-40)
	case 2:
		cDa.GeoM.Translate(s.xPos-20, s.yPos-20)
	case 3:
		cDa.GeoM.Translate(s.xPos, s.yPos-20)
	}
	return nil
}

func (s *Semaphore) queue(c *Car) {
	s.g.hud.totalCars++
	s.g.hud.currentCars++
	s.cars = append(s.cars, c)
	s.queueW(c)
}

func (s *Semaphore) queueW(c *Car) {
	s.carsAtLight = append(s.carsAtLight, c)
}
