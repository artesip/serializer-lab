package ui

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"serializer/point"
	"serializer/serializer"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const serMod = "ser"
const deSerMod = "deser"

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

func serializePoints(format string, points *[]point.PointInt, textSerialized *widget.Entry) {
	var bytes []byte
	var err error
	switch format {
	case "JSON":
		bytes, err = serializer.SerializeJSON(*points)
	case "XML":
		bytes, err = serializer.SerializeXML(*points)
	case "BIN":
		bytes, err = serializer.SerializeBIN(*points)
	case "SOAP":
		bytes, err = serializer.SerializeSOAP(*points)
	case "YAML":
		bytes, err = serializer.SerializeYAML(*points)
	case "CUSTOM":
		bytes, err = serializer.SerializeCUSTOM(*points)
	}

	if err != nil {
		textSerialized.SetText("Ошибка сериализации " + err.Error())
		return
	}

	textSerialized.SetText(string(bytes))
}

func deserializePoints(format string, points *[]point.PointInt, text []byte, textFieldBottom *widget.Entry) {
	var newPoints []point.PointInt
	var err error
	switch format {
	case "JSON":
		newPoints, err = serializer.DeserializeJSON(text)
	case "XML":
		newPoints, err = serializer.DeserializeXML(text)
	case "BIN":
		newPoints, err = serializer.DeserializeBIN(text)
	case "SOAP":
		newPoints, err = serializer.DeserializeSOAP(text)
	case "YAML":
		newPoints, err = serializer.DeserializeYAML(text)
	case "CUSTOM":
		newPoints, err = serializer.DeserializeCUSTOM(text)
	}

	if err != nil {
		textFieldBottom.SetText("Ошибка десериализации " + err.Error())
		return
	}

	*points = newPoints

	var pointsStr string
	for _, p := range newPoints {
		pointsStr += fmt.Sprintln(p.ToString())
	}
	textFieldBottom.SetText(pointsStr)
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

func selectFileWithExtension(extension string, myWindow fyne.Window, callback func(string, error)) {
	dialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			callback("", fmt.Errorf("ошибка при выборе файла: %w", err))
			return
		}
		if reader == nil {
			callback("", errors.New("файл не выбран"))
			return
		}

		filePath := reader.URI().Path()
		if !strings.HasSuffix(filePath, extension) {
			callback("", errors.New("выбран файл с неправильным расширением"))
			return
		}

		data, err := io.ReadAll(reader)
		if err != nil {
			callback("", fmt.Errorf("ошибка при чтении файла: %w", err))
			return
		}

		callback(string(data), nil)
	}, myWindow)
	dialog.Show()
}

func switchModeToDefault(currentMode *string,
	selectedFormat string,
	modeLabel *widget.Label,
	points *[]point.PointInt,
	textFieldHigh *widget.Entry,
	textFieldBottom *widget.Entry) {
	if *currentMode != serMod {
		*currentMode = serMod
		modeLabel.SetText("Сериализация:")

		var pointStr string
		for _, p := range *points {
			pointStr += fmt.Sprintln(p.ToString())
		}

		textFieldHigh.SetText(pointStr)

		serializePoints(selectedFormat, points, textFieldBottom)
	}
}

func createSidebarButtons(
	selectedFormat *string,
	points *[]point.PointInt,
	currentMode *string,
	modeLabel *widget.Label,
	textFieldHigh *widget.Entry,
	textFieldBottom *widget.Entry) *fyne.Container {
	formats := []string{"JSON", "XML", "BIN", "SOAP", "YAML", "CUSTOM"}
	formatButtons := container.NewVBox()
	var buttons []*widget.Button

	for _, format := range formats {
		format := format
		var btn *widget.Button
		btn = widget.NewButton(format, func() {
			*selectedFormat = format
			switchModeToDefault(currentMode, *selectedFormat, modeLabel, points, textFieldHigh, textFieldBottom)
			for _, b := range buttons {
				b.Importance = widget.MediumImportance
				b.Refresh()
			}
			btn.Importance = widget.HighImportance
			if len(*points) != 0 {
				serializePoints(*selectedFormat, points, textFieldBottom)
			}
		})
		if *selectedFormat == format {
			btn.Importance = widget.HighImportance
		}
		buttons = append(buttons, btn)
		formatButtons.Add(btn)
	}
	return formatButtons
}

