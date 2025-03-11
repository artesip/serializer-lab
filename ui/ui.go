package ui

import (
	"fmt"
	"math/rand"
	"serializer/point"
	"serializer/serializer"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func createPoints(len int) []point.PointInt {
	if len <= 0 {
		return []point.PointInt{}
	}

	points := []point.PointInt{}
	for i := 1; i <= len; i++ {
		randNum := rand.Intn(2)

		if randNum%2 == 0 {
			points = append(points, point.RandPoint2DConstructor())
		} else {
			points = append(points, point.RandPoint3DConstructor())
		}
	}

	return points
}

func createDialog(myWindow fyne.Window, points []point.PointInt, textOutput *widget.Entry) *dialog.CustomDialog {
	buttonsContainer := container.NewVBox()
	myDialog := dialog.NewCustom("Диалог с 4 кнопками", "Закрыть", buttonsContainer, myWindow)

	button1 := widget.NewButton("JSON", func() {
		serializer.JsonSerialize(myDialog, points, textOutput)
	})

	button2 := widget.NewButton("XML", func() {
		serializer.XmlSerialize(myDialog, points, textOutput)
	})

	button3 := widget.NewButton("BIN", func() {
		serializer.BinarySerialize(myDialog, points, textOutput)
	})

	button4 := widget.NewButton("SOAP", func() {
		serializer.SoapSerialize(myDialog, points, textOutput)
	})

	button5 := widget.NewButton("YAML", func() {
		serializer.YamlSerialize(myDialog, points, textOutput)
	})

	button6 := widget.NewButton("CUSTOM", func() {})

	buttonsContainer.Add(button1)
	buttonsContainer.Add(button2)
	buttonsContainer.Add(button3)
	buttonsContainer.Add(button4)
	buttonsContainer.Add(button5)
	buttonsContainer.Add(button6)

	return myDialog
}

func CreateUI() {
	myApp := app.New()
	myWindow := myApp.NewWindow("MySerializer")

	myWindow.Resize(fyne.NewSize(400, 500))

	textOutput := widget.NewMultiLineEntry()
	textOutput.Wrapping = fyne.TextWrapWord
	textOutput.SetMinRowsVisible(10)

	scrollContainer := container.NewVScroll(textOutput)
	scrollContainer.SetMinSize(fyne.NewSize(400, 200))
	textOutput.Wrapping = fyne.TextWrapWord

	var points []point.PointInt

	button1 := widget.NewButton("Create", func() {
		points = createPoints(5)

		var pointStr string
		for _, p := range points {
			pointStr += fmt.Sprintln(p.ToString())
		}

		textOutput.SetText("Созданы точки\n" + pointStr)
	})

	button2 := widget.NewButton("Sort", func() {

		if len(points) == 0 {
			textOutput.SetText("Сначала создайте точки для сортировки!\n")
			return
		}

		sort.Sort(point.ByMetric(points))

		var pointStr string
		for _, p := range points {
			pointStr += fmt.Sprintln(p.ToString())
		}

		textOutput.SetText("Отсортированные точки\n" + pointStr)
	})

	button3 := widget.NewButton("Serialize", func() {
		dialog := createDialog(myWindow, points, textOutput)
		dialog.Show()
	})

	button4 := widget.NewButton("Deserialize", func() {
		textOutput.SetText("Нажата Кнопка 4\nСтрока 2\nСтрока 3\nСтрока 4\nСтрока 5\nСтрока 6")
	})

	buttons := container.NewHBox(
		button1,
		button2,
		button3,
		button4,
	)

	content := container.NewVBox(
		textOutput,
		buttons,
	)

	myWindow.SetContent(content)

	myWindow.ShowAndRun()
}
