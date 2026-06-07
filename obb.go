package main

import "math"

type OBB struct {
	Center Vector2
	HW     float64 // Half Width
	HH     float64 // Half Height
	Angle  float64 // Radian
}

func NewOBB(center Vector2, hw, hh, angle float64) *OBB {
	return &OBB{
		Center: center,
		HW:     hw, HH: hh,
		Angle: angle,
	}
}

// P0: RightBottom, P1: LeftBottom
// P2: LeftTop, P3: RightTop
func (o *OBB) GetVertices() [4]Vector2 {
	localVertices := [4]Vector2{
		{X: o.HW, Y: o.HH},
		{X: -o.HW, Y: o.HH},
		{X: -o.HW, Y: -o.HH},
		{X: o.HW, Y: -o.HH},
	}

	var worldVertices [4]Vector2
	for i, localPt := range localVertices {
		rotated := localPt.Rotate(o.Angle)
		worldVertices[i] = o.Center.Add(rotated)
	}

	return worldVertices
}

// Unit Vectors
// axis0: Horizontal(P0-P1)
// axis1: Vertical(P2-P1)
func (o *OBB) GetAxes() [2]Vector2 {
	vertices := o.GetVertices()

	axis0 := vertices[0].Sub(vertices[1]).Normalized()
	axis1 := vertices[2].Sub(vertices[1]).Normalized()

	return [2]Vector2{axis0, axis1}
}

func OBBCheckCollision(a, b *OBB) bool {

	axesA := a.GetAxes()
	axesB := b.GetAxes()
	axes := [4]Vector2{
		axesA[0], axesA[1],
		axesB[0], axesB[1],
	}

	verA := a.GetVertices()
	verB := b.GetVertices()

	for _, axis := range axes {
		minA, maxA := getProjectionRange(verA, axis)
		minB, maxB := getProjectionRange(verB, axis)

		if (maxA+EPSILON) < minB || (maxB+EPSILON) < minA {
			return false
		}
	}

	return true
}

func OBBCheckCircleCollision(a *OBB,
	circleCenter Vector2, circleRadius float64) bool {
	localCircleCenter := circleCenter.Sub(a.Center).Rotate(-a.Angle)
	closestX := math.Max(-a.HW, math.Min(a.HW, localCircleCenter.X))
	closestY := math.Max(-a.HH, math.Min(a.HH, localCircleCenter.Y))
	closestPt := Vector2{X: closestX, Y: closestY}

	distVector := localCircleCenter.Sub(closestPt)
	dist := distVector.Magnitude()
	return dist < (circleRadius + EPSILON)
}

func getProjectionRange(vertices [4]Vector2, axis Vector2) (float64, float64) {
	min := math.Inf(1)
	max := math.Inf(-1)

	for _, v := range vertices {
		val := axis.Dot(v)
		if val < min {
			min = val
		}
		if val > max {
			max = val
		}
	}
	return min, max
}
