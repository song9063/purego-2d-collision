package main

import (
	"fmt"
	"math"
)

var VECTOR_UP Vector2 = Vector2{X: 0.0, Y: -1.0}
var VECTOR_DOWN Vector2 = Vector2{X: 0.0, Y: 1.0}
var VECTOR_LEFT Vector2 = Vector2{X: -1.0, Y: 0.0}
var VECTOR_RIGHT Vector2 = Vector2{X: 1.0, Y: 0.0}
var VECTOR_ZERO Vector2 = Vector2{X: 0.0, Y: 0.0}

type Vector2 struct {
	X, Y float64
}

func (v Vector2) Add(other Vector2) Vector2 {
	return Vector2{
		X: v.X + other.X,
		Y: v.Y + other.Y,
	}
}

func (v Vector2) Sub(other Vector2) Vector2 {
	return Vector2{
		X: v.X - other.X,
		Y: v.Y - other.Y,
	}
}

func (v Vector2) Dot(other Vector2) float64 {
	return (v.X * other.X) + (v.Y * other.Y)
}

func (v Vector2) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v Vector2) Scale(scale float64) Vector2 {
	return Vector2{
		X: v.X * scale,
		Y: v.Y * scale,
	}
}

func (v Vector2) Rotate(angle float64) Vector2 {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	return Vector2{
		X: v.X*cos - v.Y*sin,
		Y: v.X*sin + v.Y*cos,
	}
}

func (v Vector2) Normalized() Vector2 {
	mag := v.Magnitude()
	if mag == 0 {
		return Vector2{0, 0}
	}
	return Vector2{X: v.X / mag, Y: v.Y / mag}
}

func Vector2FromAngle(angle float64) Vector2 {
	return Vector2{
		X: math.Cos(angle),
		Y: math.Sin(angle),
	}
}

func StringFromVector2(v Vector2) string {
	return fmt.Sprintf("X: %.2f, Y: %.2f", v.X, v.Y)
}
