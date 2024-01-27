package context

import "github.com/rivo/tview"

//go:generate $MOCKGEN --source=app.go --destination=../../../mocks/display/context/app_mock.go

// TviewApp is an abstraction of tview.Application with only the methods we need.
type TviewApp interface {
	QueueUpdateDraw(func()) *tview.Application
}
