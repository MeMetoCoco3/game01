package main

import (
	"fmt"
	_ "log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type state = int

const (
	IsIDLE state = iota
	IsJumping
	IsRunning
)

type GameState struct {
	pj        Character
	platforms []rl.Rectangle
	bullets   []Bullet
}
type Character struct {
	Position  rl.Rectangle
	Velocity  rl.Vector2
	State     state
	Animation Animation
}

type Bullet struct {
	Position  rl.Rectangle
	Direction rl.Vector2
}

const gravity = 1.2
const jumpForce = -20
const speed = 5

const bulletSpeed = 7

func main() {
	fmt.Println("Start")
	screenWidth := int32(800)
	screenHeight := int32(600)

	rl.InitWindow(screenWidth, screenHeight, "Jurunya Quest!")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	bg := rl.LoadTexture("assets/bg.png")
	defer rl.UnloadTexture(bg)
	pjTexture := rl.LoadTexture("assets/player/fishy.png")
	defer rl.UnloadTexture(pjTexture)

	pj := Character{
		Position: rl.Rectangle{
			X:      float32(screenWidth / 2),
			Y:      float32(screenHeight / 2),
			Width:  64,
			Height: 96,
		},
		Velocity: rl.Vector2{X: 0, Y: 0},

		//This is idle Animation, rest of animation will change the nums for last
		Animation: Animation{
			first:         0,
			last:          13,
			current:       0,
			speed:         0.1,
			duration_left: 0.1,
			aType:         REPEATING,
			direction:     RIGHT,
		},
	}

	platforms := []rl.Rectangle{
		{
			X:      0,
			Y:      0,
			Width:  800,
			Height: 20,
		},
		{
			X:      0,
			Y:      580,
			Width:  800,
			Height: 20,
		}, {
			X:      0,
			Y:      0,
			Width:  20,
			Height: 600,
		},
		{
			X:      780,
			Y:      0,
			Width:  20,
			Height: 600,
		},
	}

	GS := GameState{
		pj:        pj,
		platforms: platforms,
		bullets:   []Bullet{},
	}

	for !rl.WindowShouldClose() {
		pj.GetInput(&GS)
		pj.Velocity.Y += gravity
		pj.UpdatePosition()

		for _, platform := range platforms {
			if rl.CheckCollisionRecs(pj.Position, platform) {
				pj.Position.Y = platform.Y - pj.Position.Height
				pj.Velocity.Y = 0
				pj.State = IsIDLE
			}

		}

		// TODO: Study other ways of doing this shit
		newBullets := []Bullet{}
		for _, bullet := range GS.bullets {
			bullet.Position.X += bullet.Direction.X
			bullet.Position.Y += bullet.Direction.Y

			collision := false

			for _, platform := range platforms {
				if rl.CheckCollisionRecs(bullet.Position, platform) {
					collision = true
					break
				}
			}
			if !collision {
				newBullets = append(newBullets, bullet)
			}
		}
		GS.bullets = newBullets

		pjFrame := pj.UpdateAnimationFrame()
		pjFrame.Width *= float32(pj.Animation.direction)
		// Drawing
		rl.BeginDrawing()

		bgRec := rl.Rectangle{0, 0, float32(bg.Width), float32(bg.Height)}
		dstRec := rl.Rectangle{0, 0, float32(screenWidth), float32(screenHeight)}
		rl.DrawTexturePro(bg, bgRec, dstRec, rl.Vector2{0, 0}, 0, rl.RayWhite)

		rl.DrawRectangleRec(pj.Position, rl.Red)

		rl.DrawTexturePro(
			pjTexture,
			pjFrame,
			rl.Rectangle{pj.Position.X - 32, pj.Position.Y - 16, 128.0, 128.0},
			rl.Vector2{0.0, 0.0},
			0.0, rl.White)

		for _, platform := range GS.platforms {
			rl.DrawRectangleRec(platform, rl.Red)
		}
		for _, bullet := range GS.bullets {
			rl.DrawRectangleRec(bullet.Position, rl.RayWhite)
		}

		rl.EndDrawing()
	}

}

func (pj *Character) UpdatePosition() {
	pj.Animation.AnimationUpdate()

	pj.Position.X += pj.Velocity.X
	pj.Position.Y += pj.Velocity.Y
}

func (pj *Character) GetInput(gs *GameState) {
	if rl.IsKeyDown(rl.KeyLeft) {
		pj.Animation.direction = LEFT
		pj.Velocity.X = -speed
	} else if rl.IsKeyDown(rl.KeyRight) {
		pj.Animation.direction = RIGHT
		pj.Velocity.X = speed

	} else {
		pj.Velocity.X = 0
	}

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {

		mousePos := rl.GetMousePosition()
		dir := rl.Vector2{
			X: (mousePos.X - pj.Position.X),
			Y: (mousePos.Y - pj.Position.Y),
		}
		length := GetVectorLength(dir)
		if length != 0 {
			dir.X = (dir.X / length) * bulletSpeed
			dir.Y = (dir.Y / length) * bulletSpeed
		}

		b := Bullet{
			Position:  rl.Rectangle{X: pj.Position.X, Y: pj.Position.Y, Width: 10, Height: 10},
			Direction: dir,
		}
		gs.bullets = append(gs.bullets, b)
	}

	if rl.IsKeyPressed(rl.KeySpace) && pj.State != IsJumping {
		pj.State = IsJumping
		pj.Velocity.Y = jumpForce
	}
}

func (pj *Character) UpdateAnimationFrame() rl.Rectangle {
	if pj.State == IsJumping {
		if pj.Animation.aType != ONESHOT {
			pj.Animation.last = 2
			pj.Animation.current = 0
			pj.Animation.aType = ONESHOT
		}
		return pj.Animation.AnimationFrame(3, 64, 32, 0, 16, 408)
	}

	if pj.Velocity.X != 0 {
		if pj.Animation.last != 5 {
			pj.Animation.aType = REPEATING
			pj.Animation.last = 5
			pj.Animation.current = 0
		}
		return pj.Animation.AnimationFrame(6, 64, 32, 0, 16, 88)
	}
	// This check exists to see if it was a change between animations
	if pj.Animation.last != 13 {
		pj.Animation.aType = REPEATING
		pj.Animation.last = 13
		pj.Animation.current = 0
	}
	return pj.Animation.AnimationFrame(14, 64, 32, 0, 16, 8)
}
