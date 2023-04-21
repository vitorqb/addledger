package controller

import "io"

// Opts represents all options for an InputController
type Opts struct {
	// Where to write journal entries to.
	output io.Writer
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
