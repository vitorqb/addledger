package printer

import (
	"embed"
	"fmt"
	"io"
	"strings"
	"text/template"
	"time"

	"github.com/vitorqb/addledger/internal/journal"
	"github.com/vitorqb/addledger/internal/userinput"
)

//go:embed template.txt
var templates embed.FS

// IPrinter is an interface for printing transactions.
type IPrinter interface {
	// Print prints the provided transaction to the provided writer.
	Print(writer io.Writer, transaction journal.Transaction) error
}

// Printer is a default implementation of IPrinter.
type Printer struct {
	NumLineBreaksBefore int // Number of empty lines to print before.
	NumLineBreaksAfter  int // Number of empty lines to print after.
}

// TemplateData is the data that will be used to fill the template.
// It is pretty similar to journal.Transaction but prepares some extra formatting.
type TemplateData struct {
	Description string
	Date        time.Time
	Posting     []journal.Posting
	Comment     string
}

func (p *Printer) Print(writer io.Writer, transaction journal.Transaction) error {
	// Print the configured number of empty lines before
	for i := 0; i < p.NumLineBreaksBefore; i++ {
		_, err := io.WriteString(writer, "\n")
		if err != nil {
			return fmt.Errorf("failed to write: %w", err)
		}
	}

	// Prepare the data to be used in the template
	templateData := TemplateData{
		Description: transaction.Description,
		Date:        transaction.Date,
		Posting:     transaction.Posting,
		Comment:     transaction.Comment,
	}
	tagsStrList := userinput.TagsToText(transaction.Tags)
	if len(tagsStrList) > 0 {
		if templateData.Comment != "" {
			templateData.Comment += " "
		}
		templateData.Comment += strings.Join(tagsStrList, " ")
	}

	tmpl, err := template.ParseFS(templates, "template.txt")
	if err != nil {
		return err
	}

	if err := tmpl.Execute(writer, templateData); err != nil {
		return err
	}

	// Print the configured number of empty lines after
	for i := 0; i < p.NumLineBreaksAfter; i++ {
		_, err := io.WriteString(writer, "\n")
		if err != nil {
			return fmt.Errorf("failed to write: %w", err)
		}
	}

	return nil
}

// New creates a new instance of Printer that implements IPrinter.
func New(numLineBreaksBefore, numLineBreaksAfter int) IPrinter {
	return &Printer{
		NumLineBreaksBefore: numLineBreaksBefore,
		NumLineBreaksAfter:  numLineBreaksAfter,
	}
}
