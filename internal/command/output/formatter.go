package output

import (
	"fmt"
	"time"
)

type Format string

const (
	FormatJSON  Format = "json"
	FormatTable Format = "table"
	FormatYAML  Format = "yaml"
)

// Formatter formats data for output.
type Formatter interface {
	// Format takes data and returns formatted output as a string
	Format(data Formattable) (string, error)
}

type Formattable interface {
	Headers() []any
	Rows() [][]any
}

func NewFormatter(format Format) (Formatter, error) {
	switch format {
	case FormatJSON:
		return &JSONFormatter{Indent: "  "}, nil
	case FormatTable:
		return &TableFormatter{
			MaxColWidth: 50,
			Separator:   "   ",
		}, nil
	case FormatYAML:
		return &YAMLFormatter{Indent: 2}, nil
	default:
		return nil, fmt.Errorf("unsupported output format: %s", format)
	}
}

func FormatTime(v time.Time) string {
	return v.Local().Format("2006-01-02 15:04:05")
}
