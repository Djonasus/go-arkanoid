package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	width  = 640
	height = 480
)

var blockPull []*Block

var (
	circle = Circle{40, height / 2.5, 10, 5, -5}
	player = Player{Block: Block{10, height - 60, 100, 20}, speed: 10}
	paused = false
)

type Game struct{}

type Circle struct {
	posX, posY, radius, vectorX, vectorY float64
}

type Block struct {
	posX, posY, width, height float64
}

type Player struct {
	Block
	speed float64
}

// USEFUL UTILS
func Clamp(f, low, high float64) float64 {
	if f < low {
		return low
	}
	if f > high {
		return high
	}
	return f
}

func removeAndCompress(key int, sl []*Block) {
	//var loc []*Block

	sl[key] = nil

	//loc = sl[:key]
	//loc = append(loc, sl[key+1:]...)
	//return loc
}

// LOGIC LOOP
func (g *Game) Update() error {
	if !paused {
		//Sphere move
		if circle.posX-circle.radius <= 0 || circle.posX+circle.radius >= width {
			circle.vectorX = -circle.vectorX
		}
		if circle.posY-circle.radius <= 0 || circle.posY+circle.radius >= height {
			circle.vectorY = -circle.vectorY
		}

		if circle.posX-circle.radius < player.posX+player.width && circle.posX+circle.radius > player.posX && circle.posY-circle.radius < player.posY+player.height && circle.posY+circle.radius > player.posY {
			//circle.vectorX = -circle.vectorX
			circle.vectorY = -circle.vectorY
		}

		for i, v := range blockPull {
			if v != nil {
				if circle.posX-circle.radius < v.posX+v.width && circle.posX+circle.radius > v.posX && circle.posY-circle.radius < v.posY+player.height && circle.posY+circle.radius > v.posY {
					circle.vectorY = -circle.vectorY
					removeAndCompress(i, blockPull)
					continue
				}
			}
		}
		if circle.posY+circle.radius >= height-5 {
			paused = true
		}

		circle.posX += circle.vectorX
		circle.posY += circle.vectorY

		//player move
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
			player.posX += player.speed
		} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
			player.posX -= player.speed
		}
		player.posX = Clamp(player.posX, 0, width-player.width)
	}
	return nil
}

// RENDER LOOP
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	//ebitenutil.DebugPrint(screen, strconv.Itoa(score))

	for _, v := range blockPull {
		if v != nil {
			//print(v.posX)
			ebitenutil.DrawRect(screen, v.posX, v.posY, v.width, v.height, color.White)
		}
	}

	ebitenutil.DrawRect(screen, player.posX, player.posY, player.width, player.height, color.White)

	ebitenutil.DrawCircle(screen, circle.posX, circle.posY, circle.radius, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
}

func main() {

	blockPull = make([]*Block, 20*5)

	for y := 0; y <= 5; y++ {
		for x := 0; x <= 20; x++ {
			blockPull = append(blockPull, &Block{float64(x) * 35, float64(y) * 10, 35, 10})
		}
	}

	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Arcanoid")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
