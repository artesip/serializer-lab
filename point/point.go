package point

import (
	"fmt"
	"math"
	"math/rand"
)

type Metric interface {
	Metric() float64
}

type ToString interface {
	ToString() string
}

type PointInt interface {
	Metric
	ToString
}

type Point2D struct {
	X float64 `json:"x" yaml:"x"`
	Y float64 `json:"y" yaml:"y"`
}

type Point3D struct {
	Point2D `yaml:"Point,inline"`
	Z       float64 `json:"z" yaml:"z"`
}

func RandPoint2DConstructor() Point2D {
	return Point2D{
		X: rand.Float64() * 10,
		Y: rand.Float64() * 10,
	}
}

func RandPoint3DConstructor() Point3D {
	return Point3D{
		Point2D: Point2D{
			X: rand.Float64() * 10,
			Y: rand.Float64() * 10,
		},
		Z: rand.Float64() * 10,
	}
}

func (p Point2D) ToString() string {
	return "Point{" + fmt.Sprintf("%.2f", p.X) + ", " + fmt.Sprintf("%.2f", p.Y) + "}"
}

func (p Point3D) ToString() string {
	return "Point{" + fmt.Sprintf("%.2f", p.X) + ", " + fmt.Sprintf("%.2f", p.Y) + ", " + fmt.Sprintf("%.2f", p.Z) + "}"
}

func (p Point2D) Metric() float64 {
	return math.Sqrt(math.Pow(p.X, 2) + math.Pow(p.Y, 2))
}

func (p Point3D) Metric() float64 {
	return math.Sqrt(math.Pow(p.X, 2) + math.Pow(p.Y, 2) + math.Pow(p.Z, 2))
}

type ByMetric []PointInt

func (e ByMetric) Len() int {
	return len(e)
}

func (a ByMetric) Less(i, j int) bool {
	return a[i].Metric() < a[j].Metric()
}
func (a ByMetric) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
