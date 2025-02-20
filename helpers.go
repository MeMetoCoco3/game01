package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

func GetVectorLength(v rl.Vector2) float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))

}

func Clamp(val, min, max float32) float32 {
	return float32(math.Max(float64(min), math.Min(float64(val), float64(max))))
}
