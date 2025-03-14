package point

import (
	"encoding/xml"
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

type GetType interface {
	GetType() string
}

type PointInt interface {
	Metric
	ToString
	GetType
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

func (p Point2D) GetType() string {
	return "2D"
}

func (p Point3D) GetType() string {
	return "3D"
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

type SOAPEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	XMLNS   string   `xml:"xmlns:soap,attr"`
	Body    SOAPBody `xml:"Body"`
}

type SOAPBody struct {
	XMLName xml.Name  `xml:"Body"`
	Points  XMLPoints `xml:"Points"`
}

type XMLPoints struct {
	XMLName xml.Name       `xml:"Points"`
	Points  []WrappedPoint `xml:"Point"`
}

type WrappedPoint struct {
	Type string    `json:"type" xml:"type" yaml:"type"`
	Data PointData `json:"data" xml:"data" yaml:"data"`
}

type PointData struct {
	Point2D *Point2D `json:"Point2D,omitempty" xml:"Point2D,omitempty" yaml:"Point2D,omitempty" `
	Point3D *Point3D `json:"Point3D,omitempty" xml:"Point3D,omitempty" yaml:"Point3D,omitempty" `
}
