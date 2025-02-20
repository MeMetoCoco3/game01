package main

import (
	"fmt"
	_ "log"
	"math"

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
	Position     rl.Rectangle
	Velocity     rl.Vector2
	Acceleration rl.Vector2
	State        state
	Animation    Animation
}

type Bullet struct {
	Position  rl.Rectangle
	Direction rl.Vector2
}

const gravity = 1.2
const jumpForce = -20
const speed = 5
const friction = 0.8

const bulletSpeed = 7

const PJ_WIDTH = 64
const PJ_HEIGHT = 64
const PJ_MAXSPEED = 5
const PJ_ACCELERATION = 0.4

func main() {
	fmt.Println("Start")
	screenWidth := int32(800)
	screenHeight := int32(600)

	rl.InitWindow(screenWidth, screenHeight, "Jurunya Quest!")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	bg := rl.LoadTexture("assets/bg.png")
	pjTexture := rl.LoadTexture("assets/player/fishy.png")
	bulletTexture := rl.LoadTexture("assets/shoot.png")
	defer rl.UnloadTexture(bg)
	defer rl.UnloadTexture(pjTexture)
	defer rl.UnloadTexture(bulletTexture)

	pj := Character{
		Position: rl.Rectangle{
			X:      float32(screenWidth / 2),
			Y:      float32(screenHeight / 2),
			Width:  PJ_WIDTH,
			Height: PJ_HEIGHT,
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

		// TODO: Cut this shit, and check just if the character is colisioning/
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
			pj.Position,
			rl.Vector2{0.0, 0.0},
			0.0, rl.White)

		for _, platform := range GS.platforms {
			rl.DrawRectangleRec(platform, rl.Red)
		}

		bulletFrame := rl.Rectangle{0, 0, float32(bulletTexture.Width), float32(bulletTexture.Height)}
		for _, bullet := range GS.bullets {
			angle := float32(math.Atan2(float64(bullet.Direction.Y), float64(bullet.Direction.X)) * (180.0 / math.Pi))
			rl.DrawTexturePro(bulletTexture,
				bulletFrame,
				bullet.Position,
				rl.Vector2{0.0, 0.0},
				angle, rl.White)
		}

		rl.EndDrawing()
	}

}

func (pj *Character) UpdatePosition() {
	pj.Animation.AnimationUpdate()

}

func (pj *Character) GetInput(gs *GameState) {
	pj.Acceleration.X = 0

	if rl.IsKeyDown(rl.KeyLeft) {
		pj.Animation.direction = LEFT
		pj.Acceleration.X -= PJ_ACCELERATION
	} else if rl.IsKeyDown(rl.KeyRight) {
		pj.Animation.direction = RIGHT
		pj.Acceleration.X += PJ_ACCELERATION
	}

	pj.Velocity.X += pj.Acceleration.X
	pj.Velocity.X = Clamp(pj.Velocity.X, -PJ_MAXSPEED, PJ_MAXSPEED)

	if !rl.IsKeyDown(rl.KeyRight) && !rl.IsKeyDown(rl.KeyLeft) {
		pj.Velocity.X *= friction

		if math.Abs(float64(pj.Velocity.X)) < 0.15 {
			pj.Velocity.X = 0
		}
	}

	pj.Position.X += pj.Velocity.X
	pj.Position.Y += pj.Velocity.Y

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
			Position:  rl.Rectangle{X: pj.Position.X + pj.Position.Width/2, Y: pj.Position.Y + pj.Position.Height/2, Width: 32, Height: 32},
			Direction: dir,
		}
		gs.bullets = append(gs.bullets, b)
	}

	if rl.IsKeyPressed(rl.KeySpace) && pj.State != IsJumping {
		pj.State = IsJumping
		pj.Velocity.Y = jumpForce
	}

	// State Managment

	if pj.State == IsJumping {

	} else {
		if math.Abs(float64(pj.Velocity.X)) > 0.1 {
			pj.State = IsRunning
		} else {
			pj.State = IsIDLE
		}

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
