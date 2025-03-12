package ui

import (
	"fmt"
	"math/rand"
	"os"
	"serializer/point"
	"serializer/serializer"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
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

func serializePoints(format string, points []point.PointInt, textSerialized *widget.Entry) {
	switch format {
	case "JSON":
		serializer.JsonSerialize(points, textSerialized)
	case "XML":
		serializer.XmlSerialize(points, textSerialized)
	case "BIN":
		serializer.BinarySerialize(points, textSerialized)
	case "SOAP":
		serializer.SoapSerialize(points, textSerialized)
	case "YAML":
		serializer.YamlSerialize(points, textSerialized)
	}
}

func saveToFile(content string, format string) {
	filename := "output." + format
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()
	file.WriteString(content)
}

func CreateUI() {
	myApp := app.New()
	myWindow := myApp.NewWindow("MySerializer")
	myWindow.Resize(fyne.NewSize(600, 500))

	textResult := widget.NewMultiLineEntry()
	textResult.SetMinRowsVisible(12)
	textResult.Wrapping = fyne.TextWrapWord

	textSerialized := widget.NewMultiLineEntry()
	textSerialized.SetMinRowsVisible(6)
	textSerialized.Wrapping = fyne.TextWrapWord

	var points []point.PointInt
	selectedFormat := "JSON"

	formats := []string{"JSON", "XML", "BIN", "SOAP", "YAML", "CUSTOM"}
	formatButtons := container.NewVBox()
	var buttons []*widget.Button

	for _, format := range formats {
		format := format
		var btn *widget.Button
		btn = widget.NewButton(format, func() {
			selectedFormat = format
			for _, b := range buttons {
				b.Importance = widget.MediumImportance
				b.Refresh()
			}
			btn.Importance = widget.HighImportance
			if len(points) != 0 {
				serializePoints(selectedFormat, points, textResult)
			}
		})
		if selectedFormat == format {
			btn.Importance = widget.HighImportance
		}
		buttons = append(buttons, btn)
		formatButtons.Add(btn)
	}

	btnCreate := widget.NewButton("Создать точки", func() {
		points = createPoints(5)
		var pointStr string
		for _, p := range points {
			pointStr += fmt.Sprintln(p.ToString())
		}
		textSerialized.SetText("Созданы точки:\n" + pointStr)
		serializePoints(selectedFormat, points, textResult)
	})

	btnSort := widget.NewButton("Сортировать", func() {
		if len(points) == 0 {
			textSerialized.SetText("Сначала создайте точки!")
			return
		}
		sort.Sort(point.ByMetric(points))
		var pointStr string
		for _, p := range points {
			pointStr += fmt.Sprintln(p.ToString())
		}
		textSerialized.SetText("Отсортированные точки:\n" + pointStr)
		serializePoints(selectedFormat, points, textResult)
	})

	btnSave := widget.NewButton("Сохранить в файл", func() {
		saveToFile(textResult.Text, selectedFormat)
	})

	btnSwitch := widget.NewButton("Перевернуть", func() {
		textSerialized.SetText("Десериализация в формате: " + selectedFormat)
	})

	mainButtons := container.NewHBox(btnCreate, btnSort, btnSave, btnSwitch)

	mainContent := container.NewVBox(
		widget.NewLabel("Сериализация:"),
		textSerialized,
		widget.NewLabel("Результат:"),
		textResult,
		mainButtons,
	)

	content := container.NewHSplit(
		container.NewVBox(widget.NewLabel("Выберите формат:"), formatButtons),
		mainContent,
	)
	content.SetOffset(0.2)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
