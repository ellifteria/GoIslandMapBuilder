package main

import (
	"image"
	"log"

	"github.com/ellifteria/GoIslandMapBuilder/perlin"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screen_width  = 1000
	screen_height = 1000
)

type Game struct {
	screen_image *image.RGBA
	pixel_array  []int
}

func (g *Game) Update() error {
	length := screen_width * screen_height
	for i := 0; i < length; i++ {
		g.screen_image.Pix[4*i] = uint8(g.pixel_array[4*i+0])
		g.screen_image.Pix[4*i+1] = uint8(g.pixel_array[4*i+1])
		g.screen_image.Pix[4*i+2] = uint8(g.pixel_array[4*i+2])
		g.screen_image.Pix[4*i+3] = uint8(g.pixel_array[4*i+3])
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.WritePixels(g.screen_image.Pix)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screen_width, screen_height
}

func height_to_color(height float64) [4]int {
	switch {
	case height < -0.5:
		return [4]int{1, 1, 122, 255}
	case height < 0:
		return [4]int{3, 138, 255, 255}
	case height < 0.25:
		return [4]int{243, 225, 107, 255}
	case height < 0.5:
		return [4]int{22, 160, 133, 255}
	case height < 0.75:
		return [4]int{108, 122, 137, 255}
	default:
		return [4]int{255, 255, 255, 255}
	}
}

func main() {
	var height_array [screen_width][screen_height]float64

	permutation_array := perlin.GeneratePermutation()

	for x := 0; x < screen_width; x++ {
		for y := 0; y < screen_height; y++ {
			height_array[x][y] = perlin.FractalBrownianMotion(
				float64(x),
				float64(y),
				2,
				permutation_array,
			)
		}
	}

	var flattened_heigh_array [screen_width * screen_height * 4]int

	for x := 0; x < screen_width; x++ {
		for y := 0; y < screen_height; y++ {
			color := height_to_color(height_array[x][y])
			flattened_heigh_array[(x+y*screen_width)*4+0] = color[0]
			flattened_heigh_array[(x+y*screen_width)*4+1] = color[1]
			flattened_heigh_array[(x+y*screen_width)*4+2] = color[2]
			flattened_heigh_array[(x+y*screen_width)*4+3] = color[3]
		}
	}

	ebiten.SetWindowSize(screen_width, screen_height)
	ebiten.SetWindowTitle("Go Island Map Builder")

	g := &Game{
		screen_image: image.NewRGBA(image.Rect(0, 0, screen_width, screen_height)),
		pixel_array:  flattened_heigh_array[:],
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
