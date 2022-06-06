package trafik

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Hud struct {
	game        *Game
	currentCars int
	totalCars   int
	q           []int
	time        int
}

var (
	mplusNormalFont font.Face
	mplusBigFont    font.Face
)

func CreateHud(game *Game) *Hud {
	hud := Hud{
		game:        game,
		currentCars: 0,
		totalCars:   0,
	}

	hud.q = []int{0, 0, 0, 0}

	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		fmt.Println(err)
	}
	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		fmt.Println(err)
	}
	mplusBigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	return &hud
}

func (h *Hud) Draw(screen *ebiten.Image) error {
	text.Draw(screen, "Total cars: "+strconv.Itoa(h.totalCars), mplusNormalFont, 20, 60, color.Black)
	text.Draw(screen, "Current cars: "+strconv.Itoa(h.currentCars), mplusNormalFont, 20, 110, color.Black)

	return nil
}
