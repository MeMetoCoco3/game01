package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

func GetVectorLength(v rl.Vector2) float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))

}
