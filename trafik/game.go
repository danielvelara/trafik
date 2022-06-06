package trafik

import (
	"math/rand"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

/* Orientation
0 - WEST
1 - SOUTH
2 - EAST
3 - NORTH
*/

type Game struct {
	semaphores       []*Semaphore
	carQueues        [][]*Car
	hud              *Hud
	isActive         bool
	numCars          int
	carChan          chan int
	dTime            int
	bg               ebiten.Image
	currentSemaphore int
}

func NewGame(ncars int) Game {
	game := Game{isActive: true, numCars: ncars, dTime: 0, currentSemaphore: 0}
	img, _, _ := ebitenutil.NewImageFromFile("assets/bg.png", ebiten.FilterDefault)

	var wg sync.WaitGroup

	game.bg = *img
	game.hud = CreateHud(&game)
	game.carChan = make(chan int)
	game.semaphores = make([]*Semaphore, 4)

	game.carQueues = [][]*Car{{}, {}, {}, {}}

	rand.Seed(time.Now().Unix())

	wg.Add(4) // Start 4 semaphores
	for i := 0; i < 4; i++ {
		go SemInit(&game, i, &wg)
	}
	wg.Wait() // Wait for creation

	game.semaphores[0].toggleLight()
	go game.handleLights()

	return game
}

func (g *Game) handleLights() {
	for true {
		time.Sleep(10 * time.Second)

		g.semaphores[g.currentSemaphore].toggleLight()
		g.currentSemaphore = (g.currentSemaphore + 1) % 4
		time.Sleep(1250 * time.Millisecond)
		g.semaphores[g.currentSemaphore].toggleLight()

	}
}

func (g *Game) Update() error {
	if g.isActive {
		g.dTime = (g.dTime + 1) % 20
		for i := 0; i < 4; i++ {
			if err := g.semaphores[i].Update(g.dTime, 2); err != nil {
				g.carChan <- g.dTime
			}
		}
	}
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) error {
	cDo := &ebiten.DrawImageOptions{}
	cDo.GeoM.Translate(0, 0)
	screen.DrawImage(&g.bg, cDo)
	for i := 0; i < 4; i++ {
		if err := g.semaphores[i].Draw(screen); err != nil {
			return err
		}
	}

	if err := g.hud.Draw(screen); err != nil {
		return err
	}

	return nil
}
