package printer

import (
	"embed"
	"io"
	"text/template"

	"github.com/vitorqb/addledger/internal/journal"
)

//go:embed template.txt
var templates embed.FS

// IPrinter is an interface for printing transactions.
type IPrinter interface {
	// Print prints the provided transaction to the provided writer.
	Print(writer io.Writer, transaction journal.Transaction) error
}

// Printer is a default implementation of IPrinter.
type Printer struct{}

func (p *Printer) Print(writer io.Writer, transaction journal.Transaction) error {
	tmpl, err := template.ParseFS(templates, "template.txt")
	if err != nil {
		return err
	}

	if err := tmpl.Execute(writer, transaction); err != nil {
		return err
	}

	return nil
}

// Ensure Printer implements IPrinter interface
var _ IPrinter = (*Printer)(nil)
