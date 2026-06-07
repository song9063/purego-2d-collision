package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const DT = 1.0 / 60.0
const TRUCK_SPEED = 300
const TRUCK_ROTATION = 3.0

const SCREEN_W = 640
const SCREEN_H = 480

const EPSILON = 1e-9 // 0.000000001
func FloatEquals(a, b float64) bool {
	return math.Abs(a-b) < EPSILON
}

type SpriteShape int

const (
	SpriteShapeRect SpriteShape = iota
	SpriteShapeCircle
)

type Sprite struct {
	*OBB

	shape   SpriteShape
	image   *ebiten.Image
	imageOp *ebiten.DrawImageOptions
	color   color.Color
}

func NewSprite(center Vector2, hw, hh, angle float64,
	c color.Color,
	shape SpriteShape) *Sprite {
	return &Sprite{
		OBB:   NewOBB(center, hw, hh, angle),
		color: c,
		shape: shape,
	}
}

func (d *Sprite) GetImage() *ebiten.Image {
	if d.image == nil {
		if d.shape == SpriteShapeRect {
			d.image = ebiten.NewImage(int(d.OBB.HW*2), int(d.OBB.HH*2))
			vector.FillRect(d.image, 0, 0,
				float32(d.OBB.HW*2), float32(d.OBB.HH*2),
				d.color, true)
		} else {
			diameter := int(d.OBB.HW * 2)
			d.image = ebiten.NewImage(diameter, diameter)
			vector.FillCircle(d.image,
				float32(d.OBB.HW),
				float32(d.OBB.HW),
				float32(d.OBB.HW), d.color, true)
		}
	}
	return d.image
}

func (d *Sprite) DrawImageOps() *ebiten.DrawImageOptions {
	if d.imageOp == nil {
		d.imageOp = &ebiten.DrawImageOptions{}
	}

	d.imageOp.GeoM.Reset()
	d.imageOp.GeoM.Translate(-d.OBB.HW, -d.OBB.HH)
	d.imageOp.GeoM.Rotate(d.OBB.Angle)
	d.imageOp.GeoM.Translate(d.OBB.Center.X, d.OBB.Center.Y)

	return d.imageOp
}

type Zombie struct {
	*Sprite
	Speed float64
}

type Truck struct {
	*Sprite
	Speed         float64
	RotationSpeed float64
}

type Building struct {
	*Sprite
}

type Game struct {
	truck     *Truck
	buildings []*Building
	zombies   []*Zombie

	debugStr string
}

func (g *Game) Update() error {
	if g.truck != nil {
		curSpeed := 0.0
		prevCenter := g.truck.OBB.Center
		prevAngle := g.truck.OBB.Angle
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
			curSpeed = g.truck.Speed
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
			curSpeed = -g.truck.Speed
		}

		if !FloatEquals(curSpeed, 0) {
			forward := 1.0
			if curSpeed < 0 {
				forward = -1.0
			}

			if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
				g.truck.OBB.Angle += -g.truck.RotationSpeed * DT * forward
			}
			if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
				g.truck.OBB.Angle += +g.truck.RotationSpeed * DT * forward
			}
		}

		dir := Vector2FromAngle(g.truck.OBB.Angle)
		speed := dir.Scale(curSpeed * DT)
		g.truck.Center = g.truck.Center.Add(speed)
		//g.debugStr = StringFromVector2(speed)

		isCrash := false
		for _, b := range g.buildings {
			if !isCrash && OBBCheckCollision(b.OBB, g.truck.OBB) {
				isCrash = true
				break
			}
		}
		for _, z := range g.zombies {
			if !isCrash && OBBCheckCircleCollision(g.truck.OBB, z.Center, z.HW) {
				isCrash = true
				break
			}
		}
		if isCrash {
			g.debugStr = "Crash!!!!"
			g.truck.Center = prevCenter
			g.truck.Angle = prevAngle
		} else {
			g.debugStr = ""
		}

	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, g.debugStr)

	if g.truck != nil {
		screen.DrawImage(g.truck.GetImage(), g.truck.DrawImageOps())
	}
	for _, b := range g.buildings {
		screen.DrawImage(b.GetImage(), b.DrawImageOps())
	}
	for _, z := range g.zombies {
		screen.DrawImage(z.GetImage(), z.DrawImageOps())
	}
}

func (g *Game) Layout(oW, oH int) (w, h int) {
	return SCREEN_W, SCREEN_H
}

func main() {

	truck := &Truck{
		Sprite: NewSprite(Vector2{X: SCREEN_W - 100, Y: SCREEN_H / 2},
			50, 25, -math.Pi,
			color.RGBA{0, 0, 255, 255}, SpriteShapeRect),
		Speed:         TRUCK_SPEED,
		RotationSpeed: TRUCK_ROTATION,
	}

	building := &Building{
		Sprite: NewSprite(Vector2{X: SCREEN_W / 3, Y: SCREEN_H / 2},
			100, 80, 30,
			color.RGBA{0, 255, 0, 255}, SpriteShapeRect),
	}
	buildings := []*Building{building}

	zombie := &Zombie{
		Sprite: NewSprite(Vector2{X: SCREEN_W - 100, Y: 100},
			20, 20, 0,
			color.RGBA{255, 0, 0, 255}, SpriteShapeCircle),
	}
	zombies := []*Zombie{zombie}

	myGame := &Game{
		truck:     truck,
		buildings: buildings,
		zombies:   zombies,
	}

	ebiten.SetWindowTitle("OBB")
	ebiten.SetWindowSize(SCREEN_W, SCREEN_H)

	if err := ebiten.RunGame(myGame); err != nil {
		log.Fatal(err)
	}
}
