package main

import (
	"fmt"
	"math"
	"math/rand"
)

func (v *Vector3) AddScaledVector(other *Vector3, s float64) *Vector3 {

	v.X += other.X * s
	v.Y += other.Y * s
	v.Z += other.Z * s

	return v

}

func NewVector3(x, y, z float64) *Vector3 {

	return &Vector3{x, y, z}

}

func Vector3_Max() *Vector3 {
	return &Vector3{
		math.MaxFloat64,
		math.MaxFloat64,
		math.MaxFloat64,
	}
}
func Vector3_Min() *Vector3 {
	return &Vector3{
		-math.MaxFloat64,
		-math.MaxFloat64,
		-math.MaxFloat64,
	}
}

func (v *Vector3) String() string {
	return fmt.Sprintf("&Vector3{X: %.4f, Y: %.4f, Z: %.4f}", v.X, v.Y, v.Z)
}

func (v *Vector3) MultiplyScalar(scalar float64) *Vector3 {

	if !math.IsInf(scalar, 0) {

		v.X *= scalar
		v.Y *= scalar
		v.Z *= scalar

	} else {

		v.X = 0
		v.Y = 0
		v.Z = 0

	}

	return v

}
func (v *Vector3) Copy(other *Vector3) *Vector3 {

	v.X = other.X
	v.Y = other.Y
	v.Z = other.Z

	return v

}
func randomPointCircle(r float64) *Vector3 {
	radius := r * math.Sqrt(rand.Float64())
	theta := math.Pi * 2 * rand.Float64()
	return &Vector3{
		radius * math.Cos(theta),
		20,
		radius * math.Sin(theta),
	}
}

func RemoveIndex(s []*Kaizomorph, index int) []*Kaizomorph {
	ret := make([]*Kaizomorph, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}
func ReplaceIndex(s []*Kaizomorph, index int) {
	s[index] = nil
}

func randomPointSquare(w, h, x, y float64) *Vector3 {

	return &Vector3{
		rand.Float64()*w - x,
		24,
		rand.Float64()*h - y,
	}
}
func weightedRandomSample(weights map[string]float64) string {
	totalWeight := 0.0
	for _, weight := range weights {
		totalWeight += weight
	}

	randomValue := rand.Float64() * totalWeight
	currentWeight := 0.0

	for id, weight := range weights {
		currentWeight += weight
		if randomValue <= currentWeight {
			return id
		}
	}

	// This should not happen if the totalWeight is correct,
	// but return an invalid ID (you may want to handle this case accordingly)
	return ""
}

type AABB struct {
	Min, Max Vector3
}

func isPointInAABB(p *Vector3, aabb AABB) bool {
	return p.X >= aabb.Min.X && p.X <= aabb.Max.X &&
		p.Y >= aabb.Min.Y && p.Y <= aabb.Max.Y &&
		p.Z >= aabb.Min.Z && p.Z <= aabb.Max.Z
}

type Point struct {
	x int
	y int
}
type Pixel struct {
	r uint32
	g uint32
	b uint32
	a uint32
}

/*
// Lerping function with a generic lerp method

	func (v *Vector3) Lerp(other *Vector3, alpha float64, lerpMethod func(float64) float64) *Vector3 {
		// Applying the lerping method to alpha
		alpha = lerpMethod(alpha)

		v.X += (other.X - v.X) * alpha
		v.Y += (other.Y - v.Y) * alpha
		v.Z += (other.Z - v.Z) * alpha

		return v
	}

// Example of a lerping method (InOutSine)

	func inOutSine(alpha float64) float64 {
		return 0.5 - 0.5*math.Cos(math.Pi*alpha)
	}
*/
func (v *Vector3) Lerp(other *Vector3, alpha float64) *Vector3 {

	v.X += (other.X - v.X) * alpha
	v.Y += (other.Y - v.Y) * alpha
	v.Z += (other.Z - v.Z) * alpha

	return v

}

func (v *Vector3) LerpVectors(v1, v2 *Vector3, alpha float64) *Vector3 {

	return v.SubVectors(v2, v1).MultiplyScalar(alpha).Add(v1)

}
func (v *Vector3) SubVectors(a, b *Vector3) *Vector3 {

	v.X = a.X - b.X
	v.Y = a.Y - b.Y
	v.Z = a.Z - b.Z

	return v

}
func (v *Vector3) Add(other *Vector3) *Vector3 {

	v.X += other.X
	v.Y += other.Y
	v.Z += other.Z

	return v

}
