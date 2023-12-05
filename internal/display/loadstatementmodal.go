package display

import "github.com/rivo/tview"

//go:generate $MOCKGEN --source=loadstatementmodal.go --destination=../../mocks/display/loadstatementmodal_mock.go

type (
	LoadStatementModal struct {
		*tview.Form
		controller LoadStatementModalController
	}

	LoadStatementModalController interface {
		OnLoadStatement(csvFile string, presetFile string)
	}
)

const csvFileLabel = "CSV File"
const presetLabel = "Preset"

func NewLoadStatementModal(controller LoadStatementModalController) *LoadStatementModal {
	form := &LoadStatementModal{tview.NewForm(), controller}
	form.SetBorder(true)
	form.SetTitle("Load Statement")
	form.AddInputField(csvFileLabel, "", 0, nil, nil)
	form.AddInputField(presetLabel, "", 0, nil, nil)
	form.AddButton("Load", func() {
		csvFileField := form.GetCsvInput().GetText()
		presetField := form.GetPresetInput().GetText()
		controller.OnLoadStatement(csvFileField, presetField)
	})
	return form
}

func (l *LoadStatementModal) GetCsvInput() *tview.InputField {
	return l.GetFormItemByLabel(csvFileLabel).(*tview.InputField)
}

func (l *LoadStatementModal) GetPresetInput() *tview.InputField {
	return l.GetFormItemByLabel(presetLabel).(*tview.InputField)
}
