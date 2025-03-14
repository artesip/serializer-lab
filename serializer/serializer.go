package serializer

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"serializer/myserializer"
	"serializer/point"

	"gopkg.in/yaml.v2"
)

func wrapPoints(points []point.PointInt) []point.WrappedPoint {
	var wrapped []point.WrappedPoint
	for _, p := range points {
		wp := point.WrappedPoint{Type: p.GetType()}
		switch v := p.(type) {
		case point.Point2D:
			wp.Data = point.PointData{Point2D: &v}
		case point.Point3D:
			wp.Data = point.PointData{Point3D: &v}
		}
		wrapped = append(wrapped, wp)
	}
	return wrapped
}

func unwrapPoints(rawPoints []point.WrappedPoint) []point.PointInt {
	var points []point.PointInt
	for _, rp := range rawPoints {
		if rp.Data.Point2D != nil {
			points = append(points, *rp.Data.Point2D)
		} else if rp.Data.Point3D != nil {
			points = append(points, *rp.Data.Point3D)
		}
	}
	return points
}

// --- СЕРИАЛИЗАТОРЫ ---

func SerializeJSON(points []point.PointInt) ([]byte, error) {
	return json.MarshalIndent(wrapPoints(points), "", "  ")
}

func SerializeXML(points []point.PointInt) ([]byte, error) {
	rawPoints := wrapPoints(points)
	xmlWrapper := point.XMLPoints{Points: rawPoints}
	return xml.MarshalIndent(xmlWrapper, "", "  ")
}

func SerializeSOAP(points []point.PointInt) ([]byte, error) {
	envelope := point.SOAPEnvelope{
		XMLNS: "http://schemas.xmlsoap.org/soap/envelope/",
		Body:  point.SOAPBody{Points: point.XMLPoints{Points: wrapPoints(points)}},
	}
	return xml.MarshalIndent(envelope, "", "  ")
}

func SerializeBIN(points []point.PointInt) ([]byte, error) {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	err := enc.Encode(wrapPoints(points))
	return buf.Bytes(), err
}

func SerializeYAML(points []point.PointInt) ([]byte, error) {
	return yaml.Marshal(wrapPoints(points))
}

func SerializeCUSTOM(points []point.PointInt) ([]byte, error) {
	return myserializer.Marshal(wrapPoints(points))
}

// --- ДЕСЕРИАЛИЗАТОРЫ ---

func DeserializeJSON(data []byte) ([]point.PointInt, error) {
	var rawPoints []point.WrappedPoint
	err := json.Unmarshal(data, &rawPoints)
	return unwrapPoints(rawPoints), err
}

func DeserializeXML(data []byte) ([]point.PointInt, error) {
	var xmlWrapper point.XMLPoints
	err := xml.Unmarshal(data, &xmlWrapper)
	return unwrapPoints(xmlWrapper.Points), err
}

func DeserializeSOAP(data []byte) ([]point.PointInt, error) {
	var envelope point.SOAPEnvelope
	err := xml.Unmarshal(data, &envelope)
	return unwrapPoints(envelope.Body.Points.Points), err
}

func DeserializeBIN(data []byte) ([]point.PointInt, error) {
	var rawPoints []point.WrappedPoint
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&rawPoints)
	return unwrapPoints(rawPoints), err
}

func DeserializeYAML(data []byte) ([]point.PointInt, error) {
	var rawPoints []point.WrappedPoint
	err := yaml.Unmarshal(data, &rawPoints)
	return unwrapPoints(rawPoints), err
}

func DeserializeCUSTOM(data []byte) ([]point.PointInt, error) {
	var rawPoints []point.WrappedPoint
	err := myserializer.Unmarshal(data, &rawPoints)
	return unwrapPoints(rawPoints), err
}
