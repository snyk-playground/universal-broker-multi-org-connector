package command

import (
	"context"
	"os"

	"github.com/yarlson/pin"
)

var defaultFrames = []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}

func newSpinner(ctx context.Context, message string) *pin.Pin {
	p := pin.New(message,
		pin.WithSpinnerColor(pin.ColorYellow),
		pin.WithSpinnerFrames(defaultFrames),
		pin.WithWriter(os.Stderr),
	)
	p.Start(ctx)

	return p
}
