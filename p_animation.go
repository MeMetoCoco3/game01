package main

import rl "github.com/gen2brain/raylib-go/raylib"

type AnimationType int
type AnimationDirection int

const (
	REPEATING AnimationType = iota
	ONESHOT
)
const (
	LEFT  AnimationDirection = -1
	RIGHT AnimationDirection = 1
)

type Animation struct {
	first         int
	last          int
	current       int
	direction     AnimationDirection
	aType         AnimationType
	speed         float32
	duration_left float32
}

func (a *Animation) AnimationUpdate() {
	dt := rl.GetFrameTime()
	a.duration_left -= dt
	if a.duration_left <= 0 {
		a.duration_left = a.speed
		a.current++

		if a.current > a.last {
			switch a.aType {
			case REPEATING:
				a.current = a.first
				break
			case ONESHOT:
				a.current = a.last
				break
			}

		}
	}
}

func (a *Animation) AnimationFrame(numFramesPerRow, sizeTile, xPad, yPad, xOffset, yOffset int) rl.Rectangle {
	x := int((a.current % numFramesPerRow) * (sizeTile + xPad))
	y := int((a.current / numFramesPerRow) * (sizeTile + yPad))
	return rl.Rectangle{
		X:      float32(x + xOffset),
		Y:      float32(y + yOffset),
		Width:  float32(sizeTile),
		Height: float32(sizeTile),
	}
}

func main2() {

	rl.InitWindow(600, 400, "Prueba animacion")

	pjTexture := rl.LoadTexture("assets/player/fishy.png")
	defer rl.UnloadTexture(pjTexture)

	anim := Animation{
		first:         0,
		last:          13,
		current:       0,
		speed:         0.1,
		duration_left: 0.1,
		aType:         ONESHOT,
	}

	for !rl.WindowShouldClose() {

		anim.AnimationUpdate()
		rl.BeginDrawing()
		rl.ClearBackground(rl.SkyBlue)
		rl.DrawTexturePro(
			pjTexture,
			anim.AnimationFrame(14, 64, 32, 0, 16, 8),
			rl.Rectangle{10.0, 10.0, 128.0, 128.0},
			rl.Vector2{0.0, 0.0}, 0.0, rl.White)
		rl.EndDrawing()
	}
}
