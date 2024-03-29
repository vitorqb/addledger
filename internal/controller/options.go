package controller

import (
	"io"

	"github.com/vitorqb/addledger/internal/dateguesser"
	"github.com/vitorqb/addledger/internal/eventbus"
	"github.com/vitorqb/addledger/internal/metaloader"
	printermod "github.com/vitorqb/addledger/internal/printer"
	"github.com/vitorqb/addledger/internal/usermessenger"
)

// Opts represents all options for an InputController
type Opts struct {
	// Where to write journal entries to.
	output io.Writer
	// The instance of IEventBus to use
	eventBus eventbus.IEventBus
	// The instance of DateGuesser to user
	dateGuesser dateguesser.IDateGuesser
	// The instance of IMetaLoader to use
	metaLoader metaloader.IMetaLoader
	// The instance of IPrinter to use
	printer printermod.IPrinter
	// The instance of ICSVStatementLoader to use
	csvStatementLoader StatementLoader
	// The instance of IUserMessenger to use
	userMessenger usermessenger.IUserMessenger
}

// Opt configures options for an InputController
type Opt interface {
	configure(opts *Opts) error
}

// OptFn is a function option used to configure a JetStream Subscribe.
type OptFn func(opts *Opts) error

// configure implements Opt for OptFn
func (opt OptFn) configure(opts *Opts) error {
	return opt(opts)
}

// WithOutput configures which output to use.
func WithOutput(output io.Writer) Opt {
	return OptFn(func(opts *Opts) error {
		opts.output = output
		return nil
	})
}

// WithEventBus configures which IEventBus to use.
func WithEventBus(eventBus eventbus.IEventBus) Opt {
	return OptFn(func(opts *Opts) error {
		opts.eventBus = eventBus
		return nil
	})
}

// WithDateGuesser configures which IDateGuesser to use.
func WithDateGuesser(dateGuesser dateguesser.IDateGuesser) Opt {
	return OptFn(func(opts *Opts) error {
		opts.dateGuesser = dateGuesser
		return nil
	})
}

// WithMetaLoader configures which IMetaLoader to use.
func WithMetaLoader(metaLoader metaloader.IMetaLoader) Opt {
	return OptFn(func(opts *Opts) error {
		opts.metaLoader = metaLoader
		return nil
	})
}

// WithPrinter configures which IMetaLoader to use.
func WithPrinter(printer printermod.IPrinter) Opt {
	return OptFn(func(opts *Opts) error {
		opts.printer = printer
		return nil
	})
}

// WithCSVStatementLoader configures which ICSVStatementLoader to use.
func WithCSVStatementLoader(csvStatementLoader StatementLoader) Opt {
	return OptFn(func(opts *Opts) error {
		opts.csvStatementLoader = csvStatementLoader
		return nil
	})
}

// WithUserMessenger configures which IUserMessenger to use.
func WithUserMessenger(userMessenger usermessenger.IUserMessenger) Opt {
	return OptFn(func(opts *Opts) error {
		opts.userMessenger = userMessenger
		return nil
	})
}
