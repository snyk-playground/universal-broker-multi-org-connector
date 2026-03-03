package output

import "github.com/gosuri/uitable"

var _ Formatter = (*TableFormatter)(nil)

// TableFormatter formats data as a table.
type TableFormatter struct {
	MaxColWidth uint
	Separator   string
}

func (f *TableFormatter) Format(data Formattable) (string, error) {
	table := uitable.New()
	table.MaxColWidth = f.MaxColWidth
	table.Separator = f.Separator

	// add header
	table.AddRow(data.Headers()...)

	// add data rows
	rows := data.Rows()
	for _, row := range rows {
		table.AddRow(row...)
	}

	return table.String(), nil
}
