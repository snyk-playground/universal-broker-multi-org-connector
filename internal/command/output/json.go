package output

import (
	"encoding/json"
)

var _ Formatter = (*JSONFormatter)(nil)

type JSONFormatter struct {
	Indent string
}

func (f *JSONFormatter) Format(data Formattable) (string, error) {
	bytes, err := json.MarshalIndent(data, "", f.Indent)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