func createMainScreenButtons(
	points *[]point.PointInt, currentMode *string,
	selectedFormat *string,
	textFieldHigh *widget.Entry,
	textFieldBottom *widget.Entry,
	modeLabel *widget.Label,
	window fyne.Window) (*widget.Button, *widget.Button, *widget.Button, *widget.Button) {
	btnCreate := widget.NewButton("Создать точки", func() {
		*points = createPoints(5)
		var pointStr string
		for _, p := range *points {
			pointStr += fmt.Sprintln(p.ToString())
		}
		textFieldHigh.SetText("Созданы точки:\n" + pointStr)
		serializePoints(*selectedFormat, points, textFieldBottom)
	})

	btnSort := widget.NewButton("Сортировать", func() {
		if len(*points) == 0 {
			textFieldHigh.SetText("Сначала создайте точки!")
			return
		}
		sort.Sort(point.ByMetric(*points))
		var pointStr string
		for _, p := range *points {
			pointStr += fmt.Sprintln(p.ToString())
		}
		textFieldHigh.SetText("Отсортированные точки:\n" + pointStr)
		serializePoints(*selectedFormat, points, textFieldBottom)
	})

	btnSave := widget.NewButton("Сохранить в файл", func() {
		if *currentMode == serMod {
			saveToFile(textFieldBottom.Text, *selectedFormat)
		} else {
			selectFileWithExtension(*selectedFormat, window, func(content string, err error) {
				if err != nil {
					textFieldHigh.SetText("Ошибка " + err.Error())
				}

				textFieldHigh.SetText(content)

				deserializePoints(*selectedFormat, points, []byte(content), textFieldBottom)
			})
		}

	})

	btnSwitch := widget.NewButton("Переключить", func() {
		if *currentMode == serMod {
			*currentMode = deSerMod
			modeLabel.SetText("Десериализация:")

			btnSave.SetText("Выбрать файл")

			textFieldHigh.SetText(textFieldBottom.Text)

			deserializePoints(*selectedFormat, points, []byte(textFieldBottom.Text), textFieldBottom)
		} else {
			*currentMode = serMod
			modeLabel.SetText("Сериализация:")
			btnSave.SetText("Сохранить в файл")
			var pointStr string
			for _, p := range *points {
				pointStr += fmt.Sprintln(p.ToString())
			}

			textFieldHigh.SetText(pointStr)

			serializePoints(*selectedFormat, points, textFieldBottom)
		}
	})
	return btnCreate, btnSort, btnSave, btnSwitch
}

func CreateUI() {
	myApp := app.NewWithID("1")
	myWindow := myApp.NewWindow("MySerializer")
	myWindow.Resize(fyne.NewSize(600, 500))

	textFieldBottom := widget.NewMultiLineEntry()
	textFieldBottom.SetMinRowsVisible(12)
	textFieldBottom.Wrapping = fyne.TextWrapWord

	textFieldHigh := widget.NewMultiLineEntry()
	textFieldHigh.SetMinRowsVisible(9)
	textFieldHigh.Wrapping = fyne.TextWrapWord

	var points []point.PointInt
	selectedFormat := "JSON"
	currentMode := serMod

	modeLabel := widget.NewLabel("Сериализация:")
	resultLabel := widget.NewLabel("Результат:")

	sidebarButons := createSidebarButtons(&selectedFormat, &points, &currentMode, modeLabel, textFieldHigh, textFieldBottom)

	btnCreate, btnSort, btnSave, btnSwitch := createMainScreenButtons(&points, &currentMode, &selectedFormat, textFieldHigh, textFieldBottom, modeLabel, myWindow)

	mainButtons := container.NewHBox(btnCreate, btnSort, btnSave, btnSwitch)

	mainContent := container.NewVBox(
		modeLabel,
		textFieldHigh,
		resultLabel,
		textFieldBottom,
		mainButtons,
	)

	content := container.NewHSplit(
		container.NewVBox(widget.NewLabel("Выберите формат:"), sidebarButons),
		mainContent,
	)
	content.SetOffset(0.2)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
