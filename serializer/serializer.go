package serializer

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"serializer/point"

	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/yaml.v2"
)

func JsonSerialize(dialog *dialog.CustomDialog, points []point.PointInt, textOutput *widget.Entry) {
	jsonData, err := json.MarshalIndent(points, "", "  ")

	if err != nil {
		textOutput.SetText("Ошибка сериализации JSON: " + err.Error())
		return
	}

	textOutput.SetText(string(jsonData))
	dialog.Hide()
}

func XmlSerialize(dialog *dialog.CustomDialog, points []point.PointInt, textOutput *widget.Entry) {
	xmlData, err := xml.MarshalIndent(points, "", "  ")

	if err != nil {
		textOutput.SetText("Ошибка сериализации XML: " + err.Error())
		return
	}

	textOutput.SetText(string(xmlData))
	dialog.Hide()
}

func BinarySerialize(dialog *dialog.CustomDialog, points []point.PointInt, textOutput *widget.Entry) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	gob.Register(point.Point2D{})
	gob.Register(point.Point3D{})

	err := encoder.Encode(points)

	if err != nil {
		textOutput.SetText("Ошибка сериализации BIN: " + err.Error())
		return
	}

	textOutput.SetText(fmt.Sprintf("%x", buffer.Bytes()))
	dialog.Hide()
}

func YamlSerialize(dialog *dialog.CustomDialog, points []point.PointInt, textOutput *widget.Entry) {
	yamlData, err := yaml.Marshal(points)
	if err != nil {
		textOutput.SetText("Ошибка сериализации YAML: " + err.Error())
		return
	}

	textOutput.SetText(string(yamlData))
	dialog.Hide()
}

type Envelope struct {
	XMLName    xml.Name `xml:"soapenv:Envelope"`
	XMLNS      string   `xml:"xmlns:soapenv,attr"`
	XMLNSPoint string   `xml:"xmlns:point,attr"`
	Body       Body     `xml:"soapenv:Body"`
}

type Body struct {
	XMLName   xml.Name  `xml:"soapenv:Body"`
	GetPoints GetPoints `xml:"point:GetPoints"`
}

type GetPoints struct {
	XMLName xml.Name         `xml:"point:GetPoints"`
	Points  []point.PointInt `xml:"point:Point"`
}

func SoapSerialize(dialog *dialog.CustomDialog, points []point.PointInt, textOutput *widget.Entry) {
	envelope := Envelope{
		XMLNS:      "http://schemas.xmlsoap.org/soap/envelope/",
		XMLNSPoint: "http://example.com/point",
		Body: Body{
			GetPoints: GetPoints{
				Points: points,
			},
		},
	}

	xmlData, err := xml.MarshalIndent(envelope, "", "  ")
	if err != nil {
		textOutput.SetText("Ошибка сериализации SOAP: " + err.Error())
		return
	}

	textOutput.SetText(string(xmlData))
	dialog.Hide()
}
